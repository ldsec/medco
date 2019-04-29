#!/usr/bin/env bash
set -Eeuo pipefail

# wait for postgres to be available
export PGPASSWORD="${PICSURE2_PGDB_PW}"
PSQL_PARAMS="-v ON_ERROR_STOP=1 -h ${PICSURE2_PGDB_ADDRESS} -p ${PICSURE2_PGDB_PORT} -U ${PICSURE2_PGDB_USER}"
until psql ${PSQL_PARAMS} -d postgres -c '\q'; do
  >&2 echo "Waiting for postgresql..."
  sleep 1
done

# load initial data if database does not exist (credentials must be valid and have create database right)
DB_CHECK=`psql ${PSQL_PARAMS} -d postgres -X -A -t -c "select count(*) from pg_database where datname = '${PICSURE2_PGDB_DB}';"`
if [[ "$DB_CHECK" -ne "1" ]]; then
echo "Initialising PIC-SURE database"
    psql ${PSQL_PARAMS} -d postgres <<-EOSQL
        CREATE DATABASE ${PICSURE2_PGDB_DB};
        ALTER DATABASE ${PICSURE2_PGDB_DB} OWNER TO ${PICSURE2_PGDB_USER};
EOSQL

    psql ${PSQL_PARAMS} -d ${PICSURE2_PGDB_DB} <<-EOSQL
        grant all on schema public to ${PICSURE2_PGDB_DB};
        grant all privileges on all tables in schema public to ${PICSURE2_PGDB_USER};
EOSQL

    psql ${PSQL_PARAMS} -d ${PICSURE2_PGDB_DB} -f "/sql/picsure2_db_ddl.sql"
    psql ${PSQL_PARAMS} -d ${PICSURE2_PGDB_DB} -f "/sql/medco_resource_function.sql"
fi

# register in the database the medco resources
for (( IDX=0; ; IDX++ )); do
    MEDCO_NODE_NAME=$(echo "${MEDCO_NODES_NAME}" | cut -f$(($IDX+1)) -d,)
    MEDCO_NODE_CONNECTOR_URL=$(echo "${MEDCO_NODES_CONNECTOR_URL}" | cut -f$(($IDX+1)) -d,)

    if [[ -z ${MEDCO_NODE_NAME} ]]; then
        break
    fi

    UUID=$(uuidgen)
    DESC="MedCo node ${IDX} (${MEDCO_NODE_NAME}) from network ${MEDCO_NETWORK_NAME}"
    NAME="MEDCO_${MEDCO_NETWORK_NAME}_${IDX}_${MEDCO_NODE_NAME}"
    RSPATH="${MEDCO_NODE_CONNECTOR_URL}"
    TARGETURL=""
    TOKEN=""
    echo "select add_or_update_medco_resource('${UUID}', '${DESC}', '${NAME}', '${RSPATH}', '${TARGETURL}', '${TOKEN}');"
    psql ${PSQL_PARAMS} -d ${PICSURE2_PGDB_DB} <<-EOSQL
        select add_or_update_medco_resource('${UUID}', '${DESC}', '${NAME}', '${RSPATH}', '${TARGETURL}', '${TOKEN}');
EOSQL
done

exec /opt/jboss/wildfly/bin/standalone.sh -b 0.0.0.0 -bmanagement 0.0.0.0
