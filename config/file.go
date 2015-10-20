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

// FindConfigFile returns the config file for CoyIM
func FindConfigFile() string {
	dir := configDir()
	ensureDir(dir, 0700)
	return filepath.Join(dir, "accounts.json")
}

// Save will save the given config to the file
func (c *Config) Save() error {
	contents, err := c.Serialize()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.Filename, contents, 0600)
}

// Serialize will serialize the config
func (c *Config) Serialize() ([]byte, error) {
	c.serializeFingerprints()
	return json.MarshalIndent(c, "", "\t")
}
