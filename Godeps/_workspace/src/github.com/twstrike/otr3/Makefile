GO_VERSION=$(shell go version | grep  -o 'go[[:digit:]]\.[[:digit:]]')


default: deps lint test

lint:
ifeq ($(GO_VERSION), go1.3)
	echo "Your version of Go is too old for running lint"
else
ifeq ($(GO_VERSION), go1.4)
	echo "Your version of Go is too old for running lint"
else
	golint ./...
endif
endif

test:
	go test -cover -v ./...

test-slow:
	make -C ./compat libotr-compat

ci: lint test test-slow

deps:
ifeq ($(GO_VERSION), go1.3)
else
ifeq ($(GO_VERSION), go1.4)
else
	go get github.com/golang/lint/golint
endif
endif
	go get golang.org/x/tools/cmd/cover

cover:
	go test . -coverprofile=coverage.out
	go tool cover -html=coverage.out
