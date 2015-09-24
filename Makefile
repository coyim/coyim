default: deps lint test

build: build-cli build-gui

build-gui:
	go build -tags nocli  -o bin/coyim-gui

build-cli:
	go build -o bin/coyim

lint:
	golint ./...

test:
	go test -cover -v ./...

ci: get default

get:
	go get -t ./...

deps:
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover

cover:
	go test . -coverprofile=coverage.out
	go tool cover -html=coverage.out

