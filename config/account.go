package config

import (
	"strconv"
	"time"

	"github.com/twstrike/coyim/xmpp"
)

// Account contains the configuration for one account
type Account struct {
	id string `json:"-"`

	Account                 string
	Server                  string   `json:",omitempty"`
	Proxies                 []string `json:",omitempty"`
	Password                string   `json:",omitempty"`
	Port                    int      `json:",omitempty"`
	PrivateKey              []byte
	KnownFingerprints       []KnownFingerprint
	HideStatusUpdates       bool
	RequireTor              bool
	OTRAutoTearDown         bool
	OTRAutoAppendTag        bool
	OTRAutoStartSession     bool
	ServerCertificateSHA256 string   `json:",omitempty"`
	AlwaysEncrypt           bool     `json:",omitempty"`
	AlwaysEncryptWith       []string `json:",omitempty"`
	InstanceTag             uint32   `json:",omitempty"`
}

// Is returns true if this account represents the same identity as the given JID
func (a *Account) Is(jid string) bool {
	return a.Account == xmpp.RemoveResourceFromJid(jid)
}

// ShouldEncryptTo returns true if the connection to the given peer should be encrypted
func (a *Account) ShouldEncryptTo(jid string) bool {
	if a.AlwaysEncrypt {
		return true
	}

	bareJid := xmpp.RemoveResourceFromJid(jid)
	for _, contact := range a.AlwaysEncryptWith {
		if contact == bareJid {
			return true
		}
	}

	return false
}

// ID returns the unique identifier for this account
func (a *Account) ID() string {
	if len(a.id) == 0 {
		a.id = strconv.FormatUint(uint64(time.Now().UnixNano()), 10)
	}

	return a.id
}
