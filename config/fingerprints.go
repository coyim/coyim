package config

import (
	"bytes"
	"encoding/hex"
	"errors"

	"github.com/twstrike/coyim/i18n"
)

// KnownFingerprint represents one fingerprint
type KnownFingerprint struct {
	UserID         string
	FingerprintHex string
	Fingerprint    []byte `json:"-"`
	Untrusted      bool
}

func parseFingerprints(a *Account) error {
	var err error
	for i, known := range a.KnownFingerprints {
		a.KnownFingerprints[i].Fingerprint, err = hex.DecodeString(known.FingerprintHex)
		if err != nil {
			return errors.New("xmpp: failed to parse hex fingerprint for " + known.UserID + ": " + err.Error())
		}
	}

	return nil
}

func (a *Account) serializeFingerprints() {
	for i, known := range a.KnownFingerprints {
		a.KnownFingerprints[i].FingerprintHex = hex.EncodeToString(known.Fingerprint)
	}
}

// UserIDForVerifiedFingerprint returns the user ID for the given verified fingerprint
func (a *Account) UserIDForVerifiedFingerprint(fpr []byte) string {
	for _, known := range a.KnownFingerprints {
		if bytes.Equal(fpr, known.Fingerprint) && !known.Untrusted {
			return known.UserID
		}
	}

	return ""
}

// AddFingerprint adds a new fingerprint for the given user
func (a *Account) AddFingerprint(fpr []byte, uid string) {
	a.KnownFingerprints = append(a.KnownFingerprints, KnownFingerprint{Fingerprint: fpr, UserID: uid, Untrusted: false})
}

// HasFingerprint returns true if we have the fingerprint for the given user
func (a *Account) HasFingerprint(uid string) bool {
	for _, known := range a.KnownFingerprints {
		if uid == known.UserID {
			return true
		}
	}

	return false
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

	a.AddFingerprint(fingerprint, uid)
	return nil
}
