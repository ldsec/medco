#!/bin/bash
set -e

### in scenario of adding medco to an existing i2b2 installation: this is the additional shrine stuff

# db lookups
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" <<-EOSQL
    insert into i2b2hive.ont_db_lookup (c_domain_id, c_project_path, c_owner_id, c_db_fullschema, c_db_datasource, c_db_servertype, c_db_nicename)
    values ('$I2B2_DOMAIN_NAME', 'MedCo-SHRINE/', '@', 'shrine_ont', 'java:/OntologyShrineDS', 'POSTGRESQL', 'MedCo-SHRINE')
    on conflict do nothing;

    insert into i2b2hive.crc_db_lookup (c_domain_id, c_project_path, c_owner_id, c_db_fullschema, c_db_datasource, c_db_servertype, c_db_nicename)
    values ('$I2B2_DOMAIN_NAME', '/MedCo-SHRINE/', '@', 'i2b2demodata', 'java:/QueryToolDemoDS', 'POSTGRESQL', 'MedCo-SHRINE')
    on conflict do nothing;

EOSQL

####################################################################################
######################### information about password hash ##########################
####################################################################################

### how to generate the hash (from shrine sources folder)
# cd "install/i2b2-1.7/i2b2"
# javac ./I2b2PasswordCryptor.java
# SHRINE_PW=$(java -classpath ./ I2b2PasswordCryptor <thepassword>)

### some encrypted versions:
# demouser=             9117d59a69dc49807671a51f10ab7f
# prigen2017=           f8eb764674b57b5710e3c1665464e29
# pFjy3EjDVwLfT2rB9xkK= 7cb1ac9deab165535494d60da1d3d7e
####################################################################################

# pm data
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" <<-EOSQL



    insert into i2b2pm.pm_cell_data (cell_id, project_path, name, method_cd, url, can_override, status_cd)
    values ('CRC', '/MedCo-SHRINE', 'MedCo-SHRINE Federated Query', 'REST', 'https://shrine-server:6443/shrine/rest/i2b2/', 1, 'A');

    insert into i2b2pm.pm_project_data (project_id, project_name, project_wiki, project_path, status_cd)
    values ('MedCo-SHRINE', 'MedCo-SHRINE', 'https://github.com/lca1/medco', '/MedCo-SHRINE', 'A');

    INSERT INTO i2b2pm.PM_USER_DATA (USER_ID, FULL_NAME, PASSWORD, STATUS_CD)
        VALUES('medcoshrineuser', 'MedCo SHRINE User', 'f8eb764674b57b5710e3c1665464e29', 'A');
    insert into i2b2pm.pm_project_user_roles (project_id, user_id, user_role_cd, status_cd)
    values ('MedCo-SHRINE', 'medcoshrineuser', 'USER', 'A');
    insert into i2b2pm.pm_project_user_roles (project_id, user_id, user_role_cd, status_cd)
    values ('MedCo-SHRINE', 'medcoshrineuser', 'DATA_OBFSC', 'A');


    INSERT INTO i2b2pm.pm_user_params(datatype_cd, user_id, param_name_cd, value, change_date, entry_date, status_cd) VALUES
    ('T', 'AGG_SERVICE_ACCOUNT', 'qep', 'true', 'NOW()', 'NOW()', 'A');
    INSERT INTO i2b2pm.pm_user_params(datatype_cd, user_id, param_name_cd, value, change_date, entry_date, status_cd) VALUES
    ('T', 'medcoadmin', 'DataSteward', 'true', 'NOW()', 'NOW()', 'A');
EOSQL

