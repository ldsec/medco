EXCLUDE_LINT = "_test.go"

# generate/update go server based on swagger specifications
swagger-gen:
	swagger validate ./swagger/medco-connector-server.yml
	swagger generate server \
		--server-package=restapi/server \
		--model-package=restapi/models \
		--principal=models.User \
		--target=./ \
		--spec=./swagger/medco-connector-server.yml \
		--name=medco-connector
	swagger validate ./swagger/medco-cli-client.yml
	swagger generate client \
		--client-package=restapi/client \
		--existing-models=github.com/lca1/medco-connector/restapi/models \
		--skip-models \
		--principal=models.User \
		--target=./ \
		--spec=./swagger/medco-cli-client.yml \
		--name=medco-cli \
		--default-scheme=https

test_lint:
	@echo Checking linting of files
	@{ \
		GO111MODULE=off go get -u golang.org/x/lint/golint; \
		el=$(EXCLUDE_LINT); \
		lintfiles=$$( golint ./... | egrep -v "$$el" ); \
		if [ -n "$$lintfiles" ]; then \
		echo "Lint errors:"; \
		echo "$$lintfiles"; \
		exit 1; \
		fi \
	}
