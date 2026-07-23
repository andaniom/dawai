# DAWAI Rewrite Design Specification

**Date:** July 23, 2026  
**Status:** Approved  
**Scope:** MVP — Multi-subject assessment platform with multi-tenant isolation

---

## 1. Product Overview

**DAWAI v3** — Generic multi-subject assessment platform. Teachers score students on any subject (violin, math, English, etc.) using customizable rubrics. Multi-tenant (schools isolated), multi-school teachers supported. Targets: school admins, teachers, students.

**MVP Scope:**
- ✅ Auth (email/password + Google OAuth)
- ✅ Multi-tenant by school
- ✅ Multi-subject per school
- ✅ Teacher assessment entry (scores + rubric components)
- ✅ Student dashboard (view own assessments)
- ✅ Admin dashboard (manage users, subjects, rubrics)

**Phase 2 (post-launch):**
- Parent dashboard, Export (PDF/Excel), Offline entry, Analytics

---

## 2. Data Model

### Core Entities

```
School (tenant root)
├── Subject (Math, English, Violin, etc.)
│   ├── Rubric (active rubric for subject)
│   │   └── RubricComponent (e.g., "Intonation", "Rhythm")
│   └── Assessment (teacher scores student)
│       └── AssessmentComponent (score per rubric component)
├── User (email, password, roles)
│   ├── UserRole (super_admin, school_admin, teacher, student, parent)
│   └── UserSchool (junction: teacher can belong to multiple schools)
└── Student (per school, linked to User)
    └── ParentStudent (parent-child relationship)
```

### Key Relations

| Relation | Cardinality | Notes |
|----------|-------------|-------|
| School → Subject | 1:N | One school has many subjects |
| Subject → Rubric | 1:1 | One active rubric per subject |
| Rubric → RubricComponent | 1:N | Rubric defines 3–10 components |
| Teacher → School (UserSchool) | N:M | Teacher can teach at multiple schools |
| Teacher → Assessment | 1:N | Teacher submits many assessments |
| Student → Assessment | 1:N | Student has many assessments (per subject/date) |
| Assessment → AssessmentComponent | 1:N | One assessment scores all rubric components |

### Multi-Tenant Isolation

Every table has `school_id` column. All queries filtered by `school_id` from JWT claims. No cross-school data leaks.

Exception: `User` table shared across schools, but `UserSchool` junction enforces school assignment.

---

## 3. Frontend Architecture

**Stack:** Next.js 15 (App Router) + TypeScript + shadcn/ui + TanStack Query + Zustand + react-hook-form

### Directory Structure

```
app/
├── (auth)/
│   ├── login/page.tsx
│   ├── register/page.tsx
│   └── reset-password/page.tsx
├── (dashboard)/
│   ├── layout.tsx               # Sidebar + header
│   ├── subjects/                # Admin: manage subjects + rubrics
│   │   ├── page.tsx
│   │   ├── [id]/edit/page.tsx
│   │   └── new/page.tsx
│   ├── users/                   # Admin: manage teachers/students
│   │   ├── page.tsx
│   │   └── [id]/edit/page.tsx
│   ├── assessments/             # Teacher: submit assessments
│   │   ├── page.tsx             # Student list + inline form
│   │   └── [id]/edit/page.tsx
│   └── results/                 # Student: view own assessments
│       ├── page.tsx
│       └── [id]/page.tsx
├── api/
│   ├── auth/[...nextauth]/route.ts
│   └── (all other routes proxy to Go backend)
└── layout.tsx                   # Global layout + theme provider

components/
├── ui/                          # shadcn/ui (Button, Card, Table, Dialog, etc.)
├── forms/
│   ├── LoginForm.tsx
│   ├── SubjectForm.tsx
│   ├── AssessmentForm.tsx
│   └── UserForm.tsx
├── assessment/
│   ├── StudentList.tsx          # Table with TanStack Table
│   ├── AssessmentModal.tsx      # Inline assessment form
│   └── AssessmentResult.tsx     # View single assessment
└── shared/
    ├── Header.tsx
    ├── Sidebar.tsx
    ├── SchoolSwitcher.tsx       # Multi-school teacher nav
    └── ErrorBoundary.tsx

hooks/
├── useAuth.ts                   # Zustand auth store
├── useAssessments.ts            # TanStack Query assessments
├── useSubjects.ts               # TanStack Query subjects
└── useSchool.ts                 # Current school from Zustand

lib/
├── api.ts                       # Axios client + JWT interceptor
├── store.ts                     # Zustand stores (auth, ui)
├── types.ts                     # Shared TypeScript types
└── utils.ts                     # Helpers (validation, formatting)

messages/
├── en.json                      # English i18n (Phase 2)
└── id.json                      # Indonesian i18n (Phase 2)
```

