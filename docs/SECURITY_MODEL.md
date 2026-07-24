# DAWAI Security Model

**Phase 13** — Complete security architecture documentation for multi-tenant violin assessment platform.

---

## 1. Authentication Flow

### 1.1 OAuth + Credential Path

```
User (Frontend)
  ↓
NextAuth.js (Google OAuth / Email+Password)
  ↓
POST /api/auth/login (Go Backend) ← Credentials only
  ↓
AuthService.Login() ← bcrypt verify
  ↓
JWT generated w/ school_id + roles
  ↓
NextAuth session (httpOnly cookie)
  ↓
Automatic header injection: Authorization: Bearer <JWT>
  ↓
All subsequent API calls (Go)
```

**Key Points:**
- NextAuth runs client-side (Next.js) and manages session
- Go backend receives email/password only; NextAuth handles provider OAuth
- JWT issued by Go only after credential verification
- JTI (JWT ID) generated per login for logout blacklisting

### 1.2 NextAuth Integration

- **File:** `frontend/next.config.js` (embedded NextAuth config)
- **Callback:** POST `/api/auth/login` with `{email, password}`
- **Response:** `{success: true, data: {accessToken, user: {id, email, name, roles, school_id}}}`
- **Session Storage:** httpOnly cookie (secure, httpOnly flags set)
- **Auto Refresh:** Not implemented (7-day expiry; manual re-login required)

---

## 2. JWT Structure & Validation

### 2.1 Token Claims

```go
type CustomClaims struct {
  SchoolID string   `json:"school_id"`     // Multi-tenant root
  Roles    []string `json:"roles"`         // Array: super_admin, school_admin, etc.
  jwt.RegisteredClaims {
    Subject:  userID                       // User UUID
    ID:       jti                          // Blacklist key
    IssuedAt: now
    ExpiresAt: now + 7 days                // 7-day lifetime
  }
}
```

### 2.2 Validation Middleware (JWTGuard)

**Location:** `backend/internal/middleware/jwt.go`

```go
// Validates:
// 1. Authorization header present ("Bearer <token>")
// 2. Signature valid (HS256, JWT_SECRET env var)
// 3. Token not expired
// 4. Claims parseable

// Injects into c.Locals:
// - user_id (Subject)
// - school_id (CustomClaims.SchoolID)
// - roles (CustomClaims.Roles)
// - jti (for logout)

// Returns 401 if any check fails
```

**Security Properties:**
- ✅ Signature verified (prevents tampering)
- ✅ Expiry checked (7 days)
- ✅ No RLS needed; school_id extracted from JWT, not from body/query/path
- ❌ No JWT blacklist check *yet* (logout blacklist table exists, not validated on every request — **TO FIX**)

### 2.3 Logout & Token Revocation

**Endpoint:** POST `/api/auth/logout`

```go
// Extracts jti from c.Locals("jti")
// Calls AuthService.Logout(ctx, jti)
// Inserts into jwt_blacklist(jti, expires_at)
// Returns 200

// Problem: next request with same JWT still accepted
// Fix: JWTGuard must check blacklist before validating claims
```

**Database Table:**
```sql
CREATE TABLE jwt_blacklist (
  jti VARCHAR(255) PRIMARY KEY,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);
```

---

## 3. Multi-Tenant Isolation

### 3.1 Isolation Model

**Approach:** Application-layer filtering via `school_id` column.

**Why Not RLS (Row Level Security)?**
- Simpler debugging (explicit filters in code, not hidden in policy)
- Easier cross-tenant queries for super_admin (override header, not policy bypass)
- Performance (no policy evaluation per row)
- Trade-off: relies on developer discipline to filter every query

### 3.2 Tenant Guard (TenantGuard)

**Location:** `backend/internal/middleware/tenant.go`

```go
// Validates:
// 1. c.Locals("school_id") exists (set by JWTGuard)
// 2. Not empty string

// Returns 403 if missing/empty
```

**Applied to:** All protected routes via `api.Group("/api", middleware.JWTGuard, middleware.TenantGuard)`

