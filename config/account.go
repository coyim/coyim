package config

import (
	"crypto/rand"
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
	DontEncryptWith         []string `json:",omitempty"`
	InstanceTag             uint32   `json:",omitempty"`
	ConnectAutomatically    bool
}

// NewAccount creates a new account
func NewAccount() (*Account, error) {
	var priv otr3.PrivateKey

	err := priv.Generate(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Account{
		RequireTor:          true,
		PrivateKey:          priv.Serialize(),
		AlwaysEncrypt:       true,
		OTRAutoStartSession: true,
		OTRAutoTearDown:     true, //See #48
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
			continue
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
			return nil, errors.New("ServerCertificateSHA256 is not 32 bytes long")
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

func (a *Account) allowsOTR(version int) bool {
	return version == 2 || version == 3
}

func (a *Account) shouldSendWhitespace() bool {
	return true
}

func (a *Account) shouldStartAKEAutomatically() bool {
	return true
}

// SetOTRPoliciesFor will set the OTR policies on the given conversation based on the users settings
func (a *Account) SetOTRPoliciesFor(jid string, c *otr3.Conversation) {
	if a.allowsOTR(2) {
		c.Policies.AllowV2()
	}
	if a.allowsOTR(3) {
		c.Policies.AllowV3()
	}
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

	if len(a.PrivateKey) != 0 {
		return false, nil
	}

	log.Printf("[%s] - No private key available. Generating...\n", a.Account)
	var priv otr3.PrivateKey

	if err := priv.Generate(rand.Reader); err != nil {
		return false, err
	}

	a.PrivateKey = priv.Serialize()

	return true, nil
}
