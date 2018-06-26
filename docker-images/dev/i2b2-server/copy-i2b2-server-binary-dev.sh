#!/bin/bash
set -e

rm -rf $BUILD_DIR/i2b2-server/
mkdir $BUILD_DIR/i2b2-server/
# copy pm cell deployment data to shared volume
cp -a -R $JBOSS_HOME/standalone $BUILD_DIR/i2b2-server/
cp -a /opt/jboss/wildfly/bin/standalone.sh $BUILD_DIR/i2b2-server/