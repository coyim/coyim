package config

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/coyim/coyim/internal/util"
)

const lockExtension = ".lock"
const lockTimeout = 30 * time.Second
const staleLockAge = 60 * time.Second

type fileLock struct {
	path string
	file *os.File
}

func acquireFileLock(targetPath string) (*fileLock, error) {
	lockPath := targetPath + lockExtension
	startTime := time.Now()

	for {
		file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err == nil {
			pid := os.Getpid()
			_, err = file.WriteString(fmt.Sprintf("%d\n%d\n", pid, time.Now().Unix()))
			util.LogIgnoredError(err, nil, "writing lock data")
			util.LogIgnoredError(file.Sync(), nil, "syncing written lock data")

			return &fileLock{
				path: lockPath,
				file: file,
			}, nil
		}

		if os.IsExist(err) {
			if isLockStale(lockPath) {
				util.LogIgnoredError(os.Remove(lockPath), nil, "removing stale lock file")
				continue
			}

			if time.Since(startTime) >= lockTimeout {
				return nil, fmt.Errorf("timeout waiting for file lock on %s", targetPath)
			}

			time.Sleep(100 * time.Millisecond)
			continue
		}

		return nil, fmt.Errorf("failed to acquire lock on %s: %w", targetPath, err)
	}
}

func isLockStale(lockPath string) bool {
	data, err := os.ReadFile(lockPath)
	if err != nil {
		util.LogIgnoredError(err, nil, "reading lock file")
		// Can't read lock file, consider it stale
		return true
	}

	var pid int
	var timestamp int64
	_, err = fmt.Sscanf(string(data), "%d\n%d\n", &pid, &timestamp)
	if err != nil {
		util.LogIgnoredError(err, nil, "scanning values from lock file")
		// Invalid lock file format, consider it stale
		return true
	}

	lockAge := time.Since(time.Unix(timestamp, 0))
	if lockAge > staleLockAge {
		return true
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		util.LogIgnoredError(err, nil, "finding process from lock file")
		// Process doesn't exist, consider the lock stale
		return true
	}

	// On Unix, FindProcess always succeeds, so send signal 0 to check if process exists
	// Signal 0 doesn't actually send a signal, it just checks if we *can* send one
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		util.LogIgnoredError(err, nil, "signalling process from lock file")
		// Process doesn't exist or we don't have permission (both mean stale lock)
		return true
	}

	return false
}

func (fl *fileLock) release() error {
	if fl.file != nil {
		util.LogIgnoredError(fl.file.Close(), nil, "closing lock file")
		fl.file = nil
	}
	return os.Remove(fl.path)
}
