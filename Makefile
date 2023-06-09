SHELL=/bin/bash -e -o pipefail
PWD = $(shell pwd)
GO_BUILD= go build
GOFLAGS= CGO_ENABLED=0
REGISTRY_IMAGE_NAME=go-movies

## help: Print this help message
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

## test: Run tests and show coverage result
.PHONY: test
test:
	go test -race -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

## tidy: Cleanup and download missing dependencies
.PHONY: tidy
tidy:
	go mod tidy
	go mod verify

## vet: Examine Go source code and reports suspicious constructs
.PHONY: vet
vet:
	go vet ./...

## fmt: Format all go source files
.PHONY: fmt
fmt:
	go fmt ./...

## build: Build binary into bin/ directory
.PHONY: build
build:
	$(GOFLAGS) $(GO_BUILD) -a -v -ldflags="-w -s" -o bin/app cmd/main.go

## docker-build: Builds a docker image
.PHONY: docker-build
docker-build:
	docker build . -t $(REGISTRY_IMAGE_NAME)

## docker-run: Runs the docker image built by [make docker-build]
.PHONY: docker-run
docker-run:
	docker run $(REGISTRY_IMAGE_NAME)

## docker: Builds and runs the docker image
.PHONY: docker
docker: docker-build docker-run
