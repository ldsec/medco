#!/usr/bin/env bash
set -Eeuo pipefail

if [[ "$1" = "medco-connector-server" ]]; then
# trust the certificates of other nodes
if [[ `ls -1 /medco-configuration/srv*-certificate.crt 2>/dev/null | wc -l` != 0 ]]; then
    /bin/cp -f /medco-configuration/srv*-certificate.crt /usr/local/share/ca-certificates/
    update-ca-certificates
fi

# wait for postgres to be available
export PGPASSWORD="$GA_DB_PW"
export PSQL_PARAMS="-v ON_ERROR_STOP=1 -h ${GA_DB_HOST} -p ${GA_DB_PORT} -U ${GA_DB_USER}"
until psql $PSQL_PARAMS -d postgres -c '\q'; do
  >&2 echo "Waiting for postgresql..."
  sleep 1
done

# initialize database if it does not exist (credentials must be valid and have create database right)
DB_CHECK=$(psql ${PSQL_PARAMS} -d postgres -X -A -t -c "select count(*) from pg_database where datname = '${GA_DB_NAME}';")
if [[ "$DB_CHECK" -ne "1" ]]; then
echo "Initialising genomic_annotations database"
    psql $PSQL_PARAMS -d postgres <<-EOSQL
        CREATE DATABASE ${GA_DB_NAME};
EOSQL
psql $PSQL_PARAMS -d "$GA_DB_NAME" <<-EOSQL
        CREATE SCHEMA genomic_annotations;
EOSQL
fi
fi

exec $@
