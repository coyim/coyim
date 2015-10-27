package config

import (
	"os"
	"path/filepath"
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func ensureDir(dirname string, perm os.FileMode) {
	if !fileExists(dirname) {
		os.MkdirAll(dirname, perm)
	}
}

func findConfigFile() (string, bool) {
	dir := configDir()
	ensureDir(dir, 0700)
	basePath := filepath.Join(dir, "accounts.json")
	if fileExists(basePath + encryptedFileEnding) {
		return basePath + encryptedFileEnding, true
	}
	return basePath, false
}
