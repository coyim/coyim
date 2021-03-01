package roster

import (
	"github.com/coyim/coyim/xmpp/jid"

	g "gopkg.in/check.v1"
)

type ListSuite struct{}

var _ = g.Suite(&ListSuite{})

func (s *ListSuite) Test_New_returnsANewList(c *g.C) {
	l := New()
	c.Assert(l, g.Not(g.IsNil))
	c.Assert(l.peers, g.Not(g.IsNil))
}

func (s *ListSuite) Test_Remove_doesNothingWhenAskedToRemoveEntryNotInList(c *g.C) {
	l := New()
	l.peers["foo@bar.com"] = &Peer{}

	res, rem := l.Remove(tj("bar@foo.com"))

	c.Assert(rem, g.Equals, false)
	c.Assert(res, g.IsNil)
	c.Assert(len(l.peers), g.Equals, 1)
}

func (s *ListSuite) Test_Remove_removesAnEntryIfInTheList(c *g.C) {
	l := New()
	l.peers["foo@bar.com"] = &Peer{}
	l.peers["bar@foo.com"] = &Peer{Name: "me"}

	res, rem := l.Remove(tj("bar@foo.com/somewhere"))

	c.Assert(rem, g.Equals, true)
	c.Assert(res.Name, g.Equals, "me")
	c.Assert(len(l.peers), g.Equals, 1)
}

func (s *ListSuite) Test_AddOrReplace_addsTheEntryIfNotInTheList(c *g.C) {
	l := New()
	p := &Peer{Jid: tj("somewhere"), Name: "something"}

	res := l.AddOrReplace(p)

	c.Assert(res, g.Equals, true)
	c.Assert(len(l.peers), g.Equals, 1)
	c.Assert(l.peers["somewhere"], g.Equals, p)
}

func (s *ListSuite) Test_AddOrReplace_replacesTheEntryIfInTheList(c *g.C) {
	l := New()
	p1 := &Peer{Jid: tj("somewhere"), Name: "something", Groups: toSet("hello"), Subscription: "from"}
	l.peers["somewhere"] = p1

	p2 := &Peer{Jid: tj("somewhere"), Name: "something2", Groups: toSet("goodbye")}
	res := l.AddOrReplace(p2)

	c.Assert(res, g.Equals, false)
	c.Assert(len(l.peers), g.Equals, 1)
	c.Assert(l.peers["somewhere"], g.Equals, p2)
}

func (s *ListSuite) Test_AddOrMerge_addsTheEntryIfNotInTheList(c *g.C) {
	l := New()
	p := &Peer{Jid: tj("somewhere"), Name: "something"}

	res := l.AddOrMerge(p)

	c.Assert(res, g.Equals, true)
	c.Assert(len(l.peers), g.Equals, 1)
	c.Assert(l.peers["somewhere"], g.Equals, p)
}

func (s *ListSuite) Test_AddOrReplace_mergesTheEntriesIfInTheList(c *g.C) {
	l := New()
	p1 := &Peer{Jid: tj("somewhere"), Name: "something", Groups: toSet("hello"), Subscription: "from"}
	l.peers["somewhere"] = p1

	p2 := &Peer{Jid: tj("somewhere"), Name: "something2", Groups: toSet("goodbye")}
	res := l.AddOrMerge(p2)

	c.Assert(res, g.Equals, false)
	c.Assert(len(l.peers), g.Equals, 1)
	c.Assert(l.peers["somewhere"], g.DeepEquals, &Peer{Jid: tj("somewhere"), Name: "something2", Groups: toSet("goodbye"), Subscription: "from", resources: make(map[string]Status)})
}

func (s *ListSuite) Test_ToSlice_createsASliceOfTheContentSortedAlphabetically(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: tj("foo@somewhere.com")})
	l.AddOrMerge(&Peer{Jid: tj("foo@somewhen.com")})
	l.AddOrMerge(&Peer{Jid: tj("bar@somewhere.com")})

	c.Assert(l.ToSlice(), g.DeepEquals, []*Peer{
		&Peer{Jid: tj("bar@somewhere.com")},
		&Peer{Jid: tj("foo@somewhen.com")},
		&Peer{Jid: tj("foo@somewhere.com")},
	})
}

