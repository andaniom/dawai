# Advanced Agent Teams Patterns & Troubleshooting

Advanced techniques, gotchas, and recovery strategies for complex team scenarios.

---

## Architecture Deep Dive

### File Locations & Persistence

```
~/.claude/teams/
├── session-abc12345/          # Team workspace (deleted on session end)
│   ├── config.json            # Runtime state (DO NOT edit)
│   └── inboxes/
│       ├── lead.json          # Messages TO lead
│       ├── teammate1.json     # Messages TO teammate1
│       └── teammate2.json
│
~/.claude/tasks/
├── session-abc12345/          # Task list (PERSISTS across sessions)
│   ├── task-1.json
│   ├── task-2.json
│   └── task-3.json
```

**Key insight:** Task list survives session resume. Team config is rebuilt. If you resume session, teammates are gone but task list remains—must spawn new teammates.

### Mailbox Format

Each message in `~/.claude/teams/.../inboxes/teammate.json` is JSON:

```json
{
  "from": "lead",
  "timestamp": "2026-07-24T10:30:00Z",
  "body": "Tell researcher to focus on database logs",
  "message_id": "msg-xyz"
}
```

Invalid messages are reported + removed on next read. Malformed JSON blocks delivery until fixed or file deleted.

---

## Advanced Patterns

### Phased Team Handoff

Task dependencies enforce sequential work across phases.

```text
Phase 1: Database design (database-expert)
Phase 2: API implementation (api-expert) — depends on Phase 1
Phase 3: Frontend (frontend-expert) — depends on Phase 2
Phase 4: Tests (test-expert) — depends on Phase 3

Create 4 tasks with dependencies. Spawn 1 teammate per task.
Lead: create task 1. Teammate claims, completes, task 2 unblocks.
Teammate 2 spawns automatically or lead manually assigns.
```

Benefit: Sequential phases still use teams (teammates work on parallel aspects within phase).

### Dynamic Team Resizing

Start small, scale up if needed.

```text
Start with 2-person team. If work expands:

"Spawn a performance testing teammate to join the current team
and focus on load testing scenarios."

Lead: update task list, new teammate claims available work.
No need to rebuild team—just add to existing one.
```

### Conversation Stitching

When resuming session, catch new teammates up without full history.

```text
Resumed Phase 13 security sprint.

Spawn replacement researchers: [previous session found 3 SQL 
injection vectors at endpoints /assessments, /users, /reports]

Investigate: are there others? Test all endpoints for similar patterns.

Report: total SQL injection vectors found + fix recommendations.
```

Provide relevant prior findings in spawn prompt. New teammates don't see old conversation.

### Task Delegation Chains

Teammates delegate to each other via messaging.

```text
Lead creates central task list. Teammates message each other:

[Security reviewer to performance reviewer]
"I found 5 SQL queries without indexes. Can you measure their impact?"

[Performance reviewer responds]
"Query 1-3 fine. Query 4-5 are bottlenecks at >100ms. 
See details in shared doc."

[Back to lead]
"Security team + performance team converged on same 2 queries
as priority fixes."
```

Benefit: teammates collaborate without lead intermediating every message.

---

## Handling Common Issues

### Issue: Teammate Stopped Working on Error

Teammate encountered API error, gave up instead of retrying.

**Solution:**

```text
Tell backend-expert: "You hit an API timeout. The service is back up now.
Can you retry the query audit from where you stopped?"

Or if it's stuck:
"Spawn a replacement researcher to continue backend-expert's work.
Previous teammate found 12 queries missing school_id filters.
Find remaining queries."
```

As of v2.1.198, teammates receiving messages wake up and retry failed API calls.

### Issue: Task Marked Complete But Work Incomplete

Teammate claimed task done but didn't actually finish.

**Solution:**

```text
Check task status: "Mark task #3 as pending (it's not actually done).
Tell the researcher to complete it or update the task description
so next teammate knows what's left."

Or manually update:
In ~/.claude/tasks/session-abc/task-3.json, change status from "completed" to "in_progress"
Teammate notices, resumes work.
```

