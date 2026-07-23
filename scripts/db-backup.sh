#!/usr/bin/env bash
set -euo pipefail

BACKUP_DIR="${BACKUP_DIR:-/var/backups/dawai}"
DB_USER="${DB_USER:-dawai}"
DB_NAME="${DB_NAME:-dawai}"
CONTAINER="${CONTAINER:-dawai-postgres-1}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
KEEP_DAYS="${KEEP_DAYS:-7}"

mkdir -p "$BACKUP_DIR"

echo "[$(date)] Backing up $DB_NAME from $CONTAINER..."
docker exec "$CONTAINER" pg_dump -U "$DB_USER" -d "$DB_NAME" --format=custom \
  > "$BACKUP_DIR/${DB_NAME}_${TIMESTAMP}.dump"

echo "[$(date)] Cleaning backups older than $KEEP_DAYS days..."
find "$BACKUP_DIR" -name "${DB_NAME}_*.dump" -mtime +"$KEEP_DAYS" -delete

echo "[$(date)] Backup complete: ${DB_NAME}_${TIMESTAMP}.dump"
ls -lh "$BACKUP_DIR/${DB_NAME}_${TIMESTAMP}.dump"
