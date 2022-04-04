#!/usr/bin/env bash
set -Eeuo pipefail

# wait for postgres to be available
export PGPASSWORD="$I2B2_DB_PW"
export PSQL_PARAMS="-v ON_ERROR_STOP=1 -h ${I2B2_DB_HOST} -p ${I2B2_DB_PORT} -U ${I2B2_DB_USER}"
until psql $PSQL_PARAMS -d postgres -c '\q'; do
  >&2 echo "Waiting for postgresql..."
  sleep 1
done

# load initial data if database does not exist (credentials must be valid and have create database right)
DB_CHECK=$(psql ${PSQL_PARAMS} -d postgres -X -A -t -c "select count(*) from pg_database where datname = '${I2B2_DB_NAME}';")
if [[ "$DB_CHECK" -ne "1" ]]; then
echo "Initialising i2b2 database"
    psql $PSQL_PARAMS -d postgres <<-EOSQL
        CREATE DATABASE ${I2B2_DB_NAME};
EOSQL

    # uncompress the i2b2 data
    tar -xf "$I2B2_COMPRESSED_DATA_DIR/i2b2-data.tar.gz" -C "$I2B2_DATA_DIR" --strip-components 1

    # run loading scripts
    for f in "$I2B2_SQL_DIR"/*.sh; do
        bash "$f"
    done

    # delete loaded data
    rm -rf "$I2B2_DATA_DIR"/.git "$I2B2_DATA_DIR"/*

fi

# execute pre-init scripts & run wildfly
for f in "$PRE_INIT_SCRIPT_DIR"/*.sh; do
    bash "$f"
done
exec /opt/jboss/wildfly/bin/standalone.sh -b 0.0.0.0 -bmanagement 0.0.0.0
