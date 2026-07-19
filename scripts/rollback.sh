#!/bin/bash
set -euo pipefail

IMAGE="ghcr.io/runfridge/wedding-card"
SERVICE="wedding-card"

usage() {
  echo "usage: $(basename "$0") [version]" >&2
  echo "       with no argument, reads .previous_version written by deploy.sh" >&2
  exit 1
}

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
  exit 1
fi

if [[ ! -f .env ]]; then
  echo "Error: .env not found in $(pwd)" >&2
  exit 1
fi

if [[ $# -eq 1 ]]; then
  TARGET="$1"
elif [[ $# -eq 0 && -f .previous_version ]]; then
  TARGET=$(cat .previous_version)
else
  usage
fi

if [[ -z "$TARGET" ]]; then
  echo "Error: target version is empty" >&2
  exit 1
fi

CURRENT=$(grep -E '^WEDDING_VERSION=' .env | cut -d= -f2- || true)
if [[ "$CURRENT" == "$TARGET" ]]; then
  echo "Already at $TARGET, nothing to do."
  exit 0
fi

echo "Rolling back: $CURRENT → $TARGET"

docker pull "$IMAGE:$TARGET"
sed -i "s/^WEDDING_VERSION=.*/WEDDING_VERSION=${TARGET}/" .env
docker compose up -d "$SERVICE"

echo "Done: rolled back to $TARGET"
