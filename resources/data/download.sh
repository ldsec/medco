#!/usr/bin/env bash
set -Eeuo pipefail

### script to download the MedCo example datasets

SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
REPO_URL="https://github.com/ldsec/projects-data/blob/master/medco/datasets"

# genomic (v0) dataset
mkdir -p ${SCRIPT_FOLDER}/genomic/tcga_cbio/
wget -O ${SCRIPT_FOLDER}/genomic/tcga_cbio/mutation_data.csv ${REPO_URL}/genomic/tcga_cbio/mutation_data.csv?raw=true
wget -O ${SCRIPT_FOLDER}/genomic/tcga_cbio/clinical_data.csv ${REPO_URL}/genomic/tcga_cbio/clinical_data.csv?raw=true
wget -O ${SCRIPT_FOLDER}/genomic/tcga_cbio/8_mutation_data.csv ${REPO_URL}/genomic/tcga_cbio/8_mutation_data.csv?raw=true
wget -O ${SCRIPT_FOLDER}/genomic/tcga_cbio/8_clinical_data.csv ${REPO_URL}/genomic/tcga_cbio/8_clinical_data.csv?raw=true
wget -O ${SCRIPT_FOLDER}/genomic/sensitive.txt ${REPO_URL}/genomic/sensitive.txt?raw=true

# i2b2 demo (v1) dataset
mkdir -p ${SCRIPT_FOLDER}/i2b2/original/
mkdir -p ${SCRIPT_FOLDER}/i2b2/converted/
wget -O ${SCRIPT_FOLDER}/i2b2/original/birn.csv ${REPO_URL}/i2b2/original/birn.csv?raw=true
wget -O ${SCRIPT_FOLDER}/i2b2/original/concept_dimension.csv ${REPO_URL}/i2b2/original/concept_dimension.csv?raw=true
wget -O ${SCRIPT_FOLDER}/i2b2/original/custom_meta.csv ${REPO_URL}/i2b2/original/custom_meta.csv?raw=true
wget -O ${SCRIPT_FOLDER}/i2b2/original/dummy_to_patient.csv ${REPO_URL}/i2b2/original/dummy_to_patient.csv?raw=true
wget -O ${SCRIPT_FOLDER}/i2b2/original/i2b2.csv ${REPO_URL}/i2b2/original/i2b2.csv?raw=true
wget -O ${SCRIPT_FOLDER}/i2b2/original/icd10_icd9.csv ${REPO_URL}/i2b2/original/icd10_icd9.csv?raw=true
wget -O ${SCRIPT_FOLDER}/i2b2/original/observation_fact.csv ${REPO_URL}/i2b2/original/observation_fact.csv?raw=true
wget -O ${SCRIPT_FOLDER}/i2b2/original/patient_dimension.csv ${REPO_URL}/i2b2/original/patient_dimension.csv?raw=true
wget -O ${SCRIPT_FOLDER}/i2b2/original/table_access.csv ${REPO_URL}/i2b2/original/table_access.csv?raw=true
wget -O ${SCRIPT_FOLDER}/i2b2/original/visit_dimension.csv ${REPO_URL}/i2b2/original/visit_dimension.csv?raw=true
wget -O ${SCRIPT_FOLDER}/i2b2/original/concept_dimension.csv ${REPO_URL}/i2b2/original/concept_dimension.csv?raw=true
wget -O ${SCRIPT_FOLDER}/i2b2/sensitive.txt ${REPO_URL}/i2b2/sensitive.txt?raw=true
wget -O ${SCRIPT_FOLDER}/i2b2/files.toml ${REPO_URL}/i2b2/files.toml?raw=true
