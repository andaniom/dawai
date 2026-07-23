# DAWAI Team Onboarding

Welcome to DAWAI. This guide helps new team members understand roles, workflows, and how to contribute.

---

## Quick Start (Day 1)

1. **Read these files in order:**
   - `CLAUDE.md` — Project overview, architecture, key patterns
   - `AGENTS.md` — Your role, decision authority, who to escalate to
   - `TEAM_ROSTER.md` — Team members, skills, contact info
   - `TEAM_WORKFLOWS.md` — How features get approved, deployed, and how incidents work

2. **Join Slack channels:**
   - #dawai-general (everyone)
   - #[your-function] (backend, frontend, devops, qa-testing, product, design)
   - #dawai-incidents (read-only, learn from past incidents)
   - #dawai-decisions (async decision approvals)

3. **Add to calendar:**
   - Daily standup: 9:00 AM GMT+7 (15 min)
   - Your function's sync: 1x/week (see TEAM_WORKFLOWS.md §2)
   - On-call rotation (if backend/devops/qa): See TEAM_ROSTER.md

4. **Get access:**
   - GitHub repo: Ask Hendra for SSH key setup
   - Local dev: Run `docker compose up -d` (see CLAUDE.md)
   - Database: `psql postgresql://dawai:change_me@localhost:5432/dawai`
   - MinIO console: http://localhost:9001 (minioadmin / minioadmin)

5. **First week goals:**
   - Get local dev running
   - Deploy to staging in Docker
   - Attend 1 decision gate / handoff meeting
   - Review 1 pull request in your area

---

## Understanding Your Role

### Find Your Role

1. Look yourself up in `TEAM_ROSTER.md` → Find your assigned agent role(s)
2. Open `AGENTS.md` → Read your role's section
   - What's your decision authority?
   - Who do you escalate to?
   - Who reports to you?
   - Who do you collaborate with?

### Example: You're Ahmad (Backend Software Engineer)

From `TEAM_ROSTER.md`:
- Agent Role: **Staff Software Engineer (Backend)**
- Reports to: Iwan Setyawan
- Collaborates with: Frontend Engineer, QA Architect, SRE

From `AGENTS.md` → Staff Software Engineer section:
- **Authority:** Implementation approach, code standards, tech debt calls
- **Success:** Zero bugs in production, readable code, <2 reviews per PR, clean tests
- **Owns:** Implementation, testing, code review, documentation
- **Escalates to:** Engineering Manager (Yuli)
- **Collaborates with:** Architecture team, QA, SRE

**What this means:**
- You decide HOW to implement features (approach, patterns, tech choices)
- You own code quality (review peers, write tests, clean code)
- Yuli (Engineering Manager) can override your decisions
- Iwan (Architect) designed the system; you implement to that design
- Nanda (QA) owns testing strategy; you ensure your code is testable

---

## Feature Development Workflow

### Phase 1: Feature Request (Days 1-2)

**Owner:** Rania (Principal PM)  
**Your involvement (by role):**

- **Backend Engineer:** Consulted on feasibility. Rania may ping: "Can we do X in 2 weeks?" → Give honest estimate.
- **Frontend Engineer:** Consulted on UI complexity.
- **Architect:** Consulted on system impact.
- **QA:** Consulted on testing scope.

**Handoff:** Rania presents PRD (1-2 pages) to Siti → Siti approves/defers → Feature moves to Phase 2

---

### Phase 2: Architecture Design (Days 3-5)

**Owner:** Iwan (Solution Architect)  
**Your involvement (by role):**

- **Backend Engineer:** Iwan designs API + DB schema. You ask clarifying questions: "Should this endpoint paginate?" → Iwan answers or updates design.
- **Frontend Engineer:** Iwan defines API contract. You plan React components.
- **DevOps:** Iwan may require infra changes. You ask: "Do we need new secrets?"
- **QA:** Nanda designs test plan in parallel.

**Handoff:** Iwan presents architecture to Rania → Rania confirms feature still fits → Feature moves to Phase 3

