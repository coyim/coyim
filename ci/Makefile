GO_VERSION=$(shell go version | grep  -o 'go[[:digit:]]\.[[:digit:]]')

default: doctor test coveralls


deps:
	make -C ../ deps
	go get -u github.com/mattn/goveralls

test:
	make check -C ../

doctor: check-go-imports

check-go-imports:
	go get -u golang.org/x/tools/cmd/goimports
	goimports -w .
	git diff --exit-code .

coveralls:
	make run-cover -C ../
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
	make -C ../ build-gui
	iconutil -c icns -o ../build/osx/mac-bundle/coyim.icns ../build/osx/mac-bundle/coyim.iconset
	cd .. && build/osx/build-bundle bin/mac-bundle/CoyIM.app

release-dmg: clean-mac-bundle mac-bundle
	mkdir -p ../release
	create-dmg --volname "CoyIM Installer" --volicon ../build/osx/mac-bundle/coyim.icns --window-pos 200 120 --window-size 800 400 --icon-size 100 --custom-icon "CoyIM.app" ../build/osx/mac-bundle/coyim.icns 200 190 --hide-extension "CoyIM.app" --app-drop-link 600 185 ../release/coyim.dmg ../bin/mac-bundle/CoyIM.app
	fileicon set ../release/coyim.dmg ../build/osx/mac-bundle/coyim.icns
	openssl dgst -sha256 ../release/coyim.dmg >> ../release/checksum

../bin/coyim:
	make -C ../ build-gui

build-linux:
	make -C ../reproducible/docker create-image
	make -C ../reproducible/docker build

release-gui: ../bin/coyim
	mkdir -p ../release
	cp ../bin/coyim ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH)
	openssl dgst -sha256 ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH) >> ../release/checksum

release-gui-linux: build-linux
	mkdir -p ../release
	cp ../bin/coyim ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH)
	cp ../bin/build_info ../release/build_info_$(shell go env GOOS)_$(shell go env GOARCH)
	openssl dgst -sha256 ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH) >> ../release/checksum

../bin/coyim.exe:
	make -C ../ build-gui-win

release-gui-win: ../bin/coyim.exe
	mkdir -p ../release
	ls -alF ../
	ls -alF ../build
	cp ../build/windows/toast.exe ../release/toast.exe
	cp ../bin/coyim.exe ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH).exe
	openssl dgst -sha256 ../release/coyim_$(shell go env GOOS)_$(shell go env GOARCH).exe >> ../release/checksum
