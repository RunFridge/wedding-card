#!/bin/bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
COUNT=5
PASSWORD="test1234"
IMAGE_DIR=""

NAMES=(
  "김민수" "이서연" "박지훈" "최유진" "정하늘"
  "강도윤" "조예린" "윤시우" "한소희" "임재현"
  "송다은" "오현우" "배수빈" "신지민" "문태양"
  "황채원" "류건우" "전하린" "안서준" "권나윤"
)

while [[ $# -gt 0 ]]; do
  case $1 in
    -n|--n) COUNT="$2"; shift 2 ;;
    -i|--image-dir) IMAGE_DIR="$2"; shift 2 ;;
    *) echo "Usage: $0 [-n COUNT] [-i IMAGE_DIR]"; exit 1 ;;
  esac
done

if ! command -v curl &>/dev/null; then
  echo "Error: curl is not installed" >&2
  exit 1
fi

# Collect image files from directory
IMAGES=()
if [[ -n "$IMAGE_DIR" ]]; then
  if [[ ! -d "$IMAGE_DIR" ]]; then
    echo "Error: $IMAGE_DIR is not a directory" >&2
    exit 1
  fi
  while IFS= read -r -d '' f; do
    IMAGES+=("$f")
  done < <(find "$IMAGE_DIR" -maxdepth 1 -type f \( -iname '*.png' -o -iname '*.jpg' -o -iname '*.jpeg' \) -print0)
  if [[ ${#IMAGES[@]} -eq 0 ]]; then
    echo "Error: no image files found in $IMAGE_DIR" >&2
    exit 1
  fi
  echo "Found ${#IMAGES[@]} images in $IMAGE_DIR"
fi

echo "Uploading $COUNT test photos to $BASE_URL/api/photos/upload"
echo ""

success=0
fail=0

for i in $(seq 1 "$COUNT"); do
  name="${NAMES[$((RANDOM % ${#NAMES[@]}))]}"
  image="${IMAGES[$((RANDOM % ${#IMAGES[@]}))]}"
  mime="image/jpeg"
  [[ "$image" == *.png ]] && mime="image/png"

  headers=$(mktemp)
  status=$(curl -s -o /dev/null -D "$headers" -w "%{http_code}" \
    -X POST "$BASE_URL/api/photos/upload" \
    -F "name=$name" \
    -F "password=$PASSWORD" \
    -F "image=@$image;type=$mime")

  # Retry on 429 with Retry-After
  while [[ "$status" == "429" ]]; do
    retry_after=$(grep -i 'retry-after' "$headers" | tr -d '\r' | awk '{print $2}')
    retry_after="${retry_after:-2}"
    echo "[$i/$COUNT] Rate limited, waiting ${retry_after}s..."
    sleep "$retry_after"
    status=$(curl -s -o /dev/null -D "$headers" -w "%{http_code}" \
      -X POST "$BASE_URL/api/photos/upload" \
      -F "name=$name" \
      -F "password=$PASSWORD" \
      -F "image=@$image;type=$mime")
  done
  rm -f "$headers"

  if [[ "$status" == "201" ]]; then
    echo "[$i/$COUNT] $name ($(basename "$image")) - OK ($status)"
    success=$((success + 1))
  else
    echo "[$i/$COUNT] $name ($(basename "$image")) - FAILED ($status)"
    fail=$((fail + 1))
  fi
done

echo ""
echo "Done: $success succeeded, $fail failed"
