package roster

import (
	"github.com/twstrike/coyim/xmpp"

	g "gopkg.in/check.v1"
)

type PeerXmppSuite struct{}

var _ = g.Suite(&PeerXmppSuite{})

func (s *PeerXmppSuite) Test_PeerFrom_returnsANewPeerWithTheSameInformation(c *g.C) {
	re := xmpp.RosterEntry{
		Jid:          "foo@bar.com",
		Subscription: "from",
		Name:         "someone",
		Group: []string{
			"onegroup",
			"twogroup",
		},
	}

	p := PeerFrom(re)

	c.Assert(p.jid, g.Equals, "foo@bar.com")
	c.Assert(p.subscription, g.Equals, "from")
	c.Assert(p.name, g.Equals, "someone")
	c.Assert(p.groups, g.DeepEquals, toSet("onegroup", "twogroup"))
}
