GO_VERSION=$(shell go version | grep  -o 'go[[:digit:]]\.[[:digit:]]*')

default: deps lint test

lint:
ifeq ($(GO_VERSION), go1.3)
else ifeq ($(GO_VERSION), go1.4)
else ifeq ($(GO_VERSION), go1.5)
else ifeq ($(GO_VERSION), go1.6)
else ifeq ($(GO_VERSION), go1.7)
else ifeq ($(GO_VERSION), go1.8)
else
	golint . ./compat ./sexp
endif

test:
	go test -cover -v ./...

test-slow:
	make -C ./compat libotr-compat

ci: lint test test-slow

deps:
ifeq ($(GO_VERSION), go1.3)
else ifeq ($(GO_VERSION), go1.4)
else ifeq ($(GO_VERSION), go1.5)
else ifeq ($(GO_VERSION), go1.6)
else ifeq ($(GO_VERSION), go1.7)
else ifeq ($(GO_VERSION), go1.8)
else
	go get golang.org/x/lint/golint
endif
	go get golang.org/x/tools/cmd/cover
#	go get github.com/golangci/golangci-lint/...
#	go get github.com/securego/gosec/cmd/gosec...

deps-ci: deps
	go get -u github.com/mattn/goveralls

run-cover:
	go test . -coverprofile=coverage.out

coveralls: run-cover
	goveralls -coverprofile=coverage.out

cover: run-cover
	go tool cover -html=coverage.out

lint-aggregator:
	golangci-lint run

gosec:
	gosec ./...
