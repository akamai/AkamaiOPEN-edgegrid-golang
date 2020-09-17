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

.PHONY: coverage-v2
coverage-v2:
	go test -coverprofile coverage.out ./pkg/...
	go tool cover -func coverage.out | grep total

.PHONY: lint
lint:
	golint -set_exit_status ./...

.PHONY: lint-v2
lint-v2:
	golint -set_exit_status ./pkg/...

.PHONY: vet
vet:
	go vet ./...

.PHONY: vet-v2
vet-v2:
	go vet ./pkg/...

.PHONY: fmt
fmt:
	go fmt ./...