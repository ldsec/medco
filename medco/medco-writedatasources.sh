#!/bin/bash
set -e

#TODO: check file OK + password
# meant to be called by Dockerfile of medco
# env var used: I2B2_DOMAIN_NAME

cat > ch.epfl.lca1.medco/etc/jboss/medco-ds.xml <<EOL
<?xml version="1.0" encoding="UTF-8"?>
<datasources xmlns="http://www.jboss.org/ironjacamar/schema">
    <datasource jta="false" jndi-name="java:/MedCoBootStrapDS"
            pool-name="MedCoBootStrapDS" enabled="true" use-ccm="false">
                <connection-url>jdbc:postgresql://i2b2-database:5432/$I2B2_DOMAIN_NAME</connection-url>
                <driver-class>org.postgresql.Driver</driver-class>
                <driver>postgresql-9.2-1002.jdbc4.jar</driver>
                <security>
                        <user-name>medco</user-name>
                        <password>demouser</password>
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
