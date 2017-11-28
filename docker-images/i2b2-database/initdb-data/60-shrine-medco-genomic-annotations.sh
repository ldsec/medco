#!/bin/bash
set -e
# load the genomic annotations schema

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$I2B2_MEDCO_DB_NAME" <<-EOSQL

    -- genomic entries in shrine/medco ontology
    insert into shrine_ont.genomic (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
        c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
        download_date, import_date, valuetype_cd, m_applied_path) values
        ('2', '\medco\genomic\annotations_Hugo_Symbol\', 'Gene Name', 'N', 'LA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
        'T', 'LIKE', '\medco\genomic\annotations_Hugo_Symbol\', 'Gene Name', '\medco\genomic\annotations_Hugo_Symbol\',
        'NOW()', 'NOW()', 'NOW()', 'GEN', '@');
    insert into shrine_ont.genomic (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
        c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
        download_date, import_date, valuetype_cd, m_applied_path) values
        ('2', '\medco\genomic\annotations_Protein_position\', 'Protein Position', 'N', 'LA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
        'T', 'LIKE', '\medco\genomic\annotations_Protein_position\', 'Protein Position', '\medco\genomic\annotations_Protein_position\',
        'NOW()', 'NOW()', 'NOW()', 'GEN', '@');
    insert into shrine_ont.genomic (c_hlevel, c_fullname, c_name, c_synonym_cd, c_visualattributes, c_totalnum,
        c_facttablecolumn, c_tablename, c_columnname, c_columndatatype, c_operator, c_dimcode, c_comment, c_tooltip, update_date,
        download_date, import_date, valuetype_cd, m_applied_path) values
        ('2', '\medco\genomic\variant\', 'Variant Name', 'N', 'LA', '0', 'concept_cd', 'concept_dimension', 'concept_path',
        'T', 'LIKE', '\medco\genomic\variant\', 'Variant Name', '\medco\genomic\variant\',
        'NOW()', 'NOW()', 'NOW()', 'GEN', '@');

    -- genomic tables
    create table genomic_annotations.genomic_annotations(
        variant_id character varying(255) NOT NULL PRIMARY KEY,
        variant_name character varying(255) NOT NULL,
        annotations text NOT NULL,
        t_depth numeric NOT NULL
    );
    create table genomic_annotations.annotation_names(
        annotation_name character varying(255) NOT NULL PRIMARY KEY
    );
    create table genomic_annotations.gene_values(
        gene_value character varying(255) NOT NULL PRIMARY KEY
    );

    -- permissions
    ALTER TABLE genomic_annotations.genomic_annotations OWNER TO genomic_annotations;
    ALTER TABLE genomic_annotations.annotation_names OWNER TO genomic_annotations;
    ALTER TABLE genomic_annotations.gene_values OWNER TO genomic_annotations;
    grant all on schema genomic_annotations to genomic_annotations;
    grant all privileges on all tables in schema genomic_annotations to genomic_annotations;
EOSQL
