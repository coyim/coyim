package config

import (
	"bytes"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"sort"

	"github.com/coyim/coyim/digests"
)

// CertificatePinForSerialization represents a certificate pin in its serialized form
type CertificatePinForSerialization struct {
	Subject         string `json:",omitempty"`
	Issuer          string `json:",omitempty"`
	FingerprintHex  string
	FingerprintType string
}

// CertificatePin represents a known certificate hash to accept as a given
type CertificatePin struct {
	Subject         string
	Issuer          string
	Fingerprint     []byte
	FingerprintType string
}

// MarshalJSON is used to create a JSON representation of this certificate pin
func (v *CertificatePin) MarshalJSON() ([]byte, error) {
	return json.Marshal(CertificatePinForSerialization{
		Subject:         v.Subject,
		Issuer:          v.Issuer,
		FingerprintHex:  hex.EncodeToString(v.Fingerprint),
		FingerprintType: v.FingerprintType,
	})
}

// UnmarshalJSON is used to parse the JSON representation of a certificate pin
func (v *CertificatePin) UnmarshalJSON(data []byte) error {
	vz := CertificatePinForSerialization{}
	err := json.Unmarshal(data, &vz)
	if err != nil {
		return err
	}

	v.Fingerprint, err = hex.DecodeString(vz.FingerprintHex)
	if err != nil {
		return nil
	}

	v.Subject = vz.Subject
	v.Issuer = vz.Issuer
	v.FingerprintType = vz.FingerprintType

	return nil
}

// CertificatePinsByNaturalOrder sorts certificate pins by the fingerprints
type CertificatePinsByNaturalOrder []*CertificatePin

func (s CertificatePinsByNaturalOrder) Len() int { return len(s) }
func (s CertificatePinsByNaturalOrder) Less(i, j int) bool {
	return bytes.Compare(s[i].Fingerprint, s[j].Fingerprint) == -1
}
func (s CertificatePinsByNaturalOrder) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (a *Account) updateCertificatePins() bool {
	if len(a.LegacyServerCertificateSHA256) == 0 {
		return false
	}

	certSHA256, err := hex.DecodeString(a.LegacyServerCertificateSHA256)
	if err == nil && len(certSHA256) == 32 {
		a.Certificates = append(a.Certificates, &CertificatePin{
			Subject:         "",
			Issuer:          "",
			FingerprintType: "SHA256",
			Fingerprint:     certSHA256,
		})
		sort.Sort(CertificatePinsByNaturalOrder(a.Certificates))
		a.LegacyServerCertificateSHA256 = ""
		return true
	}

	return false
}

// Matches returns true if this pin matches the given certificate
func (v *CertificatePin) Matches(cert *x509.Certificate) bool {
	r := cert.Raw
	var dig []byte
	switch v.FingerprintType {
	case "SHA1":
		dig = digests.Sha1(r)
	case "SHA256":
		dig = digests.Sha256(r)
	case "SHA3-256":
		dig = digests.Sha3_256(r)
	default:
		return false
	}

	return bytes.Equal(dig, v.Fingerprint)
}
