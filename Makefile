# Makefile for managing the application

.PHONY: build all init docker-up docker-down generated

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o main ./cmd

clean:
	rm -rf generated

init: clean generated
	go mod tidy
	go mod vendor

docker-up: generated
	docker-compose up --build -d

docker-down:
	docker-compose down --volumes

test:
	go clean -testcache
	go test -short -cover -coverprofile=coverage.out ./handler ./repository ./tests
	go tool cover -html=coverage.out -o coverage.html

test_api:
	go clean -testcache
	go test ./tests/...

generate: generated generate_mocks

generated: api.yml
	@echo "Generating files..."
	mkdir -p generated
	oapi-codegen --package generated -generate types,server,spec api.yml > generated/api.gen.go

INTERFACES_GO_FILES := $(shell find repository -name "interfaces.go")
INTERFACES_GEN_GO_FILES := $(INTERFACES_GO_FILES:%.go=mocks/%.mock.gen.go)

generate_mocks: $(INTERFACES_GEN_GO_FILES)
mocks/%.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))
