package plain

import (
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/sasl"
	"github.com/coyim/gotk3adapter/glib_mock"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
}

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

	expected = sasl.Token(nil)

	t, err = client.Step(nil)
	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, false)
	c.Check(t, DeepEquals, expected)
}
