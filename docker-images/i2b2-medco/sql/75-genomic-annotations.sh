#!/bin/bash
set -Eeuo pipefail
# set up common medco ontology

psql $PSQL_PARAMS -d "$I2B2_DB_NAME" <<-EOSQL

    -- genomic entries in i2b2/medco ontology
    insert into medco_ont.genomic (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
        c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
        download_date, import_date, valuetype_cd, m_applied_path) values
        ('2', '\medco\genomic\annotations_Hugo_Symbol\', 'Gene Name', 'N', 'LA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
        'T', 'LIKE', '\medco\genomic\annotations_Hugo_Symbol\', 'Gene Name', '\medco\genomic\annotations_Hugo_Symbol\',
        'NOW()', 'NOW()', 'NOW()', 'GEN', '@');
    insert into medco_ont.genomic (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
        c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
        download_date, import_date, valuetype_cd, m_applied_path) values
        ('2', '\medco\genomic\annotations_Protein_position\', 'Protein Position', 'N', 'LA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
        'T', 'LIKE', '\medco\genomic\annotations_Protein_position\', 'Protein Position', '\medco\genomic\annotations_Protein_position\',
        'NOW()', 'NOW()', 'NOW()', 'GEN', '@');
    insert into medco_ont.genomic (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
        c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
        download_date, import_date, valuetype_cd, m_applied_path) values
        ('2', '\medco\genomic\variant\', 'Variant Name', 'N', 'LA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
        'T', 'LIKE', '\medco\genomic\variant\', 'Variant Name', '\medco\genomic\variant\',
        'NOW()', 'NOW()', 'NOW()', 'GEN', '@');

    -- genomic tables
    create table genomic_annotations.genomic_annotations(
        variant_id character varying(255) NOT NULL,
        variant_name character varying(255) NOT NULL,
        protein_change character varying(255) NOT NULL,
        hugo_gene_symbol character varying(255) NOT NULL,
        annotations text NOT NULL
    );
    create table genomic_annotations.annotation_names(
        annotation_name character varying(255) NOT NULL PRIMARY KEY
    );
    create table genomic_annotations.gene_values(
        gene_value character varying(255) NOT NULL PRIMARY KEY
    );

    -- permissions
    ALTER TABLE genomic_annotations.genomic_annotations OWNER TO $I2B2_DB_USER;
    ALTER TABLE genomic_annotations.annotation_names OWNER TO $I2B2_DB_USER;
    ALTER TABLE genomic_annotations.gene_values OWNER TO $I2B2_DB_USER;
    grant all on schema genomic_annotations to $I2B2_DB_USER;
    grant all privileges on all tables in schema genomic_annotations to $I2B2_DB_USER;
EOSQL
