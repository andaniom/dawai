# DAWAI Team Workflows & Decision Chains

---

## 1. Role Mapping: Roster → Agent Roles

| Agent Role | Assigned Person | Title | FTE | Escalation |
|---|---|---|---|---|
| **CEO** | Andi Nugroho | CEO/Founder | 100% | Board/Investors |
| **Founder** | Andi Nugroho | CEO/Founder | (same) | Board/Investors |
| **Product Director** | Siti Nurhaliza | Product Director | 100% | CEO |
| **Principal Product Manager** | Rania Pratiwi | Principal PM | 100% | Siti Nurhaliza |
| **Principal Business Analyst** | Bambang Sutrisno | Principal BA | 100% | Siti Nurhaliza |
| **Solution Architect** | Iwan Setyawan | Solution Arch / Backend Arch | 100% | CEO (tech) |
| **Frontend Architect** | Dewi Kusuma | Frontend Arch / Staff Eng | 100% | Iwan Setyawan |
| **Backend Architect** | Iwan Setyawan | Solution Arch / Backend Arch | (same) | CEO (tech) |
| **Database Architect** | Ratna Wijaya | Database Architect | 80% | Iwan Setyawan |
| **Security Architect** | Arief Gunawan | Security Architect | 60% | Iwan Setyawan |
| **Staff Product Designer** | Maya Handoko | Staff Designer / DS Lead | 100% | Siti Nurhaliza |
| **Design System Lead** | Maya Handoko | Staff Designer / DS Lead | (same) | Siti Nurhaliza |
| **UX Researcher** | Lena Pratama | UX Researcher / Accessibility | 80% | Siti Nurhaliza (advisory) |
| **Accessibility Specialist** | Lena Pratama | UX Researcher / Accessibility | (same) | Siti Nurhaliza (advisory) |
| **Staff Software Engineer (Backend)** | Ahmad Hidayat | Backend Software Engineer | 100% | Iwan Setyawan |
| **Engineering Manager** | Yuli Hartono | Engineering Manager | 100% | Iwan Setyawan + Siti Nurhaliza |
| **Platform Engineer** | Hendra Wijaya | DevOps / Platform Eng | 100% | Iwan Setyawan |
| **DevOps Architect** | Hendra Wijaya | DevOps / Platform Eng | (same) | Iwan Setyawan |
| **QA Architect** | Nanda Kusuma | QA Architect / SRE Lead | 100% | Siti Nurhaliza |
| **SRE Lead** | Nanda Kusuma | QA Architect / SRE Lead | (same) | Siti Nurhaliza |
| **Music Education Consultant** | Prof. Budi Santosa | Music Education Consultant | 20% | Siti Nurhaliza (advisory) |
| **Indonesia Education Consultant** | Dr. Nur Cahyadi | Indonesia Education Consultant | 20% | Siti Nurhaliza (advisory) |

**Dual Roles (cost optimization):**
- Andi Nugroho: CEO + Founder
- Iwan Setyawan: Solution Architect + Backend Architect
- Dewi Kusuma: Frontend Architect + Staff Software Engineer
- Ratna Wijaya: Database Architect + Platform/DevOps Advisory
- Hendra Wijaya: Platform Engineer + DevOps Architect
- Nanda Kusuma: QA Architect + SRE Lead
- Maya Handoko: Staff Product Designer + Design System Lead
- Lena Pratama: UX Researcher + Accessibility Specialist

---

## 2. Decision Workflows & Approval Chains

### Feature Request → Launch