### 3.3 Query Filters (The Core)

**CRITICAL RULE:** Every query MUST filter by `school_id` from `c.Locals("school_id")`, never from request body/query/path.

**Example (✅ Correct):**
```go
func (h *StudentHandler) ListStudents(c *fiber.Ctx) error {
  schoolID := c.Locals("school_id").(string)
  
  students, err := h.studentService.ListStudents(c.Context(), schoolID)
  // Query: SELECT * FROM students WHERE school_id = $1
  // school_id bound from JWT, not user input
}
```

**Anti-Pattern (❌ Wrong):**
```go
func (h *StudentHandler) ListStudents(c *fiber.Ctx) error {
  // Query: SELECT * FROM students
  // No filter → returns all students from all schools
}
```

### 3.4 Cross-Tenant Reference Validation

When a request references an entity (e.g., `student_id` in assessment creation), verify the entity belongs to the request's school.

**Pattern:**
```go
func (h *AssessmentHandler) Create(c *fiber.Ctx) error {
  schoolID := c.Locals("school_id").(string)
  
  var req CreateAssessmentReq
  c.BodyParser(&req)
  
  // Validate student belongs to this school
  student, err := h.service.GetStudent(c.Context(), req.StudentID)
  if err != nil || student.SchoolID != schoolID {
    return fiber.NewError(403, "Student not found in your school")
  }
  
  // Validate subject belongs to this school
  subject, err := h.service.GetSubject(c.Context(), req.SubjectID)
  if err != nil || subject.SchoolID != schoolID {
    return fiber.NewError(403, "Subject not found in your school")
  }
  
  // Now safe to create assessment
}
```

### 3.5 Affected Tables

All tables with `school_id` column must filter by it:

- **students** — index: `idx_students_school_id`
- **subjects** — index: `idx_subjects_school_id`
- **assessments** — index: `idx_assessments_school_id`
- **rubric_components** — inherits via subject (must check subject.school_id)
- **user_roles** — index: `idx_user_roles_school_id`

**Note:** `users` table has NO `school_id` (users exist at platform level); isolation enforced via `user_roles` junction.

---

## 4. Role-Based Access Control (RBAC)

### 4.1 Role Definitions

| Role | Scope | Endpoint Access | Notes |
|------|-------|---|---|
| `super_admin` | Platform | Create/manage schools; impersonate via `x-school-id` header | Can read/write all schools |
| `school_admin` | 1 school | Create users, manage rubric, assessments | School-scoped |
| `teacher` | 1 school | Submit/edit assessments, read students | Can see students in own school |
| `student` | Self | Read own data + assessments | Can only access own student record |
| `parent` | 1 child | Read child's data via `parent_students` | Can only access linked student |

### 4.2 Role Guard (RoleGuard)

**Location:** `backend/internal/middleware/role.go`

```go
func RoleGuard(requiredRole string) fiber.Handler {
  return func(c *fiber.Ctx) error {
    roles := c.Locals("roles").([]string)
    
    // Match any role in user's roles array
    for _, role := range roles {
      if role == requiredRole || role == "super_admin" {
        return c.Next()
      }
    }
    
    return fiber.NewError(403, "Insufficient permissions")
  }
}
```

**Applied as:**
```go
usersGroup := api.Group("/users", middleware.RoleGuard("school_admin"))
usersGroup.Post("", userHandler.CreateUser)
```

### 4.3 User Roles Storage

**Tables:**
```sql
-- Store role assignment
CREATE TABLE user_roles (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  role_id UUID NOT NULL,
  school_id UUID NOT NULL,    -- Role scoped to school
  UNIQUE(user_id, role_id, school_id)  -- 1 role per user per school
);

-- Store school membership (convenience)
CREATE TABLE user_schools (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  school_id UUID NOT NULL,
  UNIQUE(user_id, school_id)
);
```

### 4.4 Multi-School Users

**Current Implementation:** Single school per JWT token (ponytail note in code).

