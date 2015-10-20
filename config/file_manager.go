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

func NewFileManager(configFile string) *FileManager {
	if len(configFile) == 0 {
		configFile = FindConfigFile()
	}

	return &FileManager{
		Filename: configFile,
	}
}

func (fileManager *FileManager) ParseConfigFile() error {
	var err error

	fileManager.MultiAccount, err = ParseConfig(fileManager.Filename)
	if err != nil {
		return errInvalidConfigFile
	}

	return nil
}

func (fileManager *FileManager) Add(conf Config) error {
	if fileManager.keepXmppClientCompat {
		return errors.New("Cant add accounts while in compat mode")
	}

	fileManager.MultiAccount.Add(conf)

	return nil
}

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
