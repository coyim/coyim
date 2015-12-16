package digestmd5

import (
	"testing"

	"../../sasl"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type DigestMD5 struct{}

var _ = Suite(&DigestMD5{})

func (s *DigestMD5) Test(c *C) {
	client := Mechanism.NewClient()
	c.Check(client.NeedsMore(), Equals, true)

	client.SetProperty(sasl.AuthID, "chris")
	client.SetProperty(sasl.Password, "secret")
	client.SetProperty(sasl.Service, "imap")
	//client.SetProperty(sasl.Realm, "elwood.innosoft.com")
	client.SetProperty(sasl.QOP, "auth")

	client.SetProperty(sasl.ClientNonce, "OA6MHXh6VqTrRk")

	t, err := client.Step(nil)
	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, true)
	c.Check(t, IsNil)

	rec := sasl.Token(`realm="elwood.innosoft.com",nonce="OA6MG9tEQGm2hh",qop="auth",algorithm=md5-sess,charset=utf-8`)
	t, err = client.Step(rec)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, true)
	c.Check(t, DeepEquals, sasl.Token(`charset=utf-8,username="chris",realm="elwood.innosoft.com",nonce="OA6MG9tEQGm2hh",nc=00000001,cnonce="OA6MHXh6VqTrRk",digest-uri="imap/elwood.innosoft.com",response=d388dad90d4bbd760a152321f2143af7,qop=auth`))

	rec = sasl.Token("rspauth=ea40f60335c427b5527b84dbabcdfffd")
	t, err = client.Step(rec)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, true)
	c.Check(t, IsNil)

	t, err = client.Step(nil)
	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, false)
	c.Check(t, IsNil)
}
