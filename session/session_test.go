package session

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/otrclient"
	"github.com/coyim/coyim/roster"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/otr3"

	. "gopkg.in/check.v1"
)

type SessionSuite struct{}

var _ = Suite(&SessionSuite{})

func (s *SessionSuite) Test_NewSession_returnsANewSession(c *C) {
	sess := Factory(&config.ApplicationConfig{}, &config.Account{}, xmpp.DialerFactory)
	c.Assert(sess, Not(IsNil))
}

const testTimeout = time.Duration(5) * time.Second

func (s *SessionSuite) Test_iqReceived_publishesIQReceivedEvent(c *C) {
	sess := &session{
		log: log.StandardLogger(),
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 1)
	sess.eventsReachedZero = eventsDone

	sess.iqReceived(jid.NR("someone@somewhere"))

	select {
	case <-eventsDone:
		close(observer)
		select {
		case ev := <-observer:
			c.Assert(ev, Equals, events.Peer{
				Type: events.IQReceived,
				From: jid.NR("someone@somewhere"),
			})
		default:
			c.Errorf("did not receive event")
		}
	case <-time.After(testTimeout):
		c.Errorf("test timed out")
	}
}

func (s *SessionSuite) Test_WatchStanzas_warnsAndExitsOnBadStanza(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<clientx:message xmlns:client='jabber:client' to='fo@bar.com' from='bar@foo.com' type='chat'><client:body>something</client:body></client:message>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	l, hook := test.NewNullLogger()
	sess := &session{
		log:        l,
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 1)
	sess.eventsReachedZero = eventsDone
	done := make(chan bool)
	sess.doneBadStanza = done

	sess.watchStanzas()

	<-done

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.ErrorLevel)
	c.Assert(hook.LastEntry().Message, Equals, "error reading XMPP message")
	c.Assert(hook.LastEntry().Data["error"].(error).Error(), Equals, "unexpected XMPP message clientx <message/>")
}

func (s *SessionSuite) Test_WatchStanzas_handlesUnknownMessage(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<bind:bind xmlns:bind='urn:ietf:params:xml:ns:xmpp-bind'></bind:bind>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	l, hook := test.NewNullLogger()

	sess := &session{
		log:        l,
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.watchStanzas()

	e := checkLogHasAny(hook, log.InfoLevel, "unhandled stanza")
	c.Assert(e, Not(IsNil))
	c.Assert(e.Data["name"].(xml.Name), Equals, xml.Name{Space: "urn:ietf:params:xml:ns:xmpp-bind", Local: "bind"})
	c.Assert(*(e.Data["value"].(*data.BindBind)), Equals, data.BindBind{XMLName: xml.Name{Space: "urn:ietf:params:xml:ns:xmpp-bind", Local: "bind"}, Resource: "", Jid: ""})
}

type checker func(interface{}) bool

func assertReceivesEvent(c *C, eventsDone chan bool, observer <-chan interface{}, exp checker) {
	select {
	case <-eventsDone:
		for {
			select {
			case ev := <-observer:
				if exp(ev) {
					return
				}
			case <-eventsDone:
			case <-time.After(testTimeout):
				c.Errorf("test timed out")
				return
			}
		}
	case <-time.After(testTimeout):
		c.Errorf("test timed out")
	}
}

func (s *SessionSuite) Test_WatchStanzas_handlesStreamError_withText(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:error xmlns:stream='http://etherx.jabber.org/streams'><stream:text>bad horse showed up</stream:text></stream:error>")}
	l, hook := test.NewNullLogger()
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		log:        l,
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.watchStanzas()

	e := checkLogHasAny(hook, log.ErrorLevel, "Exiting in response to fatal error from server")
	c.Assert(e, Not(IsNil))
	c.Assert(e.Data["stanza"], Equals, "bad horse showed up")
}

func checkLogHasAny(hook *test.Hook, level log.Level, message string) *log.Entry {
	for _, e := range hook.Entries {
		if e.Level == level && e.Message == message {
			return &e
		}
	}

	return nil
}

func (s *SessionSuite) Test_WatchStanzas_handlesStreamError_withEmbeddedTag(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:error xmlns:stream='http://etherx.jabber.org/streams'><not-well-formed xmlns='urn:ietf:params:xml:ns:xmpp-streams'/></stream:error>")}
	l, hook := test.NewNullLogger()
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		log:        l,
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)

	sess.watchStanzas()

	e := checkLogHasAny(hook, log.ErrorLevel, "Exiting in response to fatal error from server")
	c.Assert(e, Not(IsNil))
	c.Assert(e.Data["stanza"], Equals, "{urn:ietf:params:xml:ns:xmpp-streams not-well-formed}")
}

func (s *SessionSuite) Test_WatchStanzas_receivesAMessage(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:message xmlns:client='jabber:client' type='chat' to='some@one.org/foo' from='bla@hmm.org/somewhere'><client:body>well, hello there</client:body></client:message>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := Factory(
		&config.ApplicationConfig{},
		&config.Account{InstanceTag: uint32(42)},
		xmpp.DialerFactory,
	).(*session)

	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.watchStanzas()

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Message)
		if !ok {
			return false
		}
		c.Assert(t.Encrypted, Equals, false)
		c.Assert(t.From, Equals, jid.R("bla@hmm.org/somewhere"))
		c.Assert(string(t.Body), Equals, "well, hello there")
		return true
	})
}

func (s *SessionSuite) Test_WatchStanzas_failsOnUnrecognizedIQ(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='something'></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	l, hook := test.NewNullLogger()
	sess := &session{
		log:        l,
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.watchStanzas()

	e := checkLogHasAny(hook, log.InfoLevel, "unrecognized iq")
	c.Assert(e, Not(IsNil))
	c.Assert(e.Data["stanza"], DeepEquals, &data.ClientIQ{XMLName: xml.Name{Space: "jabber:client", Local: "iq"}, Type: "something", Query: []uint8{}})
}

func (s *SessionSuite) Test_WatchStanzas_getsDiscoInfoIQ(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte(
		"<client:iq xmlns:client='jabber:client' type='get' from='abc' to='cde'>" +
			"<query xmlns='http://jabber.org/protocol/disco#info'/>" +
			"</client:iq>",
	)}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		log:    log.StandardLogger(),
		config: &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "foo.bar@somewhere.org",
		},
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	stanzaChan := make(chan data.Stanza, 1000)
	stanza, _ := conn.Next()
	stanzaChan <- stanza

	sess.receiveStanza(stanzaChan)

	c.Assert(string(mockIn.write), Equals, ""+
		"<iq xmlns='jabber:client' to='abc' from='some@one.org/foo' type='result' id=''>"+
		"<query xmlns=\"http://jabber.org/protocol/disco#info\">"+
		"<identity xmlns=\"http://jabber.org/protocol/disco#info\" category=\"client\" type=\"pc\" name=\"foo.bar@somewhere.org\"></identity>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/disco#info\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/disco#items\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"urn:xmpp:bob\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"urn:xmpp:ping\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/caps\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"jabber:iq:version\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"vcard-temp\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"jabber:x:data\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/si\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/si/profile/file-transfer\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/si/profile/directory-transfer\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/si/profile/encrypted-data-transfer\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/bytestreams\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"urn:xmpp:eme:0\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/muc\"></feature>"+
		"</query>"+
		"</iq>")
}

