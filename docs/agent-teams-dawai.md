# Agent Teams for DAWAI Development

Tailored use cases and patterns for DAWAI's multi-tenant architecture, security requirements, and phased development.

---

## Key Context for DAWAI Teams

### Architecture Constraints

- **Multi-tenant**: Every endpoint must filter by `school_id` (application layer)
- **Shared database**: No PostgreSQL RLS, isolation via app code only
- **Security-critical**: Data leak of one school into another = critical bug
- **JWT + RBAC**: Roles stored in token, validated per-endpoint
- **Offline-first**: Students can work offline, sync when connected

### Team Communication Guidelines

- Always mention `school_id` filtering when reviewing code
- Cross-tenant access tests are non-negotiable
- Every merge requires isolation verification
- Audit logs must be immutable + complete

---

## Phase 13: Security Hardening Sprint

Hardening + compliance audit. 3-5 days. Use full team approach.

### Launch Command

```text
DAWAI Phase 13: Security hardening sprint.

Spawn 4 teammates to conduct comprehensive security audit:

1. Query auditor (Andani's audit):
   Call this teammate: "query-auditor"
   Task: grep all backend Go handlers for database queries.
   For each query, verify school_id filter is present and correct.
   Create spreadsheet at docs/PHASE_13_QUERY_AUDIT.md:
   | File | Function | Query | school_id Filter | Status |
   Report: which queries missing school_id filters (critical findings first).
   Own files: docs/PHASE_13_QUERY_AUDIT.md (read-only grep analysis)

2. Access control tester:
   Call this teammate: "access-tester"
   Task: test each endpoint for cross-tenant access violations.
   Test matrix:
   - User from School A accessing data from School B (must fail)
   - Teacher endpoints accessed by student (must fail)
   - Admin endpoints accessed by non-admin (must fail)
   - JWT token tampering (change school_id claim, must fail)
   Create report at docs/PHASE_13_ACCESS_TESTS.md with pass/fail for each test.
   Own files: docs/PHASE_13_ACCESS_TESTS.md

3. Test coverage builder:
   Call this teammate: "test-writer"
   Task: write Go tests for cross-tenant access scenarios.
   For each endpoint:
   - Write test: user A from School A cannot read School B data
   - Write test: JWT school_id mismatch is rejected
   - Write test: audit log records access attempt
   Create tests file: backend/internal/handlers/cross_tenant_access_test.go
   Run tests, report: coverage % + pass/fail status.
   Own files: backend/internal/handlers/cross_tenant_access_test.go

4. Documentation writer:
   Call this teammate: "doc-writer"
   Task: document security model for Phase 13 closure.
   Write: docs/SECURITY_MODEL.md covering:
   - Auth flow: OAuth → NextAuth → Go /api/auth/token → JWT
   - Isolation model: school_id in JWT claims + application filtering
   - Role model: roles array in JWT + RoleGuard middleware
   - Cross-tenant testing: how isolation is verified
   - Audit trail: immutable write logging
   Add checklist: SECURITY_CHECKLIST.md (for Phase 13 signoff)
   Own files: docs/SECURITY_MODEL.md, docs/SECURITY_CHECKLIST.md

Coordinate:
- Query auditor: report findings immediately, others use to guide tests
- Test writer: wait for query auditor findings, write tests for identified gaps
- Access tester: test the gaps reported by query auditor
- Doc writer: consolidate findings into compliance docs

Timeline: 2-3 days parallel work, 1 day review + fixes.

Success criteria:
- All queries have school_id filters (or documented exception + reason)
- Cross-tenant access tests all pass
- Audit trail complete for all write operations
- Phase 13 checklist signed off
```

### After Teammates Finish

```text
Review all Phase 13 outputs:
1. docs/PHASE_13_QUERY_AUDIT.md — any queries missing filters?
2. docs/PHASE_13_ACCESS_TESTS.md — any failed tests?
3. backend/internal/handlers/cross_tenant_access_test.go — run locally
4. docs/SECURITY_MODEL.md — accurate description?
5. docs/SECURITY_CHECKLIST.md — all items satisfied?

If any gaps:
- Query auditor: add missing filter + test
- Access tester: re-test gap, confirm fix
- Test writer: add test for gap
- Doc writer: update checklist status

Repeat until all pass.
```

