#!/bin/bash
set -e

# copy unlynx binary in the configuration folder
rm -f $UNLYNX_BIN_EXPORT_PATH
cp -a $(which app) $UNLYNX_BIN_EXPORT_PATH
chmod 777 $UNLYNX_BIN_EXPORT_PATH
