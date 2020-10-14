#!/bin/bash
set -Eeuo pipefail
# create schemas for i2b2

function initSchema {
    DB_NAME="$1"
    SCHEMA_NAME="$2"
    psql $PSQL_PARAMS -d "$DB_NAME" <<-EOSQL
        create schema $SCHEMA_NAME;
        grant all on schema $SCHEMA_NAME to $I2B2_DB_USER;
        grant all privileges on all tables in schema $SCHEMA_NAME to $I2B2_DB_USER;
EOSQL
}

# init demo i2b2 database
initSchema $I2B2_DB_NAME i2b2demodata_i2b2_non_sensitive
initSchema $I2B2_DB_NAME i2b2demodata_i2b2_sensitive
initSchema $I2B2_DB_NAME i2b2imdata
initSchema $I2B2_DB_NAME i2b2metadata_i2b2
initSchema $I2B2_DB_NAME i2b2workdata
initSchema $I2B2_DB_NAME i2b2pm
initSchema $I2B2_DB_NAME i2b2hive
