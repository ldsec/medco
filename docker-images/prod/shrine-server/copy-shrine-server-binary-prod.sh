#!/bin/bash
set -e

# copy pm cell deployment to the right place in the container
cp -a -rf $BUILD_DIR/shrine-server/* $CATALINA_HOME

cp /root/shrine.conf "$CATALINA_HOME/lib/"
cp /root/server.xml /root/context.xml "$CATALINA_HOME/conf/"

wget "$SHRINE_MYSQL_JAR_URL" -P "$CATALINA_HOME/lib/"
wget "$SHRINE_ADAPTER_MAPPINGS_URL" -O "$CATALINA_HOME/lib/AdapterMappings.xml"
sed -i "s#SHRINE_DOWNSTREAM_NODES_FILE_PATH#$CONF_DIR/shrine_downstream_nodes.conf#g" "$CATALINA_HOME/lib/shrine.conf" && \
sed -i "s#SHRINE_CA_CERT_ALIASES_FILE_PATH#$CONF_DIR/shrine_ca_cert_aliases.conf#g" "$CATALINA_HOME/lib/shrine.conf" && \
sed -i "s/SHRINE_KEYSTORE_PASSWORD/$ADMIN_PASSWORD/g" "$CATALINA_HOME/conf/server.xml" && \
sed -i "s/SHRINE_DB_PASSWORD/$DB_PASSWORD/g" "$CATALINA_HOME/conf/context.xml"

sed -i "s/SHRINE_KEYSTORE_PRIVATE_KEY_ALIAS/srv$NODE_IDX-private/g" "$CATALINA_HOME/conf/server.xml"
sed -i "s#SHRINE_KEYSTORE_FILE_PATH#$CONF_DIR/srv$NODE_IDX.keystore#g" "$CATALINA_HOME/conf/server.xml"
sed -i "s#FINE#$SHRINE_DEBUG_LEVEL#g" "$CATALINA_HOME/conf/logging.properties"
sed -i "s#INFO#$SHRINE_DEBUG_LEVEL#g" "$CATALINA_HOME/conf/logging.properties"