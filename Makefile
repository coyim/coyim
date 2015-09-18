default: deps lint test

build:
	go build -o bin/coyim

lint:
	golint ./...

test:
	go test -v ./... -cover

ci: get default

get:
	go get -t ./...

deps:
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover

cover:
	go test . -coverprofile=coverage.out
	go tool cover -html=coverage.out