```
Stakeholder Request
    ↓
Rania Pratiwi (Principal PM)
  ├─ Consult: Bambang (BA), Lena (UX), Budi (Music), Nur (Indonesia Education)
  └─ Output: PRD with acceptance criteria, edge cases, success metrics
    ↓
Siti Nurhaliza (Product Director)
  ├─ Gate: Roadmap fit? Business value? Resource availability?
  └─ Decision: APPROVE / DEFER / REJECT
    ↓ [APPROVE]
Iwan Setyawan (Solution Architect)
  ├─ Consult: Dewi (Frontend), Ahmad (Backend), Hendra (DevOps), Arief (Security)
  └─ Output: Architecture design, API contract, implementation estimate
    ↓
Rania Pratiwi (Principal PM) [Re-gate]
  ├─ Confirm: Still fits roadmap? Scope clear? Risks captured?
  └─ Decision: READY TO IMPLEMENT / RENEGOTIATE
    ↓ [READY]
Yuli Hartono (Engineering Manager)
  ├─ Assign: Dewi + Ahmad (frontend/backend), Hendra (if infra change)
  ├─ Estimate: Sprint capacity
  └─ Decision: SPRINT ASSIGNMENT / BACKLOG
    ↓
Nanda Kusuma (QA Architect)
  ├─ Design: Test plan, E2E coverage, performance benchmarks
  └─ Gate: Release-ready criteria defined before coding starts
    ↓ [Development Cycle]
Yuli Hartono (Engineering Manager)
  ├─ Review: Daily standup, blocker resolution, quality checks
  └─ Gate: Code review + QA pass before merge
    ↓
Nanda Kusuma (QA Architect)
  ├─ Test: E2E, regression, accessibility, performance
  └─ Gate: Release quality criteria met?
    ↓ [PASS]
Hendra Wijaya (DevOps)
  ├─ Deploy: Staging → production
  └─ Monitor: Error rates, performance metrics
    ↓
Siti Nurhaliza (Product Director) [Final gate]
  ├─ Verify: Feature works as designed, metrics track, customer ready
  └─ Publish: Announce to customers

Escalations:
- Scope bloated → Rania escalates to Siti
- Arch concerns → Iwan escalates to Andi (CEO)
- Timeline slip → Yuli escalates to Siti + Iwan
- QA blockers → Nanda escalates to Yuli + Siti
- Production issue → Nanda escalates to Iwan → Andi
```

### Security Review Workflow

```
Any code change affecting:
  - Authentication / Authorization
  - Data access (school_id filtering)
  - Encryption / Secrets
  - Third-party integrations
  - User data handling

Trigger: AUTOMATIC on PR
    ↓
Arief Gunawan (Security Architect)
  ├─ Review: Threat model, OWASP coverage, JWT validation, tenant isolation
  ├─ Consult: Iwan (Backend), Ratna (DB), Hendra (Infra)
  └─ Gate: APPROVE / REQUEST CHANGES / ESCALATE
    ↓ [REQUEST CHANGES]
Ahmad Hidayat (Backend Engineer)
  ├─ Fix: Implement security recommendations
  └─ Resubmit: Arief re-reviews
    ↓
Iwan Setyawan (Solution Architect)
  ├─ Final tech gate: Architecture still sound?
  └─ Approve merge if all security + arch concerns resolved
```

### Database Schema Change Workflow

```
Schema modification (new table, new column, index, constraint)
    ↓
Ratna Wijaya (Database Architect)
  ├─ Design: Normalization, indexing, multi-tenant isolation, performance
  ├─ Consult: Iwan (Backend), Arief (Security - if PII/sensitive data)
  └─ Output: Migration DDL + rollback script
    ↓
Iwan Setyawan (Solution Architect)
  ├─ Gate: API impact? Transaction boundaries? Performance OK?
  └─ Approve if safe
    ↓
Ahmad Hidayat (Backend Engineer)
  ├─ Test: Migration on staging DB, verify data integrity
  └─ Gate: Passes test before production deployment
    ↓
Hendra Wijaya (DevOps)
  ├─ Deploy: Run migration on production with backup
  ├─ Monitor: Query performance, error rates
  └─ Rollback capability confirmed before release

Escalation: Schema issue blocks 2+ features → Ratna escalates to Iwan → Andi
```

### Design Review Workflow

```
UI/Component change
    ↓
Maya Handoko (Staff Designer + Design System Lead)
  ├─ Design: Figma mockup, component definition, accessibility review
  ├─ Consult: Lena (UX research), Rania (product requirements match)
  └─ Output: Design spec + component handoff to frontend
    ↓
Lena Pratama (UX Researcher)
  ├─ Research: User testing (if major change), WCAG audit
  └─ Gate: Accessibility AAA? Usability OK? Task completion time acceptable?
    ↓ [PASS]
Dewi Kusuma (Frontend Architect)
  ├─ Review: Component implementation, performance, responsive behavior
  ├─ Consult: Maya on design fidelity, Lena on a11y
  └─ Gate: Code review before merge
    ↓
Nanda Kusuma (QA Architect)
  ├─ Test: Visual regression, a11y regression, responsive testing
  └─ Gate: Release quality criteria met

Escalation: Design blocks feature → Maya escalates to Siti → Andi
```

### Infrastructure / DevOps Change Workflow

