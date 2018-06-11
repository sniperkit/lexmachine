.PHONY: all test clean man glide fast release install

# env
GO15VENDOREXPERIMENT=1

# app
PROG_NAME ?= lexc

# dirs
DIST_DIR ?= ./dist
WRK_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

# vcs
GIT_BRANCH := $(subst heads/,,$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null))

# pkgs
SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')
PACKAGES = $(shell go list ./... | grep -v /vendor/)

# VERSION = $(shell cat "$(WRK_DIR)/VERSION" | tr '\n' '')
VERSION ?= $(shell git describe --tags)
VERSION_INCODE = $(shell perl -ne '/^var version.*"([^"]+)".*$$/ && print "v$$1\n"' main.go)
VERSION_INCHANGELOG = $(shell perl -ne '/^\# Release (\d+(\.\d+)+) / && print "$$1\n"' CHANGELOG.md | head -n1)

VCS_GIT_REMOTE_URL = $(shell git config --get remote.origin.url)
VCS_GIT_VERSION ?= $(VERSION)

BUILD_LDFLAGS := "-X lexmachine.Version=\"$(VERSION)\""

all: prepare build install examples
 
prepare: deps test

version:
	echo "ldflags=$(BUILD_LDFLAGS)"

build: version
	echo "ldflags=$(BUILD_LDFLAGS)"
	@go build -o ./bin/$(PROG_NAME) -ldflags "-X pkg.Version=\"$(VERSION)\"" ./cmd/$(PROG_NAME)/*.go
	@./bin/$(PROG_NAME) --version

install: version
	@go install -ldflags "$(BUILD_LDFLAGS)" ./cmd/$(PROG_NAME)/*.go
	@$(PROG_NAME) --version

fast: version
	@go build -i -ldflags "$(BUILD_LDFLAGS)-dev" -o ./bin/$(PROG_NAME) ./cmd/$(PROG_NAME)/*.go
	@./bin/$(PROG_NAME) --version

examples: version
	@go build -ldflags "$(BUILD_LDFLAGS)" -o ./bin/sensors ./examples/sensors/*.go
	@go build -ldflags "$(BUILD_LDFLAGS)" -o ./bin/sensors-parser ./examples/sensors-parser/*.go

glide:
	@go get -v -u github.com/Masterminds/glide

deps: glide
	@glide install --strip-vendor

test:
	@go test -cover $(PACKAGES)

clean:
	@go clean
	@rm -fr ./dist
	@rm -fr ./bin

release: lexc
	@git tag -a `cat VERSION`
	@git push origin `cat VERSION`

check: check-deps vet lint errcheck interfacer aligncheck structcheck varcheck unconvert gosimple staticcheck unused vendorcheck prealloc test

vet:
	@go vet $(PACKAGES)

lint:
	@golint -set_exit_status $(PACKAGES)

errcheck:
	@errcheck $(PACKAGES)

interfacer:
	@interfacer $(PACKAGES)

aligncheck:
	@aligncheck $(PACKAGES)

structcheck:
	@structcheck $(PACKAGES)

varcheck:
	@varcheck $(PACKAGES)

unconvert:
	@unconvert -v $(PACKAGES)

gosimple:
	@gosimple $(PACKAGES)

staticcheck:
	@staticcheck $(PACKAGES)

unused:
	@unused $(PACKAGES)

vendorcheck:
	@vendorcheck $(PACKAGES)
	@vendorcheck -u $(PACKAGES)

prealloc:
	@prealloc $(PACKAGES)

coverage:
	@echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	@go tool cover -html=coverage-all.out

check-deps:
	@go get -v -u github.com/alexkohler/prealloc
	@go get -v -u github.com/FiloSottile/vendorcheck
	@go get -v -u github.com/golang/dep/cmd/dep
	@go get -v -u github.com/golang/lint/golint
	@go get -v -u github.com/kisielk/errcheck
	@go get -v -u github.com/mdempsky/unconvert
	@go get -v -u github.com/opennota/check/...
	@go get -v -u honnef.co/go/tools/...
	@go get -v -u mvdan.cc/interfacer