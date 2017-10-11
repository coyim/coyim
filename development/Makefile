GTK_VERSION=$(shell pkg-config --modversion gtk+-3.0 | tr . _ | cut -d '_' -f 1-2)
GTK_BUILD_TAG="gtk_$(GTK_VERSION)"

GOTEST=govendor test +local
GOLIST=govendor list -no-status +local

default: build check
check: lint test

# TODO: include i18n?
build: gen-ui-defs gen-schema-defs

check-deps:
	@type esc >/dev/null 2>&1 || (echo "The program 'esc' is required but not available. Please install it by running 'make deps-dev'." && exit 1)

gen-ui-defs: check-deps
	make -C ../gui

gen-schema-defs:
	make -C ../gui/settings/definitions

i18n:
	make -C ../i18n
.PHONY: i18n

lint:
	for pkg in $$($(GOLIST) ./...) ; do \
		golint $$pkg ; \
    done

test:
	$(GOTEST) -cover -v -tags $(GTK_BUILD_TAG) ./...
.PHONY: test
test:

test-named:
	$(GOTEST) -v -tags $(GTK_BUILD_TAG) ./cli ./client ./config ./config/importer ./event ./gui ./i18n ./net ./roster ./sasl ./sasl/digestmd5 ./sasl/plain ./sasl/scram ./servers ./session ./xmpp ./xmpp/data ./xmpp/utils ./ui

deps:
	go get -u github.com/kardianos/govendor
	go get -u github.com/golang/lint/golint
	go get -u github.com/modocache/gover
	go get -u github.com/rosatolen/esc


run-cover: clean-cover
	mkdir -p .coverprofiles
	go test -coverprofile=.coverprofiles/cli.coverprofile     ./cli
	go test -coverprofile=.coverprofiles/client.coverprofile  ./client
	go test -coverprofile=.coverprofiles/config.coverprofile  ./config
	go test -coverprofile=.coverprofiles/config_importer.coverprofile  ./config/importer
	go test -coverprofile=.coverprofiles/event.coverprofile   ./event
	go test -coverprofile=.coverprofiles/gui.coverprofile  ./gui
	go test -coverprofile=.coverprofiles/i18n.coverprofile    ./i18n
	go test -coverprofile=.coverprofiles/net.coverprofile     ./net
	go test -coverprofile=.coverprofiles/roster.coverprofile  ./roster
	go test -coverprofile=.coverprofiles/sasl.coverprofile    ./sasl
	go test -coverprofile=.coverprofiles/sasl_digestmd5.coverprofile    ./sasl/digestmd5
	go test -coverprofile=.coverprofiles/sasl_plain.coverprofile        ./sasl/plain
	go test -coverprofile=.coverprofiles/sasl_scram.coverprofile        ./sasl/scram
	go test -coverprofile=.coverprofiles/servers.coverprofile ./servers
	go test -coverprofile=.coverprofiles/session.coverprofile ./session
	go test -coverprofile=.coverprofiles/xmpp.coverprofile    ./xmpp
	go test -coverprofile=.coverprofiles/xmpp_data.coverprofile    ./xmpp/data
	go test -coverprofile=.coverprofiles/xmpp_utils.coverprofile    ./xmpp/utils
	go test -coverprofile=.coverprofiles/ui.coverprofile      ./ui
	go test -tags $(GTK_BUILD_TAG) -coverprofile=.coverprofiles/main.coverprofile
	gover .coverprofiles .coverprofiles/gover.coverprofile

clean-cover:
	$(RM) -rf .coverprofiles

# generats an HTML report with coverage information
cover: run-cover
	go tool cover -html=.coverprofiles/gover.coverprofile
