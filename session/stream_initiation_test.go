package session

import (
	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/data"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"
)

type StreamInitiationSuite struct{}

var _ = Suite(&StreamInitiationSuite{})

func (s *StreamInitiationSuite) Test_streamInitIQ_works(c *C) {
	called := false
	supportedSIProfiles["temporary-bla-streaminit-test-profile"] = func(access.Session, *data.ClientIQ, data.SI) (interface{}, string, bool) {
		called = true
		return nil, "", false
	}

	defer func() {
		delete(supportedSIProfiles, "temporary-bla-streaminit-test-profile")
	}()

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{log: l}

	iq := &data.ClientIQ{
		Query: []byte(`<si xmlns="http://jabber.org/protocol/si" profile="temporary-bla-streaminit-test-profile"/>`),
	}

	_, _, _ = streamInitIQ(sess, iq)

	c.Assert(called, Equals, true)

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.InfoLevel)
	c.Assert(hook.LastEntry().Message, Equals, "IQ: stream initiation")
}

func (s *StreamInitiationSuite) Test_streamInitIQ_failsOnBadXML(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{log: l}

	iq := &data.ClientIQ{
		Query: []byte(`<si `),
	}

	_, _, ignore := streamInitIQ(sess, iq)

	c.Assert(ignore, Equals, false)

	c.Assert(len(hook.Entries), Equals, 2)
	c.Assert(hook.LastEntry().Level, Equals, log.WarnLevel)
	c.Assert(hook.LastEntry().Message, Equals, "Failed to parse stream initiation")
	c.Assert(hook.LastEntry().Data["error"], ErrorMatches, "XML syntax error on line.*")
}

func (s *StreamInitiationSuite) Test_streamInitIQ_failsOnUnknownProfile(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{log: l}

	iq := &data.ClientIQ{
		Query: []byte(`<si xmlns="http://jabber.org/protocol/si" profile="unknown-test-profile"/>`),
	}

	_, _, _ = streamInitIQ(sess, iq)

	c.Assert(len(hook.Entries), Equals, 2)
	c.Assert(hook.LastEntry().Level, Equals, log.WarnLevel)
	c.Assert(hook.LastEntry().Message, Equals, "Unsupported SI profile")
	c.Assert(hook.LastEntry().Data["profile"], Equals, "unknown-test-profile")
}
