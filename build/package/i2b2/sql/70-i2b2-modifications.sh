#!/bin/bash
set -Eeuo pipefail
# set up common medco ontology

psql $PSQL_PARAMS -d "$I2B2_DB_NAME" <<-EOSQL

    -- add encrypted dummy flags for patient_dimension in crc schema
    ALTER TABLE i2b2demodata_i2b2.patient_dimension ADD COLUMN enc_dummy_flag_cd character(88);
    COMMENT ON COLUMN i2b2demodata_i2b2.patient_dimension.enc_dummy_flag_cd IS 'base64-encoded encrypted dummy flag (0 or 1)';
    INSERT INTO i2b2demodata_i2b2.code_lookup VALUES
        ('patient_dimension', 'enc_dummy_flag_cd', 'CRC_COLUMN_DESCRIPTOR', 'Encrypted Dummy Flag', NULL, NULL, NULL,
        NULL, 'NOW()', NULL, 1);


    -- increase size of modifier_path in the modifier_dimension table
    ALTER TABLE i2b2demodata_i2b2.modifier_dimension ALTER COLUMN modifier_path TYPE varchar(2000);

    -- increase size of concept_path in the concept_dimension table
    ALTER TABLE i2b2demodata_i2b2.concept_dimension ALTER COLUMN concept_path TYPE varchar(2000);

    -- increase size of encounter_num values (too large for type INT)
    ALTER TABLE i2b2demodata_i2b2.visit_dimension ALTER COLUMN encounter_num TYPE bigint;
    ALTER TABLE i2b2demodata_i2b2.observation_fact ALTER COLUMN encounter_num TYPE bigint;
    ALTER TABLE i2b2demodata_i2b2.encounter_mapping ALTER COLUMN encounter_num TYPE bigint;

    -- change tval_char to type text
    ALTER TABLE i2b2demodata_i2b2.observation_fact ALTER COLUMN tval_char TYPE text;

EOSQL
