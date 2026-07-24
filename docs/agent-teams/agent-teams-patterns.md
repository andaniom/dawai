# Agent Teams Patterns for Software Development

Proven patterns for spawning effective teams on common software tasks. Copy/adapt these prompts.

---

## Code Review Patterns

### Parallel Domain Review

Split review by expertise area to hit all angles simultaneously:

```text
Review PR #<number> using 3 independent reviewers:

- Security reviewer: 
  Focus on security vulnerabilities, authentication, authorization,
  data validation, injection attacks, cryptography, secrets management.
  Rate severity: critical/high/medium/low.

- Performance reviewer:
  Focus on algorithmic complexity, database queries (N+1s), caching strategy,
  memory usage, unnecessary allocations, bulk operations.
  Suggest benchmarks if relevant.

- Maintainability reviewer:
  Focus on code clarity, error handling, logging, testability,
  adherence to project conventions, documentation completeness.
  Flag technical debt or architectural concerns.

Each review independently. When done, compile findings into single report.
Report format: [Severity] Issue - Description - Suggested Fix.
```

### Feature Review with Competing Perspectives

Four specialized reviewers from different angles:

```text
Review the new assessment submission flow (PR #<number>) from 4 angles:

- User advocate: Is UX smooth? Are error messages helpful? Can user recover?
  Test the happy path + common error cases. Flag confusing UX patterns.

- Data integrity specialist: Will data stay consistent? Check transaction handling,
  race conditions, offline sync conflicts, data validation. Focus on school_id isolation.

- Performance auditor: Will this scale? Check query efficiency, caching, connection pooling,
  background job handling. Estimate throughput limits.

- Reliability engineer: Will this stay up? Check error handling, retry logic,
  timeout behavior, monitoring/logging, graceful degradation.

Each review independently, report findings with severity + suggested fix.
```

### Specialized Language Review

When PR touches multiple languages/stacks:

```text
Review auth refactor (PR #<number>) with 3 specialists:

- Go specialist (backend): Review Go code for idioms, error handling, goroutine safety,
  memory leaks, standard library usage. Check sqlc query generation.

- TypeScript specialist (frontend): Review TS types, React hooks, state management,
  async patterns, null coalescing, interface design.

- SQL specialist (database): Review migrations, query efficiency, indexes,
  foreign keys, constraint safety. Check for data races.

Each review their stack. Share findings across team when done.
```

---

## Bug Investigation Patterns

### Competing Hypotheses Debate

Investigate ambiguous failures by spawning teams with competing theories:

```text
Users report: assessment submissions sometimes disappear after being saved.
Last seen yesterday, only in production, not reproducible in staging.

Spawn 5 teammates to investigate competing root causes. Each takes one hypothesis:

1. Database constraint: submissions rolled back due to FK constraint violation
2. Network timeout: submission API times out but client doesn't retry
3. Race condition: offline sync conflicts with concurrent submission from other device
4. Caching bug: stale cache showing "submission exists" locally but not in DB
5. Job queue failure: background job to write to MinIO is failing, marking submission as deleted

Each teammate:
- Investigate logs/metrics/code for their hypothesis
- Trace the failure scenario end-to-end for their theory
- Check evidence supporting/contradicting their theory
- Message other teammates to debate which theory survives scrutiny

When consensus emerges, summarize root cause + evidence + recommended fix.
```

### Multi-Layer Debugging

When bug spans multiple layers:

```text
POST /api/assessments returns 500 inconsistently (2-3 times/minute during peak load).

Spawn 4 teammates:

- Backend troubleshooter: Check Go logs for panics, timeouts, DB connection errors.
  Review error handling in CreateAssessment handler. Check rate limiting.

- Database troubleshooter: Check PostgreSQL logs for slow queries, connection pool exhaustion,
  constraint violations, deadlocks, autovacuum load.

- Network troubleshooter: Check request latency distribution, timeout settings,
  reverse proxy logs, connection drains, TLS handshakes.

- Load testing specialist: Reproduce under load. Check behavior at 50 rps, 100 rps, 500 rps.
  Identify at what load the failures start.

Each investigate their layer. Share findings. Merge insights into root cause.
```

---

## Feature Development Patterns

### Parallel Module Implementation

Team each owns separate module, coordinate via shared task list:

