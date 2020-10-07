MEDCO_VERSION := $(shell git describe --tags --always --dirty)

# build commands
.PHONY: docker_images_build_dev
docker_images_build_dev:
	cd "deployments/dev-local-3nodes"; \
	docker-compose -f docker-compose.yml -f docker-compose.tools.yml build

# test commands
.PHONY: test_go test_go_fmt test_go_lint test_go_unit
test_go: test_go_fmt test_go_lint test_go_unit

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

test_codecov:
	./test/coveralls.sh

# utility commands
.PHONY:	test_unlynx_loop swagger swagger-gen download_test_data version
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

version:
	@echo $(MEDCO_VERSION)