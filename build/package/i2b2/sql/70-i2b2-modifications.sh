#!/bin/bash
set -Eeuo pipefail
# set up common medco ontology

psql $PSQL_PARAMS -d "$I2B2_DB_NAME" <<-EOSQL

    -- add encrypted dummy flags for patient_dimension in crc schema
    ALTER TABLE i2b2demodata_i2b2_sensitive.patient_dimension ADD COLUMN enc_dummy_flag_cd character(88);
    COMMENT ON COLUMN i2b2demodata_i2b2_sensitive.patient_dimension.enc_dummy_flag_cd IS 'base64-encoded encrypted dummy flag (0 or 1)';
    INSERT INTO i2b2demodata_i2b2_sensitive.code_lookup VALUES
        ('patient_dimension', 'enc_dummy_flag_cd', 'CRC_COLUMN_DESCRIPTOR', 'Encrypted Dummy Flag', NULL, NULL, NULL,
        NULL, 'NOW()', NULL, 1);

EOSQL
