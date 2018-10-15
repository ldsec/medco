#!/bin/bash
set -Eeuo pipefail
# set up local medco ontology

psql $PSQL_PARAMS -d "$I2B2_DB_NAME" <<-EOSQL

    -- root nodes in table_access
    insert into i2b2metadata_i2b2.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
        c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
        c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
        ('SENSITIVE_TAGGED', 'SENSITIVE_TAGGED', 'N', 1, '\medco\tagged\', 'MedCo Sensitive Tagged Ontology',
        'N', 'CH', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\medco\tagged\', 'MedCo Sensitive Tagged Ontology');
    insert into i2b2metadata_i2b2.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
        c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
        c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
        ('NON_SENSITIVE_CLEAR', 'NON_SENSITIVE_CLEAR', 'N', 2, '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology',
        'N', 'CH', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology');

    -- medco-specific schemes (prefix to the concept codes)
    insert into i2b2metadata_i2b2.schemes(c_key, c_name, c_description) values('TAG_ID:', 'TAG_ID', 'MedCo tag identifier');
    insert into i2b2metadata_i2b2.schemes(c_key, c_name, c_description) values('ENC_ID:', 'ENC_ID', 'MedCo sensitive concept identifier (to be encrypted)');
    insert into i2b2metadata_i2b2.schemes(c_key, c_name, c_description) values('CLEAR:', 'CLEAR', 'MedCo clear value');
    insert into i2b2metadata_i2b2.schemes(c_key, c_name, c_description) values('GEN:', 'GEN', 'MedCo genomic annotations');

    -- clear non sensitive ontology
    CREATE TABLE i2b2metadata_i2b2.non_sensitive_clear (
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
    ALTER TABLE ONLY i2b2metadata_i2b2.non_sensitive_clear ADD CONSTRAINT fullname_pk_10 PRIMARY KEY (c_fullname);
    ALTER TABLE ONLY i2b2metadata_i2b2.non_sensitive_clear ADD CONSTRAINT basecode_un_10 UNIQUE (c_basecode);
    --insert into i2b2metadata_i2b2.non_sensitive_clear (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
    --c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
    --download_date, import_date, valuetype_cd, m_applied_path) values
    --('2', '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology', 'N', 'CA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
    --'T', 'LIKE', '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology', '\medco\clinical\nonsensitive\',
    --'NOW()', 'NOW()', 'NOW()', 'CLEAR', '@');


    -- tagged sensitive ontology
    CREATE TABLE i2b2metadata_i2b2.sensitive_tagged (
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
    ALTER TABLE ONLY i2b2metadata_i2b2.sensitive_tagged ADD CONSTRAINT fullname_pk_11 PRIMARY KEY (c_fullname);
    ALTER TABLE ONLY i2b2metadata_i2b2.sensitive_tagged ADD CONSTRAINT basecode_un_11 UNIQUE (c_basecode);
    --insert into i2b2metadata_i2b2.sensitive_tagged (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
    --c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
    --download_date, import_date, valuetype_cd, m_applied_path) values
    --('1', '\medco\tagged\', 'MedCo Sensitive Tagged Ontology', 'N', 'CA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
    --'T', 'LIKE', '\medco\tagged\', 'MedCo Sensitive Tagged Ontology', '\medco\tagged\', 'NOW()', 'NOW()', 'NOW()', 'TAG_ID', '@');


    -- permissions
    alter table i2b2metadata_i2b2.sensitive_tagged owner to $I2B2_DB_USER;
    alter table i2b2metadata_i2b2.non_sensitive_clear owner to $I2B2_DB_USER;
    grant all on schema i2b2metadata_i2b2 to $I2B2_DB_USER;
    grant all privileges on all tables in schema i2b2metadata_i2b2 to $I2B2_DB_USER;
    grant all privileges on all sequences in schema i2b2metadata_i2b2 to $I2B2_DB_USER;
    grant all privileges on all functions in schema i2b2metadata_i2b2 to $I2B2_DB_USER;
EOSQL