Alternatively, tell lead: `Update task #3 status to pending and reassign to another teammate.`

### Issue: Two Teammates Edited Same File

Merge conflict. Git shows file changed by both.

**Prevention (next time):**
- Give clear file ownership: "Backend-expert: owns backend/internal/handlers/auth*.go"
- Frontend-expert owns separate files
- Test-expert owns test files only

**Recovery (now):**
```bash
# Review conflict
git diff backend/internal/handlers/conflicted.go

# Keep backend-expert's version (assume it's more recent)
git checkout --theirs backend/internal/handlers/conflicted.go

# Or manual merge: take parts from each
# Then test
go test ./...

# Commit resolution
git add backend/internal/handlers/conflicted.go
git commit -m "Resolve merge conflict: kept auth handler from backend-expert, 
tests from test-expert"
```

Ask lead: "Merge the files and test. Then tell each teammate about the conflict
so they coordinate better next time."

### Issue: Lead Started Working Instead of Delegating

Lead implemented feature instead of waiting for teammate to finish.

**Solution:**

```text
Tell lead: "You started implementing the feature, but the frontend 
team is already working on it. Let's wait for them to finish before 
you add more. What should the frontend team do next?"

Or if work diverged:
"The frontend team implemented ReportViewer as a React component.
You implemented it as a plain HTML page. Decide which approach to use,
then merge or discard one."
```

### Issue: In-Process Teammate Hidden After Going Idle

Teammate not visible in agent panel, appears gone.

**Solution:**

Teammate is hidden (idle), not stopped. Still running.

- Send it a message: `Tell researcher to wake up and continue investigating`
- Or select from "N idle agents" collapsed row (press Enter to expand)

Idle teammates reappear when given work or after next turn starts.

### Issue: Orphaned Tmux Session

Split-pane team was closed but tmux session persists.

**Solution:**

```bash
tmux ls
# Shows: session-abc12345: 4 windows...

tmux kill-session -t session-abc12345

# Verify it's gone
tmux ls
# Shows: no servers running
```

### Issue: Mailbox Corruption (Malformed JSON)

Messages stopped delivering to a teammate. Mailbox has bad JSON.

**Solution:**

```bash
# Check mailbox
cat ~/.claude/teams/session-abc/inboxes/teammate.json

# If malformed:
# Option 1: Delete the file (teammate misses those messages)
rm ~/.claude/teams/session-abc/inboxes/teammate.json

# Option 2: Fix the JSON (if salvageable)
# Remove trailing comma, quote properly, etc

# Option 3: Message teammate directly to resend
"Tell researcher: the last message you received was 'investigate auth'.
Can you confirm you got it or should I resend?"
```

---

## Performance Tuning

### Reducing Token Cost

Agent teams are expensive. Minimize without losing parallelism.

**Technique 1: Narrow focus**
```text
DON'T: "Review the entire codebase for security issues"
DO:    "Review auth module (src/auth/) for JWT handling vulnerabilities"
```

Narrower scope = shorter context = fewer tokens.

**Technique 2: Async results**
```text
Spawn researcher, but don't wait for full report.
Lead: continue own work.
When researcher finishes, incorporate findings.
```

Less context switching = fewer token exchanges.

**Technique 3: Batch messages**
```text
DON'T: Send 5 separate messages to researcher
DO:    Send 1 message with all 5 questions
```

Fewer round-trips = fewer tokens.

**Technique 4: Limit team size**
```text
Optimal: 3-5 teammates
Avoid: 10+ teammates
```

Each extra teammate = full context window. Diminishing returns after 5.

### Monitoring Token Usage

No built-in meter, but estimate:
- Each teammate = ~1 full Claude context window per turn
- 5 teammates = 5x the cost of single session
- Iterative refine (multiple turns) multiplies cost

**Decision rule:** Use teams only if the time saved (parallelism) justifies the cost.

---

