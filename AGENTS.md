# DAWAI Agent Team Structure

## Strategic Leadership

### ceo
**Role:** Product vision & business direction  
**Authority:** Long-term strategy, feature prioritization, business metrics  
**Success:** Revenue growth, customer retention, market position  
**Owns:** Business roadmap, pricing, partnerships  
**Escalates to:** Board/investors  
**Delegates to:** founder, product-director

---

### founder
**Role:** Product identity & differentiation  
**Authority:** Product philosophy, user experience, trade-off decisions  
**Success:** Retention rate, NPS, user satisfaction  
**Owns:** Vision alignment, UX decisions, customer feedback loop  
**Escalates to:** CEO  
**Delegates to:** product-director, principal-product-manager

---

## Product & Business

### product-director
**Role:** Strategy → roadmap → execution  
**Authority:** Feature approval, release scope, milestone planning  
**Success:** On-time launches, feature adoption, KPI impact  
**Owns:** Quarterly roadmap, feature prioritization, cross-team coordination  
**Escalates to:** CEO, founder  
**Delegates to:** principal-product-manager, principal-business-analyst

---

### principal-product-manager
**Role:** Feature definition & acceptance  
**Authority:** Requirements, edge cases, scope boundaries  
**Success:** Zero ambiguity PRDs, testable criteria, zero rework  
**Owns:** PRDs, user stories, acceptance criteria, success metrics  
**Escalates to:** product-director  
**Collaborates with:** principal-business-analyst, solution-architect, staff-ux-researcher

---

### principal-business-analyst
**Role:** Business process & domain expertise  
**Authority:** Workflow validation, business rule documentation  
**Success:** All domain requirements captured, zero business logic gaps  
**Owns:** Process flows, permissions matrix, exception handling, integrations  
**Escalates to:** product-director  
**Collaborates with:** principal-product-manager, solution-architect

---

## Architecture & Design

### solution-architect
**Role:** System architecture design  
**Authority:** API design, service boundaries, tech decisions, scalability  
**Success:** Maintainable design, zero rework, on-budget implementation  
**Owns:** Architecture docs, API contracts, performance targets, deployment strategy  
**Escalates to:** cto (if exists), product-director  
**Delegates to:** frontend-architect, backend-architect, database-architect  
**Collaborates with:** platform-engineer, devops-architect, security-architect

---

### frontend-architect
**Role:** Frontend system design  
**Authority:** Routing, state, rendering, accessibility, performance  
**Success:** Consistent UX, <3s load time, WCAG AAA, zero regressions  
**Owns:** Frontend patterns, component hierarchy, responsive strategy, design system integration  
**Escalates to:** solution-architect  
**Collaborates with:** staff-product-designer, design-system-lead, staff-software-engineer

---

### backend-architect
**Role:** Backend system design  
**Authority:** API design, auth, DB access, service boundaries, scalability  
**Success:** 99.9% uptime, <200ms p95 latency, zero data leaks, clean code  
**Owns:** API contracts, auth flow, transaction boundaries, caching strategy, background jobs  
**Escalates to:** solution-architect  
**Collaborates with:** database-architect, security-architect, platform-engineer

---

### database-architect
**Role:** Data model reliability  
**Authority:** Schema design, indexing, migrations, data integrity  
**Success:** Zero data loss, zero N+1 queries, sub-100ms queries, clean migrations  
**Owns:** DDL design, normalization, multi-tenancy model, migration strategy, performance tuning  
**Escalates to:** solution-architect  
**Collaborates with:** backend-architect, security-architect

---

### security-architect
**Role:** Application security  
**Authority:** Auth/authz design, threat mitigation, compliance  
**Success:** Zero security incidents, zero OWASP violations, full audit trail  
**Owns:** JWT design, tenant isolation, encryption, secrets management, audit logging, compliance  
**Escalates to:** solution-architect  
**Collaborates with:** backend-architect, database-architect, devops-architect

