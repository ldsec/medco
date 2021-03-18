#!/usr/bin/env bash
set -Eeuo pipefail
# example: bash medco-demo.epfl.ch i2b2medcosrv0 medcoconnectorsrv0

DB_HOST=$1
I2B2MEDCODB=$2
MEDCOCONNECTORDB=$3
SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPT_FOLDER/../test/data/spo-synthetic"
PGPASSWORD=i2b2 psql -v ON_ERROR_STOP=1 -h "${DB_HOST}" -U "i2b2" -p 5432 -d "${I2B2MEDCODB}" <<-EOSQL
BEGIN;
TRUNCATE TABLE i2b2demodata_i2b2.patient_mapping;
TRUNCATE TABLE i2b2demodata_i2b2.encounter_mapping;
TRUNCATE TABLE i2b2demodata_i2b2.concept_dimension;
TRUNCATE TABLE i2b2demodata_i2b2.modifier_dimension;
TRUNCATE TABLE i2b2demodata_i2b2.patient_dimension;
TRUNCATE TABLE i2b2demodata_i2b2.visit_dimension;
TRUNCATE TABLE i2b2demodata_i2b2.provider_dimension;
TRUNCATE TABLE i2b2demodata_i2b2.observation_fact;

\copy i2b2demodata_i2b2.patient_mapping FROM './PATIENT_MAPPING.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2demodata_i2b2.encounter_mapping FROM './ENCOUNTER_MAPPING.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2demodata_i2b2.concept_dimension FROM './CONCEPT_DIMENSION.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2demodata_i2b2.modifier_dimension FROM './MODIFIER_DIMENSION.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2demodata_i2b2.patient_dimension FROM './PATIENT_DIMENSION.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2demodata_i2b2.visit_dimension FROM './VISIT_DIMENSION.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2demodata_i2b2.provider_dimension FROM './PROVIDER_DIMENSION.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2demodata_i2b2.observation_fact FROM './observation_fact_spo.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;


TRUNCATE TABLE medco_ont.table_access;
\copy medco_ont.table_access FROM './table_access_clean.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;

--c_hlevel,c_fullname,c_name,c_synonym_cd,c_visualattributes,c_basecode,c_facttablecolumn,c_tablename,c_columnname,c_columndatatype,c_operator,c_comment,c_dimcode,c_tooltip,m_applied_path,c_totalnum,update_date,download_date,
--import_date,sourcesystem_cd,valuetype_cd,m_exclusion_cd,c_path,c_symbol,c_metadataxml

CREATE TABLE IF NOT EXISTS medco_ont.sphn (
        				C_HLEVEL NUMERIC(22,0),
        				C_FULLNAME VARCHAR(2000),
        				C_NAME VARCHAR(1000),
        				C_SYNONYM_CD CHAR(1),
        				C_VISUALATTRIBUTES CHAR(3),
        				C_BASECODE VARCHAR(450),
        				C_FACTTABLECOLUMN VARCHAR(50),
        				C_TABLENAME VARCHAR(50),
					C_COLUMNNAME VARCHAR(50),
        				C_COLUMNDATATYPE VARCHAR(50),
        				C_OPERATOR VARCHAR(10),
        				C_COMMENT TEXT,
        				C_DIMCODE VARCHAR(2000),
        				C_TOOLTIP VARCHAR(1000),
        				M_APPLIED_PATH VARCHAR(2000),
        				C_TOTALNUM NUMERIC(22,0),
        				UPDATE_DATE DATE,
					DOWNLOAD_DATE DATE,
        				IMPORT_DATE DATE,
        				SOURCESYSTEM_CD VARCHAR(50),
        				VALUETYPE_CD VARCHAR(50),
        				M_EXCLUSION_CD VARCHAR(1000),
					C_PATH VARCHAR(2000),
					C_SYMBOL VARCHAR(50),
        				C_METADATAXML TEXT);

ALTER TABLE medco_ont.sphn OWNER TO i2b2;
TRUNCATE TABLE medco_ont.sphn;
\copy medco_ont.sphn FROM './spo_onto_db.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;

COMMIT;
EOSQL

PGPASSWORD=medcoconnector psql -v ON_ERROR_STOP=1 -h "${DB_HOST}" -U "medcoconnector" -p 5432 -d "${MEDCOCONNECTORDB}" <<-EOSQL
BEGIN;
TRUNCATE TABLE query_tools.explore_query_results, query_tools.saved_cohorts;
COMMIT;
EOSQL
popd