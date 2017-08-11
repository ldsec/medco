#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_DOMAIN_NAME" <<-EOSQL
-- todo: integrate medco ontology into shrine

-- ################### CRC DB ###################

-- add custom column in the patient dimension table containing the dummy flag
ALTER TABLE i2b2demodata.patient_dimension ADD COLUMN enc_dummy_flag_cd character(88);
COMMENT ON COLUMN i2b2demodata.patient_dimension.enc_dummy_flag_cd IS 'base64-encoded encrypted dummy flag (0 or 1)';
INSERT INTO i2b2demodata.code_lookup VALUES ('patient_dimension', 'enc_dummy_flag_cd', 'CRC_COLUMN_DESCRIPTOR', 'Encrypted Dummy Flag', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1);



-- ################### ONT DB ###################
-- i2b2 ontology medco: in postgresql i2b2 DB!
-- todo: currently added in i2b2metadata demo table

-- add medco-specific schemes (prefix to the concept codes)
insert into i2b2metadata.schemes(c_key, c_name, c_description) values
    ('TAG_ID:', 'TAG_ID', 'MedCo tag identifier');
insert into i2b2metadata.schemes(c_key, c_name, c_description) values
    ('ENC_ID:', 'ENC_ID', 'MedCo sensitive concept identifier (to be encrypted)');
insert into i2b2metadata.schemes(c_key, c_name, c_description) values
    ('CLEAR:', 'CLEAR', 'MedCo clear value');

-- clinical sensitive ontology
CREATE TABLE i2b2metadata.clinical_sensitive (
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
ALTER TABLE ONLY i2b2metadata.clinical_sensitive
    ADD CONSTRAINT fullname_pk PRIMARY KEY (c_fullname);
ALTER TABLE ONLY i2b2metadata.clinical_sensitive
    ADD CONSTRAINT basecode_un UNIQUE (c_basecode);

-- clinical non-sensitive ontology
CREATE TABLE i2b2metadata.clinical_non_sensitive (
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
ALTER TABLE ONLY i2b2metadata.clinical_non_sensitive
    ADD CONSTRAINT fullname_pk2 PRIMARY KEY (c_fullname);
ALTER TABLE ONLY i2b2metadata.clinical_non_sensitive
    ADD CONSTRAINT basecode_un2 UNIQUE (c_basecode);

-- genomic ontology
CREATE TABLE i2b2metadata.genomic (
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
ALTER TABLE ONLY i2b2metadata.genomic
    ADD CONSTRAINT fullname_pk3 PRIMARY KEY (c_fullname);
ALTER TABLE ONLY i2b2metadata.genomic
    ADD CONSTRAINT basecode_un3 UNIQUE (c_basecode);


-- table_access table
insert into i2b2metadata.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
    c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
    c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
    ('CLINICAL_SENSITIVE', 'CLINICAL_SENSITIVE', 'N', 2, '\\medco\\clinical\\sensitive\\', 'MedCo Clinical Sensitive Ontology',
    'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\sensitive\\', 'MedCo Clinical Sensitive Ontology');
insert into i2b2metadata.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
    c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
    c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
    ('CLINICAL_NON_SENSITIVE', 'CLINICAL_NON_SENSITIVE', 'N', 2, '\\medco\\clinical\\nonsensitive\\', 'MedCo Clinical Non-Sensitive Ontology',
    'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\nonsensitive\\', 'MedCo Clinical Non-Sensitive Ontology');
insert into i2b2metadata.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
    c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
    c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
    ('GENOMIC', 'GENOMIC', 'N', 2, '\\medco\\genomic\\', 'MedCo Genomic Entries',
    'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\genomic\\', 'MedCo Genomic Entries');


-- tables ont init
insert into i2b2metadata.clinical_sensitive (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
download_date, import_date, valuetype_cd, m_applied_path) values
('2', '\medco\clinical\sensitive\', 'MedCo Clinical Sensitive Ontology', 'N', 'CA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
'T', 'LIKE', '\medco\clinical\sensitive\', 'MedCo Clinical Sensitive Ontology', '\medco\clinical\sensitive\',
'NOW()', 'NOW()', 'NOW()', 'MEDCO_ENC', '@');
insert into i2b2metadata.clinical_non_sensitive (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
download_date, import_date, valuetype_cd, m_applied_path) values
('2', '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology', 'N', 'CA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
'T', 'LIKE', '\medco\clinical\nonsensitive\', 'MedCo Clinical Non-Sensitive Ontology', '\medco\clinical\nonsensitive\',
'NOW()', 'NOW()', 'NOW()', 'MEDCO_CLEAR', '@');
insert into i2b2metadata.genomic (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
download_date, import_date, valuetype_cd, m_applied_path) values
('1', '\medco\genomic\', 'MedCo Genomic Entries', 'N', 'CA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
'T', 'LIKE', '\medco\genomic\', 'MedCo Genomic Entries', '\medco\genomic\',
'NOW()', 'NOW()', 'NOW()', 'MEDCO_GEN', '@');


-- fix permissions
alter table i2b2metadata.clinical_sensitive owner to i2b2metadata;
alter table i2b2metadata.clinical_non_sensitive owner to i2b2metadata;
alter table i2b2metadata.genomic owner to i2b2metadata;
grant all on schema i2b2metadata to i2b2metadata;
grant all privileges on all tables in schema i2b2metadata to i2b2metadata;
--alter table medco_data.enc_observation_fact owner to medco_data;


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

EOSQL