```go
// GetMe/Login queries GetUserSchools, returns first school only
schoolID = schools[0].String()  // ← hardcoded to first
```

**Future:** To support users with roles in multiple schools:
1. Extend JWT to include array of `school_id` + `role_per_school`
2. Require explicit school selection (header: `x-school-id`) on request
3. Validate selected school in user's allowed schools

---

## 5. Password Security

### 5.1 Hashing Algorithm

**Library:** `golang.org/x/crypto/bcrypt`

**Cost:** 12 (OWASP standard for 2024)

```go
hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
```

**Verification:**
```go
err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// Returns nil if match, error if no match
```

### 5.2 Password Requirements

**Current:** No minimum length or complexity (accept any non-empty password).

**Recommended for production:**
- Minimum 8 characters
- No obvious patterns (e.g., "password", "123456")
- No reuse of last N passwords
- Expiry reminder every 90 days

---

## 6. Audit Logging

### 6.1 Scope

**Logged Operations:**
- User creation, role assignment, deletion
- Student creation, modification
- Assessment creation, update, deletion (with score deltas)
- Subject/rubric changes
- Failed access attempts (401, 403 responses)
- Logout events

**Not Logged:**
- Read-only API calls (GET, list)
- Successful auth (implicit in JWT issuance)

### 6.2 Audit Trail Table

```sql
CREATE TABLE audit_logs (
  id UUID PRIMARY KEY,
  school_id UUID NOT NULL,      -- Tenant scoped
  user_id UUID NOT NULL,        -- Who performed action
  action VARCHAR(50) NOT NULL,  -- create_student, update_assessment
  entity_type VARCHAR(50),      -- students, assessments
  entity_id UUID,               -- ID of affected resource
  old_data JSONB,              -- Before state (null for create)
  new_data JSONB,              -- After state (null for delete)
  created_at TIMESTAMP,
  
  INDEX idx_audit_school_id ON audit_logs(school_id)
  INDEX idx_audit_action ON audit_logs(action)
);
```

### 6.3 Immutability

- Audit logs are **append-only** (no UPDATE, no DELETE)
- Retained for 7 years (or per local law)
- Encrypted at rest if sensitive (depends on deployment)

### 6.4 Integration Pattern

**Handlers call AuditLog after every write:**

```go
func (h *StudentHandler) CreateStudent(c *fiber.Ctx) error {
  schoolID := c.Locals("school_id").(string)
  userID := c.Locals("user_id").(string)
  
  // Create student
  student, err := h.service.CreateStudent(ctx, schoolID, req)
  
  // Log
  auditLog(ctx, schoolID, userID, "create_student", "students", 
           student.ID, nil, student)
}
```

---

## 7. Cross-Tenant Testing

### 7.1 Test Scenarios

**Setup:** Seed test DB with:
- School A (schoolA UUID)
- School B (schoolB UUID)
- User UA in A (teacher role)
- User UB in B (teacher role)
- Student SA in A
- Student SB in B
- Subject SubA in A
- Subject SubB in B

**Test Cases:**

1. **Student List Isolation**
   - UA calls `GET /api/students` → sees only SA
   - UB calls `GET /api/students` → sees only SB
   - Cross-check response bodies don't leak other school's IDs

2. **Cross-Tenant Student Fetch**
   - UA calls `GET /api/students/SB` → 403 Forbidden
   - Verify error doesn't leak that SB exists

3. **Cross-Tenant Assessment Create**
   - UA calls `POST /api/assessments` with student_id=SB → 403
   - Verify student ownership validated

4. **Subject List Isolation**
   - UA calls `GET /api/subjects` → sees only SubA
   - UB calls `GET /api/subjects` → sees only SubB

5. **JWT Tampering**
   - Modify JWT `school_id` claim → signature invalid → 401
   - Verify signature cannot be forged without JWT_SECRET

6. **Logout & Reuse**
   - UA logs out → JTI blacklisted
   - UA reuses same JWT → 401 (if blacklist check implemented) OR 200 (bug)

