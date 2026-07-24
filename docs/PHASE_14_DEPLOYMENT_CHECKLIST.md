# Phase 14 Deployment Checklist

**Date:** 2026-07-24  
**Status:** 🔄 In Progress

## Environment Variables

- [x] POSTGRES_PASSWORD in .env
- [x] JWT_SECRET in .env (32+ chars)
- [x] AUTH_SECRET in .env (32+ chars)
- [x] NEXTAUTH_URL in .env
- [x] NEXT_PUBLIC_API_URL in .env
- [ ] MINIO vars configured
- [ ] Google OAuth vars configured (optional)

## Dockerfiles

- [x] backend/Dockerfile uses go install for migrate (arch-agnostic)
- [x] backend/Dockerfile multi-stage build
- [x] frontend/Dockerfile node 20-alpine
- [x] frontend/Dockerfile npm build + npm start
- [ ] All Dockerfiles pass `docker compose build` without errors

## Docker Compose

- [x] postgres service configured
- [x] api service configured + depends_on postgres
- [x] web service configured + depends_on api
- [ ] minio service configured (optional)
- [x] All env vars passed correctly via `environment:`
- [x] Volumes mounted for migrations

## Services Health

- [x] API responds: `curl http://localhost:8080/api/health` → 401 (auth required, OK)
- [x] Frontend responds: `curl http://localhost:3000` → 200
- [x] Both containers running in docker compose
- [ ] All services pass smoke test

## Woodpecker CI/CD

- [x] .woodpecker.yml defines test-backend step
- [x] .woodpecker.yml defines deploy-backend step
- [x] .woodpecker.yml defines deploy-frontend step
- [x] deploy-backend uses `docker compose up -d --build --no-deps api`
- [x] deploy-frontend uses `docker compose up -d --build --no-deps web`
- [x] Secret references: db_password, jwt_secret, auth_secret, smtp_pass, etc.
- [ ] Secrets correctly injected in Woodpecker settings

## Integration

- [ ] Frontend can reach API at NEXT_PUBLIC_API_URL
- [ ] JWT auth flow works end-to-end
- [ ] NextAuth session/JWT integrated
- [ ] Database migrations run on startup

## Issues Found & Fixed

### Blocker #1: Dockerfile migrate binary arch mismatch
- **Issue:** `curl | tar xvz` downloaded wrong architecture
- **Fix:** Changed to `go install github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0`
- **Status:** ✓ Fixed, API builds successfully

### Blocker #2: docker-compose.yml missing web service
- **Issue:** frontend not in docker-compose, only api+postgres
- **Fix:** Added web service with proper env vars + depends_on api
- **Status:** ✓ Fixed, both api + web running

### Blocker #3: AUTH_SECRET missing in .env
- **Issue:** docker-compose env interpolation failed on AUTH_SECRET
- **Fix:** Added AUTH_SECRET to .env
- **Status:** ✓ Fixed

## Final Status

**Deployment ready:** Pending teammate verification  
**Next:** Lighthouse audit, E2E tests, offline validation

