# Variables
PROJECT = endpoint
VERSION = 1.0.0

clean:
	go clean -mod cache
dep:
	go get ./...
	go mod download
	go mod vendor
	go mod tidy

CGO_ENABLED=0
build:
	go build -mod=vendor -a -ldflags '-w -extldflags "-static" -X main.BuildVersion=$(VERSION)' -o bin/endpoint .

TEST_TIMEOUT = 60
TEST_PKGS ?= .
TEST_TARGETS := test-short test-verbose test-race test-cover
.PHONY: $(TEST_TARGETS) test
test-short:   TEST_ARGS=-short
test-verbose: TEST_ARGS=-v
test-race:    TEST_ARGS=-race
test-cover:   TEST_ARGS=-cover
$(TEST_TARGETS): test

test:
	go test -mod=vendor -timeout $(TEST_TIMEOUT)s $(TEST_ARGS) $(TEST_PKGS)

docker.build:
	docker build --build-arg VERSION=$(VERSION) --no-cache --tag endpoint .

