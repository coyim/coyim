package config

import (
	"io/ioutil"
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
	switch {
	case fileExists(basePath + encryptedFileEnding):
		return basePath + encryptedFileEnding, true
	case fileExists(basePath + encryptedFileEnding + tmpExtension):
		return basePath + encryptedFileEnding, true
	}
	return basePath, false
}

const tmpExtension = ".000~"

func safeWrite(name string, data []byte, perm os.FileMode) error {
	tempName := name + tmpExtension
	err := ioutil.WriteFile(tempName, data, perm)
	if err != nil {
		return err
	}

	if fileExists(name) {
		os.Remove(name)
	}

	return os.Rename(tempName, name)
}

func readFileOrTemporaryBackup(name string) (data []byte, e error) {
	if fileExists(name) {
		data, e = ioutil.ReadFile(name)
		if len(data) == 0 && fileExists(name+tmpExtension) {
			data, e = ioutil.ReadFile(name + tmpExtension)
		}
		return
	}
	return ioutil.ReadFile(name + tmpExtension)
}
