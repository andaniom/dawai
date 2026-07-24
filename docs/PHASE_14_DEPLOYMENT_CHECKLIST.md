# Phase 14 Deployment Checklist

## Environment Variables

| Var | Required | Status | Notes |
|---|---|---|---|
| `POSTGRES_PASSWORD` | ✅ | ✅ PASS | Present in .env |
| `JWT_SECRET` | ✅ | ✅ PASS | Present in .env (32+ chars) |
| `AUTH_SECRET` | ✅ | ✅ PASS | Present in .env (32+ chars) |
| `NEXTAUTH_URL` | ✅ | ❌ FAIL | Missing from .env — required for NextAuth |
| `NEXT_PUBLIC_API_URL` | ✅ | ❌ FAIL | Missing from .env — required for frontend to reach API |
| `POSTGRES_USER` | ✅ | ⚠️  PARTIAL | Not in .env; hardcoded in docker-compose.yml as `dawai` |
| `POSTGRES_DB` | ✅ | ⚠️  PARTIAL | Not in .env; hardcoded in docker-compose.yml as `dawai` |
| `API_URL_INTERNAL` | ✅ | ✅ PASS | In docker-compose, defaults to `http://api:8080` |

**Action:** Add missing vars to .env:
```bash
NEXTAUTH_URL=http://localhost:3000
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

## Docker Build Tests

| Service | Build | Status | Notes |
|---|---|---|---|
| Backend (api) | `docker compose build api` | ✅ PASS | Multi-stage, golang:1.25-alpine → alpine |
| Frontend (web) | `docker compose build web` | ✅ PASS | Build succeeds; non-critical Next.js metadata warnings |

---

## Service Connectivity

| Endpoint | Expected | Actual | Status |
|---|---|---|---|
| `http://localhost:8080/api/health` | 401 Unauthorized (auth required) | Auth error (expected) | ✅ PASS |
| `http://localhost:3000` | 200 + login page | 200 + default template | ❌ FAIL |

**Bug:** frontend/app/page.tsx is still default Create Next App template, not DAWAI login page.

---

## Docker Compose Status

| Service | Container | Status |
|---|---|---|
| PostgreSQL | Running | ✅ Port 127.0.0.1:5432 |
| API (Go) | Running | ✅ Port 127.0.0.1:8080 |
| Frontend (Next.js) | Running | ✅ Port 127.0.0.1:3000 |

✅ All services running.

---

## Dockerfile Validation

### Backend ✅ PASS
- Multi-stage build (golang:1.25-alpine → alpine)
- Installs golang-migrate in builder
- Runs migrations before starting API

### Frontend ⚠️ PARTIAL
- Build succeeds but `NEXT_PUBLIC_API_URL` should be build arg, not runtime env
- Update docker-compose.yml to pass as build args

---

## CI/CD Pipeline (`.woodpecker.yml`)

| Step | Status | Issue |
|---|---|---|
| test-backend | ✅ PASS | - |
| deploy-backend | ✅ PASS | - |
| deploy-frontend | ⚠️ PARTIAL | Env vars passed at runtime, not build time |
| test-frontend | ❌ MISSING | No tests |

---

## Critical Issues

| ID | Issue | Fix |
|---|---|---|
| ENV-001 | Missing NEXTAUTH_URL | Add to .env |
| ENV-002 | Missing NEXT_PUBLIC_API_URL | Add to .env |
| FE-001 | Root page is default template | Redirect to /login |
| DOCKER-001 | No build args for frontend | Add args: section in docker-compose.yml |

---

Generated: 2026-07-24
