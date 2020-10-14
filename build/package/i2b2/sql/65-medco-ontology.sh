#!/bin/bash
set -Eeuo pipefail
# set up common medco ontology

psql $PSQL_PARAMS -d "$I2B2_DB_NAME" <<-EOSQL

    -- table_access
    CREATE TABLE medco_ont_sensitive.table_access(
        C_TABLE_CD VARCHAR(50) PRIMARY KEY,
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
    ALTER TABLE medco_ont_sensitive.table_access OWNER TO $I2B2_DB_USER;

    -- table_access
    CREATE TABLE medco_ont_non_sensitive.table_access(
        C_TABLE_CD VARCHAR(50) PRIMARY KEY,
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
    ALTER TABLE medco_ont_non_sensitive.table_access OWNER TO $I2B2_DB_USER;

    insert into medco_ont_sensitive.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
        c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
        c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
        ('SENSITIVE_TAGGED', 'SENSITIVE_TAGGED', 'N', 1, '\medco\tagged\', 'MedCo Sensitive Tagged Ontology',
        'N', 'CH', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\medco\tagged\', 'MedCo Sensitive Tagged Ontology');

    -- schemes
    CREATE TABLE medco_ont_sensitive.schemes(
        C_KEY VARCHAR(50) NOT NULL,
        C_NAME VARCHAR(50) NOT NULL,
        C_DESCRIPTION VARCHAR(100),
        CONSTRAINT SCHEMES_PK PRIMARY KEY (C_KEY)
    );
    ALTER TABLE medco_ont_sensitive.schemes OWNER TO $I2B2_DB_USER;

    -- todo: revise those schemes
    insert into medco_ont_sensitive.schemes(c_key, c_name, c_description) values('TAG_ID:', 'TAG_ID', 'MedCo tag identifier');
    insert into medco_ont_sensitive.schemes(c_key, c_name, c_description) values('ENC_ID:', 'ENC_ID', 'MedCo sensitive concept identifier (to be encrypted)');
    insert into medco_ont_sensitive.schemes(c_key, c_name, c_description) values('CLEAR:', 'CLEAR', 'MedCo clear value');
    insert into medco_ont_sensitive.schemes(c_key, c_name, c_description) values('GEN:', 'GEN', 'MedCo genomic annotations');

    -- tagged sensitive ontology
    CREATE TABLE medco_ont_sensitive.sensitive_tagged (
        c_hlevel numeric(22,0) not null,
        c_fullname character varying(900) not null,
        c_name character varying(2000) not null,
        c_synonym_cd character(1) not null,
        c_visualattributes character(3) not null,
        c_totalnum numeric(22,0),
        c_basecode character varying(450),
        c_metadataxml text,
        c_facttablecolumn character varying(50) not null,
        c_tablename character varying(50) not null,
        c_columnname character varying(50) not null,
        c_columndatatype character varying(50) not null,
        c_operator character varying(10) not null,
        c_dimcode character varying(900) not null,
        c_comment text,
        c_tooltip character varying(900),
        update_date date not null,
        download_date date,
        import_date date,
        sourcesystem_cd character varying(50),
        valuetype_cd character varying(50),
        m_applied_path character varying(900) not null,
        m_exclusion_cd character varying(900),
        c_path character varying(700),
        c_symbol character varying(50),
        pcori_basecode character varying(50)
    );
    ALTER TABLE medco_ont_sensitive.sensitive_tagged OWNER TO $I2B2_DB_USER;
    CREATE INDEX META_FULLNAME_IDX_sensitive_tagged ON medco_ont_sensitive.sensitive_tagged(C_FULLNAME);
    CREATE INDEX META_APPLIED_PATH_IDX_sensitive_tagged ON medco_ont_sensitive.sensitive_tagged(M_APPLIED_PATH);
    CREATE INDEX META_EXCLUSION_IDX_sensitive_tagged ON medco_ont_sensitive.sensitive_tagged(M_EXCLUSION_CD);
    CREATE INDEX META_HLEVEL_IDX_sensitive_tagged ON medco_ont_sensitive.sensitive_tagged(C_HLEVEL);
    CREATE INDEX META_SYNONYM_IDX_sensitive_tagged ON medco_ont_sensitive.sensitive_tagged(C_SYNONYM_CD);

    -- permissions
    grant all on schema medco_ont_sensitive to $I2B2_DB_USER;
    grant all privileges on all tables in schema medco_ont_sensitive to $I2B2_DB_USER;
    grant all privileges on all sequences in schema medco_ont_sensitive to $I2B2_DB_USER;
    grant all privileges on all functions in schema medco_ont_sensitive to $I2B2_DB_USER;

    grant all on schema medco_ont_non_sensitive to $I2B2_DB_USER;
    grant all privileges on all tables in schema medco_ont_non_sensitive to $I2B2_DB_USER;
    grant all privileges on all sequences in schema medco_ont_non_sensitive to $I2B2_DB_USER;
    grant all privileges on all functions in schema medco_ont_non_sensitive to $I2B2_DB_USER;
EOSQL