func (s *ListSuite) Test_Iter_yieldsEachEntry(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: tj("foo@somewhere.com")})
	l.AddOrMerge(&Peer{Jid: tj("foo@somewhen.com")})
	l.AddOrMerge(&Peer{Jid: tj("bar@somewhere.com")})

	called := 0

	l.Iter(func(ix int, p *Peer) {
		called++

		switch ix {
		case 0:
			c.Assert(p, g.DeepEquals, &Peer{Jid: tj("bar@somewhere.com")})
		case 1:
			c.Assert(p, g.DeepEquals, &Peer{Jid: tj("foo@somewhen.com")})
		case 2:
			c.Assert(p, g.DeepEquals, &Peer{Jid: tj("foo@somewhere.com")})
		}
	})

	c.Assert(called, g.DeepEquals, 3)
}

func (s *ListSuite) Test_Unsubscribed_whenDoesntExist(c *g.C) {
	l := New()
	l.Unsubscribed(tj("foo@bar.com"))
}

func (s *ListSuite) Test_Unsubscribed_whenExist(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: tj("foo@bar.com"), Subscription: "both", Asked: true, PendingSubscribeID: "foo"})
	l.AddOrMerge(&Peer{Jid: tj("foo2@bar.com"), Subscription: "to"})
	l.AddOrMerge(&Peer{Jid: tj("foo3@bar.com"), Subscription: "from"})

	l.Unsubscribed(tj("foo@bar.com/123"))
	c.Assert(l.peers["foo@bar.com"].Subscription, g.Equals, "from")
	c.Assert(l.peers["foo@bar.com"].Asked, g.Equals, false)
	c.Assert(l.peers["foo@bar.com"].PendingSubscribeID, g.Equals, "")

	l.Unsubscribed(tj("foo2@bar.com/123"))
	c.Assert(l.peers["foo2@bar.com"].Subscription, g.Equals, "none")

	l.Unsubscribed(tj("foo3@bar.com/123"))
	c.Assert(l.peers["foo3@bar.com"].Subscription, g.Equals, "from")
}

func (s *ListSuite) Test_Subscribed_whenDoesntExist(c *g.C) {
	l := New()
	l.Subscribed(tj("foo@bar.com"))
	c.Assert(len(l.peers), g.Equals, 0)
}

func (s *ListSuite) Test_Subscribed_whenExist(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: tj("foo@bar.com"), Subscription: "from", Asked: true, PendingSubscribeID: "foo"})
	l.AddOrMerge(&Peer{Jid: tj("foo2@bar.com"), Subscription: "none"})
	l.AddOrMerge(&Peer{Jid: tj("foo3@bar.com"), Subscription: ""})
	l.AddOrMerge(&Peer{Jid: tj("foo4@bar.com"), Subscription: "both"})

	l.Subscribed(tj("foo@bar.com/123"))
	c.Assert(l.peers["foo@bar.com"].Subscription, g.Equals, "both")
	c.Assert(l.peers["foo@bar.com"].Asked, g.Equals, false)
	c.Assert(l.peers["foo@bar.com"].PendingSubscribeID, g.Equals, "")

	l.Subscribed(tj("foo2@bar.com/123"))
	c.Assert(l.peers["foo2@bar.com"].Subscription, g.Equals, "to")

	l.Subscribed(tj("foo3@bar.com/123"))
	c.Assert(l.peers["foo3@bar.com"].Subscription, g.Equals, "to")

	l.Subscribed(tj("foo4@bar.com/123"))
	c.Assert(l.peers["foo4@bar.com"].Subscription, g.Equals, "both")
}

func (s *ListSuite) Test_GetPendingSubscribe_returnsThePendingSubscribeIfExists(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: tj("foo@bar.com"), PendingSubscribeID: "foo"})
	l.AddOrMerge(&Peer{Jid: tj("foo2@bar.com")})

	v, k := l.GetPendingSubscribe(tj("none@foo.com"))
	c.Assert(k, g.Equals, false)
	c.Assert(v, g.Equals, "")

	v, k = l.GetPendingSubscribe(tj("foo@bar.com/bar"))
	c.Assert(k, g.Equals, true)
	c.Assert(v, g.Equals, "foo")

	v, k = l.GetPendingSubscribe(tj("foo2@bar.com/bar"))
	c.Assert(k, g.Equals, false)
	c.Assert(v, g.Equals, "")
}