func (s *SessionSuite) Test_WatchStanzas_getsVersionInfoIQ(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='get' from='abc' to='cde'><query xmlns='jabber:iq:version'/></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		log:    log.StandardLogger(),
		config: &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "foo.bar@somewhere.org",
		},
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	stanzaChan := make(chan data.Stanza, 1000)
	stanza, _ := conn.Next()
	stanzaChan <- stanza

	sess.receiveStanza(stanzaChan)

	c.Assert(string(mockIn.write), Equals, ""+
		"<iq xmlns='jabber:client' to='abc' from='some@one.org/foo' type='result' id=''>"+
		"<query xmlns=\"jabber:iq:version\">"+
		"<name>testing</name>"+
		"<version>version</version>"+
		"<os>none</os>"+
		"</query>"+
		"</iq>")
}

func (s *SessionSuite) Test_WatchStanzas_getsUnknown(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='get' from='abc' to='cde'><query xmlns='jabber:iq:somethingStrange'/></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)
	l, hook := test.NewNullLogger()

	sess := &session{
		log:    l,
		config: &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "foo.bar@somewhere.org",
		},
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)

	sess.watchStanzas()

	e := checkLogHasAny(hook, log.InfoLevel, "Unknown IQ: <query xmlns='jabber:iq:somethingStrange'/>")
	c.Assert(e, Not(IsNil))
}

func (s *SessionSuite) Test_WatchStanzas_iq_set_roster_withBadFrom(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='set' from='some2@one.org' to='cde'><query xmlns='jabber:iq:roster'/></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	l, hook := test.NewNullLogger()
	sess := &session{
		log:    l,
		config: &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)

	stanzaChan := make(chan data.Stanza, 1000)
	stanza, _ := conn.Next()
	stanzaChan <- stanza

	sess.receiveStanza(stanzaChan)

	c.Assert(len(hook.Entries), Equals, 2)
	c.Assert(hook.LastEntry().Level, Equals, log.WarnLevel)
	c.Assert(hook.LastEntry().Message, Equals, "Ignoring roster IQ from bad address")
	c.Assert(hook.LastEntry().Data["from"], Equals, "some2@one.org")

	c.Assert(string(mockIn.write), Equals, "")
}

func (s *SessionSuite) Test_WatchStanzas_iq_set_roster_withFromContainingJid(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='set' from='some@one.org/foo' to='cde'><query xmlns='jabber:iq:roster'/></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	l, hook := test.NewNullLogger()
	sess := &session{
		log:    l,
		config: &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)

	sess.watchStanzas()

	e := checkLogHasAny(hook, log.WarnLevel, "Failed to parse roster push IQ")
	c.Assert(e, Not(IsNil))
}

func (s *SessionSuite) Test_WatchStanzas_iq_set_roster_addsANewRosterItem(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='set' to='cde'><query xmlns='jabber:iq:roster'>" +
		"<item jid='romeo@example.net' name='Romeo' subscription='both'>" +
		"<group>Friends</group>" +
		"</item>" +
		"</query></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		log:    log.StandardLogger(),
		config: &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		r:          roster.New(),
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	sess.watchStanzas()

	c.Assert(sess.r.ToSlice(), DeepEquals, []*roster.Peer{
		peerFrom(data.RosterEntry{Jid: "romeo@example.net", Subscription: "both", Name: "Romeo", Group: []string{"Friends"}}, sess.GetConfig())})
}

func (s *SessionSuite) Test_WatchStanzas_iq_set_roster_setsExistingRosterItem(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='set' to='cde'><query xmlns='jabber:iq:roster'>" +
		"<item jid='romeo@example.net' name='Romeo' subscription='both'>" +
		"<group>Friends</group>" +
		"</item>" +
		"</query></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	called := 0

	sess := &session{
		log:    log.StandardLogger(),
		config: &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		r:          roster.New(),
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	sess.r.AddOrReplace(peerFrom(data.RosterEntry{Jid: "jill@example.net", Subscription: "both", Name: "Jill", Group: []string{"Foes"}}, sess.GetConfig()))
	sess.r.AddOrReplace(peerFrom(data.RosterEntry{Jid: "romeo@example.net", Subscription: "both", Name: "Mo", Group: []string{"Foes"}}, sess.GetConfig()))

	sess.watchStanzas()

	c.Assert(called, Equals, 0)
	c.Assert(sess.r.ToSlice(), DeepEquals, []*roster.Peer{
		peerFrom(data.RosterEntry{Jid: "jill@example.net", Subscription: "both", Name: "Jill", Group: []string{"Foes"}}, sess.GetConfig()),
		peerFrom(data.RosterEntry{Jid: "romeo@example.net", Subscription: "both", Name: "Romeo", Group: []string{"Friends"}}, sess.GetConfig()),
	})
}

func (s *SessionSuite) Test_WatchStanzas_iq_set_roster_removesRosterItems(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='set' to='cde'><query xmlns='jabber:iq:roster'>" +
		"<item jid='romeo@example.net' name='Romeo' subscription='remove'>" +
		"<group>Friends</group>" +
		"</item>" +
		"</query></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		log:    log.StandardLogger(),
		config: &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		r:          roster.New(),
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	sess.r.AddOrReplace(peerFrom(data.RosterEntry{Jid: "romeo@example.net", Subscription: "both", Name: "Mo", Group: []string{"Foes"}}, sess.GetConfig()))
	sess.r.AddOrReplace(peerFrom(data.RosterEntry{Jid: "jill@example.net", Subscription: "both", Name: "Jill", Group: []string{"Foes"}}, sess.GetConfig()))
	sess.r.AddOrReplace(peerFrom(data.RosterEntry{Jid: "romeo@example.net", Subscription: "both", Name: "Mo", Group: []string{"Foes"}}, sess.GetConfig()))

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)

	sess.watchStanzas()

	c.Assert(sess.r.ToSlice(), DeepEquals, []*roster.Peer{
		peerFrom(data.RosterEntry{Jid: "jill@example.net", Subscription: "both", Name: "Jill", Group: []string{"Foes"}}, sess.GetConfig()),
	})

	for {
		select {
		case ev := <-observer:
			switch ev.(type) {
			case events.Peer:
				c.Error("Received peer event")
				return
			default:
				// ignore
				continue
			}
		default:
			// Test succeded if we get here and no events happened
			return
		}
	}
}

