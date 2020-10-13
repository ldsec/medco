MEDCO_VERSION := $(shell git describe --tags --always --dirty)

# build commands
.PHONY: docker_images_build_dev
docker_images_build_dev:
	cd "deployments/dev-local-3nodes"; \
	docker-compose -f docker-compose.yml -f docker-compose.tools.yml build

# test commands
.PHONY: test test_go_fmt test_go_lint test_codecov_unit test_codecov_e2e
test_local: test_go_fmt test_go_lint test_go_unit

test_go_fmt:
	@echo Checking correct formatting of files
	@{ \
		files=$$( go fmt ./... ); \
		if [ -n "$$files" ]; then \
		echo "Files not properly formatted: $$files"; \
		exit 1; \
		fi; \
		if ! go vet ./...; then \
		exit 1; \
		fi \
	}

test_go_lint:
	@echo Checking linting of files
	@{ \
		GO111MODULE=off go get -u golang.org/x/lint/golint; \
		el="_test.go"; \
		lintfiles=$$( golint ./... | egrep -v "$$el" ); \
		if [ -n "$$lintfiles" ]; then \
		echo "Lint errors:"; \
		echo "$$lintfiles"; \
		exit 1; \
		fi \
	}

test_go_unit:
	go test -v -race -short -p=1 ./...

test_codecov_unit:
	./test/coveralls.sh "./connector/wrappers/i2b2 ./connector/wrappers/unlynx ./connector/server/handlers"

test_codecov_e2e:
	./test/coveralls.sh "" "./connector/wrappers/i2b2 ./connector/wrappers/unlynx ./connector/server/handlers"

# utility commands
.PHONY:	test_unlynx_loop swagger swagger-gen download_test_data load_test_data version
test_unlynx_loop:
	for i in $$( seq 100 ); \
		do echo "******* Run $$i"; echo; \
		go test -v -short -p=1 -run Agg -count 10 ./unlynx/services/ > run.log || \
		( cat run.log; exit 1 ) || exit 1; \
	done

swagger-gen: swagger
	swagger validate ./connector/swagger/medco-connector.yml
	swagger generate server \
		--server-package=connector/restapi/server \
		--model-package=connector/restapi/models \
		--principal=github.com/ldsec/medco/connector/restapi/models.User \
		--target=./ \
		--spec=./connector/swagger/medco-connector.yml \
		--name=medco-connector
	swagger generate client \
		--client-package=connector/restapi/client \
		--existing-models=github.com/ldsec/medco/connector/restapi/models \
		--skip-models \
		--principal=github.com/ldsec/medco/connector/restapi/models.User \
		--target=./ \
		--spec=./connector/swagger/medco-connector.yml \
		--name=medco-cli \
		--default-scheme=https

swagger:
	@if ! which swagger >/dev/null; then \
		go install github.com/go-swagger/go-swagger/cmd/swagger && \
		echo "swagger installed"; \
	fi

download_test_data:
	./test/data/download.sh genomic_small
	./test/data/download.sh i2b2

load_test_data:
	docker-compose -f deployments/dev-local-3nodes/docker-compose.tools.yml run medco-loader-srv0 v0 --ont_clinical /data/genomic/tcga_cbio/8_clinical_data.csv --sen /data/genomic/sensitive.txt --ont_genomic /data/genomic/tcga_cbio/8_mutation_data.csv --clinical /data/genomic/tcga_cbio/8_clinical_data.csv --genomic /data/genomic/tcga_cbio/8_mutation_data.csv --output /data/
	docker-compose -f deployments/dev-local-3nodes/docker-compose.tools.yml run medco-loader-srv1 v0 --ont_clinical /data/genomic/tcga_cbio/8_clinical_data.csv --sen /data/genomic/sensitive.txt --ont_genomic /data/genomic/tcga_cbio/8_mutation_data.csv --clinical /data/genomic/tcga_cbio/8_clinical_data.csv --genomic /data/genomic/tcga_cbio/8_mutation_data.csv --output /data/
	docker-compose -f deployments/dev-local-3nodes/docker-compose.tools.yml run medco-loader-srv2 v0 --ont_clinical /data/genomic/tcga_cbio/8_clinical_data.csv --sen /data/genomic/sensitive.txt --ont_genomic /data/genomic/tcga_cbio/8_mutation_data.csv --clinical /data/genomic/tcga_cbio/8_clinical_data.csv --genomic /data/genomic/tcga_cbio/8_mutation_data.csv --output /data/

version:
	@echo $(MEDCO_VERSION)
