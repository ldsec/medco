#!/usr/bin/env bash
./medcoLoader -debug 2 v0 -g /home/jagomes/medco-deployment/configuration-profiles/prod/3nodes-samehost/group.toml --entryPointIdx 0 \
--ont_clinical ../data/genomic/tcga_cbio/8_clinical_data.csv --sen ../data/genomic/sensitive.txt  \
--ont_genomic ../data/genomic/tcga_cbio/8_mutation_data.csv  \
--clinical ../data/genomic/tcga_cbio/8_clinical_data.csv \
--genomic ../data/genomic/tcga_cbio/8_mutation_data.csv \
--dbHost localhost --dbPort 5432 --dbName i2b2medcosrv0 --dbUser i2b2 --dbPassword i2b2

./medcoLoader -debug 2 v0 -g /home/jagomes/medco-deployment/configuration-profiles/prod/3nodes-samehost/group.toml --entryPointIdx 1 \
--ont_clinical ../data/genomic/tcga_cbio/8_clinical_data.csv --sen ../data/genomic/sensitive.txt \
--ont_genomic ../data/genomic/tcga_cbio/8_mutation_data.csv  \
--clinical ../data/genomic/tcga_cbio/8_clinical_data.csv \
--genomic ../data/genomic/tcga_cbio/8_mutation_data.csv \
--dbHost localhost --dbPort 5432 --dbName i2b2medcosrv2 --dbUser i2b2 --dbPassword i2b2

./medcoLoader -debug 2 v0 -g /home/jagomes/medco-deployment/configuration-profiles/prod/3nodes-samehost/group.toml --entryPointIdx 2 \
--ont_clinical ../data/genomic/tcga_cbio/8_clinical_data.csv --sen ../data/genomic/sensitive.txt \
--ont_genomic ../data/genomic/tcga_cbio/8_mutation_data.csv  \
--clinical ../data/genomic/tcga_cbio/8_clinical_data.csv \
--genomic ../data/genomic/tcga_cbio/8_mutation_data.csv \
--dbHost localhost --dbPort 5432 --dbName i2b2medcosrv2 --dbUser i2b2 --dbPassword i2b2