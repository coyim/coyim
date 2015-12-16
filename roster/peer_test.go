package roster

import (
	"../xmpp"

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

	p := PeerFrom(re, "")

	c.Assert(p.Jid, g.Equals, "foo@bar.com")
	c.Assert(p.Subscription, g.Equals, "from")
	c.Assert(p.Name, g.Equals, "someone")
	c.Assert(p.Groups, g.DeepEquals, toSet("onegroup", "twogroup"))
}

func (s *PeerXmppSuite) Test_fromSet_returnsTheKeysInTheSet(c *g.C) {
	c.Assert(fromSet(toSet("helo")), g.DeepEquals, []string{"helo"})
	c.Assert(fromSet(toSet()), g.DeepEquals, []string{})
	c.Assert(len(fromSet(toSet("helo1", "helo2", "helo1"))), g.DeepEquals, 2)
}

func (s *PeerXmppSuite) Test_toEntry_ReturnsAnEntryWithTheInformation(c *g.C) {
	p := &Peer{
		Jid:          "foo@bar.com",
		Name:         "something",
		Subscription: "from",
		Groups:       toSet("hello::bar"),
	}

	c.Assert(p.ToEntry().Jid, g.Equals, "foo@bar.com")
	c.Assert(p.ToEntry().Name, g.Equals, "something")
	c.Assert(p.ToEntry().Subscription, g.Equals, "from")
	c.Assert(p.ToEntry().Group, g.DeepEquals, []string{"hello::bar"})
}

func (s *PeerXmppSuite) Test_PeerWithState_createsANewPeer(c *g.C) {
	p := PeerWithState("bla@foo.com/1234", "hmm", "no", "")
	c.Assert(p.Jid, g.Equals, "bla@foo.com")
	c.Assert(p.Name, g.Equals, "")
	c.Assert(p.Status, g.Equals, "hmm")
	c.Assert(p.StatusMsg, g.Equals, "no")
}

func (s *PeerXmppSuite) Test_PeerWithPendingSubscribe_createsNewPeer(c *g.C) {
	p := peerWithPendingSubscribe("bla@foo.com/1234", "223434", "")
	c.Assert(p.Jid, g.Equals, "bla@foo.com")
	c.Assert(p.Name, g.Equals, "")
	c.Assert(p.PendingSubscribeID, g.Equals, "223434")
}

func (s *PeerXmppSuite) Test_NameForPresentation_returnsTheNameIfItExistsButJidOtherwise(c *g.C) {
	c.Assert((&Peer{Name: "foo", Jid: "foo@bar.com"}).NameForPresentation(), g.Equals, "foo")
	c.Assert((&Peer{Jid: "foo@bar.com"}).NameForPresentation(), g.Equals, "foo@bar.com")
}

func (s *PeerXmppSuite) Test_MergeWith_takesTheFirstGroupsIfExists(c *g.C) {
	p1 := &Peer{
		Groups: toSet("one"),
	}
	p2 := &Peer{}

	c.Assert(fromSet(p1.MergeWith(p2).Groups)[0], g.Equals, "one")
}
