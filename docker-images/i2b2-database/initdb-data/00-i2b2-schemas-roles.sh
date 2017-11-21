#!/bin/bash
set -e
# create schemas and database users for i2b2

function initSchema {
    DB_NAME="$1"
    SCHEMA_NAME="$2"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$DB_NAME" <<-EOSQL
        create schema '$SCHEMA_NAME';
        create role '$SCHEMA_NAME' login password '$DB_PASSWORD';
        grant all on schema '$SCHEMA_NAME' to '$SCHEMA_NAME';
        grant all privileges on all tables in schema '$SCHEMA_NAME' to '$SCHEMA_NAME';
EOSQL
}

# create the databases
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
        CREATE DATABASE '$I2B2_DEMO_DB_NAME';
        CREATE DATABASE '$I2B2_MEDCO_DB_NAME';
EOSQL

# init demo i2b2 database
initSchema $I2B2_DEMO_DB_NAME i2b2demodata
initSchema $I2B2_DEMO_DB_NAME i2b2imdata
initSchema $I2B2_DEMO_DB_NAME i2b2metadata
initSchema $I2B2_DEMO_DB_NAME i2b2workdata
initSchema $I2B2_DEMO_DB_NAME i2b2pm
initSchema $I2B2_DEMO_DB_NAME i2b2hive

# init medco i2b2 database
initSchema $I2B2_MEDCO_DB_NAME i2b2demodata
initSchema $I2B2_MEDCO_DB_NAME i2b2imdata
initSchema $I2B2_MEDCO_DB_NAME i2b2metadata
initSchema $I2B2_MEDCO_DB_NAME i2b2workdata
initSchema $I2B2_MEDCO_DB_NAME shrine_ont
initSchema $I2B2_MEDCO_DB_NAME genomic_annotations
