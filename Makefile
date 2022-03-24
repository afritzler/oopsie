BINARY_NAME=oopsie

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Image URL to use all building/pushing image targets
IMG ?= oopsie:latest

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

all: test build

build: ## Build the oopsie binary.
		go build -o $(BINARY_NAME) -v main.go

test: ## Run tests.
		go test -v ./...

clean: ## Clean build artefacts.
		go clean
		rm -f $(BINARY_NAME)

lint: ## Run golangci-lint against code.
		golangci-lint run ./...

run: ## Run oopsie.
		$(GOBUILD) -o $(BINARY_NAME) -v main.go
		./$(BINARY_NAME)

docker-build: test ## Build docker image with the manager.
		docker build -t ${IMG} .
