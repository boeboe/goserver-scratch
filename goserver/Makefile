# Makefile

PROGRAM=goserver

.PHONY: help

help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

build:  ## Build binary
	go build -o .bin/${PROGRAM} main.go

run:  ## Run program
	go run main.go

clean: ## Clean binaries
	rm -rf ./.bin/*

compile:  ## Compile for other architectures
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/arm/${PROGRAM} main.go
	GOOS=linux GOARCH=arm64 go build -o bin/arm64/${PROGRAM} main.go
	GOOS=freebsd GOARCH=386 go build -o bin/386/${PROGRAM} main.go

all: build compile  ## Build for all architectures