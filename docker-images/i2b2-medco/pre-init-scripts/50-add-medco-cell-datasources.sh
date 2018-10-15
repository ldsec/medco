#!/bin/bash
set -Eeuo pipefail

cat > "$JBOSS_WAR_DEPLOYMENTS/medco-ds.xml" <<EOL
<?xml version="1.0" encoding="UTF-8"?>
<datasources xmlns="http://www.jboss.org/ironjacamar/schema">
    <datasource jta="false" jndi-name="java:/OntologyMedCoDS"
                pool-name="OntologyMedCoDS" enabled="true" use-ccm="false">
        <connection-url>jdbc:postgresql://$I2B2_DB_HOST:$I2B2_DB_PORT/$I2B2_DB_NAME?currentSchema=medco_ont</connection-url>
        <driver-class>org.postgresql.Driver</driver-class>
        <driver>$PG_JDBC_JAR</driver>
        <security>
            <user-name>$I2B2_DB_USER</user-name>
            <password>$I2B2_DB_PW</password>
        </security>
        <validation>
            <validate-on-match>false</validate-on-match>
            <background-validation>false</background-validation>
        </validation>
        <statement>
            <share-prepared-statements>false</share-prepared-statements>
        </statement>
    </datasource>
</datasources>

EOL
