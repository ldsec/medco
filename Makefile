EXCLUDE_LINT = "_test.go"

# generate/update go server based on swagger specifications
swagger-gen:
	swagger validate ./swagger/swagger.yml
	swagger generate server \
		--principal=models.User \
		--target=./swagger/ \
		--spec=./swagger/swagger.yml \
		--name=medco-connector
	swagger generate client \
		--principal=models.User \
		--target=./swagger/ \
		--spec=./swagger/swagger.yml \
		--name=medco-cli-client \
		--existing-models=github.com/lca1/medco-connector/swagger/models \
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