```text
Implement new KurMer report feature. Break into parallel modules:

Spawn 4 teammates, each owning their module:

1. Database & API teammate (backend):
   - Schema: add report_generation_jobs table, rubric_summary view
   - Migrations: 000X_add_kurmer_tables.up/down.sql
   - API: GET /api/kurmer/{assessmentId}, POST /api/kurmer/batch
   - Service: KurMerService with SQL queries
   - No frontend or test code.

2. Frontend UI teammate:
   - Components: ReportPreview, ReportDownload, DownloadFormat selector
   - Pages: /assessments/[id]/report
   - State: report fetch + caching in Zustand
   - No backend or database code.

3. Background Job teammate:
   - Worker: process KurMer generation jobs
   - Queue: Redis job queue integration
   - Retry logic + error handling
   - Logging to audit_logs table
   - No API or frontend code.

4. Test teammate:
   - Unit: test KurMerService calculations
   - Integration: test job processing end-to-end
   - API: test endpoint auth + school_id filtering
   - Frontend: test component loading states + downloads

Each own separate files. Coordinate:
- Database teammate: post schema when ready, others wait for migration
- API teammate: post endpoint signature when ready, frontend teammate uses it
- Background teammate: integrate with API after endpoints exist
- Test teammate: can start after each other teammate has working code

Lead: create 4 tasks. Teammates self-claim + report blockers.
```

### Cross-Tenant Safety Implementation

New feature must be reviewed for isolation bugs:

```text
New student roster feature. Spawn team to implement + verify isolation:

1. Implementation teammate:
   - Add StudentRoster model, queries, handlers
   - Add RBAC: school_admin + teacher can list
   - POST /api/schools/{schoolId}/rosters

2. Isolation verification teammate (starts after impl):
   - Test: Can User A from School A see roster from School B? (must fail)
   - Test: Can User A as teacher see rosters outside their school? (must fail)
   - Test: SQL injection attempt on query filters
   - Test: Privilege escalation (student trying teacher endpoint)
   - Test: JWT tampering (change school_id in token)
   - Report: pass/fail + fix recommendations

Lead: don't merge until isolation teammate reports all pass.
```

---

## Refactoring Patterns

### Parallel Refactor with Verification

Break refactor into parallel pieces, each verified:

```text
Refactor Auth middleware to extract school_id once, not per-handler.

Spawn 3 teammates:

1. Middleware refactorer:
   - Extract school_id + roles into middleware, inject into c.Locals()
   - Test: middleware correctly parses JWT
   - Test: JWT validation failures return 401
   - Test: school_id correctly injected for valid tokens
   - Own files: middleware/auth.go, middleware/tenant.go

2. Handler refactorer (starts after middleware ready):
   - Remove duplicate school_id extraction from all handlers
   - Update handlers to use c.Locals("school_id")
   - Own files: handlers/*.go (except middleware changes)
   - Test: each handler still works after refactor

3. Integration verifier (starts after both ready):
   - End-to-end: call each endpoint, verify isolation
   - Test: cross-tenant access attempts fail
   - Test: performance unchanged
   - Load test: verify no new bottlenecks

Lead: create tasks in order. Middleware first, then handlers, then verify.
```

---

## Security Hardening Patterns

### Multi-Angle Security Audit

Same code, different attack angles:

```text
Security audit of the assessment submission flow.

Spawn 4 security specialists:

1. Input validation specialist:
   - Check every API input: type validation, size limits, allowed chars
   - Check: SQL injection attempts
   - Check: XSS payload attempts
   - Check: malformed JSON/form data
   - Report: missing validation + fixes

2. Authorization specialist:
   - Check: all endpoints verify JWT + school_id
   - Check: role-based access (teacher can't access admin endpoints)
   - Check: user can only access own school's data
   - Check: parent can only access own child's data
   - Report: missing checks + fixes

3. Data isolation specialist:
   - Check: every query filters by school_id
   - Check: no data leaks across schools in responses
   - Check: audit logs don't expose cross-tenant data
   - Check: error messages don't reveal other schools' existence
   - Report: leaks + fixes

4. Cryptography specialist:
   - Check: JWT signing (algorithm, key strength)
   - Check: password hashing (bcrypt cost, salt)
   - Check: HTTPS enforcement
   - Check: secure headers (CSP, HSTS, X-Frame-Options)
   - Check: token expiry + blacklist logic
   - Report: weaknesses + fixes

Each report severity rating + suggested fix. Lead compiles.
```

### Compliance & Hardening Sprint

Multi-angle approach to hardening:

