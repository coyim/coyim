package roster

import (
	"testing"

	g "gopkg.in/check.v1"
)

func Test(t *testing.T) { g.TestingT(t) }

type ListXmppSuite struct{}

var _ = g.Suite(&ListXmppSuite{})

func (s *ListXmppSuite) Test_New_returnsANewList(c *g.C) {
	l := New()
	c.Assert(l, g.Not(g.IsNil))
	c.Assert(l.peers, g.Not(g.IsNil))
}

func (s *ListXmppSuite) Test_Remove_doesNothingWhenAskedToRemoveEntryNotInList(c *g.C) {
	l := New()
	l.peers["foo@bar.com"] = &Peer{}

	res, rem := l.Remove("bar@foo.com")

	c.Assert(rem, g.Equals, false)
	c.Assert(res, g.IsNil)
	c.Assert(len(l.peers), g.Equals, 1)
}

func (s *ListXmppSuite) Test_Remove_removesAnEntryIfInTheList(c *g.C) {
	l := New()
	l.peers["foo@bar.com"] = &Peer{}
	l.peers["bar@foo.com"] = &Peer{Name: "me"}

	res, rem := l.Remove("bar@foo.com/somewhere")

	c.Assert(rem, g.Equals, true)
	c.Assert(res.Name, g.Equals, "me")
	c.Assert(len(l.peers), g.Equals, 1)
}

func (s *ListXmppSuite) Test_AddOrReplace_addsTheEntryIfNotInTheList(c *g.C) {
	l := New()
	p := &Peer{Jid: "somewhere", Name: "something"}

	res := l.AddOrReplace(p)

	c.Assert(res, g.Equals, true)
	c.Assert(len(l.peers), g.Equals, 1)
	c.Assert(l.peers["somewhere"], g.Equals, p)
}

func (s *ListXmppSuite) Test_AddOrReplace_replacesTheEntryIfInTheList(c *g.C) {
	l := New()
	p1 := &Peer{Jid: "somewhere", Name: "something", Groups: toSet("hello"), Subscription: "from"}
	l.peers["somewhere"] = p1

	p2 := &Peer{Jid: "somewhere", Name: "something2", Groups: toSet("goodbye")}
	res := l.AddOrReplace(p2)

	c.Assert(res, g.Equals, false)
	c.Assert(len(l.peers), g.Equals, 1)
	c.Assert(l.peers["somewhere"], g.Equals, p2)
}

func (s *ListXmppSuite) Test_AddOrMerge_addsTheEntryIfNotInTheList(c *g.C) {
	l := New()
	p := &Peer{Jid: "somewhere", Name: "something"}

	res := l.AddOrMerge(p)

	c.Assert(res, g.Equals, true)
	c.Assert(len(l.peers), g.Equals, 1)
	c.Assert(l.peers["somewhere"], g.Equals, p)
}

func (s *ListXmppSuite) Test_AddOrReplace_mergesTheEntriesIfInTheList(c *g.C) {
	l := New()
	p1 := &Peer{Jid: "somewhere", Name: "something", Groups: toSet("hello"), Subscription: "from"}
	l.peers["somewhere"] = p1

	p2 := &Peer{Jid: "somewhere", Name: "something2", Groups: toSet("goodbye")}
	res := l.AddOrMerge(p2)

	c.Assert(res, g.Equals, false)
	c.Assert(len(l.peers), g.Equals, 1)
	c.Assert(*l.peers["somewhere"], g.DeepEquals, Peer{Jid: "somewhere", Name: "something2", Groups: toSet("goodbye"), Subscription: "from"})
}

func (s *ListXmppSuite) Test_ToSlice_createsASliceOfTheContentSortedAlphabetically(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: "foo@somewhere.com"})
	l.AddOrMerge(&Peer{Jid: "foo@somewhen.com"})
	l.AddOrMerge(&Peer{Jid: "bar@somewhere.com"})

	c.Assert(l.ToSlice(), g.DeepEquals, []*Peer{
		&Peer{Jid: "bar@somewhere.com"},
		&Peer{Jid: "foo@somewhen.com"},
		&Peer{Jid: "foo@somewhere.com"},
	})
}

