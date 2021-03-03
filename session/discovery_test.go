package session

import (
	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/xmpp/data"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"
)

type DiscoverySuite struct{}

var _ = Suite(&DiscoverySuite{})

func (s *DiscoverySuite) Test_discoIQ_works(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	conf := &config.Account{Account: "foo.bar@hello.com"}
	sess := &session{accountConfig: conf, log: l}

	iq := &data.ClientIQ{
		Query: []byte(`
<query xmlns="http://jabber.org/protocol/disco#info">
  <node/>
</query>
		`),
	}

	ret, tp, ignore := discoIQ(sess, iq)

	c.Assert(ret, FitsTypeOf, data.DiscoveryInfoQuery{})
	c.Assert(tp, Equals, "")
	c.Assert(ignore, Equals, false)

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.InfoLevel)
	c.Assert(hook.LastEntry().Message, Equals, "IQ: http://jabber.org/protocol/disco#info query")
}

func (s *DiscoverySuite) Test_discoIQ_failsOnBadXML(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	conf := &config.Account{Account: "foo.bar@hello.com"}
	sess := &session{accountConfig: conf, log: l}

	iq := &data.ClientIQ{
		Query: []byte(`<query xmlns="http://jabber.org/protocol/disco#info"`),
	}

	ret, tp, ignore := discoIQ(sess, iq)

	c.Assert(ret, IsNil)
	c.Assert(tp, Equals, "error")
	c.Assert(ignore, Equals, false)

	c.Assert(len(hook.Entries), Equals, 2)
	c.Assert(hook.LastEntry().Level, Equals, log.ErrorLevel)
	c.Assert(hook.LastEntry().Message, Equals, "error on parsing disco#info query")
	c.Assert(hook.LastEntry().Data["error"], ErrorMatches, "XML syntax error on line.*")
}

func (s *DiscoverySuite) Test_discoItemsIQ_failsOnBadXML(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	conf := &config.Account{Account: "foo.bar@hello.com"}
	sess := &session{accountConfig: conf, log: l}

	iq := &data.ClientIQ{
		Query: []byte(`<query xmlns="http://jabber.org/protocol/disco#info"`),
	}

	ret, tp, ignore := discoItemsIQ(sess, iq)

	c.Assert(ret, IsNil)
	c.Assert(tp, Equals, "error")
	c.Assert(ignore, Equals, false)

	c.Assert(len(hook.Entries), Equals, 2)
	c.Assert(hook.LastEntry().Level, Equals, log.ErrorLevel)
	c.Assert(hook.LastEntry().Message, Equals, "error on parsing disco#items query")
	c.Assert(hook.LastEntry().Data["error"], ErrorMatches, "XML syntax error on line.*")
}

func (s *DiscoverySuite) Test_discoItemsIQ_worksForMUCRooms(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	conf := &config.Account{Account: "foo.bar@hello.com"}
	sess := &session{accountConfig: conf, log: l}

	iq := &data.ClientIQ{
		Query: []byte(`
<query xmlns="http://jabber.org/protocol/disco#items">
  <node>http://jabber.org/protocol/muc#rooms</node>
</query>
		`),
	}

	ret, tp, ignore := discoItemsIQ(sess, iq)

	c.Assert(ret, FitsTypeOf, data.DiscoveryItemsQuery{})
	c.Assert(ret.(data.DiscoveryItemsQuery).DiscoveryItems, HasLen, 0)
	c.Assert(tp, Equals, "")
	c.Assert(ignore, Equals, false)

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.InfoLevel)
	c.Assert(hook.LastEntry().Message, Equals, "IQ: http://jabber.org/protocol/disco#items query")
}

func (s *DiscoverySuite) Test_discoItemsIQ_returnsErrorForOtherNode(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	conf := &config.Account{Account: "foo.bar@hello.com"}
	sess := &session{accountConfig: conf, log: l}

	iq := &data.ClientIQ{
		Query: []byte(`
<query xmlns="http://jabber.org/protocol/disco#items">
  <node/>
</query>
		`),
	}

	ret, tp, ignore := discoItemsIQ(sess, iq)

	c.Assert(ret, FitsTypeOf, data.ErrorReply{})
	c.Assert(ret.(data.ErrorReply).Type, Equals, "cancel")
	c.Assert(ret.(data.ErrorReply).Error, DeepEquals, data.ErrorServiceUnavailable{})
	c.Assert(tp, Equals, "")
	c.Assert(ignore, Equals, false)

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.InfoLevel)
	c.Assert(hook.LastEntry().Message, Equals, "IQ: http://jabber.org/protocol/disco#items query")
}
