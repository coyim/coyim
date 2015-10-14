package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type MultiAccountConfig struct {
	keepXmppClientCompat bool
	Accounts             []Config
}

func (multiAccountConfig *MultiAccountConfig) Add(conf Config) {
	multiAccountConfig.Accounts = append(multiAccountConfig.Accounts, conf)
}

func (multiAccountConfig *MultiAccountConfig) Serialize() ([]byte, error) {
	for _, account := range multiAccountConfig.Accounts {
		account.SerializeFingerprints()
	}

	return json.MarshalIndent(multiAccountConfig, "", "\t")
}

func ParseMultiConfig(filename string) (*MultiAccountConfig, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c, err := parseMultiConfig(contents)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func parseMultiConfig(conf []byte) (m *MultiAccountConfig, err error) {
	m = &MultiAccountConfig{}
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

func fallbackToSingleAccountConfig(conf []byte) (*MultiAccountConfig, error) {
	c, err := parseSingleConfig(conf)
	if err != nil {
		return nil, err
	}

	return &MultiAccountConfig{
		keepXmppClientCompat: true,
		Accounts:             []Config{*c},
	}, nil
}

func ParseConfig(filename string) (*MultiAccountConfig, error) {
	m, err := ParseMultiConfig(filename)
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
