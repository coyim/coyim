package xmpp

import (
	"encoding/base64"

	. "gopkg.in/check.v1"
)

type DigestMD5Suite struct{}

var _ = Suite(&DigestMD5Suite{})

func (s *DigestMD5Suite) TestWithRFC2831TestVector(c *C) {
	digest := digestMD5{
		user:        "chris",
		password:    "secret",
		clientNonce: "OA6MHXh6VqTrRk",

		servType: "imap",
		encoder:  base64.StdEncoding,
	}

	enc := base64.StdEncoding.EncodeToString([]byte(`realm="elwood.innosoft.com",nonce="OA6MG9tEQGm2hh",qop="auth",algorithm=md5-sess,charset=utf-8`))
	err := digest.receive(enc)
	c.Check(err, IsNil)

	dec, _ := base64.StdEncoding.DecodeString(digest.send())

	c.Check(string(dec), Equals, `charset=utf-8,username="chris",realm="elwood.innosoft.com",nonce="OA6MG9tEQGm2hh",nc=00000001,cnonce="OA6MHXh6VqTrRk",digest-uri="imap/elwood.innosoft.com",response=d388dad90d4bbd760a152321f2143af7,qop=auth`)

	enc = base64.StdEncoding.EncodeToString([]byte("rspauth=ea40f60335c427b5527b84dbabcdfffd"))
	err = digest.verifyResponse(enc)
	c.Check(err, IsNil)
}
