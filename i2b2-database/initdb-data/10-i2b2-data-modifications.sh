#!/bin/bash
set -e

### in scenario of adding medco to an existing i2b2 installation: this is the existing installation

# bug fixes
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" <<-EOSQL
    update i2b2hive.crc_db_lookup set c_db_fullschema = 'i2b2demodata' where c_domain_id = 'i2b2demo';
    update i2b2hive.im_db_lookup set c_db_fullschema = 'i2b2imdata' where c_domain_id = 'i2b2demo';
    update i2b2hive.ont_db_lookup set c_db_fullschema = 'i2b2metadata' where c_domain_id = 'i2b2demo';
    update i2b2hive.work_db_lookup set c_db_fullschema = 'i2b2workdata' where c_domain_id = 'i2b2demo';

    insert into i2b2pm.pm_cell_params (datatype_cd, cell_id, project_path, param_name_cd, value, changeby_char, status_cd) values
        ('T', 'FRC', '/', 'DestDir', '$I2B2_FR_FILES_DIR', 'i2b2', 'A');
EOSQL

# update hive data (DB lookups)
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" <<-EOSQL
    update i2b2hive.crc_db_lookup SET
        C_DOMAIN_ID = '$I2B2_DOMAIN_NAME', C_PROJECT_PATH = '/MedCo/', C_DB_NICENAME = 'MedCo'
        WHERE C_DOMAIN_ID = 'i2b2demo';
    UPDATE i2b2hive.im_db_lookup SET
        C_DOMAIN_ID = '$I2B2_DOMAIN_NAME', C_PROJECT_PATH = '/MedCo/'
        WHERE C_DOMAIN_ID = 'i2b2demo';
    UPDATE i2b2hive.ont_db_lookup SET
        C_DOMAIN_ID = '$I2B2_DOMAIN_NAME', C_PROJECT_PATH = '/MedCo/'
        WHERE C_DOMAIN_ID = 'i2b2demo';
    UPDATE i2b2hive.work_db_lookup SET
        C_DOMAIN_ID = '$I2B2_DOMAIN_NAME', C_PROJECT_PATH = '/MedCo/'
        WHERE C_DOMAIN_ID = 'i2b2demo';
EOSQL


# load i2b2 PM data, medco version
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" <<-EOSQL

    UPDATE i2b2pm.pm_hive_data SET
        DOMAIN_NAME = '$I2B2_DOMAIN_NAME',
        DOMAIN_ID = '$I2B2_DOMAIN_NAME',
        HELPURL = 'https://github.com/lca1/medco'
        WHERE DOMAIN_ID = 'i2b2';
    UPDATE i2b2pm.PM_USER_DATA SET
        USER_ID = 'medcoadmin', FULL_NAME = 'MedCo Admin', PASSWORD = 'f8eb764674b57b5710e3c1665464e29'
        WHERE USER_ID = 'i2b2';
    UPDATE i2b2pm.PM_USER_DATA SET
        PASSWORD = '7cb1ac9deab165535494d60da1d3d7e'
        WHERE USER_ID = 'AGG_SERVICE_ACCOUNT';
    UPDATE i2b2pm.PM_PROJECT_USER_ROLES SET
        USER_ID = 'medcoadmin'
        WHERE USER_ID = 'i2b2';

    INSERT INTO i2b2pm.PM_USER_DATA (USER_ID, FULL_NAME, PASSWORD, STATUS_CD)
        VALUES('medcouser', 'MedCo User', 'f8eb764674b57b5710e3c1665464e29', 'A');
    insert into i2b2pm.pm_project_data (project_id, project_name, project_wiki, project_path, status_cd)
        values ('MedCo', 'MedCo', 'https://github.com/lca1/medco', '/MedCo/', 'A');

    INSERT INTO i2b2pm.PM_CELL_DATA (CELL_ID, PROJECT_PATH, NAME, METHOD_CD, URL, CAN_OVERRIDE, STATUS_CD)
        VALUES('CRC', '/', 'Data Repository', 'REST', 'http://i2b2-server:8080/i2b2/services/QueryToolService/', 1, 'A');
    INSERT INTO i2b2pm.PM_CELL_DATA(CELL_ID, PROJECT_PATH, NAME, METHOD_CD, URL, CAN_OVERRIDE, STATUS_CD)
        VALUES('FRC', '/', 'File Repository ', 'SOAP', 'http://i2b2-server:8080/i2b2/services/FRService/', 1, 'A');
    INSERT INTO i2b2pm.PM_CELL_DATA(CELL_ID, PROJECT_PATH, NAME, METHOD_CD, URL, CAN_OVERRIDE, STATUS_CD)
        VALUES('ONT', '/', 'Ontology Cell', 'REST', 'http://i2b2-server:8080/i2b2/services/OntologyService/', 1, 'A');
    INSERT INTO i2b2pm.PM_CELL_DATA(CELL_ID, PROJECT_PATH, NAME, METHOD_CD, URL, CAN_OVERRIDE, STATUS_CD)
        VALUES('WORK', '/', 'Workplace Cell', 'REST', 'http://i2b2-server:8080/i2b2/services/WorkplaceService/', 1, 'A');
    INSERT INTO i2b2pm.PM_CELL_DATA(CELL_ID, PROJECT_PATH, NAME, METHOD_CD, URL, CAN_OVERRIDE, STATUS_CD)
        VALUES('IM', '/', 'IM Cell', 'REST', 'http://i2b2-server:8080/i2b2/services/IMService/', 1, 'A');

    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD)
        VALUES('MedCo', 'AGG_SERVICE_ACCOUNT', 'USER', 'A');
    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD)
        VALUES('MedCo', 'AGG_SERVICE_ACCOUNT', 'MANAGER', 'A');
    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD)
        VALUES('MedCo', 'AGG_SERVICE_ACCOUNT', 'DATA_OBFSC', 'A');
    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD)
        VALUES('MedCo', 'AGG_SERVICE_ACCOUNT', 'DATA_AGG', 'A');
    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD)
        VALUES('MedCo', 'medcoadmin', 'MANAGER', 'A');
    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD)
        VALUES('MedCo', 'medcoadmin', 'USER', 'A');
    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD)
        VALUES('MedCo', 'medcoadmin', 'DATA_OBFSC', 'A');
    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD)
        VALUES('MedCo', 'medcouser', 'USER', 'A');
    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD)
        VALUES('MedCo', 'medcouser', 'DATA_DEID', 'A');
    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD)
        VALUES('MedCo', 'medcouser', 'DATA_OBFSC', 'A');
    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD)
        VALUES('MedCo', 'medcouser', 'DATA_AGG', 'A');
    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD)
        VALUES('MedCo', 'medcouser', 'DATA_LDS', 'A');
    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD)
        VALUES('MedCo', 'medcouser', 'EDITOR', 'A');
    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD)
        VALUES('MedCo', 'medcouser', 'DATA_PROT', 'A');
EOSQL

