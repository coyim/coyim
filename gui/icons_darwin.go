package gui

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func getActualRootFolder() string {
	wd, _ := os.Getwd()
	if strings.HasSuffix(wd, "/gui") {
		return filepath.Join(wd, "../")
	}
	return wd
}

func (i *icon) getPath() string {
	iconPath := filepath.Join(getActualRootFolder(), i.path)
	if fileNotFound(iconPath) {
		tmpIconPath := filepath.Join(filepath.Join(os.TempDir(), "coyim"), i.name)
		if fileNotFound(tmpIconPath) {
			os.MkdirAll(filepath.Join(os.TempDir(), "coyim"), 0750)
			bytes, _ := i.get()
			ioutil.WriteFile(tmpIconPath, bytes, 0600)
			log.WithFields(log.Fields{
				"name": i.name,
				"path": tmpIconPath,
			}).Debug("gui/icons: wrote file")
		}
		return tmpIconPath
	}
	return iconPath
}
