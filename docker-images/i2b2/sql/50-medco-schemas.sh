#!/bin/bash
set -Eeuo pipefail
# create schemas for medco ontology and genomic annotations

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
initSchema $I2B2_DB_NAME medco_ont
initSchema $I2B2_DB_NAME genomic_annotations