func (s *SessionSuite) Test_WatchStanzas_presence_unavailable_forNoneKnownUser(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org' type='unavailable'><client:status>going on vacation</client:status></client:presence>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		log:        log.StandardLogger(),
		r:          roster.New(),
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)

	sess.watchStanzas()

	for {
		select {
		case ev := <-observer:
			switch ev.(type) {
			case events.Presence:
				c.Error("Received presence event")
				return
			default:
				// ignore
				continue
			}
		default:
			// Test succeded if we get here and no events happened
			return
		}
	}
}

func (s *SessionSuite) Test_WatchStanzas_presence_unavailable_forKnownUser(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org' type='unavailable'><client:status>going on vacation</client:status></client:presence>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		log:           log.StandardLogger(),
		config:        &config.ApplicationConfig{},
		accountConfig: &config.Account{},
		r:             roster.New(),
		connStatus:    DISCONNECTED,
	}
	sess.conn = conn
	sess.r.AddOrReplace(roster.PeerWithState(jid.NR("some2@one.org"), "somewhere", "", "", jid.NewResource("balcony")))

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.watchStanzas()

	p, _ := sess.r.Get(jid.NR("some2@one.org"))
	c.Assert(p.IsOnline(), Equals, false)

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Presence)
		if !ok {
			return false
		}

		c.Assert(t.Gone, Equals, true)
		return true
	})
}

func (s *SessionSuite) Test_WatchStanzas_presence_subscribe(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org' type='subscribe' id='adf12112'/>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		log:           log.StandardLogger(),
		config:        &config.ApplicationConfig{},
		accountConfig: &config.Account{},
		r:             roster.New(),
		connStatus:    DISCONNECTED,
	}
	sess.conn = conn

	sess.watchStanzas()

	v, _ := sess.r.GetPendingSubscribe(jid.NR("some2@one.org"))
	c.Assert(v, Equals, "adf12112")
}

func (s *SessionSuite) Test_WatchStanzas_presence_unknown(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org' type='weird'/>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		log:           log.StandardLogger(),
		config:        &config.ApplicationConfig{},
		accountConfig: &config.Account{},
		connStatus:    DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)

	sess.watchStanzas()

	for {
		select {
		case ev := <-observer:
			switch t := ev.(type) {
			case events.Presence:
				c.Error("Received presence event")
				return
			case events.Peer:
				if t.Type == events.SubscriptionRequest {
					c.Error("Received subscription request event")
				}
				return
			default:
				// ignore
				continue
			}
		default:
			// Test succeded if we get here and no events happened
			return
		}
	}
}

func (s *SessionSuite) Test_WatchStanzas_presence_regularPresenceIsAdded(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org'><client:show>dnd</client:show></client:presence>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		log:           log.StandardLogger(),
		config:        &config.ApplicationConfig{},
		accountConfig: &config.Account{},
		r:             roster.New(),
		connStatus:    DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.watchStanzas()

	pp, _ := sess.r.Get(jid.NR("some2@one.org"))
	st := pp.MainStatus()
	c.Assert(st, Equals, "dnd")

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Presence)
		if !ok {
			return false
		}

		c.Assert(t.Gone, Equals, false)
		return true
	})
}

func (s *SessionSuite) Test_WatchStanzas_presence_ignoresSameState(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org'><client:show>dnd</client:show></client:presence>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		log:           log.StandardLogger(),
		config:        &config.ApplicationConfig{},
		accountConfig: &config.Account{},
		r:             roster.New(),
		connStatus:    DISCONNECTED,
	}
	sess.conn = conn
	sess.r.AddOrReplace(roster.PeerWithState(jid.NR("some2@one.org"), "dnd", "", "", jid.NewResource("main")))

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)

	sess.watchStanzas()

	pp, _ := sess.r.Get(jid.NR("some2@one.org"))
	st := pp.MainStatus()
	c.Assert(st, Equals, "dnd")

	// In this loop we will drain all events from the observer.
	// If we ever get a presence event, we will fail the test
	// However, if the observer channel is empty, we know that
	// no presence events would be sent - since above we already
	// checked that the update has happened. We don't need
	// to do a timeout or anything like that.
	for {
		select {
		case ev := <-observer:
			switch ev.(type) {
			case events.Presence:
				c.Error("Received presence event")
				return
			default:
				// ignore
				continue
			}
		default:
			// Test succeded if we get here and no events happened
			return
		}
	}
}

func (s *SessionSuite) Test_HandleConfirmOrDeny_failsWhenNoPendingSubscribeIsWaiting(c *C) {
	l, hook := test.NewNullLogger()
	sess := &session{
		log: l,
		r:   roster.New(),
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.HandleConfirmOrDeny(jid.NR("foo@bar.com"), true)

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.WarnLevel)
	c.Assert(hook.LastEntry().Message, Equals, "No pending subscription")
	c.Assert(hook.LastEntry().Data["from"], Equals, "foo@bar.com")
}

func (s *SessionSuite) Test_HandleConfirmOrDeny_succeedsOnNotAllowed(c *C) {
	mockIn := &mockConnIOReaderWriter{}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	called := 0

	sess := &session{
		log:                 log.StandardLogger(),
		conn:                conn,
		r:                   roster.New(),
		sessionEventHandler: &mockSessionEventHandler{
			//warn: func(v string) {
			//	called++
			//},
		},
	}

	expectedPresence := `<presence xmlns="jabber:client" id="123" to="foo@bar.com" type="unsubscribed"></presence>`
	sess.r.SubscribeRequest(jid.NR("foo@bar.com"), "123", "")

	sess.HandleConfirmOrDeny(jid.NR("foo@bar.com"), false)

	c.Assert(called, Equals, 0)
	c.Assert(string(mockIn.write), Equals, expectedPresence)
	_, inMap := sess.r.GetPendingSubscribe(jid.NR("foo@bar.com"))
	c.Assert(inMap, Equals, false)
}

