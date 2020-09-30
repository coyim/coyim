GTK_VERSION=$(shell pkg-config --modversion gtk+-3.0 | tr . _ | cut -d '_' -f 1-2)
GTK_BUILD_TAG="gtk_$(GTK_VERSION)"

GIT_VERSION := $(shell git rev-parse HEAD)
GIT_SHORT_VERSION := $(shell git rev-parse --short HEAD)
TAG_VERSION := $(shell git tag -l --contains $$GIT_VERSION | tail -1)
CURRENT_DATE := $(shell date "+%Y-%m-%d")
BUILD_TIMESTAMP := $(shell TZ='GMT' date '+%Y-%m-%d %H:%M:%S')

GO_VERSION := $(shell go version | grep  -o 'go[[:digit:]]\.[[:digit:]]*')

BUILD_DIR := bin
BUILD_TOOLS_DIR := .build-tools
COVERPROFILES := .coverprofiles

PKGS := $(shell go list ./... | grep -v /vendor)
SRC_DIRS := . $(addprefix .,$(subst github.com/coyim/coyim,,$(PKGS)))
SRC_TEST := $(foreach sdir,$(SRC_DIRS),$(wildcard $(sdir)/*_test.go))
SRC_ALL := $(foreach sdir,$(SRC_DIRS),$(wildcard $(sdir)/*.go))
SRC := $(filter-out $(SRC_TEST), $(SRC_ALL))

PREF := PKG_CONFIG_PATH=/usr/local/opt/libffi/lib/pkgconfig:$$PKG_CONFIG_PATH
GO := $(PREF) go
GOBUILD := $(GO) build
GOTEST := $(GO) test
TAGS := -tags $(GTK_BUILD_TAG)

AUTOGEN := gui/settings/definitions/schemas.go gui/definitions.go

LDFLAGS := -ldflags "-X 'main.BuildTimestamp=$(BUILD_TIMESTAMP)' -X 'main.BuildCommit=$(GIT_VERSION)' -X 'main.BuildShortCommit=$(GIT_SHORT_VERSION)' -X 'main.BuildTag=$(TAG_VERSION)'"

.PHONY: default check autogen build build-gui build-gui-memory-analyzer build-gui-address-san build-gui-win build-debug debug win-ci-deps reproducible-linux-create-image reproducible-linux-build sign-reproducible upload-reproducible-signature send-reproducible-signature check-reproducible-signatures clean clean-cache update-vendor gosec ineffassign i18n lint test test-named dep-supported-only deps run-cover clean-cover cover all

default: check
check: lint test

$(BUILD_DIR)/coyim: $(AUTOGEN) $(SRC)
	$(GOBUILD) $(LDFLAGS) $(TAGS) -o $@

$(BUILD_DIR)/coyim-ma: $(AUTOGEN) $(SRC)
	$(GOBUILD) $(LDFLAGS) -x -msan $(TAGS) -o $@

# run with: export ASAN_OPTIONS=detect_stack_use_after_return=1:check_initialization_order=1:strict_init_order=1:verbosity=1:handle_segv=0
$(BUILD_DIR)/coyim-aa: $(AUTOGEN) $(SRC)
	CC="clang" CGO_CFLAGS="-fsanitize=address -fsanitize-address-use-after-scope -g -O1 -fno-omit-frame-pointer" CGO_LDFLAGS="-fsanitize=address" $(GOBUILD) $(LDFLAGS) -x -ldflags '-extldflags "-fsanitize=address"' $(TAGS) -o $@

$(BUILD_DIR)/coyim.exe: $(AUTOGEN) $(SRC)
	CGO_LDFLAGS_ALLOW=".*" CGO_CFLAGS_ALLOW=".*" CGO_CXXFLAGS_ALLOW=".*" CGO_CPPFLAGS_ALLOW=".*" $(GOBUILD) $(LDFLAGS) $(TAGS) -ldflags "-H windowsgui" -o $@

$(BUILD_DIR)/coyim-debug: $(AUTOGEN) $(SRC)
	$(GOBUILD) $(LDFLAGS) -v -gcflags "-N -l" $(TAGS) -o $@

build: build-gui
build-gui: $(BUILD_DIR)/coyim
build-gui-memory-analyzer: $(BUILD_DIR)/coyim-ma
build-gui-address-san: $(BUILD_DIR)/coyim-aa
build-gui-win: $(BUILD_DIR)/coyim.exe
build-debug: $(BUILD_DIR)/coyim-debug

debug: $(BUILD_DIR)/coyim-debug
	GDK_DEBUG=nograbs gdb -d $(shell go env GOROOT) --args $(BUILD_DIR)/coyim-debug -debug

win-ci-deps:
	go get -u github.com/rosatolen/esc

reproducible-linux-create-image:
	make -C ./reproducible/docker create-image

reproducible-linux-build:
	make -C ./reproducible/docker build

sign-reproducible:
	make -C ./reproducible sign-reproducible

upload-reproducible-signature:
	make -C ./reproducible upload-reproducible-signature

send-reproducible-signature:
	make -C ./reproducible send-reproducible-signature

check-reproducible-signatures:
	make -C ./reproducible check-reproducible-signatures

clean:
	go clean -i -x
	$(RM) -rf $(BUILD_DIR)
	$(RM) -rf $(BUILD_TOOLS_DIR)

clean-cache:
	go clean -i -cache -x
	$(RM) -rf $(BUILD_DIR)
	$(RM) -rf $(BUILD_TOOLS_DIR)

$(BUILD_TOOLS_DIR):
	mkdir -p $@

$(BUILD_TOOLS_DIR)/esc: $(BUILD_TOOLS_DIR)
	./find_esc.sh $(BUILD_TOOLS_DIR)

gui/definitions.go: $(BUILD_TOOLS_DIR)/esc gui/definitions/*.xml
	(cd gui; go generate -x ui_reader.go)

gui/muc/definitions.go: $(BUILD_TOOLS_DIR)/esc gui/muc/definitions/*.xml
	(cd gui/muc; go generate -x ui_reader.go)

gui/authors.go: authors.rb
	rm -rf gui/authors.go
	./authors.rb > gui/authors.go
	gofmt -w gui/authors.go

golangci:
	golangci-lint run

gosec:
	gosec -conf .gosec.config.json ./...

ineffassign:
	ineffassign .

gui/settings/definitions/gschemas.compiled: gui/settings/definitions/*.xml
	(cd gui/settings/definitions; glib-compile-schemas .)

gui/settings/definitions/schemas.go: gui/settings/definitions/gschemas.compiled
	(cd gui/settings/definitions; ruby ./generate.rb)

i18n:
	make -C ./i18n

lint: $(AUTOGEN)
	golint -set_exit_status $(SRC_DIRS)

test: $(AUTOGEN)
	$(GOTEST) -cover -v $(TAGS) ./...

test-named: $(AUTOGEN)
	$(GOTEST) -v $(TAGS) $(SRC_DIRS)

deps:
	go get -u golang.org/x/lint/golint
	go get -u github.com/rosatolen/esc

$(COVERPROFILES):
	mkdir -p $@

$(COVERPROFILES)/all.coverprofile: $(COVERPROFILES) $(SRC_ALL) $(AUTOGEN)
	$(GOTEST) $(TAGS) -coverprofile=$@ $(SRC_DIRS)

run-cover: $(COVERPROFILES)/all.coverprofile

clean-cover:
	$(RM) -rf $(COVERPROFILES)

# generats an HTML report with coverage information
cover: run-cover
	go tool cover -html=$(COVERPROFILES)/all.coverprofile

cover-tails: run-cover
	go tool cover -html=$(COVERPROFILES)/all.coverprofile -o ~/Tor\ Browser/coverage.html
	xdg-open ~/Tor\ Browser/coverage.html
