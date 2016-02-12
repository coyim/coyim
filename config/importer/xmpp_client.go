package importer

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"

	"github.com/twstrike/coyim/config"
)

type xmppClientConfig struct {
	Account                       string
	Server                        string   `json:",omitempty"`
	Proxies                       []string `json:",omitempty"`
	Password                      string   `json:",omitempty"`
	Port                          int      `json:",omitempty"`
	PrivateKey                    []byte
	KnownFingerprints             []xmppClientKnownFingerprint
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

type xmppClientKnownFingerprint struct {
	UserID         string `json:"UserId"`
	FingerprintHex string
}

type xmppClientImporter struct{}

func (x *xmppClientImporter) importFrom(f string) (*config.ApplicationConfig, bool) {
	contents, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, false
	}

	c := new(xmppClientConfig)
	err = json.Unmarshal(contents, c)
	if err != nil {
		return nil, false
	}

	a := new(config.ApplicationConfig)
	ac, err := a.AddNewAccount()
	if err != nil {
		return nil, false
	}

	ac.Account = c.Account
	ac.Server = c.Server
	ac.Proxies = c.Proxies
	ac.Password = c.Password
	ac.Port = c.Port
	ac.HideStatusUpdates = c.HideStatusUpdates
	ac.OTRAutoStartSession = c.OTRAutoStartSession
	ac.OTRAutoTearDown = c.OTRAutoTearDown
	ac.OTRAutoAppendTag = c.OTRAutoAppendTag
	ac.ServerCertificateSHA256 = c.ServerCertificateSHA256
	ac.PrivateKeys = [][]byte{c.PrivateKey}
	ac.AlwaysEncryptWith = c.AlwaysEncryptWith
	ac.Peers = nil
	for _, kfpr := range c.KnownFingerprints {
		fp, _ := hex.DecodeString(kfpr.FingerprintHex)
		fpr := ac.EnsurePeer(kfpr.UserID).EnsureHasFingerprint(fp)
		fpr.Trusted = true
	}

	c.Proxies = append(c.Proxies, "tor-auto://")

	a.NotifyCommand = c.NotifyCommand
	a.Bell = c.Bell
	a.RawLogFile = c.RawLogFile
	a.IdleSecondsBeforeNotification = c.IdleSecondsBeforeNotification

	return a, true
}

func (x *xmppClientImporter) findFiles() []string {
	var res []string

	res = ifExists(res, config.WithHome(".xmpp-client"))
	res = ifExists(res, config.WithHome("Persistent/.xmpp-client"))
	res = ifExistsDir(res, config.WithHome(".xmpp-client"))
	res = ifExistsDir(res, config.WithHome(".xmpp-clients"))

	return res
}

func (x *xmppClientImporter) TryImport() []*config.ApplicationConfig {
	var res []*config.ApplicationConfig

	for _, f := range x.findFiles() {
		ac, ok := x.importFrom(f)
		if ok {
			res = append(res, ac)
		}
	}

	return res
}
