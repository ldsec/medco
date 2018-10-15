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

    -- add i2b2 test query term for shrine (ontology + database)
    --INSERT INTO i2b2metadata_i2b2.sensitive_tagged VALUES
    --    (2, '\medco\tagged\TESTKEY\', '', 'N', 'LH ', NULL, 'TAG_ID:TESTKEY', NULL, 'concept_cd', 'concept_dimension',
    --    'concept_path', 'T', 'LIKE', '\medco\tagged\TESTKEY\', NULL, NULL, 'NOW()', NULL, NULL, NULL, 'TAG_ID', '@',
    --    NULL, NULL, NULL, NULL);
    --INSERT INTO i2b2demodata_i2b2.concept_dimension VALUES
    --    ('\medco\tagged\TESTKEY\', 'TAG_ID:TESTKEY', NULL, NULL, NULL, NULL, 'NOW()', NULL, -1);
    --INSERT INTO i2b2demodata_i2b2.patient_mapping VALUES
    --    ('TESTPATIENT', 'TESTSITE', -1, NULL, 'MedCo', NULL, NULL, NULL, 'NOW()', NULL, -1);
    --INSERT INTO i2b2demodata_i2b2.patient_mapping VALUES
    --    ('-1', 'HIVE', -1, 'A', 'HIVE', NULL, 'NOW()', 'NOW()', 'NOW()', 'edu.harvard.i2b2.crc', -1);
    --INSERT INTO i2b2demodata_i2b2.patient_dimension VALUES
    --    (-1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'NOW()', NULL, -1,
    --    'FzXxSbBn86gMmF7WT6a4kHDcHrOg3SEkaojcPm7U3qsQp0bhzaLZLYenL/+yNS5j39TFcLU1uSUE5I8tD3Qryw==');
    --INSERT INTO i2b2demodata_i2b2.encounter_mapping VALUES
    --    ('TESTVISIT', 'TESTSITE', 'MedCo', -1, 'TESTPATIENT', 'TESTSITE', NULL, NULL, NULL, NULL, 'NOW()', NULL, -1);
    --INSERT INTO i2b2demodata_i2b2.encounter_mapping VALUES
    --    ('-1', 'HIVE', 'HIVE', -1, 'TESTPATIENT', 'TESTSITE', 'A', NULL, 'NOW()', 'NOW()', 'NOW()', 'edu.harvard.i2b2.crc', -1);
    --INSERT INTO i2b2demodata_i2b2.visit_dimension VALUES
    --    (-1, -1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'NOW()', 'TESTSITE', -1);
    --INSERT INTO i2b2demodata_i2b2.provider_dimension VALUES
    --    ('TESTSITE', '\medco\institutions\TESTSITE\', 'TESTSITE', NULL, NULL, NULL, 'NOW()', NULL, -1);
    --INSERT INTO i2b2demodata_i2b2.observation_fact VALUES
    --    (-1, -1, 'TAG_ID:TESTKEY', 'TESTSITE', 'NOW()', '@', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'TESTSITE',
    --    NULL, NULL, NULL, NULL, 'NOW()', NULL, -1, -1);

EOSQL