```text
DAWAI Phase 13: security hardening + compliance audit.

Spawn team:

1. Query auditor (grep + code review):
   - Grep all .go files for database operations
   - Verify every SELECT/INSERT/UPDATE/DELETE filters by school_id
   - Create spreadsheet: query → school_id filter (yes/no)
   - Report: queries missing filter + location

2. Access control tester (manual testing):
   - Test each endpoint as different users (student, teacher, admin)
   - Attempt privilege escalation (student as teacher)
   - Attempt cross-tenant access (School A user accessing School B data)
   - Report: access violations + severity

3. Test coverage builder (write tests):
   - Write cross-tenant access tests for every endpoint
   - Test: user A from School A cannot read School B data
   - Test: teacher cannot access admin endpoints
   - Test: parent can only see child's data
   - Report: coverage % + gaps

4. Documentation writer (create compliance docs):
   - Document auth flow: OAuth + JWT + server auth
   - Document isolation model: school_id field + app-layer filtering
   - Document role model: roles + permissions matrix
   - Create: "SECURITY.md" with hardening checklist
   - Report: compliance checklist + status

Lead: review each teammate's output. Merge findings into Phase 13 completion.
```

---

## Testing Patterns

### Parallel Test Suite Implementation

Split test types across team:

```text
Write comprehensive test suite for new assessment endpoints.

Spawn 3 teammates:

1. Unit test specialist:
   - AssessmentService tests (business logic)
   - Test: score calculations
   - Test: rubric component handling
   - Test: idempotency_key deduplication
   - File: backend/internal/services/assessment_service_test.go

2. Integration test specialist:
   - Database + Service integration
   - Test: create assessment → saves to DB + creates audit log
   - Test: update assessment → old data accessible via audit
   - Test: delete assessment → soft delete, audit trail preserved
   - File: backend/internal/handlers/assessment_integration_test.go

3. API + Security test specialist:
   - Handler + middleware tests
   - Test: POST /assessments with valid JWT
   - Test: POST /assessments cross-tenant (must fail)
   - Test: POST /assessments missing school_id (must fail)
   - Test: offline idempotency (resubmit same request)
   - File: backend/internal/handlers/assessment_api_test.go

Each own test file. Coordinate: unit < integration < API (dependencies).
All tests must pass before merge.
```

---

## Review Synthesis Pattern

After parallel review, synthesize findings:

```text
All reviewers done. Now:

Create synthesis document:
- Blocker issues (must fix before merge)
- Should-fix issues (high priority)
- Nice-to-have improvements
- Questions for author (may need clarification)

For each issue:
- Severity (critical/high/medium/low)
- Which reviewer found it
- Description
- Suggested fix

Recommend: merge only after blockers resolved.
```

---

## Debugging with Teams Pattern

```text
Strange race condition in offline sync only on Android after 10+ minutes.

Spawn 4 teammates to debug:

1. Android teammate: reproduce on Android, check logcat, profile battery drain
2. Sync logic teammate: review offline sync code, check state machine logic
3. Network teammate: monitor network requests during failure, check retry logic
4. Database teammate: check SQLite write patterns, journal, locking behavior

Coordinate: reproduce → isolate variable → test hypothesis → verify fix.
Share findings as discovered. Debate which variable is root cause.
```

---

## Team Size Guide

| Task Size | Optimal Team | Reasoning |
|-----------|-------------|-----------|
| Review small PR | 2-3 reviewers | Domain experts only |
| Review medium PR | 3-4 reviewers | Multiple perspectives |
| New feature (module) | 3-4 devs | DB + API + frontend + tests |
| Bug with unclear root cause | 3-5 investigators | Multiple angles, debate |
| Security audit | 3-4 specialists | Input, auth, data, crypto |
| Large refactor | 2-3 specialists | Sequential modules |
| Performance optimization | 2-3 specialists | DB + API + frontend |

Avoid > 5 teammates: coordination overhead explodes.

---

## Common Mistakes to Avoid

❌ **Same-file editing**: Two teammates edit `auth.go` → overwrites. Split by file.

❌ **Unclear task dependencies**: Reviewer starts before implementation. Create task list with dependencies.

❌ **No context in spawn**: "Review the code" vs "Review auth module focusing on JWT token handling, 7-day expiry, blacklist on logout." Second is better.

❌ **Too many tasks**: 20 tasks for 3 teammates = multitasking thrashing. Aim for 5-6 tasks per teammate.

❌ **Unattended team**: Spawn team then disappear. Check in, steer approaches, synthesize findings.

❌ **Ignoring blockers**: Implementer waiting for schema. Lead must notice + unblock.

❌ **No verification step**: Refactor done, but no one verified it works. Build in verification task.

---

## Checklist Before Spawning

- [ ] Task has clear parallel components? (not sequential)
- [ ] Each component can be worked independently? (no blocking)
- [ ] Do components need to communicate? (if yes → teams, if no → subagents)
- [ ] Each teammate owns different files? (no conflicts)
- [ ] Enough context in spawn prompt? (include task-specific details)
- [ ] Named teammates? (easier to reference)
- [ ] Created task list with dependencies? (if complex)
- [ ] Clear deliverables? (each teammate knows what done looks like)
- [ ] Token budget allows it? (teams are expensive)
