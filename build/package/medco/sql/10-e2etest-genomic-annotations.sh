#!/bin/bash
set -Eeuo pipefail
# set up the genomic annotations database

psql $PSQL_PARAMS -d "$MC_DB_NAME" <<-EOSQL
BEGIN;

  -- base structure
  CREATE SCHEMA genomic_annotations;
  grant all on schema genomic_annotations to $MC_DB_USER;
  grant all privileges on all tables in schema genomic_annotations to $MC_DB_USER;

	CREATE TABLE IF NOT EXISTS genomic_annotations.genomic_annotations(
    variant_id character varying(255) NOT NULL,
    variant_id_enc character varying(255) NOT NULL,
    variant_name character varying(255) NOT NULL,
    protein_change character varying(255) NOT NULL,
    hugo_gene_symbol character varying(255) NOT NULL,
    annotations text NOT NULL
  );

  CREATE TABLE IF NOT EXISTS genomic_annotations.annotation_names(
    annotation_name character varying(255) NOT NULL PRIMARY KEY
  );

  CREATE TABLE IF NOT EXISTS genomic_annotations.gene_values(
    gene_value character varying(255) NOT NULL PRIMARY KEY
  );

	-- plpgsql functions
  CREATE OR REPLACE FUNCTION genomic_annotations.ga_getvalues(annotation varchar, val varchar, lim int) RETURNS SETOF varchar AS \$\$
  BEGIN
    RETURN QUERY EXECUTE format(
      'SELECT annotation_value FROM genomic_annotations.%I
        WHERE annotation_value ~* \$1
        ORDER BY annotation_value LIMIT \$2',
      annotation
    )
    USING val, lim;
  END;
  \$\$ LANGUAGE plpgsql;

  CREATE OR REPLACE FUNCTION genomic_annotations.ga_getvariants(annotation varchar, val varchar, zygosity varchar, enc bool) RETURNS SETOF varchar AS \$\$
	DECLARE col varchar;
  BEGIN
    IF enc
    THEN
      col := 'variant_id_enc';
    ELSE
      col := 'variant_id';
    END IF;

    RETURN QUERY EXECUTE format(
      'SELECT %I FROM genomic_annotations.genomic_annotations
        WHERE lower(%I) = lower(\$1)
        AND annotations ~* \$2 ORDER BY variant_id',
      col,
      annotation
    )
    USING val, zygosity;
  END;
  \$\$ LANGUAGE plpgsql;

  CREATE OR REPLACE FUNCTION genomic_annotations.ga_annotationexists(annotation varchar) RETURNS boolean AS \$\$
  BEGIN
    RETURN EXISTS(
      SELECT 1 FROM pg_tables where
      schemaname = 'genomic_annotations' and
      tablename = annotation
    );
  END;
  \$\$ LANGUAGE plpgsql;

  -- genomic annotations data
