package config

import (
	"bytes"
	"encoding/hex"
	"errors"
)

// KnownFingerprint represents one fingerprint
type KnownFingerprint struct {
	UserID         string
	FingerprintHex string
	Fingerprint    []byte `json:"-"`
	// TODO: add a verified/trusted parameter for each fingerprint
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

// UserIDForFingerprint returns the user ID for the given fingerprint
func (a *Account) UserIDForFingerprint(fpr []byte) string {
	for _, known := range a.KnownFingerprints {
		if bytes.Equal(fpr, known.Fingerprint) {
			return known.UserID
		}
	}

	return ""
}

// AddFingerprint adds a new fingerprint for the given user
func (a *Account) AddFingerprint(fpr []byte, uid string) {
	a.KnownFingerprints = append(a.KnownFingerprints, KnownFingerprint{Fingerprint: fpr, UserID: uid})
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
