#!/usr/bin/env bash
set -Eeuo pipefail

EXEC=$@

# trust the certificates of other nodes
if [[ -d /medco-configuration ]]; then
  NB_CA_CERTS=$(find /medco-configuration -maxdepth 1 -name '*.crt' | wc -l)
  if [[ "$NB_CA_CERTS" != 0 ]]; then
    cp -f /medco-configuration/*.crt /usr/local/share/ca-certificates/
    >&2 update-ca-certificates
    >&2 echo "MedCo ($EXEC): $NB_CA_CERTS CA certificates added"
  else
    >&2 echo "MedCo ($EXEC): no CA certificate added"
  fi
fi

if [[ "$1" = "medco-connector-server" ]]; then

  # wait for postgres to be available
  export PGPASSWORD="$MC_DB_PASSWORD"
  export PSQL_PARAMS="-v ON_ERROR_STOP=1 -h ${MC_DB_HOST} -p ${MC_DB_PORT} -U ${MC_DB_USER}"
  until psql $PSQL_PARAMS -d postgres -c '\q'; do
    >&2 echo "Waiting for postgresql..."
    sleep 1
  done

  # initialize database if it does not exist (credentials must be valid and have create database right)
  DB_CHECK=$(psql ${PSQL_PARAMS} -d postgres -X -A -t -c "select count(*) from pg_database where datname = '${MC_DB_NAME}';")
  if [[ "$DB_CHECK" -ne "1" ]]; then
    echo "Initialising medco_connector database"
    psql $PSQL_PARAMS -d postgres <<-EOSQL
          CREATE DATABASE ${MC_DB_NAME};
EOSQL

    # run loading scripts
    for f in "$MC_SQL_DIR"/*.sh; do
        bash "$f"
    done
  fi

  EXEC="${EXEC} --write-timeout=${SERVER_HTTP_WRITE_TIMEOUT_SECONDS}s"

elif [[ "$1" = "medco-unlynx" ]]; then
  # export environment variables
  export  UNLYNX_KEY_FILE_PATH="/medco-configuration/srv${MEDCO_NODE_IDX}-private.toml" \
          UNLYNX_DDT_SECRETS_FILE_PATH="/medco-configuration/srv${MEDCO_NODE_IDX}-ddtsecrets.toml"

  # run unlynx
  if [[ $# -eq 1 ]]; then
      EXEC="${EXEC} -d $UNLYNX_DEBUG_LEVEL server -c $UNLYNX_KEY_FILE_PATH"
  fi
fi

exec $EXEC
