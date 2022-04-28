MODULE   = $(shell $(GO) list -m)
COMMIT_SHA=$(shell git rev-parse --short HEAD)
VERSION ?= $(shell git describe --tags --always | grep '^v\d' || \
			echo $(FILEVERSION)-$(COMMIT_SHA))
BIN      = $(CURDIR)/bin
GOLANGCI_LINT_VERSION = v1.41.1
GO      = go
TIMEOUT = 15
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell echo ">")

clean-tools:
	@rm -rf $(BIN)/go*

$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) Building $(PACKAGE)...)
	$Q tmp=$$(mktemp -d); \
	   env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) $(GO) get $(PACKAGE) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

GOLINT = $(BIN)/golint
$(BIN)/golint: PACKAGE=golang.org/x/lint/golint

GOCOV = $(BIN)/gocov
$(BIN)/gocov: PACKAGE=github.com/axw/gocov/...

GOCOVXML = $(BIN)/gocov-xml
$(BIN)/gocov-xml: PACKAGE=github.com/AlekSi/gocov-xml

GOJUNITREPORT = $(BIN)/go-junit-report
$(BIN)/go-junit-report: PACKAGE=github.com/jstemmer/go-junit-report

GOIMPORTS = $(BIN)/goimports
$(BIN)/goimports: PACKAGE=golang.org/x/tools/cmd/goimports

GOLANGCILINT = $(BIN)/golangci-lint
$(BIN)/golangci-lint: ; $(info $(M) Installing golangci-lint...) @
	$Q curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN) $(GOLANGCI_LINT_VERSION)


.PHONY: all
all: clean fmt-check lint test-verbose create-junit-report create-coverage-files clean-tools

# Tests

TEST_TARGETS := test-default test-bench test-short test-verbose test-race
.PHONY: $(TEST_TARGETS) check test tests
test-bench:   ARGS=-run=__absolutelynothing__ -bench=. ## Run benchmarks
test-short:   ARGS=-short        ## Run only short tests
test-verbose: ARGS=-v            ## Run tests in verbose mode with coverage reporting
test-race:    ARGS=-race         ## Run tests with race detector

COVERAGE_MODE    = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML     = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML    = $(COVERAGE_DIR)/index.html

$(TEST_TARGETS): COVERAGE_DIR := $(CURDIR)/test/coverage
$(TEST_TARGETS): ; $(info $(M) Running tests with coverage...) @ ## Run coverage tests
	$Q mkdir -p $(COVERAGE_DIR)
	$Q mkdir -p test
	$Q $(GO) test -timeout $(TIMEOUT)s $(ARGS) \
		-coverpkg=./... \
		-covermode=$(COVERAGE_MODE) \
		-coverprofile="$(COVERAGE_PROFILE)" ./... | tee test/tests.output

.PHONY: create-junit-report
create-junit-report: | $(GOJUNITREPORT) ; $(info $(M) Creating juint xml report) @
	$Q cat $(CURDIR)/test/tests.output | $(GOJUNITREPORT) > $(CURDIR)/test/tests.xml
	$Q sed -i -e 's/skip=/skipped=/g' $(CURDIR)/test/tests.xml

.PHONY: create-coverage-files
create-coverage-files: COVERAGE_DIR := $(CURDIR)/test/coverage
create-coverage-files: | $(GOCOV) $(GOCOVXML); $(info $(M) Creating coverage files...) @ ## Run coverage tests
	$Q $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$Q $(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)

.PHONY: lint
lint: | $(GOLANGCILINT) ; $(info $(M) Running golangci-lint...) @
	$Q $(BIN)/golangci-lint run

.PHONY: fmt
fmt:  | $(GOIMPORTS); $(info $(M) Running goimports...) @ ## Run goimports on all source files
	$Q $(GOIMPORTS) -w .

.PHONY: fmt-check
fmt-check: | $(GOIMPORTS); $(info $(M) Running format and imports check...) @ ## Run goimports on all source files
	$(eval OUTPUT = $(shell $(GOIMPORTS) -l .))
	@if [ "$(OUTPUT)" != "" ]; then\
		echo "Found following files with incorrect format and/or imports:";\
		echo "$(OUTPUT)";\
		false;\
	fi

# Misc

.PHONY: clean
clean: ; $(info $(M) Cleaning...)	@ ## Cleanup everything
	@rm -rf $(BIN)
	@rm -rf test/tests.* test/coverage

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)
