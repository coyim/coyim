package xmpp

import . "gopkg.in/check.v1"

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

func (s *JidXmppSuite) Test_domainFromJid_returnsTheDomain(c *C) {
	c.Assert(domainFromJid("foo@bar/blarg"), Equals, "bar")
	c.Assert(domainFromJid("foo@bar2"), Equals, "bar2")
	c.Assert(domainFromJid("foobar2/blarg"), Equals, "foobar2")
}