---

## Phase 14: PWA + Deployment

Split PWA audit and deployment preparation across team.

### Launch Command

```text
DAWAI Phase 14: PWA audit + deployment preparation.

Spawn 3 teammates:

1. PWA auditor:
   Call this teammate: "pwa-auditor"
   Task: run Lighthouse PWA audit on frontend.
   npm run build && npm start
   Open http://localhost:3000 in browser
   Run Lighthouse PWA audit, report:
   - Installability score + issues
   - Performance score + bottlenecks
   - Best practices violations
   - Accessibility violations (WCAG)
   Create report: docs/PHASE_14_LIGHTHOUSE_REPORT.md
   Recommendations for each issue.

2. Deployment readiness checker:
   Call this teammate: "deployment-checker"
   Task: verify .env variables, secrets, CI/CD pipeline.
   Checklist:
   - .env template matches docker-compose.yml? (DATABASE_URL, JWT_SECRET, etc)
   - Woodpecker pipeline (.woodpecker.yml) runs on main branch
   - Secret injection working (db_password, jwt_secret, etc)
   - Build validation: docker build backend, docker build frontend success?
   - Health checks: GET /api/health returns 200?
   Create checklist: docs/PHASE_14_DEPLOYMENT_CHECKLIST.md (pass/fail)

3. Smoke test builder:
   Call this teammate: "smoke-tester"
   Task: write E2E smoke tests for critical flows.
   Test paths:
   - Login flow: email/password or OAuth
   - Create assessment: fill form, submit, verify DB save
   - View report: KurMer generation, Excel download
   - Student portal: view own assessments
   Create file: backend/tests/e2e_smoke_test.go
   Run against docker compose stack, report: pass/fail for each flow.

Coordinate:
- PWA auditor: fix issues found, re-audit until score good
- Deployment checker: ensure all env vars + secrets ready before deploy
- Smoke tester: run against production-like docker compose config

Timeline: 2-3 days
Success: Lighthouse scores > 90, all smoke tests pass, deployment checklist complete
```

---

## Feature Development: New Rubric Components

Team development of new rubric scoring feature.

### Launch Command

```text
New rubric component feature. Parallel implementation:

Spawn 3 teammates:

1. Database + API teammate (call: "backend-dev"):
   Backend structure + queries.
   Create:
   - Migration: 000X_add_rubric_components_v2.up.sql (new fields)
   - Model: RubricComponentV2 struct
   - Queries: backend/queries/rubric.sql (CRUD + filtering)
   - Service: backend/internal/services/rubric_service.go
   - Handler: backend/internal/handlers/rubric_handler.go (POST/PUT/GET)
   Test: POST /api/schools/{schoolId}/rubric-components with JWT + school_id
   Own files: migrations/*, queries/rubric.sql, backend/internal/{services,handlers}/rubric*
   No frontend code.

2. Frontend UI teammate (call: "frontend-dev"):
   Components + UI.
   Create:
   - Component: components/RubricComponentForm.tsx
   - Page: app/schools/[id]/rubric/[componentId]/page.tsx
   - State: lib/hooks/useRubricComponents.ts (TanStack Query)
   - Styles: Tailwind + shadcn/ui
   Test: form renders, can fill + submit (API call fails if backend not ready, that's OK)
   Own files: frontend/components/Rubric*, frontend/app/rubric*, frontend/lib/hooks/*
   No backend or database code.

3. Test teammate (call: "test-dev"):
   Wait for backend, then test.
   Create:
   - Unit: backend/internal/services/rubric_service_test.go
   - Integration: backend/internal/handlers/rubric_handler_test.go
   - E2E: backend/tests/rubric_e2e_test.go
   Test: create rubric component as teacher, verify school_id isolation
   Run: go test ./..., report coverage + pass/fail

Workflow:
- Backend: schema done by day 1 PM
- Frontend: waiting for schema, starts day 1 evening
- Tests: waiting for backend, starts day 2 AM
- Lead: coordinate, unblock, integrate by day 2 PM

Success: all tests pass, cross-tenant access fails, code review approved
```

