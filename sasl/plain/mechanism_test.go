package plain

import (
	"github.com/coyim/coyim/sasl"
	. "gopkg.in/check.v1"
)

type SASLPlain struct{}

var _ = Suite(&SASLPlain{})

func (s *SASLPlain) Test(c *C) {
	expected := sasl.Token("\x00foo\x00bar")

	client := Mechanism.NewClient()
	c.Check(client.NeedsMore(), Equals, true)

	_ = client.SetProperty(sasl.AuthID, "foo")
	_ = client.SetProperty(sasl.Password, "bar")

	t, err := client.Step(nil)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, true)
	c.Check(t, DeepEquals, expected)

	expected = sasl.Token(nil)

	t, err = client.Step(nil)
	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, false)
	c.Check(t, DeepEquals, expected)
}

func (s *SASLPlain) Test_Register(c *C) {
	Register()
	c.Assert(sasl.ClientSupport("PLAIN"), Equals, true)
}

func (s *SASLPlain) Test_plain_SetProperty_failsOnUnknownProperty(c *C) {
	p := &plain{}
	c.Assert(p.SetProperty(sasl.Property(42), "else"), Equals, sasl.ErrUnsupportedProperty)
}

func (s *SASLPlain) Test_plain_SetChannelBinding_doesNothing(c *C) {
	p := &plain{}
	p.SetChannelBinding(nil)
}
