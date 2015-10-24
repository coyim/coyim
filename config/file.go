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

func findConfigFile() string {
	dir := configDir()
	ensureDir(dir, 0700)
	return filepath.Join(dir, "accounts.json")
}
