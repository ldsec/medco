#!/bin/bash
set -e

# copy pm cell deployment to the right place in the container
cp -a -rf $BUILD_DIR/shrine-server/* $CATALINA_HOME
cp -a -rf $BUILD_DIR/shrine-server/src/. /opt/shrine-src/

sed -i "s/SHRINE_KEYSTORE_PRIVATE_KEY_ALIAS/srv$NODE_IDX-private/g" "$CATALINA_HOME/conf/server.xml"
sed -i "s#SHRINE_KEYSTORE_FILE_PATH#$CONF_DIR/srv$NODE_IDX.keystore#g" "$CATALINA_HOME/conf/server.xml"
sed -i "s#FINE#$SHRINE_DEBUG_LEVEL#g" "$CATALINA_HOME/conf/logging.properties"
sed -i "s#INFO#$SHRINE_DEBUG_LEVEL#g" "$CATALINA_HOME/conf/logging.properties"