func (s *ListXmppSuite) Test_Iter_yieldsEachEntry(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: "foo@somewhere.com"})
	l.AddOrMerge(&Peer{Jid: "foo@somewhen.com"})
	l.AddOrMerge(&Peer{Jid: "bar@somewhere.com"})

	called := 0

	l.Iter(func(ix int, p *Peer) {
		called++

		switch ix {
		case 0:
			c.Assert(p, g.DeepEquals, &Peer{Jid: "bar@somewhere.com"})
		case 1:
			c.Assert(p, g.DeepEquals, &Peer{Jid: "foo@somewhen.com"})
		case 2:
			c.Assert(p, g.DeepEquals, &Peer{Jid: "foo@somewhere.com"})
		}
	})

	c.Assert(called, g.DeepEquals, 3)
}

func (s *ListXmppSuite) Test_Unsubscribed_whenDoesntExist(c *g.C) {
	l := New()
	l.Unsubscribed("foo@bar.com")
}

func (s *ListXmppSuite) Test_Unsubscribed_whenExist(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: "foo@bar.com", Subscription: "both", Asked: true, PendingSubscribeId: "foo"})
	l.AddOrMerge(&Peer{Jid: "foo2@bar.com", Subscription: "to"})
	l.AddOrMerge(&Peer{Jid: "foo3@bar.com", Subscription: "from"})

	l.Unsubscribed("foo@bar.com/123")
	c.Assert(l.peers["foo@bar.com"].Subscription, g.Equals, "from")
	c.Assert(l.peers["foo@bar.com"].Asked, g.Equals, false)
	c.Assert(l.peers["foo@bar.com"].PendingSubscribeId, g.Equals, "")

	l.Unsubscribed("foo2@bar.com/123")
	c.Assert(l.peers["foo2@bar.com"].Subscription, g.Equals, "none")

	l.Unsubscribed("foo3@bar.com/123")
	c.Assert(l.peers["foo3@bar.com"].Subscription, g.Equals, "from")
}

func (s *ListXmppSuite) Test_Subscribed_whenDoesntExist(c *g.C) {
	l := New()
	l.Subscribed("foo@bar.com")
	c.Assert(len(l.peers), g.Equals, 0)
}

func (s *ListXmppSuite) Test_Subscribed_whenExist(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: "foo@bar.com", Subscription: "from", Asked: true, PendingSubscribeId: "foo"})
	l.AddOrMerge(&Peer{Jid: "foo2@bar.com", Subscription: "none"})
	l.AddOrMerge(&Peer{Jid: "foo3@bar.com", Subscription: ""})
	l.AddOrMerge(&Peer{Jid: "foo4@bar.com", Subscription: "both"})

	l.Subscribed("foo@bar.com/123")
	c.Assert(l.peers["foo@bar.com"].Subscription, g.Equals, "both")
	c.Assert(l.peers["foo@bar.com"].Asked, g.Equals, false)
	c.Assert(l.peers["foo@bar.com"].PendingSubscribeId, g.Equals, "")

	l.Subscribed("foo2@bar.com/123")
	c.Assert(l.peers["foo2@bar.com"].Subscription, g.Equals, "to")

	l.Subscribed("foo3@bar.com/123")
	c.Assert(l.peers["foo3@bar.com"].Subscription, g.Equals, "to")

	l.Subscribed("foo4@bar.com/123")
	c.Assert(l.peers["foo4@bar.com"].Subscription, g.Equals, "both")
}

func (s *ListXmppSuite) Test_GetPendingSubscribe_returnsThePendingSubscribeIfExists(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: "foo@bar.com", PendingSubscribeId: "foo"})
	l.AddOrMerge(&Peer{Jid: "foo2@bar.com"})

	v, k := l.GetPendingSubscribe("none@foo.com")
	c.Assert(k, g.Equals, false)

	v, k = l.GetPendingSubscribe("foo@bar.com/bar")
	c.Assert(k, g.Equals, true)
	c.Assert(v, g.Equals, "foo")

	v, k = l.GetPendingSubscribe("foo2@bar.com/bar")
	c.Assert(k, g.Equals, false)
	c.Assert(v, g.Equals, "")
}

