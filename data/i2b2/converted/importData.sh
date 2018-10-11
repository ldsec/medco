#!/usr/bin/env bash
PGPASSWORD=prigen2017 psql -v ON_ERROR_STOP=1 -h "localhost" -U "postgres" -p 5434 -d "i2b2medco" <<-EOSQL
BEGIN;
TRUNCATE TABLE i2b2demodata.concept_dimension;
TRUNCATE TABLE i2b2demodata.patient_dimension;
TRUNCATE TABLE i2b2demodata.visit_dimension;
TRUNCATE TABLE i2b2demodata.observation_fact;

TRUNCATE TABLE i2b2metadata.table_access;
TRUNCATE TABLE i2b2metadata.schemes;
TRUNCATE TABLE i2b2metadata.i2b2;
TRUNCATE TABLE i2b2metadata.birn;
TRUNCATE TABLE i2b2metadata.icd10_icd9;
TRUNCATE TABLE i2b2metadata.custom_meta;
TRUNCATE TABLE i2b2metadata.sensitive_tagged;

CREATE TABLE shrine_ont.i2b2 AS SELECT * FROM i2b2metadata.i2b2 WHERE 1=2;
CREATE TABLE shrine_ont.birn AS SELECT * FROM i2b2metadata.birn WHERE 1=2;
CREATE TABLE shrine_ont.custom_meta AS SELECT * FROM i2b2metadata.custom_meta WHERE 1=2;

\copy shrine_ont.table_access FROM 'shrine_table_access.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy shrine_ont.i2b2 FROM 'shrine_i2b2.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy shrine_ont.birn FROM 'shrine_birn.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy shrine_ont.custom_meta FROM 'shrine_custom_meta.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;

\copy i2b2metadata.table_access FROM 'local_table_access.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2metadata.i2b2 FROM 'local_i2b2.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2metadata.birn FROM 'local_birn.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2metadata.custom_meta FROM 'local_custom_meta.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2metadata.sensitive_tagged FROM 'sensitive_tagged.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;

\copy i2b2demodata.concept_dimension FROM 'concept_dimension.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2demodata.patient_dimension FROM 'patient_dimension.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2demodata.visit_dimension FROM 'visit_dimension.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2demodata.observation_fact FROM 'observation_fact.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
COMMIT;
EOSQL