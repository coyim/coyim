package config

import (
	"bytes"
	"encoding/hex"
	"errors"
)

type KnownFingerprint struct {
	UserId         string
	FingerprintHex string
	Fingerprint    []byte `json:"-"`
}

func parseFingerprints(c *Config) error {
	var err error
	for i, known := range c.KnownFingerprints {
		c.KnownFingerprints[i].Fingerprint, err = hex.DecodeString(known.FingerprintHex)
		if err != nil {
			return errors.New("xmpp: failed to parse hex fingerprint for " + known.UserId + ": " + err.Error())
		}
	}

	return nil
}

func (c *Config) SerializeFingerprints() {
	for i, known := range c.KnownFingerprints {
		c.KnownFingerprints[i].FingerprintHex = hex.EncodeToString(known.Fingerprint)
	}
}

func (c *Config) UserIdForFingerprint(fpr []byte) string {
	for _, known := range c.KnownFingerprints {
		if bytes.Equal(fpr, known.Fingerprint) {
			return known.UserId
		}
	}

	return ""
}

func (c *Config) AddFingerprint(fpr []byte, uid string) {
	c.KnownFingerprints = append(c.KnownFingerprints, KnownFingerprint{Fingerprint: fpr, UserId: uid})
}

func (c *Config) HasFingerprint(uid string) bool {
	for _, known := range c.KnownFingerprints {
		if uid == known.UserId {
			return true
		}
	}

	return false
}
