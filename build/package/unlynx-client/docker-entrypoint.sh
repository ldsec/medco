#!/usr/bin/env bash
set -Eeuo pipefail

EXEC=$@
# export environment variables
export  UNLYNX_KEY_FILE_PATH="/medco-configuration/group.toml"

EXEC="${EXEC} -d $UNLYNX_DEBUG_LEVEL run"
exec $EXEC