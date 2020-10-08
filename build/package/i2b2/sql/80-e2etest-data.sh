#!/bin/bash
set -Eeuo pipefail
# set up data in the database for end-to-end tests

### description of the test data
# 4 patients: 1 (real), 2 (real), 3 (real), 4 (dummy)
# 3 concepts: 1, 2, 3
# observation_fact: p1: c1; p2: c1, c2; p3: c2, c3; p4: c1, c2, c3
# the same data is replicated on the 3 different nodes

psql $PSQL_PARAMS -d "$I2B2_DB_NAME" <<-EOSQL

    -- medco_ont.table_access
    insert into medco_ont.table_access (c_table_cd, c_table_name, c_protected_access, c_hlevel, c_fullname, c_name,
        c_synonym_cd, c_visualattributes, c_facttablecolumn, c_dimtablename,
        c_columnname, c_columndatatype, c_operator, c_dimcode, c_tooltip) VALUES
        ('E2ETEST', 'E2ETEST', 'N', '1', '\e2etest\', 'End-To-End Test',
        'N', 'CA', 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\e2etest\', 'End-To-End Test');

    -- medco_ont.e2etest
    CREATE TABLE medco_ont.e2etest(
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
    ALTER TABLE ONLY medco_ont.e2etest ADD CONSTRAINT fullname_pk_10 PRIMARY KEY (c_fullname);
    ALTER TABLE ONLY medco_ont.e2etest ADD CONSTRAINT basecode_un_10 UNIQUE (c_basecode);
    ALTER TABLE medco_ont.e2etest OWNER TO $I2B2_DB_USER;

    insert into medco_ont.e2etest
        (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
        c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator,
        c_dimcode, c_comment, c_tooltip, update_date, download_date, import_date,
        valuetype_cd, m_applied_path, c_basecode) values
            (
                '0', '\e2etest\', 'End-To-End Test', 'N', 'CA', '0',
                'concept_cd', 'concept_dimension', 'concept_path',
                'T', 'LIKE', '\e2etest\', 'End-To-End Test', '\e2etest\',
                'NOW()', 'NOW()', 'NOW()', 'ENC_ID', '@', ''
            ), (
                '1', '\e2etest\1\', 'E2E Concept 1', 'N', 'LA', '0',
                'concept_cd', 'concept_dimension', 'concept_path',
                'T', 'LIKE', '\e2etest\1\', 'E2E Concept 1', '\e2etest\1\',
                'NOW()', 'NOW()', 'NOW()', 'ENC_ID', '@', 'ENC_ID:1'
            ), (
                '1', '\e2etest\2\', 'E2E Concept 2', 'N', 'LA', '0',
                'concept_cd', 'concept_dimension', 'concept_path',
                'T', 'LIKE', '\e2etest\2\', 'E2E Concept 2', '\e2etest\2\',
                'NOW()', 'NOW()', 'NOW()', 'ENC_ID', '@', 'ENC_ID:2'
            ), (
                '1', '\e2etest\3\', 'E2E Concept 3', 'N', 'LA', '0',
                'concept_cd', 'concept_dimension', 'concept_path',
                'T', 'LIKE', '\e2etest\3\', 'E2E Concept 3', '\e2etest\3\',
                'NOW()', 'NOW()', 'NOW()', 'ENC_ID', '@', 'ENC_ID:3'
            ), (
                '1', '\modifiers\', 'E2E Modifiers test', 'N', 'DA', '0',
                'modifier_cd', 'modifier_dimension', 'modifier_path',
                'T', 'LIKE', '\modifiers\', 'End-To-End Modifiers Test', '\modifiers\',
                'NOW()', 'NOW()', 'NOW()', 'ENC_ID', '\e2etest\%', 'ENC_ID:4'
            ), (
                '2', '\modifiers\1\', 'E2E Modifier 1', 'N', 'RA', '0',
                'modifier_cd', 'modifier_dimension', 'modifier_path',
                'T', 'LIKE', '\modifiers\1\', 'E2E Modifier 1', '\modifiers\1\',
                'NOW()', 'NOW()', 'NOW()', 'ENC_ID', '\e2etest\1\', 'ENC_ID:5'
            ), (
                '2', '\modifiers\2\', 'E2E Modifier 2', 'N', 'RA', '0',
                'modifier_cd', 'modifier_dimension', 'modifier_path',
                'T', 'LIKE', '\modifiers\2\', 'E2E Modifier 2', '\modifiers\2\',
                'NOW()', 'NOW()', 'NOW()', 'ENC_ID', '\e2etest\2\', 'ENC_ID:6'
            ), (
                '2', '\modifiers\3\', 'E2E Modifier 3', 'N', 'RA', '0',
                'modifier_cd', 'modifier_dimension', 'modifier_path',
                'T', 'LIKE', '\modifiers\3\', 'E2E Modifier 3', '\modifiers\3\',
                'NOW()', 'NOW()', 'NOW()', 'ENC_ID', '\e2etest\3\', 'ENC_ID:7'
            );

    -- medco_ont.sensitive_tagged
    insert into medco_ont.sensitive_tagged
        (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
        c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator,
        c_dimcode, c_comment, c_tooltip, update_date, download_date, import_date,
        valuetype_cd, m_applied_path, c_basecode) values
            (
                '1', '\medco\tagged\', 'MedCo Sensitive Tagged Ontology', 'N', 'CA', '0',
                'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE',
                '\medco\tagged\', 'MedCo Sensitive Tagged Ontology', '\medco\tagged\', 'NOW()', 'NOW()', 'NOW()',
                'TAG_ID', '@', ''
            ), (
                '2', '\medco\tagged\8d3533369426ae172271e98cef8be2bbfe9919087c776083b1ea1de803fc87aa\', '', 'N', 'LA', '0',
                'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE',
                '\medco\tagged\8d3533369426ae172271e98cef8be2bbfe9919087c776083b1ea1de803fc87aa\', '', '', 'NOW()', 'NOW()', 'NOW()',
                'TAG_ID', '@', 'TAG_ID:11'
            ), (
                '2', '\medco\tagged\c75af24ed416c61b67011eb91aa852f5069c020c4bd8c1e64a07c7fb061d8ace\', '', 'N', 'LA', '0',
                'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE',
                '\medco\tagged\c75af24ed416c61b67011eb91aa852f5069c020c4bd8c1e64a07c7fb061d8ace\', '', '', 'NOW()', 'NOW()', 'NOW()',
                'TAG_ID', '@', 'TAG_ID:12'
            ), (
                '2', '\medco\tagged\46f75970444851cb64f1b940ef9205a20b65d6fd5bb68a250b0f52f07f6da9a3\', '', 'N', 'LA', '0',
                'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE',
                '\medco\tagged\46f75970444851cb64f1b940ef9205a20b65d6fd5bb68a250b0f52f07f6da9a3\', '', '', 'NOW()', 'NOW()', 'NOW()',
                'TAG_ID', '@', 'TAG_ID:13'
            ), (
                '2', '\medco\tagged\6522b1ad2d46cf1e44fb80a9b32aefb0b93d91e612b35ea45e8074ac49a4c714\', '', 'N', 'LA', '0',
                'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE',
                '\medco\tagged\6522b1ad2d46cf1e44fb80a9b32aefb0b93d91e612b35ea45e8074ac49a4c714\', '', '', 'NOW()', 'NOW()', 'NOW()',
                'TAG_ID', '@', 'TAG_ID:11'
            ), (
                '2', '\medco\tagged\14afc3c7eabd32bd188cea384276257cc352ce217bb67eccf19572b8527b4525\', '', 'N', 'LA', '0',
                'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE',
                '\medco\tagged\14afc3c7eabd32bd188cea384276257cc352ce217bb67eccf19572b8527b4525\', '', '', 'NOW()', 'NOW()', 'NOW()',
                'TAG_ID', '@', 'TAG_ID:12'
            ), (
                '2', '\medco\tagged\09bc15e0d90046c102199f1b4d20eef9ee91b2ea3fd4608303775d000dd1248c\', '', 'N', 'LA', '0',
                'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE',
                '\medco\tagged\09bc15e0d90046c102199f1b4d20eef9ee91b2ea3fd4608303775d000dd1248c\', '', '', 'NOW()', 'NOW()', 'NOW()',
                'TAG_ID', '@', 'TAG_ID:13'
            ), (
                '2', '\medco\tagged\15a9a62d50d8239a1a133544403980fc97468b8479917f96b24373cdf4397e11\', '', 'N', 'LA', '0',
                'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE',
                '\medco\tagged\15a9a62d50d8239a1a133544403980fc97468b8479917f96b24373cdf4397e11\', '', '', 'NOW()', 'NOW()', 'NOW()',
                'TAG_ID', '@', 'TAG_ID:11'
            ), (
                '2', '\medco\tagged\a757856f859fa5f82f5164b33459e89a1a84a213cdb11a3e3bc7df460a495b3e\', '', 'N', 'LA', '0',
                'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE',
                '\medco\tagged\a757856f859fa5f82f5164b33459e89a1a84a213cdb11a3e3bc7df460a495b3e\', '', '', 'NOW()', 'NOW()', 'NOW()',
                'TAG_ID', '@', 'TAG_ID:12'
            ), (
                '2', '\medco\tagged\84ff65ad621ebeba9b7ef1c68967ae023cb0487415a0b2061baecefbd0da67ba\', '', 'N', 'LA', '0',
                'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE',
                '\medco\tagged\84ff65ad621ebeba9b7ef1c68967ae023cb0487415a0b2061baecefbd0da67ba\', '', '', 'NOW()', 'NOW()', 'NOW()',
                'TAG_ID', '@', 'TAG_ID:13'
            );

    -- i2b2demodata_i2b2.concept_dimension
    insert into i2b2demodata_i2b2.concept_dimension
        (concept_path, concept_cd, import_date, upload_id) values
            ('\medco\tagged\8d3533369426ae172271e98cef8be2bbfe9919087c776083b1ea1de803fc87aa\', 'TAG_ID:11', 'NOW()', '1'),
            ('\medco\tagged\c75af24ed416c61b67011eb91aa852f5069c020c4bd8c1e64a07c7fb061d8ace\', 'TAG_ID:12', 'NOW()', '1'),
            ('\medco\tagged\46f75970444851cb64f1b940ef9205a20b65d6fd5bb68a250b0f52f07f6da9a3\', 'TAG_ID:13', 'NOW()', '1'),
            ('\medco\tagged\6522b1ad2d46cf1e44fb80a9b32aefb0b93d91e612b35ea45e8074ac49a4c714\', 'TAG_ID:11', 'NOW()', '1'),
            ('\medco\tagged\14afc3c7eabd32bd188cea384276257cc352ce217bb67eccf19572b8527b4525\', 'TAG_ID:12', 'NOW()', '1'),
            ('\medco\tagged\09bc15e0d90046c102199f1b4d20eef9ee91b2ea3fd4608303775d000dd1248c\', 'TAG_ID:13', 'NOW()', '1'),
            ('\medco\tagged\15a9a62d50d8239a1a133544403980fc97468b8479917f96b24373cdf4397e11\', 'TAG_ID:11', 'NOW()', '1'),
            ('\medco\tagged\a757856f859fa5f82f5164b33459e89a1a84a213cdb11a3e3bc7df460a495b3e\', 'TAG_ID:12', 'NOW()', '1'),
            ('\medco\tagged\84ff65ad621ebeba9b7ef1c68967ae023cb0487415a0b2061baecefbd0da67ba\', 'TAG_ID:13', 'NOW()', '1');

    -- i2b2demodata_i2b2.modifier_dimension
    insert into i2b2demodata_i2b2.modifier_dimension
        (modifier_path, modifier_cd, import_date, upload_id) values
            ('\modifiers\', 'ENC_ID:4', 'NOW()', '1'),
            ('\modifiers\1\', 'ENC_ID:5', 'NOW()', '1'),
            ('\modifiers\2\', 'ENC_ID:6', 'NOW()', '1'),
            ('\modifiers\3\', 'ENC_ID:7', 'NOW()', '1');

    -- i2b2demodata_i2b2.provider_dimension
    insert into i2b2demodata_i2b2.provider_dimension
        (provider_id, provider_path, name_char, import_date, upload_id) values
            ('e2etest', '\e2etest\', 'e2etest', 'NOW()', '1');

    -- i2b2demodata_i2b2.patient_dimension
    insert into i2b2demodata_i2b2.patient_dimension
        (patient_num, import_date, upload_id, enc_dummy_flag_cd) values
            ('1', 'NOW()', '1', '3xTFQjXiSLkWvczWppUlGHUXo8ibioGiLg6jy7pWC8NkOI07y9wD-637Smtb1bpVh2shgb9WgSr_eJN0u3WJ-A=='),
            ('2', 'NOW()', '1', '6xQAq3jj9OVb0E5Tce44TWKGZzOXiAb6xoxxiEF3ZUzyVzGEeo7u3iDx9uLIQtamrAypL0MbTyCTqFdtGmYWGg=='),
            ('3', 'NOW()', '1', 'odD--LyYEQkOg4iZvQw_S0evR2G1-5YNF8kD8tLzVAZjP7XItD678vr2ze-DZxY0AruhagqMa45zD_AXi3m3Gg=='),
            ('4', 'NOW()', '1', 'pBcUzT4V-sawblUZhWpT3xYCsNetCLeEspUqpChabxxaI2FQNSsgBzkv1xBi2VCwYaRPvf5fvJO_YpjTl8pZ0w==');

    -- i2b2demodata_i2b2.patient_mapping
    insert into i2b2demodata_i2b2.patient_mapping
        (patient_ide, patient_ide_source, patient_num, project_id, import_date, upload_id) values
            ('e2etest1', 'e2etest', '1', 'MedCo', 'NOW()', '1'),
            ('e2etest2', 'e2etest', '2', 'MedCo', 'NOW()', '1'),
            ('e2etest3', 'e2etest', '3', 'MedCo', 'NOW()', '1'),
            ('e2etest4', 'e2etest', '4', 'MedCo', 'NOW()', '1');

    -- i2b2demodata_i2b2.visit_dimension
    insert into i2b2demodata_i2b2.visit_dimension
        (encounter_num, patient_num, import_date, upload_id) values
            ('1', '1', 'NOW()', '1'),
            ('2', '2', 'NOW()', '1'),
            ('3', '3', 'NOW()', '1'),
            ('4', '4', 'NOW()', '1');

    -- i2b2demodata_i2b2.encounter_mapping
    insert into i2b2demodata_i2b2.encounter_mapping
        (encounter_ide, encounter_ide_source, project_id, encounter_num, patient_ide, patient_ide_source, import_date, upload_id) values
            ('e2etest1', 'e2etest', 'MedCo', '1', 'e2etest1', 'e2etest', 'NOW()', '1'),
            ('e2etest2', 'e2etest', 'MedCo', '2', 'e2etest2', 'e2etest', 'NOW()', '1'),
            ('e2etest3', 'e2etest', 'MedCo', '3', 'e2etest3', 'e2etest', 'NOW()', '1'),
            ('e2etest4', 'e2etest', 'MedCo', '4', 'e2etest4', 'e2etest', 'NOW()', '1');

    -- i2b2demodata_i2b2.observation_fact
    insert into i2b2demodata_i2b2.observation_fact
        (encounter_num, patient_num, concept_cd, provider_id, start_date, modifier_cd, instance_num, import_date, upload_id) values
            ('1', '1', 'TAG_ID:11', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1'),
            ('2', '2', 'TAG_ID:11', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1'),
            ('2', '2', 'TAG_ID:12', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1'),
            ('3', '3', 'TAG_ID:12', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1'),
            ('3', '3', 'TAG_ID:13', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1'),
            ('4', '4', 'TAG_ID:11', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1'),
            ('4', '4', 'TAG_ID:12', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1'),
            ('4', '4', 'TAG_ID:13', 'e2etest', 'NOW()', '@', '1', 'NOW()', '1');
EOSQL
