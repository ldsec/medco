#!/usr/bin/env bash
set -Eeuo pipefail

if [[ "$1" = "medco-connector-server" ]]; then
# trust the certificates of other nodes
if [[ `ls -1 /medco-configuration/srv*-certificate.crt 2>/dev/null | wc -l` != 0 ]]; then
    /bin/cp -f /medco-configuration/srv*-certificate.crt /usr/local/share/ca-certificates/
    update-ca-certificates
fi

# wait for postgres to be available
export SCHEMA="query_tools"
export PGPASSWORD="$MC_DB_PW"
export PSQL_PARAMS="-v ON_ERROR_STOP=1 -h ${MC_DB_HOST} -p ${MC_DB_PORT} -U ${MC_DB_USER}"
until psql $PSQL_PARAMS -d postgres -c '\q'; do
  >&2 echo "Waiting for postgresql..."
  sleep 1
done

# initialize database if it does not exist (credentials must be valid and have create database right)
DB_CHECK=$(psql ${PSQL_PARAMS} -d postgres -X -A -t -c "select count(*) from pg_database where datname = '${MC_DB_NAME}';")
if [[ "$DB_CHECK" -ne "1" ]]; then
echo "Initialising medco connector database"
psql $PSQL_PARAMS -d postgres <<-EOSQL
    CREATE DATABASE ${MC_DB_NAME};
EOSQL
psql $PSQL_PARAMS -d "$MC_DB_NAME" <<-EOSQL
      CREATE SCHEMA genomic_annotations;
      CREATE SCHEMA '${SCHEMA}';
EOSQL
fi

#create the enum type
psql $PSQL_PARAMS -d "$MC_DB_NAME" <<-EOSQL
DO
\$\$
BEGIN
IF  NOT EXISTS (SELECT * FROM pg_type where typname='status_enum') AS type_count  THEN
CREATE TYPE status_enum AS ENUM ('running','completed','error');
END IF;
END
\$\$
EOSQL

#create explore query results
psql $PSQL_PARAMS -d "$MC_DB_NAME" <<-EOSQL
CREATE TABLE IF NOT EXISTS "${SCHEMA}".explore_query_results
(
    query_id serial NOT NULL,
    query_name character varying(255) NOT NULL,
    user_id character varying(255) NOT NULL,
    query_status status_enum  NOT  NULL,
    clear_result_set_size integer,
    clear_result_set integer[],
    query_definition text,
    i2b2_encrypted_patient_set_id integer,
    i2b2_non_encrypte_patient_set_id integer,
    PRIMARY KEY (query_id),
    UNIQUE (query_name)

)
EOSQL

#create cohorts table
psql $PSQL_PARAMS -d "$MC_DB_NAME"<<-EOSQL
CREATE TABLE IF NOT EXISTS "${SCHEMA}".saved_cohorts
(
    cohort_id serial NOT NULL,
    user_id character varying(255) NOT NULL,
    cohort_name character varying(255) NOT NULL,
    query_id INTEGER NOT  NULL,
    create_date TIMESTAMP WITHOUT TIME ZONE,
    update_date TIMESTAMP WITHOUT TIME ZONE,
    CONSTRAINT saved_cohorts_pkey PRIMARY KEY (cohort_id),
    CONSTRAINT saved_cohorts_user_id_cohort_name_key UNIQUE (user_id, cohort_name),
    CONSTRAINT query_tool_fk_psc_ri FOREIGN KEY (query_id)
        REFERENCES "${SCHEMA}".explore_query_results (query_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

EOSQL
fi


exec $@