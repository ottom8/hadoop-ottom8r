# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.PHONY: all build package lint test gotest help
.DEFAULT_GOAL := help

REPO = $(eval REPO := $$(shell go list -f '{{.ImportPath}}' .))$(value REPO)
NAME = $(eval NAME := $$(shell basename ${REPO}))$(value NAME)
VERSION = $(eval VERSION := $$(shell git describe --dirty))$(value VERSION)
PACKAGE_VERSION = $(eval PACKAGE_VERSION := $$(subst -dirty,,$${VERSION}))$(value PACKAGE_VERSION)
# LDFLAGS = -ldflags "-X ${REPO}/ottom8lib.Version=${VERSION}"
LDFLAGS = 

# sources to evaluate
SRCDIRS := $(shell find . -maxdepth 1 -mindepth 1 -type d -not -path './vendor')
SRCFILES := main.go $(shell find ${SRCDIRS} -name '*.go')

hadoop-ottom8r: vendor ${SRCFILES} ## Build hadoop-ottom8r binary
	go build ${LDFLAGS}

# GLIDE := $(shell command -v glide 2> /dev/null)
	# ifndef GLIDE
	# 	$(error "glide is not available. Install using `curl https://glide.sh/get | sh`")
	# endif
vendor: glide.yaml ## Install vendor dependencies
	glide install -v

test: ## Run all tests
	vendor gotest blackbox

gotest: ## Run Go tests
	go test $(shell glide novendor)

LINTDIRS = $(eval LINTDIRS := $(shell find ${SRCDIRS} -type d -not -path './rpc/pb' -not -path './docs*'))$(value LINTDIRS)
lint: ## Run linters
	@echo '=== golint ==='
	@for dir in ${LINTDIRS}; do golint $${dir}; done # github.com/golang/lint/golint

	@echo '=== gosimple ==='
	@gosimple ${LINTDIRS} # honnef.co/go/simple/cmd/gosimple

	@echo '=== unconvert ==='
	@unconvert ${LINTDIRS} # github.com/mdempsky/unconvert

	@echo '=== structcheck ==='
	@structcheck ${LINTDIRS} # github.com/opennota/check/cmd/structcheck

	@echo '=== varcheck ==='
	@varcheck ${LINTDIRS} # github.com/opennota/check/cmd/varcheck

	@echo '=== gas ==='
	@gas ${LINTDIRS} # github.com/HewlettPackard/gas

# packaging
	# GOX := $(shell command -v gox 2> /dev/null)
	# ifndef GOX
	# 	$(error "gox is not available. Install using `go get github.com/mitchellh/gox`")
	# endif
build: vendor test ## Build hadoop-ottom8r
	@echo "set version to ${PACKAGE_VERSION}"

	@rm -rf build/
	@mkdir -p build/
	gox \
			-osarch="darwin/386" \
			-osarch="darwin/amd64" \
			-os="linux" \
			-output="build/${NAME}_${PACKAGE_VERSION}_{{.OS}}_{{.Arch}}/${NAME}"
	find build -type f -execdir /bin/bash -c 'shasum -a 256 $$0 > $$0.sha256sum' \{\} \;

package: build ## Build and package hadoop-ottom8r
	@mkdir -p build/tgz
	for f in $(shell find build -name ${NAME} | cut -d/ -f2); do \
			(cd $(shell pwd)/build/$$f && tar -zcvf ../tgz/$$f.tar.gz *); \
		done
	(cd build/tgz; shasum -a 512 * > tgz.sha512sum)

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'