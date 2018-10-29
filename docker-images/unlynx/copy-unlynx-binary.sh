#!/usr/bin/env bash
set -Eeuo pipefail

# copy unlynx binary in the configuration folder
if [[ (-f "$UNLYNX_BIN_EXPORT_PATH" && $(md5 /go/bin/unlynxMedCo) != $(md5 $UNLYNX_BIN_EXPORT_PATH) ) ]]; then
    rm -f $UNLYNX_BIN_EXPORT_PATH
fi

if [[ ! -f "$UNLYNX_BIN_EXPORT_PATH" ]]; then
    cp -a /go/bin/unlynxMedCo $UNLYNX_BIN_EXPORT_PATH
    chmod 777 $UNLYNX_BIN_EXPORT_PATH
fi
