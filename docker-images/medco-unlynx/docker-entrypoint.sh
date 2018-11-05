#!/usr/bin/env bash
set -Eeuo pipefail

# export environment variables
export  UNLYNX_KEY_FILE_PATH="$MEDCO_CONF_DIR/srv$NODE_IDX-private.toml" \
        UNLYNX_DDT_SECRETS_FILE_PATH="$MEDCO_CONF_DIR/srv$NODE_IDX-ddtsecrets.toml"

# copy unlynx binary in the configuration folder
if [[ (-f "$UNLYNX_BIN_EXPORT_PATH" && $(md5 /go/bin/unlynxMedCo) != $(md5 $UNLYNX_BIN_EXPORT_PATH) ) ]]; then
    rm -f $UNLYNX_BIN_EXPORT_PATH
fi

if [[ ! -f "$UNLYNX_BIN_EXPORT_PATH" ]]; then
    cp -a /go/bin/unlynxMedCo $UNLYNX_BIN_EXPORT_PATH
    chmod 777 $UNLYNX_BIN_EXPORT_PATH
fi

# run unlynx
exec /go/bin/unlynxMedCo -d $UNLYNX_DEBUG_LEVEL server -c $UNLYNX_KEY_FILE_PATH
