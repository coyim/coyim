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
	l.peers["bar@foo.com"] = &Peer{name: "me"}

	res, rem := l.Remove("bar@foo.com/somewhere")

	c.Assert(rem, g.Equals, true)
	c.Assert(res.name, g.Equals, "me")
	c.Assert(len(l.peers), g.Equals, 1)
}

func (s *ListXmppSuite) Test_AddOrReplace_addsTheEntryIfNotInTheList(c *g.C) {
	l := New()
	p := &Peer{jid: "somewhere", name: "something"}

	res := l.AddOrReplace(p)

	c.Assert(res, g.Equals, true)
	c.Assert(len(l.peers), g.Equals, 1)
	c.Assert(l.peers["somewhere"], g.Equals, p)
}

func (s *ListXmppSuite) Test_AddOrReplace_replacesTheEntryIfInTheList(c *g.C) {
	l := New()
	p1 := &Peer{jid: "somewhere", name: "something", groups: toSet("hello"), subscription: "from"}
	l.peers["somewhere"] = p1

	p2 := &Peer{jid: "somewhere", name: "something2", groups: toSet("goodbye")}
	res := l.AddOrReplace(p2)

	c.Assert(res, g.Equals, false)
	c.Assert(len(l.peers), g.Equals, 1)
	c.Assert(l.peers["somewhere"], g.Equals, p2)
}

func (s *ListXmppSuite) Test_AddOrMerge_addsTheEntryIfNotInTheList(c *g.C) {
	l := New()
	p := &Peer{jid: "somewhere", name: "something"}

	res := l.AddOrMerge(p)

	c.Assert(res, g.Equals, true)
	c.Assert(len(l.peers), g.Equals, 1)
	c.Assert(l.peers["somewhere"], g.Equals, p)
}

func (s *ListXmppSuite) Test_AddOrReplace_mergesTheEntriesIfInTheList(c *g.C) {
	l := New()
	p1 := &Peer{jid: "somewhere", name: "something", groups: toSet("hello"), subscription: "from"}
	l.peers["somewhere"] = p1

	p2 := &Peer{jid: "somewhere", name: "something2", groups: toSet("goodbye")}
	res := l.AddOrMerge(p2)

	c.Assert(res, g.Equals, false)
	c.Assert(len(l.peers), g.Equals, 1)
	c.Assert(*l.peers["somewhere"], g.DeepEquals, Peer{jid: "somewhere", name: "something2", groups: toSet("hello", "goodbye"), subscription: "from"})
}
