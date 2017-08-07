#!/bin/bash
set -e

# bug fix
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" <<-EOSQL
    update i2b2hive.crc_db_lookup set c_db_fullschema = 'i2b2demodata' where c_domain_id = 'i2b2demo';
    update i2b2hive.im_db_lookup set c_db_fullschema = 'i2b2imdata' where c_domain_id = 'i2b2demo';
    update i2b2hive.ont_db_lookup set c_db_fullschema = 'i2b2metadata' where c_domain_id = 'i2b2demo';
    update i2b2hive.work_db_lookup set c_db_fullschema = 'i2b2workdata' where c_domain_id = 'i2b2demo';

    insert into i2b2pm.pm_cell_params (datatype_cd, cell_id, project_path, param_name_cd, value, changeby_char, status_cd) values
        ('T', 'FRC', '/', 'DestDir', '$I2B2_FR_FILES_DIR', 'i2b2', 'A');
EOSQL

# update domain name and url ports
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" <<-EOSQL
    update i2b2hive.crc_db_lookup SET C_DOMAIN_ID = '$I2B2_DOMAIN_NAME' WHERE C_DOMAIN_ID = 'i2b2demo';
    UPDATE i2b2hive.im_db_lookup SET C_DOMAIN_ID = '$I2B2_DOMAIN_NAME' WHERE C_DOMAIN_ID = 'i2b2demo';
    UPDATE i2b2hive.ont_db_lookup SET C_DOMAIN_ID = '$I2B2_DOMAIN_NAME' WHERE C_DOMAIN_ID = 'i2b2demo';
    UPDATE i2b2hive.work_db_lookup SET C_DOMAIN_ID = '$I2B2_DOMAIN_NAME' WHERE C_DOMAIN_ID = 'i2b2demo';
    UPDATE i2b2pm.pm_hive_data SET DOMAIN_NAME = '$I2B2_DOMAIN_NAME' WHERE DOMAIN_NAME = 'i2b2demo';
    UPDATE i2b2pm.pm_hive_data SET DOMAIN_ID = '$I2B2_DOMAIN_NAME' WHERE DOMAIN_ID = 'i2b2';

    update i2b2pm.pm_cell_data set url = 'http://i2b2-server:8080/i2b2/services/QueryToolService/' where url = 'http://localhost:9090/i2b2/services/QueryToolService/';
    update i2b2pm.pm_cell_data set url = 'http://i2b2-server:8080/i2b2/services/FRService/' where url = 'http://localhost:9090/i2b2/services/FRService/';
    update i2b2pm.pm_cell_data set url = 'http://i2b2-server:8080/i2b2/services/OntologyService/' where url = 'http://localhost:9090/i2b2/services/OntologyService/';
    update i2b2pm.pm_cell_data set url = 'http://i2b2-server:8080/i2b2/services/WorkplaceService/' where url = 'http://localhost:9090/i2b2/services/WorkplaceService/';
    update i2b2pm.pm_cell_data set url = 'http://i2b2-server:8080/i2b2/services/IMService/' where url = 'http://localhost:9090/i2b2/services/IMService/';
EOSQL


