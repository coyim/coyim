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
	golint ./...
endif

test:
	go test -cover -v ./...

test-slow:
	make -C ./compat libotr-compat

ci: lint test test-slow

deps:
	echo "BUILDING ON: $(GO_VERSION)"
ifeq ($(GO_VERSION), go1.3)
else ifeq ($(GO_VERSION), go1.4)
else ifeq ($(GO_VERSION), go1.5)
else ifeq ($(GO_VERSION), go1.6)
else ifeq ($(GO_VERSION), go1.7)
else ifeq ($(GO_VERSION), go1.8)
else
	go get -u golang.org/x/lint/golint
endif
	go get golang.org/x/tools/cmd/cover

cover:
	go test . -coverprofile=coverage.out
	go tool cover -html=coverage.out

lint-aggregator:
	go get -u github.com/golangci/golangci-lint/...
	golangci-lint run

gosec:
	go get -u github.com/securego/gosec/cmd/gosec...
	gosec ./...
