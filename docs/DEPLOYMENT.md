# DAWAI Deployment Guide

Stack: Docker Compose + Caddy reverse proxy (auto-HTTPS) + PostgreSQL 16 + MinIO.

**Target**: VPS with Ubuntu 22.04+, 1 vCPU, 1 GB RAM (t2.small or equivalent).

---

## 1. Prerequisites

| Item | Version | Check |
|------|---------|-------|
| Docker | 24+ | `docker --version` |
| Docker Compose | v2+ | `docker compose version` |
| Git | 2.30+ | `git --version` |
| Domain | — | DNS A record pointing to VPS IP |
| Ports | — | 80, 443 open in firewall |

```bash
# Ubuntu quick install
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
# Log out / log back in for group to take effect
```

---

## 2. Server Setup

### 2.1 Clone and Configure

```bash
sudo apt update && sudo apt install -y git
git clone https://github.com/yourorg/dawai.git /opt/dawai
cd /opt/dawai
cp .env.example .env
```

### 2.2 Fill `.env`

```bash
# Generate secrets
openssl rand -base64 32  # → JWT_SECRET
openssl rand -base64 32  # → AUTH_SECRET

# Edit .env with real values
# Key changes from defaults:
#   POSTGRES_PASSWORD  — strong random password (≥24 chars)
#   JWT_SECRET         — random 32+ char string
#   AUTH_SECRET        — random 32+ char string
#   NEXTAUTH_URL       — https://dawai.yourschool.edu
#   NEXT_PUBLIC_API_URL — https://dawai.yourschool.edu/api
#   API_URL_INTERNAL   — http://api:8080  (docker network, keep as-is)
```

**Security rule**: Never commit `.env` to git. The `.gitignore` already excludes it.

---

## 3. Production Deploy

```bash
cd /opt/dawai
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
```

Verify:
```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml ps
# All services should show "Up" status

# Check logs
docker compose logs -f api    # Go API
docker compose logs -f frontend  # Next.js
docker compose logs -f postgres  # Database
```

---

## 4. Reverse Proxy + SSL (Caddy)

Caddy auto-provisions Let's Encrypt certificates.

### 4.1 Install Caddy

```bash
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install caddy
```

### 4.2 Configure Caddy

Create `/etc/caddy/Caddyfile`:
```caddyfile
# Replace dawai.yourschool.edu with your actual domain
dawai.yourschool.edu {
    # Frontend (Next.js)
    handle /_next/static/* {
        reverse_proxy frontend:3000
    }

    # API routes → Go backend
    handle /api/* {
        reverse_proxy api:8080
    }

    # Everything else → Next.js frontend
    handle {
        reverse_proxy frontend:3000
    }

    # Security headers
    header {
        X-Content-Type-Options nosniff
        X-Frame-Options DENY
        Referrer-Policy strict-origin-when-cross-origin
        Strict-Transport-Security "max-age=31536000; includeSubDomains"
    }
}
```

### 4.3 Start Caddy

```bash
sudo systemctl enable --now caddy
sudo systemctl status caddy
```

Caddy will auto-provision SSL via Let's Encrypt on first request.

---

## 5. Database Management

### 5.1 Connect to PostgreSQL

```bash
# Via docker
docker compose -f docker-compose.yml -f docker-compose.prod.yml \
  exec postgres psql -U dawai -d dawai
```

### 5.2 Run Migrations

Migrations run automatically on `docker compose up`. To run manually:
```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml \
  exec api migrate -path /app/migrations -database "$DATABASE_URL" up
```

### 5.3 Backup (manual)

```bash
# Dump database
docker compose -f docker-compose.yml -f docker-compose.prod.yml \
  exec postgres pg_dump -U dawai -d dawai -F c > backup_$(date +%Y%m%d).dump

# Restore
cat backup_20260723.dump | docker compose -f docker-compose.yml -f docker-compose.prod.yml \
  exec -T postgres pg_restore -U dawai -d dawai
```

---

## 6. Secrets Management

### 6.1 Woodpecker CI (Recommended)

Set secrets in Woodpecker project settings (Settings → Secrets):
- `db_password` → same as `POSTGRES_PASSWORD`
- `jwt_secret` → same as `JWT_SECRET`
- `auth_secret` → same as `AUTH_SECRET`
- `smtp_password` → same as `SMTP_PASS`
- `google_client_id` → `GOOGLE_CLIENT_ID`
- `google_client_secret` → `GOOGLE_CLIENT_SECRET`

CI injects these at build time — they are never in code or `.env` files in git.

### 6.2 Manual VPS (No CI)

