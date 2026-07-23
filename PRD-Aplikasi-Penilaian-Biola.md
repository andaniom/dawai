# Product Requirements Document (PRD)
## DAWAI — Sistem Penilaian Biola Berbasis Web, Multi-Tenant

**Versi:** 2.1  
**Status:** Draft  
**Terakhir Diperbarui:** Juli 2026  
**Changelog:**
- v2.1 — Nama produk ditetapkan: **DAWAI**. Arsitektur diubah ke decoupled: Go + Fiber sebagai backend API, Next.js sebagai frontend UI saja. Auth berubah dari Auth.js session ke JWT stateless (Go issue token) + blacklist di PostgreSQL.
- v2.0 — Arsitektur diubah dari single-tenant ke multi-tenant (Shared Database + Shared Schema). Ditambahkan role `super_admin` dan `school_admin`. Onboarding sekolah via super_admin dashboard (manual).

---

## Daftar Isi

1. [Latar Belakang & Problem Statement](#1-latar-belakang--problem-statement)
2. [Tujuan Produk](#2-tujuan-produk)
3. [Asumsi & Batasan](#3-asumsi--batasan)
4. [Model Multi-Tenant](#4-model-multi-tenant)
5. [Pengguna & Peran (RBAC)](#5-pengguna--peran-rbac)
6. [User Stories & Acceptance Criteria](#6-user-stories--acceptance-criteria)
7. [Fitur & Spesifikasi Fungsional](#7-fitur--spesifikasi-fungsional)
8. [Logika Bisnis & Aturan Penilaian](#8-logika-bisnis--aturan-penilaian)
9. [Arsitektur Teknis](#9-arsitektur-teknis)
10. [Schema Database](#10-schema-database)
11. [Tenant Isolation & Row-Level Security](#11-tenant-isolation--row-level-security)
12. [API Design](#12-api-design)
13. [Non-Functional Requirements](#13-non-functional-requirements)
14. [Out of Scope (MVP)](#14-out-of-scope-mvp)
15. [Risiko & Mitigasi](#15-risiko--mitigasi)
16. [Milestone & Prioritas](#16-milestone--prioritas)
17. [Keputusan Teknis Terbuka](#17-keputusan-teknis-terbuka)

---

## 1. Latar Belakang & Problem Statement

### Konteks

Guru biola di sekolah musik / ekstrakurikuler mengelola puluhan siswa dengan level kemampuan berbeda. Proses penilaian saat ini bergantung pada catatan manual (buku nilai, spreadsheet), yang mengakibatkan:

- **Inkonsistensi rubrik** antar sesi atau antar guru.
- **Tidak ada visibilitas progres** bagi siswa dan orang tua di antara ujian.
- **Laporan rapor memakan waktu** karena deskripsi ditulis manual tiap siswa.
- **Tidak ada rekam jejak** untuk mengidentifikasi pola kelemahan siswa dari waktu ke waktu.

Sistem ini akan dipakai di **beberapa sekolah sekaligus** dalam satu platform terpusat, sehingga membutuhkan arsitektur multi-tenant yang menjamin isolasi data antar sekolah.

### Problem Statement

> Guru biola tidak memiliki alat digital yang memadai untuk mencatat penilaian secara konsisten, melacak perkembangan siswa, dan menghasilkan laporan Kurikulum Merdeka secara otomatis — berlaku untuk banyak sekolah sekaligus dalam satu platform.

---

## 2. Tujuan Produk

| Tujuan | Metrik Keberhasilan |
|---|---|
| Mempercepat proses input nilai di kelas | Waktu input per siswa < 60 detik |
| Konsistensi rubrik antar sesi penilaian | 100% penilaian menggunakan rubrik terstandar |
| Transparansi progres ke orang tua | Orang tua dapat melihat nilai & komentar kapan saja |
| Otomasi laporan Kurikulum Merdeka | Deskripsi rapor ter-generate tanpa pengetikan manual |
| Operasi saat offline | Penilaian dapat dilakukan tanpa koneksi internet |
| Skalabilitas multi-sekolah | Satu platform bisa melayani banyak sekolah dengan isolasi data penuh |
| Visibilitas lintas sekolah | Super admin bisa monitor semua sekolah dari satu dashboard |

---

## 3. Asumsi & Batasan

### Asumsi

- Target awal: 3–10 sekolah, masing-masing dengan 1–5 guru dan 50–200 siswa.
- Setiap sekolah dikelola oleh 1 `school_admin` yang ditunjuk.
- Onboarding sekolah baru dilakukan manual oleh `super_admin` (bukan self-registration).
- Guru menggunakan HP saat di kelas; desktop/laptop untuk rekap dan laporan.
- Koneksi internet sekolah tidak bisa diandalkan (WiFi kadang mati).
- Sekolah menggunakan Kurikulum Merdeka sebagai acuan rapor.
- Data satu sekolah **tidak boleh terlihat** oleh sekolah lain dalam kondisi apapun.

### Batasan

- Tidak ada self-registration untuk sekolah — semua diinisiasi oleh super_admin.
- Tidak ada fitur live recording audio/video.
- Tidak ada engine rekomendasi AI otomatis di MVP.
- Rubrik per sekolah bisa dikustomisasi, tapi berbagi template rubrik antar sekolah belum ada di MVP.

---

## 4. Model Multi-Tenant

### Strategi: Shared Database + Shared Schema

Semua sekolah berbagi satu database dan satu set tabel. Isolasi data dijamin melalui kolom `school_id` di setiap tabel yang relevan, dikombinasikan dengan enforcement di application layer.

```
┌─────────────────────────────────────────────────────┐
│                  PostgreSQL Database                │
│                                                     │
│  schools table                                      │
│  ┌──────────┬──────────┬──────────┐                 │
│  │school_id │school_id │school_id │                 │
│  │  = 'A'   │  = 'B'   │  = 'C'  │                 │
│  └──────────┴──────────┴──────────┘                 │
│                                                     │
│  students, assessments, songs, dll — semua punya   │
│  kolom school_id. Query SELALU di-filter by         │
│  school_id dari session user yang login.            │
└─────────────────────────────────────────────────────┘
```

**Kenapa bukan Separate Database per Tenant?**

| Aspek | Shared Schema (dipilih) | Separate DB |
|---|---|---|
| Kompleksitas operasional | Rendah — 1 DB, 1 migration | Tinggi — N DB, N migration |
| Isolasi data | Application-level (cukup untuk use case ini) | Database-level (overkill) |
| Onboarding sekolah baru | Instant — insert ke tabel `schools` | Perlu provision DB baru |
| Cost infrastruktur | Rendah | Linear dengan jumlah sekolah |
| Risiko data leak | Ada jika ada bug di layer isolasi | Minimal |

Keputusan: **Shared Schema** cukup untuk skala awal. Mitigasi risiko leak via strict server-side filtering + audit log.

### Hierarki Tenant

```
Super Admin (platform owner)
    │
    ├── Sekolah A
    │       ├── School Admin A
    │       ├── Guru A1, A2
    │       └── Siswa A1, A2, ...
    │
    ├── Sekolah B
    │       ├── School Admin B
    │       ├── Guru B1
    │       └── Siswa B1, B2, ...
    │
    └── Sekolah C
            └── ...
```

---

## 5. Pengguna & Peran (RBAC)

### 5.1 Role Hierarchy

Relasi user↔role many-to-many via tabel `user_roles` — 1 user bisa punya lebih dari 1 role (mis. teacher yang juga parent). `school_id` tetap melekat di `users`, bukan di role.

| Role | Scope | Deskripsi |
|---|---|---|
| `super_admin` | Platform (lintas sekolah) | Pemilik platform. Membuat & mengelola sekolah, bisa akses semua data untuk keperluan support. |
| `school_admin` | 1 Sekolah | Admin per sekolah. Mengelola guru, siswa, konfigurasi rubrik & lagu milik sekolahnya. |
| `teacher` | 1 Sekolah | Guru. Input nilai, lihat rekap, generate laporan untuk siswa di sekolahnya. |
| `student` | Diri sendiri | Siswa. Read-only — hanya bisa lihat data diri sendiri. |
| `parent` | 1 Siswa | Orang tua. Read-only — hanya bisa lihat data anak yang terhubung. |

### 5.2 Permission Matrix

| Aksi | super_admin | school_admin | teacher | student | parent |
|---|---|---|---|---|---|
| Buat / hapus sekolah | ✅ | ❌ | ❌ | ❌ | ❌ |
| Lihat semua sekolah | ✅ | ❌ | ❌ | ❌ | ❌ |
| Buat akun guru / school_admin | ✅ | ✅ (sekolahnya) | ❌ | ❌ | ❌ |
| Buat akun siswa | ✅ | ✅ | ❌ | ❌ | ❌ |
| CRUD data siswa | ✅ | ✅ | Edit terbatas | ❌ | ❌ |
| Input penilaian | ✅ | ✅ | ✅ | ❌ | ❌ |
| Lihat nilai semua siswa | ✅ | ✅ | ✅ | ❌ | ❌ |
| Lihat nilai diri sendiri | — | — | — | ✅ | ✅ (anaknya) |
| Konfigurasi rubrik & lagu | ✅ | ✅ | ❌ | ❌ | ❌ |
| Export laporan | ✅ | ✅ | ✅ | ❌ | ❌ |
| Naik level siswa | ✅ | ✅ | ✅ | ❌ | ❌ |

### 5.3 Karakteristik Pengguna

**Super Admin**
- Jumlah: 1 (pemilik platform).
- Akses via desktop, tidak perlu mobile-optimized.
- Kebutuhan: dashboard overview sekolah, onboarding sekolah baru, support troubleshooting.

**School Admin**
- Jumlah: 1 per sekolah.
- Akses via desktop/HP.
- Kebutuhan: setup awal sekolah, manajemen guru & siswa, laporan per semester.

**Guru**
- Literasi digital: menengah.
- Konteks: di kelas, HP di tangan, waktu terbatas.
- Pain point: input manual lambat, takut salah ketik.

**Siswa & Orang Tua**
- Konteks: di rumah, via HP.
- Kebutuhan: lihat nilai & komentar terbaru.

---

## 6. User Stories & Acceptance Criteria

### SUPER ADMIN — Platform Management

**US-00: Dashboard Lintas Sekolah**
```
Sebagai super_admin, saya ingin melihat ringkasan semua sekolah
dalam satu dashboard agar bisa memonitor kesehatan platform.
```
**AC:**
- Tampil: jumlah sekolah aktif, total siswa, total penilaian bulan ini.
- List sekolah dengan: nama, jumlah guru, jumlah siswa, tanggal terakhir ada aktivitas.
- Bisa klik masuk ke dashboard sekolah tertentu (impersonate view, bukan impersonate user).

**US-01: Onboarding Sekolah Baru**
```
Sebagai super_admin, saya ingin menambahkan sekolah baru ke platform
dan membuat akun school_admin untuk mereka.
```
**AC:**
- Form: nama sekolah, kota, nomor telepon, nama school_admin, email school_admin.
- Sistem generate `school_id` unik.
- Akun school_admin dibuat otomatis, password sementara di-generate (random 12 karakter).
- Password sementara **ditampilkan sekali di UI** (response page dengan tombol copy) + dikirim via email ke email school_admin. Tidak ada WA otomatis — super_admin forward manual jika perlu.
- Sekolah baru langsung aktif tanpa approval tambahan.
- School_admin wajib ganti password saat login pertama (flag `force_password_reset BOOLEAN DEFAULT false` di tabel `users`, set `true` saat akun dibuat via super_admin).

**US-02: Nonaktifkan Sekolah**
```
Sebagai super_admin, saya ingin menonaktifkan sekolah
tanpa menghapus data historis mereka.
```
**AC:**
- Status sekolah berubah ke `inactive`.
- Semua user sekolah tersebut tidak bisa login.
- Data tetap ada di database untuk keperluan audit.
- Bisa diaktifkan kembali kapan saja.

---

### SCHOOL ADMIN — Manajemen Sekolah

**US-03: Setup Konfigurasi Awal Sekolah**
```
Sebagai school_admin, saya ingin mengatur konfigurasi dasar sekolah
(rubrik, daftar lagu) sebelum guru mulai menggunakan aplikasi.
```
**AC:**
- Rubrik default sudah tersedia (4 komponen × 25 poin).
- Bisa mengubah nama komponen rubrik (tanpa mengubah bobot di MVP).
- Bisa menambah / menonaktifkan lagu dari daftar master.

**US-04: Tambah Akun Guru**
```
Sebagai school_admin, saya ingin menambahkan guru baru
ke sekolah saya.
```
**AC:**
- Form: nama, email, nomor HP.
- Guru otomatis terikat ke `school_id` sekolah yang bersangkutan.
- Guru tidak bisa melihat data sekolah lain.
- Password sementara di-generate (random 12 karakter), ditampilkan sekali di UI + dikirim ke email guru. Flag `force_password_reset = true` agar ganti password saat login pertama.

**US-05: Tambah Siswa Baru**
```
Sebagai school_admin, saya ingin menambahkan siswa
dan membuatkan akun login untuk mereka.
```
**AC:**
- Form: nama lengkap, NIS, email (opsional), level awal, nama orang tua, nomor HP orang tua.
- NIS unik dalam scope sekolah (boleh sama dengan sekolah lain).
- Level default: Bronze.
- Akun siswa terikat ke `school_id`.

---

### GURU — Autentikasi

**US-06: Login via Google**
```
Sebagai guru, saya ingin login dengan akun Google saya
agar tidak perlu mengingat password tambahan.
```
**AC:**
- Hanya email yang sudah terdaftar oleh school_admin yang bisa login.
- Setelah login, session berisi `school_id` dari sekolah guru tersebut.
- Sesi bertahan minimal 7 hari.

**US-07: Login Manual**
```
Sebagai guru yang tidak punya akun Google,
saya ingin login dengan email dan password yang diberikan school_admin.
```
**AC:**
- Password di-hash bcrypt sebelum disimpan.
- Terdapat fitur "Lupa Password" via email.

---

### GURU — Manajemen Siswa & Penilaian

**US-08: Melihat Daftar Siswa**
```
Sebagai guru, saya ingin melihat semua siswa di sekolah saya
beserta level dan nilai terkini mereka.
```
**AC:**
- Hanya siswa dari `school_id` yang sama yang tampil — tidak ada bocoran data sekolah lain.
- Tabel memuat: nama, NIS, level, nilai sumatif terakhir, tanggal penilaian terakhir.
- Bisa difilter berdasarkan level.
- Bisa dicari berdasarkan nama atau NIS.

**US-09: Input Nilai via Mobile**
```
Sebagai guru, saya ingin menilai siswa langsung dari HP
dalam waktu < 60 detik per siswa.
```
**AC:**
- Tampilan mobile-first: satu layar tanpa scroll horizontal.
- Toggle Formatif/Sumatif terlihat jelas.
- Dropdown pilihan lagu hanya menampilkan lagu yang dikonfigurasi sekolah ini.
- 4 slider (0–25) untuk komponen rubrik sekolah ini.
- Total skor ter-update real-time.
- Kolom komentar opsional.
- Tombol "Simpan" satu kali klik.

**US-10: Input Nilai saat Offline**
```
Sebagai guru, saya ingin tetap bisa menilai siswa
meskipun WiFi sekolah sedang mati.
```
**AC:**
- Aplikasi terdeteksi offline → banner "Mode Offline".
- Data tersimpan di IndexedDB lokal (termasuk `school_id`).
- Saat koneksi kembali, data otomatis di-sync ke server.
- Sinkronisasi memverifikasi `school_id` sebelum INSERT.
- Banner berubah: "Sinkronisasi berhasil — X data telah diunggah."

**US-11: Memberikan Tantangan Ekstra**
```
Sebagai guru, saya ingin memberikan tantangan tambahan
kepada siswa berprestasi tanpa menaikkan levelnya.
```
**AC:**
- Checkbox "Tantangan Ekstra" di panel penilaian.
- Jika dicentang, muncul textarea untuk deskripsi tantangan.
- Tampil di portal siswa sebagai label khusus.

**US-12: Naik Level Siswa**
```
Sebagai guru, saya ingin menaikkan level siswa
setelah mereka lulus Mastery Test.
```
**AC:**
- Tombol "Naik Level" di profil siswa.
- Konfirmasi dialog sebelum eksekusi.
- Riwayat perubahan level tercatat (student_id, from, to, teacher_id, timestamp).

---

### GURU — Laporan & Ekspor

**US-13: Rekap Nilai Semua Siswa**
```
Sebagai guru, saya ingin melihat ringkasan nilai semua siswa sekolah saya
dalam satu tampilan tabel di desktop.
```
**AC:**
- Tabel matriks: baris = siswa, kolom = komponen rubrik + total.
- Filter: semester, level, lagu.
- Kolom dengan nilai terendah di-highlight.
- Hanya data siswa `school_id` yang sama.

**US-14: Generate Deskripsi Kurikulum Merdeka**
```
Sebagai guru, saya ingin sistem menghasilkan deskripsi otomatis untuk rapor.
```
**AC:**
- Deskripsi berbahasa Indonesia yang natural.
- Bisa diedit manual sebelum ekspor.
- Tombol "Regenerate" tersedia.

**US-15: Ekspor Nilai ke Excel**
```
Sebagai guru, saya ingin mengekspor nilai ke file Excel.
```
**AC:**
- File `.xlsx` dengan nama: `Nilai_Biola_[NamaSekolah]_[Semester]_[Tahun].xlsx`.
- Memuat: nama, NIS, nilai per komponen, total, deskripsi KurMer, level.
- Hanya data siswa sekolah guru yang bersangkutan.

---

### SISWA / ORANG TUA — Portal

**US-16: Lihat Profil & Level**
```
Sebagai siswa, saya ingin melihat level dan badge saya saat ini.
```
**AC:**
- Badge visual Bronze/Silver/Gold di halaman utama.
- Nama sekolah tampil di profil (untuk konfirmasi konteks).

**US-17: Lihat Grafik Perkembangan**
```
Sebagai siswa, saya ingin melihat grafik nilai latihan dari waktu ke waktu.
```
**AC:**
- Line/bar chart nilai Formatif per minggu.
- Nilai Sumatif ditandai berbeda.
- Hanya data diri sendiri — tidak ada akses ke data siswa lain.

**US-18: Lihat Histori Penilaian & Komentar**
```
Sebagai siswa/orang tua, saya ingin melihat riwayat nilai dan komentar guru.
```
**AC:**
- Daftar kronologis: tanggal, lagu, skor per komponen, komentar guru.
- Nilai Sumatif dibedakan secara visual dari Formatif.

---

## 7. Fitur & Spesifikasi Fungsional

### 7.1 Module Map

```
┌──────────────────────────────────────────────────────────────────┐
│                      APLIKASI BIOLA                              │
├──────────────────────┬───────────────────────────────────────────┤
│  SUPER ADMIN MODULE  │           SCHOOL MODULE                   │
│                      │                                           │
│ · Dashboard platform │  AUTH      STUDENT     ASSESSMENT         │
│ · Onboarding sekolah │  ────────  ──────────  ──────────────     │
│ · Manajemen sekolah  │  · OAuth   · CRUD      · Panel Cepat      │
│ · Akses lintas tenant│  · Email/  · Level     · Formatif/Sum     │
│                      │    Pass    · Profil    · Offline Sync     │
│                      │  · RBAC    · Portal    · Tantangan Extra  │
│                      │            View                           │
│                      │                                           │
│                      │  REPORT    CURRICULUM  ADMIN SCHOOL       │
│                      │  ────────  ──────────  ───────────────    │
│                      │  · Rekap   · Rubrik    · User Mgmt        │
│                      │  · KurMer  · Lagu      · Audit Log        │
│                      │  · Export  · Konfigur. · Settings         │
└──────────────────────┴───────────────────────────────────────────┘
```

### 7.2 Routing

| Route | Akses | Deskripsi |
|---|---|---|
| `/login` | Public | Halaman login |
| `/super-admin` | super_admin | Dashboard platform |
| `/super-admin/schools` | super_admin | Daftar & manajemen sekolah |
| `/super-admin/schools/new` | super_admin | Form onboarding sekolah baru |
| `/super-admin/schools/[id]` | super_admin | Detail & stats sekolah |
| `/dashboard` | school_admin, teacher | Overview sekolah aktif |
| `/students` | school_admin, teacher | Daftar siswa sekolah aktif |
| `/students/[id]` | school_admin, teacher | Profil detail siswa |
| `/students/[id]/assess` | teacher, school_admin | Panel penilaian |
| `/students/new` | school_admin | Form tambah siswa |
| `/reports` | school_admin, teacher | Rekap & ekspor nilai |
| `/curriculum` | school_admin | Konfigurasi rubrik & lagu |
| `/admin/users` | school_admin | Manajemen guru & akun |
| `/portal` | student, parent | Dashboard siswa (read-only) |
| `/portal/history` | student, parent | Histori penilaian |

---

## 8. Logika Bisnis & Aturan Penilaian

### 8.1 Struktur Level

| Level | Fokus | Repertoar |
|---|---|---|
| **Bronze** | Fondasi: postur, pegangan, intonasi dasar | Suzuki Buku 1 (awal) |
| **Silver** | Pengembangan: dinamika, shifting, slur | Suzuki Buku 1 (akhir) – Buku 2 |
| **Gold** | Kemahiran: vibrato, tempo cepat, teknik lanjut | Suzuki Buku 3+ |

**Aturan Promosi:**
- Hanya melalui **Mastery Test** yang ditandai guru saat submit penilaian.
- Riwayat tiap perubahan level dicatat lengkap.

### 8.2 Rubrik Penilaian

| Komponen | Bobot Maks |
|---|---|
| Postur & Teknik | 25 |
| Intonasi | 25 |
| Ritme & Tempo | 25 |
| Musikalitas & Ekspresi | 25 |

**Total: 100 poin.** Nama komponen bisa dikustomisasi per sekolah oleh school_admin.

### 8.3 Konversi Nilai

| Rentang | Indeks | Label |
|---|---|---|
| 90 – 100 | A | Sangat Baik |
| 80 – 89 | B | Baik |
| 70 – 79 | C | Cukup |
| < 70 | D | Perlu Bimbingan |

### 8.4 Tipe Penilaian

- **Formatif:** Progres harian/mingguan. Tidak dihitung ke nilai rapor. Divisualisasikan dalam grafik.
- **Sumatif:** Penilaian resmi. Rata-rata Sumatif = nilai akhir rapor. Basis promosi level.

### 8.5 Algoritma Auto-Deskripsi Kurikulum Merdeka

```
INPUT: skor[postur, intonasi, ritme, musikalitas]

1. komponen_kuat  = komponen dengan skor TERTINGGI
2. komponen_lemah = komponen dengan skor TERENDAH
3. indeks         = konversi total skor (A/B/C/D)

OUTPUT:
"Siswa menunjukkan penguasaan {komponen_kuat} yang {label_indeks},
dan membutuhkan pendampingan lebih lanjut dalam {komponen_lemah}."

EDGE CASE — semua skor sama:
"Siswa menunjukkan perkembangan yang merata di semua aspek penilaian."

EDGE CASE — total < 70:
Tambah kalimat: "Disarankan untuk meningkatkan intensitas latihan
di seluruh aspek penilaian."
```

---

## 9. Arsitektur Teknis

### 9.1 Tech Stack

**Arsitektur: Decoupled — Backend API terpisah dari Frontend UI.**

| Layer | Teknologi | Alasan |
|---|---|---|
| **Backend** | Go + Fiber | Efisiensi CPU/RAM jauh lebih baik dari Node.js; throughput API optimal |
| Query Layer | sqlc + pgx | Typesafe SQL queries tanpa ORM overhead; compile-time safety |
| DB Migration | golang-migrate | Migrasi berbasis file SQL, tidak terikat ORM |
| Auth | JWT (Go sign/verify) + PostgreSQL blacklist | Stateless, cocok decoupled; invalidasi via blacklist table |
| **Frontend** | Next.js 14 (App Router) + NextAuth | SSR/SSG untuk UI; NextAuth handles OAuth + CredentialsProvider; tukar JWT dengan Go |
| i18n | next-intl | Default: English (en). Supported: Indonesian (id). Locale dari browser preference + user setting. |
| UI | Tailwind CSS + shadcn/ui | Cepat, konsisten |
| Form | React Hook Form + Zod | Validasi typesafe di sisi klien |
| **Database** | PostgreSQL | Relasional, JSONB, battle-tested |
| Offline | Dexie.js (IndexedDB) | API ergonomis; request sync langsung ke Go API |
| Export | SheetJS (client-side) | Hemat server resource |
| Deployment | Docker Compose + VPS | Dua container: Go API + Next.js |

### 9.2 Arsitektur Deployment

```
                         VPS (Ubuntu)
┌────────────────────────────────────────────────────────┐
│                   Docker Compose                       │
│                                                        │
│  ┌───────────────────┐   ┌──────────────────────────┐  │
│  │  Go + Fiber API   │   │      PostgreSQL          │  │
│  │  (Port 8080)      │──▶│      (Port 5432)         │  │
│  │                   │   │                          │  │
│  │  Middleware:      │   │  Shared Database         │  │
│  │  · JWT auth       │   │  Shared Schema           │  │
│  │  · school_id      │   │  + school_id columns     │  │
│  │    injection      │   │  + jwt_blacklist table   │  │
│  │  · RBAC guard     │   └──────────────────────────┘  │
│  └─────────┬─────────┘             ▲                   │
│            │                       │                   │
│  ┌─────────┴─────────┐             │                   │
│  │  Next.js Frontend │─────────────┘                   │
│  │  (Port 3000)      │                                 │
│  │  UI only          │                                 │
│  └───────────────────┘                                 │
│                                                        │
└────────────────────────────────────────────────────────┘
          │ HTTPS
          ▼
    Browser / PWA
```

### 9.3 JWT & Tenant Context

Go API meng-issue JWT saat login. Token disimpan di browser sebagai **httpOnly cookie** (bukan localStorage — tidak accessible JavaScript).

```go
// JWT Claims
type JWTClaims struct {
    Sub             string   `json:"sub"`              // user id
    Email           string   `json:"email"`
    Name            string   `json:"name"`
    Roles           []string `json:"roles"`            // 1+ dari: super_admin | school_admin | teacher | student | parent
    SchoolID        string   `json:"school_id"`        // "" untuk super_admin
    JTI             string `json:"jti"`              // JWT ID untuk blacklist lookup
    PreferredLocale string `json:"preferred_locale"` // "en" | "id" — NextAuth set locale cookie dari sini
    Exp             int64  `json:"exp"`
}
```

Setiap request ke Go API, `school_id` diambil dari JWT claims yang sudah diverifikasi — **tidak pernah dari request body atau query param** — untuk mencegah tenant spoofing.

**Invalidasi token (logout paksa, user/sekolah dinonaktifkan):**

```sql
CREATE TABLE jwt_blacklist (
  jti        UUID PRIMARY KEY,
  expired_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_jwt_blacklist_expired_at ON jwt_blacklist(expired_at);
-- Cleanup periodic di run pada Go API startup secara background process (goroutine)
```

JWT middleware Go cek `jti` di tabel ini sebelum memproses request. Tanpa Redis di MVP — PostgreSQL cukup untuk skala awal.

### 9.4 Middleware Architecture (Go Fiber)

```
Request masuk ke Go Fiber
     │
     ▼
[JWT Middleware] — verifikasi signature, cek exp, cek jti di jwt_blacklist
     │
     ▼
[Role Guard] — verifikasi role punya akses ke route ini
     │
     ▼
[Tenant Guard] — inject school_id ke Fiber context
     │           super_admin: baca x-school-id header (Wajib: reject header jika role != super_admin)
     │           lainnya: school_id dari JWT claims (immutable)
     │           Log audit untuk super_admin akan merecord explicit `user_id=$super_admin_id` dan `school_id=$requested_school_id` untuk validasi silang.
     ▼
[Handler] — semua query WAJIB include WHERE school_id = c.Locals("school_id")
```

### 9.5 PWA & Offline Architecture

```
Browser
├── Service Worker
│   ├── Cache: static assets + App Shell
│   ├── Cache: GET /api/students (stale-while-revalidate)
│   └── Online detection → trigger sync
│
├── Dexie.js (IndexedDB)
│   ├── pending_assessments[] — {school_id, student_id, ...data, idempotency_key}
│   └── cached_students[]    — snapshot daftar siswa
│
└── React App
[Fact-Forcing Gate]
1. `CLAUDE.md` references `PRD-Aplikasi-Penilaian-Biola.md`.
2. No public functions/classes affected (markdown document).
3. No data files read/written.
4. User instruction: "yes please, i prefer use minio for storage"

    ├── Online:  fetch API → update IndexedDB cache
    ├── Offline: read IndexedDB → push ke pending_assessments
    ├── Reconnect (JWT Valid): flush pending_assessments → POST /api/assessments/batch
    └── Reconnect (JWT Expired): tangkap 401, simpan antrean, prompt re-login, lalu flush

```

**Offline idempotency key behavior:**

`idempotency_key` adalah UUID yang di-generate client saat assessment dibuat di IndexedDB, disimpan bersama data lokal. Server UNIQUE constraint mencegah duplikat insert. 
Server response 201 mengembalikan `idempotency_key`. Client menulis ini ke persistent config/queue.
Saat client resume, call `GET /api/assessments?idempotency_key=X` untuk cek apakah sudah terkirim. Jika ya, hapus dari queue pending. Clear cache aman karena state terverifikasi di server.

---

## 10. Schema Database

### 10.1 Entity Relationship Overview

```
schools
  │
  ├── users (school_id FK, nullable untuk super_admin)
  │       └── parent_students (parent_id FK → users, student_id FK → students)
  │
  ├── students (school_id FK)
  │       └── level_history (school_id FK)
  │
  ├── songs (school_id FK — per-sekolah atau shared template)
  │
  ├── rubric_components (school_id FK — per-sekolah)
  │
  └── assessments (school_id FK)
          ├── assessment_components
          └── assessment_notes
```

### 10.2 DDL Lengkap

```sql
-- ============================================================
-- TABEL: schools
-- ============================================================
CREATE TABLE schools (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name         VARCHAR(255) NOT NULL,
  city         VARCHAR(100),
  phone        VARCHAR(20),
  is_active    BOOLEAN DEFAULT true,
  created_by   UUID,  -- super_admin user id
  created_at   TIMESTAMPTZ DEFAULT NOW(),
  updated_at   TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- TABEL: users
-- Mencakup semua role: super_admin, school_admin, teacher, student, parent
-- Role dipindah ke tabel roles + user_roles (many-to-many, lihat bawah)
-- ============================================================
CREATE TABLE users (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id     UUID REFERENCES schools(id),  -- NULL untuk super_admin
  email         VARCHAR(255) UNIQUE,
  password_hash VARCHAR(255),                  -- NULL jika OAuth
  name          VARCHAR(255) NOT NULL,
  google_id     VARCHAR(255) UNIQUE,
  avatar_url    TEXT,
  is_active     BOOLEAN DEFAULT true,
  force_password_reset BOOLEAN DEFAULT false,  -- set true saat akun dibuat via admin. Middleware akan intercept request ke /reset-password jika true, redirect jika di frontend.
  preferred_locale     VARCHAR(5) DEFAULT 'en', -- i18n: 'en' | 'id'
  created_at    TIMESTAMPTZ DEFAULT NOW(),
  updated_at    TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- TABEL: roles
-- Daftar role yang bisa dimiliki user (lookup table)
-- ============================================================
CREATE TABLE roles (
  id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name  VARCHAR(20) UNIQUE NOT NULL
          CHECK (name IN ('super_admin','school_admin','teacher','student','parent'))
);

-- ============================================================
-- TABEL: user_roles
-- Junction many-to-many: 1 user bisa punya lebih dari 1 role
-- (contoh: teacher yang juga parent dari siswa lain)
-- ============================================================
CREATE TABLE user_roles (
  user_id  UUID REFERENCES users(id) NOT NULL,
  role_id  UUID REFERENCES roles(id) NOT NULL,
  PRIMARY KEY (user_id, role_id)
);
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);

-- ============================================================
-- TABEL: students
-- Data siswa terpisah dari users untuk fleksibilitas
-- (siswa bisa ada tanpa akun login)
-- ============================================================
CREATE TABLE students (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id     UUID REFERENCES schools(id) NOT NULL,
  user_id       UUID REFERENCES users(id),   -- NULL jika belum punya login
  nis           VARCHAR(50) NOT NULL,
  name          VARCHAR(255) NOT NULL,
  current_level VARCHAR(20) DEFAULT 'bronze'
                  CHECK (current_level IN ('bronze', 'silver', 'gold')),
  parent_name   VARCHAR(255),
  parent_phone  VARCHAR(20),
  notes         TEXT,
  created_at    TIMESTAMPTZ DEFAULT NOW(),
  updated_at    TIMESTAMPTZ DEFAULT NOW(),
  -- NIS unik dalam scope sekolah (boleh sama di sekolah berbeda)
  UNIQUE (school_id, nis)
);

-- ============================================================
-- TABEL: level_history
-- ============================================================
CREATE TABLE level_history (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id     UUID REFERENCES schools(id) NOT NULL,
  student_id    UUID REFERENCES students(id) NOT NULL,
  from_level    VARCHAR(20),
  to_level      VARCHAR(20) NOT NULL,
  changed_by    UUID REFERENCES users(id) NOT NULL,
  reason        TEXT,
  changed_at    TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- TABEL: rubric_components
-- Dikonfigurasi per sekolah. Default: 4 komponen × 25 poin.
-- ============================================================
CREATE TABLE rubric_components (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id     UUID REFERENCES schools(id) NOT NULL,
  name          VARCHAR(100) NOT NULL,
  max_score     INTEGER NOT NULL DEFAULT 25,
  order_index   INTEGER DEFAULT 0,
  is_active     BOOLEAN DEFAULT true
);

-- ============================================================
-- TABEL: songs
-- Daftar lagu per sekolah. school_admin bisa tambah/nonaktifkan.
-- ============================================================
CREATE TABLE songs (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id     UUID REFERENCES schools(id) NOT NULL,
  title         VARCHAR(255) NOT NULL,
  level         VARCHAR(20) NOT NULL
                  CHECK (level IN ('bronze', 'silver', 'gold')),
  book          VARCHAR(100),
  order_index   INTEGER DEFAULT 0,
  is_active     BOOLEAN DEFAULT true,
  created_at    TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- TABEL: assessments
-- Header penilaian per sesi
-- ============================================================
CREATE TABLE assessments (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id         UUID REFERENCES schools(id) NOT NULL,
  student_id        UUID REFERENCES students(id) NOT NULL,
  song_id           UUID REFERENCES songs(id) NOT NULL,
  teacher_id        UUID REFERENCES users(id) NOT NULL,
  type              VARCHAR(20) NOT NULL
                      CHECK (type IN ('formative', 'summative')),
  is_mastery_test   BOOLEAN DEFAULT false,
  total_score       INTEGER NOT NULL CHECK (total_score >= 0 AND total_score <= 100),
  grade_index       VARCHAR(5) CHECK (grade_index IN ('A', 'B', 'C', 'D')),
  comment           TEXT,
  idempotency_key   UUID UNIQUE,  -- dedup untuk offline submission
  assessed_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- Waktu server (UTC). Frontend mengatur konversi ke zona waktu lokal.
  created_at        TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- TABEL: assessment_components
-- Skor per komponen rubrik. Denormalized school_id untuk query safety.
-- ============================================================
CREATE TABLE assessment_components (
  id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  assessment_id         UUID REFERENCES assessments(id) NOT NULL,
  school_id             UUID REFERENCES schools(id) NOT NULL,
  rubric_component_id   UUID REFERENCES rubric_components(id) NOT NULL,
  score                 INTEGER NOT NULL CHECK (score >= 0 AND score <= 25)
);
CREATE INDEX idx_assessment_components_assessment_id ON assessment_components(assessment_id);
CREATE INDEX idx_assessment_components_school_id ON assessment_components(school_id);

-- ============================================================
-- TABEL: assessment_notes
-- Catatan tantangan ekstra dan catatan khusus lainnya
-- ============================================================
CREATE TABLE assessment_notes (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  assessment_id UUID REFERENCES assessments(id) NOT NULL,
  note_type     VARCHAR(50) DEFAULT 'extra_challenge',
  content       TEXT NOT NULL,
  created_at    TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- TABEL: password_reset_tokens
-- ============================================================
CREATE TABLE password_reset_tokens (
  token        VARCHAR(255) PRIMARY KEY,   -- random 32-char, bcrypt-hashed
  user_id      UUID REFERENCES users(id),
  expires_at   TIMESTAMPTZ,
  created_at   TIMESTAMPTZ DEFAULT NOW()
);

-- ============================================================
-- TABEL: parent_students
-- Relasi parent↔student agar portal ortu bisa filter "anak saya".
-- ============================================================
CREATE TABLE parent_students (
  parent_id  UUID REFERENCES users(id) NOT NULL,
  student_id UUID REFERENCES students(id) NOT NULL,
  school_id  UUID REFERENCES schools(id) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (parent_id, student_id)
);
CREATE INDEX idx_parent_students_parent_id ON parent_students(parent_id);
CREATE INDEX idx_parent_students_school_id ON parent_students(school_id);
```

### 10.3 Index

```sql
-- Performance indexes
CREATE INDEX idx_users_school_id ON users(school_id);
CREATE INDEX idx_students_school_id ON students(school_id);
CREATE INDEX idx_students_school_nis ON students(school_id, nis);
CREATE INDEX idx_assessments_school_id ON assessments(school_id);
CREATE INDEX idx_assessments_student_id ON assessments(student_id);
CREATE INDEX idx_assessments_school_type ON assessments(school_id, type);
CREATE INDEX idx_assessments_assessed_at ON assessments(assessed_at DESC);
CREATE INDEX idx_level_history_student_id ON level_history(student_id);
CREATE INDEX idx_level_history_school_id ON level_history(school_id);
CREATE INDEX idx_rubric_components_school_id ON rubric_components(school_id);
CREATE INDEX idx_songs_school_id ON songs(school_id);
CREATE INDEX idx_songs_school_level ON songs(school_id, level);
CREATE INDEX idx_assessment_components_assessment_id ON assessment_components(assessment_id);
```

---

## 11. Tenant Isolation & Row-Level Security

### 11.1 Application-Layer Enforcement (Primary)

Setiap query yang melibatkan data tenant WAJIB menyertakan `school_id` dari Fiber context:

```go
// ✅ BENAR — selalu filter by school_id dari JWT context
func GetStudents(c *fiber.Ctx) error {
    schoolID := c.Locals("school_id").(string) // di-inject oleh Tenant Guard

    rows, err := db.Query(ctx,
        `SELECT id, name, nis, current_level FROM students
         WHERE school_id = $1`,
        schoolID,
    )
    // ...
}

// ❌ SALAH — tidak ada filter school_id (security bug)
func GetStudents(c *fiber.Ctx) error {
    rows, err := db.Query(ctx, `SELECT * FROM students`)
    // ...
}
```

Pola wajib: `school_id` selalu berasal dari `c.Locals("school_id")` — tidak pernah dari `c.Body()`, `c.Query()`, atau `c.Params()`.

### 11.2 Cross-Tenant Validation

Setiap kali ada relasi antar entitas (misal: `student_id` dalam assessment), validasi bahwa entitas yang direferensikan memiliki `school_id` yang sama:

```go
// Validasi saat submit assessment
var studentSchoolID string
err := db.QueryRow(ctx,
    `SELECT school_id FROM students WHERE id = $1`,
    body.StudentID,
).Scan(&studentSchoolID)

if err != nil || studentSchoolID != c.Locals("school_id").(string) {
    return fiber.NewError(fiber.StatusForbidden, "Student not found in your school")
}
```

### 11.3 Super Admin Access Pattern

Super admin bisa akses data sekolah manapun dengan menyertakan `school_id` secara eksplisit di request header — bukan dari session. Middleware menolak header `x-school-id` jika token JWT yang dikirim tidak memiliki role `super_admin`.

```typescript
// Khusus super_admin
const schoolId = session.user.roles.includes('super_admin')
  ? request.headers.get('x-school-id')  // eksplisit dari header
  : session.user.school_id;              // dari session (immutable)
```

### 11.4 Audit Log

Semua operasi write (INSERT, UPDATE, DELETE) pada data tenant dicatat:

```sql
CREATE TABLE audit_logs (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id   UUID REFERENCES schools(id),
  user_id     UUID REFERENCES users(id),
  action      VARCHAR(50) NOT NULL,   -- 'create_student', 'submit_assessment', dll
  entity      VARCHAR(50) NOT NULL,   -- nama tabel
  entity_id   UUID,
  old_data    JSONB,
  new_data    JSONB,
  ip_address  INET,
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_school_id ON audit_logs(school_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
```

---

## 12. API Design

### 12.1 Konvensi URL

Semua endpoint di-serve oleh **Go + Fiber** (port 8080).

- `/api/super-admin/*` — hanya accessible oleh `super_admin`
- `/api/*` — accessible oleh user dengan JWT valid; `school_id` dari JWT claims otomatis
- Tidak ada `school_id` di URL path untuk endpoint non-super-admin (diambil dari JWT)
- Next.js **tidak punya API routes** — murni UI renderer

### 12.2 Super Admin Endpoints

```
GET    /api/super-admin/schools              — List semua sekolah + stats
POST   /api/super-admin/schools              — Buat sekolah baru + school_admin
GET    /api/super-admin/schools/:id          — Detail sekolah
PUT    /api/super-admin/schools/:id          — Update data sekolah
PATCH  /api/super-admin/schools/:id/status   — Aktif/nonaktifkan sekolah
GET    /api/super-admin/stats                — Platform-wide statistics
```

### 12.3 Auth Endpoints

```
-- Auth dihandle NextAuth (next-auth.js.org) di Next.js.
-- NextAuth mengelola Google OAuth + CredentialsProvider (email/password).
-- Setelah NextAuth verify credentials, ia POST ke Go untuk tukar JWT:
POST   /api/auth/token           — NextAuth memanggil ini; Go verifikasi user → issue JWT
                                   Response: { token, expires_at, school_id, role }
                                   JWT disimpan di NextAuth session (httpOnly cookie Next.js-managed)
-- Go API endpoints menerima JWT dari Authorization header yang di-inject NextAuth
POST   /api/auth/logout          — Blacklist jti di PostgreSQL, clear NextAuth session
POST   /api/auth/forgot-password — Kirim reset link via email
POST   /api/auth/reset-password  — Reset dengan token
-- Flow: Browser → NextAuth (Next.js :3000) → POST /api/auth/token (Go :8080) → JWT
-- Semua /api/* lain: NextAuth inject JWT ke Authorization header → Go verify → handle
```

### 12.4 Students Endpoints

```
GET    /api/students             — List siswa (school dari session)
POST   /api/students             — Tambah siswa [school_admin]
GET    /api/students/:id         — Detail siswa
PUT    /api/students/:id         — Update siswa [school_admin]
DELETE /api/students/:id         — Nonaktifkan [school_admin]
POST   /api/students/:id/promote — Naik level [teacher, school_admin]
GET    /api/students/:id/history — Riwayat level
```

### 12.5 Assessments Endpoints

```
GET    /api/assessments          — List penilaian (filter: student, type, date)
POST   /api/assessments          — Submit penilaian baru
GET    /api/assessments/:id      — Detail penilaian
PUT    /api/assessments/:id      — Edit (same-day only) [teacher]
DELETE /api/assessments/:id      — Hapus [school_admin]
POST   /api/assessments/batch    — Batch submit untuk offline sync
```

### 12.6 Reports Endpoints

```
GET    /api/reports/summary        — Rekap semua siswa sekolah aktif (Returns JSON: {students:[], assessments:[], rubric_components:[]} untuk diolah SheetJS di frontend)
GET    /api/reports/student/:id    — Rekap per siswa
POST   /api/reports/generate-desc  — Generate deskripsi KurMer
```

### 12.7 Curriculum Endpoints

```
GET    /api/songs                       — Daftar lagu sekolah aktif
POST   /api/songs                       — Tambah lagu [school_admin]
PUT    /api/songs/:id                   — Edit lagu [school_admin]
DELETE /api/songs/:id                   — Nonaktifkan [school_admin]
GET    /api/rubric-components           — Komponen rubrik sekolah aktif
PUT    /api/rubric-components/:id       — Edit nama komponen [school_admin]
```

### 12.8 School Admin — User Management Endpoints

```
GET    /api/admin/users              — List semua user sekolah aktif
POST   /api/admin/users              — Buat akun guru/siswa
PUT    /api/admin/users/:id          — Update akun
PATCH  /api/admin/users/:id/status   — Aktif/nonaktifkan akun
POST   /api/admin/users/:id/link-student   — Link parent ke student
DELETE /api/admin/users/:id/link-student/:student_id — Unlink
```

### 12.9 Portal Endpoints (Siswa/Ortu)

```
GET    /api/portal/me             — Profil siswa yang login (student role: diri sendiri via JWT sub;
                                    parent role: list anak via parent_students JOIN students)
GET    /api/portal/progress       — Data grafik perkembangan (student: filter student_id = self;
                                    parent: requires ?student_id=, validate via parent_students)
GET    /api/portal/assessments    — Histori penilaian (same scoping logic as /portal/progress)
POST   /api/admin/users/:id/link-student   — Link parent ke student (Requires school_admin)
DELETE /api/admin/users/:id/link-student/:student_id — Unlink
-- parent_students tabel menjadi sumber kebenaran untuk relasi parent↔student.
```

### 12.10 Contoh Request/Response

**POST /api/assessments** (school_id diambil dari session, bukan body)

```json
Request body:
{
  "student_id": "uuid-siswa",
  "song_id": "uuid-lagu",
  "type": "summative",
  "is_mastery_test": false,
  "components": [
    { "rubric_component_id": "uuid-postur", "score": 20 },
    { "rubric_component_id": "uuid-intonasi", "score": 18 },
    { "rubric_component_id": "uuid-ritme", "score": 22 },
    { "rubric_component_id": "uuid-musikalitas", "score": 19 }
  ],
  "comment": "Intonasi nada tinggi masih perlu latihan.",
  "extra_challenge": "Tambahkan vibrato di birama 8–12.",
  "assessed_at": "2026-07-18T09:30:00Z",
  "idempotency_key": "uuid-client-generated"
}

Response 201:
{
  "id": "uuid-assessment",
  "school_id": "uuid-sekolah",       -- dikembalikan untuk konfirmasi
  "total_score": 79,
  "grade_index": "C",
  "created_at": "2026-07-18T09:30:05Z"
}

Response 409 (duplikat idempotency_key — offline resubmit):
{
  "error": "DUPLICATE_SUBMISSION",
  "existing_id": "uuid-assessment-lama"
}
```

**GET /api/super-admin/schools** 

```json
Response 200:
{
  "schools": [
    {
      "id": "uuid-sekolah-a",
      "name": "SDN Merdeka 1",
      "city": "Bandung",
      "is_active": true,
      "stats": {
        "teacher_count": 3,
        "student_count": 87,
        "assessment_count_this_month": 124,
        "last_activity_at": "2026-07-17T14:22:00Z"
      }
    }
  ],
  "total": 7
}
```

---

## 13. Non-Functional Requirements

### 13.1 Performance

| Metrik | Target |
|---|---|
| First Contentful Paint | < 2 detik (koneksi 4G) |
| Load halaman daftar siswa | < 1 detik |
| Simpan penilaian (online) | < 500ms |
| Sync offline queue (10 entries) | < 3 detik |
| Query rekap nilai (100 siswa) | < 2 detik |

### 13.2 Ketersediaan

- **PWA Manifest & Service Worker:** Frontend Next.js menggunakan package `next-pwa` untuk manifest.json generation dan Service Worker orchestration.
- Target uptime: 99% (~7 jam downtime/bulan).
- Backup database harian otomatis (pg_dump).
- Zero-downtime deployment via Docker rolling update.

### 13.3 Keamanan & Isolasi Tenant

- HTTPS wajib untuk semua komunikasi.
- `school_id` TIDAK PERNAH diambil dari request body/query param untuk operasi non-super-admin.
- Cross-tenant validation pada setiap relasi FK.
- Audit log untuk semua operasi write.
- RBAC enforcement di server-side.
- Rate limiting: auth endpoints max 10 req/menit per IP.
- bcrypt cost factor ≥ 12 untuk password hash.
- Session database-backed (bukan pure JWT) agar bisa invalidate.

### 13.4 Aksesibilitas & Responsivitas

- Mobile-first (breakpoint: 375px, 768px, 1280px).
- Slider penilaian one-finger friendly.
- Kontras WCAG AA minimum.

### 13.5 Maintainability

- Semua schema change via **golang-migrate** (file SQL murni). Tidak ada Prisma di stack ini.
- Environment variables untuk semua config sensitif.
- Docker Compose untuk reproducibility.
- `c.Locals("school_id")` wajib disertakan di setiap query Go — tidak boleh ada query tabel tenant tanpa filter `school_id`.

### 13.6 Internationalization (i18n)

- Library: **next-intl** (App Router native, no pages/ workaround needed).
- Supported locales: `en` (English, default), `id` (Bahasa Indonesia).
- Locale detection order: user profile setting → `Accept-Language` header → fallback `en`.
- URL strategy: **path prefix** — `/en/*` and `/id/*`. Root `/` redirects to detected locale.
- All UI copy stored in `/messages/en.json` and `/messages/id.json`.
- KurMer auto-description (§8.5) always outputs Indonesian regardless of UI locale — it is a report artifact, not UI copy.
- Dates: formatted per locale (`en`: `Jul 18, 2026` / `id`: `18 Jul 2026`).
- Numbers: no locale number formatting — all scores use `Geist Mono`, locale-neutral.
- Excel export filenames stay Indonesian (`Nilai_Biola_...`) — these are artifacts consumed by Indonesian schools.
- `users.preferred_locale VARCHAR(5) DEFAULT 'en'` — add column to `users` table. Server reads this after login to set cookie; client doesn't need to detect again.

### 13.7 Theme Switcher

- Supported themes: `light` (default), `dark`.
- Implementation: **`next-themes`** — wraps `<html>` with `data-theme` + Tailwind `darkMode: 'class'`. Zero custom code.
- Persistence: `localStorage` via next-themes built-in. No DB column — theme is a display preference, not account data.
- Default: `light`. System preference (`prefers-color-scheme`) respected on first visit before user sets explicit preference.
- Dark palette (Tailwind `dark:` classes):
  - Background: `#1C1917` (Ink Deep becomes canvas)
  - Surface: `#292524` (Stone-800)
  - Text primary: `#FAF8F5` (Parchment becomes text)
  - Text secondary: `#A8A29E` (Stone-400)
  - Border: `rgba(87,83,78,0.5)` (Stone-600 at 50%)
  - Accent Rosewood: unchanged `#B5603C` — holds on dark bg, no adjustment needed.
- Theme toggle: icon button (sun/moon), top-right nav. No label text. 44px touch target.
- SSR: `suppressHydrationWarning` on `<html>` — next-themes standard pattern, one line.

---

## 14. Out of Scope (MVP)

- **Self-registration sekolah** — semua onboarding manual oleh super_admin.
- **Notifikasi push / WhatsApp otomatis** saat nilai baru masuk.
[Fact-Forcing Gate]
1. `CLAUDE.md` references `PRD-Aplikasi-Penilaian-Biola.md`.
2. No public functions/classes affected (markdown document).
3. No data files read/written.
4. User instruction: "yes please, i prefer use minio for storage"

- **Avatar / File Storage:** MinIO (S3-compatible) untuk avatar/file, di-deploy via Docker Compose bersama API & DB.
- **AI rekomendasi latihan** berdasarkan pola nilai.
- **Shared song template lintas sekolah** — tiap sekolah manage lagu sendiri di MVP.
- **Billing / subscription management** — belum ada di MVP.
- **Kalender jadwal les** dan pengingat.
- **Pembayaran SPP** atau fitur keuangan.
- **Analytics lintas sekolah** (perbandingan performa antar sekolah).
- **Separate database per tenant** — bisa dimigrasi nanti jika regulasi data mengharuskan.

---

## 15. Risiko & Mitigasi

| Risiko | Probabilitas | Dampak | Mitigasi |
|---|---|---|---|
| Data satu sekolah bocor ke sekolah lain | Rendah | Sangat Tinggi | Tenant Guard middleware + `c.Locals("school_id")` di setiap query + cross-tenant validation + audit log |
| Super admin salah akses data sekolah yang salah | Sedang | Tinggi | Explicit x-school-id header; UI konfirmasi "Anda sedang melihat data [NamaSekolah]" |
| Guru tidak konsisten menggunakan aplikasi | Tinggi | Tinggi | UX sesederhana mungkin; onboarding awal semester |
| Kehilangan data offline jika browser di-clear | Sedang | Tinggi | Banner jumlah pending entries; ingatkan sync sebelum clear cache |
| Konflik data saat sync offline | Rendah | Sedang | Idempotency key + 409 response handling |
| VPS down | Rendah | Sangat Tinggi | Backup harian + monitoring uptime sederhana |
| School admin menonaktifkan akun guru yang sedang aktif | Rendah | Sedang | Session invalidation segera saat status user berubah |
| Satu sekolah nonaktif tapi datanya diakses via URL langsung | Rendah | Tinggi | Middleware cek `school.is_active` sebelum setiap request |

---

## 16. Milestone & Prioritas

### Phase 1 — Foundation & Multi-Tenant Core (Minggu 1–3)
- [ ] Setup project: Go + Fiber (backend) + Next.js (frontend) + PostgreSQL + Docker Compose
- [ ] Schema database lengkap dengan `school_id` di semua tabel + `jwt_blacklist`
- [ ] golang-migrate: migrasi database dari file SQL
- [ ] JWT middleware + Tenant Guard + RBAC middleware di Go Fiber
- [ ] Auth: Google OAuth callback di Go + Email/Password + JWT issue/invalidate
- [ ] Super admin: onboarding sekolah + buat school_admin
- [ ] Seed data: 1 sekolah demo + rubrik default + lagu Suzuki

### Phase 2 — School Admin & Guru (Minggu 4–5)
- [ ] School admin dashboard: manajemen guru & siswa
- [ ] Konfigurasi rubrik & lagu per sekolah
- [ ] CRUD siswa dengan isolasi school_id

### Phase 3 — Core Assessment (Minggu 6–8)
- [ ] Panel penilaian mobile (slider, toggle, dropdown)
- [ ] Submit penilaian + cross-tenant validation
- [ ] Offline mode (Dexie.js) + sync dengan idempotency key
- [ ] Naik level siswa + riwayat level

### Phase 4 — Reporting (Minggu 9–10)
- [ ] Rekap nilai per sekolah (desktop view)
- [ ] Generator deskripsi Kurikulum Merdeka
- [ ] Export Excel per sekolah

### Phase 5 — Portal Siswa (Minggu 11)
- [ ] Dashboard siswa/ortu (read-only)
- [ ] Grafik perkembangan
- [ ] Histori penilaian & komentar

### Phase 6 — Polish & Deploy (Minggu 12–13)
- [ ] Super admin dashboard: stats lintas sekolah
- [ ] Audit log viewer (untuk super_admin)
- [ ] PWA manifest + service worker
- [ ] Security hardening + penetration test ringan
- [ ] Deployment ke VPS + domain + SSL
- [ ] UAT dengan 1–2 sekolah pilot

---

## 17. Keputusan Teknis Terbuka

| # | Keputusan | Opsi | Status |
|---|---|---|---|
| 1 | Query layer Go | sqlc vs pgx langsung | **sqlc** — generate typesafe code dari SQL; lebih maintainable jangka panjang |
| 2 | Auth mechanism | JWT + PG blacklist vs database sessions | ✅ **DIPUTUSKAN: JWT + PostgreSQL blacklist** — stateless, cocok decoupled, invalidasi tetap bisa dilakukan tanpa Redis |
| 3 | Offline library | idb vs Dexie.js | **Dexie.js** — API lebih ergonomis; request sync langsung ke Go API via `/api/*` |
| 4 | Chart library | Recharts vs Chart.js | **Recharts** — paling kompatibel dengan React |
| 5 | Export Excel | Client (SheetJS) vs Server (Go) | **SheetJS client-side** — hemat server resource, cukup untuk ukuran data ini |
| 6 | Redis? | JWT blacklist cache, rate limiting | **Tidak di MVP** — PostgreSQL cukup untuk blacklist + rate limiting sederhana |
| 7 | Row-Level Security di DB | PostgreSQL RLS vs app-layer filtering | **App-layer dulu** (handler filter by `school_id`) — lebih mudah debug. RLS bisa ditambahkan sebagai safety net nanti |
| 8 | Song library | Shared template global vs per-sekolah | **Per-sekolah di MVP** — lebih simpel. Global template bisa di-backlog |
| 9 | Impersonate sekolah untuk super_admin | x-school-id header vs query param | ✅ **DIPUTUSKAN: x-school-id header** — tidak tercatat di server logs URL |
| 10 | DB migration tool | golang-migrate vs goose | **golang-migrate** — file SQL murni, tidak tied ke Go struct |
|| 11 | OAuth di arsitektur decoupled | Google OAuth di Go vs Next.js (NextAuth) | **NextAuth di Next.js** — NextAuth handles OAuth + CredentialsProvider; POST /api/auth/token ke Go untuk issue JWT |

---

*Dokumen ini adalah living document. Update setiap kali ada keputusan baru atau perubahan scope.*
