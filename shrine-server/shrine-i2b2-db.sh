#!/bin/bash
set -e

# db lookups TODO: password hardcoded, use compose secrets
psql -v ON_ERROR_STOP=1 postgresql://i2b2hive:demouser@i2b2-database/$I2B2_DOMAIN_NAME <<-EOSQL
    insert into i2b2hive.ont_db_lookup (c_domain_id, c_project_path, c_owner_id, c_db_fullschema, c_db_datasource, c_db_servertype, c_db_nicename)
    values ($I2B2_DOMAIN_NAME, 'SHRINE/', '@', 'shrine_ont', 'java:/OntologyShrineDS', 'POSTGRESQL', 'SHRINE')
    on conflict do nothing;

    insert into i2b2hive.crc_db_lookup (c_domain_id, c_project_path, c_owner_id, c_db_fullschema, c_db_datasource, c_db_servertype, c_db_nicename)
    values ($I2B2_DOMAIN_NAME, '/SHRINE/', '@', 'i2b2demodata', 'java:/QueryToolDemoDS', 'POSTGRESQL', 'SHRINE')
    on conflict do nothing;

EOSQL

# add shrine user in pm TODO: password hardcoded, use compose secrets
cd "install/i2b2-1.7/i2b2"
javac ./I2b2PasswordCryptor.java
SHRINE_PW=$(java -classpath ./ I2b2PasswordCryptor demouser)
cd "$SHRINE_SRC_DIR"

# TODO: password hardcoded, use compose secrets
# TODO: check if the external address is OK (maybe pass by argument), port?
psql -v ON_ERROR_STOP=1 postgresql://i2b2pm:demouser@i2b2-database/$I2B2_DOMAIN_NAME <<-EOSQL
    insert into i2b2pm.pm_user_data (user_id, full_name, password, status_cd)
    values ('shrine', 'shrine', $SHRINE_PW, 'A');

    insert into i2b2pm.pm_project_data (project_id, project_name, project_wiki, project_path, status_cd)
    values ('SHRINE', 'SHRINE', 'http://open.med.harvard.edu/display/SHRINE', '/SHRINE', 'A');

    insert into i2b2pm.pm_project_user_roles (project_id, user_id, user_role_cd, status_cd)
    values ('SHRINE', 'shrine', 'USER', 'A');

    insert into i2b2pm.pm_project_user_roles (project_id, user_id, user_role_cd, status_cd)
    values ('SHRINE', 'shrine', 'DATA_OBFSC', 'A');

    insert into i2b2pm.pm_cell_data (cell_id, project_path, name, method_cd, url, can_override, status_cd)
    values ('CRC', '/SHRINE', 'SHRINE Federated Query', 'REST', 'https://shrine-server:6443/shrine/rest/i2b2/', 1, 'A');


EOSQL