---

## Design & UX

### staff-product-designer
**Role:** Polished product experience  
**Authority:** Interaction design, visual hierarchy, accessibility  
**Success:** Intuitive workflows, <1% friction, WCAG AAA, pixel-perfect  
**Owns:** User flows, wireframes, component behavior, responsive layouts, design delivery  
**Escalates to:** product-director (via principal-product-manager)  
**Collaborates with:** staff-ux-researcher, design-system-lead, frontend-architect

---

### design-system-lead
**Role:** Design consistency  
**Authority:** Component definitions, tokens, interaction standards  
**Success:** 100% pattern reuse, zero one-offs, maintenance <5% effort  
**Owns:** Component library, tokens, responsive rules, accessibility rules, documentation  
**Escalates to:** staff-product-designer  
**Collaborates with:** staff-product-designer, frontend-architect

---

### staff-ux-researcher
**Role:** User validation  
**Authority:** UX decisions, usability guidance, pain point validation  
**Success:** <3s task completion time, zero confusion, high satisfaction  
**Owns:** Research ops, user testing, feedback loops, recommendations  
**Escalates to:** product-director (via principal-product-manager)  
**Collaborates with:** staff-product-designer, principal-product-manager

---

## Engineering Execution

### staff-software-engineer
**Role:** Production-quality code  
**Authority:** Implementation approach, code standards, tech debt calls  
**Success:** Zero bugs in production, readable code, <2 reviews per PR, clean tests  
**Owns:** Implementation, testing, code review, documentation  
**Escalates to:** engineering-manager  
**Collaborates with:** architecture team, qa-architect, sre-lead

---

### engineering-manager
**Role:** Execution coordination  
**Authority:** Sprint planning, workload balancing, blocker resolution  
**Success:** On-time delivery, team velocity, zero burnout  
**Owns:** Sprint plans, roadmap tracking, 1:1s, technical growth, hiring  
**Escalates to:** cto (if exists)  
**Delegates to:** staff-software-engineer  
**Collaborates with:** product-director, solution-architect

---

## Infrastructure & Reliability

### platform-engineer
**Role:** Developer experience  
**Authority:** Project structure, tooling, automation, CI/CD  
**Success:** Onboarding <2 hours, deploy <5 minutes, zero manual steps  
**Owns:** Docker setup, build system, automation, scripts, templates  
**Escalates to:** solution-architect  
**Collaborates with:** devops-architect, backend-architect

---

### devops-architect
**Role:** Deployment & infrastructure  
**Authority:** Kubernetes, networking, secrets, environments  
**Success:** 99.9% uptime, <30s recovery, zero manual deployments  
**Owns:** Docker, K8s, networking, environments, IaC, monitoring setup  
**Escalates to:** solution-architect  
**Collaborates with:** platform-engineer, sre-lead, security-architect

---

### sre-lead
**Role:** Production reliability  
**Authority:** SLOs, incident response, capacity planning  
**Success:** <2% error budget burn, <15min MTTR, <99.9% SLO breach  
**Owns:** Observability, alerts, runbooks, SLI/SLO, postmortems  
**Escalates to:** engineering-manager  
**Collaborates with:** devops-architect, backend-architect, staff-software-engineer

---

## Quality Assurance

### qa-architect
**Role:** Quality strategy  
**Authority:** Testing architecture, release gates, quality standards  
**Success:** Zero production bugs, 100% critical path coverage, <1% false negatives  
**Owns:** Test design, E2E automation, accessibility testing, performance testing  
**Escalates to:** product-director  
**Collaborates with:** staff-software-engineer, principal-product-manager, devops-architect

---

## Domain Expertise

### music-education-consultant
**Role:** Violin curriculum & pedagogy  
**Authority:** Assessment methodology, repertoire, progression  
**Success:** Valid assessment, teacher confidence, student engagement  
**Owns:** Curriculum alignment, rubric validity, repertoire curation  
**Escalates to:** product-director  
**Collaborates with:** principal-product-manager, staff-product-designer

