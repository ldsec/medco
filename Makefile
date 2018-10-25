EXCLUDE_LINT = "_test.go"

test_fmt:
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

test_lint:
	@echo Checking linting of files
	@{ \
		go get -u golang.org/x/lint/golint; \
		el=$(EXCLUDE_LINT); \
		lintfiles=$$( golint ./... | egrep -v "$$el" ); \
		if [ -n "$$lintfiles" ]; then \
		echo "Lint errors:"; \
		echo "$$lintfiles"; \
		exit 1; \
		fi \
	}

test_verbose:
	go test -v -race -short -p=1 ./...

test_playground:
	cd protocols; \
	for a in $$( seq 10 ); do \
	  go test -v -race -p=1 || exit 1 ; \
	done;

test_goveralls:
	go get github.com/mattn/goveralls
	./coveralls.sh

test: test_fmt test_lint test_goveralls

local: test_fmt test_lint test_verbose

all: install test