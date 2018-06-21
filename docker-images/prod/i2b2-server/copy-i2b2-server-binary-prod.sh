#!/bin/bash
set -e

# copy pm cell deployment data
cp -a -R $CONF_DIR/i2b2-server/deployments/* $JBOSS_HOME/standalone/deployments

# copy i2b2-server binary in the configuration folder (environment variables are available)
cp -a $CONF_DIR/i2b2-server/standalone.sh /opt/jboss/wildfly/bin/standalone.sh