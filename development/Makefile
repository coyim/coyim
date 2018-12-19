GTK_VERSION=$(shell pkg-config --modversion gtk+-3.0 | tr . _ | cut -d '_' -f 1-2)
GTK_BUILD_TAG="gtk_$(GTK_VERSION)"
GO_VERSION=$(shell go version | grep  -o 'go[[:digit:]]\.[[:digit:]]')

GOTEST=govendor test +local
GOLIST=govendor list -no-status +local

default: autogen check
check: lint test

# TODO: include i18n?
autogen: gen-ui-defs gen-schema-defs

check-deps:
	@type esc >/dev/null 2>&1 || (echo "The program 'esc' is required but not available. Please install it by running 'make deps'." && exit 1)

gen-ui-defs: check-deps
	make -C ../gui

gen-schema-defs:
	make -C ../gui/settings/definitions

i18n:
	make -C ../i18n
.PHONY: i18n

lint:
ifeq ($(GO_VERSION), go1.7)
	echo "$(GO_VERSION) is not a supported Go release. Skipping golint."
else ifeq ($(GO_VERSION), go1.8)
	echo "$(GO_VERSION) is not a supported Go release. Skipping golint."
else
	for pkg in $$($(GOLIST) ./...) ; do \
		golint $$pkg ; \
	done
endif

test:
	$(GOTEST) -cover -v -tags $(GTK_BUILD_TAG) ./...
.PHONY: test

test-named:
	$(GOTEST) -v -tags $(GTK_BUILD_TAG) ./cli ./client ./config ./config/importer ./event ./gui ./i18n ./net ./roster ./sasl ./sasl/digestmd5 ./sasl/plain ./sasl/scram ./servers ./session ./xmpp ./xmpp/data ./xmpp/utils ./ui

dep-supported-only:
ifeq ($(GO_VERSION), go1.7)
	echo "$(GO_VERSION) is not a supported Go release. Skipping golint."
else ifeq ($(GO_VERSION), go1.8)
	echo "$(GO_VERSION) is not a supported Go release. Skipping golint."
else
	go get -u github.com/golang/lint/golint
endif

deps: dep-supported-only
	go get -u github.com/kardianos/govendor
	go get -u github.com/modocache/gover
	go get -u github.com/rosatolen/esc

run-cover: clean-cover
	mkdir -p ../.coverprofiles
	go test -coverprofile=../.coverprofiles/cli.coverprofile     ./../cli
	go test -coverprofile=../.coverprofiles/config.coverprofile  ./../config
	go test -coverprofile=../.coverprofiles/config_importer.coverprofile  ./../config/importer
	go test -coverprofile=../.coverprofiles/gui.coverprofile  ./../gui
	go test -coverprofile=../.coverprofiles/i18n.coverprofile    ./../i18n
	go test -coverprofile=../.coverprofiles/net.coverprofile     ./../net
	go test -coverprofile=../.coverprofiles/otr_client.coverprofile  ./../otr_client
	go test -coverprofile=../.coverprofiles/roster.coverprofile  ./../roster
	go test -coverprofile=../.coverprofiles/sasl.coverprofile    ./../sasl
	go test -coverprofile=../.coverprofiles/sasl_digestmd5.coverprofile    ./../sasl/digestmd5
	go test -coverprofile=../.coverprofiles/sasl_plain.coverprofile        ./../sasl/plain
	go test -coverprofile=../.coverprofiles/sasl_scram.coverprofile        ./../sasl/scram
	go test -coverprofile=../.coverprofiles/servers.coverprofile ./../servers
	go test -coverprofile=../.coverprofiles/session.coverprofile ./../session
	go test -coverprofile=../.coverprofiles/xmpp.coverprofile    ./../xmpp
	go test -coverprofile=../.coverprofiles/xmpp_data.coverprofile    ./../xmpp/data
	go test -coverprofile=../.coverprofiles/ui.coverprofile      ./../ui
	go test -tags $(GTK_BUILD_TAG) -coverprofile=../.coverprofiles/main.coverprofile ./..
	gover ../.coverprofiles ../.coverprofiles/gover.coverprofile

clean-cover:
	$(RM) -rf ../.coverprofiles

# generats an HTML report with coverage information
cover: run-cover
	go tool cover -html=../.coverprofiles/gover.coverprofile

