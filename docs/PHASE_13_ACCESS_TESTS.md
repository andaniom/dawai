# Phase 13: Access Control Tester
**Cross-Tenant Access Verification for DAWAI**

Status: IN PROGRESS  
Date: 2026-07-24  
Objective: Validate multi-tenant isolation + role-based access control across all endpoints.

---

## Executive Summary

Phase 13 verifies that DAWAI's multi-tenant architecture properly isolates tenant data and enforces role-based access control. Security tests validate:

1. **Cross-tenant access is blocked** — users from School A cannot read/write School B data
2. **Role-based access is enforced** — students cannot perform admin/teacher actions
3. **JWT claims are validated** — expired/tampered tokens are rejected
4. **School ID extraction is secure** — always from JWT, never from user input

---

## Test Environment

### Deployment
```bash
docker compose up -d postgres api
# postgres: localhost:5432
# api: localhost:8080
```

### Database Seeding
Two schools created with complete hierarchies:
- **School A**: admin-a, teacher-a, student-a, Violin subject
- **School B**: admin-b, teacher-b, student-b, Piano subject

Seeded by `seed()` function in `backend/internal/handlers/isolation_test.go`.

### Running Tests
```bash
cd backend

# Run all isolation tests (requires live DB)
DAWAI_SKIP_DB_TESTS=0 go test ./internal/handlers -v -run "Isolation"

# Run single test
DAWAI_SKIP_DB_TESTS=0 go test ./internal/handlers -v -run "TestCrossTenantStudentFetchForbidden"
```

---

## Test Matrix

### Category A: Cross-Tenant Data Isolation

| Test ID | Scenario | Expected | Status |
|---------|----------|----------|--------|
| A1 | School A user lists students → sees only School A | 200 + 1 student | [ ] |
| A2 | School A user GETs School B student → 404 | 404 + not_found error | [ ] |
| A3 | School A admin lists subjects → sees only School A | 200 + 1 subject | [ ] |
| A4 | School A teacher creates assessment with School B student → rejected | 403 + validation_error | [ ] |
| A5 | School A user lists assessments → sees only School A | 200 + filtered results | [ ] |

**Root Causes Blocked:**
- School ID extracted from JWT `school_id` claim ✅
- Service layer validates entity school_id matches request context
- No unguarded table scans (`WHERE school_id = $1` enforced)

---

### Category B: Role-Based Access Control

| Test ID | Scenario | Expected | Status |
|---------|----------|----------|--------|
| B1 | Student accesses /api/assessments (POST) | 403 or validation failure | [ ] |
| B2 | Student accesses /api/students (POST) | 403 or validation failure | [ ] |
| B3 | Teacher accesses /api/users (POST, admin-only) | 403 + authorization_error | [ ] |
| B4 | Student accesses /api/me (own data) | 200 + own user info | [ ] |
| B5 | Parent accesses child's data via parent_students | 200 if child linked | [ ] |

**Enforcement Points:**
- `RoleGuard("school_admin")` on user management endpoints
- Service-layer business logic (e.g., only teacher can submit assessment)
- No role = no tenant access (TenantGuard requires school_id in JWT)

---

### Category C: JWT Validation

| Test ID | Scenario | Expected | Status |
|---------|----------|----------|--------|
| C1 | Missing Authorization header | 401 + auth_error | [ ] |
| C2 | Invalid Bearer format (no "Bearer" prefix) | 401 + auth_error | [ ] |
| C3 | Tampered token (modify school_id claim) | 401 + invalid token | [ ] |
| C4 | Expired token | 401 + invalid token | [ ] |
| C5 | Blacklisted token (after logout) | 401 + invalid token | [ ] |

