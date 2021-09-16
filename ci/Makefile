GO_VERSION=$(shell go version | grep  -o 'go[[:digit:]]\.[[:digit:]]')
AS_CREATE_DMG=/usr/local/Cellar/create-dmg/1.0.9/bin/create-dmg

PLATFORM := _$(shell go env GOOS)_$(shell go env GOARCH)
RELEASE_DIR := ../release
MAC_BUNDLE_DEF := ../build/osx/mac-bundle
MAC_BUNDLE_OUT := ../bin/mac-bundle

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
	$(RM) -r $(MAC_BUNDLE_OUT)

$(RELEASE_DIR):
	mkdir -p $@

$(MAC_BUNDLE_OUT)/CoyIM.app:
	cd .. && build/osx/build-bundle bin/mac-bundle/CoyIM.app

make-mac-bundle: $(MAC_BUNDLE_OUT)/CoyIM.app
#	make -C ../ build-gui

make-dmg: $(RELEASE_DIR)
#	$(AS_CREATE_DMG) --sandbox-safe --volname "CoyIM Installer" --volicon ../build/osx/mac-bundle/coyim.icns --window-pos 200 120 --window-size 800 400 --icon-size 100 --hide-extension "CoyIM.app" --app-drop-link 600 185 ../release/coyim.dmg ../bin/mac-bundle/CoyIM.app
#	create-dmg --help
#	create-dmg --dmg-title="CoyIM Installer" ../bin/mac-bundle/CoyIM.app ../release
#	ls -alF ../release
#	mv ../release/*.dmg ../release/coyim.dmg
#	fileicon set ../release/coyim.dmg ../build/osx/mac-bundle/coyim.icns
#	openssl dgst -sha256 ../release/coyim.dmg >> ../release/coyim.dmg_checksum
#	(cd ../release; zip coyim_installer.zip coyim.dmg)
	appdmg ../build/osx/dmg-spec.json ../release/coyim-appdmg.dmg
	(cd ../release; zip coyim_installer_appdmg.zip coyim-appdmg.dmg)
	cp ../release/coyim-appdmg.dmg ../release/coyim-appdmg2.dmg
	fileicon set ../release/coyim-appdmg2.dmg ../build/osx/mac-bundle/coyim.icns
	(cd ../release; zip coyim_installer_appdmg.zip coyim-appdmg.dmg)
	(cd ../release; zip coyim_installer_appdmg2.zip coyim-appdmg2.dmg)

../bin/coyim:
	make -C ../ build-gui

../bin/coyim.exe:
	make -C ../ build-gui-win

build-linux:
	make -C ../reproducible/docker create-image
	make -C ../reproducible/docker build

make-gui-for-release: $(RELEASE_DIR)/coyim$(PLATFORM) $(RELEASE_DIR)/coyim$(PLATFORM)_checksum

make-gui-for-release-linux: build-linux $(RELEASE_DIR)
	cp ../bin/coyim $(RELEASE_DIR)/coyim$(PLATFORM)
	cp ../bin/build_info $(RELEASE_DIR)/build_info$(PLATFORM)
	openssl dgst -sha256 $(RELEASE_DIR)/coyim$(PLATFORM) >> $(RELEASE_DIR)/coyim$(PLATFORM)_checksum

$(RELEASE_DIR)/toast.exe: $(RELEASE_DIR)
	cp ../build/windows/toast.exe $@

$(RELEASE_DIR)/coyim$(PLATFORM): ../bin/coyim $(RELEASE_DIR)
	cp ../bin/coyim $@

$(RELEASE_DIR)/coyim$(PLATFORM).exe: ../bin/coyim.exe $(RELEASE_DIR)
	cp ../bin/coyim.exe $@

%_checksum: %
	openssl dgst -sha256 $< >> $@

make-gui-for-release-win: $(RELEASE_DIR)/coyim$(PLATFORM).exe $(RELEASE_DIR)/toast.exe $(RELEASE_DIR)/coyim$(PLATFORM).exe_checksum
