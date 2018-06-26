#!/bin/bash
set -e

rm -rf $BUILD_DIR/i2b2-web/
mkdir $BUILD_DIR/i2b2-web/
# copy pm cell deployment data to shared volume
cp -a -R $LIGHTTPD_WEB_ROOT $BUILD_DIR/i2b2-web/