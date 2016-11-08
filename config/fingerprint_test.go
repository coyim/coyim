package config

import (
	"encoding/hex"
	"encoding/json"

	. "gopkg.in/check.v1"
)

type FingerprintXMPPSuite struct{}

var _ = Suite(&FingerprintXMPPSuite{})

func (s *FingerprintXMPPSuite) Test_formatFingerprint(c *C) {
	testVal := []byte{0x5d, 0xfc, 0x9e, 0x41, 0x6b, 0xf7, 0x83, 0xea, 0x14, 0x90, 0xb8, 0x16, 0x9b, 0x86, 0x68, 0x21, 0xb5, 0x2e, 0xbb, 0xb7}

	res := FormatFingerprint(testVal)

	c.Assert(res, Equals, "5DFC9E41 6BF783EA 1490B816 9B866821 B52EBBB7")
}

func (s *FingerprintXMPPSuite) Test_SerializeAndDeserialize(c *C) {
	var jsonBlob = []byte(`{
	"UserID": "user@coy.im",
	"FingerprintHex": "a0cbc473380411c659e87088031035ed464c9270",
	"Untrusted": true
}`)

	fp, _ := hex.DecodeString("a0cbc473380411c659e87088031035ed464c9270")
	expected := KnownFingerprint{
		UserID:      "user@coy.im",
		Fingerprint: fp,
		Untrusted:   true,
	}

	dest := KnownFingerprint{}
	err := json.Unmarshal(jsonBlob, &dest)

	c.Check(err, IsNil)
	c.Check(dest, DeepEquals, expected)

	marshal, err := json.MarshalIndent(dest, "", "\t")
	c.Check(err, IsNil)
	c.Check(string(marshal), Equals, string(jsonBlob))
}
