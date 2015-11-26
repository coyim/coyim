package xmpp

import (
	"encoding/base64"

	. "gopkg.in/check.v1"
)

type ScramSuite struct{}

var _ = Suite(&ScramSuite{})

func (s *ScramSuite) TestScramNormalizesPassword(c *C) {
	// From: libidn-1.9/tests/tst_stringprep.c
	// See RFC 4013, section 3
	testCases := []struct {
		raw        string
		normalized string
	}{
		{"I\xC2\xADX", "IX"},
		{"user", "user"},
		{"USER", "USER"},
		{"\xC2\xAA", "a"},
		{"x\xC2\xADy", "xy"},
		{"\xE2\x85\xA3", "IV"},
		{"\xE2\x85\xA8", "IX"},
		//They should error becuase they have forbidden chars
		//{"\x07", ""},      //should error
		//{"\xD8\xA71", ""}, //shold error
	}

	for _, test := range testCases {
		scram := scramClient{
			password: test.raw,
		}
		normalized, _ := scram.normalizedPassword()
		c.Check(normalized, Equals, test.normalized)
	}
}

func (s *ScramSuite) TestScramWithRFC5802TestVector(c *C) {

	scram := scramClient{
		user:        "user",
		password:    "pencil",
		clientNonce: "fyko+d2lbbFgONRv9qkxdawL",
	}

	encoding := base64.StdEncoding
	dec, _ := encoding.DecodeString(scram.firstMessage())
	c.Check(string(dec), Equals, "n,,n=user,r=fyko+d2lbbFgONRv9qkxdawL")

	enc := encoding.EncodeToString([]byte("r=fyko+d2lbbFgONRv9qkxdawL3rfcNHYJY1ZVvWVs7j,s=QSXCR+Q6sek8bf92,i=4096"))
	err := scram.receive(enc)
	c.Check(err, IsNil)

	reply, serverAuth, err := scram.reply()

	c.Check(err, IsNil)

	dec, _ = encoding.DecodeString(reply)
	c.Check(string(dec), Equals, "c=biws,r=fyko+d2lbbFgONRv9qkxdawL3rfcNHYJY1ZVvWVs7j,p=v0X8v3Bz2T0CJGbJQyF0X+HI4Ts=")

	dec, _ = encoding.DecodeString(serverAuth)
	c.Check(string(dec), Equals, "v=rmF9pqV8S7suAoZWja4dJRkFsKQ=")
}
