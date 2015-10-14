package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func FindConfigFile(homeDir string) (*string, error) {
	if len(homeDir) == 0 {
		return nil, errHomeDirNotSet
	}

	persistentDir := filepath.Join(homeDir, "Persistent")
	if stat, err := os.Lstat(persistentDir); err == nil && stat.IsDir() {
		// Looks like Tails.
		homeDir = persistentDir
	}

	configFile := filepath.Join(homeDir, ".xmpp-client")
	return &configFile, nil
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