**Your checklist:**
- [ ] I understand the API contract
- [ ] I understand the DB schema and how my code queries it
- [ ] I know what existing code patterns to follow
- [ ] I have 2-3 clarifying questions answered
- [ ] I can estimate my implementation effort

---

### Phase 3: Implementation (Days 6-12)

**Owner:** Yuli (Engineering Manager)  
**Your involvement (by role):**

- **Backend Engineer:** 
  - [ ] Create feature branch from `main`
  - [ ] Follow code patterns in `CLAUDE.md` (JWT, school_id filtering, audit logging)
  - [ ] Write unit tests (target: >80% coverage of business logic)
  - [ ] Create PR with clear description + link to Iwan's architecture doc
  - [ ] Request review from Iwan (architecture) + peer (code quality)
  - [ ] Address review comments within 24 hours
  - [ ] Merge when both reviews approve

- **Frontend Engineer:**
  - [ ] Create feature branch from `main`
  - [ ] Implement Maya's design using design system components
  - [ ] Write tests for complex interactions (state changes, validation)
  - [ ] Create PR with screenshot/demo
  - [ ] Request review from Dewi (frontend arch) + Maya (design fidelity)
  - [ ] Test on multiple screen sizes (mobile + desktop)
  - [ ] Merge when reviews approve

- **All:** Daily standup at 9:00 AM
  - What you completed yesterday
  - What you're doing today
  - Any blockers (asks Yuli for help)

**Your checklist:**
- [ ] Code follows CLAUDE.md patterns
- [ ] Tests written and passing
- [ ] Code reviewed and approved
- [ ] No merge conflicts
- [ ] Ready for QA

---

### Phase 4: QA Testing (Days 13-15)

**Owner:** Nanda (QA Architect)  
**Your involvement (by role):**

- **Backend Engineer:**
  - [ ] Deploy to staging Docker environment
  - [ ] Answer QA questions: "Does this error message make sense?" → Fix if needed
  - [ ] Run load test (if Nanda asks): `load-test.sh` in `/backend/scripts/`
  - [ ] Monitor performance metrics during testing

- **Frontend Engineer:**
  - [ ] Deploy to staging frontend
  - [ ] Test on staging (QA tests your feature)
  - [ ] Answer UI questions: "Should this button be here?" → Consult Maya if major change
  - [ ] Verify responsive behavior on multiple devices

- **QA Architect (Nanda):**
  - [ ] Execute test plan (E2E, regression, accessibility)
  - [ ] Report bugs to your Slack thread
  - [ ] Gate: "Release ready" or "Back to engineering"

**Your checklist:**
- [ ] All critical bugs fixed
- [ ] No regressions in existing features
- [ ] Performance metrics acceptable
- [ ] Ready for production

---

### Phase 5: Deployment (Day 16)

**Owner:** Hendra (DevOps)  
**Your involvement (by role):**

