#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PRE_LOAD_ENV="${PROJECT_ROOT}/.env.pre"

echo "→ Loading '${PRE_LOAD_ENV}' ..."
set -a
source "${PRE_LOAD_ENV}"
set +a

if [[ -z "${ENV_TAG}" ]]; then
  echo "ERROR: ENV_TAG is not set in .env" >&2
  exit 1
fi

SETTINGS_FILE="${PROJECT_ROOT}/settings/settings.${ENV_TAG}.jsonc"

# ── 1. Build ──────────────────────────────────────────────────────────────────
echo "→ Running build.sh ..."
bash "${PROJECT_ROOT}/bashscripts/build.sh"

# ── 2. Load .env ──────────────────────────────────────────────────────────────
if [[ ! -f "${SETTINGS_FILE}" ]]; then
  echo "ERROR: Settings file not found at ${SETTINGS_FILE}" >&2
  exit 1
fi

ENV_FILE_NAME=$(grep -E '^\s*"use_env"\s*:' "${SETTINGS_FILE}" \
  | sed 's/.*"use_env"\s*:\s*"\([^"]*\)".*/\1/' \
  | tr -d '[:space:]')

if [[ -z "${ENV_FILE_NAME}" ]]; then
  echo "ERROR: 'use_env' field not found in ${SETTINGS_FILE}" >&2
  exit 1
fi

ENV_FILE="${PROJECT_ROOT}/${ENV_FILE_NAME}"

if [[ ! -f "${ENV_FILE}" ]]; then
  echo "ERROR: env file '${ENV_FILE_NAME}' not found in project root" >&2
  exit 1
fi

echo "→ Loading '${ENV_FILE_NAME}' ..."
set -a
# shellcheck source=/dev/null
source "${ENV_FILE}"
set +a

if [[ -z "${TAG:-}" ]]; then
  echo "ERROR: TAG is not set in ${ENV_FILE_NAME}" >&2
  exit 1
fi

if [[ -z "${USERNAME:-}" ]]; then
  echo "ERROR: USERNAME is not set in ${ENV_FILE_NAME}" >&2
  exit 1
fi

# ── 3. Docker build ───────────────────────────────────────────────────────────
echo "🔨 Building image with tag: ${TAG}"
docker build \
  -t "${USERNAME}/ots:${TAG}" \
  -f "${PROJECT_ROOT}/Dockerfile" \
  "${PROJECT_ROOT}"

echo "✅ Done! Image built: ${USERNAME}/ots:${TAG}"


# ── 4. Docker deploy ───────────────────────────────────────────────────────────
echo "📤 Pushing image to docker hub..."

docker push ${USERNAME}/ots:${TAG}

echo "✅ Done! Image pushed to docker hub"