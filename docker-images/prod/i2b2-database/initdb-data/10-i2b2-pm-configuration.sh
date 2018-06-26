#!/bin/bash
set -e
# configuration of the i2b2 project management cell, that is common to the demo and the medco projects

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

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DEMO_DB_NAME" <<-EOSQL

    -- cell parameters
    insert into i2b2pm.pm_cell_params (datatype_cd, cell_id, project_path, param_name_cd, value, changeby_char, status_cd) values
        ('T', 'FRC', '/', 'DestDir', '$I2B2_FR_FILES_DIR', 'i2b2', 'A');


    -- database lookups
    update i2b2hive.crc_db_lookup SET C_DOMAIN_ID = '$I2B2_DOMAIN_NAME' WHERE C_DOMAIN_ID = 'i2b2demo';
    UPDATE i2b2hive.im_db_lookup SET C_DOMAIN_ID = '$I2B2_DOMAIN_NAME' WHERE C_DOMAIN_ID = 'i2b2demo';
    UPDATE i2b2hive.ont_db_lookup SET C_DOMAIN_ID = '$I2B2_DOMAIN_NAME' WHERE C_DOMAIN_ID = 'i2b2demo';
    UPDATE i2b2hive.work_db_lookup SET C_DOMAIN_ID = '$I2B2_DOMAIN_NAME' WHERE C_DOMAIN_ID = 'i2b2demo';

    INSERT INTO i2b2hive.CRC_DB_LOOKUP(c_domain_id, c_project_path, c_owner_id, c_db_fullschema, c_db_datasource, c_db_servertype,
        c_db_nicename, c_db_tooltip, c_comment, c_entry_date, c_change_date, c_status_cd)
        VALUES('$I2B2_DOMAIN_NAME', '/MedCo/', '@', 'i2b2demodata', 'java:/QueryToolMedCoDS', 'POSTGRESQL', 'MedCo', NULL, NULL, NULL, NULL, NULL);
    INSERT INTO i2b2hive.IM_DB_LOOKUP(c_domain_id, c_project_path, c_owner_id, c_db_fullschema, c_db_datasource, c_db_servertype,
        c_db_nicename, c_db_tooltip, c_comment, c_entry_date, c_change_date, c_status_cd)
        VALUES('$I2B2_DOMAIN_NAME', 'MedCo/', '@', 'i2b2imdata', 'java:/IMMedCoDS', 'POSTGRESQL', 'IM', NULL, NULL, NULL, NULL, NULL);
    INSERT INTO i2b2hive.ONT_DB_LOOKUP(c_domain_id, c_project_path, c_owner_id, c_db_fullschema, c_db_datasource, c_db_servertype,
        c_db_nicename, c_db_tooltip, c_comment, c_entry_date, c_change_date, c_status_cd)
        VALUES('$I2B2_DOMAIN_NAME', 'MedCo/', '@', 'i2b2metadata', 'java:/OntologyMedCoDS', 'POSTGRESQL', 'Metadata', NULL, NULL, NULL, NULL, NULL);
    INSERT INTO i2b2hive.WORK_DB_LOOKUP(c_domain_id, c_project_path, c_owner_id, c_db_fullschema, c_db_datasource, c_db_servertype,
        c_db_nicename, c_db_tooltip, c_comment, c_entry_date, c_change_date, c_status_cd)
        VALUES('$I2B2_DOMAIN_NAME', 'MedCo/', '@', 'i2b2workdata', 'java:/WorkplaceMedCoDS', 'POSTGRESQL', 'Workplace', NULL, NULL, NULL, NULL, NULL);


    -- hive & users data
    UPDATE i2b2pm.pm_hive_data SET DOMAIN_ID = '$I2B2_DOMAIN_NAME', DOMAIN_NAME = '$I2B2_DOMAIN_NAME',
        HELPURL = 'https://github.com/lca1/medco' WHERE DOMAIN_ID = 'i2b2';
    insert into i2b2pm.pm_project_data (project_id, project_name, project_wiki, project_path, status_cd) values
        ('MedCo', 'MedCo', 'https://github.com/lca1/medco', '/MedCo', 'A');

    UPDATE i2b2pm.PM_CELL_DATA SET URL = 'http://i2b2-server:8080/i2b2/services/QueryToolService/' WHERE CELL_ID = 'CRC';
    UPDATE i2b2pm.PM_CELL_DATA SET URL = 'http://i2b2-server:8080/i2b2/services/FRService/' WHERE CELL_ID = 'FRC';
    UPDATE i2b2pm.PM_CELL_DATA SET URL = 'http://i2b2-server:8080/i2b2/services/OntologyService/' WHERE CELL_ID = 'ONT';
    UPDATE i2b2pm.PM_CELL_DATA SET URL = 'http://i2b2-server:8080/i2b2/services/WorkplaceService/' WHERE CELL_ID = 'WORK';
    UPDATE i2b2pm.PM_CELL_DATA SET URL = 'http://i2b2-server:8080/i2b2/services/IMService/' WHERE CELL_ID = 'IM';

    UPDATE i2b2pm.PM_USER_DATA SET PASSWORD = 'f8eb764674b57b5710e3c1665464e29' WHERE USER_ID = 'i2b2';
    UPDATE i2b2pm.PM_USER_DATA SET PASSWORD = 'f8eb764674b57b5710e3c1665464e29' WHERE USER_ID = 'demo';
    UPDATE i2b2pm.PM_USER_DATA SET PASSWORD = '7cb1ac9deab165535494d60da1d3d7e' WHERE USER_ID = 'AGG_SERVICE_ACCOUNT';
    INSERT INTO i2b2pm.PM_USER_DATA (USER_ID, FULL_NAME, PASSWORD, STATUS_CD) VALUES('medcoi2b2user', 'MedCo I2b2 User', 'f8eb764674b57b5710e3c1665464e29', 'A');
    INSERT INTO i2b2pm.PM_USER_DATA (USER_ID, FULL_NAME, PASSWORD, STATUS_CD) VALUES('medcoadmin', 'MedCo Admin', 'f8eb764674b57b5710e3c1665464e29', 'A');

    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD) VALUES
        ('MedCo', 'AGG_SERVICE_ACCOUNT', 'USER', 'A'),
        ('MedCo', 'AGG_SERVICE_ACCOUNT', 'MANAGER', 'A'),
        ('MedCo', 'AGG_SERVICE_ACCOUNT', 'DATA_OBFSC', 'A'),
        ('MedCo', 'AGG_SERVICE_ACCOUNT', 'DATA_AGG', 'A'),
        ('MedCo', 'medcoadmin', 'MANAGER', 'A'),
        ('MedCo', 'medcoadmin', 'USER', 'A'),
        ('MedCo', 'medcoadmin', 'DATA_OBFSC', 'A'),
        ('MedCo', 'medcoi2b2user', 'USER', 'A'),
        ('MedCo', 'medcoi2b2user', 'DATA_DEID', 'A'),
        ('MedCo', 'medcoi2b2user', 'DATA_OBFSC', 'A'),
        ('MedCo', 'medcoi2b2user', 'DATA_AGG', 'A'),
        ('MedCo', 'medcoi2b2user', 'DATA_LDS', 'A'),
        ('MedCo', 'medcoi2b2user', 'EDITOR', 'A'),
        ('MedCo', 'medcoi2b2user', 'DATA_PROT', 'A');
EOSQL