COPY genomic_annotations.genomic_annotations (variant_id, variant_id_enc, variant_name, protein_change, hugo_gene_symbol, annotations) FROM stdin;
-2823470849823937376	VZO4Wk9OaiNixXrksO8dTNipzOi7gAYc_p-20OqNbY6XbCeh95RyhRqkexaDLSy-neOX71afsXo_xw9ZjjgIXA==	22:54792418:TTAGTT>CAGGAA	N232Y	BLK	Heterozygous;Cancer Study Id=209;Cancer Study Identifier=luad_tcga;Tumor Alt Count=87;Tumor Ref Count=939;Normal Alt Count=-26;Normal Ref Count=1374;End Position=74019605;Mutation Type=Nonstop_Mutation;Ncbi Build=;Strand=;Variant Type=DNP;Entrez Gene Id=126
-8874641984867408336	3sMH3Gh9eH2G-bPYfSTO5W5tJd7OP72H8_vtY-SuY0al5pmDxMQSTTZAeaBgx7y-u5Hcv1QCJIbE3p7gVzo9CA==	1:56344713:CGTGGC>TGACAA	P1461L	NBPF15	Heterozygous;Cancer Study Id=210;Cancer Study Identifier=lusc_tcga;Tumor Alt Count=486;Tumor Ref Count=1029;Normal Alt Count=36;Normal Ref Count=569;End Position=124401152;Mutation Type=Targeted_Region;Ncbi Build=37;Strand=;Variant Type=TNP;Entrez Gene Id=72954
-8223566780127745187	unlCj3QUlPP9XANU02y16bwcWLcK611TsUOm38cgi78lF9_o4-1vs6tDCoXKGTax8wNZtVGEGRHFQ7onmkynlA==	3:125834837:TCGTAA>GCTTCT	K438=	TPRXL	Heterozygous;Cancer Study Id=208;Cancer Study Identifier=lihc_tcga;Tumor Alt Count=373;Tumor Ref Count=1523;Normal Alt Count=29;Normal Ref Count=288;End Position=31401569;Mutation Type=In_Frame_Ins;Ncbi Build=;Strand=1;Variant Type=DNP;Entrez Gene Id=5654
-7455563962931223533	vtuBCedKEoXoGOs0u_s63vVbNd6fgv04MSzACrEjRGBdCTB_Vq-Hih7oU3clqO-1s_070C7BkfGkvyzvB74Iuw==	6:35786830:GGGACC>TAATAC	F137V	CCDC140	Heterozygous;Cancer Study Id=215;Cancer Study Identifier=skcm_tcga;Tumor Alt Count=471;Tumor Ref Count=216;Normal Alt Count=99;Normal Ref Count=2056;End Position=50659412;Mutation Type=Nonstop_Mutation;Ncbi Build=hg19;Strand=1;Variant Type=TNP;Entrez Gene Id=7907
-5164302573854136353	3DOcYug-U-0dIi9tT3f8DN8DanA0nOCZXMonds0QY_VPdFeDzRPSrISR9LqnFLOcEW4tweLwiCSOYAlyrQhqSw==	14:22206638:GTGTAC>TCCTCC	E482*	IL3	Heterozygous;Cancer Study Id=201;Cancer Study Identifier=hnsc_tcga;Tumor Alt Count=462;Tumor Ref Count=1462;Normal Alt Count=26;Normal Ref Count=490;End Position=115936860;Mutation Type=Splice_Site;Ncbi Build=GRCh37;Strand=nu;Variant Type=TNP;Entrez Gene Id=105
-6148177330222111359	nUkkI_y43jANpuJocFw2gSsGm9PWSAVGtrrDJ4rOuajDpwK_9gWtvRjYsbnmzPk12yKcpIZiOGAOTJguucw8Dw==	10:179643691:AAAAA>TTGAAT	L447R	APRT	Heterozygous;Cancer Study Id=227;Cancer Study Identifier=ov_tcga;Tumor Alt Count=305;Tumor Ref Count=1048;Normal Alt Count=-76;Normal Ref Count=1267;End Position=12325306;Mutation Type=Frame_Shift_Del;Ncbi Build=37;Strand=nu;Variant Type=;Entrez Gene Id=2897
-7158674830913413776	qHmWosVvojlyf6gYQ_ug83FrigDSAZzHk956Xy1HJGaTeYrh871JKCg9iMd6tq-2pUV_y7NRoF6Tfvei0AlYIQ==	7:43850925:ATCAA>CTTCA	M1527I	MIR1283-2	Heterozygous;Cancer Study Id=216;Cancer Study Identifier=stad_tcga;Tumor Alt Count=396;Tumor Ref Count=323;Normal Alt Count=-32;Normal Ref Count=90;End Position=107416898;Mutation Type=Nonstop_Mutation;Ncbi Build=37;Strand=nu;Variant Type=TNP;Entrez Gene Id=7908
-8233372561981092240	etMFhskSrROZSa7ghanl4vcmzQ-eCLryGzpV_gt5wJIEslEzm6P1UKZJAscLhXGoi8zXfwsBkScobu4Mz-ccZg==	3:116702491:TCTTTG>CGTC	A2390T	LOC284009	Heterozygous;Cancer Study Id=214;Cancer Study Identifier=sarc_tcga;Tumor Alt Count=393;Tumor Ref Count=1152;Normal Alt Count=2;Normal Ref Count=1398;End Position=110086737;Mutation Type=Frame_Shift_Del;Ncbi Build=37;Strand=nu;Variant Type=DNP;Entrez Gene Id=2350
-2821912056169469430	cxQMXY4sCmUCixVKohrDoNPDOAL2NwAdHyzJTA0LmoQtxBfJJFVSSnAtFq2o_J6sKzQF10sfvVUMRCSMCNt3GQ==	22:56244158:ATTA>GGAAGG	T1003M	OPA1-AS1	Heterozygous;Cancer Study Id=210;Cancer Study Identifier=lusc_tcga;Tumor Alt Count=39;Tumor Ref Count=30;Normal Alt Count=-6;Normal Ref Count=171;End Position=12396889;Mutation Type=Translation_Start_Site;Ncbi Build=37;Strand=nu;Variant Type=INS;Entrez Gene Id=2623
-7121901993980174104	hGyT6sji86pEK_SlhkCUojNa8Qq3YgW8I6rxr7_jUf3nd5bXvSa2RZ_NuT_oq2cfOkjsOScDx18nPJkXrb1BtA==	7:78098298:TCTTTA>AACGGA	E1037A	NAV3	Heterozygous;Cancer Study Id=219;Cancer Study Identifier=ucec_tcga;Tumor Alt Count=388;Tumor Ref Count=673;Normal Alt Count=-71;Normal Ref Count=477;End Position=136071142;Mutation Type=Missense_Mutation;Ncbi Build=37;Strand=;Variant Type=ONP;Entrez Gene Id=975
-4898572880864589696	6FWPKxuwXMU9T2jrKDJA1B39PYeFTYtJ9lLFSCLIny3_ljdUHE9zLndDcK6uZbu8I4BgASRYXRelyNwanwzWsA==	7:78098298:TCTTTA>AACGGA	E1037A	NAV3	Heterozygous;Cancer Study Id=219;Cancer Study Identifier=ucec_tcga;Tumor Alt Count=388;Tumor Ref Count=673;Normal Alt Count=-71;Normal Ref Count=477;End Position=136071142;Mutation Type=Missense_Mutation;Ncbi Build=37;Strand=;Variant Type=ONP;Entrez Gene Id=975
-6271408487767448598	DNBlQque-I3m-r74qhqkB7ejw9Z5gyErsaPz2_ifMCGV73ZhJZxSirQNCgyI4BYOJzT-g8cvrawExPu59EMNcQ==	7:78098298:TCTTTA>AACGGA	E1037A	NAV3	Unknown;Cancer Study Id=219;Cancer Study Identifier=ucec_tcga;Tumor Alt Count=388;Tumor Ref Count=673;Normal Alt Count=-71;Normal Ref Count=477;End Position=136071142;Mutation Type=Missense_Mutation;Ncbi Build=37;Strand=;Variant Type=ONP;Entrez Gene Id=975
-2541144436738335408	vEdzkDO-EYU0NVdAprG83EbreJP0l9Tbz8AkOTUqiKnPMmmw8wZnQ0h0Xjr76GNwKycAODiy2FhthLyuQXju6w==	23:49293924:GGGAA>CTTTA	R675K	FA2H	Heterozygous;Cancer Study Id=219;Cancer Study Identifier=ucec_tcga;Tumor Alt Count=326;Tumor Ref Count=321;Normal Alt Count=-29;Normal Ref Count=899;End Position=33388070;Mutation Type=Translation_Start_Site;Ncbi Build=37;Strand=+;Variant Type=TNP;Entrez Gene Id=407
-8644437442588316507	Ce_39lNbnyZiNPHFE9sJIBnUIOMO8OZ3ZulQKvxgcVJgt8IT4t6g5pStcXjXj6bB3vanf_CzC1YX9WMcX6UbfA==	2:2303944:CGCTCC>CAGGTT	V1310I	FAM24A	Heterozygous;Cancer Study Id=219;Cancer Study Identifier=ucec_tcga;Tumor Alt Count=206;Tumor Ref Count=257;Normal Alt Count=13;Normal Ref Count=809;End Position=68938097;Mutation Type=Missense_Mutation;Ncbi Build=GRCh37;Strand=;Variant Type=DEL;Entrez Gene Id=16553
-8884062052056279280	nfMYwgiJWqz7Q8EcbIDl14NjuZW6tfEPEmXkZbJKY5nzz1UIyNFMNeblat4nWBN93uhlekmW04APuMdM_etQwA==	1:47571592:ATTTA>TCATA	R5689W	C17ORF75	Heterozygous;Cancer Study Id=216;Cancer Study Identifier=stad_tcga;Tumor Alt Count=219;Tumor Ref Count=1078;Normal Alt Count=8;Normal Ref Count=1434;End Position=21414799;Mutation Type=In_Frame_Ins;Ncbi Build=hg19;Strand=+;Variant Type=SNP;Entrez Gene Id=39027
-2829226452809038848	u4k8LJBTopSXsdx8BFQ_PtMkgrAIYPmKtrg-LZ2pDvHtfbMKdQI-5Ze4lt8XP7J1jV88cguBsPG2uWHKyrBbOg==	22:49432095:TGTTGC>TAAA	G578S	CHST7	Heterozygous;Cancer Study Id=205;Cancer Study Identifier=kirc_tcga;Tumor Alt Count=23;Tumor Ref Count=1158;Normal Alt Count=43;Normal Ref Count=1210;End Position=40760969;Mutation Type=Targeted_Region;Ncbi Build=;Strand=-1;Variant Type=DEL;Entrez Gene Id=8465
-6855388190066128176	no8Io_fUhCQ9lu0DEJt51bnUYxBg4Hhg6ZhYeEkPkQpQsXH6oVgqLi8bhvqnYD3AiFJfY93PLo5WU46ONlVUMg==	8:57873164:GCTGTG>GGCT	Q1032R	ATRAID	Heterozygous;Cancer Study Id=190;Cancer Study Identifier=coadread_tcga;Tumor Alt Count=224;Tumor Ref Count=816;Normal Alt Count=68;Normal Ref Count=169;End Position=24533482;Mutation Type=Targeted_Region;Ncbi Build=GRCh37;Strand=-1;Variant Type=DEL;Entrez Gene Id=5
-4856934075737223744	lP-lKSo8pzkL70QUD3t2a7ajCHdNfGzWzWgmHaVx8YV_qNNCgyMZALIzdZKboenLtaRaV7DArcTQ8bFzRQEC9A==	15:40030403:CATGAG>CTCA	G824R	TMEM75	Heterozygous;Cancer Study Id=215;Cancer Study Identifier=skcm_tcga;Tumor Alt Count=328;Tumor Ref Count=312;Normal Alt Count=67;Normal Ref Count=28;End Position=100668298;Mutation Type=Nonsense_Mutation;Ncbi Build=GRCh37;Strand=;Variant Type=DNP;Entrez Gene Id=49386
-2835277828882931569	J5f0cf3BmiAc9EnP5lIakR0wLFbLHKNMdi4dpk6bCqOX7-Nm0vN7IsW87nKvtXvn4NDPeDkW1GwBpZl_xDao_w==	22:43796312:ACCAG>AAGACC	H277Y	IGHMBP2	Heterozygous;Cancer Study Id=216;Cancer Study Identifier=stad_tcga;Tumor Alt Count=285;Tumor Ref Count=1124;Normal Alt Count=2;Normal Ref Count=246;End Position=34130197;Mutation Type=Translation_Start_Site;Ncbi Build=37;Strand=nu;Variant Type=SNP;Entrez Gene Id=928
-4898572880864589696	6FWPKxuwXMU9T2jrKDJA1B39PYeFTYtJ9lLFSCLIny3_ljdUHE9zLndDcK6uZbu8I4BgASRYXRelyNwanwzWsA==	15:1251244:ACTGG>GAGA	R53S	HEBP2	Heterozygous;Cancer Study Id=208;Cancer Study Identifier=lihc_tcga;Tumor Alt Count=344;Tumor Ref Count=1523;Normal Alt Count=27;Normal Ref Count=185;End Position=52217174;Mutation Type=Splice_Region;Ncbi Build=;Strand=1;Variant Type=DEL;Entrez Gene Id=40081
-6585443479478848512	a-mvne2elFTsv-qMfcMYcLd6VAAtPfOGrycoSgXDHemJeBEY5wUc3cOdbBbtYGklRdaY-XLmWzpYIMJI5UMQyw==	9:40843311:GGGAA>TAA	A4667E	CD3D	Heterozygous;Cancer Study Id=219;Cancer Study Identifier=ucec_tcga;Tumor Alt Count=695;Tumor Ref Count=898;Normal Alt Count=-120;Normal Ref Count=673;End Position=62362777;Mutation Type=Frame_Shift_Del;Ncbi Build=hg19;Strand=+;Variant Type=DNP;Entrez Gene Id=17026
-5414149846471650656	iVOc5hN-_Y91wQJk4TKTzgYTKA9iO8-1td1BmZkw_7Oe3Y8mc4u_BFC0kjrMMtWPwH10Pwd3j68pkk6HKqKe7g==	13:57953689:AAGGT>AGGGA	L61Sfs*54	AANAT	Heterozygous;Cancer Study Id=192;Cancer Study Identifier=dlbc_tcga;Tumor Alt Count=286;Tumor Ref Count=1553;Normal Alt Count=-75;Normal Ref Count=1327;End Position=3341899;Mutation Type=Targeted_Region;Ncbi Build=;Strand=nu;Variant Type=DNP;Entrez Gene Id=28430
-4677558035137305897	YFPvf1Ga-JUNVf3oYI2Hu66flYK5-68U1DG4Jx7POp8nnARDsb8oE3nI5vYNrsJSocdL_YlY8A1QZ41pnd6_aQ==	15:207087359:TGGCGC>AGCTTC	D4283N	LOC285556	Heterozygous;Cancer Study Id=215;Cancer Study Identifier=skcm_tcga;Tumor Alt Count=195;Tumor Ref Count=630;Normal Alt Count=204;Normal Ref Count=380;End Position=1881297;Mutation Type=In_Frame_Ins;Ncbi Build=37;Strand=;Variant Type=DEL;Entrez Gene Id=49411
-8339306202921406714	M6lbAkv7qveo0PBP4bw3nDoatfXNRhgcZ9QpdQNTSoEl55vFlhkDl2_-I_IC7sO__2XMM9bfjw1E4GiSEmFv3w==	3:18044100:GCTT>CCAATG	Y1062H	TRIM8	Heterozygous;Cancer Study Id=218;Cancer Study Identifier=thca_tcga;Tumor Alt Count=63;Tumor Ref Count=511;Normal Alt Count=-48;Normal Ref Count=30;End Position=136691503;Mutation Type=In_Frame_Ins;Ncbi Build=hg19;Strand=+;Variant Type=SNP;Entrez Gene Id=34668
-6271408487767448598	DNBlQque-I3m-r74qhqkB7ejw9Z5gyErsaPz2_ifMCGV73ZhJZxSirQNCgyI4BYOJzT-g8cvrawExPu59EMNcQ==	10:64875732:GGGAA>ACCGGG	F377S	JDP2	Heterozygous;Cancer Study Id=215;Cancer Study Identifier=skcm_tcga;Tumor Alt Count=338;Tumor Ref Count=960;Normal Alt Count=23;Normal Ref Count=615;End Position=42232034;Mutation Type=Nonstop_Mutation;Ncbi Build=37;Strand=+;Variant Type=INS;Entrez Gene Id=171
\.

  -- autocompletion tables
  CREATE TABLE genomic_annotations.hugo_gene_symbol as
    select distinct hugo_gene_symbol as annotation_value from genomic_annotations.genomic_annotations;
  CREATE TABLE genomic_annotations.protein_change as
    select distinct protein_change as annotation_value from genomic_annotations.genomic_annotations;
  CREATE TABLE genomic_annotations.variant_name as
    select distinct variant_name as annotation_value from genomic_annotations.genomic_annotations;

  -- permissions
  ALTER TABLE genomic_annotations.genomic_annotations OWNER TO $MC_DB_USER;
  ALTER TABLE genomic_annotations.annotation_names OWNER TO $MC_DB_USER;
  ALTER TABLE genomic_annotations.gene_values OWNER TO $MC_DB_USER;
  ALTER TABLE genomic_annotations.hugo_gene_symbol OWNER TO $MC_DB_USER;
  ALTER TABLE genomic_annotations.protein_change OWNER TO $MC_DB_USER;
  ALTER TABLE genomic_annotations.variant_name OWNER TO $MC_DB_USER;
  GRANT ALL on schema genomic_annotations to $MC_DB_USER;
  GRANT ALL privileges on all tables in schema genomic_annotations to $MC_DB_USER;

COMMIT;
EOSQL
