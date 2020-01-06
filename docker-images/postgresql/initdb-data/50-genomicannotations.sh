#!/bin/bash
set -e

IFS=' ' read -r -a db_array <<< "$GA_DATABASES"

for element in "${db_array[@]}"
do
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE "$element";
    ALTER DATABASE "$element" OWNER TO genomicannotations;
EOSQL
psql -v ON_ERROR_STOP=1 -U genomicannotations -d "$element" <<-EOSQL
    CREATE SCHEMA genomic_annotations;
    GRANT ALL ON SCHEMA genomic_annotations TO genomicannotations;
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA genomic_annotations to genomicannotations;
EOSQL
done