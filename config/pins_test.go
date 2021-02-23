package config

import (
	"crypto/x509"
	"sort"

	. "gopkg.in/check.v1"
)

type CertificatePinsSuite struct{}

var _ = Suite(&CertificatePinsSuite{})

func (s *CertificatePinsSuite) Test_CertificatePin_MarshalJSON(c *C) {
	cp := &CertificatePin{
		Subject:         "one",
		Issuer:          "two",
		Fingerprint:     []byte{0x01, 0x02, 0x55},
		FingerprintType: "three",
	}
	res, e := cp.MarshalJSON()
	c.Assert(e, IsNil)
	c.Assert(string(res), Equals, "{\"Subject\":\"one\",\"Issuer\":\"two\",\"FingerprintHex\":\"010255\",\"FingerprintType\":\"three\"}")
}

func (s *CertificatePinsSuite) Test_CertificatePin_UnmarshalJSON_simpleUnmarshallingWorks(c *C) {
	cp := &CertificatePin{}
	e := cp.UnmarshalJSON([]byte("{\"Subject\":\"one\",\"Issuer\":\"two\",\"FingerprintHex\":\"010255\",\"FingerprintType\":\"three\"}"))
	c.Assert(e, IsNil)
	c.Assert(cp.Subject, Equals, "one")
	c.Assert(cp.Issuer, Equals, "two")
	c.Assert(cp.Fingerprint, DeepEquals, []byte{0x01, 0x02, 0x55})
	c.Assert(cp.FingerprintType, Equals, "three")
}

func (s *CertificatePinsSuite) Test_CertificatePin_UnmarshalJSON_invalidJSONGeneratesError(c *C) {
	cp := &CertificatePin{}
	e := cp.UnmarshalJSON([]byte("{\"Subject\":"))
	c.Assert(e, ErrorMatches, ".*unexpected end of JSON input.*")
}

func (s *CertificatePinsSuite) Test_CertificatePin_UnmarshalJSON_invalidFingerprintGeneratesError(c *C) {
	cp := &CertificatePin{}
	e := cp.UnmarshalJSON([]byte("{\"Subject\":\"one\",\"Issuer\":\"two\",\"FingerprintHex\":\"010xQ255\",\"FingerprintType\":\"three\"}"))
	c.Assert(e, ErrorMatches, ".*invalid byte.*")
}

func (s *CertificatePinsSuite) Test_CertificatePin_Matches_comparesASHA1Certificate(c *C) {
	cp := &CertificatePin{
		Subject:         "one",
		Issuer:          "two",
		Fingerprint:     []byte{0x67, 0x37, 0x29, 0x24, 0xc2, 0x70, 0xd4, 0xc2, 0xb4, 0x11, 0x87, 0xfa, 0x5c, 0x54, 0xf7, 0x6c, 0x77, 0x28, 0x12, 0x68},
		FingerprintType: "SHA1",
	}
	cert := &x509.Certificate{Raw: []byte{0x02, 0x3, 0x01, 0x12, 0x75}}
	res := cp.Matches(cert)
	c.Assert(res, Equals, true)
}

func (s *CertificatePinsSuite) Test_CertificatePin_Matches_comparesASHA256Certificate(c *C) {
	cp := &CertificatePin{
		Subject:         "one",
		Issuer:          "two",
		Fingerprint:     []byte{0xd0, 0x4b, 0xb0, 0xdb, 0xbd, 0xe9, 0xf8, 0xa7, 0xe0, 0x68, 0x27, 0xa2, 0x5e, 0x44, 0xa9, 0x58, 0xa7, 0x2d, 0x96, 0x79, 0xa9, 0x3c, 0xe1, 0x24, 0xb2, 0x36, 0x86, 0x02, 0xc7, 0x41, 0xf1, 0xdb},
		FingerprintType: "SHA256",
	}
	cert := &x509.Certificate{Raw: []byte{0x02, 0x3, 0x01, 0x12, 0x75}}
	res := cp.Matches(cert)
	c.Assert(res, Equals, true)
}

