package config

import (
	"errors"
	"io/ioutil"
)

var (
	errInvalidConfigFile = errors.New("Failed to parse config file")
)

// FileManager contains the information about several accounts
type FileManager struct {
	Filename string
	*MultiAccount
}

// NewFileManager returns a new file manager from the given configuration file
func NewFileManager(configFile string) *FileManager {
	if len(configFile) == 0 {
		configFile = FindConfigFile()
	}

	return &FileManager{
		Filename: configFile,
	}
}

// ParseConfigFile will parse the config file
func (fileManager *FileManager) ParseConfigFile() error {
	var err error

	fileManager.MultiAccount, err = ParseConfig(fileManager.Filename)
	if err != nil {
		return errInvalidConfigFile
	}

	return nil
}

// Add will add the given configuration to the file manager
func (fileManager *FileManager) Add(conf Config) error {
	if fileManager.keepXmppClientCompat {
		return errors.New("Cant add accounts while in compat mode")
	}

	fileManager.MultiAccount.Add(conf)

	return nil
}

// Save will save the file manager
func (fileManager *FileManager) Save() error {
	if fileManager.keepXmppClientCompat {
		account := fileManager.MultiAccount.Accounts[0]
		account.Filename = fileManager.Filename
		return account.Save()
	}

	contents, err := fileManager.Serialize()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileManager.Filename, contents, 0600)
}