func (s *ListSuite) Test_RemovePendingSubscribe_removesThePendingSubscribe(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: tj("foo@bar.com"), PendingSubscribeID: "foo"})
	l.AddOrMerge(&Peer{Jid: tj("foo2@bar.com")})

	v, k := l.RemovePendingSubscribe(tj("none@foo.com"))
	c.Assert(k, g.Equals, false)
	c.Assert(v, g.Equals, "")

	v, k = l.RemovePendingSubscribe(tj("foo@bar.com/bar"))
	c.Assert(k, g.Equals, true)
	c.Assert(v, g.Equals, "foo")
	c.Assert(l.peers["foo@bar.com"].PendingSubscribeID, g.Equals, "")

	v, k = l.RemovePendingSubscribe(tj("foo2@bar.com/bar"))
	c.Assert(k, g.Equals, false)
	c.Assert(v, g.Equals, "")
}

func (s *ListSuite) Test_SubscribeRequest_addsTheSubscribeID(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: tj("foo@bar.com")})

	l.SubscribeRequest(tj("fox@bar.com/hmm"), "something", "")
	c.Assert(l.peers["fox@bar.com"].PendingSubscribeID, g.Equals, "something")

	l.SubscribeRequest(tj("foo@bar.com/hmm2"), "something3", "")
	c.Assert(l.peers["foo@bar.com"].PendingSubscribeID, g.Equals, "something3")
}

func (s *ListSuite) Test_PeerBecameUnavailable_setsTheOfflineState(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: tj("foo@bar.com"), resources: map[string]Status{"foo2": Status{"a", "b"}}})

	res := l.PeerBecameUnavailable(jid.Parse("hmm@bar.com/foo"))
	c.Assert(res, g.Equals, false)

	res = l.PeerBecameUnavailable(jid.Parse("foo@bar.com/foo2"))
	c.Assert(res, g.Equals, true)
	c.Assert(l.peers["foo@bar.com"].IsOnline(), g.Equals, false)
}

func (s *ListSuite) Test_PeerBecameUnavailable_clearsResources(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: tj("foo@bar.com"), resources: map[string]Status{"foo2": Status{"a", "b"}}})

	// res := l.PeerBecameUnavailable(jid.Parse("hmm@bar.com"))
	// c.Assert(res, g.Equals, false)

	res := l.PeerBecameUnavailable(jid.Parse("foo@bar.com"))
	c.Assert(res, g.Equals, true)
	c.Assert(l.peers["foo@bar.com"].IsOnline(), g.Equals, false)
	c.Assert(l.peers["foo@bar.com"].resources, g.DeepEquals, map[string]Status{})

}

func (s *ListSuite) Test_PeerPresenceUpdate_sometimesUpdatesNonExistantPeers(c *g.C) {
	l := New()

	res := l.PeerPresenceUpdate(tjr("foo@bar.com/hmm"), "hello", "goodbye", "")
	c.Assert(res, g.Equals, true)
	c.Assert(l.peers["foo@bar.com"].MainStatus(), g.Equals, "hello")
	c.Assert(l.peers["foo@bar.com"].MainStatusMsg(), g.Equals, "goodbye")

	res = l.PeerPresenceUpdate(tjr("foo2@bar.com/hmm"), "xa", "goodbye", "")
	c.Assert(res, g.Equals, true)
	c.Assert(l.peers["foo2@bar.com"].MainStatus(), g.Equals, "xa")
	c.Assert(l.peers["foo2@bar.com"].MainStatusMsg(), g.Equals, "goodbye")

	res = l.PeerPresenceUpdate(tjr("foo3@bar.com/hmm"), "away", "goodbye", "")
	c.Assert(res, g.Equals, true)
	c.Assert(l.peers["foo3@bar.com"].MainStatus(), g.Equals, "away")
	c.Assert(l.peers["foo3@bar.com"].MainStatusMsg(), g.Equals, "goodbye")

}