func (s *CertificatePinsSuite) Test_CertificatePin_Matches_comparesASHA3_256Certificate(c *C) {
	cp := &CertificatePin{
		Subject:         "one",
		Issuer:          "two",
		Fingerprint:     []byte{0x65, 0x80, 0x02, 0x49, 0xf9, 0xf5, 0xc3, 0x74, 0xaa, 0x90, 0x44, 0x2d, 0xd4, 0x19, 0xb0, 0x51, 0x3d, 0x1a, 0xb8, 0xec, 0x4e, 0xa4, 0x8a, 0x5d, 0x4b, 0x7b, 0xbd, 0xcd, 0x2c, 0x90, 0xad, 0xf5},
		FingerprintType: "SHA3-256",
	}
	cert := &x509.Certificate{Raw: []byte{0x02, 0x3, 0x01, 0x12, 0x75}}
	res := cp.Matches(cert)
	c.Assert(res, Equals, true)
}

func (s *CertificatePinsSuite) Test_CertificatePin_Matches_cantCompareUnknownType(c *C) {
	cp := &CertificatePin{
		Subject:         "one",
		Issuer:          "two",
		Fingerprint:     []byte{0x65, 0x80, 0x02, 0x49, 0xf9, 0xf5, 0xc3, 0x74, 0xaa, 0x90, 0x44, 0x2d, 0xd4, 0x19, 0xb0, 0x51, 0x3d, 0x1a, 0xb8, 0xec, 0x4e, 0xa4, 0x8a, 0x5d, 0x4b, 0x7b, 0xbd, 0xcd, 0x2c, 0x90, 0xad, 0xf5},
		FingerprintType: "MD5",
	}
	cert := &x509.Certificate{Raw: []byte{0x02, 0x3, 0x01, 0x12, 0x75}}
	res := cp.Matches(cert)
	c.Assert(res, Equals, false)
}

func (s *CertificatePinsSuite) Test_CertificatePin_Matches_failsWhenDigestIsIncorrect(c *C) {
	cp := &CertificatePin{
		Subject:         "one",
		Issuer:          "two",
		Fingerprint:     []byte{0x99, 0x99, 0x99, 0x99},
		FingerprintType: "SHA3-256",
	}
	cert := &x509.Certificate{Raw: []byte{0x02, 0x3, 0x01, 0x12, 0x75}}
	res := cp.Matches(cert)
	c.Assert(res, Equals, false)
}

func (s *CertificatePinsSuite) Test_Account_justReturnsIfTheresNoLegacyInformation(c *C) {
	a := &Account{
		LegacyServerCertificateSHA256: "",
	}
	res := a.updateCertificatePins()
	c.Assert(res, Equals, false)
	c.Assert(a.Certificates, HasLen, 0)
}

func (s *CertificatePinsSuite) Test_Account_justReturnsIfTheLegacyCertificateIsNotHex(c *C) {
	a := &Account{
		LegacyServerCertificateSHA256: "q",
	}
	res := a.updateCertificatePins()
	c.Assert(res, Equals, false)
	c.Assert(a.Certificates, HasLen, 0)
}

func (s *CertificatePinsSuite) Test_Account_justReturnsIfTheLegacyCertificateIsNotTheRightLength(c *C) {
	a := &Account{
		LegacyServerCertificateSHA256: "01",
	}
	res := a.updateCertificatePins()
	c.Assert(res, Equals, false)
	c.Assert(a.Certificates, HasLen, 0)
}

func (s *CertificatePinsSuite) Test_Account_addsTheLegacyCertificateToNewCertificates(c *C) {
	a := &Account{
		LegacyServerCertificateSHA256: "0001020304050607101112131415161720121223242526273031323334353637",
	}
	res := a.updateCertificatePins()
	c.Assert(res, Equals, true)
	c.Assert(a.LegacyServerCertificateSHA256, Equals, "")
	c.Assert(a.Certificates, HasLen, 1)
	c.Assert(a.Certificates[0].FingerprintType, Equals, "SHA256")
	c.Assert(a.Certificates[0].Fingerprint, DeepEquals, []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x20, 0x12, 0x12, 0x23, 0x24, 0x25, 0x26, 0x27,
		0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37})
}

func (s *CertificatePinsSuite) Test_CertificatePinsByNaturalOrder_sortsProperlyBasedOnFingerprint(c *C) {
	cpa := &CertificatePin{Fingerprint: []byte{0x01, 0x02}}
	cpb := &CertificatePin{Fingerprint: []byte{0x01, 0x01}}
	cpc := &CertificatePin{Fingerprint: []byte{0x02, 0x01}}
	one := []*CertificatePin{cpa, cpb, cpc}
	sort.Sort(CertificatePinsByNaturalOrder(one))
	c.Assert(one[0], Equals, cpb)
	c.Assert(one[1], Equals, cpa)
	c.Assert(one[2], Equals, cpc)
}