func (s *SessionSuite) Test_HandleConfirmOrDeny_succeedsOnAllowedAndAskBack(c *C) {
	mockIn := &mockConnIOReaderWriter{}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	called := 0

	sess := &session{
		log:                 log.StandardLogger(),
		conn:                conn,
		r:                   roster.New(),
		sessionEventHandler: &mockSessionEventHandler{
			//warn: func(v string) {
			//	called++
			//},
		},
	}

	expectedPresence := `<presence xmlns="jabber:client" id="123" to="foo@bar.com" type="subscribed"></presence>` +
		`<presence xmlns="jabber:client" id="[0-9]+" to="foo@bar.com" type="subscribe"></presence>`

	sess.r.SubscribeRequest(jid.NR("foo@bar.com"), "123", "")
	sess.HandleConfirmOrDeny(jid.NR("foo@bar.com"), true)

	c.Assert(called, Equals, 0)
	c.Assert(string(mockIn.write), Matches, expectedPresence)
	_, inMap := sess.r.GetPendingSubscribe(jid.NR("foo@bar.com"))
	c.Assert(inMap, Equals, false)
}

func (s *SessionSuite) Test_HandleConfirmOrDeny_handlesSendPresenceError(c *C) {
	mockIn := &mockConnIOReaderWriter{}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		&mockConnIOReaderWriter{err: errors.New("foo bar")},
		"some@one.org/foo",
	)

	l, hook := test.NewNullLogger()
	sess := &session{
		log: l,
		r:   roster.New(),
	}
	sess.conn = conn
	sess.r.SubscribeRequest(jid.NR("foo@bar.com"), "123", "")

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.HandleConfirmOrDeny(jid.NR("foo@bar.com"), true)

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.WarnLevel)
	c.Assert(hook.LastEntry().Message, Equals, "Error sending presence stanza")
	c.Assert(hook.LastEntry().Data["error"].(error).Error(), Equals, "foo bar")
}

func (s *SessionSuite) Test_watchTimeouts_cancelsTimedoutRequestsAndForgetsAboutThem(c *C) {
	tickInterval = time.Millisecond
	defer func() {
		tickInterval = time.Second
	}()

	now := time.Now()
	timeouts := map[data.Cookie]time.Time{
		data.Cookie(1): now.Add(-1 * time.Second),
		data.Cookie(2): now.Add(10 * time.Second),
	}

	sess := &session{
		log:        log.StandardLogger(),
		connStatus: CONNECTED,
		timeouts:   timeouts,
		conn:       xmpp.NewConn(nil, nil, ""),
	}

	go func() {
		time.Sleep(time.Duration(100) * time.Millisecond)
		sess.setConnStatus(DISCONNECTED)
	}()

	sess.watchTimeout()
	c.Assert(sess.timeouts, HasLen, 1)

	_, ok := sess.timeouts[data.Cookie(2)]
	c.Assert(ok, Equals, true)
}

type mockConvManager struct {
	getConversationWith    func(jid.Any) (otrclient.Conversation, bool)
	ensureConversationWith func(jid.Any, []byte) (otrclient.Conversation, bool)
	terminateAll           func()
}

func (mcm *mockConvManager) GetConversationWith(peer jid.Any) (otrclient.Conversation, bool) {
	return mcm.getConversationWith(peer)
}

func (mcm *mockConvManager) EnsureConversationWith(peer jid.Any, msg []byte) (otrclient.Conversation, bool) {
	return mcm.ensureConversationWith(peer, msg)
}

func (mcm *mockConvManager) TerminateAll() {
	mcm.terminateAll()
}

type mockConv struct {
	receive     func([]byte) ([]byte, error)
	isEncrypted func() bool
}

func (mc *mockConv) Receive(s []byte) ([]byte, error) {
	return mc.receive(s)
}

func (mc *mockConv) IsEncrypted() bool {
	return mc.isEncrypted()
}

func (mc *mockConv) Send([]byte) (trace int, err error) {
	return 0, nil
}

func (mc *mockConv) StartEncryptedChat() error {
	return nil
}

func (mc *mockConv) EndEncryptedChat() error {
	return nil
}

func (mc *mockConv) ProvideAuthenticationSecret([]byte) error {
	return nil
}

func (mc *mockConv) StartAuthenticate(string, []byte) error {
	return nil
}

func (mc *mockConv) AbortAuthentication() error {
	return nil
}

func (mc *mockConv) GetSSID() [8]byte {
	return [8]byte{}
}

func (mc *mockConv) EventHandler() *otrclient.EventHandler {
	return &otrclient.EventHandler{}
}

func (mc *mockConv) OurFingerprint() []byte {
	return nil
}

func (mc *mockConv) TheirFingerprint() []byte {
	return nil
}

func (mc *mockConv) CreateExtraSymmetricKey() ([]byte, error) {
	return nil, nil
}

func (mc *mockConv) GetAndWipeLastExtraKey() (usage uint32, usageData []byte, symkey []byte) {
	return 0, nil, nil
}

// func otrEventHandlerWith(s string, eh *otrclient.EventHandler) *otrclient.EventHandlers {
// 	ehs := otrclient.NewEventHandlers("one", func(jid.Any, *otrclient.EventHandler, chan string, chan int) {})
// 	ehs.Add(jid.Parse(s), eh)
// 	return ehs
// }

func (s *SessionSuite) Test_receiveClientMessage_willNotProcessBRTagsWhenNotEncrypted(c *C) {
	mcm := &mockConvManager{}
	sess := &session{
		log:         log.StandardLogger(),
		connStatus:  CONNECTED,
		convManager: mcm,
		config:      &config.ApplicationConfig{},
	}

	mc := &mockConv{}

	mc.receive = func(s3 []byte) ([]byte, error) {
		return s3, nil
	}

	mc.isEncrypted = func() bool {
		return false
	}

	mcm.ensureConversationWith = func(jid.Any, []byte) (otrclient.Conversation, bool) {
		return mc, false
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	go sess.receiveClientMessage(jid.R("someone@some.org/something"), time.Now(), "hello<br>ola<BR/>wazup?")

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Message)
		if !ok {
			return false
		}

		c.Assert(string(t.Body), Equals, "hello<br>ola<BR/>wazup?")
		c.Assert(t.Encrypted, Equals, false)
		return true
	})
}

func (s *SessionSuite) Test_receiveClientMessage_willProcessBRTagsWhenEncrypted(c *C) {
	mcm := &mockConvManager{}
	sess := &session{
		log:         log.StandardLogger(),
		connStatus:  CONNECTED,
		convManager: mcm,
		config:      &config.ApplicationConfig{},
	}

	mc := &mockConv{}
	mc.receive = func(s []byte) ([]byte, error) { return s, nil }
	mc.isEncrypted = func() bool { return true }
	mcm.ensureConversationWith = func(jid.Any, []byte) (otrclient.Conversation, bool) {
		return mc, false
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	go sess.receiveClientMessage(jid.R("someone@some.org/something"), time.Now(), "hello<br>ola<br/><BR/>wazup?")

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Message)
		if !ok {
			return false
		}

		c.Assert(string(t.Body), Equals, "hello\nola\n\nwazup?")
		c.Assert(t.Encrypted, Equals, true)
		return true
	})
}

