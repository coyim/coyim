package data

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glib_mock"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
}

type JidXMPPSuite struct{}

var _ = Suite(&JidXMPPSuite{})

func (s *JidXMPPSuite) Test_ParseJID(c *C) {
	c.Assert(ParseJID(""), DeepEquals, JIDDomain(""))
	c.Assert(ParseJID("foo.bar"), DeepEquals, JIDDomain("foo.bar"))
	c.Assert(ParseJID("ola@foo.bar2"), DeepEquals, bareJID("ola@foo.bar2"))
	c.Assert(ParseJID("ola@foo.bar3/foo"), DeepEquals, fullJID("ola@foo.bar3/foo"))
	c.Assert(ParseJID("foo/bar"), DeepEquals, domainWithResource("foo/bar"))
	c.Assert(ParseJID("foo3/bar/flux"), DeepEquals, domainWithResource("foo3/bar/flux"))
	c.Assert(ParseJID("foo3/bar/flux").EnsureNoResource(), DeepEquals, JIDDomain("foo3"))
	c.Assert(ParseJID("foo3/bar/flux").(JIDWithResource).Resource(), DeepEquals, JIDResource("bar/flux"))
}
