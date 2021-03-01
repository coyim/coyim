package digestmd5

import (
	"github.com/coyim/coyim/sasl"

	. "gopkg.in/check.v1"
)

type DigestMD5Suite struct{}

var _ = Suite(&DigestMD5Suite{})

func (s *DigestMD5Suite) Test(c *C) {
	client := Mechanism.NewClient()
	c.Check(client.NeedsMore(), Equals, true)

	_ = client.SetProperty(sasl.AuthID, "chris")
	_ = client.SetProperty(sasl.Password, "secret")
	_ = client.SetProperty(sasl.Service, "imap")
	//client.SetProperty(sasl.Realm, "elwood.innosoft.com")
	_ = client.SetProperty(sasl.QOP, "auth")

	_ = client.SetProperty(sasl.ClientNonce, "OA6MHXh6VqTrRk")

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

func (s *DigestMD5Suite) Test_Register(c *C) {
	Register()
	c.Assert(sasl.ClientSupport("DIGEST-MD5"), Equals, true)
}
func (s *DigestMD5Suite) Test_digestMD5_SetChannelBinding_doesNothing(c *C) {
	p := &digestMD5{}
	p.SetChannelBinding(nil)
}
