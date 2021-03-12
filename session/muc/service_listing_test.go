package muc

import (
	"github.com/coyim/coyim/xmpp/jid"

	. "gopkg.in/check.v1"
)

func (s *MucSuite) Test_NewServiceListing_createsListing(c *C) {
	vv := jid.Parse("bla@foo.hmm")
	sl := NewServiceListing(vv, "hello")
	c.Assert(sl.Jid, Equals, vv)
	c.Assert(sl.Name, Equals, "hello")
}
