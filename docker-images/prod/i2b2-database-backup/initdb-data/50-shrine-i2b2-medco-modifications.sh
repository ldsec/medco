#!/bin/bash
set -e
# modifications in the i2b2 database for shrine/medco

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_MEDCO_DB_NAME" <<-EOSQL

    -- add encrypted dummy flags for patient_dimension in crc schema
    ALTER TABLE i2b2demodata.patient_dimension ADD COLUMN enc_dummy_flag_cd character(88);
    COMMENT ON COLUMN i2b2demodata.patient_dimension.enc_dummy_flag_cd IS 'base64-encoded encrypted dummy flag (0 or 1)';
    INSERT INTO i2b2demodata.code_lookup VALUES
        ('patient_dimension', 'enc_dummy_flag_cd', 'CRC_COLUMN_DESCRIPTOR', 'Encrypted Dummy Flag', NULL, NULL, NULL,
        NULL, 'NOW()', NULL, 1);

    -- add i2b2 test query term for shrine (ontology + database)
    INSERT INTO i2b2metadata.sensitive_tagged VALUES
        (2, '\medco\tagged\TESTKEY\', '', 'N', 'LH ', NULL, 'TAG_ID:TESTKEY', NULL, 'concept_cd', 'concept_dimension',
        'concept_path', 'T', 'LIKE', '\medco\tagged\TESTKEY\', NULL, NULL, 'NOW()', NULL, NULL, NULL, 'TAG_ID', '@',
        NULL, NULL, NULL, NULL);
    INSERT INTO i2b2demodata.concept_dimension VALUES
        ('\medco\tagged\TESTKEY\', 'TAG_ID:TESTKEY', NULL, NULL, NULL, NULL, 'NOW()', NULL, -1);
    INSERT INTO i2b2demodata.patient_mapping VALUES
        ('TESTPATIENT', 'TESTSITE', -1, NULL, 'MedCo', NULL, NULL, NULL, 'NOW()', NULL, -1);
    INSERT INTO i2b2demodata.patient_mapping VALUES
        ('-1', 'HIVE', -1, 'A', 'HIVE', NULL, 'NOW()', 'NOW()', 'NOW()', 'edu.harvard.i2b2.crc', -1);
    INSERT INTO i2b2demodata.patient_dimension VALUES
        (-1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'NOW()', NULL, -1,
        'FzXxSbBn86gMmF7WT6a4kHDcHrOg3SEkaojcPm7U3qsQp0bhzaLZLYenL/+yNS5j39TFcLU1uSUE5I8tD3Qryw==');
    INSERT INTO i2b2demodata.encounter_mapping VALUES
        ('TESTVISIT', 'TESTSITE', 'MedCo', -1, 'TESTPATIENT', 'TESTSITE', NULL, NULL, NULL, NULL, 'NOW()', NULL, -1);
    INSERT INTO i2b2demodata.encounter_mapping VALUES
        ('-1', 'HIVE', 'HIVE', -1, 'TESTPATIENT', 'TESTSITE', 'A', NULL, 'NOW()', 'NOW()', 'NOW()', 'edu.harvard.i2b2.crc', -1);
    INSERT INTO i2b2demodata.visit_dimension VALUES
        (-1, -1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'NOW()', 'TESTSITE', -1);
    INSERT INTO i2b2demodata.provider_dimension VALUES
        ('TESTSITE', '\medco\institutions\TESTSITE\', 'TESTSITE', NULL, NULL, NULL, 'NOW()', NULL, -1);
    INSERT INTO i2b2demodata.observation_fact VALUES
        (-1, -1, 'TAG_ID:TESTKEY', 'TESTSITE', 'NOW()', '@', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'TESTSITE',
        NULL, NULL, NULL, NULL, 'NOW()', NULL, -1, -1);

    -- insert custom shrine version entry
    --INSERT INTO shrine_ont.SHRINE (C_HLEVEL, C_FULLNAME, C_NAME, C_SYNONYM_CD, C_VISUALATTRIBUTES, C_TOTALNUM, C_BASECODE, C_METADATAXML,
    --    C_FACTTABLECOLUMN, C_TABLENAME, C_COLUMNNAME, C_COLUMNDATATYPE, C_OPERATOR, C_DIMCODE, C_COMMENT, C_TOOLTIP,
    --    UPDATE_DATE, DOWNLOAD_DATE, IMPORT_DATE, SOURCESYSTEM_CD, VALUETYPE_CD, M_APPLIED_PATH, M_EXCLUSION_CD ) VALUES
    --    (1, '\SHRINE\ONTOLOGYVERSION\', 'ONTOLOGYVERSION', 'N', 'FH', NULL, NULL, '', 'concept_cd', 'concept_dimension',
    --    'concept_path', 'T', 'LIKE', '\SHRINE\ONTOLOGYVERSION\', '', 'ONTOLOGYVERSION\', NULL, NULL, NULL, 'SHRINE', NULL,
    --    '@', NULL );
    --INSERT INTO shrine_ont.SHRINE (C_HLEVEL, C_FULLNAME, C_NAME, C_SYNONYM_CD, C_VISUALATTRIBUTES, C_TOTALNUM, C_BASECODE, C_METADATAXML,
    --    C_FACTTABLECOLUMN, C_TABLENAME, C_COLUMNNAME, C_COLUMNDATATYPE, C_OPERATOR, C_DIMCODE, C_COMMENT, C_TOOLTIP,
    --    UPDATE_DATE, DOWNLOAD_DATE, IMPORT_DATE, SOURCESYSTEM_CD, VALUETYPE_CD, M_APPLIED_PATH, M_EXCLUSION_CD ) VALUES
    --    (2, '\SHRINE\ONTOLOGYVERSION\MedCo-SHRINE_Ontology_Empty\', 'MedCo-SHRINE_Ontology_Empty', 'N', 'LH', NULL, NULL,
    --    '', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\SHRINE\ONTOLOGYVERSION\MedCo-SHRINE_Ontology_Empty\',
    --    '', 'ONTOLOGYVERSION\MedCo-SHRINE_Ontology_Empty\', NULL, NULL, NULL, 'SHRINE', NULL, '@', NULL );
EOSQL
