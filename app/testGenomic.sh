#!/usr/bin/env bash
./medco-loader -debug 2 v0 -g /opt/medco/medco-deployment/configuration-profiles/dev-3nodes-1host/group.toml --entryPointIdx 0 \
--ont_clinical ../data/genomic/tcga_cbio/clinical_data.csv --sen ../data/genomic/sensitive.txt  \
--ont_genomic ../data/genomic/tcga_cbio/mutation_data.csv  \
--clinical ../data/genomic/tcga_cbio/clinical_data.csv \
--genomic ../data/genomic/tcga_cbio/mutation_data.csv \
--output ../data/genomic/ \
--dbHost localhost --dbPort 5432 --dbName i2b2medcosrv0 --dbUser i2b2 --dbPassword i2b2

./medco-loader -debug 2 v0 -g /opt/medco/medco-deployment/configuration-profiles/dev-3nodes-1host/group.toml --entryPointIdx 1 \
--ont_clinical ../data/genomic/tcga_cbio/clinical_data.csv --sen ../data/genomic/sensitive.txt \
--ont_genomic ../data/genomic/tcga_cbio/mutation_data.csv  \
--clinical ../data/genomic/tcga_cbio/clinical_data.csv \
--genomic ../data/genomic/tcga_cbio/mutation_data.csv \
--output ../data/genomic/ \
--dbHost localhost --dbPort 5432 --dbName i2b2medcosrv1 --dbUser i2b2 --dbPassword i2b2

./medco-loader -debug 2 v0 -g /opt/medco/medco-deployment/configuration-profiles/dev-3nodes-1host/group.toml --entryPointIdx 2 \
--ont_clinical ../data/genomic/tcga_cbio/clinical_data.csv --sen ../data/genomic/sensitive.txt \
--ont_genomic ../data/genomic/tcga_cbio/mutation_data.csv  \
--clinical ../data/genomic/tcga_cbio/clinical_data.csv \
--genomic ../data/genomic/tcga_cbio/mutation_data.csv \
--output ../data/genomic/ \
--dbHost localhost --dbPort 5432 --dbName i2b2medcosrv2 --dbUser i2b2 --dbPassword i2b2