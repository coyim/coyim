GTK_VERSION=$(shell pkg-config --modversion gtk+-3.0 | tr . _ | cut -d '_' -f 1-2)
GTK_BUILD_TAG="gtk_$(GTK_VERSION)"
GIT_VERSION=$(shell git rev-parse HEAD)

default: deps gen-ui-defs lint test
.PHONY: test

build: build-cli build-gui

# This should not be added as a requirement to build-gui because it may hide
# build problems. build-gui is exactly what `go get` will do on a clean repo
gen-ui-defs:
	make -C ./gui/definitions

build-gui: generate-version-file
	go build -tags $(GTK_BUILD_TAG) -o bin/coyim

build-cli: generate-version-file
	go build -tags cli -o bin/coyim-cli

build-debug:
	go build -gcflags "-N -l" -tags $(GTK_BUILD_TAG) -o bin/coyim-debug

debug: build-debug
	gdb bin/coyim-debug -d $(shell go env GOROOT) -x build/debug

i18n:
	make -C i18n
.PHONY: i18n

lint:
	golint ./...

test:
	go test -cover -v -tags $(GTK_BUILD_TAG) ./...

clean-gui-test:
	$(RM) gui-test/*

#TODO: this should only be called on a linux environment
gui-test: clean-gui-test
ifeq ($(shell uname), Linux)
	git clone https://github.com/twstrike/coyim-testing.git gui-test
	echo $$COYIM_PATH
	cd gui-test && behave --stop
endif

generate-version-file:
	./gen_version_file.sh $(GIT_VERSION)

run-cover: clean-cover
	go test -coverprofile=xmpp.coverprofile ./xmpp
	go test -coverprofile=session.coverprofile ./session
	go test -coverprofile=event.coverprofile ./event
	go test -coverprofile=config.coverprofile ./config
	go test -coverprofile=ui.coverprofile ./ui
	go test -tags $(GTK_BUILD_TAG) -coverprofile=gui.coverprofile ./gui
	go test -coverprofile=roster.coverprofile ./roster
	go test -coverprofile=main.coverprofile
	gover .

clean-cover:
	$(RM) *.coverprofile

# generats an HTML report with coverage information
cover: run-cover
	go tool cover -html=gover.coverprofile

get:
	go get -t -tags $(GTK_BUILD_TAG) ./...

deps-u:
	go get -u github.com/golang/lint/golint
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/modocache/gover
	go get -u -tags $(GTK_BUILD_TAG) github.com/gotk3/gotk3/gtk
	go get -u github.com/twstrike/otr3
	go get -u github.com/twstrike/otr3/sexp
	go get -u golang.org/x/crypto/ssh/terminal
	go get -u golang.org/x/net/html
	go get -u golang.org/x/net/html/atom
	go get -u golang.org/x/net/proxy
	go get -u golang.org/x/text/transform
	go get -u gopkg.in/check.v1
	go get -u github.com/miekg/dns
	go get -u golang.org/x/crypto/scrypt
	go get -u github.com/hydrogen18/stalecucumber
	go get -u github.com/DHowett/go-plist

deps-dev:
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover
	go get github.com/modocache/gover

deps: deps-dev
	go get -tags $(GTK_BUILD_TAG) github.com/gotk3/gotk3/gtk
	go get github.com/twstrike/otr3
	go get github.com/twstrike/otr3/sexp
	go get golang.org/x/crypto/ssh/terminal
	go get golang.org/x/net/html
	go get golang.org/x/net/html/atom
	go get golang.org/x/net/proxy
	go get golang.org/x/text/transform
	go get gopkg.in/check.v1
	go get github.com/miekg/dns
	go get golang.org/x/crypto/scrypt
	go get github.com/hydrogen18/stalecucumber
	go get github.com/DHowett/go-plist
