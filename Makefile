PKG            := github.com/504dev/logr
PKG_LIST       := $(shell go list ${PKG}/... | grep -v /vendor/)
SERVICE_SERVER := server

all: setup test build

setup: ## Installing all service dependencies
	echo "Setup..."
	GO111MODULE=on go mod vendor

.PHONY: config
config: ## Creating the local config yml.
	echo "Creating local config.yml ..."
	cp templates/config.yml.template config.yml

env: ## Creating .env file.
	echo "Creating .env file ..."
	cat templates/.env.template | sed "s/OAUTH_JWT_SECRET=/OAUTH_JWT_SECRET=$$(openssl rand -hex 12 | awk '{ print $1 }')/" > .env

build: ## Build the executable file of service.
	echo "Building backend..."
	go build -o logr-server ./cmd/$(SERVICE_SERVER)/main.go

front: ## Run service.
	echo "Building frontend..."
	cd ./frontend && npm i --only=production && npm run build

run: build ## Run service.
	echo "Running..."
	./logr-server

help: ## Display this help screen.
	grep -E '^[a-zA-Z_\-\:]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ": .*?## "}; {gsub(/[\\]*/,""); printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
