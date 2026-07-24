# DAWAI Integration: Agent Teams × Project Files

How `CLAUDE.md` and `DESIGN.md` constraints feed into agent team spawn commands.

---

## CLAUDE.md → Agent Team Spawn Rules

CLAUDE.md defines DAWAI's non-negotiable constraints. Every agent team spawn must reference these.

### Multi-Tenant Isolation

```
CRITICAL: Every query MUST filter by school_id from JWT claims via c.Locals("school_id").
Never from request body, query params, or URL path.
Cross-tenant validation: verify entity school_id matches request context.
```

**Spawn prompt injection:**

```text
[task description]

CRITICAL company policy:
- school_id MUST come from c.Locals("school_id") — NEVER from request body
- Cross-tenant validation required for every entity reference
- See CLAUDE.md for full isolation rules
```

### JWT Auth Flow

```
NextAuth (Google OAuth / email+password) → POST /api/auth/token → JWT w/ school_id + roles
→ JWT stored in httpOnly cookie → injected into Authorization header on API calls
→ Middleware: JWTGuard + TenantGuard + RoleGuard
```

**Spawn prompt for auth-related work:**

```text
DAWAI auth flow: NextAuth → Go /api/auth/token → JWT (school_id, roles array).
JWTGuard validates signature + expiry. TenantGuard extracts school_id.
RoleGuard checks roles array matches endpoint requirements.
Token blacklist on logout: POST /api/auth/logout inserts jti into jwt_blacklist.
```

### Tech Stack

```
Backend: Go / Fiber (port 8080) + PostgreSQL + MinIO
Frontend: Next.js (port 3000) + Tailwind + shadcn/ui
Migrations: golang-migrate (raw SQL, no Prisma)
i18n: next-intl (en, id)
```

**Spawn prompt for implementation:**

```text
DAWAI tech stack:
- Backend: Go / Fiber, PostgreSQL via pgx, sqlc for generated queries
- Frontend: Next.js 15 App Router, Tailwind, shadcn/ui
- DB migrations: golang-migrate raw SQL files in backend/migrations/
- Multi-tenant: school_id in JWT claims, app-layer filtering (no RLS)
```

### Commit Message Format

```
Co-Authored-By: Claude Haiku 4.5 <noreply@anthropic.com>
```

Add to any commit from agent team output.

---

## DESIGN.md → Agent Team Spawn Rules

DESIGN.md defines visual language. Spawn frontend/design teammates with these baked in.

### Color & Typography Constraints

```
Canvas: Parchment #FAF8F5
Accent: Rosewood #B5603C (singular — no gradients, no neon variants)
Typography: Satoshi (UI), Geist Mono (numeric data)
Banned: Inter, Roboto
Body max: 65ch
Assessments: full-width sliders, Geist Mono score numerals
```

**Spawn prompt for design work:**

```text
DAWAI design system constraints (DESIGN.md):
- Palette: Parchment #FAF8F5 canvas, Rosewood #B5603C accent, Ink Deep #1C1917 text
- Typography: Satoshi headings/body, Geist Mono for scores/timestamps
- Cards: rounded-xl (1.5rem), ivory background, whisper border
- Sliders: 44px touch zone, Geist Mono score counter, Rosewood active track
- Dark mode: #1C1917 background, #292524 surface, Rosewood accent unchanged
```

### Dark Mode Tokens

```
Background: #1C1917
Surface: #292524
Text primary: #FAF8F5
Text secondary: #A8A29E
Border: rgba(87,83,78,0.5)
Rosewood accent: unchanged
```

Include in any spawn touching component CSS:

```text
Dark mode tokens also defined. Mandatory for all new components to support both modes.
```

---

## TEAM_WORKFLOWS.md → Agent Team Decision Chains

Workflows define which roles approve which outputs. Use as spawn coordination rules.

### Feature Development Chain

| Step | Decision Maker | Input | Output |
|------|---------------|-------|--------|
| Scope | principal-pm | Stakeholder request | PRD |
| Approval | product-director | PRD + roadmap | APPROVED/DEFERRED |
| Architecture | solution-architect | PRD | API spec + design |
| Design | staff-designer | PRD + spec | Figma + components |
| Implementation | staff-software-engineer | Design + spec | Code + tests |
| QA | qa-architect | Feature + criteria | Test results |
| Deploy | devops-architect | Release notes | Production |

**Spawn coordination rule for multi-role teams:**

```text
Each teammate produces a specific artifact. No downstream teammate starts until
the upstream artifact is ready. Chain:
1. principal-pm writes PRD → product-director approves
2. solution-architect designs API → 3. staff-designer creates components
4. staff-software-engineer implements → 5. qa-architect validates
6. devops-architect deploys
```