```
Infrastructure modification (Docker, K8s, CI/CD, secrets, networking)
    ↓
Hendra Wijaya (DevOps Architect)
  ├─ Design: Change plan, rollback strategy, monitoring setup
  ├─ Consult: Iwan (service impact), Nanda (SLO impact), Arief (security)
  └─ Output: Deployment plan + runbook
    ↓
Iwan Setyawan (Solution Architect)
  ├─ Gate: System architecture still sound? Scalability maintained?
  └─ Approve if safe
    ↓
Nanda Kusuma (SRE Lead)
  ├─ Monitor plan: What metrics? Alert thresholds? Rollback triggers?
  └─ Gate: Observability ready before deployment
    ↓
Hendra Wijaya (DevOps)
  ├─ Deploy: Staging first, then production
  ├─ Monitor: Real-time dashboards, alert response
  └─ Gate: SLO still met? No incidents?
    ↓
Iwan Setyawan (Solution Architect) [Post-flight check]
  ├─ Verify: Change successful, no regressions
  └─ Approve change closure
```

---

## 3. Handoff Protocols (When & How Agents Delegate)

### Decision Handoff

**Scenario 1: Feature Request → Product Director → Architect**

When: PRD is complete and roadmap-approved
Who: Rania (Principal PM) → Siti (Product Director) → Iwan (Architect)

**Handoff packet required:**
```
- PRD (1-2 pages): user story, acceptance criteria, success metrics
- Edge cases: Known scope boundaries, what's explicitly NOT included
- Constraints: Timeline, tech stack assumptions, budget
- Dependencies: Other features, external systems, data migrations
- Risks: Known unknowns, technical concerns
- Contacts: DRI (Rania), stakeholders (consultants if domain-specific)
```

**Handoff meeting (30 min):**
- Rania presents PRD to Iwan + Dewi + Ahmad
- Iwan asks clarifying questions on scope/constraints
- Iwan estimates architecture effort
- Iwan identifies missing details → Rania refines PRD if needed
- Outcome: "Ready for architecture design" or "Need more detail"

---

**Scenario 2: Architecture → Engineering Manager → Engineers**

When: Architecture design is complete and approved
Who: Iwan (Architect) → Yuli (Engineering Manager) → Dewi + Ahmad (Engineers)

**Handoff packet required:**
```
- Architecture design doc: Service boundaries, API contracts, DB schema
- Implementation plan: Breakdown into tasks, estimated story points
- Code patterns: Link to existing examples in codebase (auth, validation, queries)
- Testing strategy: Unit test coverage targets, E2E scenarios
- Dependencies: Other work that must complete first
- Deployment plan: Infrastructure changes needed (if any)
- Contacts: DRI (Iwan), on-call for questions during implementation
```

**Handoff meeting (45 min):**
- Iwan walks through architecture diagram
- Iwan explains key decisions: why these boundaries, why this API design
- Dewi + Ahmad ask implementation questions
- Yuli estimates sprint capacity
- Outcome: Sprint assignment + detailed task breakdown

---

**Scenario 3: Engineering → QA Architect → QA Team**

When: Feature is code-complete and ready for testing
Who: Yuli (Engineering Manager) → Nanda (QA Architect) → Nanda (QA execution)

**Handoff packet required:**
```
- Feature summary: What changed, user-facing behavior
- Test plan: E2E scenarios to cover (happy path + edge cases)
- Acceptance criteria: What "done" means (performance targets, a11y, etc.)
- Performance benchmarks: Expected response times, load thresholds
- Regression scope: What existing features might be affected
- Environment: Where to test (staging, test data setup)
- Contacts: DRI (Yuli), engineer on-call if questions arise
```

**Handoff meeting (30 min):**
- Engineer demos feature
- Nanda discusses test scenarios
- Nanda clarifies edge cases
- Outcome: QA can execute test plan without engineer involvement

---

**Scenario 4: QA → DevOps → Production**

When: QA passes and feature is release-ready
Who: Nanda (QA Architect) → Hendra (DevOps) → Production

**Handoff packet required:**
```
- QA sign-off: All tests passed, no blockers, accessibility verified
- Deployment plan: Changes to config, secrets, infrastructure
- Rollback plan: How to revert if issues detected
- Monitoring plan: What metrics to watch, alert thresholds
- Customer communication: Release notes, feature documentation
- On-call contact: Who to escalate to if production issues
```

**Handoff meeting (30 min):**
- Nanda confirms QA completion
- Hendra confirms deployment readiness
- Hendra confirms monitoring/rollback is ready
- Outcome: Deployment can proceed

---

### Information Handoff (Async)

**Daily standup (9:00 AM, 15 min):**
- Each person: What done yesterday, what doing today, blockers
- Owner: Yuli (Engineering Manager)
- Async option: Slack thread if timezone conflict

**Weekly Sync by Function:**
- **Backend** (Wed 3 PM): Iwan, Ahmad, Ratna, Arief
- **Frontend** (Wed 4 PM): Dewi, Maya, Lena
- **DevOps/Infra** (Thu 10 AM): Hendra, Nanda, Iwan
- **Product** (Tue 2 PM): Siti, Rania, Bambang, Nanda, Maya

