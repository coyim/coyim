package config

import (
	"encoding/hex"
	"errors"
	"log"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/otr3"
)

var (
	errCertificateSizeMismatch = errors.New("ServerCertificateSHA256 is not 32 bytes long")
)

// Account contains the configuration for one account
type Account struct {
	id string `json:"-"`

	Account                 string
	Server                  string   `json:",omitempty"`
	Proxies                 []string `json:",omitempty"`
	Password                string   `json:",omitempty"`
	Port                    int      `json:",omitempty"`
	DeprecatedPrivateKey    []byte   `json:"PrivateKey,omitempty"`
	PrivateKeys             [][]byte `json:",omitempty"`
	KnownFingerprints       []KnownFingerprint
	HideStatusUpdates       bool
	RequireTor              bool
	OTRAutoTearDown         bool
	OTRAutoAppendTag        bool
	OTRAutoStartSession     bool
	ServerCertificateSHA256 string   `json:",omitempty"`
	AlwaysEncrypt           bool     `json:",omitempty"`
	AlwaysEncryptWith       []string `json:",omitempty"`
	DontEncryptWith         []string `json:",omitempty"`
	InstanceTag             uint32   `json:",omitempty"`
	ConnectAutomatically    bool
}

// AllPrivateKeys returns all private keys for this account
func (a *Account) AllPrivateKeys() [][]byte {
	if len(a.DeprecatedPrivateKey) > 0 {
		return append(a.PrivateKeys, a.DeprecatedPrivateKey)
	}
	return a.PrivateKeys
}

// SerializedKeys will generate a new slice of a byte slice containing serializations of all keys given
func SerializedKeys(keys []otr3.PrivateKey) [][]byte {
	var result [][]byte

	for _, k := range keys {
		result = append(result, k.Serialize())
	}

	return result
}

// NewAccount creates a new account
func NewAccount() (*Account, error) {
	pkeys, err := otr3.GenerateMissingKeys([][]byte{})
	if err != nil {
		return nil, err
	}

	return &Account{
		RequireTor:          true,
		PrivateKeys:         SerializedKeys(pkeys),
		AlwaysEncrypt:       true,
		OTRAutoStartSession: true,
		OTRAutoTearDown:     true, //See #48
		OTRAutoAppendTag:    true,
	}, nil
}

// EnsureTorProxy makes sure the account has a Tor Proxy configured
func (a *Account) EnsureTorProxy(torAddress string) {
	if !a.RequireTor {
		return
	}

	if a.Proxies == nil {
		a.Proxies = make([]string, 0, 1)
	}

	for _, proxy := range a.Proxies {
		p, err := url.Parse(proxy)
		if err != nil {
			continue
		}

		//Already configured
		if p.Host == torAddress {
			return
		}
	}

	//Tor refuses to connect to any other proxy at localhost/127.0.0.1 in the
	//chain, so we remove them
	allowedProxies := make([]string, 0, len(a.Proxies))
	for _, proxy := range a.Proxies {
		p, err := url.Parse(proxy)
		if err != nil {
			continue
		}

		host, _, err := net.SplitHostPort(p.Host)
		if err != nil {
			host = p.Host
		}

		if host == "localhost" || host == "127.0.0.1" {
			continue
		}

		allowedProxies = append(allowedProxies, proxy)
	}

	torProxy := newTorProxy(torAddress)
	allowedProxies = append(allowedProxies, torProxy)
	a.Proxies = allowedProxies
}

// ServerCertificateHash returns the hash for the server certificate
func (a *Account) ServerCertificateHash() ([]byte, error) {
	var certSHA256 []byte
	var err error
	if len(a.ServerCertificateSHA256) > 0 {
		certSHA256, err = hex.DecodeString(a.ServerCertificateSHA256)
		if err != nil {
			return nil, errors.New("Failed to parse ServerCertificateSHA256 (should be hex string): " + err.Error())
		}

		if len(certSHA256) != 32 {
			return nil, errCertificateSizeMismatch
		}
	}

	return certSHA256, err
}

// Is returns true if this account represents the same identity as the given JID
func (a *Account) Is(jid string) bool {
	return a.Account == xmpp.RemoveResourceFromJid(jid)
}

// ShouldEncryptTo returns true if the connection with this peer should be encrypted
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

// ToggleAlwaysEncrypt toggles the state of AlwaysEncrypt config
func (a *Account) ToggleAlwaysEncrypt() {
	a.AlwaysEncrypt = !a.AlwaysEncrypt
}

// ToggleConnectAutomatically toggles the state of ConnectAutomatically config
func (a *Account) ToggleConnectAutomatically() {
	a.ConnectAutomatically = !a.ConnectAutomatically
}

func (a *Account) allowsOTR(version string) bool {
	return version == "2" || version == "3" // || version == "J"
}

func (a *Account) shouldSendWhitespace() bool {
	return a.OTRAutoAppendTag
}

func (a *Account) shouldStartAKEAutomatically() bool {
	return true
}

// SetOTRPoliciesFor will set the OTR policies on the given conversation based on the users settings
func (a *Account) SetOTRPoliciesFor(jid string, c *otr3.Conversation) {
	if a.allowsOTR("2") {
		c.Policies.AllowV2()
	}
	if a.allowsOTR("3") {
		c.Policies.AllowV3()
	}
	// if a.allowsOTR("J") {
	// 	c.Policies.AllowVExtensionJ()
	// }
	if a.shouldSendWhitespace() {
		c.Policies.SendWhitespaceTag()
	}
	if a.shouldStartAKEAutomatically() {
		c.Policies.WhitespaceStartAKE()
	}
	if a.ShouldEncryptTo(jid) {
		c.Policies.RequireEncryption()
		c.Policies.ErrorStartAKE()
	}
}

// ID returns the unique identifier for this account
func (a *Account) ID() string {
	if len(a.id) == 0 {
		a.id = strconv.FormatUint(uint64(time.Now().UnixNano()), 10)
	}

	return a.id
}

// EnsurePrivateKey generates a private key for the account in case it's missing
func (a *Account) EnsurePrivateKey() (hasUpdate bool, e error) {
	log.Printf("[%s] ensureConfigHasKey()\n", a.Account)

	prevKeys := a.AllPrivateKeys()
	newKeys, err := otr3.GenerateMissingKeys(prevKeys)

	if err != nil {
		return false, err
	}
	if len(newKeys) == 0 {
		return false, nil
	}

	a.DeprecatedPrivateKey = nil
	a.PrivateKeys = append(prevKeys, SerializedKeys(newKeys)...)

	return true, nil
}
