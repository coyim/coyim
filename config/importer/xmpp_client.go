package importer

import (
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

func (x *xmppClientImporter) importFrom(f string) (*config.Accounts, bool) {
	contents, _ := ioutil.ReadFile(f)

	c := new(xmppClientConfig)
	json.Unmarshal(contents, c)

	a := new(config.Accounts)
	ac, _ := a.AddNewAccount()

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
	ac.PrivateKey = c.PrivateKey
	ac.AlwaysEncryptWith = c.AlwaysEncryptWith
	ac.KnownFingerprints = make([]config.KnownFingerprint, len(c.KnownFingerprints))
	for ix, kf := range c.KnownFingerprints {
		ac.KnownFingerprints[ix] = config.KnownFingerprint{
			UserID:         kf.UserID,
			FingerprintHex: kf.FingerprintHex,
			Untrusted:      false,
		}
	}

	ac.RequireTor = len(c.Proxies) > 0

	a.NotifyCommand = c.NotifyCommand
	a.Bell = c.Bell
	a.RawLogFile = c.RawLogFile
	a.IdleSecondsBeforeNotification = c.IdleSecondsBeforeNotification

	return a, true
}
func (x *xmppClientImporter) findFiles() []string {
	var res []string
	// look at .xmpp-client
	// look at Persistent/.xmpp-client
	// look at all files in .xmpp-client
	// look at all files in .xmpp-clients
	// persistentDir := filepath.Join(homeDir, "Persistent")
	// if stat, err := os.Lstat(persistentDir); err == nil && stat.IsDir() {
	// 	// Looks like Tails.
	// 	homeDir = persistentDir
	// }
	// *configFile = filepath.Join(homeDir, ".xmpp-client")

	return res
}

func (x *xmppClientImporter) TryImport() []*config.Accounts {
	var res []*config.Accounts
	return res
}
