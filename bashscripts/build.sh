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

SETTINGS_FILE="settings/settings.${ENV_TAG}.jsonc"

# ── 1. Read env_file_name from settings ──────────────────────────────────────
if [[ ! -f "${PROJECT_ROOT}/${SETTINGS_FILE}" ]]; then
  echo "ERROR: Settings file not found at ${PROJECT_ROOT}/${SETTINGS_FILE}" >&2
  exit 1
fi

# Extract the value of env_file_name (handles both 'key: value' and 'key:value')
ENV_FILE_NAME=$(grep -E '^\s*"use_env"\s*:' "${PROJECT_ROOT}/${SETTINGS_FILE}" \
  | sed 's/.*"use_env"\s*:\s*"\([^"]*\)".*/\1/' \
  | tr -d '[:space:]')

if [[ -z "${ENV_FILE_NAME}" ]]; then
  echo "ERROR: 'use_env' field not found in ${SETTINGS_FILE}" >&2
  exit 1
fi

if [[ ! -f "${PROJECT_ROOT}/${ENV_FILE_NAME}" ]]; then
  echo "ERROR: env file '${ENV_FILE_NAME}' not found in project root" >&2
  exit 1
fi

echo "→ Copying 'settings.${ENV_TAG}.jsonc' to 'settings.jsonc' ..."
cp "${PROJECT_ROOT}/settings/settings.${ENV_TAG}.jsonc" "${PROJECT_ROOT}/settings/settings.jsonc"
echo "  Done."

# ── 2. Build ──────────────────────────────────────────────────────────────────
mkdir -p "${PROJECT_ROOT}/bin"

case "$(uname -s)" in
  MINGW*|MSYS*|CYGWIN*|Windows_NT)
    echo "→ Building ots (windows/amd64) ..."
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o "${PROJECT_ROOT}/bin/ots.exe"
    echo "  Build complete → bin/ots.exe"
    ;;
  *)
    echo "→ Building ots (linux/amd64) ..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "${PROJECT_ROOT}/bin/ots"
    echo "  Build complete → bin/ots"
    ;;
esac