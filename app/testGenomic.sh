#!/usr/bin/env bash
./medcoLoader -debug 2 v0 -g /home/jagomes/medco-deployment/configuration-profiles/prod/3nodes-samehost/group.toml --entryPointIdx 0 \
--ont_clinical ../data/genomic/tcga_cbio/clinical_data.csv --sensitive CANCER_TYPE_DETAILED \
--ont_genomic ../data/genomic/tcga_cbio/mutation_data.csv  \
--clinical ../data/genomic/tcga_cbio/manipulations/80_node0_clinical_data.csv \
--genomic ../data/genomic/tcga_cbio/manipulations/80_node0_mutation_data.csv \
--dbHost postgresql --dbPort 5432 --dbName i2b2medcosrv0 --dbUser i2b2 --dbPassword i2b2

./medcoLoader -debug 2 v0 -g /home/jagomes/medco-deployment/configuration-profiles/prod/3nodes-samehost/group.toml --entryPointIdx 1 \
--ont_clinical ../data/genomic/tcga_cbio/clinical_data.csv --sensitive CANCER_TYPE_DETAILED \
--ont_genomic ../data/genomic/tcga_cbio/mutation_data.csv  \
--clinical ../data/genomic/tcga_cbio/manipulations/80_node1_clinical_data.csv \
--genomic ../data/genomic/tcga_cbio/manipulations/80_node1_mutation_data.csv \
--dbHost postgresql --dbPort 5432 --dbName i2b2medcosrv2 --dbUser i2b2 --dbPassword i2b2

./medcoLoader -debug 2 v0 -g /home/jagomes/medco-deployment/configuration-profiles/prod/3nodes-samehost/group.toml --entryPointIdx 2 \
--ont_clinical ./clinical_data.csv --sensitive CANCER_TYPE_DETAILED \
--ont_genomic ./mutation_data.csv  \
--clinical ../data/genomic/tcga_cbio/manipulations/80_node2_clinical_data.csv \
--genomic ../data/genomic/tcga_cbio/manipulations/80_node2_mutation_data.csv \
--dbHost postgresql --dbPort 5432 --dbName i2b2medcosrv2 --dbUser i2b2 --dbPassword i2b2