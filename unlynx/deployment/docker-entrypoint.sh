#!/usr/bin/env bash
set -Eeuo pipefail

# export environment variables
export  UNLYNX_KEY_FILE_PATH="$MEDCO_CONF_DIR/srv$NODE_IDX-private.toml" \
        UNLYNX_DDT_SECRETS_FILE_PATH="$MEDCO_CONF_DIR/srv$NODE_IDX-ddtsecrets.toml"

# run unlynx
if [[ $# -eq 0 ]]; then
    ARGS="-d $UNLYNX_DEBUG_LEVEL server -c $UNLYNX_KEY_FILE_PATH"
else
    ARGS=$@
fi

exec medco-unlynx ${ARGS}
