GO_VERSION=$(shell go version | grep  -o 'go[[:digit:]]\.[[:digit:]]')
AS_CREATE_DMG=/usr/local/Cellar/create-dmg/1.0.9/bin/create-dmg

default: doctor test coveralls

deps:
	make -C ../ deps
	go get -u github.com/mattn/goveralls

test:
	${RUN_WITH} make check -C ../

doctor: check-go-imports

check-go-imports:
	go get -u golang.org/x/tools/cmd/goimports
	goimports -w .
	git diff --exit-code .

coveralls:
	${RUN_WITH} make run-cover -C ../
ifeq ($(CI),true)
ifeq ($(TRAVIS),true)
	@echo "We only send coverage information from Github Actions these days, so skipping it"
#	goveralls -coverprofile=../.coverprofiles/all.coverprofile -service=travis-ci
else ifeq ($(GITHUB_ACTIONS),true)
	goveralls -coverprofile=../.coverprofiles/all.coverprofile
else
	@echo "Can't send coverage results to Coveralls since we don't seem to be in neither Travis nor Github Actions"
endif
else
	@echo "Not running in CI, so skipping sending coverage results to Coveralls"
endif

clean-release:
	$(RM) bin/*

clean-mac-bundle:
	$(RM) -r ../bin/mac-bundle

mac-bundle:
#	make -C ../ build-gui
	iconutil -c icns -o ../build/osx/mac-bundle/coyim.icns ../build/osx/mac-bundle/coyim.iconset
	cd .. && build/osx/build-bundle bin/mac-bundle/CoyIM.app

release-dmg:
	mkdir -p ../release
#	$(AS_CREATE_DMG) --sandbox-safe --volname "CoyIM Installer" --volicon ../build/osx/mac-bundle/coyim.icns --window-pos 200 120 --window-size 800 400 --icon-size 100 --hide-extension "CoyIM.app" --app-drop-link 600 185 ../release/coyim.dmg ../bin/mac-bundle/CoyIM.app
#	create-dmg --help
#	create-dmg --dmg-title="CoyIM Installer" ../bin/mac-bundle/CoyIM.app ../release
#	ls -alF ../release
#	mv ../release/*.dmg ../release/coyim.dmg
#	fileicon set ../release/coyim.dmg ../build/osx/mac-bundle/coyim.icns
#	openssl dgst -sha256 ../release/coyim.dmg >> ../release/coyim.dmg_checksum
	(cd ../release; zip coyim_installer.zip coyim.dmg)
	appdmg ../build/osx/dmg-spec.json ../release/coyim-appdmg.dmg
	(cd ../release; zip coyim_installer_appdmg.zip coyim-appdmg.dmg)
	cp ../release/coyim-appdmg.dmg ../release/coyim-appdmg2.dmg
	fileicon set ../release/coyim-appdmg2.dmg ../build/osx/mac-bundle/coyim.icns
	(cd ../release; zip coyim_installer_appdmg.zip coyim-appdmg.dmg)
	(cd ../release; zip coyim_installer_appdmg2.zip coyim-appdmg2.dmg)

../bin/coyim:
	make -C ../ build-gui

build-linux:
	make -C ../reproducible/docker create-image
	make -C ../reproducible/docker build

release-gui: ../bin/coyim
	mkdir -p ../release
	cp ../bin/coyim ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH)
	openssl dgst -sha256 ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH) >> ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH)_checksum

release-gui-linux: build-linux
	mkdir -p ../release
	cp ../bin/coyim ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH)
	cp ../bin/build_info ../release/build_info_$(shell go env GOOS)_$(shell go env GOARCH)
	openssl dgst -sha256 ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH) >> ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH)_checksum

../bin/coyim.exe:
	make -C ../ build-gui-win

release-gui-win: ../bin/coyim.exe
	mkdir -p ../release
	cp ../build/windows/toast.exe ../release/toast.exe
	cp ../bin/coyim.exe ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH).exe
	openssl dgst -sha256 ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH).exe >> ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH).exe_checksum