# add demo shrine ontology structure
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" <<-EOSQL

    CREATE TABLE shrine_ont.SHRINE
    (
        C_HLEVEL NUMERIC(22,0),
        C_FULLNAME VARCHAR(900),
        C_NAME VARCHAR(2000),
        C_SYNONYM_CD CHAR(1),
        C_VISUALATTRIBUTES CHAR(3),
        C_TOTALNUM NUMERIC(22,0),
        C_BASECODE VARCHAR(450),
        C_METADATAXML TEXT,
        C_FACTTABLECOLUMN VARCHAR(50),
        C_TABLENAME VARCHAR(50),
        C_COLUMNNAME VARCHAR(50),
        C_COLUMNDATATYPE VARCHAR(50),
        C_OPERATOR VARCHAR(10),
        C_DIMCODE VARCHAR(900),
        C_COMMENT TEXT,
        C_TOOLTIP VARCHAR(900),
        UPDATE_DATE DATE,
        DOWNLOAD_DATE DATE,
        IMPORT_DATE DATE,
        SOURCESYSTEM_CD VARCHAR(50),
        VALUETYPE_CD VARCHAR(50),
        M_APPLIED_PATH VARCHAR(900),
        M_EXCLUSION_CD VARCHAR(900)
    );

    grant all privileges on all tables in schema shrine_ont to shrine_ont;
    grant all privileges on all sequences in schema shrine_ont to shrine_ont;
    grant all privileges on all functions in schema shrine_ont to shrine_ont;
    grant all privileges on all tables in schema shrine_ont to i2b2metadata;
    grant all privileges on all sequences in schema shrine_ont to i2b2metadata;
    grant all privileges on all functions in schema shrine_ont to i2b2metadata;
EOSQL

# add encrypted dummy flags for patient_dimension in crc schema
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" <<-EOSQL
    ALTER TABLE i2b2demodata.patient_dimension ADD COLUMN enc_dummy_flag_cd character(88);
    COMMENT ON COLUMN i2b2demodata.patient_dimension.enc_dummy_flag_cd IS 'base64-encoded encrypted dummy flag (0 or 1)';
    INSERT INTO i2b2demodata.code_lookup VALUES
        ('patient_dimension', 'enc_dummy_flag_cd', 'CRC_COLUMN_DESCRIPTOR', 'Encrypted Dummy Flag', NULL, NULL, NULL,
        NULL, 'NOW()', NULL, 1);
EOSQL


