#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" <<-EOSQL
-- ################### I2B2 ONT DB ###################

-- table_access table
insert into i2b2metadata.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
    c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
    c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
    ('SENSITIVE_TAGGED', 'SENSITIVE_TAGGED', 'N', 1, '\medco\tagged\', 'MedCo Sensitive Tagged Ontology',
    'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\medco\tagged\', 'MedCo Sensitive Tagged Ontology');
insert into i2b2metadata.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
    c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
    c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
    ('NON_SENSITIVE_CLEAR', 'NON_SENSITIVE_CLEAR', 'N', 2, '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology',
    'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology');


-- add medco-specific schemes (prefix to the concept codes)
insert into i2b2metadata.schemes(c_key, c_name, c_description) values
    ('TAG_ID:', 'TAG_ID', 'MedCo tag identifier');
insert into i2b2metadata.schemes(c_key, c_name, c_description) values
    ('ENC_ID:', 'ENC_ID', 'MedCo sensitive concept identifier (to be encrypted)');
insert into i2b2metadata.schemes(c_key, c_name, c_description) values
    ('CLEAR:', 'CLEAR', 'MedCo clear value');
insert into i2b2metadata.schemes(c_key, c_name, c_description) values
    ('GEN:', 'GEN', 'MedCo genomic annotations');


-- clear non sensitive ontology
CREATE TABLE i2b2metadata.non_sensitive_clear (
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
ALTER TABLE ONLY i2b2metadata.non_sensitive_clear
    ADD CONSTRAINT fullname_pk_10 PRIMARY KEY (c_fullname);
ALTER TABLE ONLY i2b2metadata.non_sensitive_clear
    ADD CONSTRAINT basecode_un_10 UNIQUE (c_basecode);
insert into i2b2metadata.non_sensitive_clear (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
download_date, import_date, valuetype_cd, m_applied_path) values
('2', '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology', 'N', 'CA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
'T', 'LIKE', '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology', '\medco\clinical\nonsensitive\',
'NOW()', 'NOW()', 'NOW()', 'CLEAR', '@');


-- tagged sensitive ontology
CREATE TABLE i2b2metadata.sensitive_tagged (
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
ALTER TABLE ONLY i2b2metadata.sensitive_tagged
    ADD CONSTRAINT fullname_pk_11 PRIMARY KEY (c_fullname);
ALTER TABLE ONLY i2b2metadata.sensitive_tagged
    ADD CONSTRAINT basecode_un_11 UNIQUE (c_basecode);
insert into i2b2metadata.sensitive_tagged (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
download_date, import_date, valuetype_cd, m_applied_path) values
('1', '\medco\tagged\', 'MedCo Sensitive Tagged Ontology', 'N', 'CA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
'T', 'LIKE', '\medco\tagged\', 'MedCo Sensitive Tagged Ontology', '\medco\tagged\', 'NOW()', 'NOW()', 'NOW()', 'ENC_ID', '@');


-- fix permissions
alter table i2b2metadata.sensitive_tagged owner to i2b2metadata;
alter table i2b2metadata.non_sensitive_clear owner to i2b2metadata;
grant all on schema i2b2metadata to i2b2metadata;
grant all privileges on all tables in schema i2b2metadata to i2b2metadata;
grant all privileges on all sequences in schema i2b2metadata to i2b2metadata;
grant all privileges on all functions in schema i2b2metadata to i2b2metadata;


-- ################### OLD ###################
--create schema medco_data;

--create table medco_data.enc_observation_fact (
--    encounter_id character varying(50) not null,
--    patient_id character varying(50) not null,
--    provider_id character varying(50) not null,
--    enc_concept_id character varying(500) not null
--);


-- access rights
--create role medco_data login password 'demouser';
--grant all on schema medco_data to medco_data;
--grant all privileges on all tables in schema medco_data to medco_data;


-- next usable id init
--insert into i2b2metadata.clinical_sensitive (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
--c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
--download_date, import_date, valuetype_cd, m_applied_path) values
--('3', '\medco\clinical\sensitive\next_usable_id\', 'nextEncUsableId', 'N', 'LH', '0', 'concept_cd', 'concept_dimension', 'concept_path',
--'T', 'LIKE', '\medco\clinical\sensitive\next_usable_id\', 'The next usable id for Unlynx encryption.', '\medco\clinical\sensitive\next_usable_id\',
--'NOW()', 'NOW()', 'NOW()', 'MEDCO_ADMIN', '@');
--insert into i2b2metadata.clinical_non_sensitive (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
--c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
--download_date, import_date, valuetype_cd, m_applied_path) values
--('3', '\medco\clinical\nonsensitive\next_usable_id\', 'nextClearUsableId', 'N', 'LH', '0', 'concept_cd', 'concept_dimension', 'concept_path',
--'T', 'LIKE', '\medco\clinical\nonsensitive\next_usable_id\', 'The next usable id for non sensitive clinical data.', '\medco\clinical\nonsensitive\next_usable_id\',
--'NOW()', 'NOW()', 'NOW()', 'MEDCO_ADMIN', '@');

--insert into i2b2metadata.clinical_sensitive (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum, c_basecode,
--c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, update_date,
--download_date, import_date, valuetype_cd, m_applied_path) values
--('4', '\medco\clinical\sensitive\next_usable_id\1\', '1', 'N', 'LH', '0', 'MEDCO_ADMIN:nextEncUsableId', 'concept_cd', 'concept_dimension', 'concept_path',
--'T', 'LIKE', '\medco\clinical\sensitive\next_usable_id\1\', 'The next usable id for Unlynx encryption.',
--'NOW()', 'NOW()', 'NOW()', 'MEDCO_ADMIN', '@');
--insert into i2b2metadata.clinical_non_sensitive (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum, c_basecode,
--c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, update_date,
--download_date, import_date, valuetype_cd, m_applied_path) values
--('4', '\medco\clinical\nonsensitive\next_usable_id\1\', '1', 'N', 'LH', '0', 'MEDCO_ADMIN:nextClearUsableId', 'concept_cd', 'concept_dimension', 'concept_path',
--'T', 'LIKE', '\medco\clinical\nonsensitive\next_usable_id\1\', 'The next usable id for non sensitive clinical data.',
--'NOW()', 'NOW()', 'NOW()', 'MEDCO_ADMIN', '@');

--alter table medco_data.enc_observation_fact owner to medco_data;


EOSQL