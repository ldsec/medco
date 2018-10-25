#!/bin/bash
set -Eeuo pipefail
# set up PM for medco

# generate password hashes
I2B2_SERVICE_PASSWORD_HASH=$(java -classpath "$JBOSS_HOME/I2b2PasswordHash/" I2b2PasswordHash "$I2B2_SERVICE_PASSWORD")
DEFAULT_USER_PASSWORD_HASH=$(java -classpath "$JBOSS_HOME/I2b2PasswordHash/" I2b2PasswordHash "$DEFAULT_USER_PASSWORD")

psql $PSQL_PARAMS -d "$I2B2_DB_NAME" <<-EOSQL
    -- database lookups
    insert into i2b2hive.ont_db_lookup (c_domain_id, c_project_path, c_owner_id, c_db_fullschema, c_db_datasource, c_db_servertype, c_db_nicename)
        values ('$I2B2_DOMAIN_NAME', 'MedCo/', '@', 'medco_ont', 'java:/OntologyMedCoDS', 'POSTGRESQL', 'MedCo');
    insert into i2b2hive.crc_db_lookup (c_domain_id, c_project_path, c_owner_id, c_db_fullschema, c_db_datasource, c_db_servertype, c_db_nicename)
        values ('$I2B2_DOMAIN_NAME', '/MedCo/', '@', 'i2b2demodata_i2b2', 'java:/QueryToolDemoDS', 'POSTGRESQL', 'MedCo');

    -- hive & users data
    insert into i2b2pm.pm_project_data (project_id, project_name, project_wiki, project_path, status_cd) values
        ('MedCo', 'MedCo', 'https://github.com/lca1/medco', '/MedCo', 'A');
    --insert into i2b2pm.pm_cell_data (cell_id, project_path, name, method_cd, url, can_override, status_cd)
    --    values ('CRC', '/MedCo', 'MedCo Federated Query', 'REST', 'https://tbd:6443/shrine/rest/i2b2/', 1, 'A');

    INSERT INTO i2b2pm.PM_USER_DATA (USER_ID, FULL_NAME, PASSWORD, STATUS_CD)
        VALUES('medcoadmin', 'MedCo Admin', '$DEFAULT_USER_PASSWORD_HASH', 'A');
    INSERT INTO i2b2pm.PM_USER_DATA (USER_ID, FULL_NAME, PASSWORD, STATUS_CD)
        VALUES('medcouser', 'MedCo User', '$DEFAULT_USER_PASSWORD_HASH', 'A');
    INSERT INTO i2b2pm.PM_USER_DATA (USER_ID, FULL_NAME, PASSWORD, STATUS_CD)
        VALUES('medcoservice', 'MedCo Service User', '$I2B2_SERVICE_PASSWORD_HASH', 'A');

    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD) VALUES
        ('MedCo', 'medcoservice', 'USER', 'A'),
        ('MedCo', 'medcoservice', 'DATA_DEID', 'A'),
        ('MedCo', 'medcoservice', 'DATA_OBFSC', 'A'),
        ('MedCo', 'medcoservice', 'DATA_AGG', 'A'),
        ('MedCo', 'medcoservice', 'DATA_LDS', 'A'),
        ('MedCo', 'medcoservice', 'EDITOR', 'A'),
        ('MedCo', 'medcoservice', 'DATA_PROT', 'A'),
        ('MedCo', 'medcoservice', 'MANAGER', 'A'),

        ('MedCo', 'AGG_SERVICE_ACCOUNT', 'USER', 'A'),
        ('MedCo', 'AGG_SERVICE_ACCOUNT', 'MANAGER', 'A'),
        ('MedCo', 'AGG_SERVICE_ACCOUNT', 'DATA_OBFSC', 'A'),
        ('MedCo', 'AGG_SERVICE_ACCOUNT', 'DATA_AGG', 'A'),

        ('MedCo', 'medcoadmin', 'MANAGER', 'A'),
        ('MedCo', 'medcoadmin', 'USER', 'A'),
        ('MedCo', 'medcoadmin', 'DATA_OBFSC', 'A'),

        ('MedCo', 'medcouser', 'USER', 'A'),
        ('MedCo', 'medcouser', 'DATA_DEID', 'A'),
        ('MedCo', 'medcouser', 'DATA_OBFSC', 'A'),
        ('MedCo', 'medcouser', 'DATA_AGG', 'A'),
        ('MedCo', 'medcouser', 'DATA_LDS', 'A'),
        ('MedCo', 'medcouser', 'EDITOR', 'A'),
        ('MedCo', 'medcouser', 'DATA_PROT', 'A'),

        -- give medcouser rights to Demo project
        ('Demo', 'medcouser', 'USER', 'A'),
        ('Demo', 'medcouser', 'DATA_DEID', 'A'),
        ('Demo', 'medcouser', 'DATA_OBFSC', 'A'),
        ('Demo', 'medcouser', 'DATA_AGG', 'A'),
        ('Demo', 'medcouser', 'DATA_LDS', 'A'),
        ('Demo', 'medcouser', 'EDITOR', 'A'),
        ('Demo', 'medcouser', 'DATA_PROT', 'A');
EOSQL
