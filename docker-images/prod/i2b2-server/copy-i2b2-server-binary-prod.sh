#!/bin/bash
set -e

# copy pm cell deployment to the right place in the container
cp -a -rf $BUILD_DIR/i2b2-server/standalone/* $JBOSS_HOME/standalone
cp -a $BUILD_DIR/i2b2-server/standalone.sh /opt/jboss/wildfly/bin/standalone.sh