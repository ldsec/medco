#!/bin/bash
set -Eeuo pipefail

# set i2b2 service password (in config)
pushd "$JBOSS_HOME/standalone/configuration"
sed -i "/edu.harvard.i2b2.crc.pm.serviceaccount.password/c\edu.harvard.i2b2.crc.pm.serviceaccount.password=$I2B2_SERVICE_PASSWORD" crcapp/crc.properties
sed -i "/edu.harvard.i2b2.ontology.pm.serviceaccount.password/c\edu.harvard.i2b2.ontology.pm.serviceaccount.password=$I2B2_SERVICE_PASSWORD" ontologyapp/ontology.properties
popd

# set i2b2 service password (in jar)
pushd "$JBOSS_HOME/standalone/deployments/i2b2.war/WEB-INF/lib"

jar -xf CRC-Server.jar crc.properties
sed -i "/edu.harvard.i2b2.crc.pm.serviceaccount.password/c\edu.harvard.i2b2.crc.pm.serviceaccount.password=$I2B2_SERVICE_PASSWORD" crc.properties
jar -uf CRC-Server.jar crc.properties
rm crc.properties

jar -xf Ontology-Server.jar ontology.properties
sed -i "/edu.harvard.i2b2.ontology.pm.serviceaccount.password/c\edu.harvard.i2b2.ontology.pm.serviceaccount.password=$I2B2_SERVICE_PASSWORD" ontology.properties
jar -uf Ontology-Server.jar ontology.properties
rm ontology.properties

popd

# set i2b2 log level
pushd  "$JBOSS_HOME/standalone/deployments/i2b2.war/WEB-INF/classes"
sed -i "/^log4j.rootCategory=/c\log4j.rootCategory=$AXIS2_LOGLEVEL, CONSOLE" log4j.properties
popd
