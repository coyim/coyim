package xmpp

import (
	"errors"
	"time"

	"github.com/coyim/coyim/xmpp/data"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"
)

type PingSuite struct{}

var _ = Suite(&PingSuite{})

func (s *PingSuite) Test_conn_SendPing(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		log: testLogger(),
		out: mockOut,
		jid: "juliet@example.com/chamber",
	}
	conn.inflights = make(map[data.Cookie]inflight)

	reply, cookie, err := conn.SendPing()
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Matches, "<iq xmlns='jabber:client' from='juliet@example.com/chamber' type='get' id='.+'><ping xmlns=\"urn:xmpp:ping\"></ping></iq>")
	c.Assert(reply, NotNil)
	c.Assert(cookie, NotNil)
}

func (s *PingSuite) Test_conn_SendPingReply(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		log: testLogger(),
		out: mockOut,
		jid: "juliet@example.com/chamber",
	}
	conn.inflights = make(map[data.Cookie]inflight)

	err := conn.SendPingReply("huff")
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Matches, "<iq xmlns='jabber:client' from='juliet@example.com/chamber' type='result' id='huff'></iq>")
}

func (s *PingSuite) Test_conn_ReceivePong(c *C) {
	t1 := time.Now().Add(-time.Hour * 1)
	conn := conn{
		lastPongResponse: t1,
	}

	conn.ReceivePong()
	c.Assert(t1, Not(Equals), conn.lastPongResponse)
}

func (s *PingSuite) Test_ParsePong_failsIfStanzaIsNotIQ(c *C) {
	rp := data.Stanza{
		Value: "hello",
	}

	e := ParsePong(rp)
	c.Assert(e, ErrorMatches, "xmpp: ping request resulted in tag of type.*")
}

func (s *PingSuite) Test_ParsePong_failsIfIQTypeIsError(c *C) {
	rp := data.Stanza{
		Value: &data.ClientIQ{
			Type: "error",
			Error: data.StanzaError{
				Text: "oh no",
			},
		},
	}

	e := ParsePong(rp)
	c.Assert(e, ErrorMatches, "xmpp: ping request resulted in an error: oh no")
}

func (s *PingSuite) Test_ParsePong_failsIfIQTypeIsUnknown(c *C) {
	rp := data.Stanza{
		Value: &data.ClientIQ{
			Type: "bla",
		},
	}

	e := ParsePong(rp)
	c.Assert(e, ErrorMatches, "xmpp: ping request resulted in an unexpected type")
}

func (s *PingSuite) Test_ParsePong_worksWithCorrectType(c *C) {
	rp := data.Stanza{
		Value: &data.ClientIQ{
			Type: "result",
		},
	}

	e := ParsePong(rp)
	c.Assert(e, IsNil)
}

func (s *PingSuite) Test_conn_watchPings_returnsWhenClosed(c *C) {
	orgPingInterval := pingInterval
	defer func() {
		pingInterval = orgPingInterval
	}()

	pingInterval = 1 * time.Millisecond

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	cn := &conn{
		closed: true,
		log:    l,
	}

	cn.watchPings()
	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "xmpp: trying to send ping on closed connection")
}

func (s *PingSuite) Test_conn_watchPings_returnsWhenFailingToSendPing(c *C) {
	orgPingInterval := pingInterval
	defer func() {
		pingInterval = orgPingInterval
	}()

	pingInterval = 1 * time.Millisecond

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockOut := &mockConnIOReaderWriter{
		err: errors.New("marker error"),
	}

	cn := &conn{
		closed: false,
		log:    l,
		out:    mockOut,
	}

	cn.watchPings()
	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "xmpp: error when sending ping")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["error"], ErrorMatches, "marker error")
}

func (s *PingSuite) Test_conn_watchPings_succeedsOneRound(c *C) {
	orgPingInterval := pingInterval
	defer func() {
		pingInterval = orgPingInterval
	}()

	pingInterval = 1 * time.Millisecond

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockOut := &mockConnIOReaderWriter{}

	cn := &conn{
		closed:    false,
		log:       l,
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan bool, 1)
	go func() {
		cn.watchPings()
		done <- true
	}()

	inf := waitForInflightTo(cn, "")
	c.Assert(inf, Not(IsNil))

	cn.closed = true

	inf.replyChan <- data.Stanza{
		Value: &data.ClientIQ{
			Type: "result",
		},
	}

	<-done

	c.Assert(hook.Entries, HasLen, 1)
}

func (s *PingSuite) Test_conn_watchPings_nonIQResult(c *C) {
	orgPingInterval := pingInterval
	defer func() {
		pingInterval = orgPingInterval
	}()

	pingInterval = 1 * time.Millisecond

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockOut := &mockConnIOReaderWriter{}

	cn := &conn{
		closed:    false,
		log:       l,
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan bool, 1)
	go func() {
		cn.watchPings()
		done <- true
	}()

	inf := waitForInflightTo(cn, "")
	c.Assert(inf, Not(IsNil))

	cn.closed = true

	inf.replyChan <- data.Stanza{
		Value: "something",
	}

	<-done

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "xmpp: received invalid result to ping")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["value"], Equals, "something")
}

func (s *PingSuite) Test_conn_watchPings_IQerrorResult(c *C) {
	orgPingInterval := pingInterval
	defer func() {
		pingInterval = orgPingInterval
	}()

	pingInterval = 1 * time.Millisecond

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockOut := &mockConnIOReaderWriter{}

	cn := &conn{
		closed:    false,
		log:       l,
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan bool, 1)
	go func() {
		cn.watchPings()
		done <- true
	}()

	inf := waitForInflightTo(cn, "")
	c.Assert(inf, Not(IsNil))

	cn.closed = true

	inf.replyChan <- data.Stanza{
		Value: &data.ClientIQ{
			Type: "error",
		},
	}

	<-done

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "xmpp: received invalid result to ping")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
}

func (s *PingSuite) Test_conn_watchPings_closedResultChannel(c *C) {
	orgPingInterval := pingInterval
	defer func() {
		pingInterval = orgPingInterval
	}()

	pingInterval = 1 * time.Millisecond

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockOut := &mockConnIOReaderWriter{}

	cn := &conn{
		closed:    false,
		log:       l,
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan bool, 1)
	go func() {
		cn.watchPings()
		done <- true
	}()

	inf := waitForInflightTo(cn, "")
	c.Assert(inf, Not(IsNil))

	cn.closed = true

	close(inf.replyChan)

	<-done

	c.Assert(hook.Entries, HasLen, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "xmpp: ping result channel closed")
	c.Assert(hook.Entries[0].Data, HasLen, 0)
}

func (s *PingSuite) Test_conn_watchPings_timesOutAndFails(c *C) {
	orgPingInterval := pingInterval
	orgPingTimeout := pingTimeout
	orgStreamClosedTimeout := streamClosedTimeout
	defer func() {
		pingInterval = orgPingInterval
		pingTimeout = orgPingTimeout
		streamClosedTimeout = orgStreamClosedTimeout
	}()

	pingInterval = 1 * time.Millisecond
	pingTimeout = 1 * time.Millisecond
	streamClosedTimeout = 1 * time.Millisecond

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockOut := &mockConnIOReaderWriter{}

	cn := &conn{
		closed:    false,
		log:       l,
		out:       mockOut,
		rawOut:    mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	cn.watchPings()

	c.Assert(hook.Entries, HasLen, 4)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "xmpp: ping failures reached threshold")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["threshold"], Equals, 2)
}
