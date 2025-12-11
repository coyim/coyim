GLIB_VERSION=$(shell pkg-config --modversion glib-2.0 | tr . _ | cut -d '_' -f 1-2)
GLIB_BUILD_TAG="glib_$(GLIB_VERSION)"

GTK_VERSION_FULL=$(shell pkg-config --modversion gtk+-3.0)
GTK_VERSION_PATCH=$(shell echo $(GTK_VERSION_FULL) | cut -f3 -d.)
GTK_VERSION=$(shell echo $(GTK_VERSION_FULL) | tr . _ | cut -d '_' -f 1-2)
GTK_BUILD_TAG="gtk_$(GTK_VERSION)"

# All this is necessary to downgrade the gtk version used to 3.22 if the
# 3.24 patch level is lower than 14. The reason for that is that
# a new variable was introduced at 3.24.14, and older patch levels
# won't compile with gotk3

GTK_VERSION_PATCH_LESS14=$(shell expr $(GTK_VERSION_PATCH) \< 14)
ifeq ($(GTK_BUILD_TAG),"gtk_3_24")
ifeq ($(GTK_VERSION_PATCH_LESS14),1)
GTK_BUILD_TAG="gtk_3_22"
endif
endif

PANGO_VERSION=$(shell pkg-config --modversion pango | tr . _ | cut -d '_' -f 1-2)
PANGO_BUILD_TAG="pango_$(PANGO_VERSION)"

CAIRO_VERSION=$(shell pkg-config --modversion cairo | tr . _ | cut -d '_' -f 1-2)
CAIRO_BUILD_TAG="cairo_$(CAIRO_VERSION)"

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
GCC_IGNORE_DEPRECATED_WARNINGS := CGO_CFLAGS_ALLOW="-Wno-deprecated-declarations" CGO_CFLAGS="-Wno-deprecated-declarations"
LD_IGNORE_DUPLICATED_LIBS := -Wl,-no_warn_duplicate_libraries
GO := $(PREF) $(GCC_IGNORE_DEPRECATED_WARNINGS) ${BUILD_RUN_PREFIX} go
GOBUILD := $(GO) build
GOTEST := $(GO) test
TAGS_EXTRA :=
TAGS := -tags $(GLIB_BUILD_TAG),$(GTK_BUILD_TAG),$(PANGO_BUILD_TAG),$(CAIRO_BUILD_TAG)$(TAGS_EXTRA)

AUTOGEN := gui/settings/definitions/schemas.go

LDFLAGS_VARS := -X 'main.BuildTimestamp=$(BUILD_TIMESTAMP)' -X 'main.BuildCommit=$(GIT_VERSION)' -X 'main.BuildShortCommit=$(GIT_SHORT_VERSION)' -X 'main.BuildTag=$(TAG_VERSION)'
LDFLAGS_REGULAR = -ldflags "$(LDFLAGS_VARS)"
LDFLAGS_MAC = -ldflags "$(LDFLAGS_VARS) -extldflags $(LD_IGNORE_DUPLICATED_LIBS)"
LDFLAGS_WINDOWS = -ldflags "$(LDFLAGS_VARS) -H windowsgui"

LDF = $(LDFLAGS_REGULAR)

ifeq ($(OS),Windows_NT)
LDF = $(LDFLAGS_WINDOWS)
else
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
LDF = $(LDFLAGS_MAC)
endif
endif

.PHONY: default check autogen build build-gui build-gui-memory-analyzer build-gui-address-san build-gui-win build-debug debug win-ci-deps reproducible-linux-create-image reproducible-linux-build sign-reproducible send-reproducible-signature check-reproducible-signatures clean clean-cache update-vendor gosec ineffassign i18n lint test test-named dep-supported-only deps run-cover clean-cover cover all authors

default: check
check: lint test

$(BUILD_DIR)/coyim: $(AUTOGEN) $(SRC)
	$(GOBUILD) $(LDF) $(TAGS) -o $@

$(BUILD_DIR)/coyim-ma: $(AUTOGEN) $(SRC)
	$(GOBUILD) $(LDF) -x -msan $(TAGS) -o $@

# run with: export ASAN_OPTIONS=detect_stack_use_after_return=1:check_initialization_order=1:strict_init_order=1:verbosity=1:handle_segv=0
$(BUILD_DIR)/coyim-aa: $(AUTOGEN) $(SRC)
	CC="clang" CGO_CFLAGS="-fsanitize=address -fsanitize-address-use-after-scope -g -O1 -fno-omit-frame-pointer" CGO_LDFLAGS="-fsanitize=address" $(GOBUILD) $(LDFLAGS_REGULAR) -x -ldflags '-extldflags "-fsanitize=address"' $(TAGS) -o $@

$(BUILD_DIR)/coyim.exe: $(AUTOGEN) $(SRC)
	CGO_LDFLAGS_ALLOW=".*" CGO_CFLAGS_ALLOW=".*" CGO_CXXFLAGS_ALLOW=".*" CGO_CPPFLAGS_ALLOW=".*" $(GOBUILD) $(LDF) $(TAGS) -o $@

$(BUILD_DIR)/coyim-debug: $(AUTOGEN) $(SRC)
	$(GOBUILD) $(LDF) -v -gcflags "-N -l" $(TAGS) -o $@

build: build-gui
build-gui: $(BUILD_DIR)/coyim
build-gui-memory-analyzer: $(BUILD_DIR)/coyim-ma
build-gui-address-san: $(BUILD_DIR)/coyim-aa
build-gui-win: $(BUILD_DIR)/coyim.exe
build-debug: $(BUILD_DIR)/coyim-debug

debug: $(BUILD_DIR)/coyim-debug
	GDK_DEBUG=nograbs gdb -d $(shell go env GOROOT) --args $(BUILD_DIR)/coyim-debug -debug

win-ci-deps:
	go install github.com/rosatolen/esc@v0.0.0-20170322162328-d21c3d2332cb

reproducible-linux-create-image:
	make -C ./reproducible/docker create-image

reproducible-linux-build:
	make -C ./reproducible/docker build

sign-reproducible:
	make -C ./reproducible sign-reproducible

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
	./build/find_esc.sh $(BUILD_TOOLS_DIR)

gui/muc/definitions.go: $(BUILD_TOOLS_DIR)/esc gui/muc/definitions/*.xml
	(cd gui/muc; go generate -x ui_reader.go)

gui/authors.go: build/authors.rb
	rm -rf $@
	./build/authors.rb > $@
	gofmt -w $@

touch-authors:
	touch build/authors.rb

authors: touch-authors gui/authors.go

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
	make -C ./i18n generate

lint: $(AUTOGEN)
	golint -set_exit_status $(SRC_DIRS)

test: $(AUTOGEN)
	$(GOTEST) $(LDF) -cover -v $(TAGS) ./...

test-named: $(AUTOGEN)
	$(GOTEST) -v $(TAGS) $(SRC_DIRS)

deps:
	go install golang.org/x/lint/golint@v0.0.0-20210508222113-6edffad5e616
	go install github.com/rosatolen/esc@v0.0.0-20170322162328-d21c3d2332cb

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
