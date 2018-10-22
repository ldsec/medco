
SET search_path TO genomic_annotations;

-- --- variants table ---
create table genomic_annotations
(
    variant_id character varying(255) NOT NULL PRIMARY KEY,
    variant_name character varying(255) NOT NULL,
    annotations text NOT NULL,
    t_depth numeric NOT NULL
);

ALTER TABLE genomic_annotations OWNER TO genomic_annotations;
\copy genomic_annotations FROM 'annotation_tables/SHRINE_ONT_GENOMIC_ANNOTATIONS_NEW.csv' ESCAPE '"' DELIMITER ',' CSV;

-- --- annotation_names ---
create table annotation_names
(   
    annotation_name character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE annotation_names OWNER TO genomic_annotations;
\copy annotation_names FROM 'annotation_tables/annotation_names' ESCAPE '"' DELIMITER ',' CSV;

-- --- GENE VALUES ---
create table gene_values
(
    gene_value character varying(255) NOT NULL PRIMARY KEY
);
   
ALTER TABLE gene_values OWNER TO genomic_annotations;
\copy gene_values FROM 'annotation_tables/SHRINE_GENES.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- AA_MAF ---
create table AA_MAF
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE AA_MAF OWNER TO genomic_annotations;
\copy AA_MAF FROM 'annotation_tables/AA_MAF.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- AFR_MAF ---
create table AFR_MAF
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE AFR_MAF OWNER TO genomic_annotations;
\copy AFR_MAF FROM 'annotation_tables/AFR_MAF.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- all_effects ---
create table all_effects
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE all_effects OWNER TO genomic_annotations;
\copy all_effects FROM 'annotation_tables/all_effects.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Allele ---
create table Allele
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Allele OWNER TO genomic_annotations;
\copy Allele FROM 'annotation_tables/Allele.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- ALLELE_NUM ---
create table ALLELE_NUM
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE ALLELE_NUM OWNER TO genomic_annotations;
\copy ALLELE_NUM FROM 'annotation_tables/ALLELE_NUM.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Amino_acids ---
create table Amino_acids
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Amino_acids OWNER TO genomic_annotations;
\copy Amino_acids FROM 'annotation_tables/Amino_acids.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- AMR_MAF ---
create table AMR_MAF
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE AMR_MAF OWNER TO genomic_annotations;
\copy AMR_MAF FROM 'annotation_tables/AMR_MAF.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- BIOTYPE ---
create table BIOTYPE
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE BIOTYPE OWNER TO genomic_annotations;
\copy BIOTYPE FROM 'annotation_tables/BIOTYPE.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- CANONICAL ---
create table CANONICAL
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE CANONICAL OWNER TO genomic_annotations;
\copy CANONICAL FROM 'annotation_tables/CANONICAL.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- CCDS ---
create table CCDS
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE CCDS OWNER TO genomic_annotations;
\copy CCDS FROM 'annotation_tables/CCDS.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- cDNA_position ---
create table cDNA_position
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE cDNA_position OWNER TO genomic_annotations;
\copy cDNA_position FROM 'annotation_tables/cDNA_position.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- CDS_position ---
create table CDS_position
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE CDS_position OWNER TO genomic_annotations;
\copy CDS_position FROM 'annotation_tables/CDS_position.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Center ---
create table Center
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Center OWNER TO genomic_annotations;
\copy Center FROM 'annotation_tables/Center.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Chromosome ---
create table Chromosome
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Chromosome OWNER TO genomic_annotations;
\copy Chromosome FROM 'annotation_tables/Chromosome.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- CLIN_SIG ---
create table CLIN_SIG
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE CLIN_SIG OWNER TO genomic_annotations;
\copy CLIN_SIG FROM 'annotation_tables/CLIN_SIG.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Codons ---
create table Codons
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Codons OWNER TO genomic_annotations;
\copy Codons FROM 'annotation_tables/Codons.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Consequence ---
create table Consequence
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Consequence OWNER TO genomic_annotations;
\copy Consequence FROM 'annotation_tables/Consequence.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- dbSNP_RS ---
create table dbSNP_RS
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE dbSNP_RS OWNER TO genomic_annotations;
\copy dbSNP_RS FROM 'annotation_tables/dbSNP_RS.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- dbSNP_Val_Status ---
create table dbSNP_Val_Status
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE dbSNP_Val_Status OWNER TO genomic_annotations;
\copy dbSNP_Val_Status FROM 'annotation_tables/dbSNP_Val_Status.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- DISTANCE ---
create table DISTANCE
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE DISTANCE OWNER TO genomic_annotations;
\copy DISTANCE FROM 'annotation_tables/DISTANCE.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- DOMAINS ---
create table DOMAINS
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE DOMAINS OWNER TO genomic_annotations;
\copy DOMAINS FROM 'annotation_tables/DOMAINS.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- EA_MAF ---
create table EA_MAF
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE EA_MAF OWNER TO genomic_annotations;
\copy EA_MAF FROM 'annotation_tables/EA_MAF.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- EAS_MAF ---
create table EAS_MAF
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE EAS_MAF OWNER TO genomic_annotations;
\copy EAS_MAF FROM 'annotation_tables/EAS_MAF.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- End_Position ---
create table End_Position
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE End_Position OWNER TO genomic_annotations;
\copy End_Position FROM 'annotation_tables/End_Position.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- ENSP ---
create table ENSP
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE ENSP OWNER TO genomic_annotations;
\copy ENSP FROM 'annotation_tables/ENSP.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Entrez_Gene_Id ---
create table Entrez_Gene_Id
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Entrez_Gene_Id OWNER TO genomic_annotations;
\copy Entrez_Gene_Id FROM 'annotation_tables/Entrez_Gene_Id.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- EUR_MAF ---
create table EUR_MAF
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE EUR_MAF OWNER TO genomic_annotations;
\copy EUR_MAF FROM 'annotation_tables/EUR_MAF.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- ExAC_AF_AFR ---
create table ExAC_AF_AFR
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE ExAC_AF_AFR OWNER TO genomic_annotations;
\copy ExAC_AF_AFR FROM 'annotation_tables/ExAC_AF_AFR.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- ExAC_AF_AMR ---
create table ExAC_AF_AMR
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE ExAC_AF_AMR OWNER TO genomic_annotations;
\copy ExAC_AF_AMR FROM 'annotation_tables/ExAC_AF_AMR.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- ExAC_AF ---
create table ExAC_AF
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE ExAC_AF OWNER TO genomic_annotations;
\copy ExAC_AF FROM 'annotation_tables/ExAC_AF.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- ExAC_AF_EAS ---
create table ExAC_AF_EAS
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE ExAC_AF_EAS OWNER TO genomic_annotations;
\copy ExAC_AF_EAS FROM 'annotation_tables/ExAC_AF_EAS.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- ExAC_AF_FIN ---
create table ExAC_AF_FIN
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE ExAC_AF_FIN OWNER TO genomic_annotations;
\copy ExAC_AF_FIN FROM 'annotation_tables/ExAC_AF_FIN.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- ExAC_AF_NFE ---
create table ExAC_AF_NFE
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE ExAC_AF_NFE OWNER TO genomic_annotations;
\copy ExAC_AF_NFE FROM 'annotation_tables/ExAC_AF_NFE.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- ExAC_AF_OTH ---
create table ExAC_AF_OTH
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE ExAC_AF_OTH OWNER TO genomic_annotations;
\copy ExAC_AF_OTH FROM 'annotation_tables/ExAC_AF_OTH.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- ExAC_AF_SAS ---
create table ExAC_AF_SAS
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE ExAC_AF_SAS OWNER TO genomic_annotations;
\copy ExAC_AF_SAS FROM 'annotation_tables/ExAC_AF_SAS.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Existing_variation ---
create table Existing_variation
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Existing_variation OWNER TO genomic_annotations;
\copy Existing_variation FROM 'annotation_tables/Existing_variation.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- EXON ---
create table EXON
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE EXON OWNER TO genomic_annotations;
\copy EXON FROM 'annotation_tables/EXON.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Exon_Number ---
create table Exon_Number
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Exon_Number OWNER TO genomic_annotations;
\copy Exon_Number FROM 'annotation_tables/Exon_Number.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Feature ---
create table Feature
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Feature OWNER TO genomic_annotations;
\copy Feature FROM 'annotation_tables/Feature.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Feature_type ---
create table Feature_type
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Feature_type OWNER TO genomic_annotations;
\copy Feature_type FROM 'annotation_tables/Feature_type.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- FILTER ---
create table FILTER
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE FILTER OWNER TO genomic_annotations;
\copy FILTER FROM 'annotation_tables/FILTER.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Gene ---
create table Gene
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Gene OWNER TO genomic_annotations;
\copy Gene FROM 'annotation_tables/Gene.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- GENE_PHENO ---
create table GENE_PHENO
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE GENE_PHENO OWNER TO genomic_annotations;
\copy GENE_PHENO FROM 'annotation_tables/GENE_PHENO.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- GMAF ---
create table GMAF
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE GMAF OWNER TO genomic_annotations;
\copy GMAF FROM 'annotation_tables/GMAF.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- HGNC_ID ---
create table HGNC_ID
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE HGNC_ID OWNER TO genomic_annotations;
\copy HGNC_ID FROM 'annotation_tables/HGNC_ID.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- HGVSc ---
create table HGVSc
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE HGVSc OWNER TO genomic_annotations;
\copy HGVSc FROM 'annotation_tables/HGVSc.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- HGVS_OFFSET ---
create table HGVS_OFFSET
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE HGVS_OFFSET OWNER TO genomic_annotations;
\copy HGVS_OFFSET FROM 'annotation_tables/HGVS_OFFSET.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- HGVSp ---
create table HGVSp
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE HGVSp OWNER TO genomic_annotations;
\copy HGVSp FROM 'annotation_tables/HGVSp.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- HGVSp_Short ---
create table HGVSp_Short
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE HGVSp_Short OWNER TO genomic_annotations;
\copy HGVSp_Short FROM 'annotation_tables/HGVSp_Short.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- HIGH_INF_POS ---
create table HIGH_INF_POS
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE HIGH_INF_POS OWNER TO genomic_annotations;
\copy HIGH_INF_POS FROM 'annotation_tables/HIGH_INF_POS.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Hugo_Symbol ---
create table Hugo_Symbol
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Hugo_Symbol OWNER TO genomic_annotations;
\copy Hugo_Symbol FROM 'annotation_tables/Hugo_Symbol.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- IMPACT ---
create table IMPACT
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE IMPACT OWNER TO genomic_annotations;
\copy IMPACT FROM 'annotation_tables/IMPACT.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- INTRON ---
create table INTRON
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE INTRON OWNER TO genomic_annotations;
\copy INTRON FROM 'annotation_tables/INTRON.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- MA:FImpact ---
create table MA:FImpact
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE MA:FImpact OWNER TO genomic_annotations;
\copy MA:FImpact FROM 'annotation_tables/MA:FImpact.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- MA:FIS ---
create table MA:FIS
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE MA:FIS OWNER TO genomic_annotations;
\copy MA:FIS FROM 'annotation_tables/MA:FIS.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- MA:link.MSA ---
create table MA:link.MSA
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE MA:link.MSA OWNER TO genomic_annotations;
\copy MA:link.MSA FROM 'annotation_tables/MA:link.MSA.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- MA:link.PDB ---
create table MA:link.PDB
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE MA:link.PDB OWNER TO genomic_annotations;
\copy MA:link.PDB FROM 'annotation_tables/MA:link.PDB.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- MA:link.var ---
create table MA:link.var
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE MA:link.var OWNER TO genomic_annotations;
\copy MA:link.var FROM 'annotation_tables/MA:link.var.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- MA:protein.change ---
create table MA:protein.change
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE MA:protein.change OWNER TO genomic_annotations;
\copy MA:protein.change FROM 'annotation_tables/MA:protein.change.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Matched_Norm_Sample_Barcode ---
create table Matched_Norm_Sample_Barcode
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Matched_Norm_Sample_Barcode OWNER TO genomic_annotations;
\copy Matched_Norm_Sample_Barcode FROM 'annotation_tables/Matched_Norm_Sample_Barcode.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Match_Norm_Seq_Allele1 ---
create table Match_Norm_Seq_Allele1
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Match_Norm_Seq_Allele1 OWNER TO genomic_annotations;
\copy Match_Norm_Seq_Allele1 FROM 'annotation_tables/Match_Norm_Seq_Allele1.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Match_Norm_Seq_Allele2 ---
create table Match_Norm_Seq_Allele2
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Match_Norm_Seq_Allele2 OWNER TO genomic_annotations;
\copy Match_Norm_Seq_Allele2 FROM 'annotation_tables/Match_Norm_Seq_Allele2.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- MINIMISED ---
create table MINIMISED
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE MINIMISED OWNER TO genomic_annotations;
\copy MINIMISED FROM 'annotation_tables/MINIMISED.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- MOTIF_NAME ---
create table MOTIF_NAME
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE MOTIF_NAME OWNER TO genomic_annotations;
\copy MOTIF_NAME FROM 'annotation_tables/MOTIF_NAME.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- MOTIF_POS ---
create table MOTIF_POS
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE MOTIF_POS OWNER TO genomic_annotations;
\copy MOTIF_POS FROM 'annotation_tables/MOTIF_POS.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- MOTIF_SCORE_CHANGE ---
create table MOTIF_SCORE_CHANGE
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE MOTIF_SCORE_CHANGE OWNER TO genomic_annotations;
\copy MOTIF_SCORE_CHANGE FROM 'annotation_tables/MOTIF_SCORE_CHANGE.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- n_alt_count ---
create table n_alt_count
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE n_alt_count OWNER TO genomic_annotations;
\copy n_alt_count FROM 'annotation_tables/n_alt_count.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- NCBI_Build ---
create table NCBI_Build
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE NCBI_Build OWNER TO genomic_annotations;
\copy NCBI_Build FROM 'annotation_tables/NCBI_Build.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- n_depth ---
create table n_depth
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE n_depth OWNER TO genomic_annotations;
\copy n_depth FROM 'annotation_tables/n_depth.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- n_ref_count ---
create table n_ref_count
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE n_ref_count OWNER TO genomic_annotations;
\copy n_ref_count FROM 'annotation_tables/n_ref_count.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- PHENO ---
create table PHENO
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE PHENO OWNER TO genomic_annotations;
\copy PHENO FROM 'annotation_tables/PHENO.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- PICK ---
create table PICK
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE PICK OWNER TO genomic_annotations;
\copy PICK FROM 'annotation_tables/PICK.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- PolyPhen ---
create table PolyPhen
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE PolyPhen OWNER TO genomic_annotations;
\copy PolyPhen FROM 'annotation_tables/PolyPhen.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Protein_position ---
create table Protein_position
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Protein_position OWNER TO genomic_annotations;
\copy Protein_position FROM 'annotation_tables/Protein_position.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- PUBMED ---
create table PUBMED
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE PUBMED OWNER TO genomic_annotations;
\copy PUBMED FROM 'annotation_tables/PUBMED.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Reference_Allele ---
create table Reference_Allele
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Reference_Allele OWNER TO genomic_annotations;
\copy Reference_Allele FROM 'annotation_tables/Reference_Allele.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- RefSeq ---
create table RefSeq
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE RefSeq OWNER TO genomic_annotations;
\copy RefSeq FROM 'annotation_tables/RefSeq.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- SAS_MAF ---
create table SAS_MAF
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE SAS_MAF OWNER TO genomic_annotations;
\copy SAS_MAF FROM 'annotation_tables/SAS_MAF.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Sequencer ---
create table Sequencer
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Sequencer OWNER TO genomic_annotations;
\copy Sequencer FROM 'annotation_tables/Sequencer.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- SIFT ---
create table SIFT
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE SIFT OWNER TO genomic_annotations;
\copy SIFT FROM 'annotation_tables/SIFT.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- SOMATIC ---
create table SOMATIC
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE SOMATIC OWNER TO genomic_annotations;
\copy SOMATIC FROM 'annotation_tables/SOMATIC.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Start_Position ---
create table Start_Position
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Start_Position OWNER TO genomic_annotations;
\copy Start_Position FROM 'annotation_tables/Start_Position.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Strand ---
create table Strand
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Strand OWNER TO genomic_annotations;
\copy Strand FROM 'annotation_tables/Strand.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- SWISSPROT ---
create table SWISSPROT
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE SWISSPROT OWNER TO genomic_annotations;
\copy SWISSPROT FROM 'annotation_tables/SWISSPROT.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- SYMBOL ---
create table SYMBOL
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE SYMBOL OWNER TO genomic_annotations;
\copy SYMBOL FROM 'annotation_tables/SYMBOL.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- SYMBOL_SOURCE ---
create table SYMBOL_SOURCE
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE SYMBOL_SOURCE OWNER TO genomic_annotations;
\copy SYMBOL_SOURCE FROM 'annotation_tables/SYMBOL_SOURCE.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- t_alt_count ---
create table t_alt_count
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE t_alt_count OWNER TO genomic_annotations;
\copy t_alt_count FROM 'annotation_tables/t_alt_count.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- t_depth ---
create table t_depth
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE t_depth OWNER TO genomic_annotations;
\copy t_depth FROM 'annotation_tables/t_depth.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Transcript_ID ---
create table Transcript_ID
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Transcript_ID OWNER TO genomic_annotations;
\copy Transcript_ID FROM 'annotation_tables/Transcript_ID.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- t_ref_count ---
create table t_ref_count
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE t_ref_count OWNER TO genomic_annotations;
\copy t_ref_count FROM 'annotation_tables/t_ref_count.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- TREMBL ---
create table TREMBL
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE TREMBL OWNER TO genomic_annotations;
\copy TREMBL FROM 'annotation_tables/TREMBL.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Tumor_Sample_Barcode ---
create table Tumor_Sample_Barcode
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Tumor_Sample_Barcode OWNER TO genomic_annotations;
\copy Tumor_Sample_Barcode FROM 'annotation_tables/Tumor_Sample_Barcode.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Tumor_Seq_Allele1 ---
create table Tumor_Seq_Allele1
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Tumor_Seq_Allele1 OWNER TO genomic_annotations;
\copy Tumor_Seq_Allele1 FROM 'annotation_tables/Tumor_Seq_Allele1.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Tumor_Seq_Allele2 ---
create table Tumor_Seq_Allele2
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Tumor_Seq_Allele2 OWNER TO genomic_annotations;
\copy Tumor_Seq_Allele2 FROM 'annotation_tables/Tumor_Seq_Allele2.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- UNIPARC ---
create table UNIPARC
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE UNIPARC OWNER TO genomic_annotations;
\copy UNIPARC FROM 'annotation_tables/UNIPARC.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- VARIANT_CLASS ---
create table VARIANT_CLASS
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE VARIANT_CLASS OWNER TO genomic_annotations;
\copy VARIANT_CLASS FROM 'annotation_tables/VARIANT_CLASS.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Variant_Classification ---
create table Variant_Classification
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Variant_Classification OWNER TO genomic_annotations;
\copy Variant_Classification FROM 'annotation_tables/Variant_Classification.csv' ESCAPE '"' DELIMITER ',' CSV;



-- --- Variant_Type ---
create table Variant_Type
(
    annotation_value character varying(255) NOT NULL PRIMARY KEY
);

ALTER TABLE Variant_Type OWNER TO genomic_annotations;
\copy Variant_Type FROM 'annotation_tables/Variant_Type.csv' ESCAPE '"' DELIMITER ',' CSV;