type convManagerWithoutConversation struct{}

func (ncm *convManagerWithoutConversation) GetConversationWith(peer jid.Any) (otrclient.Conversation, bool) {
	return nil, false
}

func (ncm *convManagerWithoutConversation) EnsureConversationWith(peer jid.Any, msg []byte) (otrclient.Conversation, bool) {
	return nil, false
}

func (ncm *convManagerWithoutConversation) TerminateAll() {
}

func sessionWithConvMngrWithoutConvs() *session {
	return &session{
		log:         log.StandardLogger(),
		connStatus:  CONNECTED,
		convManager: &convManagerWithoutConversation{},
		config:      &config.ApplicationConfig{},
	}
}

func (s *SessionSuite) Test_logsError_whenWeStartSMPWithoutAConversation(c *C) {
	sess := sessionWithConvMngrWithoutConvs()
	l, hook := test.NewNullLogger()
	sess.log = l

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 1)
	sess.eventsReachedZero = eventsDone

	sess.StartSMP(jid.R("someone's peer/resource"), "Im a question", "im an answer")

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.ErrorLevel)
	c.Assert(hook.LastEntry().Message, Equals, "tried to start SMP when a conversation does not exist")
}

func (s *SessionSuite) Test_logsError_whenWeFinishSMPWithoutAConversation(c *C) {
	sess := sessionWithConvMngrWithoutConvs()
	l, hook := test.NewNullLogger()
	sess.log = l

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 1)
	sess.eventsReachedZero = eventsDone

	sess.FinishSMP(jid.R("someone's peer/resource"), "im an answer")

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.ErrorLevel)
	c.Assert(hook.LastEntry().Message, Equals, "tried to finish SMP when a conversation does not exist")
}

func (s *SessionSuite) Test_logsError_whenWeAbortSMPWithoutAConversation(c *C) {
	sess := sessionWithConvMngrWithoutConvs()
	l, hook := test.NewNullLogger()
	sess.log = l

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 1)
	sess.eventsReachedZero = eventsDone

	sess.AbortSMP(jid.R("someone's peer/resource"))

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.ErrorLevel)
	c.Assert(hook.LastEntry().Message, Equals, "tried to abort SMP when a conversation does not exist")
}

func (s *SessionSuite) Test_session_GetInMemoryLog_works(c *C) {
	b := &bytes.Buffer{}

	sess := &session{
		inMemoryLog: b,
	}

	c.Assert(sess.GetInMemoryLog(), Equals, b)
}

func (s *SessionSuite) Test_parseFromConfig_works(c *C) {
	data := alicePrivateKey.Serialize()
	acc := &config.Account{
		PrivateKeys: [][]byte{
			data,
			[]byte{0x99, 0x99, 0x99},
		},
	}

	res := parseFromConfig(acc)
	c.Assert(res, HasLen, 1)
	c.Assert(res[0], DeepEquals, alicePrivateKey)
}

func (s *SessionSuite) Test_CreateXMPPLogger_works(c *C) {
	orgDebug := *config.DebugFlag
	defer func() {
		*config.DebugFlag = orgDebug
	}()
	*config.DebugFlag = false

	inm, l := CreateXMPPLogger("")
	c.Assert(inm, IsNil)
	c.Assert(l, IsNil)
}

func (s *SessionSuite) Test_CreateXMPPLogger_createsMultiWriterWhenDebugFlag(c *C) {
	tf, _ := ioutil.TempFile("", "")
	defer os.Remove(tf.Name())
	_ = tf.Close()

	orgDebug := *config.DebugFlag
	defer func() {
		*config.DebugFlag = orgDebug
	}()
	*config.DebugFlag = true

	inm, l := CreateXMPPLogger(tf.Name())
	c.Assert(inm, Not(IsNil))
	c.Assert(l, Not(IsNil))
	c.Assert(l, Not(FitsTypeOf), &bytes.Buffer{})
}

func (s *SessionSuite) Test_CreateXMPPLogger_usesInMemoryBufferWhenDebugFlag(c *C) {
	orgDebug := *config.DebugFlag
	defer func() {
		*config.DebugFlag = orgDebug
	}()
	*config.DebugFlag = true

	inm, l := CreateXMPPLogger("")
	c.Assert(inm, Not(IsNil))
	c.Assert(l, Not(IsNil))
	c.Assert(l.(*bytes.Buffer), Equals, inm)
}

func (s *SessionSuite) Test_session_Send_returnsOfflineError(c *C) {
	sess := &session{
		connStatus: DISCONNECTED,
	}

	res := sess.Send(jid.Parse("hello@goodbye.com"), "something", false)
	c.Assert(res, ErrorMatches, "Couldn't send message since we are not connected")
}

func (s *SessionSuite) Test_session_Send_sends(c *C) {
	mockIn := &mockConnIOReaderWriter{}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log:        l,
		conn:       conn,
		connStatus: CONNECTED,
	}

	res := sess.Send(jid.Parse("hello@goodbye.com"), "something", false)
	c.Assert(res, IsNil)
	c.Assert(string(mockIn.write), Equals, "<message to='hello@goodbye.com' from='some@one.org/foo'"+
		" type='chat'><body>something</body><nos:x xmlns:nos='google:nosave' value='enabled'/>"+
		"<arc:record xmlns:arc='http://jabber.org/protocol/archive' otr='require'/>"+
		"<no-copy xmlns='urn:xmpp:hints'/><no-permanent-store xmlns='urn:xmpp:hints'/><private xmlns='urn:xmpp:carbons:2'/></message>")

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.DebugLevel)
	c.Assert(hook.LastEntry().Message, Equals, "Send()")
	c.Assert(hook.LastEntry().Data["to"], DeepEquals, jid.Parse("hello@goodbye.com"))
	c.Assert(hook.LastEntry().Data["sentMsg"], Equals, "something")
}

func closeToNow(t time.Time) bool {
	nw := time.Now()
	late := nw.Add(10 * time.Second)
	early := nw.Add(-10 * time.Second)

	return t.After(early) && t.Before(late)
}

