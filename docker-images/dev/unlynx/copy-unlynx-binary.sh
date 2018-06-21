#!/bin/bash
set -e

# copy unlynx binary in the configuration folder (environment variables are available)
rm -f $BIN_EXPORT_PATH
cp -a $(which app) $BIN_EXPORT_PATH
chmod 777 $BIN_EXPORT_PATH
