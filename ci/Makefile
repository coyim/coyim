ci: get lint check-ui-defs test coveralls

get:
	make get -C ../

lint:
	make lint -C ../

test:
	make test -C ../

# make sure the ui definitions checked in are in sync
check-ui-defs:
	make doctor -C ../gui/definitions

# send coverage data to coveralls
coveralls:
	make run-cover -C ../
	go get github.com/mattn/goveralls
	goveralls -coverprofile=../gover.coverprofile -service=travis-ci || true

clean-release:
	$(RM) bin/*

release-cli: cross-compile

cross-compile:
	mkdir -p ../release
	go get github.com/mitchellh/gox
	gox -build-toolchain || true

	# windows does not have syscall.SIGWINCH
	cd .. && gox -os "!windows" -output "release/{{.Dir}}-cli_{{.OS}}_{{.Arch}}"

clean-mac-bundle:
	$(RM) -r ../bin/mac-bundle

mac-bundle:
	make -C ../ i18n build-gui

	iconutil -c icns -o ../build/mac-bundle/coy.icns ../build/mac-bundle/coy.iconset
	cd .. && ci/build-osx-bundle bin/mac-bundle/CoyIM.app

release-dmg: clean-mac-bundle mac-bundle
	mkdir -p ../release
	ln -s /Applications ../bin/mac-bundle/Applications
	cp ../build/mac-bundle/ds-store ../bin/mac-bundle/.DS_Store
	./make-dmg coyim ../bin/mac-bundle ../release/coyim.dmg
	openssl dgst -sha1 ../release/coyim.dmg >> ../release/checksum

../bin/coyim:
	make -C ../ build-gui

release-gui: ../bin/coyim
	mkdir -p ../release
	cp ../bin/coyim ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH)
	openssl dgst -sha1 ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH) >> ../release/checksum

../bin/coyim.exe:
	make -C ../ build-gui

release-gui-win: ../bin/coyim.exe
	mkdir -p ../release
	cp ../bin/coyim.exe ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH).exe
	openssl dgst -sha1 ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH).exe >> ../release/checksum