### State Management

**Zustand stores:**
- `authStore` — user, roles, school_id, accessToken
- `uiStore` — currentSchool, sidebarOpen, theme (light/dark)

**TanStack Query:**
- `useQuery('subjects', ...)` — cache subjects
- `useQuery('assessments', ...)` — cache assessments (per student/subject)
- `useMutation('submitAssessment', ...)` — optimistic update

### Component Props (Type-Safe)

```typescript
// No `any` type. Strict mode everywhere.
interface AssessmentFormProps {
  studentId: string;
  subjectId: string;
  rubricComponents: RubricComponent[];
  onSubmit: (data: AssessmentInput) => Promise<void>;
}

interface StudentListProps {
  schoolId: string;
  subjectId: string;
  students: Student[];
  onAssessClick: (studentId: string) => void;
}
```

---

## 4. Backend Architecture

**Stack:** Go 1.23 + Fiber + PostgreSQL + sqlc + golang-migrate

### Response Format (All Endpoints)

**Success (2xx):**
```json
{
  "success": true,
  "code": 200,
  "data": { /* response body */ },
  "error": null,
  "meta": {
    "timestamp": "2026-07-23T17:00:00Z",
    "path": "/api/assessments",
    "version": "v1"
  }
}
```

**Error (4xx/5xx):**
```json
{
  "success": false,
  "code": 400,
  "data": null,
  "error": {
    "message": "Student not found in your school",
    "type": "validation_error",
    "details": { "student_id": "invalid UUID" }
  },
  "meta": {
    "timestamp": "2026-07-23T17:00:00Z",
    "path": "/api/assessments",
    "version": "v1"
  }
}
```

### Directory Structure

```
backend/
├── cmd/api/
│   └── main.go
├── migrations/
│   └── *.up.sql, *.down.sql     # golang-migrate
├── queries/
│   └── *.sql                     # sqlc source
├── internal/
│   ├── handlers/                 # HTTP handlers
│   │   ├── auth.go
│   │   ├── subjects.go
│   │   ├── assessments.go
│   │   ├── users.go
│   │   └── middleware.go
│   ├── services/                 # Business logic
│   │   ├── auth.go
│   │   ├── assessment.go
│   │   └── subject.go
│   ├── models/                   # sqlc-generated + custom types
│   │   └── types.go
│   ├── middleware/
│   │   ├── jwt_guard.go
│   │   ├── tenant_guard.go       # school_id from JWT
│   │   └── role_guard.go
│   └── config/
│       └── config.go
├── go.mod
└── Dockerfile
```

### Key Endpoints

```
# Auth
POST   /api/auth/login              { email, password } → { accessToken, user }
POST   /api/auth/register           { email, password, name } → { user }
POST   /api/auth/logout             → blacklist JWT
POST   /api/auth/token              (called by NextAuth) → { accessToken }

# Subjects (admin only)
GET    /api/subjects                → { data: [Subject] }
POST   /api/subjects                { name, description } → { data: Subject }
PATCH  /api/subjects/:id            { name, description } → { data: Subject }

# Rubric Components (admin only)
GET    /api/subjects/:id/rubric     → { data: [RubricComponent] }
POST   /api/subjects/:id/rubric     { name, scale, weight } → { data: RubricComponent }

# Assessments (teacher/student/admin)
GET    /api/assessments             → { data: [Assessment], meta: { pagination } }
POST   /api/assessments             { studentId, subjectId, scores: [] } → { data: Assessment }
PATCH  /api/assessments/:id         { scores: [] } → { data: Assessment }
GET    /api/assessments/:id         → { data: Assessment }

# Users (admin only)
GET    /api/users                   → { data: [User] }
POST   /api/users                   { email, name, role, schools: [] } → { data: User }
PATCH  /api/users/:id               { role, schools: [] } → { data: User }

# Me
GET    /api/me                      → { data: { user, school_id, roles, schools: [] } }

# Students (teacher/admin)
GET    /api/students                → { data: [Student] } (filtered by school_id)
GET    /api/students/:id            → { data: Student }
```

### Multi-Tenant Security

**JWT claims:**
```json
{
  "sub": "user_id",
  "school_ids": ["school_1", "school_2"],
  "roles": ["teacher", "super_admin"],
  "iat": 1234567890,
  "exp": 1234567890 + 7*24*60*60
}
```

