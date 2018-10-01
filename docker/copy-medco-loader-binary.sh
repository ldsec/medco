#!/bin/bash
set -e

# copy unlynx binary in the configuration folder
rm -f $MEDCO_LOADER_BIN_EXPORT_PATH
cp -a $(which app) $MEDCO_LOADER_BIN_EXPORT_PATH
chmod 777 $MEDCO_LOADER_BIN_EXPORT_PATH
