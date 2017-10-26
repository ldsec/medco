#!/bin/bash
set -e

mysql -p$ADMIN_PASSWORD -u root <<-EOSQL

    create database shrine_query_history;
    grant all privileges on shrine_query_history.* to shrine@localhost identified by '$DB_PASSWORD';
    create database stewardDB;
    grant all privileges on stewardDB.* to shrine@localhost identified by '$DB_PASSWORD';
    create database adapterAuditDB;
    grant all privileges on adapterAuditDB.* to shrine@localhost identified by '$DB_PASSWORD';
    create database qepAuditDB;
    grant all privileges on qepAuditDB.* to shrine@localhost identified by '$DB_PASSWORD';

EOSQL

noope
 -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" <<-EOSQL

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

    INSERT into shrine_ont.TABLE_ACCESS
        ( C_TABLE_CD, C_TABLE_NAME, C_PROTECTED_ACCESS, C_HLEVEL, C_NAME, C_FULLNAME, C_SYNONYM_CD, C_VISUALATTRIBUTES, C_TOOLTIP, C_FACTTABLECOLUMN, C_DIMTABLENAME, C_COLUMNNAME, C_COLUMNDATATYPE, C_DIMCODE, C_OPERATOR)
        values ( 'SHRINE', 'SHRINE', 'N', 0, 'SHRINE Ontology', '\SHRINE\', 'N', 'CA', 'SHRINE Ontology', 'concept_cd', 'concept_dimension', 'concept_path', 'T', '\SHRINE\', 'LIKE')
        on conflict do nothing;

    grant all privileges on all tables in schema shrine_ont to shrine_ont;
    grant all privileges on all sequences in schema shrine_ont to shrine_ont;
    grant all privileges on all functions in schema shrine_ont to shrine_ont;
    grant all privileges on all tables in schema shrine_ont to i2b2metadata;
    grant all privileges on all sequences in schema shrine_ont to i2b2metadata;
    grant all privileges on all functions in schema shrine_ont to i2b2metadata;
EOSQL

