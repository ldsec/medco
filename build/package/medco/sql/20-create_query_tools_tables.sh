#!/usr/bin/env bash
set -Eeuo pipefail

# create the schema
psql $PSQL_PARAMS -d "$MC_DB_NAME" <<-EOSQL
CREATE SCHEMA query_tools;
EOSQL

#create the enum type
psql $PSQL_PARAMS -d "$MC_DB_NAME" <<-EOSQL
DO
\$\$
BEGIN
IF  NOT EXISTS (SELECT * FROM pg_type where typname='status_enum') AS type_count  THEN
CREATE TYPE status_enum AS ENUM ('running','completed','error','predefined');
END IF;
END
\$\$
EOSQL

#create explore query results
psql $PSQL_PARAMS -d "$MC_DB_NAME" <<-EOSQL
CREATE TABLE IF NOT EXISTS query_tools.explore_query_results
(
    query_id serial NOT NULL,
    query_name character varying(255) NOT NULL,
    user_id character varying(255) NOT NULL,
    query_status status_enum  NOT  NULL,
    clear_result_set_size integer,
    clear_result_set integer[],
    query_definition text,
    i2b2_encrypted_patient_set_id integer,
    i2b2_non_encrypted_patient_set_id integer,
    PRIMARY KEY (query_id),
    UNIQUE (query_name)

)
EOSQL

#create cohorts table
psql $PSQL_PARAMS -d "$MC_DB_NAME"<<-EOSQL
CREATE TABLE IF NOT EXISTS query_tools.saved_cohorts
(
    cohort_id serial NOT NULL,
    user_id character varying(255) NOT NULL,
    cohort_name character varying(255) NOT NULL,
    query_id INTEGER NOT  NULL,
    create_date TIMESTAMP WITHOUT TIME ZONE,
    update_date TIMESTAMP WITHOUT TIME ZONE,
    predefined BOOLEAN NOT NULL,
    default_flag BOOLEAN NOT NULL,
    CONSTRAINT saved_cohorts_pkey PRIMARY KEY (cohort_id),
    CONSTRAINT saved_cohorts_user_id_cohort_name_key UNIQUE (user_id, cohort_name),
    CONSTRAINT query_tool_fk_psc_ri FOREIGN KEY (query_id)
        REFERENCES query_tools.explore_query_results (query_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

EOSQL

