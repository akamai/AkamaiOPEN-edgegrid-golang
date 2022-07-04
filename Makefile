.PHONY: all
all: test

.PHONY: test
test: ; $(info > Running tests...) @
	go test -v  ./...
