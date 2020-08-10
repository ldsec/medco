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





# create medcoqt database for qeuery tool tables
export PGPASSWORD="$QT_DB_PASSWORD"
export PSQL_PARAMS="-v ON_ERROR_STOP=1 -h ${QT_DB_HOST} -p ${QT_DB_PORT} -U ${QT_DB_USER}"
until psql $PSQL_PARAMS -d postgres -c '\q'; do
  >&2 echo "Waiting for postgresql..."
  sleep 1
done

DB_CHECK=$(psql ${PSQL_PARAMS} -d postgres -X -A -t -c "select count(*) from pg_database where datname = '${QT_DB_NAME}';")

if  [[ "$DB_CHECK" -ne "1" ]]; then
echo  "Initialising query tool databse"
    psql $PSQL_PARAMS -d postgres <<-EOSQL
        CREATE DATABASE ${QT_DB_NAME};
EOSQL
SCHEMA=medco_query_tools
SCHEMA_CHECK=$(psql ${PSQL_PARAMS} -d postgres-X -A -t -c "select count(*) from information_schema.schemata where schema_name = '${SCHEMA}';")
if [[ "$SCHEMA_CHECK"  -ne "1"]];then
echo "creating medcoqt schema"
psql $PSQL_PARAMS -d postgres <<-EOSQL
       CREATE SCHEMA '${SCHEMA}'
EOSQL
fi

#create explore query results
psql $PSQL_PARAMS -d postgres <<-EOSQL
IF (SELECT COUNT(*) FROM pg_types where typname=status_enum) AS type_count =0 THEN
CREATE TYPE status_enum AS ENUM (‘running’,’completed’,’error’)
END IF;
CREATE TABLE IF NOT EXISTS '${SCHEMA}'.explore_query_results
(
    query_id serial NOT NULL,
    query_name character varying(255) NOT NULL,
    user_id character varying(255) NOT NULL,
    query_status status_enum  NOT  NULL,
    enc_result_set_size bytea,
    enc_result_set bytea,
    query_definition text,
    i2b2_encrypted_patient_set_id integer,
    i2b2_non_encrypte_patient_set_id integer,
    PRIMARY KEY (query_id),
    UNIQUE (query_name)

)
EOSQL

#create cohorts table
psql $PSQL_PARAMS -d postgres <<-EOSQL
CREATE TABLE IF NOT EXISTS '${SCHEMA}'.saved_cohorts
(
    cohort_id serial NOT NULL,
    user_id character varying(255) NOT NULL,
    cohort_name character varying(255) NOT NULL,
    query_id INTEGER NOT  NULL,
    create_date TIMESTAMP WITHOUT TIME ZONE,
    update_date TIMESTAMP WITHOUT TIME ZONE,
    CONSTRAINT saved_cohorts_pkey PRIMARY KEY (cohort_id),
    CONSTRAINT saved_cohorts_user_id_cohort_name_key UNIQUE (user_id, cohort_name),
    CONSTRAINT query_tool_fk_psc_ri FOREIGN KEY (qurey_id)
        REFERENCES '${SCHEMA}'.explore_query_results (query_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

EOSQL
fi

exec $@
