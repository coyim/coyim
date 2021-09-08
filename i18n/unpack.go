//go:generate ../.build-tools/esc -o translations.go -modtime 1489449600 -pkg i18n -private ar en_US es_EC pt_BR sv_SE zh_CN

package i18n

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/coyim/coyim/coylog"
	log "github.com/sirupsen/logrus"
)

const translationFileSuffix = ".mo"
const defaultDirectoryPermission = 0755
const defaultFilePermission = 0755

const guardFileName = "coyim.guard"

func (u *unpacker) hasCorrectGuard() bool {
	targetFile := filepath.Join(u.dir, guardFileName)
	u.log.WithField("file", targetFile).WithField("expected guard", u.guard).Debug("reading guard from file")
	content, e := ioutil.ReadFile(targetFile)
	if e != nil {
		u.log.WithError(e).WithField("file", targetFile).Debug("couldn't read translation guard file")
		return false
	}

	u.log.WithField("file", targetFile).WithField("found guard", string(content)).Debug("read guard from file")

	return u.guard == string(content)
}

// prepareDirectory returns false if it failed preparing the translation directory
// This signifies that the unpacking of translation files will not succeed
// If you call this method, you assume that the directory should be removed, and there
// will be no consequences for doing this
func (u *unpacker) prepareDirectory() bool {
	u.log.WithField("dir", u.dir).Debug("preparing translation directory")

	e := os.RemoveAll(u.dir)
	if e != nil {
		u.log.WithError(e).WithField("dir", u.dir).Error("failed removing directory")
		return false
	}

	e = os.MkdirAll(u.dir, defaultDirectoryPermission)
	if e != nil {
		u.log.WithError(e).WithField("dir", u.dir).Error("failed creating directory")
		return false
	}

	return true
}

func (u *unpacker) writeGuard() {
	targetFile := filepath.Join(u.dir, guardFileName)
	u.log.WithField("file", targetFile).WithField("guard", u.guard).Debug("writing guard to file")
	e := ioutil.WriteFile(targetFile, []byte(u.guard), defaultFilePermission)
	if e != nil {
		u.log.WithError(e).WithField("file", targetFile).Error("couldn't write translation guard file")
	}
}

func (u *unpacker) isTranslationFile(name string, entry *_escFile) bool {
	return !entry.isDir && strings.HasSuffix(name, translationFileSuffix)
}

func (u *unpacker) copyTranslationFile(name string) {
	f, e := _escStatic.Open(name)

	if e != nil {
		u.log.WithError(e).Error("couldn't open embedded file")
		return
	}

	// In theory we should close the opened file here, but
	// there's no point with the embedded file system

	df := filepath.Dir(name)
	targetDir := filepath.Join(u.dir, df)
	targetFile := filepath.Join(u.dir, name)
	e = os.MkdirAll(targetDir, defaultDirectoryPermission)
	if e != nil {
		u.log.WithError(e).WithField("dir", targetDir).Error("couldn't create target directory for translation")
		return
	}

	content, re := ioutil.ReadAll(f)
	if re != nil {
		u.log.WithError(re).WithField("file", name).Error("couldn't read file content")
		return
	}

	log.WithField("where", "i18n").WithField("file", name).WithField("target", targetFile).Info("writing translation file")
	e = ioutil.WriteFile(targetFile, content, defaultFilePermission)

	if e != nil {
		u.log.WithError(e).WithField("file", targetFile).Error("couldn't write translation file")
	}
}

func (u *unpacker) copyAllTranslationFiles() {
	for name, entry := range _escData {
		if u.isTranslationFile(name, entry) {
			u.copyTranslationFile(name)
		}
	}
}

type unpacker struct {
	dir   string
	guard string
	log   coylog.Logger
}

func (u *unpacker) unpackTranslationFiles() {
	if u.hasCorrectGuard() {
		u.log.Info("translations already exists for this version - no need to unpack again")
		return
	}

	u.log.Info("unpacking translations into")

	if !u.prepareDirectory() {
		u.log.Warn("failure to prepare the translations directory - we can't continue")
		return
	}

	u.copyAllTranslationFiles()

	u.writeGuard()
}

// UnpackTranslationFilesInto will unpack all translation files into the given directory
// unless the directory already contains a guard file with the given content. If the
// guard file doesn't exist, or exists but with the wrong content, the directory will be
// removed, translation files added to it, and then the guard file written.
func UnpackTranslationFilesInto(dir, guard string) {
	u := unpacker{
		dir:   dir,
		guard: guard,
		log:   log.WithField("where", "i18n").WithField("dir", dir),
	}
	u.unpackTranslationFiles()
}
