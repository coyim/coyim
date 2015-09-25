package config

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	errHomeDirNotSet = errors.New("$HOME not set. Please either export $HOME or use the -config-file option.\n")
)

type MultiAccountConfig struct {
	Filename string `json:"-"`
	Accounts []Config
}

type Config struct {
	Filename                      string `json:"-"`
	Account                       string
	Server                        string   `json:",omitempty"`
	Proxies                       []string `json:",omitempty"`
	Password                      string   `json:",omitempty"`
	Port                          int      `json:",omitempty"`
	PrivateKey                    []byte
	KnownFingerprints             []KnownFingerprint
	RawLogFile                    string   `json:",omitempty"`
	NotifyCommand                 []string `json:",omitempty"`
	IdleSecondsBeforeNotification int      `json:",omitempty"`
	Bell                          bool
	HideStatusUpdates             bool
	UseTor                        bool
	OTRAutoTearDown               bool
	OTRAutoAppendTag              bool
	OTRAutoStartSession           bool
	ServerCertificateSHA256       string   `json:",omitempty"`
	AlwaysEncrypt                 bool     `json:",omitempty"`
	AlwaysEncryptWith             []string `json:",omitempty"`
}

type KnownFingerprint struct {
	UserId         string
	FingerprintHex string
	Fingerprint    []byte `json:"-"`
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

	c.Filename = filename

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

func fallbackToSingleAccountConfig(conf []byte) (*MultiAccountConfig, error) {
	c, err := parseConfig(conf)
	if err != nil {
		return nil, err
	}

	//TODO: Convert from single to multi account format

	return &MultiAccountConfig{
		Accounts: []Config{*c},
	}, nil
}

func ParseConfig(filename string) (*Config, error) {
	m, err := ParseMultiConfig(filename)
	if err != nil {
		return nil, err
	}

	if len(m.Accounts) == 0 {
		return nil, errors.New("account config is missing")
	}

	c := &m.Accounts[0]
	c.Filename = filename
	return c, parseFingerprints(c)
}

func parseConfig(contents []byte) (c *Config, err error) {
	c = new(Config)
	if err = json.Unmarshal(contents, &c); err != nil {
		return
	}

	return
}

func parseFingerprints(c *Config) error {
	var err error
	for i, known := range c.KnownFingerprints {
		c.KnownFingerprints[i].Fingerprint, err = hex.DecodeString(known.FingerprintHex)
		if err != nil {
			return errors.New("xmpp: failed to parse hex fingerprint for " + known.UserId + ": " + err.Error())
		}
	}

	return nil
}

func ParseYes(input string) bool {
	switch strings.ToLower(input) {
	case "y", "yes":
		return true
	}

	return false
}

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
	for i, known := range c.KnownFingerprints {
		c.KnownFingerprints[i].FingerprintHex = hex.EncodeToString(known.Fingerprint)
	}

	contents, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.Filename, contents, 0600)
}

func (c *Config) UserIdForFingerprint(fpr []byte) string {
	for _, known := range c.KnownFingerprints {
		if bytes.Equal(fpr, known.Fingerprint) {
			return known.UserId
		}
	}

	return ""
}

func (c *Config) HasFingerprint(uid string) bool {
	for _, known := range c.KnownFingerprints {
		if uid == known.UserId {
			return true
		}
	}

	return false
}

func (c *Config) ShouldEncryptTo(uid string) bool {
	if c.AlwaysEncrypt {
		return true
	}

	for _, contact := range c.AlwaysEncryptWith {
		if contact == uid {
			return true
		}
	}
	return false
}

var (
	errInvalidConfigFile = errors.New("Failed to parse config file")
)

func Load(configFile string) (*Config, error) {
	if len(configFile) == 0 {
		c, err := FindConfigFile(os.Getenv("HOME"))
		if err != nil {
			return nil, err
		}

		configFile = *c
	}

	config, err := ParseConfig(configFile)
	if err != nil {
		return nil, errInvalidConfigFile
	}

	config.Filename = configFile
	return config, nil
}
