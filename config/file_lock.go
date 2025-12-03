package config

import (
	"fmt"
	"os"
	"time"
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
			_, _ = file.WriteString(fmt.Sprintf("%d\n%d\n", pid, time.Now().Unix()))
			_ = file.Sync()

			return &fileLock{
				path: lockPath,
				file: file,
			}, nil
		}

		if os.IsExist(err) {
			if isLockStale(lockPath) {
				_ = os.Remove(lockPath)
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
		// Can't read lock file, consider it stale
		return true
	}

	var pid int
	var timestamp int64
	_, err = fmt.Sscanf(string(data), "%d\n%d\n", &pid, &timestamp)
	if err != nil {
		// Invalid lock file format, consider it stale
		return true
	}

	lockAge := time.Since(time.Unix(timestamp, 0))
	if lockAge > staleLockAge {
		return true
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		// Process doesn't exist, consider the lock stale
		return true
	}

	_ = process

	return false
}

func (fl *fileLock) release() error {
	if fl.file != nil {
		_ = fl.file.Close()
		fl.file = nil
	}
	return os.Remove(fl.path)
}
