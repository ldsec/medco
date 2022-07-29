SHELL := /bin/bash

.PHONY: build pull up up-d stop down docker-compose docker-compose.tools load_test_data
build:
	docker-compose -f docker-compose.yml -f docker-compose.tools.yml build

pull:
	docker-compose -f docker-compose.yml -f docker-compose.tools.yml pull

up:
	docker-compose -f docker-compose.yml up

up-d:
	docker-compose -f docker-compose.yml up -d

stop:
	docker-compose -f docker-compose.yml stop

down:
	docker-compose -f docker-compose.yml down

# use environment variable ARGS to pass arguments to docker-compose: make docker-compose ARGS="up -d"
docker-compose:
	docker-compose -f docker-compose.yml $(ARGS)

# use environment variable ARGS to pass arguments to docker-compose.tools: make docker-compose.tools ARGS="run ..."
docker-compose.tools:
	docker-compose -f docker-compose.tools.yml $(ARGS)

run_test_query:
	docker-compose -f docker-compose.tools.yml run medco-cli-client --user test --password test query enc::1 OR enc::2 AND enc::3

load_test_data:
	docker-compose -f docker-compose.tools.yml run medco-loader-srv0 v0 \
		--ont_clinical /data/genomic/tcga_cbio/8_clinical_data.csv --sen /data/genomic/sensitive.txt \
		--ont_genomic /data/genomic/tcga_cbio/8_mutation_data.csv --clinical /data/genomic/tcga_cbio/8_clinical_data.csv \
		--genomic /data/genomic/tcga_cbio/8_mutation_data.csv --output /data/
	docker-compose -f docker-compose.tools.yml run medco-loader-srv1 v0 \
		--ont_clinical /data/genomic/tcga_cbio/8_clinical_data.csv --sen /data/genomic/sensitive.txt \
		--ont_genomic /data/genomic/tcga_cbio/8_mutation_data.csv --clinical /data/genomic/tcga_cbio/8_clinical_data.csv \
		--genomic /data/genomic/tcga_cbio/8_mutation_data.csv --output /data/
	docker-compose -f docker-compose.tools.yml run medco-loader-srv2 v0 \
		--ont_clinical /data/genomic/tcga_cbio/8_clinical_data.csv --sen /data/genomic/sensitive.txt \
		--ont_genomic /data/genomic/tcga_cbio/8_mutation_data.csv --clinical /data/genomic/tcga_cbio/8_clinical_data.csv \
		--genomic /data/genomic/tcga_cbio/8_mutation_data.csv --output /data/
