export const meta = {
  name: 'phase13-sprint-review',
  description: 'Parallel multi-agent review of Phase 13 auth critical fixes and deployment readiness',
  phases: [
    { title: 'Security Audit', detail: 'Verify all 5 critical auth fixes' },
    { title: 'Performance Review', detail: 'Cache efficiency + blacklist overhead' },
    { title: 'Test Coverage Analysis', detail: 'Scaffold completeness + Phase 1 gaps' },
    { title: 'Deployment Readiness', detail: 'Build clean, migrations applied, docker verified' },
  ],
}

// Run all 4 reviews in parallel (split pane mode)
const results = await parallel([
  () => agent(
    `SECURITY AUDIT: Verify Phase 13 critical auth fixes

Files modified:
- backend/cmd/api/main.go (JWT_SECRET validation)
- backend/internal/services/auth.go (roles + claim guards)
- backend/internal/middleware/jwt.go (algorithm check + type-safety)
- backend/internal/handlers/auth.go (safe jti cast)

Verify each fix:
1. JWT_SECRET: Does startup panic if env var empty?
2. Hardcoded roles: Does GetRoleNamesByUserSchool() query work? Test with 3 role types.
3. Algorithm check: Does keyfunc reject HS512/RS256? Accept HS256?
4. Type-safety: Does nil claim parsing return 401, not panic? Test 5 malformed cases.
5. Test scaffold: Do all 17 test functions compile and import correctly?

Report: ✅/❌ for each fix + any remaining gaps.`,
    { label: 'Security Audit', phase: 'Security Audit' }
  ),
  () => agent(
    `PERFORMANCE REVIEW: Measure overhead from Phase 13 changes

Current state:
- JWT_SECRET fetched via os.Getenv() per request (expensive)
- Blacklist query on every logout (no cleanup job yet)
- Bcrypt called on every login (expected, cost=12)
- Type-safety guards add 1-2 conditionals per JWT parse

Measure/estimate:
1. os.Getenv() cost per request (~20 sec/hour @ 10 QPS) — quantify actual overhead
2. Blacklist table size — current row count, growth rate, cleanup frequency needed
3. Type-safety guard latency — expected <1ms per validation check
4. Cache JWT_SECRET at init — would reduce CPU by ~2% (defer to Medium M1)

Report: Current bottlenecks ranked by impact + recommendations for sprint 2.`,
    { label: 'Performance Review', phase: 'Performance Review' }
  ),
  () => agent(
    `TEST COVERAGE ANALYSIS: Review auth_test.go scaffold

File: backend/internal/services/auth_test.go
Expected: 17 test function signatures across 4 phases

Verify:
1. Phase 1 (Login Flow): 6 functions — valid/invalid/timing/disabled/multi-school/no-school
2. Phase 2 (JWT Validation): 8 functions — valid/expired/tampered/algorithm/confusion/nil/missing/malformed
3. Phase 3 (Cross-Tenant): 5 functions — isolation boundaries + super-admin header
4. Phase 4 (Edge Cases): 8 functions — logout/rate-limit/init panic/expiry/concurrent/timing/empty-field validation

Check:
- Are all 17 function signatures present?
- Do imports compile (testify/assert, fixtures, seed.go)?
- Are test comments clear enough to implement?
- Any gaps in coverage (e.g., LDAP auth, MFA)?

Report: Completeness score (%) + phase-by-phase implementation priority.`,
    { label: 'Test Coverage Analysis', phase: 'Test Coverage Analysis' }
  ),
  () => agent(
    `DEPLOYMENT READINESS: Verify Phase 13 ready for Phase 14 handoff

Checklist:
1. Go backend builds clean: \`go build ./cmd/api\` succeeds?
2. sqlc generated code up-to-date: \`sqlc generate\` no diffs?
3. Migrations reversible: up/down.sql files for all schema changes?
4. Docker image builds: \`docker build -t dawai-backend .\` succeeds?
5. docker-compose verified: All 4 services (postgres, minio, api, frontend) start?
6. Environment setup: .env.example or docker-compose.env complete?
7. No uncommitted changes: \`git status\` clean?
8. CI/CD ready: .woodpecker.yml updated for Phase 13 changes?
9. Secrets configured: JWT_SECRET, DB_PASSWORD, etc. set in Woodpecker?
10. Rollback plan: How to revert Phase 13 if Critical #1–#5 have issues?

Report: ✅/❌ for each item + blockers for Phase 14 (PWA audit, deployment test, offline capability).`,
    { label: 'Deployment Readiness', phase: 'Deployment Readiness' }
  ),
])

// Summarize findings
log('=== PHASE 13 SPRINT REVIEW (MULTI-AGENT, SPLIT PANE) ===')
log(`Security: Complete`)
log(`Performance: Complete`)
log(`Test Coverage: Complete`)
log(`Deployment: Complete`)

return {
  security: results[0],
  performance: results[1],
  tests: results[2],
  deployment: results[3],
}
