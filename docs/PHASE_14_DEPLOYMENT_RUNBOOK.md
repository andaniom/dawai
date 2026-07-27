# Phase 14: Deployment Runbook

## 1. Pre-Deployment Checklist
- [ ] Code reviewed and merged to `main`
- [ ] Environment variables configured (`.env.production` or CI secrets)
- [ ] Database credentials rotated for production
- [ ] Audit logs and security controls verified
- [ ] No hardcoded `school_id` values in queries

## 2. Lighthouse PWA Audit (Target >= 90)
- Run audit: `npm run build && npm start` (in `frontend/`)
- Open Chrome DevTools > Lighthouse
- Target scores:
  - Performance: >= 90
  - Accessibility: >= 90
  - Best Practices: >= 90
  - SEO: >= 90
  - PWA: Installable and offline-capable

## 3. Docker Deployment
```bash
# Pull latest changes
git pull origin main

# Build images
docker compose -f docker-compose.prod.yml build

# Start services in detached mode
docker compose -f docker-compose.prod.yml up -d

# Check database migrations applied automatically
docker compose -f docker-compose.prod.yml logs api

# Smoke test
curl -I https://your-domain.com/api/health
```

## 4. E2E Testing
- [ ] Verify Super Admin login and school creation
- [ ] Verify School Admin login and user management
- [ ] Verify Teacher login and assessment submission
- [ ] Verify Student/Parent login and portal access
- [ ] Verify cross-tenant isolation (Student A cannot see Student B's data from another school)

## 5. Monitoring + Alerts
- Monitor Docker containers: `docker stats`
- View application logs: `docker compose logs -f`
- Monitor Woodpecker CI pipeline status
- (Optional) Setup Telegram/Slack alerts for 5xx errors

## 6. Rollback Procedure
```bash
# If deployment fails, revert to previous image tag or commit
git checkout <previous_stable_commit>

# Rebuild and restart
docker compose -f docker-compose.prod.yml up -d --build

# If database needs rollback (WARNING: Data loss potential)
# migrate -path migrations -database "$DATABASE_URL" down 1
```

## 7. Phase 14 Success Criteria
- [ ] Application running in production environment via Docker
- [ ] PWA installable and passes Lighthouse audit (>= 90)
- [ ] CI/CD pipeline (Woodpecker) successfully deploying changes
- [ ] Multi-tenant isolation confirmed in production
- [ ] HTTPS/SSL enabled and working