func (s *SessionSuite) Test_retrieveMessageTime_returnsEmptyTimeIfNoDelayFound(c *C) {
	c.Assert(closeToNow(retrieveMessageTime(&data.ClientMessage{})), Equals, true)
	c.Assert(closeToNow(retrieveMessageTime(&data.ClientMessage{Delay: &data.Delay{}})), Equals, true)
}

func (s *SessionSuite) Test_retrieveMessageTime_returnsEmptyTimeIfCantParseTime(c *C) {
	c.Assert(closeToNow(retrieveMessageTime(&data.ClientMessage{Delay: &data.Delay{Stamp: "qqqqqqqqqqqqqq"}})), Equals, true)
}

func (s *SessionSuite) Test_retrieveMessageTime_returnsTimestamp(c *C) {
	expTime := time.Date(2012, 2, 3, 19, 11, 2, 0, time.UTC)
	c.Assert(expTime.Equal(retrieveMessageTime(&data.ClientMessage{Delay: &data.Delay{Stamp: "2012-02-03T19:11:02Z"}})), Equals, true)
}

func (s *SessionSuite) Test_session_receivedClientMessage_processesExtensions(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
	}

	cm := &data.ClientMessage{
		Body: "",
		Extensions: data.Extensions{
			&data.Extension{
				Body: "<un-unknown xmlns='urn:test:namespace'/>",
			},
		},
	}

	res := sess.receivedClientMessage(cm)

	c.Assert(res, Equals, true)
	c.Assert(hook.Entries, HasLen, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "receivedClientMessage()")
	c.Assert(hook.Entries[0].Data["stanza"], Not(Equals), "")
	c.Assert(hook.Entries[1].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[1].Message, Equals, "Unknown extension")
	c.Assert(fmt.Sprintf("%s", hook.Entries[1].Data["extension"]), Equals, "<un-unknown xmlns='urn:test:namespace'/>")
}

func (s *SessionSuite) Test_session_receivedClientMessage_works(c *C) {
	mcm := &mockConvManager{}
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)
	sess := &session{
		log:         l,
		connStatus:  CONNECTED,
		convManager: mcm,
		config:      &config.ApplicationConfig{},
	}

	mc := &mockConv{}

	mc.receive = func(s3 []byte) ([]byte, error) {
		return s3, nil
	}

	mc.isEncrypted = func() bool {
		return false
	}

	mcm.ensureConversationWith = func(jid.Any, []byte) (otrclient.Conversation, bool) {
		return mc, false
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	cm := &data.ClientMessage{
		Body: "hello",
		From: "some@example.org/foo",
	}

	res := sess.receivedClientMessage(cm)

	c.Assert(res, Equals, true)
	c.Assert(hook.Entries, HasLen, 1)
}

func (s *SessionSuite) Test_session_receivedClientMessage_processesEncryptionTag(c *C) {
	mcm := &mockConvManager{}
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)
	sess := &session{
		log:         l,
		connStatus:  CONNECTED,
		convManager: mcm,
		config:      &config.ApplicationConfig{},
	}

	mc := &mockConv{}

	mc.receive = func(s3 []byte) ([]byte, error) {
		return s3, nil
	}

	mc.isEncrypted = func() bool {
		return false
	}

	mcm.ensureConversationWith = func(jid.Any, []byte) (otrclient.Conversation, bool) {
		return mc, false
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	cm := &data.ClientMessage{
		Body: "",
		From: "some@example.org/foo",
		Encryption: &data.Encryption{
			Namespace: otrEncryptionNamespace,
		},
	}

	res := sess.receivedClientMessage(cm)

	c.Assert(res, Equals, true)
	c.Assert(hook.Entries, HasLen, 2)
}

func (s *SessionSuite) Test_session_receivedClientMessage_handlesRegularErrorType(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
	}

	cm := &data.ClientMessage{
		Body: "",
		Type: "error",
		From: "helu@baba.gan/foo",
		Error: &data.StanzaError{
			Type: "cancel",
			Text: "bla",
		},
	}

	res := sess.receivedClientMessage(cm)

	c.Assert(res, Equals, true)
	c.Assert(hook.Entries, HasLen, 2)
	c.Assert(hook.Entries[1].Level, Equals, log.ErrorLevel)
	c.Assert(hook.Entries[1].Message, Equals, "Error reported from peer")
	c.Assert(hook.Entries[1].Data["error"], FitsTypeOf, &data.StanzaError{})
	c.Assert(fmt.Sprintf("%s", hook.Entries[1].Data["peer"]), Equals, "helu@baba.gan")
}

func (s *SessionSuite) Test_session_receivedClientMessage_handlesMUCError(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
	}

	published := []interface{}{}

	sess.muc = newMUCManager(sess.log, sess.Conn, func(ev interface{}) {
		published = append(published, ev)
	})

	cm := &data.ClientMessage{
		Body: "",
		Type: "error",
		From: "helu@baba.gan/foo",
		Error: &data.StanzaError{
			Type:             "cancel",
			Text:             "bla",
			MUCNotAcceptable: &data.MUCNotAcceptable{},
		},
	}

	res := sess.receivedClientMessage(cm)

	c.Assert(res, Equals, true)
	c.Assert(hook.Entries, HasLen, 1)

	c.Assert(published, DeepEquals, []interface{}{
		events.MUCError{
			ErrorType: events.MUCMessageNotAcceptable,
			Room:      jid.ParseBare("helu@baba.gan"),
		},
	})
}

func (s *SessionSuite) Test_session_receivedClientMessage_handlesGroupChat(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
	}

	published := []interface{}{}

	sess.muc = newMUCManager(sess.log, sess.Conn, func(ev interface{}) {
		published = append(published, ev)
	})

	cm := &data.ClientMessage{
		Body: "",
		Type: "groupchat",
		From: "helu@baba.gan/foo",
	}

	res := sess.receivedClientMessage(cm)

	c.Assert(res, Equals, true)
	c.Assert(hook.Entries, HasLen, 2)
	c.Assert(hook.Entries[1].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[1].Message, Equals, "handleMUCReceivedClientMessage()")
	c.Assert(hook.Entries[1].Data["stanza"], Not(IsNil))

	c.Assert(published, HasLen, 0)
}

func (s *SessionSuite) Test_session_SendMUCMessage_failsWhenNotConnected(c *C) {
	sess := &session{
		connStatus: DISCONNECTED,
	}

	res := sess.SendMUCMessage("to@foo.com", "from@bar.com", "hello there")
	c.Assert(res, ErrorMatches, "session is not connected")
}

