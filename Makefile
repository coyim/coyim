default: deps lint test

lint:
	golint ./...

test:
	go test -v ./... -cover

ci: default

deps:
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover

cover:
	go test . -coverprofile=coverage.out
	go tool cover -html=coverage.out
