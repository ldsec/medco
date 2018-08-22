#!/bin/bash
set -e

rm -rf $BUILD_DIR/shrine-server/
mkdir $BUILD_DIR/shrine-server/
# copy pm cell deployment data to shared volume
cp -a -R $CATALINA_HOME $BUILD_DIR/shrine-server/
ls /root/
cp -a -R /root/.m2 $BUILD_DIR/shrine-server/