func (s *SessionSuite) Test_session_SendMUCMessage_works(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockIn := &mockConnIOReaderWriter{}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		conn:       conn,
		log:        l,
		connStatus: CONNECTED,
		config:     &config.ApplicationConfig{},
	}

	published := []interface{}{}

	sess.muc = newMUCManager(sess.log, sess.Conn, func(ev interface{}) {
		published = append(published, ev)
	})

	res := sess.SendMUCMessage("to@foo.com", "from@bar.com", "hello there")

	c.Assert(res, IsNil)
	c.Assert(hook.Entries, HasLen, 0)
	c.Assert(published, HasLen, 0)
	c.Assert(string(mockIn.write), Matches, `<message xmlns="jabber:client" from="from@bar.com" id="[0-9]+" to="to@foo.com" type="groupchat"><body>hello there</body></message>`)
}

func (s *SessionSuite) Test_session_CommandManager_accessors(c *C) {
	sess := &session{}
	sess.SetCommandManager(nil)
	c.Assert(sess.CommandManager(), IsNil)

	vv := &mockCommandManager{}
	sess.SetCommandManager(vv)
	c.Assert(sess.CommandManager(), Equals, vv)
}

func (s *SessionSuite) Test_session_SetWantToBeOnline(c *C) {
	sess := &session{}

	sess.SetWantToBeOnline(true)
	c.Assert(sess.wantToBeOnline, Equals, true)

	sess.SetWantToBeOnline(false)
	c.Assert(sess.wantToBeOnline, Equals, false)
}

func (s *SessionSuite) Test_session_PrivateKeys(c *C) {
	sess := &session{}

	vv := []otr3.PrivateKey{nil, nil, nil}

	sess.privateKeys = vv

	c.Assert(sess.PrivateKeys(), DeepEquals, vv)
}

func (s *SessionSuite) Test_session_R(c *C) {
	sess := &session{}

	vv := &roster.List{}

	sess.r = vv

	c.Assert(sess.R(), Equals, vv)
}

type mockConnector struct {
	connect func()
}

func (m *mockConnector) Connect() {
	if m.connect != nil {
		m.connect()
	}
}

func (s *SessionSuite) Test_session_SetConnector(c *C) {
	sess := &session{}
	vv := &mockConnector{}
	sess.SetConnector(vv)
	c.Assert(sess.connector, Equals, vv)
}

func (s *SessionSuite) Test_session_GroupDelimiter(c *C) {
	sess := &session{
		groupDelimiter: "abcfoo",
	}
	c.Assert(sess.GroupDelimiter(), Equals, "abcfoo")
}

func (s *SessionSuite) Test_session_Config(c *C) {
	vv := &config.ApplicationConfig{}

	sess := &session{
		config: vv,
	}

	c.Assert(sess.Config(), Equals, vv)
}

func (s *SessionSuite) Test_session_SetLastActionTime(c *C) {
	vv := time.Now()

	sess := &session{}

	sess.SetLastActionTime(vv)

	c.Assert(sess.lastActionTime, Equals, vv)
}

func (s *SessionSuite) Test_session_setStatus_setsConnected(c *C) {
	sess := &session{}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.setStatus(CONNECTED)

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Event)
		if !ok {
			return false
		}

		c.Assert(t.Type, Equals, events.Connected)
		return true
	})
}

func (s *SessionSuite) Test_session_setStatus_setsDisconnected(c *C) {
	sess := &session{}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.setStatus(DISCONNECTED)

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Event)
		if !ok {
			return false
		}

		c.Assert(t.Type, Equals, events.Disconnected)
		return true
	})
}

func (s *SessionSuite) Test_session_setStatus_setsConnecting(c *C) {
	sess := &session{}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.setStatus(CONNECTING)

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Event)
		if !ok {
			return false
		}

		c.Assert(t.Type, Equals, events.Connecting)
		return true
	})
}

func (s *SessionSuite) Test_session_receivedClientPresence_subscribe_handlesAutoApprove(c *C) {
	mockIn := &mockConnIOReaderWriter{}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	appr := map[string]bool{}
	sess := &session{
		autoApproves: appr,
		conn:         conn,
	}

	stanza := &data.ClientPresence{
		From: "hello@goodbye.com/compu",
		Type: "subscribe",
		ID:   "an-id-maybe",
	}

	appr["hello@goodbye.com"] = true

	res := sess.receivedClientPresence(stanza)
	c.Assert(res, Equals, true)

	_, ok := appr["hello@goodbye.com"]
	c.Assert(ok, Equals, false)
	c.Assert(string(mockIn.write), Equals, `<presence xmlns="jabber:client" id="an-id-maybe" to="hello@goodbye.com" type="subscribed"></presence>`)
}

func (s *SessionSuite) Test_session_receivedClientPresence_unavailable_forMUC(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockIn := &mockConnIOReaderWriter{}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		conn: conn,
		log:  l,
	}

	published := []interface{}{}

	sess.muc = newMUCManager(sess.log, sess.Conn, func(ev interface{}) {
		published = append(published, ev)
	})

	stanza := &data.ClientPresence{
		From: "hello@goodbye.com/compu",
		Type: "unavailable",
		ID:   "an-id-maybe",
		MUCUser: &data.MUCUser{
			Item: &data.MUCUserItem{},
		},
	}

	res := sess.receivedClientPresence(stanza)
	c.Assert(res, Equals, true)

	c.Assert(string(mockIn.write), Equals, ``)
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.ErrorLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Trying to get a room that is not in the room manager")
	c.Assert(hook.Entries[0].Data["method"], Equals, "handleOccupantLeft")
	c.Assert(hook.Entries[0].Data["occupant"], Equals, "compu")
	c.Assert(hook.Entries[0].Data["room"], DeepEquals, jid.ParseBare("hello@goodbye.com"))
	c.Assert(published, HasLen, 0)
}

func (s *SessionSuite) Test_session_receivedClientPresence_empty_withoutResource(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
	}

	stanza := &data.ClientPresence{
		From: "hello@goodbye.com",
		Type: "",
		ID:   "an-id-maybe",
	}

	res := sess.receivedClientPresence(stanza)
	c.Assert(res, Equals, true)

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Got a presence without resource in 'from' - this is likely an error")
	c.Assert(hook.Entries[0].Data["from"], Equals, "hello@goodbye.com")
	c.Assert(hook.Entries[0].Data["stanza"], Equals, stanza)
}