If deploying manually without Woodpecker:
```bash
# Generate strong secrets
openssl rand -base64 32 | tr -d '\n'  # Use output for JWT_SECRET
openssl rand -base64 32 | tr -d '\n'  # Use output for AUTH_SECRET
openssl rand -base64 24 | tr -d '\n'  # Use output for POSTGRES_PASSWORD

# Add to .env
vi /opt/dawai/.env
```

### 6.3 Docker Secrets (Alternative)

For more secure secret handling, use Docker secrets:
```bash
echo "your_secret" | docker secret create jwt_secret -
echo "your_secret" | docker secret create postgres_password -

# Then reference in docker-compose.prod.yml:
# environment:
#   JWT_SECRET_FILE: /run/secrets/jwt_secret
# secrets:
#   - jwt_secret
```

---

## 7. Monitoring & Logs

### 7.1 View Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f api
docker compose logs -f postgres

# With timestamps
docker compose logs -f --timestamps api
```

### 7.2 Health Checks

The production compose includes health checks:
```bash
docker compose ps
# Shows "healthy" status for postgres, api, frontend
```

### 7.3 Disk Usage

```bash
# Docker disk usage
docker system df

# Clean up unused images/containers
docker system prune -f

# WARNING: Only run on maintenance windows
```

---

## 8. Scaling & Resource Limits

### 8.1 Current Limits (in docker-compose.prod.yml)

| Service | Memory Limit |
|---------|--------------|
| postgres | 512 MB |
| api | 256 MB |
| frontend | 256 MB |

### 8.2 Adjusting for Load

For >50 concurrent users:
```bash
# Edit docker-compose.prod.yml
# postgres: increase to 1024M
# api: increase to 512M
# frontend: increase to 512M

# Restart
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

For >100 concurrent users, consider:
- Separate PostgreSQL server (Supabase, Railway, Render)
- Multiple API replicas (requires load balancer)
- CDN for static assets (Cloudflare, AWS CloudFront)

---

## 9. Troubleshooting

### 9.1 Service Won't Start

```bash
# Check logs
docker compose logs api
docker compose logs postgres

# Check if port is in use
sudo lsof -i :80
sudo lsof -i :443

# Rebuild without cache
docker compose build --no-cache
docker compose up -d
```

### 9.2 Database Connection Refused

```bash
# Check postgres is running
docker compose ps postgres

# Check postgres health
docker compose exec postgres pg_isready -U dawai -d dawai

# Check DATABASE_URL matches .env
docker compose exec api printenv DATABASE_URL
```

### 9.3 SSL Certificate Issues

```bash
# Check Caddy logs
sudo journalctl -u caddy -f

# Force certificate renewal
sudo systemctl restart caddy

# Check domain DNS
dig dawai.yourschool.edu
```

### 9.4 Memory Issues (OOM Killer)

```bash
# Check if OOM killed
sudo dmesg | grep -i oom

# Check container memory usage
docker stats

# Solution: increase memory limits in docker-compose.prod.yml
```

---

## 10. Updates & Maintenance

### 10.1 Update DAWAI

```bash
cd /opt/dawai
git pull origin main
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
```

### 10.2 Database Migrations

Migrations run automatically. To verify:
```bash
docker compose exec api migrate -path /app/migrations -database "$DATABASE_URL" version
```

### 10.3 Rotate Secrets

```bash
# Generate new secrets
NEW_JWT_SECRET=$(openssl rand -base64 32 | tr -d '\n')
NEW_AUTH_SECRET=$(openssl rand -base64 32 | tr -d '\n')

# Update .env
vi /opt/dawai/.env

# Restart (existing JWTs will be invalidated)
docker compose -f docker-compose.yml -f docker-compose.prod.yml restart api frontend
```

---

## 11. Rollback

```bash
cd /opt/dawai
git log --oneline -5  # Find previous commit
git checkout <commit-hash>

# Restore database from backup if needed
cat backup_20260723.dump | docker compose -f docker-compose.yml -f docker-compose.prod.yml \
  exec -T postgres pg_restore -U dawai -d dawai

docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
```

---

## Quick Reference

| Command | Purpose |
|---------|---------|
| `docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d` | Start production |
| `docker compose logs -f api` | View API logs |
| `docker compose exec postgres psql -U dawai -d dawai` | Connect to DB |
| `docker compose ps` | Service status |
| `docker system prune -f` | Clean up disk |
| `sudo systemctl restart caddy` | Restart reverse proxy |

---

**Support**: Check `README.md` for development setup, `CLAUDE.md` for architecture details.
