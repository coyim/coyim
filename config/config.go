package config

import (
	"crypto/rand"
	"strconv"
	"time"

	"github.com/twstrike/otr3"
)

// Config contains the configuration for one account
type Config struct {
	Filename string `json:"-"`
	id       string `json:"-"`

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

// NewConfig creates a new configuration from scratch
func NewConfig() *Config {
	var torProxy []string
	torAddress := detectTor()

	if len(torAddress) != 0 {
		torProxy = []string{
			newTorProxy(torAddress),
		}
	}

	var priv otr3.PrivateKey
	priv.Generate(rand.Reader)

	return &Config{
		Proxies: torProxy,
		UseTor:  torProxy != nil,

		PrivateKey:          priv.Serialize(),
		AlwaysEncrypt:       true,
		OTRAutoStartSession: true,
		OTRAutoTearDown:     true, //See #48
	}
}

// ID returns the unique identifier for this account
func (c *Config) ID() string {
	if len(c.id) == 0 {
		c.id = strconv.FormatUint(uint64(time.Now().UnixNano()), 10)
	}

	return c.id
}
