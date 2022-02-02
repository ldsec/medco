export MEDCO_VERSION := $(shell make --no-print-directory -C ../../ medco_version)
export GLOWING_BEAR_MEDCO_VERSION := $(shell make --no-print-directory -C ../../ gb_version)

.PHONY: build pull up stop down load_test_data
build:
	docker-compose -f docker-compose.yml -f docker-compose.tools.yml build

pull:
	docker-compose -f docker-compose.yml -f docker-compose.tools.yml pull

up:
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

load_test_data:
	docker-compose -f docker-compose.tools.yml run medco-loader v0 \
		--ont_clinical /data/genomic/tcga_cbio/8_clinical_data.csv --sen /data/genomic/sensitive.txt \
		--ont_genomic /data/genomic/tcga_cbio/8_mutation_data.csv --clinical /data/genomic/tcga_cbio/8_clinical_data.csv \
		--genomic /data/genomic/tcga_cbio/8_mutation_data.csv --output /data/
