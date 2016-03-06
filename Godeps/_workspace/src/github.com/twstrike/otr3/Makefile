default: deps lint test

lint:
	golint ./...

test:
	go test -cover -v ./...

test-slow:
	make -C ./compat libotr-compat

ci: lint test test-slow

deps:
	./deps.sh

cover:
	go test . -coverprofile=coverage.out
	go tool cover -html=coverage.out