7. **Role Bypass**
   - Student calls `POST /api/users` (admin-only) → 403
   - Verify RoleGuard enforced

8. **Super Admin Header**
   - super_admin calls `POST /api/super-admin/schools` with `x-school-id: schoolA`
   - Creates/reads schoolA data
   - Regular user calls with `x-school-id: other_school` → ignored (uses JWT school_id)

---

## 8. Error Messages & Information Disclosure

### 8.1 Safe Error Responses

**Pattern:** Return generic errors that don't reveal data or business logic.

**✅ Good:**
```json
{
  "success": false,
  "code": 403,
  "error": {
    "message": "Student not found in your school",
    "type": "validation_error"
  }
}
```

**❌ Bad:**
```json
{
  "success": false,
  "code": 403,
  "error": "Student 12345 (Bob Smith) belongs to School B, not your school A"
}
```
→ Leaks existence, name, and school affiliation of other school's data.

### 8.2 Auth Errors

All auth failures return the same generic message:

```go
// ❌ Don't distinguish
"User not found"  vs  "Password incorrect"

// ✅ Do return
"Invalid email or password"
```

---

## 9. Rate Limiting

### 9.1 Current State

**Not Implemented.** Required for production.

### 9.2 Recommended Limits

| Endpoint | Limit | Window |
|---|---|---|
| `POST /api/auth/login` | 10 req/min | Per IP |
| `POST /api/auth/forgot-password` | 5 req/hour | Per email |
| `POST /api/auth/reset-password` | 10 req/day | Per token |
| All other endpoints | 100 req/min | Per user (via JWT) |

**Implementation:** Fiber middleware + Redis (or in-memory for dev).

---

## 10. HTTPS & Transport Security

### 10.1 Development

- Docker Compose runs unencrypted (HTTP on port 8080)
- localhost only, no internet exposure

### 10.2 Production

