#!/bin/bash
# scripts/migrate.sh - 数据库迁移脚本

set -e

DIRECTION=${1:-up}
STEPS=${2:-}

# 从 .env 读取 DATABASE_URL
if [ -f ".env" ]; then
  export $(grep -v '^#' .env | xargs)
fi

DB_URL=${DATABASE_URL:-"postgres://blog:blogpass@localhost:5432/blog?sslmode=disable"}

MIGRATE_BIN=$(which migrate 2>/dev/null || echo "")

if [ -z "$MIGRATE_BIN" ]; then
  echo "Error: golang-migrate not found. Install with:"
  echo "  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
  exit 1
fi

MIGRATIONS_PATH="./migrations"

echo "Running migration: ${DIRECTION} ${STEPS}"

if [ -n "$STEPS" ]; then
  migrate -path "${MIGRATIONS_PATH}" -database "${DB_URL}" "${DIRECTION}" "${STEPS}"
else
  migrate -path "${MIGRATIONS_PATH}" -database "${DB_URL}" "${DIRECTION}"
fi

echo "Migration completed."
