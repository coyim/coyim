package xmpp

import (
	. "gopkg.in/check.v1"
)

type DNSXmppSuite struct{}

var _ = Suite(&DNSXmppSuite{})

// WARNING: these tests require a real live connection to the Internet. Not so good...

func (s *DNSXmppSuite) Test_Resolve_resolvesCorrectly(c *C) {
	hostport, err := Resolve("olabini.se")
	c.Assert(err, IsNil)
	c.Assert(hostport[0], Equals, "xmpp.olabini.se:5222")
}

func (s *DNSXmppSuite) Test_Resolve_handlesErrors(c *C) {
	_, err := Resolve("doesntexist.olabini.se")
	c.Assert(err.Error(), Matches, "lookup _xmpp-client._tcp.doesntexist.olabini.se.*?: no such host")
}
