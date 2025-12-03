package config

import (
	"crypto/rand"
	"os"
)

func writeRandomData(filename string, size int64) error {
	f, err := os.OpenFile(filename, os.O_WRONLY, 0)
	if err != nil {
		return err
	}

	random := make([]byte, size)
	if _, err := rand.Read(random); err != nil {
		f.Close()
		return err
	}

	if _, err := f.Write(random); err != nil {
		f.Close()
		return err
	}

	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	f.Close()

	return nil
}

func secureRemove(filename string) error {
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}

	// Overwrite with random data 3 times
	for i := 0; i < 3; i++ {
		err = writeRandomData(filename, info.Size())
		if err != nil {
			return err
		}
	}

	return os.Remove(filename)
}
