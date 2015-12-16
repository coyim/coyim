package config

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"

	"../i18n"
)

// KnownFingerprint represents one fingerprint
type KnownFingerprint struct {
	UserID      string
	Fingerprint []byte
	Untrusted   bool
}

func (k KnownFingerprint) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		UserID         string
		FingerprintHex string
		Untrusted      bool
	}{
		UserID:         k.UserID,
		FingerprintHex: hex.EncodeToString(k.Fingerprint),
		Untrusted:      k.Untrusted,
	})
}

func (k *KnownFingerprint) UnmarshalJSON(data []byte) error {
	v := struct {
		UserID         string
		FingerprintHex string
		Untrusted      bool
	}{}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	k.Fingerprint, err = hex.DecodeString(v.FingerprintHex)
	if err != nil {
		return nil
	}

	k.UserID = v.UserID
	k.Untrusted = v.Untrusted

	return nil
}

// ByNaturalOrder sorts fingerprints according to first the user ID and then the fingerprint
type ByNaturalOrder []*KnownFingerprint

func (s ByNaturalOrder) Len() int { return len(s) }
func (s ByNaturalOrder) Less(i, j int) bool {
	if s[i].UserID == s[j].UserID {
		return bytes.Compare(s[i].Fingerprint, s[j].Fingerprint) == -1
	}

	return s[i].UserID < s[j].UserID
}

func (s ByNaturalOrder) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

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
