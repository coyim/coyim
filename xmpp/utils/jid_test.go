package utils

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/gotk3adapter/glib_mock"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
}

type JidXMPPSuite struct{}

var _ = Suite(&JidXMPPSuite{})

func (s *JidXMPPSuite) Test_RemoveResourceFromJid_returnsEverythingBeforeTheSlash(c *C) {
	c.Assert(RemoveResourceFromJid("foo/bar"), Equals, "foo")
	c.Assert(RemoveResourceFromJid("/bar"), Equals, "")
	c.Assert(RemoveResourceFromJid("foo2/"), Equals, "foo2")
	c.Assert(RemoveResourceFromJid("foo3/bar/flux"), Equals, "foo3")
}

func (s *JidXMPPSuite) Test_RemoveResourceFromJid_returnsTheWholeStringIfNoSlashesAreInIt(c *C) {
	c.Assert(RemoveResourceFromJid("foo"), Equals, "foo")
	c.Assert(RemoveResourceFromJid("barasdfgdfgdsfgdsfgsdfgdsf"), Equals, "barasdfgdfgdsfgdsfgsdfgdsf")
	c.Assert(RemoveResourceFromJid(""), Equals, "")
}

func (s *JidXMPPSuite) Test_DomainFromJid_returnsTheDomain(c *C) {
	c.Assert(DomainFromJid("foo@bar/blarg"), Equals, "bar")
	c.Assert(DomainFromJid("foo@bar2"), Equals, "bar2")
	c.Assert(DomainFromJid("foobar2/blarg"), Equals, "foobar2")
}
