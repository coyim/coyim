package utils

import (
	"io/ioutil"
	"log"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
}

type JidXmppSuite struct{}

var _ = Suite(&JidXmppSuite{})

func (s *JidXmppSuite) Test_RemoveResourceFromJid_returnsEverythingBeforeTheSlash(c *C) {
	c.Assert(RemoveResourceFromJid("foo/bar"), Equals, "foo")
	c.Assert(RemoveResourceFromJid("/bar"), Equals, "")
	c.Assert(RemoveResourceFromJid("foo2/"), Equals, "foo2")
	c.Assert(RemoveResourceFromJid("foo3/bar/flux"), Equals, "foo3")
}

func (s *JidXmppSuite) Test_RemoveResourceFromJid_returnsTheWholeStringIfNoSlashesAreInIt(c *C) {
	c.Assert(RemoveResourceFromJid("foo"), Equals, "foo")
	c.Assert(RemoveResourceFromJid("barasdfgdfgdsfgdsfgsdfgdsf"), Equals, "barasdfgdfgdsfgdsfgsdfgdsf")
	c.Assert(RemoveResourceFromJid(""), Equals, "")
}

func (s *JidXmppSuite) Test_DomainFromJid_returnsTheDomain(c *C) {
	c.Assert(DomainFromJid("foo@bar/blarg"), Equals, "bar")
	c.Assert(DomainFromJid("foo@bar2"), Equals, "bar2")
	c.Assert(DomainFromJid("foobar2/blarg"), Equals, "foobar2")
}
