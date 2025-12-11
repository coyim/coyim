package i18n

import (
	"embed"
	"io"
	"io/fs"
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

//go:embed es fr nb_NO pt sv
var files embed.FS

func (u *unpacker) hasCorrectGuard() bool {
	targetFile := filepath.Join(u.dir, guardFileName)
	u.log.WithField("file", targetFile).WithField("expected guard", u.guard).Debug("reading guard from file")
	content, e := os.ReadFile(targetFile)
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
	e := os.WriteFile(targetFile, []byte(u.guard), defaultFilePermission)
	if e != nil {
		u.log.WithError(e).WithField("file", targetFile).Error("couldn't write translation guard file")
	}
}

func (u *unpacker) isTranslationFile(name string, entry fs.DirEntry) bool {
	return !entry.IsDir() && strings.HasSuffix(name, translationFileSuffix)
}

func (u *unpacker) copyTranslationFile(name string) {
	f, e := files.Open(name)

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

	content, re := io.ReadAll(f)
	if re != nil {
		u.log.WithError(re).WithField("file", name).Error("couldn't read file content")
		return
	}

	log.WithField("where", "i18n").WithField("file", name).WithField("target", targetFile).Info("writing translation file")
	e = os.WriteFile(targetFile, content, defaultFilePermission)

	if e != nil {
		u.log.WithError(e).WithField("file", targetFile).Error("couldn't write translation file")
	}
}

func (u *unpacker) allFilesAndDirs(path string, d fs.DirEntry, err error) error {
	if u.isTranslationFile(path, d) {
		u.copyTranslationFile(path)
	}
	return nil
}

func (u *unpacker) copyAllTranslationFiles() {
	e := fs.WalkDir(files, "", u.allFilesAndDirs)
	if e != nil {
		u.log.WithError(e).Error("couldn't walk translation files")
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
