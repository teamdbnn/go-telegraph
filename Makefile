SHELL := /bin/bash
MAKEFILE_PATH := $(abspath $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
PATH := $(MAKEFILE_PATH):$(PATH)

#
PKGS = $(shell go list ./...)

# Colors
GREEN_COLOR   = "\033[0;32m"
PURPLE_COLOR  = "\033[0;35m"
DEFAULT_COLOR = "\033[m"

.PHONY: init help test coverage lint format version

all: format lint test

help: ## Show this help screen
	@printf 'Usage: make \033[36m<TARGETS>\033[0m ... \033[36m<OPTIONS>\033[0m\n\nAvailable targets are:'
	@awk 'BEGIN {FS = ":.*##"; printf "\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "    \033[36m%-17s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@printf "\nTargets run by default are: clean proto build format lint test.\n"

init: ## Install required tools
	@echo -e $(GREEN_COLOR)[INIT]$(DEFAULT_COLOR)
	@cd tools && go generate -x -tags=tools

test: ## Run unit TESTS
	@echo -e $(GREEN_COLOR)[TEST]$(DEFAULT_COLOR)
	@go test -race $(PKGS)

deps: ## Download required dependencies and remove unused
	@echo -e $(GREEN_COLOR)[RESOLVE DEPENDENCIES]$(DEFAULT_COLOR)
	go mod tidy

update: ## Update dependencies
	@echo -e $(GREEN_COLOR)[UPDATE DEPENDENCIES]$(DEFAULT_COLOR)
	go get -u

coverage: ## Report code TESTS coverage
	@echo -e $(GREEN_COLOR)[COVERAGE]$(DEFAULT_COLOR)
	@# Create the coverage files directory
	@mkdir -p $(COVERAGE_DIR)/
	@mkdir -p $(SUBCOV_DIR)/
	@# Create a coverage file for each package
	@for package in $(PKGS); do $(GOTEST) -covermode=count -coverprofile $(SUBCOV_DIR)/`basename "$$package"`.cov "$$package"; done
	@# Merge the coverage profile files
	@echo 'mode: count' > $(COVERAGE_DIR)/coverage.cov ;
	@tail -q -n +2 $(SUBCOV_DIR)/*.cov >> $(COVERAGE_DIR)/coverage.cov ;
	@go tool cover -func=$(COVERAGE_DIR)/coverage.cov ;
	@# If needed, generate HTML report
	@if [ $(html) ]; then go tool cover -html=$(COVERAGE_DIR)/coverage.cov -o coverage.html ; fi
	@# Remove the coverage files directory
	@rm -rf $(COVERAGE_DIR);

lint: ## Run linter on package sources
	@echo -e $(GREEN_COLOR)[LINT]$(DEFAULT_COLOR)
	@bin/golangci-lint run --config=$(MAKEFILE_PATH)/.golangci.yml

format: ## Format project sources
	@echo -e $(GREEN_COLOR)[FORMAT]$(DEFAULT_COLOR)
	@bin/gofumports -l -w .
	@bin/gci write --skip-generated --section Standard --section Default --section "Prefix(github.com/teamdbnn/go-telegraph)" .

version: ## Print Go version
	@echo -e $(GREEN_COLOR)[VERSION]$(DEFAULT_COLOR)
	@go version