## Session Resumption Strategy

Agent teams don't persist across resume. Recovery:

### Before Closing Session

```text
Summarize findings document:

Create file: /tmp/team-findings.md with:
1. What each teammate completed
2. What's still outstanding
3. Current task status
4. Key decisions made
5. Next steps

This gives new session context when resuming.
```

### After Resuming Session

```text
[Manual] Read /tmp/team-findings.md (your prior session's output)

[In new session]
"My previous team was working on Phase 13 security hardening.
Here's what they found: [copy findings]. 

I need to:
1. Merge their findings into compliance docs
2. Run any remaining tests
3. Fix identified issues

Spawn a team to pick up where we left off: verify all findings,
add any missing tests, update Phase 13 docs."
```

New team resumes work with full context.

---

## Quality Enforcement with Hooks

Use hooks to prevent bad work from reaching merge.

### TeammateIdle Hook

Runs when teammate is about to go idle. Exit code 2 to block:

```bash
#!/bin/bash
# ~/.claude/hooks/TeammateIdle

# Check: Did security reviewer find any SQL injection?
if grep -q "SQL injection" /tmp/team-findings.md; then
  echo "Security team found SQL injection vulnerabilities that must be fixed before approval"
  exit 2  # Block teammate from going idle, keep them working
fi

exit 0  # OK to go idle
```

### TaskCompleted Hook

Runs when teammate marks task complete. Exit code 2 to block:

```bash
#!/bin/bash
# ~/.claude/hooks/TaskCompleted

# Reject: Code review task cannot be marked complete without findings doc
if [[ "$TASK_NAME" =~ "code review" ]]; then
  if ! test -f docs/PHASE_13_CODE_REVIEW.md; then
    echo "Code review task must produce findings document: docs/PHASE_13_CODE_REVIEW.md"
    exit 2  # Block completion, send feedback
  fi
fi

exit 0
```

### TaskCreated Hook

Runs when task is being created. Exit code 2 to block:

```bash
#!/bin/bash
# ~/.claude/hooks/TaskCreated

# Reject: task with no clear deliverable
if [[ "$TASK_DESC" != *"report"* ]] && [[ "$TASK_DESC" != *"test"* ]]; then
  echo "Tasks must have clear deliverable: 'produce report' or 'write test' or 'implement X'"
  exit 2  # Block creation, send feedback
fi

exit 0
```

---

## Communication Patterns

### Status Synchronization

Regular check-ins prevent drift:

```text
[After 30 min of work]

Lead: "Quick status update from all teammates:
- Database expert: which queries have you verified so far?
- Security expert: have you found any issues?
- Test expert: are tests passing?

Report in 1-2 sentences each."

Benefit: catch blockers early, redirect if needed.
```

### Escalation Pattern

When teammate encounters unknown:

```text
[Test expert to lead]
"I tried to test POST /api/assessments but got 401 Unauthorized.
Is the endpoint deployed? Should I mock it instead?"

[Lead to API expert]
"Test expert needs POST /api/assessments deployed. Can you confirm
it's running and provide test credentials?"

[API expert responds]
"Endpoint is live. Use these test creds. I'll add test user to DB."

[Lead to test expert]
"Endpoint is ready. API expert added test user. Retry."
```

Clear escalation unblocks quickly.

### Debate Resolution

When teammates disagree:

```text
[Security expert vs performance expert]
Security: "Add bcrypt cost 14 for maximum security"
Performance: "Cost 12 is fine, reduces login latency"

[Lead intervenes]
"Security expert: what's the real-world attack cost of cost 12 vs 14?
Performance expert: how much latency difference are we talking?

Once you have numbers, I'll decide based on risk/performance tradeoff."
```

Numbers beat opinions.

---

## Anti-Patterns to Avoid

❌ **Spawning team for 1-hour task**: overhead > benefit. Use subagent instead.

❌ **No clear deliverables**: "Review the code" vs "Produce 1-page findings doc listing all SQL injection risks."

