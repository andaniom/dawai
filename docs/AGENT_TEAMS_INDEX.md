# Agent Teams Master Reference Index

Complete documentation for building, managing, and debugging agent teams in Claude Code.

---

## Quick Navigation

### 🚀 Getting Started
- **[Quick Start](agent-teams-quick-start.md)** — Setup, spawn your first team, control panel reference
- **[Reference Guide](agent-teams-reference.md)** — Complete feature overview, architecture, best practices

### 📋 Implementation Guides
- **[Patterns & Use Cases](agent-teams-patterns.md)** — Copy/adapt prompts for code review, debugging, refactoring, testing
- **[DAWAI-Specific Guide](agent-teams-dawai.md)** — Phase 13 & 14 sprints, security audits, multi-tenant patterns
- **[Advanced Techniques](agent-teams-advanced.md)** — Troubleshooting, hooks, session resumption, performance tuning

---

## By Task Type

### Code Review
- [Parallel Domain Review](agent-teams-patterns.md#parallel-domain-review) — split by expertise
- [Specialized Language Review](agent-teams-patterns.md#specialized-language-review) — Go, TS, SQL reviewers
- [Full example for DAWAI](agent-teams-dawai.md#code-review-assessment-flow-refactor)

### Bug Investigation
- [Competing Hypotheses Debate](agent-teams-patterns.md#competing-hypotheses-debate) — multiple theories in parallel
- [Multi-Layer Debugging](agent-teams-patterns.md#multi-layer-debugging) — backend, DB, network, load
- [DAWAI data leak example](agent-teams-dawai.md#bug-investigation-multi-school-data-leak)

### Feature Development
- [Parallel Module Implementation](agent-teams-patterns.md#parallel-module-implementation) — DB + API + frontend + tests
- [Cross-Tenant Safety Implementation](agent-teams-patterns.md#cross-tenant-safety-implementation)
- [DAWAI rubric feature example](agent-teams-dawai.md#feature-development-new-rubric-components)

### Refactoring
- [Parallel Refactor with Verification](agent-teams-patterns.md#parallel-refactor-with-verification)

### Testing
- [Parallel Test Suite](agent-teams-patterns.md#parallel-test-suite-implementation)
- [Offline Sync Testing](agent-teams-dawai.md#testing-offline-sync-scenarios)

### Security & Compliance
- [Security Audit](agent-teams-patterns.md#multi-angle-security-audit)
- [DAWAI Phase 13 Hardening](agent-teams-dawai.md#phase-13-security-hardening-sprint)

---

## By Skill Level

### Beginner (First Team)
1. Read: [Quick Start](agent-teams-quick-start.md)
2. Try: "Spawn 3 teammates to review PR #X"
3. Monitor: Use agent panel, send messages
4. Learn from [Basic Patterns](agent-teams-patterns.md#code-review-patterns)

### Intermediate (Complex Teams)
1. Read: [Reference Guide](agent-teams-reference.md) (full architecture)
2. Study: [DAWAI Examples](agent-teams-dawai.md) (real projects)
3. Try: Multi-phase team with task dependencies
4. Reference: [Patterns](agent-teams-patterns.md) when spawning

### Advanced (Expert Teams)
1. Understand: [Architecture Deep Dive](agent-teams-advanced.md#architecture-deep-dive)
2. Implement: [Quality Enforcement](agent-teams-advanced.md#quality-enforcement-with-hooks) with hooks
3. Optimize: [Performance Tuning](agent-teams-advanced.md#performance-tuning)
4. Recover: [Troubleshooting](agent-teams-advanced.md#handling-common-issues)

---

## Configuration Reference

### One-Time Setup
```json
~/.claude/settings.json

{
  "env": {
    "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1"
  },
  "teammateMode": "auto",
  "defaultTeammateModel": "claude-haiku-4-5"
}
```

See [Reference Guide § Enable Agent Teams](agent-teams-reference.md#enable-agent-teams)

### Display Modes
- **in-process** (default) — all teammates in single terminal
- **auto** — split panes if tmux/iTerm2 available
- **tmux** — force split panes with tmux
- **iterm2** — force split panes with iTerm2

### Hooks (Quality Gates)
- `TeammateIdle` — runs when teammate finishes turn
- `TaskCreated` — runs when task being created
- `TaskCompleted` — runs when task marked complete

See [Advanced § Quality Enforcement](agent-teams-advanced.md#quality-enforcement-with-hooks)

---

## Common Commands

### Spawn a Team
```text
Spawn 3 teammates to review PR #142
```

### Message a Teammate
```text
Tell researcher to focus on database logs
```

### Shut Down a Teammate
```text
Ask researcher to shut down
```

### Assign Work
```text
Give task #3 to security-reviewer
```

### Check Progress
- Arrow keys: select teammate in panel
- Enter: view teammate, type to message
- Escape: exit view
- Ctrl+T: toggle task list

---

## Troubleshooting Quick Links

| Problem | Solution |
|---------|----------|
| Teammates not appearing | [Quick Start § Troubleshoot](agent-teams-quick-start.md#troubleshoot) |
| Too many permission prompts | [Reference § Permissions](agent-teams-reference.md#permissions) |
| Teammate stops on error | [Advanced § Stopped on Error](agent-teams-advanced.md#issue-teammate-stopped-working-on-error) |
| Two teammates edited same file | [Advanced § File Conflicts](agent-teams-advanced.md#issue-two-teammates-edited-same-file) |
| Task appears stuck | [Advanced § Task Stuck](agent-teams-advanced.md#issue-task-marked-complete-but-work-incomplete) |
| Orphaned tmux session | [Quick Start § Troubleshoot](agent-teams-quick-start.md#troubleshoot) |
| Mailbox corruption | [Advanced § Mailbox Corruption](agent-teams-advanced.md#issue-mailbox-corruption-malformed-json) |
| Session resumption | [Advanced § Session Resumption](agent-teams-advanced.md#session-resumption-strategy) |

---

## DAWAI-Specific Checklists

### Phase 13: Security Hardening
- Spawn command: [Phase 13 Launch](agent-teams-dawai.md#phase-13-security-hardening-sprint)
- Checklist: [Before Spawning](agent-teams-dawai.md#checklist-before-spawning-team-on-dawai)
- Gotchas: [Common Pitfalls](agent-teams-dawai.md#common-dawai-team-pitfalls)

### Phase 14: PWA + Deployment
- Spawn command: [Phase 14 Launch](agent-teams-dawai.md#phase-14-pwa--deployment)

### Feature Development
- Spawn command: [Rubric Feature](agent-teams-dawai.md#feature-development-new-rubric-components)
- Isolation verification: [Safety Implementation](agent-teams-dawai.md#cross-tenant-safety-implementation)

---

## File Structure

```
docs/
├── AGENT_TEAMS_INDEX.md              (this file)
├── agent-teams-quick-start.md        (5 min read, get started)
├── agent-teams-reference.md          (20 min read, complete guide)
├── agent-teams-patterns.md           (30 min read, copy/adapt prompts)
├── agent-teams-dawai.md              (10 min read, DAWAI-specific)
└── agent-teams-advanced.md           (20 min read, troubleshooting)
```

Total: ~85 minutes reading time, 100+ copy/paste prompts ready to use.

---

## Key Concepts

**Team Lead** — main Claude Code session that spawns teammates

**Teammates** — separate full Claude sessions, fully independent, can message each other

**Task List** — shared work items that teammates claim and complete (persists across resume)

**Mailbox** — JSON message system for inter-agent communication

**In-Process Mode** — default, all teammates in single terminal

**Split Panes** — tmux or iTerm2, each teammate gets own pane

---

## When to Use Teams

✅ **Use teams when:**
- Work is parallelizable (features span frontend/backend/tests)
- Teammates need to collaborate (code review, debate)
- Task is complex (3+ days, multiple phases)
- Token budget allows (teams are expensive)

❌ **Use subagents instead when:**
- Work is sequential (not parallel)
- Single focused task (1-2 hours)
- No inter-worker communication needed

See [Reference § Compare with Subagents](agent-teams-reference.md#compare-with-subagents)

---

## External References

- **Claude Code Docs:** https://code.claude.com/docs/en/agent-teams
- **Token Costs:** https://code.claude.com/docs/en/costs#agent-team-token-costs
- **Subagents:** https://code.claude.com/docs/en/sub-agents
- **Hooks:** https://code.claude.com/docs/en/hooks
- **Skills:** https://code.claude.com/docs/en/skills

---

## Quick Stats

| Metric | Value |
|--------|-------|
| Guides | 5 |
| Copy/Paste Prompts | 15+ |
| Code Examples | 30+ |
| DAWAI Examples | 8+ |
| Troubleshooting Sections | 20+ |
| Checklists | 8+ |
| Total Content | ~9,000 words |

---

## How to Use This Index

1. **First time:** Read [Quick Start](agent-teams-quick-start.md)
2. **Ready to spawn:** Pick your task type in [By Task Type](#by-task-type), copy prompt
3. **Team struggling:** Go to [Troubleshooting](#troubleshooting-quick-links)
4. **Need to understand:** Read [Reference Guide](agent-teams-reference.md)
5. **DAWAI work:** Jump to [DAWAI Guide](agent-teams-dawai.md)
6. **Advanced needs:** Dive into [Advanced](agent-teams-advanced.md)

---

## Feedback & Updates

As you use agent teams on DAWAI, update this documentation:
- Found a good pattern? Add to [Patterns](agent-teams-patterns.md)
- Hit a new gotcha? Add to [Advanced § Common Issues](agent-teams-advanced.md)
- DAWAI-specific learnings? Update [DAWAI Guide](agent-teams-dawai.md)

Keep this index as the single source of truth for team development.