# add shrine ontology TODO: password hardcoded, use compose secrets
psql -v ON_ERROR_STOP=1 postgresql://postgres:prigen2017@i2b2-database/$I2B2_DOMAIN_NAME <<-EOSQL
    DO
    \$body\$
    BEGIN
        IF NOT EXISTS (
            SELECT *
            FROM   pg_catalog.pg_user
            WHERE  usename = shrine_ont) THEN

            create user shrine_ont with password 'demouser';
        END IF;
    END
    \$body\$;

    create schema IF NOT EXISTS authorization shrine_ont;
    grant all privileges on all tables in schema shrine_ont to shrine_ont;
    grant all privileges on all sequences in schema shrine_ont to shrine_ont;
    grant all privileges on all functions in schema shrine_ont to shrine_ont;

    CREATE TABLE IF NOT EXISTS shrine_ont.SHRINE
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

    ALTER TABLE shrine_ont.SHRINE OWNER TO shrine_ont;
    GRANT SELECT, INSERT, UPDATE, DELETE ON shrine_ont.SHRINE TO i2b2metadata;
    GRANT SELECT, INSERT, UPDATE, DELETE ON shrine_ont.SHRINE TO shrine_ont;

    CREATE TABLE IF NOT EXISTS shrine_ont.TABLE_ACCESS
    (
        C_TABLE_CD VARCHAR(50),
        C_TABLE_NAME VARCHAR(50),
        C_PROTECTED_ACCESS CHAR(1),
        C_HLEVEL NUMERIC(22,0),
        C_FULLNAME VARCHAR(900),
        C_NAME VARCHAR(2000),
        C_SYNONYM_CD CHAR(1),
        C_VISUALATTRIBUTES CHAR(3),
        C_TOTALNUM NUMERIC(22,0),
        C_BASECODE VARCHAR(450),
        C_METADATAXML TEXT,
        C_FACTTABLECOLUMN VARCHAR(50),
        C_DIMTABLENAME VARCHAR(50),
        C_COLUMNNAME VARCHAR(50),
        C_COLUMNDATATYPE VARCHAR(50),
        C_OPERATOR VARCHAR(10),
        C_DIMCODE VARCHAR(900),
        C_COMMENT TEXT,
        C_TOOLTIP VARCHAR(900),
        C_ENTRY_DATE DATE,
        C_CHANGE_DATE DATE,
        C_STATUS_CD CHAR(1),
        VALUETYPE_CD VARCHAR(50)
    );

    ALTER TABLE shrine_ont.TABLE_ACCESS OWNER TO shrine_ont;
    GRANT SELECT, INSERT, UPDATE, DELETE ON shrine_ont.TABLE_ACCESS TO i2b2metadata;
    GRANT SELECT, INSERT, UPDATE, DELETE ON shrine_ont.TABLE_ACCESS TO shrine_ont;

    CREATE TABLE IF NOT EXISTS shrine_ont.SCHEMES
    (
        C_KEY VARCHAR(50) NOT NULL,
        C_NAME VARCHAR(50) NOT NULL,
        C_DESCRIPTION VARCHAR(100),
        CONSTRAINT SCHEMES_PK PRIMARY KEY (C_KEY)
    ) ;

    ALTER TABLE shrine_ont.SCHEMES OWNER TO shrine_ont;
    GRANT SELECT, INSERT, UPDATE, DELETE ON shrine_ont.SCHEMES TO i2b2metadata;
    GRANT SELECT, INSERT, UPDATE, DELETE ON shrine_ont.SCHEMES TO shrine_ont;

    INSERT into shrine_ont.TABLE_ACCESS
        ( C_TABLE_CD, C_TABLE_NAME, C_PROTECTED_ACCESS, C_HLEVEL, C_NAME, C_FULLNAME, C_SYNONYM_CD, C_VISUALATTRIBUTES, C_TOOLTIP, C_FACTTABLECOLUMN, C_DIMTABLENAME, C_COLUMNNAME, C_COLUMNDATATYPE, C_DIMCODE, C_OPERATOR)
        values ( 'SHRINE', 'SHRINE', 'N', 0, 'SHRINE Ontology', '\SHRINE\', 'N', 'CA', 'SHRINE Ontology', 'concept_cd', 'concept_dimension', 'concept_path', 'T', '\SHRINE\', 'LIKE')
        on conflict do nothing;

    GRANT usage ON schema i2b2metadata TO shrine_ont;
    GRANT SELECT, INSERT, UPDATE, DELETE on all tables in schema i2b2metadata TO shrine_ont;
EOSQL

# load shrine ontology TODO: password hardcoded, use compose secrets
wget https://open.med.harvard.edu/svn/shrine-ontology/SHRINE_Demo_Downloads/trunk/ShrineDemo.sql
mv ShrineDemo.sql Shrine.sql
sed -i '1s/^/SET search_path TO shrine_ont;\n/' Shrine.sql
psql -v ON_ERROR_STOP=1 postgresql://shrine_ont:demouser@i2b2-database/$I2B2_DOMAIN_NAME < Shrine.sql