func (s *ListSuite) Test_PeerPresenceUpdate_updatesPreviouslyKnownPeer(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: tj("foo@bar.com"), resources: make(map[string]Status)})
	l.AddOrMerge(&Peer{Jid: tj("foo2@bar.com"), resources: map[string]Status{"hmm": Status{"dnd", "working"}}})

	res := l.PeerPresenceUpdate(tjr("foo@bar.com/hmm"), "hello", "goodbye", "")
	c.Assert(res, g.Equals, true)
	c.Assert(l.peers["foo@bar.com"].MainStatus(), g.Equals, "hello")
	c.Assert(l.peers["foo@bar.com"].MainStatusMsg(), g.Equals, "goodbye")
	c.Assert(l.peers["foo@bar.com"].IsOnline(), g.Equals, true)

	res = l.PeerPresenceUpdate(tjr("foo2@bar.com/hmm"), "dnd", "working", "")
	c.Assert(res, g.Equals, false)
	c.Assert(l.peers["foo2@bar.com"].MainStatus(), g.Equals, "dnd")
	c.Assert(l.peers["foo2@bar.com"].MainStatusMsg(), g.Equals, "working")
	c.Assert(l.peers["foo2@bar.com"].IsOnline(), g.Equals, true)
}

func (s *ListSuite) Test_Clear_clearsTheList(c *g.C) {
	l := New()
	l.AddOrMerge(&Peer{Jid: tj("foo@bar.com")})

	l.Clear()

	c.Assert(len(l.peers), g.Equals, 0)
}

func (s *ListSuite) Test_Peers_sortsByNameForPresentation(c *g.C) {
	expectedPeers := []*Peer{
		&Peer{
			Jid: tj("ba"), Name: "ab",
		},
		&Peer{
			Jid: tj("ac"), Name: "",
		},
		&Peer{
			Jid: tj("aa"), Name: "bb",
		},
		&Peer{
			Jid: tj("aa"), Name: "cb",
		},
	}

	group := &Group{
		peers: []*Peer{
			expectedPeers[2],
			expectedPeers[0],
			expectedPeers[3],
			expectedPeers[1],
		},
	}

	c.Assert(group.Peers(), g.DeepEquals, expectedPeers)
}

func (s *ListSuite) Test_LatestError_setsLatestErrorWhenExists(c *g.C) {
	l := New()
	pp := &Peer{Jid: tj("foo@bar.com")}
	l.AddOrMerge(pp)
	l.LatestError(tj("foo@bar.com/foo"), "tow", "frou", "sxi")

	c.Assert(pp.LatestError, g.DeepEquals, &PeerError{"tow", "frou", "sxi"})
}

func (s *ListSuite) Test_LatestError_doesntDoAnythingForUnexistingPeer(c *g.C) {
	l := New()
	l.LatestError(tj("foo@bar.com/foo"), "tow", "frou", "sxi")
}

func (s *ListSuite) Test_IterAll_willIterateOverAllTheListsGivenAndYieldTheirPeers(c *g.C) {
	l := New()
	l2 := New()
	pp := &Peer{Jid: tj("foo@bar.com")}
	pp2 := &Peer{Jid: tj("foo2@bar.com")}
	l.AddOrMerge(pp)
	l2.AddOrMerge(pp2)

	result := []*Peer{}
	IterAll(func(_ int, p *Peer) {
		result = append(result, p)
	}, l, l2)

	c.Assert(result, g.DeepEquals, []*Peer{pp2, pp})
}

func (s *ListSuite) Test_List_GetGroupNames(c *g.C) {
	l := New()
	c.Assert(l.GetGroupNames(), g.DeepEquals, map[string]bool{})
	l.AddOrMerge(&Peer{Jid: tj("foo@bar.com"), Groups: map[string]bool{
		"one":      true,
		"two::foo": true,
	}})

	l.AddOrMerge(&Peer{Jid: tj("bar@bar.com"), Groups: map[string]bool{
		"four": true,
	}})

	l.AddOrMerge(&Peer{Jid: tj("some@bar.com"), Groups: map[string]bool{}})
	c.Assert(l.GetGroupNames(), g.DeepEquals, map[string]bool{
		"one":      true,
		"four":     true,
		"two::foo": true,
	})
}
