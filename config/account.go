package config

import (
	"sort"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/otr3"
)

// Account contains the configuration for one account
type Account struct {
	id string `json:"-"`

	//TODO: this should be JID
	Account              string
	Nickname             string   `json:",omitempty"`
	Server               string   `json:",omitempty"`
	Proxies              []string `json:",omitempty"`
	Password             string   `json:",omitempty"`
	Port                 int      `json:",omitempty"`
	PrivateKeys          [][]byte `json:",omitempty"`
	Peers                []*Peer
	HideStatusUpdates    bool
	OTRAutoTearDown      bool
	OTRAutoAppendTag     bool
	OTRAutoStartSession  bool
	AlwaysEncrypt        bool   `json:",omitempty"`
	InstanceTag          uint32 `json:",omitempty"`
	ConnectAutomatically bool
	Certificates         []*CertificatePin `json:",omitempty"`
	PinningPolicy        string            `json:",omitempty"`

	LegacyKnownFingerprints       []KnownFingerprint `json:"KnownFingerprints,omitempty"`
	DeprecatedPrivateKey          []byte             `json:"PrivateKey,omitempty"`
	LegacyServerCertificateSHA256 string             `json:"ServerCertificateSHA256,omitempty"`

	// AlwaysEncryptWith and DontEncryptWith should be promoted to legacy and replaced with the peer settings
	AlwaysEncryptWith []string `json:",omitempty"`
	DontEncryptWith   []string `json:",omitempty"`
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
		PrivateKeys:         SerializedKeys(pkeys),
		AlwaysEncrypt:       true,
		OTRAutoStartSession: true,
		OTRAutoTearDown:     true, //See #48
		OTRAutoAppendTag:    true,
		Proxies: []string{
			"tor-auto://",
		},
	}, nil
}

// Is returns true if this account represents the same identity as the given JID
func (a *Account) Is(j string) bool {
	return a.Account == jid.Parse(j).NoResource().String()
}

// ShouldEncryptTo returns true if the connection with this peer should be encrypted
func (a *Account) ShouldEncryptTo(j string) bool {
	p, ok := a.GetPeer(j)

	if ok && p.EncryptionSettings != Default && p.EncryptionSettings != "" {
		return p.EncryptionSettings == AlwaysEncrypt
	}

	if a.AlwaysEncrypt {
		return true
	}

	bareJid := jid.Parse(j).NoResource().String()
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

// SaveCert will put the given certificate as a pinned certificate. It expects a SHA3-256 hash of the certificate.
func (a *Account) SaveCert(subject, issuer string, sha3Digest []byte) {
	a.Certificates = append(a.Certificates, &CertificatePin{
		Subject:         subject,
		Issuer:          issuer,
		FingerprintType: "SHA3-256",
		Fingerprint:     sha3Digest,
	})
	sort.Sort(CertificatePinsByNaturalOrder(a.Certificates))

}
