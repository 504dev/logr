PKG          := github.com/504dev/kidlog
PKG_LIST     := $(shell go list ${PKG}/... | grep -v /vendor/)

all: setup test build

setup: ## Installing all service dependencies
	echo "Setup..."
	GO111MODULE=on go mod vendor

.PHONY: config
config: ## Creating the local config yml.
	echo "Creating local config yml ..."
	cp config.example.yml config.yml

build\:server: ## Build the executable file of service
	echo "Building..."
	cd cmd/server && go build

build\:server\:env: ## Build the executable file of server for production env
	echo "Building server..."
	cd cmd/$(SERVER_NAME) && env GOOS=linux GOARCH=amd64 go build

test: ## Run tests for all packages
	echo "Testing..."
	go test -race -count=1 ${PKG_LIST}

coverage: ## Calculating code test coverage
	echo "Calculating coverage..."
	PKG=$(PKG) ./tools/coverage.sh

clean: ## Cleans the binary files and etc
	echo "Clean..."
	rm -f cmd/server/server

help: ## Display this help screen
	grep -E '^[a-zA-Z_\-\:]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ": .*?## "}; {gsub(/[\\]*/,""); printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
