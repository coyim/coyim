package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type MultiAccount struct {
	keepXmppClientCompat bool
	Accounts             []Config
}

func (multiAccount *MultiAccount) Add(conf Config) {
	multiAccount.Accounts = append(multiAccount.Accounts, conf)
}

func (multiAccount *MultiAccount) Serialize() ([]byte, error) {
	for _, account := range multiAccount.Accounts {
		account.SerializeFingerprints()
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
