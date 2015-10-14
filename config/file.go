package config

import (
	"encoding/json"
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

func FindConfigFile() string {
	dir := configDir()
	ensureDir(dir, 0700)
	return filepath.Join(dir, "accounts.json")
}

func (c *Config) Save() error {
	contents, err := c.Serialize()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.Filename, contents, 0600)
}

func (c *Config) Serialize() ([]byte, error) {
	c.SerializeFingerprints()
	return json.MarshalIndent(c, "", "\t")
}
