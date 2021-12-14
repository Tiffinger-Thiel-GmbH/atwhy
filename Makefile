# Based on https://about.gitlab.com/blog/2017/11/27/go-tools-and-gitlab-how-to-do-continuous-integration-like-a-boss/

PROJECT_NAME := "atwhy"
PKG := "github.com/Tiffinger-Thiel-GmbH/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
GOOS ?= linux
GOARCH ?= amd64

.PHONY: all dep build clean test race coverage coverhtml lint dist

all: build

lint: ## Lint the files
	@go vet	
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@staticcheck

test: ## Run unittests
	@go test -short ${PKG_LIST}

race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	echo TODO: ./tools/coverage.sh; 

coverhtml: ## Generate global code coverage report in HTML
	echo TODO: ./tools/coverage.sh html;

dep: ## Get the dependencies
	@go get -v -d ./...

build: dep ## Build the binary file
	@go build -v $(PKG)

dist: build
ifeq ($(GOOS),windows)
	@zip ${GOOS}-${GOARCH}-atwhy.zip ${PROJECT_NAME}.exe
else
	@chmod +x ${PROJECT_NAME}
	@tar -czf ${GOOS}-${GOARCH}-atwhy.tar.gz ${PROJECT_NAME}
endif
	
clean: ## Remove previous build
	@go clean

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
