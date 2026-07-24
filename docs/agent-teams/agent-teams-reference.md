# Agent Teams Master Reference Guide

Complete reference for building effective agent teams in Claude Code. Use this to spawn, coordinate, and manage parallel teams for complex tasks.

---

## Quick Start

### Enable Agent Teams

```json
{
  "env": {
    "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1"
  }
}
```

### Spawn Your First Team

```text
I'm designing a CLI tool. Spawn three teammates:
- UX designer for user experience
- Architect for technical design  
- Devil's advocate to poke holes
```

Claude spawns them, creates a shared task list, and coordinates work.

---

## When to Use Agent Teams

### ✅ Best For (Parallel Exploration)

- **Research & review**: Multiple reviewers check different aspects simultaneously (security, performance, tests)
- **New features**: Each teammate owns separate modules without stepping on each other
- **Competing hypotheses**: Teams test different theories in parallel, debate findings
- **Cross-layer work**: Frontend, backend, database changes owned by different teammates

### ❌ Avoid For

- **Sequential tasks**: work that must happen step-by-step
- **Same-file edits**: two teammates editing the same file = overwrites
- **Simple focused work**: overhead exceeds benefit
- **Many dependencies**: teams with heavy inter-task blocking

### Agent Teams vs Subagents

| Aspect | Subagents | Agent Teams |
|--------|-----------|------------|
| **Communication** | Report back to main agent only | Teammates message each other directly |
| **Coordination** | Main agent manages all work | Shared task list + self-coordination |
| **Context** | Shared with caller | Fully independent per teammate |
| **Best for** | Quick focused workers | Complex work requiring collaboration |
| **Token cost** | Lower (results summarized) | Higher (full separate instances) |

**Use subagents** when you need quick helpers. **Use teams** when teammates need to discuss, challenge, and coordinate.

---

## Spawning Teams

### Basic Spawn

```text
Spawn 3 teammates to refactor the auth module
```

Claude decides number. Be explicit for control:

```text
Spawn exactly 4 teammates: frontend expert, backend expert, 
security reviewer, test specialist
```

### With Subagent Definitions

Reference any subagent type when spawning:

```text
Spawn a teammate using the security-reviewer agent type to audit auth.
```

Teammate honors subagent's `tools` allowlist, `model`, and instructions.

### With Specific Instructions

```text
Spawn a teammate with this prompt: "Review the authentication module 
at src/auth/ focusing on token handling, session management, and input 
validation. Report with severity ratings."
```

Include task-specific context—teammates don't inherit lead's conversation history.

### Model & Effort Control

- **Default**: teammates use system default model (set in `/config` → "Default teammate model")
- **Explicit**: `"Use Sonnet for each teammate"`
- **Effort**: teammates inherit lead's effort level from v2.1.186+

---

## Team Architecture

### File Structure

```
~/.claude/teams/{team-name}/
├── config.json         # Runtime state (auto-managed, don't edit)
├── inboxes/
│   ├── lead.json       # Messages to lead
│   └── teammate-1.json # Messages to teammate
└── ...

~/.claude/tasks/{team-name}/
├── task-1.json
├── task-2.json
└── ...
```

**Team config deleted on session end.** Task list persists for resumed sessions.

### Components

| Component | Role |
|-----------|------|
| **Team lead** | Main session that spawns & coordinates |
| **Teammates** | Separate Claude instances, fully independent |
| **Task list** | Shared work items teammates claim & complete |
| **Mailbox** | JSON message system for inter-agent communication |

Lead auto-generates team name: `session-` + first 8 chars of session ID.

---

## Control & Coordination

### Display Modes

**In-process** (default):
- All teammates in single terminal
- Use arrow keys in agent panel to select teammate
- Press Enter to view & message
- Works everywhere, no setup

**Split panes** (tmux/iTerm2):
- Each teammate gets own pane
- See everyone's output simultaneously
- Click pane to interact directly

Set in `settings.json`:

```json
{
  "teammateMode": "auto"      // auto-detect tmux/iTerm2
  // or "tmux", "iterm2", "in-process"
}
```

Or per-session:

```bash
claude --teammate-mode auto
```

### Interacting with Teammates

**From lead terminal:**
- Arrow keys: select teammate in agent panel
- Enter: view teammate session, type to message
- Escape: close teammate view
- Ctrl+T: toggle task list
- x: stop selected teammate

**Direct messaging:** "Tell researcher to investigate X"

**Team shutdown:** "Ask researcher to shut down" → teammate approves/rejects

### Task Management

Tasks coordinate work across team. States: pending → in-progress → completed.

