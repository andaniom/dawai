# DAWAI — Sistem Penilaian Biola Berbasis Web

Multi-tenant violin assessment platform for schools. Built with Go + Fiber backend, Next.js frontend, PostgreSQL, and MinIO.

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.22+ (for local development)
- Node.js 20+ (for frontend development)

### Setup

1. **Clone and configure**
   ```bash
   cd /Users/andaniom/Documents/Martin/dawai
   cp .env.example .env
   # Edit .env with your values (or use defaults for local dev)
   ```

2. **Start services**
   ```bash
   docker-compose up -d
   ```

3. **Initialize database**
   ```bash
   docker-compose exec api migrate -path migrations -database "$DATABASE_URL" up
   ```

4. **Access services**
   - Frontend: http://localhost:3000
   - API: http://localhost:8080
   - MinIO: http://localhost:9001
   - PostgreSQL: localhost:5432

### Local Development

#### Backend
```bash
cd backend

# Install dependencies
go mod download

# Generate sqlc code
sqlc generate

# Run migrations
migrate -path migrations -database "postgres://dawai:change_me@localhost:5432/dawai?sslmode=disable" up

# Start server
go run cmd/api/main.go
```

#### Frontend
```bash
cd frontend

# Install dependencies
npm install

# Start dev server
npm run dev

# Build for production
npm run build
npm start
```

## Architecture

### Backend
- **Framework:** Go 1.22 + Fiber
- **Database:** PostgreSQL with golang-migrate
- **Query Layer:** sqlc for typesafe SQL
- **Authentication:** JWT + PostgreSQL blacklist
- **Rate Limiting:** In-memory (10 req/min on auth)

### Frontend
- **Framework:** Next.js 14 (App Router)
- **Auth:** NextAuth v5 (Google OAuth + Credentials)
- **Styling:** Tailwind CSS
- **i18n:** next-intl (English, Indonesian)
- **Theme:** next-themes (light/dark)
- **Offline:** Dexie.js (IndexedDB)
- **Charts:** Recharts
- **Export:** SheetJS

### Database
- **14 tables** with multi-tenant `school_id` isolation
- **Application-layer** enforcement (no RLS)
- **Shared schema** across all schools
- All tables indexed on `school_id` for query performance

## API Endpoints

### Authentication (Public)
```
POST   /api/auth/token            — Issue JWT
POST   /api/auth/logout           — Blacklist token
POST   /api/auth/forgot-password  — Request password reset
POST   /api/auth/reset-password   — Reset password with token
```

### Health
```
GET    /health                    — Service health check
```

### To Be Implemented (Phases 5-14)
- `/api/super-admin/*` — Super admin school management
- `/api/students` — Student CRUD
- `/api/assessments` — Assessment submission
- `/api/reports` — Recap and KurMer generation
- `/api/songs` — Curriculum management
- `/api/portal` — Student/parent read-only access

## File Structure

```
dawai/
├── backend/
│   ├── migrations/        # 14 SQL schema migrations
│   ├── queries/          # sqlc SQL query definitions
│   ├── internal/
│   │   ├── config/       # Configuration
│   │   ├── middleware/   # JWT, Rate Limit, Tenant, Role
│   │   ├── handlers/     # HTTP request handlers
│   │   ├── services/     # Business logic
│   │   └── models/       # Data structures
│   ├── cmd/api/          # Application entrypoint
│   ├── router/           # Route setup
│   ├── go.mod            # Dependencies
│   └── Dockerfile        # Backend image
│
├── frontend/
│   ├── app/[locale]/     # Next.js App Router pages
│   ├── components/       # React components
│   ├── lib/              # Utilities (auth, i18n)
│   ├── messages/         # Translations (en.json, id.json)
│   ├── package.json      # Dependencies
│   └── Dockerfile        # Frontend image
│
├── docker-compose.yml    # Service orchestration
├── .env.example          # Environment template
├── CLAUDE.md             # Project instructions
├── DESIGN.md             # Design system
├── PRD-Aplikasi-Penilaian-Biola.md  # Requirements
├── IMPLEMENTATION_PLAN.md # 14-phase plan
└── README.md             # This file
```