### Security Review Gate

```
Any change to auth/authz, data access, encryption, third-party integrations, user data:
→ AUTOMATIC trigger on PR → security-architect reviews and gates
→ If requested changes: staff-software-engineer implements fix → security-architect re-reviews
→ Final tech gate: solution-architect
```

**Spawning a security review team:**

```text
Spawn a security reviewer using the security-architect agent type.
Gate rule: must report APPROVE or provide specific fix list before merge.
Escalation: solution-architect if security and implementation disagree.
```

### Incident Escalation

| Priority | Response | Team Composition |
|----------|----------|-----------------|
| P1 (production outage) | <5 min | staff-software-engineer + devops-architect + sre-lead |
| P2 (degraded) | <30 min | relevant architect + engineer |
| P3 (minor bug) | <24 hours | single engineer |

```text
P1 incident. Spawn incident response team:
- staff-software-engineer: find root cause
- devops-architect: system resources + rollback
- sre-lead: SLO monitoring + customer comms
```

---

## AGENTS.md → Agent Team Role Assignments

AGENTS.md defines each role's authority, success metrics, and delegation structure. Cross-reference when spawning.

### Role Authority Hierarchy

| Category | Decides | Consults | Informs |
|----------|---------|----------|---------|
| Business strategy | ceo | founder | product-director |
| Product roadmap | product-director | ceo | principal-pm |
| Feature scope | principal-pm | product-director | solution-architect |
| Architecture | solution-architect | backend-architect, frontend-architect | security-architect |
| API contract | backend-architect | solution-architect | frontend-architect |
| DB schema | database-architect | solution-architect | security-architect |
| Security | security-architect | solution-architect | backend-architect |
| UX/Design | staff-designer | principal-pm | ux-researcher |
| Code quality | staff-software-engineer | engineering-manager | qa-architect |
| Deployment | devops-architect | solution-architect | platform-engineer |
| Release gates | qa-architect | product-director | engineering-manager |

Use to resolve disputes between teammates:

```text
If architecture and product disagree on scope: solution-architect decides,
principal-pm provides requirements. Escalate to product-director if blocked.
```

### Success Metrics Per Role

| Role | Primary KPI | Secondary KPI |
|------|------------|---------------|
| ceo | Revenue, Churn | Market share |
| product-director | Launch velocity | Feature adoption |
| principal-pm | Scope clarity | Zero rework rate |
| solution-architect | Design quality | Zero rework rate |
| backend-architect | API latency, uptime | Code quality |
| database-architect | Query latency | Zero data issues |
| security-architect | Incident rate | Compliance pass rate |
| staff-software-engineer | Bug rate, code review time | Test coverage |
| platform-engineer | Onboarding time | Deploy time |
| qa-architect | Production bug rate | Test automation % |
| sre-lead | Uptime SLO | MTTR |

**Spawn with success criteria:**

```text
Spawn a staff-software-engineer to implement auth endpoint.
Success: zero bugs, clean code review, tests pass.
Authority: implementation approach, code standards.
Skip: architecture decisions (belongs to solution-architect).
```

---

## Quick Reference: Spawn Templates by Document

| When task touches... | Spawn template includes... |
|---------------------|--------------------------|
| Data isolation | "school_id from c.Locals(), cross-tenant verify" |
| JWT/Auth | "NextAuth → POST /api/auth/token → JWTGuard → TenantGuard" |
| Frontend UI | "Satoshi type, Geist Mono scores, Rosewood accent, Parchment canvas" |
| Dark mode | "Dark mode required: #1C1917 bg, #292524 surface, Rosewood accent" |
| DB schema | "golang-migrate raw SQL, sqlc generate, look at existing migration pattern" |
| RBAC | "RoleGuard middleware, roles from JWT array, school_id field" |
| Code review | "CLAUDE.md isolation rules mandatory check. security-architect gate if auth/data change" |

---

## Conflict Resolution Between Teammates

When spawned teammates disagree:

```
1. Consulting rule: team asks if decision category is theirs (see authority table)
2. Escalation: technical → solution-architect, product → product-director, blended → ceo
3. Deadlock: if same role disagrees (e.g., 2 staff-software-engineer dispute), solution-architect arbitrates
4. No consensus after 2 rounds: escalate + pause that thread, work elsewhere
```

**Message to use in conflict:**

```text
Decision category check: who owns this?
- API design → solution-architect decides
- Feature scope → principal-pm decides
- Security → security-architect decides
- If unclear → escalate to product-director or ceo
```