**Middleware chain:**
1. `JWTGuard` — validate JWT signature + expiry
2. `TenantGuard` — extract school_id from JWT, set `c.Locals("school_id")`
3. `RoleGuard` — check endpoint's required role in JWT roles array

**Cross-tenant FK validation on writes:**
```go
// Before INSERT Assessment, verify student's school_id matches request school_id
var studentSchoolID string
err := db.QueryRow(ctx,
    `SELECT school_id FROM students WHERE id = $1`,
    req.StudentID,
).Scan(&studentSchoolID)

if err != nil || studentSchoolID != c.Locals("school_id").(string) {
    return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
        Message: "Student not found in your school",
    })
}
```

---

## 5. Key User Flows

### Flow 1: Teacher Submits Assessment

1. Teacher logs in → stored in Zustand authStore
2. Dashboard shows school switcher (if multiple schools)
3. Pick school → Pick subject → See StudentList table
4. Click "Assess" on student row → AssessmentModal opens
5. Form shows rubric components (e.g., Intonation 1–5, Rhythm 1–5)
6. Teacher enters scores + optional comment
7. Submit → POST /api/assessments → optimistic update → row updates with timestamp

### Flow 2: Student Views Own Assessments

1. Student logs in
2. Dashboard shows "My Assessments" (filtered by student_id)
3. Table shows: Subject, Date, Score, Status (submitted/pending)
4. Click row → Detail view shows rubric breakdown + teacher comment

### Flow 3: Admin Sets Up Subject

1. Admin logs in
2. Subjects page → "New Subject"
3. Form: name, description
4. Submit → POST /api/subjects
5. Then "Edit Rubric" → Add components (name, score scale 1–5, weight)
6. POST /api/subjects/:id/rubric for each component

---

## 6. Error Handling & Validation

**Frontend:**
- Form validation via react-hook-form + Zod schema
- shadcn/ui error messages inline on fields
- Toast notifications for API errors (optimistic update failures)

**Backend:**
- Input validation in handlers (email format, UUID, required fields)
- Business logic validation in services (cross-tenant FK, role checks)
- Error response format: type + message + details
- Rate limiting on auth endpoints (10 req/min per IP)

---

## 7. Security

**Multi-tenant isolation (critical):**
- Every query filters by school_id from JWT
- Cross-tenant FK checks before writes
- No RLS (app-layer filtering only) — code review must verify filters

**Authentication:**
- JWT (7-day lifetime) issued by Go backend
- bcrypt (cost ≥ 12) for password hashing
- NextAuth v5 handles Google OAuth + credential flow

**HTTPS:**
- Development: HTTP (docker compose)
- Production: HTTPS mandatory (Certbot for SSL cert renewal)

---

## 8. Testing

**Backend:**
- Unit tests for services (auth, assessment logic)
- Integration tests for endpoints (multi-tenant isolation, cross-tenant FK)
- Run: `go test ./...`

**Frontend:**
- Component tests (shadcn/ui components via Vitest + React Testing Library)
- Integration tests (form submission, data loading)
- E2E tests (login → assess → view results via Playwright)

**MVP minimum:**
- 80% backend coverage (critical paths)
- Component smoke tests (buttons render, forms submit)
- E2E happy path (login → assess → view)

---

## 9. Deployment

**Local development:**
```bash
docker compose up -d
# postgres, go-api, next-frontend all start
```

**CI/CD (Woodpecker):**
- test-backend: `go test ./...`
- build-backend: docker build
- build-frontend: docker build
- deploy: docker compose up (on main branch)

**Production:**
- VPS (t2.small or equivalent)
- Docker Compose orchestration
- Postgres backup strategy (TBD)
- Certbot for SSL cert renewal

---

## 10. Success Criteria

- ✅ Teacher can submit assessment in < 60 seconds
- ✅ Student sees assessment results immediately
- ✅ Admin can set up school + subjects + rubrics
- ✅ Zero cross-tenant data leaks (security audit)
- ✅ All endpoints respond with standard format
- ✅ 80%+ backend test coverage
- ✅ UI accessible (keyboard nav, screen reader support via shadcn/ui)
- ✅ Responsive on mobile + desktop

---

## 11. Open Questions / Phase 2

- Parent dashboard (when phase 2 starts)
- Offline assessment caching (Dexie.js + service worker)
- Analytics dashboard (progress tracking)
- Export formats (PDF, Excel via SheetJS)
- i18n (en/id via next-intl)
- Dark mode theme toggle

---

**Design approved by:** User (2026-07-23)