**Async Decisions (Slack #dawai-decisions):**
- Any decision affecting >1 team
- DRI posts: Decision needed, deadline (default 24 hours)
- Stakeholders review + approve/comment
- Final decision posted with rationale + implementation owner
- Outcome tracked in DECISIONS.md

---

## 4. Incident Escalation Matrix

### Priority Levels

| P1 | P2 | P3 |
|---|---|---|
| **Production outage** (>50 users affected, cannot login/assess) | **Degraded service** (some features slow, intermittent errors) | **Non-urgent bug** (feature works, minor friction, <5 users affected) |
| SLO violated | SLO at risk | No SLO impact |
| **Response: <5 min** | **Response: <30 min** | **Response: <24 hours** |

---

### P1 Incident Escalation

```
Production Issue Detected (monitoring alert OR customer report)
    ↓
On-Call Engineer (Ahmad or Hendra depending on layer)
  ├─ Acknowledge: Confirm production status, scope of impact
  ├─ Triage: Is this P1, P2, or P3?
  └─ Page: If P1 → page immediate escalation team
    ↓ [P1 Confirmed]
On-Call Engineer
  ├─ Immediate actions: Rollback if available, kill hanging processes, stop bleeding
  ├─ Page: Iwan (Backend) + Hendra (DevOps) + Nanda (SRE)
  └─ Start: War room in #dawai-incident Slack channel
    ↓
Incident Commander (Iwan or Hendra, whoever responds first)
  ├─ Establish: War room, timeline, communication cadence
  ├─ Assign: Root cause analysis (Iwan), mitigation (Hendra), monitoring (Nanda)
  └─ Broadcast: Customer communication owner (Siti or Rania)
    ↓
Root Cause Investigation (5-15 min)
  ├─ Iwan: Check logs, recent deployments, config changes
  ├─ Ahmad: Query DB, check slow queries, locks
  ├─ Hendra: Check system resources, network, Docker/K8s status
  └─ Output: "Issue is [X], fix is [Y], ETA [Z]"
    ↓
Mitigation (immediate)
  ├─ Option A: Rollback last deployment → Hendra executes
  ├─ Option B: Database query kill → Ahmad executes
  ├─ Option C: Infrastructure change → Hendra executes
  ├─ Option D: Feature flag disable → Dewi or Ahmad (requires code change)
    ↓
Restore & Verify (until SLO met)
  ├─ Hendra: Monitor system metrics
  ├─ Nanda: Confirm SLO metrics recovering
  ├─ Ahmad: Run smoke tests on critical paths
  └─ Gate: SLO restored for 5 minutes before declaring resolved
    ↓
Incident Closeout (within 24 hours)
  ├─ Incident Commander posts: Timeline, root cause, mitigation, prevention
  ├─ Team: Creates tasks for permanent fixes in backlog
  └─ Siti: Sends customer communication (what happened, why, steps taken)
```

**On-Call Rotation:**
- **Backend:** Ahmad Hidayat (primary), Iwan Setyawan (backup)
- **DevOps/Infra:** Hendra Wijaya (primary), Ratna Wijaya (backup, 80% available)
- **SRE/Monitoring:** Nanda Kusuma (primary), Ahmad Hidayat (backup)

**On-Call Coverage:**
- Weekdays: 6:00 AM - 8:00 PM GMT+7
- Weekends: 8:00 AM - 6:00 PM GMT+7
- After-hours (production only): Rotation among primary contacts

---

### P2 Incident Escalation

```
P2 Issue Detected (monitoring alert, customer report, or QA finding)
    ↓
On-Call Engineer
  ├─ Acknowledge: Confirm degradation, scope
  ├─ Triage: Estimate impact (how many users, how long?)
  └─ Page: Relevant architect (Iwan, Dewi, Hendra) if >30 min fix needed
    ↓
Engineer + Relevant Architect (within 30 min)
  ├─ Diagnose: Root cause, workaround, permanent fix
  ├─ Assign: Who will fix (engineer) vs who will design fix (architect)
  └─ Communicate: ETA to Nanda (SRE) for customer messaging
    ↓
Implement & Deploy (same day target)
  ├─ Engineer: Develop fix, test in staging
  ├─ Architect: Code review for correctness
  ├─ DevOps: Deploy to production with monitoring
  └─ SRE: Confirm metrics improving
    ↓
Post-Incident (within 48 hours)
  ├─ Create: Task in backlog for permanent fix if needed
  ├─ Document: Root cause + fix in #dawai-incidents channel
  └─ Share: Learnings with team (if pattern repeats)
```

---

### P3 Incident Escalation

```
P3 Issue Detected (QA bug, customer minor complaint, code review finding)
    ↓
Engineer (non-urgent response)
  ├─ Log: Create task in backlog, assign priority
  ├─ Plan: Fix in next sprint or two
  └─ No escalation unless it becomes P2
    ↓
During Sprint Planning
  ├─ Prioritize: Backlog grooming, engineer estimates effort
  └─ Assign: Next available engineer slot
    ↓
Normal Development Cycle
  ├─ Develop: Code + tests + code review
  ├─ QA: Test in staging
  └─ Deploy: Next release cycle
```

---

### Critical Infrastructure Incidents

**Database down:**
```
Ratna Wijaya (Database Architect) paged immediately
  ├─ Diagnose: Connectivity, replication, disk space, query locks
  ├─ Escalate: Iwan (Backend) if data corruption suspected
  └─ Recovery: Restore from backup if needed (coordinated with Hendra)
```

**Production Server down / K8s cluster issue:**
```
Hendra Wijaya (DevOps) paged immediately
  ├─ Diagnose: Resource exhaustion, networking, service mesh
  ├─ Escalate: Iwan if application-layer problem
  └─ Recovery: Restart services, rebalance load, failover if available
```

**Security breach / Data leak detected:**
```
Arief Gunawan (Security Architect) paged immediately
  ├─ Activate: Incident response protocol
  ├─ Isolate: Affected systems, disable compromised accounts
  ├─ Escalate: Andi Nugroho (CEO) + legal/compliance immediately
  └─ Coordinate: Full forensic investigation
```

---

### Incident War Room Communication

**Slack channel:** #dawai-incident (created per incident)

**Messages every 5 minutes during active incident:**
```
[Incident Status Update — 14:35 GMT+7]
Incident: Database connection pool exhausted
Status: ACTIVE
Affected: All API endpoints
Impact: 127 users unable to submit assessments
ETA Recovery: 14:55 (estimated 20 min)
Current Action: Restarting connection pool service
On-Call: Iwan, Hendra, Ratna, Nanda
```

**Escalation signals:**
- If issue unresolved after 15 minutes → escalate to Andi (CEO)
- If SLO recovery ETA exceeds 1 hour → notify customers proactively
- If root cause unknown after 30 minutes → page security (if potential breach)

---

### Post-Incident Lessons Learned

**Incident Review (within 24 hours):**
- Incident Commander: Writes timeline + root cause analysis
- Team: Discusses what worked, what didn't, how to prevent
- Output: 1-3 action items for backlog (monitoring, automation, documentation)
- No blame culture: Focus on system design that prevents recurrence

**Blameless Postmortem Template:**
```
Incident: [name]
Timeline:
  14:30 - Alert fired for high error rate
  14:35 - Engineer confirmed P1, paged team
  14:42 - Root cause identified: query timeout in GetStudents handler
  14:48 - Mitigation: Disabled slow query logging, query returned to <100ms
  14:55 - SLO restored, incident resolved

Root Cause:
  - A teacher ran batch export for 500 students simultaneously
  - Query was not indexed on (school_id, created_at)
  - No query timeout protection in code

Impact:
  - 240 minutes of degraded service
  - 45 users affected (teacher + students in that school)
  - SLO violated: 2.1% error rate vs 0.1% target

Prevention:
  [ ] Add query timeout to all DB calls (engineer task)
  [ ] Add index on (school_id, created_at) (DBA task)
  [ ] Add alerting for slow queries >1s (SRE task)
  [ ] Load test batch export at scale (QA task)
```

---

## Summary

| Step | Owner | When | Input | Output |
|---|---|---|---|---|
| **Feature Request** | Rania (PM) | Idea phase | Stakeholder request | PRD |
| **Approval Gate** | Siti (Product Director) | PRD complete | Roadmap + resources | APPROVED / DEFERRED |
| **Architecture** | Iwan (Architect) | Approved | PRD | Arch design + API spec |
| **Design** | Maya (Designer) | PRD + Arch | Requirements | Figma mocks + components |
| **Engineering** | Ahmad (Backend), Dewi (Frontend) | Approved + Designed | Arch + Design | Code + tests |
| **QA** | Nanda (QA Architect) | Code complete | Feature + acceptance criteria | Test results |
| **Deployment** | Hendra (DevOps) | QA pass | Release notes + deployment plan | Production live |
| **Incident** | On-call | Production issue | Alert / customer report | Root cause + mitigation |

