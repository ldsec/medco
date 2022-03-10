SHELL := /bin/bash

.PHONY: build pull up stop down docker-compose docker-compose.tools load_test_data
build:
	source ../../scripts/versions.sh && \
	docker-compose -f docker-compose.yml -f docker-compose.tools.yml build

pull:
	source ../../scripts/versions.sh && \
	docker-compose -f docker-compose.yml -f docker-compose.tools.yml pull

up:
	source ../../scripts/versions.sh && \
	docker-compose -f docker-compose.yml up -d

stop:
	source ../../scripts/versions.sh && \
	docker-compose -f docker-compose.yml stop

down:
	source ../../scripts/versions.sh && \
	docker-compose -f docker-compose.yml down

# use environment variable ARGS to pass arguments to docker-compose: make docker-compose ARGS="up -d"
docker-compose:
	source ../../scripts/versions.sh && \
	docker-compose -f docker-compose.yml $(ARGS)

# use environment variable ARGS to pass arguments to docker-compose.tools: make docker-compose.tools ARGS="run ..."
docker-compose.tools:
	source ../../scripts/versions.sh && \
	docker-compose -f docker-compose.tools.yml $(ARGS)

load_test_data:
	source ../../scripts/versions.sh && \
	docker-compose -f docker-compose.tools.yml run medco-loader v0 \
		--ont_clinical /data/genomic/tcga_cbio/8_clinical_data.csv --sen /data/genomic/sensitive.txt \
		--ont_genomic /data/genomic/tcga_cbio/8_mutation_data.csv --clinical /data/genomic/tcga_cbio/8_clinical_data.csv \
		--genomic /data/genomic/tcga_cbio/8_mutation_data.csv --output /data/
