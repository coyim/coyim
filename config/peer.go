package config

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/twstrike/coyim/i18n"
)

// EncryptionSettings configures the encryption setting for this peer
type EncryptionSettings string

const (
	// Default encrypts conversations with this peer depending on the account
	// configuration (config.Account#AlwaysEncrypt)
	Default EncryptionSettings = "default"
	// AlwaysEncrypt always encrypts conversations with this peer
	AlwaysEncrypt = "always"
	// NeverEncrypt never encrypts conversations with this peer
	NeverEncrypt = "never"
)

// FingerprintForSerialization represents a fingerprint in its serialized form
type FingerprintForSerialization struct {
	FingerprintHex string
	Trusted        bool
}

// Fingerprint represents a known fingerprint for a specific peer
type Fingerprint struct {
	Fingerprint []byte
	Trusted     bool
}

// Peer represents one peer
type Peer struct {
	UserID             string
	Nickname           string
	EncryptionSettings EncryptionSettings `json:",omitempty"`

	Fingerprints []*Fingerprint
}

// MarshalJSON is used to create a JSON representation of this fingerprint
func (k *Fingerprint) MarshalJSON() ([]byte, error) {
	return json.Marshal(FingerprintForSerialization{
		FingerprintHex: hex.EncodeToString(k.Fingerprint),
		Trusted:        k.Trusted,
	})
}

// UnmarshalJSON is used to parse the JSON representation of a fingerprint
func (k *Fingerprint) UnmarshalJSON(data []byte) error {
	v := FingerprintForSerialization{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	k.Fingerprint, err = hex.DecodeString(v.FingerprintHex)
	if err != nil {
		return nil
	}

	k.Trusted = v.Trusted

	return nil
}

// ByNaturalOrder sorts fingerprints according to the fingerprint
type ByNaturalOrder []*Fingerprint

func (s ByNaturalOrder) Len() int { return len(s) }
func (s ByNaturalOrder) Less(i, j int) bool {
	return bytes.Compare(s[i].Fingerprint, s[j].Fingerprint) == -1
}
func (s ByNaturalOrder) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (a *Account) updateToLatestVersion() bool {
	return a.updateFingerprintsToLatestVersion() ||
		a.removeEmptyFingerprints()
}

func (a *Account) removeEmptyFingerprints() bool {
	changed := false

	for _, p := range a.Peers {
		if a.RemoveFingerprint(p.UserID, []byte{}) {
			changed = true
		}
	}

	return changed
}

func (a *Account) updateFingerprintsToLatestVersion() bool {
	if len(a.LegacyKnownFingerprints) == 0 {
		return false
	}

	for _, kfpr := range a.LegacyKnownFingerprints {
		if len(kfpr.Fingerprint) > 0 {
			fpr := a.EnsurePeer(kfpr.UserID).EnsureHasFingerprint(kfpr.Fingerprint)
			if !kfpr.Untrusted {
				fpr.Trusted = true
			}
		}
	}

	a.LegacyKnownFingerprints = []KnownFingerprint{}

	return true
}

// EnsurePeer returns the first peer with the given uid, or creates and adds a new one if none exist
func (a *Account) EnsurePeer(uid string) *Peer {
	p, ex := a.GetPeer(uid)
	if ex {
		return p
	}
	p = &Peer{UserID: uid}
	a.Peers = append(a.Peers, p)
	return p
}

// GetPeer returns the first peer with the given uid, or false if none is found
func (a *Account) GetPeer(uid string) (*Peer, bool) {
	for _, p := range a.Peers {
		if p.UserID == uid {
			return p, true
		}
	}

	return nil, false
}

// EnsureHasFingerprint ensures that the peer has the given fingerprint and returns the Fingerprint instance
func (p *Peer) EnsureHasFingerprint(fpr []byte) *Fingerprint {
	for _, f := range p.Fingerprints {
		if bytes.Equal(f.Fingerprint, fpr) {
			return f
		}
	}
	f := &Fingerprint{Fingerprint: fpr, Trusted: false}
	p.Fingerprints = append(p.Fingerprints, f)
	return f
}

// HasTrustedFingerprint returns true if the peer has the given fingerprint and it is trusted
func (p *Peer) HasTrustedFingerprint(fpr []byte) bool {
	for _, ff := range p.Fingerprints {
		if ff.Trusted && bytes.Equal(fpr, ff.Fingerprint) {
			return true
		}
	}
	return false
}

// AddTrustedFingerprint adds a new fingerprint for the given user
func (a *Account) AddTrustedFingerprint(fpr []byte, uid string) {
	f := a.EnsurePeer(uid).EnsureHasFingerprint(fpr)
	f.Trusted = true
}

// HasFingerprint returns true if we have the fingerprint for the given user
func (a *Account) HasFingerprint(uid string) bool {
	u, ok := a.GetPeer(uid)
	if ok {
		return len(u.Fingerprints) > 0
	}
	return false
}

// UserIDForVerifiedFingerprint returns the user ID for the given verified fingerprint
func (a *Account) UserIDForVerifiedFingerprint(fpr []byte) string {
	for _, pe := range a.Peers {
		if pe.HasTrustedFingerprint(fpr) {
			return pe.UserID
		}
	}

	return ""
}

var (
	errFingerprintAlreadyAuthorized = errors.New(i18n.Local("the fingerprint is already authorized"))
)

// AuthorizeFingerprint will authorize and add the fingerprint for the given user
// or return an error if the fingerprint is already associated with another user
func (a *Account) AuthorizeFingerprint(uid string, fingerprint []byte) error {
	existing := a.UserIDForVerifiedFingerprint(fingerprint)
	if len(existing) != 0 {
		return errFingerprintAlreadyAuthorized
	}

	a.AddTrustedFingerprint(fingerprint, uid)
	return nil
}

// RemoveFingerprint removes the fingerprint for the given uid
func (a *Account) RemoveFingerprint(uid string, fpr []byte) bool {
	p, ex := a.GetPeer(uid)
	if !ex {
		return false
	}

	result := false

	newFprs := make([]*Fingerprint, 0, len(p.Fingerprints))
	for _, f := range p.Fingerprints {
		if !bytes.Equal(f.Fingerprint, fpr) {
			newFprs = append(newFprs, f)
			result = true
		}
	}
	p.Fingerprints = newFprs
	return result
}

// RemovePeer removes the given peer
func (a *Account) RemovePeer(uid string) {
	newPeers := make([]*Peer, 0, len(a.Peers))
	for _, p := range a.Peers {
		if p.UserID != uid {
			newPeers = append(newPeers, p)
		}
	}
	a.Peers = newPeers
}

// RenamePeer sets the nickname for the account
func (a *Account) RenamePeer(pid, nickname string) {
	a.EnsurePeer(pid).Nickname = nickname
}