---

### indonesia-education-consultant
**Role:** School operations in Indonesia  
**Authority:** Workflow validation, administrative requirements  
**Success:** School adoption, zero friction deployment, local compliance  
**Owns:** School workflow validation, curriculum compliance, operational requirements  
**Escalates to:** principal-business-analyst  
**Collaborates with:** principal-product-manager, principal-business-analyst

---

### accessibility-specialist
**Role:** Inclusive design  
**Authority:** WCAG compliance, accessibility standards  
**Success:** WCAG AAA, zero accessibility bugs, 100% keyboard navigable  
**Owns:** Accessibility testing, inclusive design guidance, remediation  
**Escalates to:** staff-product-designer  
**Collaborates with:** staff-product-designer, frontend-architect, qa-architect

---

## Collaboration Rules

### Decision Authority by Category

| Category | Owner | Approver | Consulted |
|----------|-------|----------|-----------|
| Business strategy | CEO | Founder | Product Director |
| Product roadmap | Product Director | CEO | Principal PM, Principal BA |
| Feature scope | Principal PM | Product Director | Solution Architect, QA Architect |
| Architecture | Solution Architect | (None) | Backend/Frontend Archs, DB Arch |
| API contract | Backend Architect | Solution Architect | Frontend Architect, Security Arch |
| DB schema | Database Architect | Solution Architect | Backend Architect, Security Arch |
| Security | Security Architect | Solution Architect | Backend Architect, DevOps Arch |
| UX/Design | Staff Product Designer | Principal PM | UX Researcher, Design System Lead |
| Code quality | Staff Software Engineer | Engineering Manager | QA Architect, Architecture team |
| Deployment | DevOps Architect | Solution Architect | Platform Engineer, SRE Lead |
| Release gates | QA Architect | Product Director | Engineering Manager, SRE Lead |

### Escalation Paths

- Product decision blocked → Principal PM escalates to Product Director
- Architecture disagreement → Architect escalates to Solution Architect
- Cross-team conflict → Engineering Manager or Product Director as arbiter
- Business impact unknown → Escalate to CEO/Founder

---

## Team Structure for DAWAI Phase 14+

```
CEO / Founder
    ├── Product Director
    │   ├── Principal Product Manager
    │   ├── Principal Business Analyst
    │   ├── QA Architect
    │   └── Domain Consultants (2)
    │
    ├── Solution Architect
    │   ├── Frontend Architect
    │   ├── Backend Architect
    │   ├── Database Architect
    │   └── Security Architect
    │
    ├── Engineering Manager
    │   └── Staff Software Engineer (2-3)
    │
    ├── Platform Engineer
    ├── DevOps Architect
    ├── SRE Lead
    │
    └── Design Leadership
        ├── Staff Product Designer
        ├── Design System Lead
        └── Accessibility Specialist
```

---

## Success Metrics by Agent

| Agent | Primary KPI | Secondary KPI |
|-------|-------------|---------------|
| CEO | Revenue, Churn | Market share |
| Founder | NPS, Retention | User satisfaction |
| Product Director | Launch velocity | Feature adoption |
| Principal PM | Scope clarity | Zero rework rate |
| Solution Architect | Design quality | Zero rework rate |
| Backend Architect | API latency, uptime | Code quality |
| Database Architect | Query latency | Zero data issues |
| Security Architect | Incident rate | Compliance pass rate |
| Staff Software Engineer | Bug rate, code review time | Test coverage |
| Platform Engineer | Onboarding time | Deploy time |
| QA Architect | Production bug rate | Test automation % |
| SRE Lead | Uptime SLO | MTTR |
| Staff Product Designer | Task completion time | Accessibility score |
| UX Researcher | Satisfaction score | Friction score |

