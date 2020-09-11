.PHONY: all
all: fmt lint vet coverage

.PHONY: test
test:
	go test -count=1 -race ./...

.PHONY: coverage-ui
coverage-ui:
	go test -covermode=count -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: coverage
coverage:
	go test -coverprofile coverage.out ./...
	go tool cover -func coverage.out | grep total

.PHONY: lint
lint:
	golint -set_exit_status ./...

.PHONY: vet
lint:
	go vet ./...

.PHONY: fmt
lint:
	go fmt ./...