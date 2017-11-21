#!/bin/bash
set -e
# configuration of the i2b2 project management cell for shrine/medco

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

    -- database lookups
    insert into i2b2hive.ont_db_lookup (c_domain_id, c_project_path, c_owner_id, c_db_fullschema, c_db_datasource, c_db_servertype, c_db_nicename)
        values ('$I2B2_DOMAIN_NAME', 'MedCo-SHRINE/', '@', 'shrine_ont', 'java:/OntologyMedCoShrineDS', 'POSTGRESQL', 'MedCo-SHRINE');
    insert into i2b2hive.crc_db_lookup (c_domain_id, c_project_path, c_owner_id, c_db_fullschema, c_db_datasource, c_db_servertype, c_db_nicename)
        values ('$I2B2_DOMAIN_NAME', '/MedCo-SHRINE/', '@', 'i2b2demodata', 'java:/QueryToolMedCoDS', 'POSTGRESQL', 'MedCo-SHRINE');

    -- hive & users data
    insert into i2b2pm.pm_project_data (project_id, project_name, project_wiki, project_path, status_cd)
        values ('MedCo-SHRINE', 'MedCo-SHRINE', 'https://github.com/lca1/medco', '/MedCo-SHRINE', 'A');
    insert into i2b2pm.pm_cell_data (cell_id, project_path, name, method_cd, url, can_override, status_cd)
        values ('CRC', '/MedCo-SHRINE', 'MedCo-SHRINE Federated Query', 'REST', 'https://shrine-server:6443/shrine/rest/i2b2/', 1, 'A');

    INSERT INTO i2b2pm.PM_USER_DATA (USER_ID, FULL_NAME, PASSWORD, STATUS_CD) VALUES('medcoshrineuser', 'MedCo SHRINE User', 'f8eb764674b57b5710e3c1665464e29', 'A');
    INSERT INTO i2b2pm.PM_USER_DATA (USER_ID, FULL_NAME, PASSWORD, STATUS_CD) VALUES('medcoservice', 'MedCo Service User', '7cb1ac9deab165535494d60da1d3d7e', 'A');

    insert into i2b2pm.pm_project_user_roles (project_id, user_id, user_role_cd, status_cd) values ('MedCo-SHRINE', 'medcoshrineuser', 'USER', 'A');
    insert into i2b2pm.pm_project_user_roles (project_id, user_id, user_role_cd, status_cd) values ('MedCo-SHRINE', 'medcoshrineuser', 'DATA_OBFSC', 'A');

    INSERT INTO i2b2pm.PM_PROJECT_USER_ROLES (PROJECT_ID, USER_ID, USER_ROLE_CD, STATUS_CD) VALUES('MedCo-SHRINE', 'medcoservice', 'USER', 'A'),
        ('MedCo-SHRINE', 'medcoservice', 'DATA_DEID', 'A'),
        ('MedCo-SHRINE', 'medcoservice', 'DATA_OBFSC', 'A'),
        ('MedCo-SHRINE', 'medcoservice', 'DATA_AGG', 'A'),
        ('MedCo-SHRINE', 'medcoservice', 'DATA_LDS', 'A'),
        ('MedCo-SHRINE', 'medcoservice', 'EDITOR', 'A'),
        ('MedCo-SHRINE', 'medcoservice', 'DATA_PROT', 'A'),
        ('MedCo-SHRINE', 'medcoservice', 'MANAGER', 'A'),
        ('MedCo', 'medcoservice', 'USER', 'A'),
        ('MedCo', 'medcoservice', 'DATA_DEID', 'A'),
        ('MedCo', 'medcoservice', 'DATA_OBFSC', 'A'),
        ('MedCo', 'medcoservice', 'DATA_AGG', 'A'),
        ('MedCo', 'medcoservice', 'DATA_LDS', 'A'),
        ('MedCo', 'medcoservice', 'EDITOR', 'A'),
        ('MedCo', 'medcoservice', 'DATA_PROT', 'A'),
        ('MedCo', 'medcoservice', 'MANAGER', 'A');

    -- user parameters
    INSERT INTO i2b2pm.pm_user_params(datatype_cd, user_id, param_name_cd, value, change_date, entry_date, status_cd)
        VALUES('T', 'medcoservice', 'qep', 'true', 'NOW()', 'NOW()', 'A');
    INSERT INTO i2b2pm.pm_user_params(datatype_cd, user_id, param_name_cd, value, change_date, entry_date, status_cd)
        VALUES('T', 'medcoadmin', 'DataSteward', 'true', 'NOW()', 'NOW()', 'A');

EOSQL
