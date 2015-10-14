default: deps lint test
.PHONY: test

build: build-cli build-gui

build-gui:
	go build -tags nocli -o bin/coyim-gui

build-cli:
	go build -o bin/coyim

lint:
	golint ./...

test:
	go test -cover -v ./...

ci: get default coveralls

coveralls:
	go get github.com/axw/gocov/gocov
	go get github.com/modocache/gover
	go get github.com/mattn/goveralls
	go test -coverprofile=xmpp.coverprofile ./xmpp
	go test -coverprofile=session.coverprofile ./session
	go test -coverprofile=event.coverprofile ./event
	go test -coverprofile=config.coverprofile ./config
	go test -coverprofile=ui.coverprofile ./ui
	go test -coverprofile=gui.coverprofile ./gui
	go test -coverprofile=main.coverprofile
	gover .
	goveralls -coverprofile=gover.coverprofile -service=travis-ci || true

get:
	go get -t ./...

deps-u:
	go get -u github.com/golang/lint/golint
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/modocache/gover
	go get -u github.com/twstrike/gotk3/gtk
	go get -u github.com/twstrike/otr3
	go get -u github.com/twstrike/otr3/sexp
	go get -u golang.org/x/crypto/ssh/terminal
	go get -u golang.org/x/net/html
	go get -u golang.org/x/net/html/atom
	go get -u golang.org/x/net/proxy
	go get -u gopkg.in/check.v1

deps:
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover
	go get github.com/modocache/gover
	go get github.com/twstrike/gotk3/gtk
	go get github.com/twstrike/otr3
	go get github.com/twstrike/otr3/sexp
	go get golang.org/x/crypto/ssh/terminal
	go get golang.org/x/net/html
	go get golang.org/x/net/html/atom
	go get golang.org/x/net/proxy
	go get gopkg.in/check.v1

cover:
	rm -f *.coverprofile
	go test -coverprofile=xmpp.coverprofile ./xmpp
	go test -coverprofile=session.coverprofile ./session
	go test -coverprofile=event.coverprofile ./event
	go test -coverprofile=config.coverprofile ./config
	go test -coverprofile=ui.coverprofile ./ui
	go test -coverprofile=gui.coverprofile ./gui
	go test -coverprofile=main.coverprofile
	gover .
	go tool cover -html=gover.coverprofile
