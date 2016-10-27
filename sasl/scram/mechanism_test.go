package scram

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/sasl"
	"github.com/twstrike/gotk3adapter/glib_mock"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
}

type ScramSuite struct{}

var _ = Suite(&ScramSuite{})

func (s *ScramSuite) TestScramWithRFC5802TestVector(c *C) {
	client := Mechanism.NewClient()

	client.SetProperty(sasl.AuthID, "user")
	client.SetProperty(sasl.Password, "pencil")
	client.SetProperty(sasl.ClientNonce, "fyko+d2lbbFgONRv9qkxdawL")

	t, err := client.Step(nil)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, true)
	c.Check(t, DeepEquals, sasl.Token(`n,,n=user,r=fyko+d2lbbFgONRv9qkxdawL`))

	rec := sasl.Token("r=fyko+d2lbbFgONRv9qkxdawL3rfcNHYJY1ZVvWVs7j,s=QSXCR+Q6sek8bf92,i=4096")
	t, err = client.Step(rec)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, true)
	c.Check(t, DeepEquals, sasl.Token(`c=biws,r=fyko+d2lbbFgONRv9qkxdawL3rfcNHYJY1ZVvWVs7j,p=v0X8v3Bz2T0CJGbJQyF0X+HI4Ts=`))

	rec = sasl.Token("v=rmF9pqV8S7suAoZWja4dJRkFsKQ=")
	t, err = client.Step(rec)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, false)
	c.Check(t, IsNil)
}
