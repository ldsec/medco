#!/bin/bash
set -e

# copy pm cell deployment data
cp -a -R $JBOSS_HOME/standalone/deployments $CONF_DIR/i2b2-server/

# copy i2b2-server binary in the configuration folder (environment variables are available)
cp -a /opt/jboss/wildfly/bin/standalone.sh $CONF_DIR/i2b2-server/