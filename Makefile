test_fmt:
	@echo Checking correct formatting of files
	@{ \
		files=$$( go fmt ./... | egrep -v "bindata.go" ); \
		if [ -n "$$files" ]; then \
		echo "Files not properly formatted: $$files"; \
		exit 1; \
		fi; \
		if ! go vet ./...; then \
		exit 1; \
		fi \
	}

test_lint:
	@echo Checking linting of files
	@{ \
		go get -u github.com/golang/lint/golint; \
		lintfiles=$$( golint ./... | egrep -v "bindata.go" ); \
		if [ -n "$$lintfiles" ]; then \
		echo "Lint errors:"; \
		echo "$$lintfiles"; \
		exit 1; \
		fi \
	}

test_go:
	go test -v -race -p=1 ./...;

local: test_fmt test_lint test_go