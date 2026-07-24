# DAWAI Phase 13 Sign-off Checklist

This checklist must be completed and all items checked before Phase 13 (Security Hardening) is considered complete.

## 1. Multi-Tenant Isolation
- [ ] All database queries reading/writing tenant data include `school_id = $1` filter
- [ ] By-ID queries updated to include school_id (GetSubjectByID, GetStudentByID, GetAssessmentByID, GetRubricComponentByID)
- [ ] Cross-tenant reference checks implemented (e.g., creating assessment validates student and subject belong to same school)
- [ ] `isolation_test.go` integration tests pass with 100% coverage of cross-tenant attempts
- [ ] `/api/students` endpoints properly restricted by role (fixing ISOLATION-006)
- [ ] Parent-scoped access validation implemented (parents can only read their linked children)

## 2. Authentication & JWT
- [ ] JWT signature validation enforced (HS256)
- [ ] JWT expiry validated (7-day lifetime)
- [ ] JWT blacklist checking implemented in `JWTGuard` (currently missing)
- [ ] Password hashing verified to use bcrypt with cost ≥ 12
- [ ] Rate limiting implemented on `/api/auth/login` (10 req/min)
- [ ] Rate limiting implemented on forgot/reset password endpoints (if applicable)

## 3. Authorization & RBAC
- [ ] `RoleGuard` middleware applied correctly to all endpoints
- [ ] Super Admin access strictly requires `x-school-id` header for tenant impersonation
- [ ] Users cannot elevate their own roles
- [ ] Multi-school user support designed (or explicitly deferred to v2)

## 4. Audit Trail
- [ ] `audit_logs` table schema verified (immutable, append-only)
- [ ] Write operations (INSERT, UPDATE, DELETE) logged successfully
- [ ] Failed authorization attempts logged
- [ ] User context (`user_id`, `school_id`, IP) accurately captured in logs
- [ ] Old/New data structures correctly serialized in `JSONB` columns

## 5. Transport & Data Protection
- [ ] HTTPS enforcement configured (HSTS, secure cookies)
- [ ] Error messages scrubbed of sensitive information (no school names or IDs leaked in generic 403s)
- [ ] `.env.example` scrubbed of any actual secrets
- [ ] Database connection configured for encryption (if applicable in prod)

## 6. Known Gaps Addressed
- [ ] JWT Blacklist check in middleware
- [ ] Rate limiting on auth endpoints
- [ ] ISOLATION-006 (Student list exposure)

---
*Checklist created: 2026-07-24*
*Status: In Progress*
