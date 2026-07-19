#!/bin/bash
set -euo pipefail

IMAGE="ghcr.io/runfridge/wedding-card"
SERVICE="wedding-card"
DB_DIR="./data/sqlite"
BACKUP_DIR="./data/backups"
HEALTH_TIMEOUT=60

usage() {
  echo "usage: $(basename "$0") <version>   e.g. 1.42.0" >&2
  exit 1
}

[[ $# -eq 1 ]] || usage
NEW_VERSION="$1"

for cmd in docker; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "Error: $cmd is not installed" >&2
    exit 1
  fi
done

if ! docker info &>/dev/null; then
  echo "Error: Docker daemon is not running" >&2
  exit 1
fi

if [[ ! -f docker-compose.yml ]]; then
  echo "Error: docker-compose.yml not found in $(pwd)" >&2
  echo "       run this script from the directory that holds docker-compose.yml" >&2
  exit 1
fi

if [[ ! -f .env ]]; then
  echo "Error: .env not found in $(pwd) (must define WEDDING_VERSION=...)" >&2
  exit 1
fi

OLD_VERSION=$(grep -E '^WEDDING_VERSION=' .env | cut -d= -f2- || true)
if [[ -z "$OLD_VERSION" ]]; then
  echo "Error: WEDDING_VERSION is not set in .env" >&2
  exit 1
fi

if [[ "$NEW_VERSION" == "$OLD_VERSION" ]]; then
  echo "Already at $NEW_VERSION, nothing to do."
  exit 0
fi

echo "Current version: $OLD_VERSION"
echo "Target version:  $NEW_VERSION"

echo "Pulling $IMAGE:$NEW_VERSION"
docker pull "$IMAGE:$NEW_VERSION"

STAMP=$(date +%Y%m%d-%H%M%S)
if [[ -f "$DB_DIR/wedding.db" ]]; then
  mkdir -p "$BACKUP_DIR"
  BACKUP_BASE="$BACKUP_DIR/wedding-${OLD_VERSION}-${STAMP}.db"
  echo "Backing up database → $BACKUP_BASE"
  cp -a "$DB_DIR/wedding.db" "$BACKUP_BASE"
  [[ -f "$DB_DIR/wedding.db-wal" ]] && cp -a "$DB_DIR/wedding.db-wal" "${BACKUP_BASE}-wal"
  [[ -f "$DB_DIR/wedding.db-shm" ]] && cp -a "$DB_DIR/wedding.db-shm" "${BACKUP_BASE}-shm"
fi

echo "$OLD_VERSION" > .previous_version

sed -i "s/^WEDDING_VERSION=.*/WEDDING_VERSION=${NEW_VERSION}/" .env

echo "Swapping $SERVICE container"
docker compose up -d "$SERVICE"

echo "Waiting up to ${HEALTH_TIMEOUT}s for healthcheck..."
CID=$(docker compose ps -q "$SERVICE")
deadline=$((SECONDS + HEALTH_TIMEOUT))
while (( SECONDS < deadline )); do
  STATUS=$(docker inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}' "$CID" 2>/dev/null || echo "unknown")
  case "$STATUS" in
    healthy)
      echo "Done: deployed $NEW_VERSION (was $OLD_VERSION)"
      exit 0
      ;;
    unhealthy)
      echo "Error: container reported unhealthy" >&2
      echo "       run scripts/rollback.sh to revert to $OLD_VERSION" >&2
      exit 1
      ;;
  esac
  sleep 2
done

echo "Warning: healthcheck did not report healthy within ${HEALTH_TIMEOUT}s" >&2
echo "         run scripts/rollback.sh to revert to $OLD_VERSION" >&2
exit 1
