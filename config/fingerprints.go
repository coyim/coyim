package config

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
)

// KnownFingerprint represents one fingerprint
type KnownFingerprint struct {
	UserID      string
	Fingerprint []byte
	Untrusted   bool
}

// MarshalJSON is used to create a JSON representation of this known fingerprint
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

// UnmarshalJSON is used to parse the JSON representation of a known fingerprint
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

// LegacyByNaturalOrder sorts fingerprints according to first the user ID and then the fingerprint
type LegacyByNaturalOrder []*KnownFingerprint

func (s LegacyByNaturalOrder) Len() int { return len(s) }
func (s LegacyByNaturalOrder) Less(i, j int) bool {
	if s[i].UserID == s[j].UserID {
		return bytes.Compare(s[i].Fingerprint, s[j].Fingerprint) == -1
	}

	return s[i].UserID < s[j].UserID
}

func (s LegacyByNaturalOrder) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
