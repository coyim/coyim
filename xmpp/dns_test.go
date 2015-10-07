package xmpp

import (
	. "gopkg.in/check.v1"
)

type DnsXmppSuite struct{}

var _ = Suite(&DnsXmppSuite{})

// WARNING: these tests require a real live connection to the Internet. Not so good...

func (s *DnsXmppSuite) Test_Resolve_resolvesCorrectly(c *C) {
	host, port, err := Resolve("olabini.se")
	c.Assert(err, IsNil)
	c.Assert(host, Equals, "xmpp.olabini.se.")
	c.Assert(port, Equals, uint16(5222))
}

func (s *DnsXmppSuite) Test_Resolve_handlesErrors(c *C) {
	_, _, err := Resolve("doesntexist.olabini.se")
	c.Assert(err.Error(), Matches, "lookup _xmpp-client._tcp.doesntexist.olabini.se.*?: no such host")
}
