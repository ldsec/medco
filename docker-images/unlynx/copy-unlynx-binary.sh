#!/bin/bash
set -e

# copy unlynx binary in the configuration folder (environment variables are available)
rm -f $UNLYNX_BIN_EXPORT_PATH
cp -a $(which unlynxMedCo) $UNLYNX_BIN_EXPORT_PATH
chmod 777 $UNLYNX_BIN_EXPORT_PATH
