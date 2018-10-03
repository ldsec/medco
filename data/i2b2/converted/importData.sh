#!/usr/bin/env bash
PGPASSWORD=prigen2017 psql -v ON_ERROR_STOP=1 -h "localhost" -U "postgres" -p 5434 -d "i2b2medco" <<-EOSQL
BEGIN;
TRUNCATE TABLE shrine_ont.shrine;
TRUNCATE TABLE i2b2metadata.i2b2;
TRUNCATE TABLE i2b2metadata.sensitive_tagged;
TRUNCATE TABLE i2b2demodata.concept_dimension;
TRUNCATE TABLE i2b2demodata.modifier_dimension;
TRUNCATE TABLE i2b2demodata.patient_dimension;
TRUNCATE TABLE i2b2demodata.visit_dimension;
TRUNCATE TABLE i2b2demodata.observation_fact;

\copy shrine_ont.shrine FROM 'shrine.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2metadata.i2b2 FROM 'i2b2.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2metadata.sensitive_tagged FROM 'sensitive_tagged.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2demodata.concept_dimension FROM 'concept_dimension.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2demodata.modifier_dimension FROM 'modifier_dimension.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2demodata.patient_dimension FROM 'patient_dimension.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2demodata.visit_dimension FROM 'visit_dimension.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
\copy i2b2demodata.observation_fact FROM 'observation_fact.csv' ESCAPE '"' DELIMITER ',' CSV HEADER;
COMMIT;
EOSQL