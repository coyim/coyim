package plain

import (
	"testing"

	"../../sasl"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type SASLPlain struct{}

var _ = Suite(&SASLPlain{})

func (s *SASLPlain) Test(c *C) {
	expected := sasl.Token("\x00foo\x00bar")

	client := Mechanism.NewClient()
	c.Check(client.NeedsMore(), Equals, true)

	client.SetProperty(sasl.AuthID, "foo")
	client.SetProperty(sasl.Password, "bar")

	t, err := client.Step(nil)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, true)
	c.Check(t, DeepEquals, expected)

	t, err = client.Step(nil)
	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, false)
}
