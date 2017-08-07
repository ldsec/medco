-- ################### medco_data ###################
create schema medco_data;

create table medco_data.enc_observation_fact (
    encounter_id character varying(50) not null,
    patient_id character varying(50) not null,
    provider_id character varying(50) not null,
    enc_concept_id character varying(500) not null
);


-- access rights
create role medco_data login password 'demouser';
grant all on schema medco_data to medco_data;
grant all privileges on all tables in schema medco_data to medco_data;

-- ################### medco_ontology ###################
-- i2b2 ontology medco: in postgresql i2b2 DB!
-- todo: currently added in i2b2metadata demo table
CREATE SCHEMA medco_ontology;

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

-- schemes table
CREATE TABLE i2b2metadata.schemes (
    c_key character varying(50) NOT NULL,
    c_name character varying(50) NOT NULL,
    c_description character varying(100)
);
ALTER TABLE ONLY i2b2metadata.schemes
    ADD CONSTRAINT schemes_pk PRIMARY KEY (c_key);

insert into i2b2metadata.schemes(c_key, c_name, c_description) values
    ('MEDCO_ENC:', 'MEDCO_ENC', 'MedCo encrypted values');
insert into i2b2metadata.schemes(c_key, c_name, c_description) values
    ('MEDCO_CLEAR:', 'MEDCO_CLEAR', 'MedCo clear values');
insert into i2b2metadata.schemes(c_key, c_name, c_description) values
    ('MEDCO_ADMIN:', 'MEDCO_ADMIN', 'MedCo administrative entries');
insert into i2b2metadata.schemes(c_key, c_name, c_description) values
    ('MEDCO_GEN:', 'MEDCO_GEN', 'MedCo genomic entries');

-- table_access table
CREATE TABLE medco_ontology.table_access (
    c_table_cd character varying(50),
    c_table_name character varying(50),
    c_protected_access character(1),
    c_hlevel numeric(22,0),
    c_fullname character varying(900),
    c_name character varying(2000),
    c_synonym_cd character(1),
    c_visualattributes character(3),
    c_totalnum numeric(22,0),
    c_basecode character varying(450),
    c_metadataxml text,
    c_facttablecolumn character varying(50),
    c_dimtablename character varying(50),
    c_columnname character varying(50),
    c_columndatatype character varying(50),
    c_operator character varying(10),
    c_dimcode character varying(900),
    c_comment text,
    c_tooltip character varying(900),
    c_entry_date date,
    c_change_date date,
    c_status_cd character(1),
    valuetype_cd character varying(50)
);

insert into medco_ontology.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
    c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
    c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
    ('CLINICAL_SENSITIVE', 'CLINICAL_SENSITIVE', 'N', 2, '\\medco\\clinical\\sensitive\\', 'MedCo Clinical Sensitive Ontology',
    'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\sensitive\\', 'MedCo Clinical Sensitive Ontology');
insert into medco_ontology.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
    c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
    c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
    ('CLINICAL_NON_SENSITIVE', 'CLINICAL_NON_SENSITIVE', 'N', 2, '\\medco\\clinical\\nonsensitive\\', 'MedCo Clinical Non-Sensitive Ontology',
    'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\nonsensitive\\', 'MedCo Clinical Non-Sensitive Ontology');
insert into medco_ontology.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
    c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
    c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
    ('GENOMIC', 'GENOMIC', 'N', 2, '\\medco\\genomic\\', 'MedCo Genomic Entries',
    'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\genomic\\', 'MedCo Genomic Entries');

-- access rights
create role medco_ontology login password 'demouser';
grant all on schema medco_ontology to medco_ontology;
grant all privileges on all tables in schema medco_ontology to medco_ontology;

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


-- next usable id init
insert into i2b2metadata.clinical_sensitive (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
download_date, import_date, valuetype_cd, m_applied_path) values
('3', '\medco\clinical\sensitive\next_usable_id\', 'nextEncUsableId', 'N', 'LH', '0', 'concept_cd', 'concept_dimension', 'concept_path',
'T', 'LIKE', '\medco\clinical\sensitive\next_usable_id\', 'The next usable id for Unlynx encryption.', '\medco\clinical\sensitive\next_usable_id\',
'NOW()', 'NOW()', 'NOW()', 'MEDCO_ADMIN', '@');
insert into i2b2metadata.clinical_non_sensitive (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
download_date, import_date, valuetype_cd, m_applied_path) values
('3', '\medco\clinical\nonsensitive\next_usable_id\', 'nextClearUsableId', 'N', 'LH', '0', 'concept_cd', 'concept_dimension', 'concept_path',
'T', 'LIKE', '\medco\clinical\nonsensitive\next_usable_id\', 'The next usable id for non sensitive clinical data.', '\medco\clinical\nonsensitive\next_usable_id\',
'NOW()', 'NOW()', 'NOW()', 'MEDCO_ADMIN', '@');

insert into i2b2metadata.clinical_sensitive (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum, c_basecode,
c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, update_date,
download_date, import_date, valuetype_cd, m_applied_path) values
('4', '\medco\clinical\sensitive\next_usable_id\1\', '1', 'N', 'LH', '0', 'MEDCO_ADMIN:nextEncUsableId', 'concept_cd', 'concept_dimension', 'concept_path',
'T', 'LIKE', '\medco\clinical\sensitive\next_usable_id\1\', 'The next usable id for Unlynx encryption.',
'NOW()', 'NOW()', 'NOW()', 'MEDCO_ADMIN', '@');
insert into i2b2metadata.clinical_non_sensitive (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum, c_basecode,
c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, update_date,
download_date, import_date, valuetype_cd, m_applied_path) values
('4', '\medco\clinical\nonsensitive\next_usable_id\1\', '1', 'N', 'LH', '0', 'MEDCO_ADMIN:nextClearUsableId', 'concept_cd', 'concept_dimension', 'concept_path',
'T', 'LIKE', '\medco\clinical\nonsensitive\next_usable_id\1\', 'The next usable id for non sensitive clinical data.',
'NOW()', 'NOW()', 'NOW()', 'MEDCO_ADMIN', '@');

-- todo: db lookup for medco?? yes in i2b2hive.ont_db_lookup add entry
-- setup i2b2-medco ontology todo: when setting up shrine, change to the schema of medco/shrine ontology !!
