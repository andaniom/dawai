# Agent Teams Quick Start Checklist

Fast reference for spawning and managing your first agent team.

## Setup (One-time)

- [ ] Enable: Add `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` to `~/.claude/settings.json` or shell env
- [ ] (Optional) Choose display mode: `"teammateMode": "auto"` in settings.json for split panes
- [ ] (Optional) Set default teammate model: `"defaultTeammateModel": "claude-haiku-4-5"`

## Spawn a Team

### Minimal

```text
Spawn 3 teammates to review PR #142
```

### With Specification

```text
Spawn 4 teammates, call them:
- security-reviewer: check token handling
- performance-reviewer: check query efficiency  
- test-reviewer: check coverage
- maintainability-reviewer: check code clarity

Each review independently and report findings.
```

### With Subagent Type

```text
Spawn a teammate using the security-reviewer agent type to audit the auth module.
```

### With Plan Approval

```text
Spawn an architect to refactor the database schema.
Require plan approval before making changes.
Only approve if the plan includes migration reversibility.
```

## Control Team

### View & Message

- **In-process mode** (default):
  - Arrow keys: select teammate in panel
  - Enter: view teammate, type to message
  - Escape: exit teammate view
  - x: stop teammate
  - Ctrl+T: toggle task list

- **Split panes mode**:
  - Click pane to interact
  - Each teammate has full terminal

### Assign Work

```text
Tell the team to:
1. Security reviewer: check endpoints for school_id leaks
2. Performance reviewer: benchmark query patterns
3. Test reviewer: add cross-tenant access tests

Break into 3 tasks. Assign task 1 to security-reviewer.
```

### Message Teammate

```text
Tell researcher to focus on database logs instead of network logs.
```

### Shut Down Teammate

```text
Ask researcher to shut down.
```

(Teammate approves/rejects. If rejected, ask why and provide new context.)

## Monitor Progress

- [ ] Check teammate panel for idle/working status
- [ ] Read through teammate outputs regularly
- [ ] Send course corrections if approach drifts
- [ ] Wait for all teammates before proceeding

## Common Patterns

### Code Review Split by Domain

```text
Spawn 3 reviewers:
- Security: token handling, auth, data leaks
- Performance: DB queries, caching, complexity
- Tests: coverage, edge cases, mock quality

Each review independently.
```

### Competing Hypotheses

```text
Spawn 4 teammates to investigate crash:
- Hypothesis 1: memory leak
- Hypothesis 2: connection timeout
- Hypothesis 3: uncleaned event listeners
- Hypothesis 4: race condition in state

Have them debate to disprove each other's theories.
```

### Parallel Feature Development

```text
Spawn 3 teammates:
- Database team: schema + migrations
- API team: endpoints + handlers
- Client team: components + state

Each owns separate files. Coordinate via task list.
```

## Troubleshoot

| Issue | Solution |
|-------|----------|
| No teammates appearing | They might be hidden (idle). Send one a message by name. |
| Too many permission prompts | Pre-approve operations in permission settings. |
| Teammate stops on error | Message it with recovery steps or spawn replacement. |
| Lead doing work instead of delegating | Tell lead: "Wait for teammates to finish first." |
| Task blocked | Check if work is done, update task status manually. |
| Two teammates editing same file | Split so each teammate owns different files. |
| Teammates not talking to each other | Message them directly: "Share findings with security-reviewer." |

## Tips

✅ **Include context in spawn**: "Review auth module at src/auth/. Focus on JWT handling, token storage in httpOnly cookies, 7-day expiry, blacklist on logout."

✅ **Name teammates**: makes them easy to reference later

✅ **3-5 teammates optimal**: balances parallelism with coordination

✅ **5-6 tasks per teammate**: keeps everyone busy

✅ **Avoid same-file edits**: prevents overwrites

✅ **Start with research/review**: lower complexity, clear value

✅ **Monitor in-process**: check in regularly, steer approaches

❌ **Don't spawn huge teams**: coordination overhead kills productivity

❌ **Don't leave team unattended forever**: check in periodically

❌ **Don't let lead do work unilaterally**: use task list to delegate

❌ **Don't expect teammates to follow entire conversation**: pass task-specific context

## For DAWAI

### Phase 13: Security Audit

```text
Spawn 3 security teammates:
- Query auditor: grep handlers for school_id filters, verify every query
- Access tester: attempt cross-tenant access, bypass attempts
- Token verifier: validate JWT isolation, blacklist logic, expiry

Report vulnerabilities with severity ratings.
```

### Parallel Module Work

```text
Spawn 4 teammates for feature X:
- Database: migrations, schema
- Backend: API endpoints, business logic
- Frontend: components, state management
- Tests: unit, integration, E2E

Each own files in separate directories. No conflicts.
```

### Multi-Perspective Review

```text
Spawn 4 reviewers for new assessment flow:
- User flow: is UX smooth end-to-end?
- Data isolation: school_id filters correct?
- Performance: any N+1 queries or memory leaks?
- Accessibility: WCAG compliance, keyboard nav?

Each review independently, report issues with severity.
```

---

**Full reference**: See `agent-teams-reference.md` for complete guide.
