#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 1 ]; then
  echo "Usage: $0 <backup-file.dump> [--drop]"
  echo "  --drop  Drop and recreate database before restore (DESTRUCTIVE)"
  exit 1
fi

BACKUP_FILE="$1"
DB_USER="${DB_USER:-dawai}"
DB_NAME="${DB_NAME:-dawai}"
CONTAINER="${CONTAINER:-dawai-postgres-1}"

if [ ! -f "$BACKUP_FILE" ]; then
  echo "Error: File not found: $BACKUP_FILE"
  exit 1
fi

echo "WARNING: This will overwrite the '$DB_NAME' database."
read -p "Continue? [y/N] " -n 1 -r
echo
[[ $REPLY =~ ^[Yy]$ ]] || exit 0

if [ "${2:-}" = "--drop" ]; then
  echo "Dropping and recreating database..."
  docker exec "$CONTAINER" psql -U "$DB_USER" -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;"
  docker exec "$CONTAINER" psql -U "$DB_USER" -d postgres -c "CREATE DATABASE $DB_NAME OWNER $DB_USER;"
fi

echo "Restoring from $BACKUP_FILE..."
docker exec -i "$CONTAINER" pg_restore -U "$DB_USER" -d "$DB_NAME" --no-owner --no-acl < "$BACKUP_FILE"

echo "Restore complete."