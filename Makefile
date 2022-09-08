BIN = $(CURDIR)/bin
GO  = go
Q   = $(if $(filter 1,$V),,@)

clean-tools:
	@rm -rf $(BIN)/go*

$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) Building $(PACKAGE)...)
	$Q tmp=$$(mktemp -d); \
	   env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) $(GO) get $(PACKAGE) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

GOIMPORTS = $(BIN)/goimports
$(BIN)/goimports: PACKAGE=golang.org/x/tools/cmd/goimports
	   

.PHONY: all
all: fmt-check test

.PHONY: test
test: ; $(info > Running tests...) @
	go test -v  ./...

.PHONY: fmt
fmt:  | $(GOIMPORTS); $(info $(M) Running goimports...) @
	$Q $(GOIMPORTS) -w .

.PHONY: fmt-check
fmt-check: | $(GOIMPORTS); $(info $(M) Running format and imports check...) @
	$(eval OUTPUT = $(shell $(GOIMPORTS) -l .))
	@if [ "$(OUTPUT)" != "" ]; then\
		echo "Found following files with incorrect format and/or imports:";\
		echo "$(OUTPUT)";\
		false;\
	fi