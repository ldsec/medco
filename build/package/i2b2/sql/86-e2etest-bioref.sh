#!/bin/bash
set -Eeuo pipefail
# Set up data in the database for end-to-end tests for bioref
# Create 10 observations that are linked to patients from a same cohort.
# 2 observations are linked to the same patients.
# The cohort is created in the file bioref-cohort.sh


psql $PSQL_PARAMS -d "$I2B2_DB_NAME" <<-EOSQL

    -- maybe a cleaner way of creating test would be to create a new table called e2etest_bioref?
    -- if this doesn't destroy other tests we are fine


    insert into medco_ont.e2etest
        (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
        c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator,
        c_dimcode, c_comment, c_tooltip, update_date, download_date, import_date,
        valuetype_cd, m_applied_path, c_basecode, c_metadataxml) values
            (
                '1', '\e2etest\bioref\', 'E2E bioref concept', 'N', 'LA', '0',
                'concept_cd', 'concept_dimension', 'concept_path',
                'T', 'LIKE', '\e2etest\bioref\', 'E2E bioref concept', '\e2etest\bioref\',
                'NOW()', 'NOW()', 'NOW()', 'ENC_ID', '@', 'BIOREF:1', '<?xml version="1.0"?><ValueMetadata></ValueMetadata>'
            );



    -- i2b2demodata_i2b2.concept_dimension
    insert into i2b2demodata_i2b2.concept_dimension
        (concept_path, concept_cd, import_date, upload_id) values
            ('\e2etest\bioref\', 'BIOREF:1', 'NOW()', '1');

    --OMITTED
    -- i2b2demodata_i2b2.visit_dimension

    -- i2b2demodata_i2b2.encounter_mapping


    -- i2b2demodata_i2b2.observation_fact
    -- we insert 10 values for the same concept BIOREF:1 those observations have numerical values ranging from 0 to 10
    insert into i2b2demodata_i2b2.observation_fact
        (encounter_num, patient_num, concept_cd, provider_id, start_date, modifier_cd, instance_num, import_date, upload_id, valtype_cd, tval_char, nval_num, units_cd) values
            ('100', '1', 'BIOREF:1', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1', 'N', 'E', '0', 'test unit'),
            ('101', '2', 'BIOREF:1', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1', 'N', 'E', '0.5', 'test unit'),
            ('102', '3', 'BIOREF:1', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1', 'N', 'E', '1', 'test unit'),

            ('103', '4', 'BIOREF:1', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1', 'N', 'E', '4', 'test unit'),
            ('104', '5', 'BIOREF:1', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1', 'N', 'E', '4', 'test unit'),
            ('105', '5', 'BIOREF:1', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1', 'N', 'E', '4.5', 'test unit'), --two values for the same patient
            ('106', '6', 'BIOREF:1', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1', 'N', 'E', '4.7', 'test unit'),


            ('107', '7', 'BIOREF:1', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1', 'N', 'E', '8', 'test unit'),
            ('108', '8', 'BIOREF:1', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1', 'N', 'E', '8.1', 'test unit'),
            ('109', '9', 'BIOREF:1', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1', 'N', 'E', '8.6', 'test unit');

EOSQL

