GIT_VERSION := $(shell git rev-parse HEAD)
GIT_SHORT_VERSION := $(shell git rev-parse --short HEAD)
TAG_VERSION := $(shell git tag -l --contains $$GIT_VERSION | tail -1)
CURRENT_DATE := $(shell date "+%Y-%m-%d")
BUILD_TIMESTAMP := $(shell TZ='GMT' date '+%Y-%m-%d %H:%M:%S')

PKGS := $(shell go list ./... | grep -v /vendor)
SRC_DIRS := . $(addprefix .,$(subst github.com/coyim/constbn,,$(PKGS)))
SRC_TEST := $(foreach sdir,$(SRC_DIRS),$(wildcard $(sdir)/*_test.go))
SRC_ALL := $(foreach sdir,$(SRC_DIRS),$(wildcard $(sdir)/*.go))
SRC := $(filter-out $(SRC_TEST), $(SRC_ALL))

GO := go
GOBUILD := $(GO) build
GOTEST := $(GO) test

default: check
check: quality test

quality: lint gosec ineffassign lint-ci

lint:
	golint -set_exit_status $(SRC_DIRS)

gosec:
	gosec ./...

ineffassign:
	ineffassign .

lint-ci:
	golangci-lint run

test:
	$(GOTEST) -cover -v ./...

run-cover:
	$(GOTEST) -coverprofile=coverage.out -v ./...

coverage: run-cover
	go tool cover -html=coverage.out

coverage-tails: run-cover
	go tool cover -html=coverage.out -o ~/Tor\ Browser/coverage.html
	xdg-open ~/Tor\ Browser/coverage.html
