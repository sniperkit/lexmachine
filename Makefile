.PHONY: all test clean man glide fast release install

# env
GO15VENDOREXPERIMENT=1

# app
WRK_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

all: lexc

lexc: deps test
	@go build -ldflags "-X main.VERSION=`cat VERSION`" -o ./bin/lexc ./cmd/lexc/*.go

fast:
	@go build -i -ldflags "-X main.VERSION=`cat VERSION`-dev" -o ./build/textql ./textql/main.go

deps: glide
	@glide install

test:
	@go test ./pkg/...

clean:
	@rm -fr ./dist
	@rm -fr ./bin

release: lexc
	@git tag -a `cat VERSION`
	@git push origin `cat VERSION`

install: deps test
	@go install -ldflags "-X main.VERSION=`cat VERSION`" ./cmd/lexc/*.go