- **Backend Engineer:**
  - [ ] Verify database migrations run cleanly
  - [ ] Monitor backend logs during deploy (Slack #dawai-incidents)
  - [ ] Run smoke tests (critical API paths)
  - [ ] Available for emergency rollback if needed

- **Frontend Engineer:**
  - [ ] Verify build succeeds
  - [ ] Test on production staging URL
  - [ ] Available for emergency rollback if needed

- **DevOps (Hendra):**
  - [ ] Run migrations on production
  - [ ] Deploy new Docker images
  - [ ] Monitor: error rates, latency, resource usage
  - [ ] Confirm SLO metrics healthy

**Your checklist:**
- [ ] Feature live in production
- [ ] No alerts in #dawai-incidents
- [ ] Customer notifications sent
- [ ] Retrospective scheduled (if issues found)

---

## Approval Gates (When You Need Permission)

### Small Decisions (You decide)
- Code style, naming, method signatures
- Refactoring that doesn't change behavior
- Test organization, comments
- **Decision owner:** You + your peer reviewer

### Medium Decisions (Your manager decides)
- Task breakdown, sprint assignment
- How to implement architectural requirement
- Trade-off between speed and code quality
- **Decision owner:** Yuli (Engineering Manager)

### Large Decisions (Architect decides)
- New service/microservice
- Major database schema change
- Significant API contract change
- New third-party library
- **Decision owner:** Iwan (Solution Architect)
- **How:** Post in #dawai-decisions with 24-hour review window

### Business Decisions (Product Director decides)
- Feature scope, what to build
- Roadmap priority, when to build
- Release timeline
- Customer-facing changes
- **Decision owner:** Siti (Product Director)
- **How:** Rania (PM) proposes, Siti approves

---

## Escalation Paths

**You hit a blocker. Who do you ask?**

| Blocker | Ask | How | Response Time |
|---|---|---|---|
| Code review feedback confusing | Peer reviewer | Slack DM | <1 hour |
| Stuck on implementation | Iwan (Architect) | Slack #backend | <2 hours |
| Need design clarification | Maya (Designer) | Slack #design | <2 hours |
| Disagree on architecture | Iwan + Yuli | Slack conversation | <4 hours |
| Need product scope clarity | Rania (PM) | Slack @rania | <4 hours |
| Two teams have conflicting needs | Siti (Product Director) | Escalate via Yuli or Rania | <24 hours |
| Production incident (customer impact) | Yuli (Engineering Manager) | Page via Slack @here in #dawai-incidents | Immediate |

---

## Code Review Standards

### You're reviewing someone's PR

**Checklist (5-10 min):**
- [ ] Solves the stated problem
- [ ] Follows CLAUDE.md patterns (school_id filtering, audit logging, JWT usage)
- [ ] Tests written (unit or E2E)
- [ ] No obvious bugs (off-by-one, nil pointer, SQL injection)
- [ ] Variable names make sense
- [ ] No security issues (hardcoded secrets, XSS, IDOR)

**Comment style:**
- ✅ "Nice, I like how you extracted the validation logic"
- ✅ "Should this filter by school_id? Current code returns all schools' data"
- ✅ "Do we need to log this to audit_logs per CLAUDE.md?"
- ❌ "This is bad" (not specific)
- ❌ "I would do this differently" (not a blocker)

**Approval levels:**
- **APPROVE:** "Ready to merge"
- **COMMENT:** "Looks good, I have suggestions (not blockers)"
- **REQUEST CHANGES:** "Fix this before merging"

---

## On-Call Rotation (Backend/DevOps/QA Only)

### When you're on-call

**Coverage:** 6 AM - 8 PM GMT+7 (weekdays), 8 AM - 6 PM (weekends)

**During your shift:**
1. Monitor #dawai-incidents channel (read it every 30 min)
2. Check production metrics dashboard (link in Slack pinned messages)
3. Keep phone nearby (for alerts)

**If production issue detected:**
1. You're paged (Slack alert)
2. Join war room #dawai-incident (created automatically)
3. Follow incident protocol in TEAM_WORKFLOWS.md §4
4. Post status updates every 5 minutes
5. Post rootcause + action items when resolved

**After your shift:**
- Handoff to next on-call person (Slack message)
- Note any ongoing issues that need monitoring

---

## Collaboration Patterns

### Daily Standup (9:00 AM, 15 min)

**Format:**
- Each person: "Yesterday I X, today I Y, blocker Z?"
- Async option: Post in #dawai-standup if timezone conflict

**If you're stuck:**
- Mention blocker: "Waiting on DB schema approval from Iwan"
- Yuli removes blockers: "I'll ping Iwan, ETA 2 hours"

### Feature Handoff Meetings (30-45 min)

**Owner** → **Engineer**

**Meeting structure:**
- Owner: "Here's the feature requirement" (5 min)
- Engineer: "Here's my questions" (10 min)
- Owner: "Clarifications" (5 min)
- Engineer: "I can deliver" or "Need more detail" (5 min)

**After meeting:** Engineer confirms in Slack "Ready to start implementation"

### Code Review (Async, within 24 hours)

**Reviewer:** "Here's my review" → **Author:** "Fixed" or "Disagree, here's why"

**If disagreement:** Tag Iwan (architect) or Yuli (manager) to break tie

---

## First Week Milestones

### Day 1: Setup
- [ ] Local dev running (`docker compose up`)
- [ ] Can access Slack, GitHub, databases
- [ ] Attended first standup
- [ ] Read CLAUDE.md + AGENTS.md + TEAM_WORKFLOWS.md

### Day 2-3: Contribute
- [ ] Reviewed 1 PR
- [ ] Fixed 1 small bug or documentation issue
- [ ] Attended 1 function sync meeting (backend/frontend/devops/qa)

### Day 4-5: Small Feature
- [ ] Implemented small feature (bugfix or minor improvement)
- [ ] Created PR with description
- [ ] Got code review feedback
- [ ] Merged code

### Week 2: Confidence
- [ ] Contributed to medium feature (full cycle: requirements → code → QA → deploy)
- [ ] Participated in 1 architectural decision
- [ ] Know your escalation path by heart
- [ ] Can explain your role to someone new

---

## FAQ

**Q: I disagree with the architecture decision. What do I do?**  
A: Talk to Iwan (Architect) directly first. If still disagree, bring to Yuli (Manager) who can escalate. See escalation table above.

**Q: Can I skip code review if I'm confident my code is correct?**  
A: No. Code review is mandatory — not just for bugs, but for knowledge sharing and consistency.

**Q: What if I'm blocked on something outside my control?**  
A: Mention in standup. Yuli removes blockers. Don't wait >2 hours without escalating.

**Q: I deployed something and it broke production. What now?**  
A: (1) Slack alert in #dawai-incidents. (2) Immediate rollback by Hendra. (3) Post-incident review with team.  
No blame, just learn. We expect this will happen; that's why we have processes.

**Q: I want to introduce a new library / framework. How?**  
A: Discuss with Iwan (Architect) first. If approved, create decision in #dawai-decisions. Iwan + Yuli must approve before you add it to package.json.

**Q: Standup at 9 AM is hard for me (timezone issue). Can I skip?**  
A: Post async in #dawai-standup instead. You still need to communicate blockers.

**Q: I need to talk to someone on another team urgently. Who do I ping?**  
A: Ping Yuli (Engineering Manager). If Yuli isn't available, ping Siti (Product Director).

---

## Success Criteria (First Month)

- ✅ Can explain your role + decision authority to a peer
- ✅ Have merged 2-3 PRs without issues
- ✅ Participated in 1 major architectural decision or feature approval
- ✅ Helped review 3+ peers' PRs (quality feedback)
- ✅ Attended all daily standups + your function sync
- ✅ Know escalation path by heart (when to ask whom)
- ✅ Can onboard next new teammate (teach them this doc)

---

## Resources

- **Technical:** `/CLAUDE.md` (patterns, architecture, commands)
- **Team:** `/AGENTS.md` (roles + decision authority)
- **People:** `/TEAM_ROSTER.md` (contacts + skills)
- **Workflows:** `/TEAM_WORKFLOWS.md` (approval chains + handoffs + incidents)
- **GitHub:** Link from TEAM_ROSTER.md
- **Slack:** Channels listed in TEAM_ROSTER.md
- **Staging:** http://localhost:3000 (after `docker compose up`)

---

## Next: Talk to Your Direct Manager

Message your manager (from TEAM_ROSTER.md):
- **Backend:** Iwan Setyawan
- **Frontend:** Iwan Setyawan (via Dewi Kusuma)
- **DevOps:** Iwan Setyawan (via Hendra Wijaya)
- **QA:** Siti Nurhaliza
- **Product/Design:** Siti Nurhaliza
- **Anyone:** Yuli Hartono (Engineering Manager)

Say: "Finished onboarding docs. What should I work on first?"

Welcome to DAWAI! 🎻

