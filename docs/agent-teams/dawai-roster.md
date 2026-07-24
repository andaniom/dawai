# DAWAI Team Roster — Agent Teams Integration

Team structure mapped to agent team roles for spawning parallel teams.

---

## Agent Team Roles by Function

### Strategic Leadership

| Person | Title | Agent Team Role | Authority | Escalation |
|--------|-------|-----------------|-----------|------------|
| Andi Nugroho | CEO/Founder | ceo | Strategy, roadmap, partnerships | Board/investors |
| Andi Nugroho | CEO/Founder | founder | Vision, UX philosophy, trade-offs | CEO |

### Product & Business

| Person | Title | Agent Team Role | Authority | Escalation |
|--------|-------|-----------------|-----------|------------|
| Siti Nurhaliza | Product Director | product-director | Feature approval, roadmap, scope | CEO, founder |
| Rania Pratiwi | Principal PM | principal-pm | PRD, acceptance criteria | product-director |
| Bambang Sutrisno | Principal BA | principal-ba | Process, workflows, rules | product-director |

### Architecture & Design

| Person | Title | Agent Team Role | Authority | Escalation |
|--------|-------|-----------------|-----------|------------|
| Iwan Setyawan | Solution Arch / Backend Arch | solution-architect | API design, service boundaries, tech decisions | CEO (tech) |
| Dewi Kusuma | Frontend Arch / Staff Eng | frontend-architect | Frontend system, patterns, a11y | solution-architect |
| Ratna Wijaya | Database Architect | database-architect | Schema, migrations, indexing | solution-architect |
| Arief Gunawan | Security Architect | security-architect | JWT, isolation, compliance | solution-architect |

### Design & UX

| Person | Title | Agent Team Role | Authority | Escalation |
|--------|-------|-----------------|-----------|------------|
| Maya Handoko | Staff Designer / DS Lead | staff-designer | Interaction, visual hierarchy, a11y | product-director |
| Maya Handoko | Staff Designer / DS Lead | design-system-lead | Components, tokens, patterns | staff-designer |
| Lena Pratama | UX Researcher / A11y Spec | ux-researcher | Research, usability, WCAG | product-director |

### Engineering Execution

| Person | Title | Agent Team Role | Authority | Escalation |
|--------|-------|-----------------|-----------|------------|
| Ahmad Hidayat | Backend Engineer | staff-software-engineer | Implementation, code standards | engineering-manager |
| Yuli Hartono | Engineering Manager | engineering-manager | Sprint planning, blocking | solution-architect + product-director |

### Infrastructure & Reliability

| Person | Title | Agent Team Role | Authority | Escalation |
|--------|-------|-----------------|-----------|------------|
| Hendra Wijaya | DevOps / Platform Eng | platform-engineer | Docker, build, automation | solution-architect |
| Hendra Wijaya | DevOps / Platform Eng | devops-architect | K8s, networking, environments | solution-architect |
| Nanda Kusuma | QA / SRE Lead | sre-lead | SLOs, incidents, observability | engineering-manager |
| Nanda Kusuma | QA / SRE Lead | qa-architect | Testing, release gates | product-director |

### Domain Expertise

| Person | Title | Agent Team Role | Authority | Escalation |
|--------|-------|-----------------|-----------|------------|
| Prof. Budi Santosa | Music Ed Consultant | music-consultant | Curriculum, rubric validity | product-director |
| Dr. Nur Cahyadi | Indonesia Ed Consultant | indonesia-consultant | School ops, compliance | product-director |

---

## Spawn Commands by Role

### Spawn {role} for Task

```text
Spawn a teammate using the ceo agent type to evaluate feature strategy.
Authority: product roadmap, business metrics, partnership decisions.
Escalation path: board/investors if resource conflict.
```

**Replace `ceo` with any role slug from table above.**

### Spawn Cross-Functional Team

```text
Spawn 4 teammates:
1. Intended as product-director: define scope + success metrics
2. solution-architect: design API + data model
3. staff-designer: create component mockups
4. staff-software-engineer: estimate implementation effort

Each own separate work product. Share findings when done.
```

### Spawn Incident Response Team

```text
Production outage suspected. Spawn 3 teammates:
- staff-software-engineer: trace logs, find root cause
- devops-architect: check system resources, recent changes
- sre-lead: monitor SLOs, coordinate rollback if needed

Coordinate immediately. Report root cause + fix within 15 min.
```

---

## Dual Roles (Cost Optimization)

Some people hold 2 agent roles. When spawning team referencing them:

```text
Give the architecture task to solution-architect (Iwan).
Note: Iwan also holds backend-architect role — do not double-assign.
```

| Person | Dual Roles | When to Spawn Each |
|--------|-----------|-------------------|
| Andi Nugroho | ceo, founder | ceo for strategy; founder for product vision |
| Iwan Setyawan | solution-architect, backend-architect | solution-architect for system design; backend-architect for API work |
| Dewi Kusuma | frontend-architect, staff-software-engineer | frontend-architect for patterns; staff-software-engineer for implementation |
| Maya Handoko | staff-designer, design-system-lead | staff-designer for mockups; design-system-lead for component tokens |
| Hendra Wijaya | platform-engineer, devops-architect | platform-engineer for automation; devops-architect for infra |
| Nanda Kusuma | qa-architect, sre-lead | qa-architect for test strategy; sre-lead for reliability |
| Lena Pratama | ux-researcher, accessibility-specialist | ux-researcher for usability; accessibility-specialist for WCAG |

---

## Timezone & Availability

- **Timezone:** GMT+7 (Jakarta)
- **Business hours:** 8:00 AM - 5:00 PM GMT+7
- **On-call:** Backend Engineer + DevOps for production incidents
- **Async protocol:** 24-hour approval window in Slack #dawai-decisions

When scheduling team sprints, respect 8-5 Jakarta hours for sync work.

---

## FTE Capacity

| Category | Count |
|----------|-------|
| Full-time (100%) | 9 |
| Part-time (60-80%) | 3 |
| Total FTE | ~11 |
| Backend/Architecture | 3 |
| Frontend | 2 |
| DevOps/Platform | 2 |
| QA/Testing | 1 |
| Design/UX | 2 |
| Product/Business | 3 |
| Consulting | 3 |

**Implication:** Max 3-4 parallel teammates for DAWAI-specific sprints. Beyond that exhausts human equivalents.
