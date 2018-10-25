#!/bin/bash
set -Eeuo pipefail
# set up common medco ontology

psql $PSQL_PARAMS -d "$I2B2_DB_NAME" <<-EOSQL

    -- table_access & root nodes
    CREATE TABLE medco_ont.table_access(
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
    --INSERT into medco_ont.table_access
    --    ( C_TABLE_CD, C_TABLE_NAME, C_PROTECTED_ACCESS, C_HLEVEL, C_NAME, C_FULLNAME, C_SYNONYM_CD, C_VISUALATTRIBUTES,
    --    C_TOOLTIP, C_FACTTABLECOLUMN, C_DIMTABLENAME, C_COLUMNNAME, C_COLUMNDATATYPE, C_DIMCODE, C_OPERATOR) values
    --    ( 'MedCo', 'MedCo', 'N', 0, 'MedCo Ontology', '\medco\', 'N', 'CA', 'MedCo Ontology', 'concept_cd',
    --    'concept_dimension', 'concept_path', 'T', '\medco\', 'LIKE');
    insert into medco_ont.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
        c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
        c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
        ('CLINICAL_SENSITIVE', 'CLINICAL_SENSITIVE', 'N', 2, '\medco\clinical\sensitive\', 'MedCo Clinical Sensitive Ontology',
        'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\medco\clinical\sensitive\', 'MedCo Clinical Sensitive Ontology');
    insert into medco_ont.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
        c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
        c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
        ('CLINICAL_NON_SENSITIVE', 'CLINICAL_NON_SENSITIVE', 'N', 2, '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology',
        'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology');
    insert into medco_ont.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
        c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
        c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
        ('GENOMIC', 'GENOMIC', 'N', 1, '\medco\genomic\', 'MedCo Genomic Ontology',
        'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\medco\genomic\', 'MedCo Genomic Ontology');
    insert into medco_ont.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
        c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
        c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
        ('SENSITIVE_TAGGED', 'SENSITIVE_TAGGED', 'N', 1, '\medco\tagged\', 'MedCo Sensitive Tagged Ontology',
        'N', 'CH', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\medco\tagged\', 'MedCo Sensitive Tagged Ontology');

    -- schemes
    CREATE TABLE medco_ont.schemes(
        C_KEY VARCHAR(50) NOT NULL,
        C_NAME VARCHAR(50) NOT NULL,
        C_DESCRIPTION VARCHAR(100),
        CONSTRAINT SCHEMES_PK PRIMARY KEY (C_KEY)
    );
    insert into medco_ont.schemes(c_key, c_name, c_description) values('TAG_ID:', 'TAG_ID', 'MedCo tag identifier');
    insert into medco_ont.schemes(c_key, c_name, c_description) values('ENC_ID:', 'ENC_ID', 'MedCo sensitive concept identifier (to be encrypted)');
    insert into medco_ont.schemes(c_key, c_name, c_description) values('CLEAR:', 'CLEAR', 'MedCo clear value');
    insert into medco_ont.schemes(c_key, c_name, c_description) values('GEN:', 'GEN', 'MedCo genomic annotations');

    -- ontology table
    CREATE TABLE medco_ont.shrine(
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

    -- clinical sensitive ontology
    CREATE TABLE medco_ont.clinical_sensitive(
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
    ALTER TABLE ONLY medco_ont.clinical_sensitive ADD CONSTRAINT fullname_pk_20 PRIMARY KEY (c_fullname);
    ALTER TABLE ONLY medco_ont.clinical_sensitive ADD CONSTRAINT basecode_un_20 UNIQUE (c_basecode);
    --insert into medco_ont.clinical_sensitive (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
    --    c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
    --    download_date, import_date, valuetype_cd, m_applied_path) values
    --    ('2', '\medco\clinical\sensitive\', 'MedCo Clinical Sensitive Ontology', 'N', 'CA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
    --    'T', 'LIKE', '\medco\clinical\sensitive\', 'MedCo Clinical Sensitive Ontology', '\medco\clinical\sensitive\',
    --    'NOW()', 'NOW()', 'NOW()', 'ENC_ID', '@');

    -- clinical non-sensitive ontology
    CREATE TABLE medco_ont.clinical_non_sensitive(
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
    ALTER TABLE ONLY medco_ont.clinical_non_sensitive ADD CONSTRAINT fullname_pk_21 PRIMARY KEY (c_fullname);
    ALTER TABLE ONLY medco_ont.clinical_non_sensitive ADD CONSTRAINT basecode_un_21 UNIQUE (c_basecode);
    --insert into medco_ont.clinical_non_sensitive (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
    --    c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
    --    download_date, import_date, valuetype_cd, m_applied_path) values
    --    ('2', '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology', 'N', 'CA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
    --    'T', 'LIKE', '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology', '\medco\clinical\nonsensitive\',
    --    'NOW()', 'NOW()', 'NOW()', 'CLEAR', '@');

    -- tagged sensitive ontology
    CREATE TABLE medco_ont.sensitive_tagged (
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
    ALTER TABLE ONLY medco_ont.sensitive_tagged ADD CONSTRAINT fullname_pk_11 PRIMARY KEY (c_fullname);
    ALTER TABLE ONLY medco_ont.sensitive_tagged ADD CONSTRAINT basecode_un_11 UNIQUE (c_basecode);
    --insert into medco_ont.sensitive_tagged (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
    --c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
    --download_date, import_date, valuetype_cd, m_applied_path) values
    --('1', '\medco\tagged\', 'MedCo Sensitive Tagged Ontology', 'N', 'CA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
    --'T', 'LIKE', '\medco\tagged\', 'MedCo Sensitive Tagged Ontology', '\medco\tagged\', 'NOW()', 'NOW()', 'NOW()', 'TAG_ID', '@');

    -- genomic ontology
    CREATE TABLE medco_ont.genomic(
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
    ALTER TABLE ONLY medco_ont.genomic ADD CONSTRAINT fullname_pk_22 PRIMARY KEY (c_fullname);
    ALTER TABLE ONLY medco_ont.genomic ADD CONSTRAINT basecode_un_22 UNIQUE (c_basecode);
    insert into medco_ont.genomic (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
        c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
        download_date, import_date, valuetype_cd, m_applied_path) values
        ('1', '\medco\genomic\', 'MedCo Genomic Ontology', 'N', 'CA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
        'T', 'LIKE', '\medco\genomic\', 'MedCo Genomic Ontology', '\medco\genomic\',
        'NOW()', 'NOW()', 'NOW()', 'GEN', '@');

    -- permissions
    ALTER TABLE medco_ont.table_access OWNER TO $I2B2_DB_USER;
    ALTER TABLE medco_ont.schemes OWNER TO $I2B2_DB_USER;
    ALTER TABLE medco_ont.shrine OWNER TO $I2B2_DB_USER;
    ALTER TABLE medco_ont.genomic OWNER TO $I2B2_DB_USER;
    ALTER TABLE medco_ont.clinical_sensitive OWNER TO $I2B2_DB_USER;
    ALTER TABLE medco_ont.clinical_non_sensitive OWNER TO $I2B2_DB_USER;
    ALTER TABLE medco_ont.sensitive_tagged OWNER TO $I2B2_DB_USER;

    grant all on schema medco_ont to $I2B2_DB_USER;
    grant all privileges on all tables in schema medco_ont to $I2B2_DB_USER;
    grant all privileges on all sequences in schema medco_ont to $I2B2_DB_USER;
    grant all privileges on all functions in schema medco_ont to $I2B2_DB_USER;

EOSQL
