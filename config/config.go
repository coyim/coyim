package config

import (
	"crypto/rand"
	"strconv"
	"time"

	"github.com/twstrike/otr3"
)

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

//TODO return a config with secure defaults
func NewConfig() *Config {
	var torProxy []string = nil
	torAddress := detectTor()

	if len(torAddress) != 0 {
		torProxy = []string{
			newTorProxy(torAddress),
		}
	}

	var priv otr3.PrivateKey
	priv.Generate(rand.Reader)

	return &Config{
		//TODO: Should those 2 setting be set on startup (or maybe every connection)?
		Proxies: torProxy,
		UseTor:  torProxy != nil,

		PrivateKey:          priv.Serialize(),
		AlwaysEncrypt:       true,
		OTRAutoStartSession: true,
		OTRAutoTearDown:     true, //See #48
	}
}

func (c *Config) Id() string {
	if len(c.id) == 0 {
		c.id = strconv.FormatUint(uint64(time.Now().UnixNano()), 10)
	}

	return c.id
}
