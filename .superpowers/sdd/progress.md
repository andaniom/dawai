# DAWAI Rewrite Progress Ledger

## Phase 1: Backend Scaffolding (2 days)
- [x] Task 1.1: Initialize Go backend project structure
- [x] Task 1.2: Set up PostgreSQL migrations (golang-migrate)
- [ ] Task 1.3: Set up sqlc for type-safe queries
- [ ] Task 1.4: Set up middleware and auth service skeleton

## Phase 2: Frontend Setup (1 day)
- [x] Task 2.1: Initialize Next.js project with shadcn/ui
- [ ] Task 2.2: Create auth store and API client

## Phase 3: Backend Features (3 days)
- [ ] Task 3.1: Implement complete auth endpoints with JWT
- [ ] Task 3.2: Implement subjects & rubric endpoints
- [ ] Task 3.3: Implement assessments endpoints
- [ ] Task 3.4: Implement students & users endpoints

## Phase 4: Frontend Features (2 days)
- [ ] Task 4.1: Implement login page
- [ ] Task 4.2: Implement dashboard & school switcher
- [ ] Task 4.3: Implement subject management (admin)
- [ ] Task 4.4: Implement assessment entry form
- [ ] Task 4.5: Implement student list with inline assessment
- [ ] Task 4.6: Implement student results view

## Phase 5: Integration & Testing (2 days)
- [ ] Task 5.1: Connect frontend to backend auth flow
- [ ] Task 5.2: Add multi-tenant isolation tests
- [ ] Task 5.3: Add cross-tenant FK validation tests
- [ ] Task 5.4: Add E2E tests (login → assess → view results)

## Phase 6: Deployment Prep (2 days)
- [ ] Task 6.1: Configure Woodpecker CI/CD
- [ ] Task 6.2: Docker Compose production config
- [ ] Task 6.3: Environment setup (staging/production)
- [ ] Task 6.4: Database backup strategy

## Completed Tasks
(none yet)
---
Task 1.1: complete (commits df3f9fd..0f279c4, review clean)
Task 1.2: complete (commit 47895e6, review clean)
Task 2.1: complete (commit 2893f7d, Next.js 16, review clean)
