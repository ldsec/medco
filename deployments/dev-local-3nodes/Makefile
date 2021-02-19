.PHONY: build up stop down load_test_data
build:
	docker-compose -f docker-compose.yml -f docker-compose.tools.yml build

up:
	docker-compose -f docker-compose.yml up

up-d:
	docker-compose -f docker-compose.yml up -d

debug:
	docker-compose -f docker-compose.yml up > logfile
	
stop:
	docker-compose -f docker-compose.yml stop

down:
	docker-compose -f docker-compose.yml down

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
