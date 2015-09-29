package config

import (
	"io/ioutil"
	"os"
)

type ConfigFileManager struct {
	Filename string
	*MultiAccountConfig
}

func NewConfigFileManager(configFile string) (*ConfigFileManager, error) {
	if len(configFile) == 0 {
		c, err := FindConfigFile(os.Getenv("HOME"))
		if err != nil {
			return nil, err
		}

		configFile = *c
	}

	return &ConfigFileManager{
		Filename: configFile,
	}, nil
}

func (configFileManager *ConfigFileManager) ParseConfigFile() error {
	var err error

	configFileManager.MultiAccountConfig, err = ParseConfig(configFileManager.Filename)
	if err != nil {
		return errInvalidConfigFile
	}

	return nil
}

func (configFileManager *ConfigFileManager) Save() error {
	if configFileManager.keepXmppClientCompat {
		account := configFileManager.MultiAccountConfig.Accounts[0]
		account.Filename = configFileManager.Filename
		return account.Save()
	}

	contents, err := configFileManager.Serialize()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFileManager.Filename, contents, 0600)
}
