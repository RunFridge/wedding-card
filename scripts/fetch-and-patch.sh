#!/bin/bash
set -euo pipefail

REPO="hwhang0917/wedding-card"
IMAGE="wedding-server"

for cmd in gh docker; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "Error: $cmd is not installed" >&2
    exit 1
  fi
done

if ! gh auth status &>/dev/null; then
  echo "Error: gh CLI is not authenticated. Run 'gh auth login'" >&2
  exit 1
fi

if ! docker info &>/dev/null; then
  echo "Error: Docker daemon is not running" >&2
  exit 1
fi

TAG=$(gh release view --repo "$REPO" --json tagName --jq '.tagName' 2>/dev/null) || {
  echo "Error: no releases found for $REPO" >&2
  exit 1
}

VERSION="${TAG#v}"
CURRENT=$(docker inspect --format '{{index .RepoDigests 0}}' "$IMAGE:latest" 2>/dev/null || echo "none")

echo "Latest release: $TAG"
echo "Current image:  $CURRENT"

OLD_IMAGE_ID=$(docker images --format '{{.ID}}' "$IMAGE:latest" 2>/dev/null || echo "")

TARBALL="/tmp/${IMAGE}-${VERSION}.tar.gz"
cleanup() { rm -f "$TARBALL"; }
trap cleanup EXIT

echo "Downloading $REPO $TAG"
gh release download "$TAG" --repo "$REPO" --pattern "*.tar.gz" --dir /tmp --clobber

if [ ! -f "$TARBALL" ]; then
  echo "Error: download failed, $TARBALL not found" >&2
  exit 1
fi

echo "Loading $TARBALL"
docker image load -i "$TARBALL"
docker tag "$IMAGE:$VERSION" "$IMAGE:latest"

echo "Updating docker-compose.yml"
curl -fsSL -o docker-compose.yml "https://raw.githubusercontent.com/${REPO}/main/docker-compose.yml"

echo "Restarting containers"
docker compose up -d

if [ -n "$OLD_IMAGE_ID" ]; then
  NEW_IMAGE_ID=$(docker images --format '{{.ID}}' "$IMAGE:latest")
  if [ "$OLD_IMAGE_ID" != "$NEW_IMAGE_ID" ]; then
    echo "Removing old image $OLD_IMAGE_ID"
    docker rmi "$OLD_IMAGE_ID" 2>/dev/null || true
  fi
fi

echo "Done: $IMAGE:$VERSION deployed"
