-- this file contains a full example set of SQL statements that should be produced by the MedCo loading tool

-- random information
-- warning about potential confusion: patient_id != sample_id in i2b2, but same in 95% of the dataset
-- row containing the name of the fields in the dataset files: the 6th row for clinical, 2nd for genomic

-- loading order
-- 1. ontology in shrine_ont.{clinical_sensitive, clinical_non_sensitive, genomic}
-- 2. ontology in i2b2metadata.{sensitive_tagged, non_sensitive_clear} + i2b2demodata.concept_dimension
-- 3. patients + samples in i2b2demodata.{patient_mapping, encounter_mapping, patient_dimension, visit_dimension} + provider_dimension
-- 4. data in observation_fact



------------------------------------------------------------------------------------------------------------------------
-----SCHEMA shrine_ont (SHRINE ontology cell ONT)----------------------------------------------------------------------------
-----TABLE clinical_sensitive-------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- ontology data for sensitive clinical attributes, from the clinical dataset
--- example with 2 clinical fields classified as sensitive (PRIMARY_TUMOR_LOCALIZATION_TYPE and CANCER_TYPE_DETAILED)

-- 1 entry per field (= 2), level 3
INSERT INTO shrine_ont.clinical_sensitive VALUES (3, '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\', 'PRIMARY_TUMOR_LOCALIZATION_TYPE', 'N', 'CA ', NULL, NULL, NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\', 'Sensitive field encrypted by Unlynx', '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\', 'NOW()', NULL, NULL, NULL, 'ENC_ID', '@', NULL, NULL, NULL, NULL);
INSERT INTO shrine_ont.clinical_sensitive VALUES (3, '\\medco\\clinical\\sensitive\\CANCER_TYPE_DETAILED\\', 'CANCER_TYPE_DETAILED', 'N', 'CA ', NULL, NULL, NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\sensitive\\CANCER_TYPE_DETAILED\\', 'Sensitive field encrypted by Unlynx', '\\medco\\clinical\\sensitive\\CANCER_TYPE_DETAILED\\', 'NOW()', NULL, NULL, NULL, 'ENC_ID', '@', NULL, NULL, NULL, NULL);

-- 1 entry per different value occurring for each field (= 5 + 1), level 4
-- c_basecode column is of the form "ENC_ID:X", with X being a random and unique (within the ENC_ID values) integer
-- notice the special "<empty>" value, which for when there is no value at all in the dataset
INSERT INTO shrine_ont.clinical_sensitive VALUES (4, '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\Mucosal\\', 'Mucosal', 'N', 'LA ', NULL, 'ENC_ID:1', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\Mucosal\\', 'Sensitive value encrypted by Unlynx', '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\Mucosal\\', 'NOW()', NULL, NULL, NULL, 'ENC_ID', '@', NULL, NULL, NULL, NULL);
INSERT INTO shrine_ont.clinical_sensitive VALUES (4, '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\<empty>\\', '<empty>', 'N', 'LA ', NULL, 'ENC_ID:2', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\<empty>\\', 'Sensitive value encrypted by Unlynx', '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\<empty>\\', 'NOW()', NULL, NULL, NULL, 'ENC_ID', '@', NULL, NULL, NULL, NULL);
INSERT INTO shrine_ont.clinical_sensitive VALUES (4, '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\Acral\\', 'Acral', 'N', 'LA ', NULL, 'ENC_ID:3', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\Acral\\', 'Sensitive value encrypted by Unlynx', '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\Acral\\', 'NOW()', NULL, NULL, NULL, 'ENC_ID', '@', NULL, NULL, NULL, NULL);
INSERT INTO shrine_ont.clinical_sensitive VALUES (4, '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\Skin\\', 'Skin', 'N', 'LA ', NULL, 'ENC_ID:4', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\Skin\\', 'Sensitive value encrypted by Unlynx', '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\Skin\\', 'NOW()', NULL, NULL, NULL, 'ENC_ID', '@', NULL, NULL, NULL, NULL);
INSERT INTO shrine_ont.clinical_sensitive VALUES (4, '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\Uveal\\', 'Uveal', 'N', 'LA ', NULL, 'ENC_ID:5', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\Uveal\\', 'Sensitive value encrypted by Unlynx', '\\medco\\clinical\\sensitive\\PRIMARY_TUMOR_LOCALIZATION_TYPE\\Uveal\\', 'NOW()', NULL, NULL, NULL, 'ENC_ID', '@', NULL, NULL, NULL, NULL);
INSERT INTO shrine_ont.clinical_sensitive VALUES (4, '\\medco\\clinical\\sensitive\\CANCER_TYPE_DETAILED\\Cutaneous Melanoma\\', 'Cutaneous Melanoma', 'N', 'LA ', NULL, 'ENC_ID:6', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\sensitive\\CANCER_TYPE_DETAILED\\Cutaneous Melanoma\\', 'Sensitive value encrypted by Unlynx', '\\medco\\clinical\\sensitive\\CANCER_TYPE_DETAILED\\Cutaneous Melanoma\\', 'NOW()', NULL, NULL, NULL, 'ENC_ID', '@', NULL, NULL, NULL, NULL);


------------------------------------------------------------------------------------------------------------------------
-----TABLE clinical_non_sensitive---------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- ontology data for non-sensitive clinical attributes, from the clinical dataset
--- example with 2 clinical fields classified as non sensitive (CANCER_TYPE and PRIMARY_DIAGNOSIS)

-- 1 entry per field (= 2), level 3
INSERT INTO shrine_ont.clinical_non_sensitive VALUES (3, '\\medco\\clinical\\nonsensitive\\CANCER_TYPE\\', 'CANCER_TYPE', 'N', 'CA ', NULL, NULL, NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\nonsensitive\\CANCER_TYPE\\', 'Non-sensitive field', '\\medco\\clinical\\nonsensitive\\CANCER_TYPE\\', 'NOW()', NULL, NULL, NULL, 'CLEAR', '@', NULL, NULL, NULL, NULL);
INSERT INTO shrine_ont.clinical_non_sensitive VALUES (3, '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\', 'PRIMARY_DIAGNOSIS', 'N', 'CA ', NULL, NULL, NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\', 'Non-sensitive field', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\', 'NOW()', NULL, NULL, NULL, 'CLEAR', '@', NULL, NULL, NULL, NULL);

-- 1 entry per different value occurring for each field (= 3 + 1), level 4
-- c_basecode column is of the form "CLEAR:X", with X being a unique integer (within the CLEAR values), not necessarly random
-- notice the special "<empty>" value, which for when there is no value at all in the dataset
INSERT INTO shrine_ont.clinical_non_sensitive VALUES (4, '\\medco\\clinical\\nonsensitive\\CANCER_TYPE\\Melanoma\\', 'Melanoma', 'N', 'LA ', NULL, 'CLEAR:1', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\nonsensitive\\CANCER_TYPE\\Melanoma\\', 'Non-sensitive value', '\\medco\\clinical\\nonsensitive\\CANCER_TYPE\\Melanoma\\', 'NOW()', NULL, NULL, NULL, 'CLEAR', '@', NULL, NULL, NULL, NULL);
INSERT INTO shrine_ont.clinical_non_sensitive VALUES (4, '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\Melanoma\\', 'Melanoma', 'N', 'LA ', NULL, 'CLEAR:2', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\Melanoma\\', 'Non-sensitive value', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\Melanoma\\', 'NOW()', NULL, NULL, NULL, 'CLEAR', '@', NULL, NULL, NULL, NULL);
INSERT INTO shrine_ont.clinical_non_sensitive VALUES (4, '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\<empty>\\', '<empty>', 'N', 'LA ', NULL, 'CLEAR:3', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\<empty>\\', 'Non-sensitive value', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\<empty>\\', 'NOW()', NULL, NULL, NULL, 'CLEAR', '@', NULL, NULL, NULL, NULL);
INSERT INTO shrine_ont.clinical_non_sensitive VALUES (4, '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\Not Melanoma\\', 'Not Melanoma', 'N', 'LA ', NULL, 'CLEAR:4', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\Not Melanoma\\', 'Non-sensitive value', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\Not Melanoma\\', 'NOW()', NULL, NULL, NULL, 'CLEAR', '@', NULL, NULL, NULL, NULL);


------------------------------------------------------------------------------------------------------------------------
-----TABLE genomic_annotations------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- ontology data for genomic attributes, from the genomic dataset
--- example with 2 variants that has 2 and 3 annotations, in the special annotations table
INSERT INTO shrine_ont.genomic_annotations VALUES ('-8976521235638865', '{Protein_Position:600/766, Hugo_Symbol:BRCA1}');
INSERT INTO shrine_ont.genomic_annotations VALUES ('-2938472982331123', '{Protein_Position:722/766, Hugo_Symbol:PRPF, atf:3.2}');



------------------------------------------------------------------------------------------------------------------------
-----SCHEMA i2b2metadata (ontology cell ONT)----------------------------------------------------------------------------
-----TABLE sensitive_tagged---------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- ontology data for all sensitive attributes, in encrypted / tagged form for i2b2 to answer queries
--- example with several tagged values (both clinical and genomic)
INSERT INTO i2b2metadata.sensitive_tagged VALUES (2, '\\medco\\tagged\\EkaojcPm7U3qsQp0bhzaLZLYenL/+yNS5j39TFcLU1u=\\', '', 'N', 'LA ', NULL, 'TAG_ID:563255632', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\tagged\\EkaojcPm7U3qsQp0bhzaLZLYenL/+yNS5j39TFcLU1u=\\', NULL, NULL, 'NOW()', NULL, NULL, NULL, 'TAG_ID', '@', NULL, NULL, NULL, NULL);
INSERT INTO i2b2metadata.sensitive_tagged VALUES (2, '\\medco\\tagged\\BMSfLSsNrDeTssfy57z5DfT8V/4u9cE7UWFjgBPpu7y=\\', '', 'N', 'LA ', NULL, 'TAG_ID:2325434152', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\tagged\\BMSfLSsNrDeTssfy57z5DfT8V/4u9cE7UWFjgBPpu7y=\\', NULL, NULL, 'NOW()', NULL, NULL, NULL, 'TAG_ID', '@', NULL, NULL, NULL, NULL);
INSERT INTO i2b2metadata.sensitive_tagged VALUES (2, '\\medco\\tagged\\eTssfFjgBPBMSfLSsNrDpu7yy57z5DfT8V/4u9cE7UW=\\', '', 'N', 'LA ', NULL, 'TAG_ID:2011256355', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\tagged\\eTssfFjgBPBMSfLSsNrDpu7yy57z5DfT8V/4u9cE7UW=\\', NULL, NULL, 'NOW()', NULL, NULL, NULL, 'TAG_ID', '@', NULL, NULL, NULL, NULL);
INSERT INTO i2b2metadata.sensitive_tagged VALUES (2, '\\medco\\tagged\\U3qsQp0bhzaEkaojcP39TFcLU1um7LZLYenL/+yNS5j=\\', '', 'N', 'LA ', NULL, 'TAG_ID:321455215', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\tagged\\U3qsQp0bhzaEkaojcP39TFcLU1um7LZLYenL/+yNS5j=\\', NULL, NULL, 'NOW()', NULL, NULL, NULL, 'TAG_ID', '@', NULL, NULL, NULL, NULL);
INSERT INTO i2b2metadata.sensitive_tagged VALUES (2, '\\medco\\tagged\\8V/4u9cE7UWFjgBPpu7yBMSfLSsNrDeTssfy57z5DfT=\\', '', 'N', 'LA ', NULL, 'TAG_ID:984949149', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\tagged\\8V/4u9cE7UWFjgBPpu7yBMSfLSsNrDeTssfy57z5DfT=\\', NULL, NULL, 'NOW()', NULL, NULL, NULL, 'TAG_ID', '@', NULL, NULL, NULL, NULL);
INSERT INTO i2b2metadata.sensitive_tagged VALUES (2, '\\medco\\tagged\\WFjgBPpu7yBMSfLSsNrDeTs8V/sfy57z5DfT4u9cE7U=\\', '', 'N', 'LA ', NULL, 'TAG_ID:1052524212', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\tagged\\WFjgBPpu7yBMSfLSsNrDeTs8V/sfy57z5DfT4u9cE7U=\\', NULL, NULL, 'NOW()', NULL, NULL, NULL, 'TAG_ID', '@', NULL, NULL, NULL, NULL);


------------------------------------------------------------------------------------------------------------------------
-----TABLE non_sensitive_clear------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- ontology data for non-sensitive clinical attributes, from the clinical dataset
--- example with 2 clinical fields classified as non sensitive (CANCER_TYPE and PRIMARY_DIAGNOSIS)

-- 1 entry per field (= 2), level 3
INSERT INTO i2b2metadata.non_sensitive_clear VALUES (3, '\\medco\\clinical\\nonsensitive\\CANCER_TYPE\\', 'CANCER_TYPE', 'N', 'CA ', NULL, NULL, NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\nonsensitive\\CANCER_TYPE\\', 'Non-sensitive field', '\\medco\\clinical\\nonsensitive\\CANCER_TYPE\\', 'NOW()', NULL, NULL, NULL, 'CLEAR', '@', NULL, NULL, NULL, NULL);
INSERT INTO i2b2metadata.non_sensitive_clear VALUES (3, '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\', 'PRIMARY_DIAGNOSIS', 'N', 'CA ', NULL, NULL, NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\', 'Non-sensitive field', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\', 'NOW()', NULL, NULL, NULL, 'CLEAR', '@', NULL, NULL, NULL, NULL);

-- 1 entry per different value occurring for each field (= 3 + 1), level 4
-- c_basecode column is of the form "CLEAR:X", with X being a unique integer (within the CLEAR values), not necessarly random
-- notice the special "<empty>" value, which for when there is no value at all in the dataset
INSERT INTO i2b2metadata.non_sensitive_clear VALUES (4, '\\medco\\clinical\\nonsensitive\\CANCER_TYPE\\Melanoma\\', 'Melanoma', 'N', 'LA ', NULL, 'CLEAR:1', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\nonsensitive\\CANCER_TYPE\\Melanoma\\', 'Non-sensitive value', '\\medco\\clinical\\nonsensitive\\CANCER_TYPE\\Melanoma\\', 'NOW()', NULL, NULL, NULL, 'CLEAR', '@', NULL, NULL, NULL, NULL);
INSERT INTO i2b2metadata.non_sensitive_clear VALUES (4, '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\Melanoma\\', 'Melanoma', 'N', 'LA ', NULL, 'CLEAR:2', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\Melanoma\\', 'Non-sensitive value', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\Melanoma\\', 'NOW()', NULL, NULL, NULL, 'CLEAR', '@', NULL, NULL, NULL, NULL);
INSERT INTO i2b2metadata.non_sensitive_clear VALUES (4, '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\<empty>\\', '<empty>', 'N', 'LA ', NULL, 'CLEAR:3', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\<empty>\\', 'Non-sensitive value', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\<empty>\\', 'NOW()', NULL, NULL, NULL, 'CLEAR', '@', NULL, NULL, NULL, NULL);
INSERT INTO i2b2metadata.non_sensitive_clear VALUES (4, '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\Not Melanoma\\', 'Not Melanoma', 'N', 'LA ', NULL, 'CLEAR:4', NULL, 'concept_cd', 'concept_dimension', 'concept_path', 'T', 'LIKE', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\Not Melanoma\\', 'Non-sensitive value', '\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\Not Melanoma\\', 'NOW()', NULL, NULL, NULL, 'CLEAR', '@', NULL, NULL, NULL, NULL);



------------------------------------------------------------------------------------------------------------------------
-----SCHEMA i2b2demodata (data repository cell CRC)---------------------------------------------------------------------
-----TABLE concept_dimension--------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- reduced set of ontology data specific to the CRC, that is joined to the observation_fact table to answer queries
-- contains the clear concept the same way they are in the ontology
-- for the sensitive ones it contains the tagged identifiers
--- example follows what is in the ontology cell

-- clear-text concepts (clinical non-sensitive)
INSERT INTO i2b2demodata.concept_dimension VALUES ('\\medco\\clinical\\nonsensitive\\CANCER_TYPE\\Melanoma\\', 'CLEAR:1', 'Melanoma', NULL, NULL, NULL, 'NOW()', NULL, NULL);
INSERT INTO i2b2demodata.concept_dimension VALUES ('\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\Melanoma\\', 'CLEAR:2', 'Melanoma', NULL, NULL, NULL, 'NOW()', NULL, NULL);
INSERT INTO i2b2demodata.concept_dimension VALUES ('\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\<empty>\\', 'CLEAR:3', '<empty>', NULL, NULL, NULL, 'NOW()', NULL, NULL);
INSERT INTO i2b2demodata.concept_dimension VALUES ('\\medco\\clinical\\nonsensitive\\PRIMARY_DIAGNOSIS\\Not Melanoma\\', 'CLEAR:4', 'Not Melanoma', NULL, NULL, NULL, 'NOW()', NULL, NULL);

-- tagged concepts (both clinical sensitive and genomic)
-- concept_cd is a unique, random 32-bits integer with the "TAG_ID" prefix
INSERT INTO i2b2demodata.concept_dimension VALUES ('\\medco\\tagged\\EkaojcPm7U3qsQp0bhzaLZLYenL/+yNS5j39TFcLU1u=\\', 'TAG_ID:563255632', NULL, NULL, NULL, NULL, 'NOW()', NULL, NULL);
INSERT INTO i2b2demodata.concept_dimension VALUES ('\\medco\\tagged\\BMSfLSsNrDeTssfy57z5DfT8V/4u9cE7UWFjgBPpu7y=\\', 'TAG_ID:2325434152', NULL, NULL, NULL, NULL, 'NOW()', NULL, NULL);
INSERT INTO i2b2demodata.concept_dimension VALUES ('\\medco\\tagged\\eTssfFjgBPBMSfLSsNrDpu7yy57z5DfT8V/4u9cE7UW=\\', 'TAG_ID:2011256355', NULL, NULL, NULL, NULL, 'NOW()', NULL, NULL);
INSERT INTO i2b2demodata.concept_dimension VALUES ('\\medco\\tagged\\U3qsQp0bhzaEkaojcP39TFcLU1um7LZLYenL/+yNS5j=\\', 'TAG_ID:321455215', NULL, NULL, NULL, NULL, 'NOW()', NULL, NULL);
INSERT INTO i2b2demodata.concept_dimension VALUES ('\\medco\\tagged\\8V/4u9cE7UWFjgBPpu7yBMSfLSsNrDeTssfy57z5DfT=\\', 'TAG_ID:984949149', NULL, NULL, NULL, NULL, 'NOW()', NULL, NULL);
INSERT INTO i2b2demodata.concept_dimension VALUES ('\\medco\\tagged\\WFjgBPpu7yBMSfLSsNrDeTs8V/sfy57z5DfT4u9cE7U=\\', 'TAG_ID:1052524212', NULL, NULL, NULL, NULL, 'NOW()', NULL, NULL);


------------------------------------------------------------------------------------------------------------------------
-----TABLE patient_mapping----------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- mappings of patient identifiers: internal i2b2 id (patient_num column) <-> external id (patient_ide)
-- if the external source (patient_ide_source) has the special value of "HIVE", it means the corresponding id is the i2b2 internal one
-- (every patient registered in i2b2 has such a record)
--- example with 2 patients: "MEL-Ma-Mel-103b" (patient_num = 40) and "MEL-Ma-Mel-102" (patient_num = 39) => 1 patient = 2 entries

-- entries containing the identifiers from the external source (here the source is defined as "chuv")
INSERT INTO i2b2demodata.patient_mapping VALUES ('MEL-Ma-Mel-103b', 'chuv', 40, NULL, 'Demo', NULL, NULL, NULL, 'NOW()', NULL, 1);
INSERT INTO i2b2demodata.patient_mapping VALUES ('MEL-Ma-Mel-102', 'chuv', 39, NULL, 'Demo', NULL, NULL, NULL, 'NOW()', NULL, 1);

-- entries containing the internal i2b2 identifiers, with the special source "HIVE"
INSERT INTO i2b2demodata.patient_mapping VALUES ('40', 'HIVE', 40, 'A', 'HIVE', NULL, 'NOW()', 'NOW()', 'NOW()', 'edu.harvard.i2b2.crc', 1);
INSERT INTO i2b2demodata.patient_mapping VALUES ('39', 'HIVE', 39, 'A', 'HIVE', NULL, 'NOW()', 'NOW()', 'NOW()', 'edu.harvard.i2b2.crc', 1);


------------------------------------------------------------------------------------------------------------------------
-----TABLE patient_dimension--------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- dimension table containing the patient data, identified by their internal i2b2 identifiers (patient_num)
-- this includes the dummy encrypted flag in the "enc_dummy_flag" column
--- example with 2 patients: same as in the patient_mapping table, 1 patient = 1 entry

INSERT INTO i2b2demodata.patient_dimension VALUES (40, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, "FzXxSbBn86gMmF7WT6a4kHDcHrOg3SEkaojcPm7U3qsQp0bhzaLZLYenL/+yNS5j39TFcLU1uSUE5I8tD3Qryw==");
INSERT INTO i2b2demodata.patient_dimension VALUES (39, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, "66xaTIbPcE8V/4u9cE7UWFjgBPpu7yBMSfLSsNrDeTssfy57z5DfTAI+ynrVMzosOapo2SqQxRrrKFSWIljEbw==");


------------------------------------------------------------------------------------------------------------------------
-----TABLE encounter_mapping--------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- mappings of the encounter (= visit) identifiers: similar use as "patient_mapping", note that it additionally contains
-- the patient identifier
--- example with 2 visits/encounters: "MEL-Ma-Mel-103b" and "MEL-Ma-Mel-102"
--- note: in the example of this specific dataset the visits are interpreted as samples, and there are as many samples
--- as patients, but this is not necessarily always the case

-- entries containing the identifiers from the external source (here the source is defined as "chuv")
INSERT INTO i2b2demodata.encounter_mapping VALUES ('MEL-Ma-Mel-103b', 'chuv', 'Demo', 36, 'MEL-Ma-Mel-103b', 'chuv', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1);
INSERT INTO i2b2demodata.encounter_mapping VALUES ('MEL-Ma-Mel-102', 'chuv', 'Demo', 30, 'MEL-Ma-Mel-102', 'chuv', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1);

-- entries containing the internal i2b2 identifiers, with the special source "HIVE"
INSERT INTO i2b2demodata.encounter_mapping VALUES ('36', 'HIVE', 'HIVE', 36, 'MEL-Ma-Mel-103b', 'chuv', 'A', NULL, 'NOW()', 'NOW()', 'NOW()', 'edu.harvard.i2b2.crc', 1);
INSERT INTO i2b2demodata.encounter_mapping VALUES ('30', 'HIVE', 'HIVE', 30, 'MEL-Ma-Mel-102', 'chuv', 'A', NULL, 'NOW()', 'NOW()', 'NOW()', 'edu.harvard.i2b2.crc', 1);


------------------------------------------------------------------------------------------------------------------------
-----TABLE visit_dimension----------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- similar to patient_dimension, contains encounters (sample) identifiers (encounter_num)
INSERT INTO i2b2demodata.visit_dimension VALUES (36, 40, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'NOW()', 'chuv', 1);
INSERT INTO i2b2demodata.visit_dimension VALUES (30, 39, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'NOW()', 'chuv', 1);


------------------------------------------------------------------------------------------------------------------------
-----TABLE provider_dimension-------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- contains the different providers defined in the system
--- here just chuv is defined and used everywhere

INSERT INTO i2b2demodata.provider_dimension VALUES ('chuv', '\\medco\\institutions\\chuv\\', 'chuv', NULL, NULL, NULL, 'NOW()', NULL, 1);


------------------------------------------------------------------------------------------------------------------------
-----TABLE observation_fact---------------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------------
-- the observation facts, that link a patient, a visit, a concept, a provider, and a date

-- clear / normal i2b2 observation facts, for clinical non-sensitive data, nb entries = nb rows in clinical dataset * nb non-sensitive columns
-- concept codes are the ones defined in i2b2metadata.clinical_non_sensitive
INSERT INTO i2b2demodata.observation_fact VALUES (39, 30, 'CLEAR:1', 'chuv', 'NOW()', '@', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'chuv', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, 1);
INSERT INTO i2b2demodata.observation_fact VALUES (39, 30, 'CLEAR:2', 'chuv', 'NOW()', '@', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'chuv', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, 1);
INSERT INTO i2b2demodata.observation_fact VALUES (40, 36, 'CLEAR:1', 'chuv', 'NOW()', '@', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'chuv', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, 1);
INSERT INTO i2b2demodata.observation_fact VALUES (40, 36, 'CLEAR:3', 'chuv', 'NOW()', '@', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'chuv', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, 1);

-- encrypted observation facts, contains both clinical sensitive and genomic tagged values (1 tagged value = 32B = 44 base64 characters)
-- here as an example we are adding 3 sensitive attributes
INSERT INTO i2b2demodata.observation_fact VALUES (39, 30, 'TAG_ID:563255632', 'chuv', 'NOW()', '@', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'chuv', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, 1);
INSERT INTO i2b2demodata.observation_fact VALUES (39, 30, 'TAG_ID:2325434152', 'chuv', 'NOW()', '@', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'chuv', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, 1);
INSERT INTO i2b2demodata.observation_fact VALUES (39, 30, 'TAG_ID:2011256355', 'chuv', 'NOW()', '@', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'chuv', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, 1);
INSERT INTO i2b2demodata.observation_fact VALUES (40, 36, 'TAG_ID:563255632', 'chuv', 'NOW()', '@', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'chuv', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, 1);
INSERT INTO i2b2demodata.observation_fact VALUES (40, 36, 'TAG_ID:2325434152', 'chuv', 'NOW()', '@', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'chuv', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, 1);
INSERT INTO i2b2demodata.observation_fact VALUES (40, 36, 'TAG_ID:984949149', 'chuv', 'NOW()', '@', 1, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'chuv', NULL, NULL, NULL, NULL, 'NOW()', NULL, 1, 1);

