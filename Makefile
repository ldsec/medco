MEDCO_VERSION := $(shell scripts/version.sh)
GB_VERSION := v3.0.0
DOCKER_REGISTRY := ghcr.io/chuv-ds

# test commands
.PHONY: test test_go_fmt test_go_lint test_codecov_unit test_codecov_e2e
test_local: test_go_fmt test_go_lint test_go_unit

test_go_fmt:
	@echo Checking correct formatting of files
	@{ \
  		GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports; \
		files=$$( goimports -w -l . ); \
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
	go test -v -race -short --tags=unit_test -p=1 ./...

test_go_integration:
	go test -v -race -short --tags=integration_test -p=1 ./...

# utility commands
.PHONY:	test_unlynx_loop swagger swagger-gen download_test_data medco_version gb_version
test_unlynx_loop:
	for i in $$( seq 100 ); \
		do echo "******* Run $$i"; echo; \
		go test -v -short -p=1 -run Agg -count 10 ./unlynx/services/ > run.log || \
		( cat run.log; exit 1 ) || exit 1; \
	done

swagger-gen: swagger
	swagger validate ./connector/swagger/medco-connector.yml
	rm -rf connector/restapi/client/* connector/restapi/models/*
	find connector/restapi/server/ -mindepth 1 -type f ! -name "configure_medco_connector.go" -delete
	find connector/restapi/server/ -mindepth 1 -type d -delete
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

medco_version:
	@echo $(MEDCO_VERSION)

gb_version:
	@echo $(GB_VERSION)

docker_registry:
	@echo $(DOCKER_REGISTRY)