## Key Decisions

### Architecture
- **Decoupled:** Go backend + Next.js frontend (better separation of concerns)
- **Shared Schema:** All schools in one DB (vs separate DB per tenant) — simpler operations
- **JWT + Blacklist:** Stateless auth with PostgreSQL invalidation (vs Redis)
- **In-Memory Rate Limiting:** Single-instance deployment acceptable (vs Redis)
- **Application-Layer Security:** `school_id` filtering in every query (vs PostgreSQL RLS)

### Tech Stack
- **Go + Fiber:** Best performance/RAM ratio for API
- **sqlc:** Type-safe SQL without ORM overhead
- **NextAuth:** Simplest OAuth + credentials handling
- **Dexie.js:** Ergonomic IndexedDB for offline queue
- **SheetJS:** Client-side Excel export (no server burden)

## Security

### Multi-Tenant Isolation
- `school_id` mandatory in every tenant-table query
- Cross-tenant FK validation (student_id must match school_id)
- JWT claims immutable for non-super_admin
- Super admin requires explicit `X-School-ID` header

### Authentication
- bcrypt cost ≥ 12 for passwords
- JWT valid 7 days (per PRD US-06)
- Token blacklist on logout/deactivation
- Rate limiting on auth endpoints (10 req/min)

### Audit
- All write operations logged to `audit_logs`
- JSONB old_data/new_data tracking
- User + school_id recorded for accountability

## Performance Targets

| Metric | Target |
|--------|--------|
| Assessment entry | < 60 seconds (human task) |
| API response | < 500ms |
| Student list load | < 1 second |
| Offline sync (10 items) | < 3 seconds |
| Report query (100 students) | < 2 seconds |

## Internationalization

**Supported Locales:**
- `en` — English (default)
- `id` — Bahasa Indonesia

**URL Strategy:** Path prefix (`/en/*`, `/id/*`)

**Scope:**
- All UI copy via `next-intl` translations
- KurMer report descriptions always Indonesian (artifact, not UI)
- Excel export filenames Indonesian (school-facing)
- User profile setting: `users.preferred_locale`

## Development

### Code Style
- Go: `gofmt` formatting, no comments unless WHY
- TypeScript: No unrequested abstractions, minimal code
- SQL: Raw only, no ORM, explicit indexed columns
- CSS: Tailwind only, DESIGN.md tokens mandatory

### Testing
- Go: `go test ./...` unit + integration tests
- Frontend: Vitest + React Testing Library (components)
- E2E: Playwright (full user journeys)
- Security: Manual tenant isolation audit (Phase 13)

## Troubleshooting

### Database connection error
```bash
# Check PostgreSQL is healthy
docker-compose logs postgres

# Verify connection string in .env
# Default: postgres://dawai:change_me@postgres:5432/dawai?sslmode=disable
```

### MinIO not accessible
```bash
# Check MinIO is running
docker-compose logs minio

# Access console: http://localhost:9001 (minioadmin/minioadmin)
```

### Frontend can't reach API
```bash
# Verify API_URL_INTERNAL in frontend container
docker-compose logs web

# Should be http://api:8080 (internal Docker network)
```

## Roadmap

### Phase 1-4: Completed ✅
- Database migrations
- Auth endpoints
- Frontend scaffold

### Phase 5-8: In Progress 🟡
- Super admin dashboard
- School admin curriculum
- Student management
- Assessment core (highest priority)

### Phase 9-11: Pending
- Offline sync
- Reporting & KurMer
- Portal (student/parent)

### Phase 12-14: Final
- Audit logging
- Security hardening
- PWA & deployment

## Support

For issues, questions, or PRD clarifications:
- See IMPLEMENTATION_PLAN.md for 14-phase breakdown
- See DESIGN.md for exact UI/color specifications
- See CLAUDE.md for architecture decisions
- See PRD for complete feature spec

---

**Status:** Phase 4 scaffolding complete  
**Next Focus:** Assessment core (Phase 8) — scoring math & mobile panel  
**Built with:** Go, Next.js, PostgreSQL, Tailwind  
**License:** Private
