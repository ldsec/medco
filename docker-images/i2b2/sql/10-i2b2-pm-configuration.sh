#!/bin/bash
set -Eeuo pipefail
# configuration of the i2b2 project management cell

# generate password hashes
I2B2_SERVICE_PASSWORD_HASH=$(java -classpath "$JBOSS_HOME/I2b2PasswordHash/" I2b2PasswordHash "$I2B2_SERVICE_PASSWORD")
DEFAULT_USER_PASSWORD_HASH=$(java -classpath "$JBOSS_HOME/I2b2PasswordHash/" I2b2PasswordHash "$DEFAULT_USER_PASSWORD")

psql $PSQL_PARAMS -d "$I2B2_DB_NAME" <<-EOSQL
    -- todo: cell data http://i2b2:8080/i2b2/services as param, and also keycloak

    -- cell parameters
    insert into i2b2pm.pm_cell_params (datatype_cd, cell_id, project_path, param_name_cd, value, changeby_char, status_cd) values
        ('T', 'FRC', '/', 'DestDir', '$I2B2_FR_FILES_DIR', 'i2b2', 'A');

    -- database lookups
    update i2b2hive.crc_db_lookup SET C_DOMAIN_ID = '$I2B2_DOMAIN_NAME' WHERE C_DOMAIN_ID = 'i2b2demo';
    UPDATE i2b2hive.im_db_lookup SET C_DOMAIN_ID = '$I2B2_DOMAIN_NAME' WHERE C_DOMAIN_ID = 'i2b2demo';
    UPDATE i2b2hive.ont_db_lookup SET C_DOMAIN_ID = '$I2B2_DOMAIN_NAME' WHERE C_DOMAIN_ID = 'i2b2demo';
    UPDATE i2b2hive.work_db_lookup SET C_DOMAIN_ID = '$I2B2_DOMAIN_NAME' WHERE C_DOMAIN_ID = 'i2b2demo';

    -- hive & users data
    UPDATE i2b2pm.pm_hive_data SET DOMAIN_ID = '$I2B2_DOMAIN_NAME', DOMAIN_NAME = '$I2B2_DOMAIN_NAME' WHERE DOMAIN_ID = 'i2b2';

    UPDATE i2b2pm.PM_CELL_DATA SET URL = 'http://i2b2:8080/i2b2/services/QueryToolService/' WHERE CELL_ID = 'CRC';
    UPDATE i2b2pm.PM_CELL_DATA SET URL = 'http://i2b2:8080/i2b2/services/FRService/' WHERE CELL_ID = 'FRC';
    UPDATE i2b2pm.PM_CELL_DATA SET URL = 'http://i2b2:8080/i2b2/services/OntologyService/' WHERE CELL_ID = 'ONT';
    UPDATE i2b2pm.PM_CELL_DATA SET URL = 'http://i2b2:8080/i2b2/services/WorkplaceService/' WHERE CELL_ID = 'WORK';
    UPDATE i2b2pm.PM_CELL_DATA SET URL = 'http://i2b2:8080/i2b2/services/IMService/' WHERE CELL_ID = 'IM';

    UPDATE i2b2pm.PM_USER_DATA SET PASSWORD = '$DEFAULT_USER_PASSWORD_HASH' WHERE USER_ID = 'i2b2';
    UPDATE i2b2pm.PM_USER_DATA SET PASSWORD = '$DEFAULT_USER_PASSWORD_HASH' WHERE USER_ID = 'demo';
    UPDATE i2b2pm.PM_USER_DATA SET PASSWORD = '$I2B2_SERVICE_PASSWORD_HASH' WHERE USER_ID = 'AGG_SERVICE_ACCOUNT';

    -- oidc user
    -- todo: separate sql file to add oidc user(s) or cell or project
    INSERT INTO i2b2pm.pm_user_data VALUES ('test', 'test', NULL, 'test@test.com', NULL, NOW(), NOW(), 'i2b2', 'A');
    INSERT INTO i2b2pm.pm_user_params VALUES (1, 'T', 'test', 'authentication_method', 'OIDC', NOW(), NOW(), 'i2b2', 'A');
    INSERT INTO i2b2pm.pm_user_params VALUES (2, 'T', 'test', 'oidc_jwks_uri', 'http://keycloak:8080/auth/realms/master/protocol/openid-connect/certs', NOW(), NOW(), 'i2b2', 'A');
    INSERT INTO i2b2pm.pm_user_params VALUES (3, 'T', 'test', 'oidc_client_id', 'i2b2-local-jwt', NOW(), NOW(), 'i2b2', 'A');
    INSERT INTO i2b2pm.pm_user_params VALUES (4, 'T', 'test', 'oidc_user_field', 'preferred_username', NOW(), NOW(), 'i2b2', 'A');
    INSERT INTO i2b2pm.pm_user_params VALUES (5, 'T', 'test', 'oidc_token_issuer', 'http://keycloak:8080/auth/realms/master', NOW(), NOW(), 'i2b2', 'A');

    INSERT INTO i2b2pm.pm_project_user_roles VALUES ('Demo', 'test', 'MANAGER', NOW(), NOW(), 'i2b2', 'A');
    INSERT INTO i2b2pm.pm_project_user_roles VALUES ('Demo', 'test', 'USER', NOW(), NOW(), 'i2b2', 'A');
    INSERT INTO i2b2pm.pm_project_user_roles VALUES ('Demo', 'test', 'DATA_OBFSC', NOW(), NOW(), 'i2b2', 'A');
EOSQL
