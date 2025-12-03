package config

import (
	"errors"
	"fmt"
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

var osRename = os.Rename

func safeWrite(name string, data []byte, perm os.FileMode) error {
	// This function will leave a backup of the config file every time it writes
	if len(data) < 10 {
		return errors.New("data amount too small - unlikely to be real data")
	}

	lock, err := acquireFileLock(name)
	if err != nil {
		return fmt.Errorf("failed to acquire lock for %s: %w", name, err)
	}
	defer lock.release()

	backupName := fmt.Sprintf("%s.backup.000~", name)
	tempName := fmt.Sprintf("%s%s", name, tmpExtension)

	// Write to temporary file first, in case we are interrupted
	err = os.WriteFile(tempName, data, perm)
	if err != nil {
		return fmt.Errorf("failed to write configuration to file %s: %w", tempName, err)
	}

	// Verify the write by reading back the size
	stat, err := os.Stat(tempName)
	if err != nil || stat.Size() != int64(len(data)) {
		_ = os.Remove(tempName)
		return fmt.Errorf("failed to verify written data: expected %d bytes, got %d", len(data), stat.Size())
	}

	if fileExists(name) {
		if fileExists(backupName) {
			_ = os.Remove(backupName)
		}

		err := osRename(name, backupName)
		if err != nil {
			_ = os.Remove(tempName)
			return fmt.Errorf("failed to rename %s to %s: %w", name, backupName, err)
		}
	}

	err = osRename(tempName, name)
	if err != nil {
		if fileExists(backupName) {
			_ = osRename(backupName, name)
		}
		return fmt.Errorf("failed to rename %s to %s: %w", tempName, name, err)
	}

	return nil
}

func readFileOrTemporaryBackup(name string) (data []byte, e error) {
	if fileExists(name) {
		data, e = os.ReadFile(filepath.Clean(name))
		if len(data) == 0 && fileExists(name+tmpExtension) {
			data, e = os.ReadFile(filepath.Clean(name + tmpExtension))
		}
		return
	}
	return os.ReadFile(filepath.Clean(name + tmpExtension))
}

func configDir() string {
	return filepath.Join(SystemConfigDir(), "coyim")
}