---

## Bug Investigation: Multi-School Data Leak

When isolation bug is suspected.

### Launch Command

```text
Suspected data leak: School A user reporting seeing School B data.

Spawn 4 teammates to investigate competing theories:

1. Query leak hypothesis:
   Teammate: "query-investigator"
   Theory: A SELECT query is missing WHERE school_id = ?
   Investigate:
   - Which endpoint did user hit? (check audit logs)
   - What was in response? (check response logs)
   - Which query generated that response? (check backend code)
   - Does query filter by school_id? (grep + code review)
   Report: which query is leaking, fix location

2. JWT tamper hypothesis:
   Teammate: "jwt-investigator"
   Theory: User modified their JWT to change school_id
   Investigate:
   - Was JWT signature valid? (check JWTGuard validation)
   - Was token expired? (check expiry time)
   - Did claims match DB data? (check user.school_id in DB)
   - Is JTI blacklisted? (check jwt_blacklist table)
   Report: JWT valid or tampered, how it leaked

3. Caching hypothesis:
   Teammate: "cache-investigator"
   Theory: Stale cache returning old data
   Investigate:
   - Is caching enabled? (check middleware)
   - Is cache keyed by school_id? (check cache key generation)
   - Was cache invalidated? (check invalidation logic)
   - Can we reproduce offline? (check with cache disabled)
   Report: cache issue or not

4. Audit trail analyzer:
   Teammate: "audit-investigator"
   Theory: audit_logs table shows anomalies
   Investigate:
   - Was data accessed? (check audit_logs for read)
   - By which user? (check user_id + school_id in audit_logs)
   - What time? (check timestamp, match user's session)
   - Was it written to? (check INSERT/UPDATE/DELETE logs)
   Report: audit trail shows what was leaked + when

Coordinate:
- All investigate simultaneously
- Query investigator: if found missing filter → CRITICAL, all stop
- JWT investigator: if JWT tamper found → security audit needed
- Cache investigator: if cache issue → purge + disable
- Audit investigator: consolidates timeline

Lead: once consensus emerges, implement fix + add regression test
```

---

## Code Review: Assessment Flow Refactor

Multi-angle review of complex business logic change.

### Launch Command

```text
Review assessment flow refactor (PR #<number>).

Spawn 4 reviewers, each focused:

1. Data integrity reviewer:
   Focus: database consistency, transaction handling, race conditions
   Check:
   - Assessment creation: is it atomic? (all components created or none)
   - Update flow: is state consistent? (can't have partial update)
   - Soft delete: is audit trail preserved?
   - Offline sync: conflict resolution logic correct?
   Report: data integrity issues

2. Isolation reviewer:
   Focus: school_id filtering, cross-tenant access, permission checks
   Check:
   - Every query filters by school_id
   - Assessment from School A cannot be accessed by School B user
   - Teacher can only update own school's assessments
   - Audit log doesn't leak other schools' data
   Report: isolation issues with severity

3. Performance reviewer:
   Focus: query efficiency, N+1 problems, caching
   Check:
   - Query plan: is DB using indexes?
   - Bulk operations: are queries batched?
   - Caching: is response cached if unchanged?
   - Component scoring: is it O(n) or O(n²)?
   Report: performance issues + suggestions

4. Error handling reviewer:
   Focus: what happens when things fail?
   Check:
   - DB connection fails: does handler retry or 500?
   - JWT invalid: does request fail gracefully?
   - Validation fails: is error message helpful?
   - Concurrency timeout: does user see clear message?
   Report: error handling gaps

Each review independently. Report findings with severity + code location.
Lead: synthesize. If critical issues → request changes. If minor → approve.
```

---

## Testing: Offline Sync Scenarios