func (s *ListXmppSuite) Test_RemovePendingSubscribe_removesThePendingSubscribe(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: "foo@bar.com", PendingSubscribeId: "foo"})
	l.AddOrMerge(&Peer{Jid: "foo2@bar.com"})

	v, k := l.RemovePendingSubscribe("none@foo.com")
	c.Assert(k, g.Equals, false)

	v, k = l.RemovePendingSubscribe("foo@bar.com/bar")
	c.Assert(k, g.Equals, true)
	c.Assert(v, g.Equals, "foo")
	c.Assert(l.peers["foo@bar.com"].PendingSubscribeId, g.Equals, "")

	v, k = l.RemovePendingSubscribe("foo2@bar.com/bar")
	c.Assert(k, g.Equals, false)
	c.Assert(v, g.Equals, "")
}

func (s *ListXmppSuite) Test_SubscribeRequest_addsTheSubscribeID(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: "foo@bar.com"})

	l.SubscribeRequest("fox@bar.com/hmm", "something")
	c.Assert(l.peers["fox@bar.com"].PendingSubscribeId, g.Equals, "something")

	l.SubscribeRequest("foo@bar.com/hmm2", "something3")
	c.Assert(l.peers["foo@bar.com"].PendingSubscribeId, g.Equals, "something3")
}

func (s *ListXmppSuite) Test_StateOf_returnsState(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: "foo@bar.com", Status: "bla", StatusMsg: "hmm"})

	st, sm, k := l.StateOf("hmm.bar@bar.com")
	c.Assert(k, g.Equals, false)

	st, sm, k = l.StateOf("foo@bar.com/aha")
	c.Assert(k, g.Equals, true)
	c.Assert(st, g.Equals, "bla")
	c.Assert(sm, g.Equals, "hmm")
}

func (s *ListXmppSuite) Test_PeerBecameUnavailable_setsTheOfflineState(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: "foo@bar.com", Offline: false})

	res := l.PeerBecameUnavailable("hmm@bar.com/foo")
	c.Assert(res, g.Equals, false)

	res = l.PeerBecameUnavailable("foo@bar.com/foo2")
	c.Assert(res, g.Equals, true)
	c.Assert(l.peers["foo@bar.com"].Offline, g.Equals, true)
}

func (s *ListXmppSuite) Test_PeerPresenceUpdate_sometimesUpdatesNonExistantPeers(c *g.C) {
	l := New()

	res := l.PeerPresenceUpdate("foo@bar.com/hmm", "hello", "goodbye")
	c.Assert(res, g.Equals, true)
	c.Assert(l.peers["foo@bar.com"].Status, g.Equals, "hello")
	c.Assert(l.peers["foo@bar.com"].StatusMsg, g.Equals, "goodbye")

	res = l.PeerPresenceUpdate("foo2@bar.com/hmm", "xa", "goodbye")
	c.Assert(res, g.Equals, false)
	c.Assert(l.peers["foo2@bar.com"], g.IsNil)

	res = l.PeerPresenceUpdate("foo3@bar.com/hmm", "away", "goodbye")
	c.Assert(res, g.Equals, false)
	c.Assert(l.peers["foo3@bar.com"], g.IsNil)

}

func (s *ListXmppSuite) Test_PeerPresenceUpdate_updatesPreviouslyKnownPeer(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: "foo@bar.com", Offline: true})
	l.AddOrMerge(&Peer{Jid: "foo2@bar.com", Offline: false, Status: "dnd", StatusMsg: "working"})

	res := l.PeerPresenceUpdate("foo@bar.com/hmm", "hello", "goodbye")
	c.Assert(res, g.Equals, true)
	c.Assert(l.peers["foo@bar.com"].Status, g.Equals, "hello")
	c.Assert(l.peers["foo@bar.com"].StatusMsg, g.Equals, "goodbye")
	c.Assert(l.peers["foo@bar.com"].Offline, g.Equals, false)

	res = l.PeerPresenceUpdate("foo2@bar.com/hmm", "dnd", "working")
	c.Assert(res, g.Equals, false)
	c.Assert(l.peers["foo2@bar.com"].Status, g.Equals, "dnd")
	c.Assert(l.peers["foo2@bar.com"].StatusMsg, g.Equals, "working")
	c.Assert(l.peers["foo2@bar.com"].Offline, g.Equals, false)
}

func (s *ListXmppSuite) Test_Clear_clearsTheList(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: "foo@bar.com"})

	l.Clear()

	c.Assert(len(l.peers), g.Equals, 0)
}