**Dependencies:** A task can depend on other tasks. Blocked tasks unblock when dependencies finish.

**Claiming:**
- Lead assigns: "Give task #3 to security-reviewer"
- Self-claim: teammate picks next unassigned, unblocked task
- File locking prevents race conditions

**Assigning work:**

```text
Break this into 5 tasks:
1. Auth module refactor (for backend-expert)
2. API client update (for frontend-expert)
3. Test suite updates (for test-specialist)
4. Documentation (for doc-writer)
5. Integration testing (for qa-specialist)

Start with task 1, assign task 2 to frontend-expert.
```

Lead creates tasks; teammates self-claim available work.

---

## Communication Patterns

### Automatic Delivery

Messages between teammates arrive immediately—no polling. Lead notified when teammates idle or fail.

### Naming Teammates

Claude auto-names: Teammate 1, Researcher, Backend Expert, etc.

Specify names for predictable references:

```text
Spawn three teammates, call them:
- researcher (investigates bugs)
- implementer (fixes code)
- tester (validates fixes)
```

Then message by name: "Tell researcher to focus on database logs"

### Message Types

- **Lead → Teammate**: instructions, clarifications, new work
- **Teammate → Lead**: progress updates, blockers, requests
- **Teammate ↔ Teammate**: findings, debate, coordination (direct messaging)

---

## Best Practices

### 1. Give Enough Context

Teammates load project settings (CLAUDE.md, MCP servers, skills) but not lead's conversation history.

Include task-specific details in spawn prompt:

```text
Spawn security reviewer with prompt: "Review auth module at src/auth/.
Focus on token handling in JWT flow: tokens stored in httpOnly cookies,
7-day expiry, JTI-based blacklist on logout. Check for timing attacks,
token reuse, session fixation."
```

### 2. Right-Size Team

- **3-5 teammates**: optimal for most tasks
- **Coordination overhead**: scales with team size
- **Token cost**: linear per teammate
- **5-6 tasks per teammate**: keeps everyone productive

3 focused teammates > 5 scattered ones.

### 3. Size Tasks Appropriately

- **Too small**: coordination overhead > benefit
- **Too large**: teammates work too long without check-ins, risk wasted effort
- **Just right**: self-contained deliverable (function, test file, feature, review)

Ask lead: "Split work into smaller tasks" if insufficient granularity.

### 4. Avoid File Conflicts

Two teammates editing same file = overwrites.

Break work so each teammate owns different files:

```text
- Frontend-expert owns src/components/
- Backend-expert owns backend/handlers/
- Tester owns tests/
```

### 5. Monitor & Steer

Don't let team run unattended indefinitely. Check in on progress, redirect failing approaches, synthesize findings.

```text
Researcher, switch to checking network logs instead.
Implementer, hold off until researcher finishes hypothesis.
```

### 6. Enforce Quality Gates with Hooks

Exit code 2 to block + send feedback:

```bash
# ~/.claude/hooks/TeammateIdle
# Runs when teammate about to go idle
# Exit 2 to send feedback + keep working
```

Hook types:
- `TeammateIdle`: teammate finishing turn
- `TaskCreated`: task being created
- `TaskCompleted`: task being marked complete

### 7. Start with Research/Review

If new to teams: pick tasks with clear boundaries (code review, bug investigation, research). Shows value of parallel exploration without coordination complexity.

### 8. Wait for Teammates

Lead sometimes starts implementing instead of delegating:

```text
Wait for your teammates to complete their tasks before proceeding.
```

---

## Require Plan Approval

For complex/risky work, require plan before implementation:

```text
Spawn architect to refactor auth module.
Require plan approval before making changes.
Only approve plans that include comprehensive tests.
```

Workflow:
1. Teammate plans in read-only mode
2. Sends plan to lead for review
3. Lead approves or rejects with feedback
4. If rejected, teammate revises & resubmits
5. Once approved, exits plan mode → implementation

Lead makes approvals autonomously per your criteria.

---

## Parallel Workflows

### Parallel Code Review

Split review by domain:

```text
Spawn 3 reviewers for PR #142:
- Security-focused: token handling, auth, data leaks
- Performance-focused: algorithm complexity, DB queries, caching
- Test-focused: coverage gaps, edge cases, mock adequacy

Each review independently and report findings.
```

Each applies different filter. Lead synthesizes after all finish.

### Competing Hypotheses

App crashes after one message. Spawn teams to debate:

