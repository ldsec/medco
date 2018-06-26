#!/bin/bash
set -e

# copy pm cell deployment data to shared volume
cp -a -R $BUILD_DIR/i2b2-web/html/* $LIGHTTPD_WEB_ROOT 