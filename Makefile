# Adapted from https://www.thapaliya.com/en/writings/well-documented-makefiles/
.PHONY: help
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.DEFAULT_GOAL:=help

# Read the version from the build/VERSION file
VERSION := $(shell cat build/VERSION)

.PHONY: version
version: ## Print the current app version
	@echo $(VERSION)

.PHONY: build
build: ## Build the inkfem binary
	go build -o bin/inkfem inkfem.go 
	
.PHONY: test
test: ## Run all the tests
	go test -v ./...
	
.PHONY: run
run: ## Run the inkfem binary
	go run inkfem.go

.PHONY: bench
bench: ## Run all the benchmarks
	go test -benchmem -run=^$$ -bench '^BenchmarkSolveStructure$$' github.com/angelsolaorbaiceta/inkfem/tests -count 4
	
.PHONY: fmt
fmt: ## Run go fmt
	go fmt ./...