**Required:**
- HTTPS only (redirect HTTP → HTTPS)
- TLS 1.2 minimum
- Valid certificate (Let's Encrypt for free)
- HSTS header: `Strict-Transport-Security: max-age=31536000`

**Nginx config example:**
```nginx
server {
  listen 443 ssl http2;
  ssl_certificate /etc/letsencrypt/live/dawai.example.com/fullchain.pem;
  ssl_certificate_key /etc/letsencrypt/live/dawai.example.com/privkey.pem;
  ssl_protocols TLSv1.2 TLSv1.3;
  add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
  
  location / {
    proxy_pass http://api:8080;
  }
}
```

---

## 11. Secrets Management

### 11.1 Environment Variables

**Must Never Commit:**
- `JWT_SECRET` (min 32 bytes, random)
- `DATABASE_URL` (with password)
- `MINIO_SECRET_KEY`
- `AUTH_SECRET` (NextAuth)

### 11.2 Deployment Secrets

Store in:
- Kubernetes Secrets (if K8s)
- AWS Secrets Manager / Parameter Store
- HashiCorp Vault
- Woodpecker CI secrets (encrypted, per-project)

**Never:**
- Hardcode in code
- Commit to .env file (use .env.example template)
- Log secret values
- Pass in URL query params

---

## 12. Known Gaps & To-Fix Items

### 12.1 Missing JWT Blacklist Check

**Status:** ❌ TO FIX

**Issue:** JWTGuard validates signature + expiry but does NOT check blacklist table.

**Impact:** User who logs out can still use old JWT until expiry (7 days).

**Fix:**
```go
// In JWTGuard, after signature validation:
blacklisted, err := db.IsJWTBlacklisted(ctx, claims.ID)
if blacklisted {
  return c.Status(401).JSON(...("Token revoked"))
}
```

**Effort:** 5–10 min (query, cache if performance needed).

---

### 12.2 No Rate Limiting

**Status:** ❌ TO FIX

**Issue:** Auth endpoints (login, forgot-password) have no brute-force protection.

**Impact:** Account takeover via credential guessing.

**Fix:** Fiber middleware (gofiber/limiter) + Redis.

**Effort:** 30 min (middleware + redis config).

---

### 12.3 Multi-School User Support

**Status:** 🟡 PARTIAL

**Issue:** Login returns first school only; users with roles in multiple schools can't switch.

**Impact:** Users must log out / log back in with different school.

**Fix:** 
1. Extend JWT to include all schools + roles
2. Add `x-school-id` header validation
3. Client-side school selector

**Effort:** 2–3 hours.

---

### 12.4 No Student Role Parent Restriction (ISOLATION-006)

**Status:** ❌ TO FIX (Confirmed bug)

**Issue:** `/api/students` endpoints have no parent-only check; any authenticated user can list all students.

**Impact:** Student/parent data exposure across school.

**Fix:** Add RoleGuard or explicit role check in handler.

**Effort:** 5 min.

---

### 12.5 HTTPS Enforcement

**Status:** ❌ NOT DEPLOYED

**Issue:** No HTTPS in development; production needs TLS + HSTS.

**Fix:** Reverse proxy (Nginx/Caddy) + Let's Encrypt.

**Effort:** Setup-dependent.

---

### 12.6 Parent-Scoped Access

**Status:** 🟡 PARTIAL

**Issue:** Parent role exists but no endpoints validate parent-child relationship.

**Impact:** Parent can only read own user record, not child assessments.

**Fix:** 
1. Add `GET /api/students/:id/assessments` with parent check
2. Validate request user is parent of :id

**Effort:** 1–2 hours.

---

### 12.7 Missing DB-Level school_id Filters in By-ID Queries

**Status:** ❌ TO FIX

**Issue:** Four by-ID lookup queries in `queries.sql` lack `school_id` parameter, allowing direct ID access bypass:
- `GetSubjectByID` (line 32) — `SELECT * FROM subjects WHERE id = $1`
- `GetStudentByID` (line 55) — `SELECT * FROM students WHERE id = $1`
- `GetAssessmentByID` (line 65) — `SELECT * FROM assessments WHERE id = $1`
- `GetRubricComponentByID` (line 78) — `SELECT * FROM rubric_components WHERE id = $1`

**Impact:** Handler must validate `school_id` manually; if missed, cross-tenant read possible.

**Fix:** Add school_id to WHERE clause (requires parameter change in calling handlers):
```sql
-- GetSubjectByID
SELECT * FROM subjects WHERE id = $1 AND school_id = $2;

-- GetStudentByID
SELECT * FROM students WHERE id = $1 AND school_id = $2;

-- GetAssessmentByID
SELECT * FROM assessments WHERE id = $1 AND school_id = $2;

-- GetRubricComponentByID
SELECT * FROM rubric_components WHERE id = $1 AND school_id = $2;
```

**Effort:** 15 min (update queries.sql, regenerate sqlc, update 4 handlers to pass school_id).

---

## 13. Deployment Checklist

**Before going to production:**

- [ ] JWT_SECRET is ≥32 bytes, random, environment-only
- [ ] DATABASE_URL uses strong password, encrypted connection
- [ ] HTTPS/TLS configured + HSTS header
- [ ] Rate limiting enabled on auth endpoints
- [ ] Audit logging tested + logs stored securely
- [ ] All school_id filters verified via query audit
- [ ] Cross-tenant tests pass (isolation_test.go)
- [ ] Error messages don't leak data
- [ ] JWT blacklist check implemented
- [ ] Bcrypt cost ≥12 verified
- [ ] MINIO_SECRET_KEY rotated + stored securely
- [ ] Backup strategy for audit_logs table
- [ ] Log retention policy documented
- [ ] Admin dashboards audit access

---

## 14. References

- **JWT Best Practices:** https://tools.ietf.org/html/rfc8725
- **OWASP Top 10:** https://owasp.org/www-project-top-ten/
- **Multi-Tenancy Isolation:** https://cheatsheetseries.owasp.org/cheatsheets/Multi_Tenant_SaaS_HTML5_Web_Application_Cheat_Sheet.html
- **Bcrypt Cost:** https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
- **HSTS:** https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security