func (s *SessionSuite) Test_session_receivedClientPresence_empty_forMUC(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockIn := &mockConnIOReaderWriter{}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		conn: conn,
		log:  l,
	}

	published := []interface{}{}

	sess.muc = newMUCManager(sess.log, sess.Conn, func(ev interface{}) {
		published = append(published, ev)
	})

	stanza := &data.ClientPresence{
		From: "hello@goodbye.com/compu",
		Type: "",
		ID:   "an-id-maybe",
		MUCUser: &data.MUCUser{
			Item: &data.MUCUserItem{},
		},
	}

	res := sess.receivedClientPresence(stanza)
	c.Assert(res, Equals, true)

	c.Assert(string(mockIn.write), Equals, ``)
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.ErrorLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Trying to get a room that is not in the room manager")
	c.Assert(hook.Entries[0].Data["method"], Equals, "handleOccupantUpdate")
	c.Assert(hook.Entries[0].Data["occupant"], Equals, "compu")
	c.Assert(hook.Entries[0].Data["room"], DeepEquals, jid.ParseBare("hello@goodbye.com"))
	c.Assert(published, HasLen, 0)
}

func (s *SessionSuite) Test_session_receivedClientPresence_subscribed(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
		r:   roster.New(),
	}

	sess.r.AddOrReplace(&roster.Peer{Jid: jid.ParseBare("hello@goodbye.com")})

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	stanza := &data.ClientPresence{
		From: "hello@goodbye.com/compu",
		Type: "subscribed",
		ID:   "an-id-maybe",
	}

	res := sess.receivedClientPresence(stanza)
	c.Assert(res, Equals, true)
	v, _ := sess.r.Get(jid.ParseBare("hello@goodbye.com"))
	c.Assert(v.Subscription, Equals, "to")

	c.Assert(hook.Entries, HasLen, 0)

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Peer)
		if !ok {
			return false
		}

		c.Assert(t.Type, Equals, events.Subscribed)
		c.Assert(t.From, Equals, jid.Parse("hello@goodbye.com"))
		return true
	})
}

func (s *SessionSuite) Test_session_receivedClientPresence_unsubscribe(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
		r:   roster.New(),
	}

	sess.r.AddOrReplace(&roster.Peer{Jid: jid.ParseBare("hello@goodbye.com")})

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	stanza := &data.ClientPresence{
		From: "hello@goodbye.com/compu",
		Type: "unsubscribe",
		ID:   "an-id-maybe",
	}

	res := sess.receivedClientPresence(stanza)
	c.Assert(res, Equals, true)
	v, _ := sess.r.Get(jid.ParseBare("hello@goodbye.com"))
	c.Assert(v.Subscription, Equals, "")

	c.Assert(hook.Entries, HasLen, 0)

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Peer)
		if !ok {
			return false
		}

		c.Assert(t.Type, Equals, events.Unsubscribe)
		c.Assert(t.From, Equals, jid.Parse("hello@goodbye.com"))
		return true
	})

}

func (s *SessionSuite) Test_session_receivedClientPresence_unsubscribed(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
		r:   roster.New(),
	}

	stanza := &data.ClientPresence{
		From: "hello@goodbye.com/compu",
		Type: "unsubscribed",
		ID:   "an-id-maybe",
	}

	res := sess.receivedClientPresence(stanza)
	c.Assert(res, Equals, true)
	c.Assert(hook.Entries, HasLen, 0)
}

func (s *SessionSuite) Test_session_receivedClientPresence_error(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
		r:   roster.New(),
	}

	sess.r.AddOrReplace(&roster.Peer{Jid: jid.ParseBare("hello@goodbye.com")})

	stanza := &data.ClientPresence{
		From: "hello@goodbye.com/compu",
		Type: "error",
		ID:   "an-id-maybe",
		Error: &data.StanzaError{
			By:   "hmm",
			Code: "12342",
			Type: "modify",
			Text: "what are you really thinking",
		},
	}

	stanza.Error.Condition.XMLName.Space = "urn:test:42"
	stanza.Error.Condition.XMLName.Local = "summathing"
	stanza.Error.Condition.Body = "<hello>foo</hello>"

	res := sess.receivedClientPresence(stanza)
	c.Assert(res, Equals, true)
	v, _ := sess.r.Get(jid.ParseBare("hello@goodbye.com"))

	c.Assert(v.LatestError, DeepEquals, &roster.PeerError{
		Code: "12342",
		Type: "modify",
		More: "urn:test:42 summathing",
	})

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.ErrorLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Got a presence error")
	c.Assert(hook.Entries[0].Data["from"], Equals, "hello@goodbye.com/compu")
	c.Assert(hook.Entries[0].Data["error"], Equals, stanza.Error)
}

func (s *SessionSuite) Test_session_receivedClientPresence_MUCerror(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
		r:   roster.New(),
	}

	published := []interface{}{}

	sess.muc = newMUCManager(sess.log, sess.Conn, func(ev interface{}) {
		published = append(published, ev)
	})

	sess.r.AddOrReplace(&roster.Peer{Jid: jid.ParseBare("hello@goodbye.com")})

	stanza := &data.ClientPresence{
		From: "hello@goodbye.com/compu",
		Type: "error",
		ID:   "an-id-maybe",
		Error: &data.StanzaError{
			By:            "hmm2",
			Code:          "12343",
			Type:          "cancel",
			Text:          "you think so?",
			MUCNotAllowed: &data.MUCNotAllowed{},
		},
	}

	stanza.Error.Condition.XMLName.Space = "urn:test:43"
	stanza.Error.Condition.XMLName.Local = "else"
	stanza.Error.Condition.Body = "<hello>foo</hello>"

	res := sess.receivedClientPresence(stanza)
	c.Assert(res, Equals, true)
	v, _ := sess.r.Get(jid.ParseBare("hello@goodbye.com"))

	c.Assert(v.LatestError, DeepEquals, &roster.PeerError{
		Code: "12343",
		Type: "cancel",
		More: "urn:test:43 else",
	})

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.ErrorLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Got a presence error")
	c.Assert(hook.Entries[0].Data["from"], Equals, "hello@goodbye.com/compu")
	c.Assert(hook.Entries[0].Data["error"], Equals, stanza.Error)

	c.Assert(published, HasLen, 1)
	c.Assert(published[0], DeepEquals, events.MUCError{
		ErrorType: events.MUCNotAllowed,
		Room:      jid.ParseBare("hello@goodbye.com"),
		Nickname:  "compu",
	})
}

func (s *SessionSuite) Test_session_receivedClientPresence_unknown_type(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
	}

	stanza := &data.ClientPresence{
		From: "hello@goodbye.com/compu",
		Type: "something else",
		ID:   "an-id-maybe",
	}

	res := sess.receivedClientPresence(stanza)
	c.Assert(res, Equals, true)

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Unrecognized presence")
	c.Assert(hook.Entries[0].Data["from"], Equals, "hello@goodbye.com/compu")
	c.Assert(hook.Entries[0].Data["stanza"], Equals, stanza)
}
