#!/bin/bash
set -Eeuo pipefail
# set up data in the database for end-to-end tests

psql $PSQL_PARAMS -d "$MC_DB_NAME" <<-EOSQL

    CREATE TABLE genomic_annotations.e2etest(
				variant_id character varying(255) NOT NULL,
				variant_id_enc character varying(255) NOT NULL,
				variant_name character varying(255) NOT NULL,
				protein_change character varying(255) NOT NULL,
				hugo_gene_symbol character varying(255) NOT NULL,
				annotations text NOT NULL);

    ALTER TABLE ONLY genomic_annotations.e2etest ADD CONSTRAINT variant_id_pk_10 PRIMARY KEY (variant_id);
    ALTER TABLE genomic_annotations.e2etest OWNER TO $MC_DB_USER;

    INSERT INTO genomic_annotations.e2etest
    (variant_id, variant_id_enc, variant_name, protein_change, hugo_gene_symbol, annotations) VALUES
        (
            'vID1', 'enc(vID1)', 'vn1', 'pc1', 'hgs1', '1111'
        ), (
            'vID2', 'enc(vID2)', 'vn2', 'pc2', 'hgs2', '2222'
        ), (
            'vID3', 'enc(vID3)', 'vn3', 'pc3', 'hgs3', '3333'
        ), (
            'vID4', 'enc(vID4)', 'vn4', 'pc4', 'hgs4', '4444'
        ), (
            'vID5', 'enc(vID5)', 'vn5', 'pc5', 'hgs5', '5555'
        );

EOSQL