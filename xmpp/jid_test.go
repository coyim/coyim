package xmpp

import (
	. "gopkg.in/check.v1"
)

type JidSuite struct{}

var _ = Suite(&JidSuite{})

func (s *JidSuite) Test_ParseJID(c *C) {
	c.Assert(ParseJID("coyim"), DeepEquals, &JID{
		DomainPart: "coyim",
	})

	c.Assert(ParseJID("local@coyim"), DeepEquals, &JID{
		LocalPart:  "local",
		DomainPart: "coyim",
	})

	c.Assert(ParseJID("coyim/resource"), DeepEquals, &JID{
		DomainPart:   "coyim",
		ResourcePart: "resource",
	})

	c.Assert(ParseJID("local@coyim/resource"), DeepEquals, &JID{
		LocalPart:    "local",
		DomainPart:   "coyim",
		ResourcePart: "resource",
	})

	c.Assert(ParseJID("local@co@yim/reso/urce"), DeepEquals, &JID{
		LocalPart:    "local",
		DomainPart:   "co@yim",
		ResourcePart: "reso/urce",
	})
}