# add i2b2 test query term for shrine (ontology + database)
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" <<-EOSQL
    INSERT INTO i2b2metadata.sensitive_tagged VALUES
        (2, '\medco\tagged\TESTKEY\', '', 'N', 'LH ', NULL, 'TAG_ID:TESTKEY', NULL, 'concept_cd', 'concept_dimension',
        'concept_path', 'T', 'LIKE', '\medco\tagged\TESTKEY\', NULL, NULL, 'NOW()', NULL, NULL, NULL, 'TAG_ID', '@',
        NULL, NULL, NULL, NULL);
    INSERT INTO i2b2demodata.concept_dimension VALUES
        ('\medco\tagged\TESTKEY\', 'TAG_ID:TESTKEY', NULL, NULL, NULL, NULL, 'NOW()', NULL, NULL);
    INSERT INTO i2b2demodata.patient_mapping VALUES
        ('TESTPATIENT', 'TESTSITE', -1, NULL, 'MedCo', NULL, NULL, NULL, 'NOW()', NULL, 1);
    INSERT INTO i2b2demodata.patient_mapping VALUES
        ('-1', 'HIVE', -1, 'A', 'HIVE', NULL, 'NOW()', 'NOW()', 'NOW()', 'edu.harvard.i2b2.crc', 1);
    INSERT INTO i2b2demodata.patient_dimension VALUES
        (-1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'NOW()', NULL, 1,
        'FzXxSbBn86gMmF7WT6a4kHDcHrOg3SEkaojcPm7U3qsQp0bhzaLZLYenL/+yNS5j39TFcLU1uSUE5I8tD3Qryw==');
    INSERT INTO i2b2demodata.encounter_mapping VALUES
        ('TESTVISIT', 'TESTSITE', 'MedCo', -1, 'TESTPATIENT', 'TESTSITE', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1);
    INSERT INTO i2b2demodata.encounter_mapping VALUES
        ('-1', 'HIVE', 'HIVE', -1, 'TESTPATIENT', 'TESTSITE', 'A', NULL, 'NOW()', 'NOW()', 'NOW()', 'edu.harvard.i2b2.crc', 1);
    INSERT INTO i2b2demodata.visit_dimension VALUES
        (-1, -1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'NOW()', 'TESTSITE', 1);
    INSERT INTO i2b2demodata.provider_dimension VALUES
        ('TESTSITE', '\medco\institutions\TESTSITE\', 'TESTSITE', NULL, NULL, NULL, 'NOW()', NULL, 1);
    INSERT INTO i2b2demodata.observation_fact VALUES
        (-1, -1, 'TAG_ID:TESTKEY', 'TESTSITE', 'NOW()', '@', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'TESTSITE',
        NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, 1);
EOSQL

# original shrine ontology: only the version key
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" <<-EOSQL
    INSERT into shrine_ont.TABLE_ACCESS
        ( C_TABLE_CD, C_TABLE_NAME, C_PROTECTED_ACCESS, C_HLEVEL, C_NAME, C_FULLNAME, C_SYNONYM_CD, C_VISUALATTRIBUTES,
        C_TOOLTIP, C_FACTTABLECOLUMN, C_DIMTABLENAME, C_COLUMNNAME, C_COLUMNDATATYPE, C_DIMCODE, C_OPERATOR) values
        ( 'SHRINE', 'SHRINE', 'N', 0, 'SHRINE Ontology', '\SHRINE\', 'N', 'CH', 'SHRINE Ontology', 'concept_cd',
        'concept_dimension', 'concept_path', 'T', '\SHRINE\', 'LIKE')
        on conflict do nothing;
    INSERT INTO shrine_ont.SHRINE (C_HLEVEL, C_FULLNAME, C_NAME, C_SYNONYM_CD, C_VISUALATTRIBUTES, C_TOTALNUM, C_BASECODE, C_METADATAXML,
        C_FACTTABLECOLUMN, C_TABLENAME, C_COLUMNNAME, C_COLUMNDATATYPE, C_OPERATOR, C_DIMCODE, C_COMMENT, C_TOOLTIP,
        UPDATE_DATE, DOWNLOAD_DATE, IMPORT_DATE, SOURCESYSTEM_CD, VALUETYPE_CD, M_APPLIED_PATH, M_EXCLUSION_CD ) VALUES
        (1, '\SHRINE\ONTOLOGYVERSION\', 'ONTOLOGYVERSION', 'N', 'FH', NULL, NULL, '', 'concept_cd', 'concept_dimension',
        'concept_path', 'T', 'LIKE', '\SHRINE\ONTOLOGYVERSION\', '', 'ONTOLOGYVERSION\', NULL, NULL, NULL, 'SHRINE', NULL,
        '@', NULL );
    INSERT INTO shrine_ont.SHRINE (C_HLEVEL, C_FULLNAME, C_NAME, C_SYNONYM_CD, C_VISUALATTRIBUTES, C_TOTALNUM, C_BASECODE, C_METADATAXML,
        C_FACTTABLECOLUMN, C_TABLENAME, C_COLUMNNAME, C_COLUMNDATATYPE, C_OPERATOR, C_DIMCODE, C_COMMENT, C_TOOLTIP,
        UPDATE_DATE, DOWNLOAD_DATE, IMPORT_DATE, SOURCESYSTEM_CD, VALUETYPE_CD, M_APPLIED_PATH, M_EXCLUSION_CD ) VALUES
        (2, '\SHRINE\ONTOLOGYVERSION\MedCo-SHRINE_Ontology_Empty\', 'MedCo-SHRINE_Ontology_Empty', 'N', 'LH', NULL, NULL,
        '', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\SHRINE\ONTOLOGYVERSION\MedCo-SHRINE_Ontology_Empty\',
        '', 'ONTOLOGYVERSION\MedCo-SHRINE_Ontology_Empty\', NULL, NULL, NULL, 'SHRINE', NULL, '@', NULL );
EOSQL
# full loading of shrine ontology disabled
#wget https://open.med.harvard.edu/svn/shrine-ontology/SHRINE_Demo_Downloads/trunk/ShrineDemo.sql
#sed -i '1s/^/SET search_path TO shrine_ont;\n/' Shrine.sql
#psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" < ShrineDemo.sql
