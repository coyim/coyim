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
}

func parseFingerprints(c *Config) error {
	var err error
	for i, known := range c.KnownFingerprints {
		c.KnownFingerprints[i].Fingerprint, err = hex.DecodeString(known.FingerprintHex)
		if err != nil {
			return errors.New("xmpp: failed to parse hex fingerprint for " + known.UserID + ": " + err.Error())
		}
	}

	return nil
}

func (c *Config) serializeFingerprints() {
	for i, known := range c.KnownFingerprints {
		c.KnownFingerprints[i].FingerprintHex = hex.EncodeToString(known.Fingerprint)
	}
}

// UserIDForFingerprint returns the user ID for the given fingerprint
func (c *Config) UserIDForFingerprint(fpr []byte) string {
	for _, known := range c.KnownFingerprints {
		if bytes.Equal(fpr, known.Fingerprint) {
			return known.UserID
		}
	}

	return ""
}

// AddFingerprint adds a new fingerprint for the given user
func (c *Config) AddFingerprint(fpr []byte, uid string) {
	c.KnownFingerprints = append(c.KnownFingerprints, KnownFingerprint{Fingerprint: fpr, UserID: uid})
}

// HasFingerprint returns true if we have the fingerprint for the given user
func (c *Config) HasFingerprint(uid string) bool {
	for _, known := range c.KnownFingerprints {
		if uid == known.UserID {
			return true
		}
	}

	return false
}
