GO_VERSION=$(shell go version | grep  -o 'go[[:digit:]]\.[[:digit:]]')

ci: deps lint doctor test coveralls

deps:
	make deps -C ../

lint:
	make lint -C ../

test:
	make test -C ../

doctor: check-go-imports check-ui-defs

check-go-imports:
ifeq ($(GO_VERSION), go1.3)
	echo "Your version of Go is too old for running goimports. Skipping."
else
ifeq ($(GO_VERSION), go1.4)
	echo "Your version of Go is too old for running goimports. Skipping."
else
	go get golang.org/x/tools/cmd/goimports
	goimports -w .
	git diff --exit-code .
endif
endif

# make sure the ui definitions checked in are in sync
check-ui-defs:
	make doctor -C ../gui/definitions

# send coverage data to coveralls
coveralls:
	make run-cover -C ../
	go get github.com/mattn/goveralls
	goveralls -coverprofile=../.coverprofiles/gover.coverprofile -service=travis-ci || true

clean-release:
	$(RM) bin/*

release-cli: cross-compile

cross-compile:
	mkdir -p ../release
	go get github.com/mitchellh/gox
	gox -build-toolchain || true

	# windows does not have syscall.SIGWINCH
	cd .. && gox -os "!windows" -tags "cli" -output "release/{{.Dir}}-cli_{{.OS}}_{{.Arch}}"

clean-mac-bundle:
	$(RM) -r ../bin/mac-bundle

mac-bundle:
	make -C ../ build-gui

	iconutil -c icns -o ../build/mac-bundle/coy.icns ../build/mac-bundle/coy.iconset
	cd .. && ci/build-osx-bundle bin/mac-bundle/CoyIM.app

release-dmg: clean-mac-bundle mac-bundle
	mkdir -p ../release
	ln -s /Applications ../bin/mac-bundle/Applications
	cp ../build/mac-bundle/ds-store ../bin/mac-bundle/.DS_Store
	./make-dmg coyim ../bin/mac-bundle ../release/coyim.dmg
	openssl dgst -sha256 ../release/coyim.dmg >> ../release/checksum

../bin/coyim:
	make -C ../ build-gui

build-linux:
	make -C ../ reproducible-linux-build

release-gui: ../bin/coyim
	mkdir -p ../release
	cp ../bin/coyim ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH)
	openssl dgst -sha256 ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH) >> ../release/checksum

release-gui-linux: build-linux
	mkdir -p ../release
	cp ../bin/coyim ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH)
	cp ../bin/coyim-cli ../release/coyim-cli_$(shell go env GOOS)_$(shell go env GOARCH)
	cp ../bin/build_info ../release/build_info_$(shell go env GOOS)_$(shell go env GOARCH)
	openssl dgst -sha256 ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH) >> ../release/checksum

../bin/coyim.exe:
	make -C ../ build-gui-win

release-gui-win: ../bin/coyim.exe
	mkdir -p ../release
	cp ../build/toast.exe ../release/toast.exe
	cp ../bin/coyim.exe ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH).exe
	openssl dgst -sha256 ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH).exe >> ../release/checksum
