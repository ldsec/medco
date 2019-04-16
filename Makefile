EXCLUDE_LINT = "_test.go"

# generate/update go server based on swagger specifications
swagger-gen:
	swagger validate ./swagger/swagger.yml
	swagger generate server \
		--principal=models.User \
		--target=./swagger/ \
		--spec=./swagger/swagger.yml \
		--name=medco-connector

test_lint:
	@echo Checking linting of files
	@{ \
		go install golang.org/x/lint/golint; \
		el=$(EXCLUDE_LINT); \
		lintfiles=$$( golint ./... | egrep -v "$$el" ); \
		if [ -n "$$lintfiles" ]; then \
		echo "Lint errors:"; \
		echo "$$lintfiles"; \
		exit 1; \
		fi \
	}