```text
Users report app exits after one message. Spawn 5 teammates to 
investigate different hypotheses:
- Memory leak (check heap)
- Connection timeout (check network logs)
- Event listener not cleaning up (check event handlers)
- Buffer overflow (check message parsing)
- Race condition in state management (check async code)

Talk to each other to try to disprove each other's theories.
Update findings doc with whatever consensus emerges.
```

Multiple independent investigators actively disproving each other → higher-quality root cause.

### Parallel Module Development

Each teammate owns separate module:

```text
Spawn 4 teammates for new payment system:
- Database: schema design, migrations
- API: endpoints, request/response handling
- Client: UI components, state management
- Tests: unit, integration, contract tests

Each own separate files. Coordinate via task list.
```

No file conflicts. Teammates self-coordinate via tasks.

---

## Troubleshooting

| Problem | Fix |
|---------|-----|
| **Teammates not appearing** | Check agent panel (in-process) or tmux panes (split). Are teammates hidden? Send one a message by name to wake it. |
| **Too many permission prompts** | Pre-approve common operations in permission settings before spawning. |
| **Teammate stops on error** | Message it directly with recovery instructions or spawn replacement. |
| **Lead does work instead of delegating** | Tell it: "Wait for teammates to finish before proceeding." |
| **Task appears stuck** | Check if work is actually done. Update task status or nudge teammate. |
| **Orphaned tmux session** | `tmux kill-session -t <session-name>` |
| **File overwrites (two teammates same file)** | Break work so each teammate owns different files. |

---

## Current Limitations (Experimental)

- ❌ **No session resumption with in-process teammates**: `/resume` doesn't restore teammates. Spawn new ones after resuming.
- ❌ **Task status lag**: teammates sometimes fail to mark tasks complete. May block dependents. Update manually if stuck.
- ❌ **Shutdown slow**: teammates finish current request/tool before exiting.
- ❌ **One team per session**: no multiple named teams or cross-session sharing.
- ❌ **No nested teams**: teammates can't spawn their own teammates.
- ❌ **No background subagents from in-process teammates**: background work can't outlive lead process.
- ❌ **Lead is fixed**: can't promote teammate or transfer leadership.
- ❌ **Permissions at spawn**: all teammates start with lead's mode. Can change after, not at spawn.
- ❌ **Split panes require tmux/iTerm2**: in-process works anywhere.

---

## Token Cost

Agent teams use **significantly more tokens** than single session. Each teammate = full separate context window.

**Worth the cost for:**
- Research where parallel exploration saves time
- New features with independent modules
- Bug investigation with competing theories
- Complex reviews with multiple perspectives

**Not worth for:**
- Simple tasks
- Sequential workflows
- Single-focused work

See [Claude docs: agent team token costs](https://code.claude.com/docs/en/costs#agent-team-token-costs) for detailed pricing.

---

## Integrating with DAWAI

### Use Case: Multi-Path Backend Testing

```text
Spawn 3 teammates to test DAWAI's multi-tenant isolation:
- Tenant isolation tester: verify school_id filtering across all queries
- Cross-tenant access tester: attempt unauthorized access patterns
- Performance tester: load test with concurrent school operations

Test against live database. Report vulnerabilities found.
```

### Use Case: Parallel Feature Development

```text
DAWAI Phase 13 security hardening. Spawn 4 teammates:
- Query auditor: grep all handlers for missing school_id filters
- JWT verifier: verify token isolation across endpoints
- Test writer: add cross-tenant access tests
- Documentation: list all security patterns used

Each own separate files in docs/ + tests/. Coordinate via task list.
```

### Use Case: Code Review + Refactor

```text
Spawn reviewers for auth module refactor:
- Security specialist: token handling, blacklist, expiry
- Performance specialist: JWT parsing overhead, caching
- Maintainability specialist: error handling, logging, tests

Each review independently. If approved, spawn implementer to refactor.
```

---

## Settings Configuration

### Enable in settings.json

```json
{
  "env": {
    "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1"
  },
  "teammateMode": "auto",
  "defaultTeammateModel": "claude-haiku-4-5"
}
```

### CLI Flag

```bash
claude --teammate-mode auto
```

---

## Related Tools

- **[Subagents](/docs/en/sub-agents)**: lightweight delegation within session
- **[Git Worktrees](/docs/en/worktrees)**: run multiple sessions manually
- **[Hooks](/docs/en/hooks)**: enforce quality gates
- **[Skills](/docs/en/skills)**: define reusable patterns

---

## References

- Claude Code Docs: https://code.claude.com/docs/en/agent-teams
- Costs & Token Usage: https://code.claude.com/docs/en/costs#agent-team-token-costs
- Subagents Reference: https://code.claude.com/docs/en/sub-agents
