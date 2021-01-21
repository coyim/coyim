package config

import (
	"errors"
	"fmt"
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
		_ = os.MkdirAll(dirname, perm)
	}
}

func findConfigFile(filename string) string {
	if len(filename) == 0 {
		dir := configDir()
		ensureDir(dir, 0700)
		basePath := filepath.Join(dir, "accounts.json")
		switch {
		case fileExists(basePath + encryptedFileEnding):
			return basePath + encryptedFileEnding
		case fileExists(basePath + encryptedFileEnding + tmpExtension):
			return basePath + encryptedFileEnding
		}
		return basePath
	}
	ensureDir(filepath.Dir(filename), 0700)
	return filename
}

const tmpExtension = ".000~"

func safeWrite(name string, data []byte, perm os.FileMode) error {
	// This function will leave a backup of the config file every time it writes

	if len(data) < 10 {
		return errors.New("data amount too small - unlikely to be real data")
	}

	backupName := fmt.Sprintf("%s.backup.000~", name)

	if fileExists(backupName) {
		_ = os.Remove(backupName)
	}

	if fileExists(name) {
		err := os.Rename(name, backupName)
		if err != nil {
			return err
		}
	}

	tempName := fmt.Sprintf("%s%s", name, tmpExtension)
	err := ioutil.WriteFile(tempName, data, perm)
	if err != nil {
		return err
	}

	return os.Rename(tempName, name)
}

func readFileOrTemporaryBackup(name string) (data []byte, e error) {
	if fileExists(name) {
		data, e = ioutil.ReadFile(filepath.Clean(name))
		if len(data) == 0 && fileExists(name+tmpExtension) {
			data, e = ioutil.ReadFile(filepath.Clean(name + tmpExtension))
		}
		return
	}
	return ioutil.ReadFile(filepath.Clean(name + tmpExtension))
}

func configDir() string {
	return filepath.Join(SystemConfigDir(), "coyim")
}
