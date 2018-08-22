#!/bin/bash
set -e

rm -rf $BUILD_DIR/shrine-server/
mkdir $BUILD_DIR/shrine-server/
# copy pm cell deployment data to shared volume
cp -a -R $CATALINA_HOME $BUILD_DIR/shrine-server/
cp -a -R $SHRINE_SRC_DIR $BUILD_DIR/shrine-server/src/