Multi-angle test of complex sync logic.

### Launch Command

```text
Write comprehensive tests for offline sync.

Spawn 3 teammates:

1. Sync logic tester:
   Unit tests: lib/services/offlineSync.ts
   Test scenarios:
   - Queue 5 assessments, then go online: all sync in order
   - Sync fails for 3rd item: others still sync, retry 3rd on next go
   - Local delete + server modification: conflict resolution
   - Idempotency: resubmit same assessment twice, only 1 saved
   Report: test coverage %, pass/fail

2. Service worker tester:
   Integration tests: service worker behavior
   Test scenarios:
   - Network down: request queued locally
   - Network up: queued requests flushed
   - Background sync: tasks persist across page close/reopen
   - Offline indicator: UI shows "offline" status when no network
   Report: pass/fail for each scenario

3. E2E scenario tester:
   Real browser tests: full user journey
   Test scenarios:
   - User offline, creates 3 assessments, goes online, all sync
   - User online, connection drops mid-sync, reconnect, retry
   - Concurrent submissions from 2 devices, no data loss
   - Offline for 1 hour, then sync with 50+ queued operations
   Report: pass/fail, any data loss

All report: coverage %, which scenarios fail, recommended fixes
Lead: integrate tests, run as part of CI
```

---

## Checklist: Before Spawning Team on DAWAI

- [ ] **Task is parallel-friendly?** (not sequential)
- [ ] **Each teammate isolated?** (no same-file edits)
- [ ] **Security-critical work?** (add isolation verification teammate)
- [ ] **Cross-tenant context mentioned?** (spawn prompt includes "school_id" keyword)
- [ ] **Deadlines reasonable?** (3-5 days typical, not 1 day)
- [ ] **Token budget allows?** (teams are expensive, feature-worthy tasks only)
- [ ] **Clear deliverables?** (each teammate knows done when they see X)
- [ ] **Isolation tests included?** (for any endpoint touching school_id)
- [ ] **Audit trail verified?** (if write operations, check audit_logs table)

---

## Common DAWAI Team Pitfalls

❌ **Forgot to add school_id filter check**: Query reviewer misses WHERE school_id = $1. Use grep sweep pre-review.

❌ **Multiple teammates edit same handler file**: auth_handler.go edited by 2 teams → conflicts. Split: one team per file.

❌ **No isolation test**: Feature ships, later discovered School A can see School B data. Make isolation test mandatory.

❌ **Audit trail incomplete**: Writes recorded but reads not logged. Check audit_logs table exhaustively.

❌ **JWT tampering not caught**: Didn't test invalid school_id claim. Add JWT manipulation test.

---

## Reference: CLAUDE.md Security Constraints

From `CLAUDE.md` § Multi-Tenant Isolation:

> **Rule:** `school_id` MUST be extracted from JWT claims via `c.Locals("school_id")` in Go. Never from request body, query params, or URL path (except super_admin using `x-school-id` header).

**Team implication:** Every code review must verify this rule. Grep for `c.Locals("school_id")` in every handler.

> **Cross-tenant validation:** When referencing entities (e.g., `student_id` in assessment), always verify the entity's `school_id` matches request context.

**Team implication:** Access tester must test: given Student X from School A, can team access other students? Must fail.

---

## Integration with Git Workflow

After teammates finish:

1. Coordinate: collect all PRs from teammates
2. Review: lead reviews each PR per CLAUDE.md
3. Test: run full suite: `go test ./...` + `npm run lint`
4. Merge: merge all PRs to master in order (database schema first if applicable)
5. Deploy: `docker compose up -d` to test
6. Tag: if successful, tag commit with phase completion

Example:

```bash
# After Phase 13 teammates finish
git log --oneline -10  # verify all 4 PRs merged
go test ./...          # all tests pass?
npm run build          # frontend builds clean?
docker compose up -d   # services start?
curl http://localhost:8080/api/health  # backend responds?

# If all pass
git tag phase-13-security-hardening
git push origin phase-13-security-hardening
```
