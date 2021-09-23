#!/usr/bin/env bash
set -Eeuo pipefail

EXEC=$@
# export environment variables
export  UNLYNX_KEY_FILE_PATH="/medco-configuration/srv${MEDCO_NODE_IDX}-private.toml" \
        UNLYNX_DDT_SECRETS_FILE_PATH="/medco-configuration/srv${MEDCO_NODE_IDX}-ddtsecrets.toml"

EXEC="${EXEC} -d $UNLYNX_DEBUG_LEVEL server -c $UNLYNX_KEY_FILE_PATH"
exec $EXEC