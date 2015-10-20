package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// MultiAccount keeps track of several account configurations
type MultiAccount struct {
	keepXmppClientCompat bool
	Accounts             []Config
}

// Add will add a new configuration to the multi account
func (multiAccount *MultiAccount) Add(conf Config) {
	multiAccount.Accounts = append(multiAccount.Accounts, conf)
}

// Serialize will serialize all the account information
func (multiAccount *MultiAccount) Serialize() ([]byte, error) {
	for _, account := range multiAccount.Accounts {
		account.serializeFingerprints()
	}

	return json.MarshalIndent(multiAccount, "", "\t")
}

func readMultiAccount(filename string) (*MultiAccount, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c, err := parseMultiAccount(contents)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func parseMultiAccount(conf []byte) (m *MultiAccount, err error) {
	m = &MultiAccount{}
	if err = json.Unmarshal(conf, &m); err != nil {
		return
	}

	if m.Accounts == nil {
		return fallbackToSingleAccountConfig(conf)
	}

	return
}

func parseSingleConfig(contents []byte) (c *Config, err error) {
	c = new(Config)
	if err = json.Unmarshal(contents, &c); err != nil {
		return
	}

	return
}

func fallbackToSingleAccountConfig(conf []byte) (*MultiAccount, error) {
	c, err := parseSingleConfig(conf)
	if err != nil {
		return nil, err
	}

	return &MultiAccount{
		keepXmppClientCompat: true,
		Accounts:             []Config{*c},
	}, nil
}

// ParseConfig will parse the config in the named file
func ParseConfig(filename string) (*MultiAccount, error) {
	m, err := readMultiAccount(filename)
	if err != nil {
		return nil, err
	}

	if len(m.Accounts) == 0 {
		return nil, errors.New("account config is missing")
	}

	for _, c := range m.Accounts {
		parseFingerprints(&c)
	}

	return m, nil
}
