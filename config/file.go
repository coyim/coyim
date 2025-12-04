package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/coyim/coyim/internal/util"
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func ensureDir(dirname string, perm os.FileMode) {
	if !fileExists(dirname) {
		util.LogIgnoredError(os.MkdirAll(dirname, perm), nil, "ensuring directory")
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
const backupExtension = ".backup" + tmpExtension

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
	defer func() {
		util.LogIgnoredError(lock.release(), nil, "releasing lock")
	}()

	backupName := fmt.Sprintf("%s%s", name, backupExtension)
	tempName := fmt.Sprintf("%s%s", name, tmpExtension)

	// Write to temporary file first, in case we are interrupted
	err = os.WriteFile(tempName, data, perm)
	if err != nil {
		return fmt.Errorf("failed to write configuration to file %s: %w", tempName, err)
	}

	// Verify the write by reading back the size
	stat, err := os.Stat(tempName)
	if err != nil || stat.Size() != int64(len(data)) {
		util.LogIgnoredError(secureRemove(tempName), nil, "securely removing temp file")
		return fmt.Errorf("failed to verify written data: expected %d bytes, got %d", len(data), stat.Size())
	}

	if fileExists(name) {
		if fileExists(backupName) {
			util.LogIgnoredError(secureRemove(backupName), nil, "securely removing backup file")
		}

		err := osRename(name, backupName)
		if err != nil {
			util.LogIgnoredError(secureRemove(tempName), nil, "securely removing temp file")
			return fmt.Errorf("failed to rename %s to %s: %w", name, backupName, err)
		}
	}

	err = osRename(tempName, name)
	if err != nil {
		if fileExists(backupName) {
			util.LogIgnoredError(osRename(backupName, name), nil, "renaming backup to real name")
		}
		return fmt.Errorf("failed to rename %s to %s: %w", tempName, name, err)
	}

	return nil
}

func readFileOrTemporaryBackup(name string) (data []byte, e error) {
	// Fast path
	if fileExists(name) {
		data, e = os.ReadFile(filepath.Clean(name))
		if e == nil && len(data) > 0 {
			return data, nil
		}
	}

	lockPath := name + lockExtension
	if fileExists(lockPath) {
		if !isLockStale(lockPath) {
			time.Sleep(200 * time.Millisecond)
			if fileExists(name) {
				data, e = os.ReadFile(filepath.Clean(name))
				if e == nil && len(data) > 0 {
					return data, nil
				}
			}
			return nil, errors.New("config file is being written, please try again")
		}
		// Lock is stale - proceed to recovery
		util.LogIgnoredError(os.Remove(lockPath), nil, "removing stale lock file")
	}

	tempName := name + tmpExtension
	backupName := name + backupExtension

	if fileExists(tempName) {
		data, e = os.ReadFile(filepath.Clean(tempName))
		if e == nil && len(data) > 0 {
			return data, nil
		}
	}

	if fileExists(backupName) {
		data, e = os.ReadFile(filepath.Clean(backupName))
		if e == nil && len(data) > 0 {
			return data, nil
		}
	}

	return nil, fmt.Errorf("no valid config file found at %s", name)
}

func configDir() string {
	return filepath.Join(SystemConfigDir(), "coyim")
}
