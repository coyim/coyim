GTK_VERSION=$(shell pkg-config --modversion gtk+-3.0 | tr . _ | cut -d '_' -f 1-2)
GTK_BUILD_TAG="gtk_$(GTK_VERSION)"
GIT_VERSION=$(shell git rev-parse HEAD)
TAG_VERSION=$(shell git tag -l --contains $$GIT_VERSION | tail -1)

BUILD_DIR=bin

# This should be the default if we want to build the app when `go get`
#default: build

default: check

check:
	make -C ./development

build: build-cli build-gui

build-gui: generate-version-file
	go build -i -tags $(GTK_BUILD_TAG) -o $(BUILD_DIR)/coyim

build-gui-memory-analyzer: generate-version-file
	go build -x -msan -i -tags $(GTK_BUILD_TAG) -o $(BUILD_DIR)/coyim-ma

# run with: export ASAN_OPTIONS=detect_stack_use_after_return=1:check_initialization_order=1:strict_init_order=1:verbosity=1:handle_segv=0
build-gui-address-san: generate-version-file
	CC="clang" CGO_CFLAGS="-fsanitize=address -fsanitize-address-use-after-scope -g -O1 -fno-omit-frame-pointer" CGO_LDFLAGS="-fsanitize=address" go build -x -i -ldflags '-extldflags "-fsanitize=address"' -tags $(GTK_BUILD_TAG) -o $(BUILD_DIR)/coyim-aa

build-gui-win: generate-version-file
	go build -i -tags $(GTK_BUILD_TAG) -ldflags -H=windowsgui -o $(BUILD_DIR)/coyim.exe

buil$(BUILD_DIR)/coyim-debugd-cli: generate-version-file
	go build -i -tags cli -o $(BUILD_DIR)/coyim-cli

build-debug:
	go build -gcflags "-N -l" -tags $(GTK_BUILD_TAG) -o $(BUILD_DIR)/coyim-debug

debug: build-debug
	GDK_DEBUG=nograbs gdb -d $(shell go env GOROOT) --args $(BUILD_DIR)/coyim-debug -debug

# TODO: We can replace this by `go build -ldflags "-X main.Version=$(TAG_VERSION)"`.
generate-version-file:
	./gen_version_file.sh $(GIT_VERSION) $(TAG_VERSION)

lint:
	make -C ./development lint

# Convenience
test:
	make -C ./development test
.PHONY: test

deps:
	make -C ./development deps

reproducible-linux-create-image:
	make -C ./reproducible/docker create-image

reproducible-linux-build: reproducible-linux-create-image
	make -C ./reproducible/docker build

sign-reproducible:
	make -C ./reproducible sign-reproducible

upload-reproducible-signature:
	make -C ./reproducible upload-reproducible-signature

send-reproducible-signature:
	make -C ./reproducible send-reproducible-signature

check-reproducible-signatures:
	make -C ./reproducible check-reproducible-signatures

# TODO: we can use `go generate` for this.
gen-authors:
	rm -rf gui/authors.go
	./authors.rb > gui/authors.go
	gofmt -w gui/authors.go

update-vendor:
	go get -u ./...
	go get -u -t ./...
	govendor update +v