❌ **Team working unattended for days**: leads to wasted effort, duplicated work. Check in daily.

❌ **10+ teammates**: coordination overhead explodes. Scale back to 3-5.

❌ **Teammates given no context**: "Review auth" vs "Review auth flow: OAuth redirect → NextAuth → Go /api/auth/token endpoint → JWT returned. Focus on token validation."

❌ **Mixing sequential + parallel work**: some tasks must happen in order. Use task dependencies, don't hide them.

❌ **No isolation verification on multi-tenant code**: teammates ship feature, later found data leak. Make isolation test mandatory in spawn prompt.

---

## When to Use Teams vs Alternatives

| Task | Teams | Subagents | Single Session |
|------|-------|-----------|----------------|
| Code review (3 reviewers) | ✅ best | ❌ no inter-review | ❌ serial |
| New feature (4 modules) | ✅ best | ⚠️ OK if independent | ❌ serial |
| Bug investigation (competing theories) | ✅ best | ⚠️ OK | ❌ serial |
| Quick task (1 helper) | ❌ overhead | ✅ best | ⚠️ OK |
| Long research (1 path) | ❌ overhead | ✅ best | ⚠️ OK |
| Urgent fix (1 person) | ❌ overhead | ❌ overhead | ✅ best |

**Decision tree:**
1. Can work be parallelized? If no → single session
2. Do teammates need to communicate? If no → subagents
3. Is the task complex (3+ days)? If yes → teams
4. Is coordination overhead acceptable? If yes → teams

If all yes → **use teams**.

---

## Checklists

### Pre-Spawn Checklist

- [ ] Task is parallelizable? (inherently parallel, not forced)
- [ ] Team size 3-5? (more → overhead explodes)
- [ ] Each teammate owns different files? (no conflicts)
- [ ] Clear deliverables? (each teammate knows when done)
- [ ] Context complete? (task-specific details in spawn prompt)
- [ ] Dependencies documented? (if task A blocks task B, explain)
- [ ] Success criteria defined? (what does "done" look like?)

### Mid-Sprint Checklist (Daily)

- [ ] All teammates still productive? (or are they stuck?)
- [ ] Blockers visible? (anyone waiting on someone else?)
- [ ] Findings being shared? (teammates talking to each other?)
- [ ] Scope creep? (task expanding beyond original?)
- [ ] Quality on track? (or will need rework?)

### End-Sprint Checklist

- [ ] All deliverables produced?
- [ ] Teammates can shut down gracefully?
- [ ] Findings synthesized into actionable items?
- [ ] Code review passed (if applicable)?
- [ ] Tests pass?
- [ ] Documentation updated?
- [ ] Lessons learned recorded (for next team)?

---

## Debugging Team Deadlock

When team stops making progress:

1. **Check task status**: are tasks marked done but actually incomplete?
   ```bash
   cat ~/.claude/tasks/session-abc/task-1.json | jq '.status'
   ```

2. **Check mailboxes**: are messages stalled?
   ```bash
   cat ~/.claude/teams/session-abc/inboxes/teammate.json | tail -5
   ```

3. **Message lead**: "What's blocking the team right now?"

4. **Send recovery instructions**:
   ```text
   Database expert: mark task #2 as pending, you're not done
   API expert: I'm waiting on schema from database expert
   Lead: create task #5 (parallel path) so API expert has work while waiting
   ```

5. **If still stuck**: spawn fresh teammate with clear instructions, let stuck one finish gracefully

---

## Cost Estimation

Rough token costs (Opus 4.8):

- 1-hour task, 1 teammate: ~40K tokens
- 1-day task, 3 teammates: ~120K tokens (40K × 3)
- 3-day sprint, 4 teammates, iterative: ~500K+ tokens

**ROI calculation:**
- If parallelism saves 2 days of serial work → 120K tokens saves 2 days
- If token cost is $1.50/100K → $1.80 for 2 days saved = worthwhile

But if no parallelism (sequential dependencies) → tokens wasted.