**Validation Stack:**
- JWTGuard: parses token, verifies signature with `JWT_SECRET`, extracts claims
- TenantGuard: ensures school_id present + non-empty
- No skip: all /api/* routes protected except /api/auth/login

---

### Category D: Super-Admin Impersonation (x-school-id header)

| Test ID | Scenario | Expected | Status |
|---------|----------|----------|--------|
| D1 | Non-admin uses x-school-id header | ignored, uses own school_id | [ ] |
| D2 | Super-admin impersonates School A | can see/manage School A data | [ ] |
| D3 | Super-admin cannot impersonate invalid school | 404 or validation error | [ ] |

**Note:** Super-admin endpoints (`/api/super-admin/*`) not yet implemented; placeholder for Phase 14.

---

## Endpoint Coverage

### Authentication (`/api/auth/`)
- POST `/api/auth/login` — public, no JWT required
- POST `/api/auth/logout` — protected, JWT required

### Users (`/api/users/`) — school_admin only
- POST `/api/users` — create user in school
- GET `/api/users` — list users in school

### Students (`/api/students/`)
- GET `/api/students` — list students, school scoped
- GET `/api/students/:id` — fetch by ID, cross-tenant check at line 77 student.go
- POST `/api/students` — create student record

### Subjects (`/api/subjects/`) — school scoped
- GET `/api/subjects` — list subjects
- POST `/api/subjects` — create subject
- GET `/api/subjects/:id/rubric` — fetch rubric components

### Assessments (`/api/assessments/`) — teacher primary actor
- POST `/api/assessments` — create assessment (implicit teacher check via business logic)
- GET `/api/assessments` — list assessments
- GET `/api/assessments/:id` — fetch assessment
- PATCH `/api/assessments/:id` — update assessment
- DELETE `/api/assessments/:id` — delete assessment

---

## Key Security Patterns

### 1. School ID Extraction (CRITICAL)
```go
// ✅ CORRECT — from JWT claim
func GetStudents(c *fiber.Ctx) error {
    schoolID := c.Locals("school_id").(string)  // from JWTGuard
    // Query with school_id = $1
}

// ❌ WRONG — from request input
func GetStudents(c *fiber.Ctx) error {
    schoolID := c.Query("school_id")  // SECURITY BUG
    // User can inject another school_id
}
```

### 2. Cross-Tenant Validation
```go
// Validate entity belongs to same school
if uuidToString(row.SchoolID) != schoolID {
    return nil, errors.New("student not found")  // 404, not "cross-tenant"
}
```

### 3. Role Enforcement
```go
// Middleware checks JWT roles array
api.Group("/users", middleware.RoleGuard("school_admin"))
```

---

## Running Live Tests

### Prerequisites
```bash
# Terminal 1: Start services
docker compose up -d

# Terminal 2: Run tests
cd backend
DAWAI_SKIP_DB_TESTS=0 go test ./internal/handlers -v -run "Isolation"
```

### Expected Output
```
=== RUN   TestStudentListIsolation
    isolation_test.go:280: school A teacher list students should succeed
    isolation_test.go:281: school A should have exactly 1 student
--- PASS: TestStudentListIsolation (0.150s)

=== RUN   TestCrossTenantStudentFetchForbidden
    isolation_test.go:308: accessing school B student from school A should return 404
--- PASS: TestCrossTenantStudentFetchForbidden (0.120s)

=== RUN   TestCrossTenantAssessmentCreateRejected
    isolation_test.go:363: cross-tenant assessment create should be forbidden
--- PASS: TestCrossTenantAssessmentCreateRejected (0.200s)

ok  	github.com/violin-assessment/dawai/internal/handlers	0.470s
```

---

## Manual Testing Checklist

### Step 1: Seed Data
```bash
# Via docker psql
docker compose exec postgres psql -U dawai -d dawai

-- Verify schools exist
SELECT id, name FROM schools;

-- Verify user roles
SELECT u.email, r.name, ur.school_id FROM users u
  JOIN user_roles ur ON ur.user_id = u.id
  JOIN roles r ON r.id = ur.role_id
  ORDER BY u.email;
```

### Step 2: Login & Get Token
```bash
# School A teacher login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"teacher-a@school-a.local", "password":"..."}'

# Response: { access_token: "eyJ...", user: {..., school_id: "..."} }
TEACHER_A_TOKEN="eyJ..."
SCHOOL_A_ID="..."
```

### Step 3: Cross-Tenant Access Test
```bash
# School A teacher tries to fetch School B student
curl -X GET http://localhost:8080/api/students/<student-b-id> \
  -H "Authorization: Bearer $TEACHER_A_TOKEN"

# Expected: 404 + "student not found"
```

### Step 4: Role Boundary Test
```bash
# Student tries to create assessment (teacher-only)
curl -X POST http://localhost:8080/api/assessments \
  -H "Authorization: Bearer $STUDENT_A_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{...}'

# Expected: 403 or validation error
```

---

## Findings & Remediation

### Critical Issues Found
_(To be updated during testing)_

- [ ] Issue #1: [description] → Requires [fix]
- [ ] Issue #2: [description] → Requires [fix]

### Code Quality Checks
- [x] All handlers extract school_id from c.Locals("school_id")
- [x] Services validate cross-tenant entity ownership
- [x] RoleGuard applied to protected endpoints
- [x] JWTGuard + TenantGuard on all /api/* except /api/auth/login
- [ ] Role enforcement consistent across all write operations

---

## Next Steps

1. **Run automated tests** — `DAWAI_SKIP_DB_TESTS=0 go test ./internal/handlers -v -run "Isolation"`
2. **Manual spot-check** — login as different roles, verify access boundaries
3. **Document failures** — record test ID, actual result, root cause
4. **Patch issues** — fix at handler/service/middleware layer
5. **Regression verify** — re-run full test suite after fixes
6. **Sign off Phase 13** — update status to COMPLETE

---

## Test Code Reference

All tests implemented in: `backend/internal/handlers/isolation_test.go`

**Key Functions:**
- `seed(t *testing.T, pool *pgxpool.Pool) *TestData` — creates two schools + users + subjects
- `appWithAuth(pool, schoolID, userID, roles)` — builds Fiber app with injected locals
- `doRequest(app, method, path, body)` — executes HTTP request, returns (status, response)

**Test Functions:**
- `TestStudentListIsolation` — A1
- `TestCrossTenantStudentFetchForbidden` — A2
- `TestCrossTenantSubjectListIsolation` — A3
- `TestCrossTenantAssessmentCreateRejected` — A4
- `TestStudentCannotAccessTeacherEndpoint` — B1
- `TestHarnessBuilds` — compile-time sanity check

---

## Approval

- [x] Security model reviewed (JWT + school_id isolation)
- [ ] All tests passing
- [ ] Manual verification complete
- [ ] Phase 13 signed off
