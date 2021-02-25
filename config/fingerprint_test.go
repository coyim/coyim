package config

import (
	"encoding/hex"
	"encoding/json"
	"sort"

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

func (s *FingerprintXMPPSuite) Test_KnownFingerprint_UnmarshalJSON_failsOnBadJSON(c *C) {
	jsonBlob := []byte(`{
	"User`)

	kf := &KnownFingerprint{}
	err := kf.UnmarshalJSON(jsonBlob)

	c.Assert(err, ErrorMatches, "unexpected end of JSON input")
}

func (s *FingerprintXMPPSuite) Test_KnownFingerprint_UnmarshalJSON_failsOnBadHexInFingerprint(c *C) {
	jsonBlob := []byte(`{
	"UserID": "user@coy.im",
	"FingerprintHex": "QQQa0cbc473380411c659e87088031035ed464c9270",
	"Untrusted": true
}`)

	kf := &KnownFingerprint{}
	err := kf.UnmarshalJSON(jsonBlob)

	c.Assert(err, ErrorMatches, "encoding/hex.*")
}

func (s *FingerprintXMPPSuite) Test_KnownFingerprint_LegacyByNaturalOrder(c *C) {
	fp1 := &KnownFingerprint{UserID: "one1", Fingerprint: []byte{0x01, 0x02}}
	fp2 := &KnownFingerprint{UserID: "two1", Fingerprint: []byte{0x03, 0x02}}
	fp3 := &KnownFingerprint{UserID: "two1", Fingerprint: []byte{0x03, 0x01}}
	fp4 := &KnownFingerprint{UserID: "four", Fingerprint: []byte{0x01, 0x02}}
	fp5 := &KnownFingerprint{UserID: "six", Fingerprint: []byte{0x01, 0x02}}

	one := []*KnownFingerprint{fp1, fp2, fp3, fp4, fp5}
	sort.Sort(LegacyByNaturalOrder(one))
	c.Assert(one[0], Equals, fp4)
	c.Assert(one[1], Equals, fp1)
	c.Assert(one[2], Equals, fp5)
	c.Assert(one[3], Equals, fp3)
	c.Assert(one[4], Equals, fp2)
}
