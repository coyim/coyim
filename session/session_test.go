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
	"github.com/coyim/coyim/tls"
	"github.com/coyim/coyim/xmpp"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/coyim/xmpp/mock"
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
	send             func([]byte) (int, error)
	receive          func([]byte) ([]byte, error)
	isEncrypted      func() bool
	endEncryptedChat func() error
	eh               *otrclient.EventHandler
}

func (mc *mockConv) Receive(s []byte) ([]byte, error) {
	return mc.receive(s)
}

func (mc *mockConv) IsEncrypted() bool {
	return mc.isEncrypted()
}

func (mc *mockConv) Send(v []byte) (trace int, err error) {
	if mc.send != nil {
		return mc.send(v)
	}
	return 0, nil
}

func (mc *mockConv) StartEncryptedChat() error {
	return nil
}

func (mc *mockConv) EndEncryptedChat() error {
	if mc.endEncryptedChat != nil {
		return mc.endEncryptedChat()
	}
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
	if mc.eh != nil {
		return mc.eh
	}
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

func (s *SessionSuite) Test_session_AwaitVersionReply_failsOnClosedChannel(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
	}

	ch := make(chan data.Stanza, 1)

	close(ch)

	sess.AwaitVersionReply(ch, "foobarium@example.org/hello")

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Version request timed out")
	c.Assert(hook.Entries[0].Data["user"], Equals, "foobarium@example.org/hello")
}

func (s *SessionSuite) Test_session_AwaitVersionReply_works(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
	}

	ch := make(chan data.Stanza, 1)

	st := data.Stanza{
		Value: &data.ClientIQ{
			Type: "result",
			Query: []byte(`<query xmlns="jabber:iq:version">
  <name>One</name>
  <version>Two</version>
  <os>Three</os>
</query>`),
		},
	}

	ch <- st

	sess.AwaitVersionReply(ch, "foobarium@example.org/hello")

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Version reply from peer")
	c.Assert(hook.Entries[0].Data["user"], Equals, "foobarium@example.org/hello")
	c.Assert(hook.Entries[0].Data["version"], DeepEquals, data.VersionReply{
		XMLName: xml.Name{Space: "jabber:iq:version", Local: "query"},
		Name:    "One",
		Version: "Two",
		OS:      "Three",
	})
}

func (s *SessionSuite) Test_session_AwaitVersionReply_failsWhenNotIQ(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
	}

	ch := make(chan data.Stanza, 1)

	st := data.Stanza{
		Value: "something else",
	}

	ch <- st

	sess.AwaitVersionReply(ch, "foobarium@example.org/hello")

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Version request resulted in bad reply type")
	c.Assert(hook.Entries[0].Data["user"], Equals, "foobarium@example.org/hello")
}

func (s *SessionSuite) Test_session_AwaitVersionReply_failsWhenStanzaError(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
	}

	ch := make(chan data.Stanza, 1)

	st := data.Stanza{
		Value: &data.ClientIQ{
			Type: "error",
			Query: []byte(`<query xmlns="jabber:iq:version">
  <name>One</name>
  <version>Two</version>
  <os>Three</os>
</query>`),
		},
	}

	ch <- st

	sess.AwaitVersionReply(ch, "foobarium@example.org/hello")

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Version request resulted in XMPP error")
	c.Assert(hook.Entries[0].Data["user"], Equals, "foobarium@example.org/hello")
}

func (s *SessionSuite) Test_session_AwaitVersionReply_failsWhenUnknownType(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
	}

	ch := make(chan data.Stanza, 1)

	st := data.Stanza{
		Value: &data.ClientIQ{
			Type: "unknownium",
			Query: []byte(`<query xmlns="jabber:iq:version">
  <name>One</name>
  <version>Two</version>
  <os>Three</os>
</query>`),
		},
	}

	ch <- st

	sess.AwaitVersionReply(ch, "foobarium@example.org/hello")

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Version request resulted in response with unknown type")
	c.Assert(hook.Entries[0].Data["user"], Equals, "foobarium@example.org/hello")
	c.Assert(hook.Entries[0].Data["type"], Equals, "unknownium")
}

func (s *SessionSuite) Test_session_AwaitVersionReply_failsWhenBadXML(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
	}

	ch := make(chan data.Stanza, 1)

	st := data.Stanza{
		Value: &data.ClientIQ{
			Type:  "result",
			Query: []byte(`<query`),
		},
	}

	ch <- st

	sess.AwaitVersionReply(ch, "foobarium@example.org/hello")

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Failed to parse version reply")
	c.Assert(hook.Entries[0].Data["user"], Equals, "foobarium@example.org/hello")
	c.Assert(hook.Entries[0].Data["error"], ErrorMatches, "XML syntax error.*")
}

func (s *SessionSuite) Test_peerFrom_works(c *C) {
	p := peerFrom(data.RosterEntry{
		Jid:          "romeo@example.net",
		Subscription: "both",
		Name:         "Mo",
		Group:        []string{"Foes"},
	}, &config.Account{})

	c.Assert(p.Jid, DeepEquals, jid.Parse("romeo@example.net"))
	c.Assert(p.Subscription, Equals, "both")
	c.Assert(p.Name, Equals, "Mo")
	c.Assert(p.Nickname, Equals, "")
	c.Assert(p.Groups, DeepEquals, map[string]bool{"Foes": true})

	ac := &config.Account{
		Peers: []*config.Peer{
			&config.Peer{
				UserID:   "romeo@example.net",
				Nickname: "blaha",
				Groups: []string{
					"something",
					"else::bar",
				},
			},
		},
	}

	p = peerFrom(data.RosterEntry{
		Jid:          "romeo@example.net",
		Subscription: "both",
		Name:         "Mo",
		Group:        []string{"Foes"},
	}, ac)
	c.Assert(p.Jid, DeepEquals, jid.Parse("romeo@example.net"))
	c.Assert(p.Subscription, Equals, "both")
	c.Assert(p.Name, Equals, "Mo")
	c.Assert(p.Nickname, Equals, "blaha")
	c.Assert(p.BelongsTo, Equals, ac.ID())
	c.Assert(p.Groups, DeepEquals, map[string]bool{"something": true, "else::bar": true})
}

func (s *SessionSuite) Test_receiveClientMessage_logsConversationReceivalError(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mcm := &mockConvManager{}
	sess := &session{
		log:         l,
		connStatus:  CONNECTED,
		convManager: mcm,
		config:      &config.ApplicationConfig{},
	}

	mc := &mockConv{}
	mc.receive = func(s []byte) ([]byte, error) { return nil, errors.New("marker error") }
	mc.isEncrypted = func() bool { return true }
	mcm.ensureConversationWith = func(jid.Any, []byte) (otrclient.Conversation, bool) {
		return mc, false
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.receiveClientMessage(jid.R("someone@some.org/something"), time.Now(), "hello")

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.ErrorLevel)
	c.Assert(hook.Entries[0].Message, Equals, "While processing message from peer")
	c.Assert(hook.Entries[0].Data["peer"], Equals, "someone@some.org/something")
	c.Assert(hook.Entries[0].Data["error"], ErrorMatches, "marker error")
}

func (s *SessionSuite) Test_receiveClientMessage_handlesNewOTRKeys(c *C) {
	eh := &otrclient.EventHandler{
		Log: log.StandardLogger(),
	}
	mcm := &mockConvManager{}
	sess := &session{
		log:         log.StandardLogger(),
		connStatus:  CONNECTED,
		convManager: mcm,
		config:      &config.ApplicationConfig{},
	}

	mc := &mockConv{}
	mc.eh = eh
	mc.receive = func(s []byte) ([]byte, error) { return nil, nil }
	mc.isEncrypted = func() bool { return true }
	mcm.ensureConversationWith = func(jid.Any, []byte) (otrclient.Conversation, bool) {
		return mc, false
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	eh.HandleSecurityEvent(otr3.GoneSecure)
	sess.receiveClientMessage(jid.R("someone@some.org/something"), time.Now(), "hello")

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Peer)
		if !ok {
			return false
		}

		c.Assert(t.From, DeepEquals, jid.Parse("someone@some.org/something"))
		c.Assert(t.Type, Equals, events.OTRNewKeys)
		return true
	})
}

func (s *SessionSuite) Test_receiveClientMessage_handlesRenewedOTRKeys(c *C) {
	eh := &otrclient.EventHandler{
		Log: log.StandardLogger(),
	}
	mcm := &mockConvManager{}
	sess := &session{
		log:         log.StandardLogger(),
		connStatus:  CONNECTED,
		convManager: mcm,
		config:      &config.ApplicationConfig{},
	}

	mc := &mockConv{}
	mc.eh = eh
	mc.receive = func(s []byte) ([]byte, error) { return nil, nil }
	mc.isEncrypted = func() bool { return true }
	mcm.ensureConversationWith = func(jid.Any, []byte) (otrclient.Conversation, bool) {
		return mc, false
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	eh.HandleSecurityEvent(otr3.StillSecure)
	sess.receiveClientMessage(jid.R("someone@some.org/something"), time.Now(), "hello")

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Peer)
		if !ok {
			return false
		}

		c.Assert(t.From, DeepEquals, jid.Parse("someone@some.org/something"))
		c.Assert(t.Type, Equals, events.OTRRenewedKeys)
		return true
	})
}

func (s *SessionSuite) Test_receiveClientMessage_handlesConversationEnded(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	eh := &otrclient.EventHandler{
		Log: l,
	}
	mcm := &mockConvManager{}
	ac := &config.Account{}
	sess := &session{
		log:           l,
		connStatus:    CONNECTED,
		convManager:   mcm,
		config:        &config.ApplicationConfig{},
		accountConfig: ac,
	}

	mc := &mockConv{}
	mc.eh = eh
	mc.receive = func(s []byte) ([]byte, error) { return nil, nil }
	mc.isEncrypted = func() bool { return true }
	mcm.ensureConversationWith = func(jid.Any, []byte) (otrclient.Conversation, bool) {
		return mc, false
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	eh.HandleSecurityEvent(otr3.GoneInsecure)
	sess.receiveClientMessage(jid.R("someone@some.org/something"), time.Now(), "hello")

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Peer)
		if !ok {
			return false
		}

		c.Assert(t.From, DeepEquals, jid.Parse("someone@some.org/something"))
		c.Assert(t.Type, Equals, events.OTREnded)
		return true
	})

	c.Assert(hook.Entries, HasLen, 2)

	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "HandleSecurityEvent()")
	c.Assert(hook.Entries[0].Data["event"], Equals, otr3.GoneInsecure)

	c.Assert(hook.Entries[1].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[1].Message, Equals, "Peer has ended the secure conversation. You should do likewise")
	c.Assert(hook.Entries[1].Data["peer"], DeepEquals, jid.Parse("someone@some.org/something"))
}

func (s *SessionSuite) Test_receiveClientMessage_handlesConversationEnded_withAutoTearDownFailing(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	eh := &otrclient.EventHandler{
		Log: l,
	}
	mcm := &mockConvManager{}
	ac := &config.Account{OTRAutoTearDown: true}
	sess := &session{
		log:           l,
		connStatus:    CONNECTED,
		convManager:   mcm,
		config:        &config.ApplicationConfig{},
		accountConfig: ac,
	}

	mc := &mockConv{}
	mc.eh = eh
	mc.receive = func(s []byte) ([]byte, error) { return nil, nil }
	mc.endEncryptedChat = func() error { return errors.New("another marker error") }
	mc.isEncrypted = func() bool { return true }
	mcm.ensureConversationWith = func(jid.Any, []byte) (otrclient.Conversation, bool) {
		return mc, false
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	eh.HandleSecurityEvent(otr3.GoneInsecure)
	sess.receiveClientMessage(jid.R("someone@some.org/something"), time.Now(), "hello")

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Peer)
		if !ok {
			return false
		}

		c.Assert(t.From, DeepEquals, jid.Parse("someone@some.org/something"))
		c.Assert(t.Type, Equals, events.OTREnded)
		return true
	})

	c.Assert(hook.Entries, HasLen, 3)

	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "HandleSecurityEvent()")
	c.Assert(hook.Entries[0].Data["event"], Equals, otr3.GoneInsecure)

	c.Assert(hook.Entries[1].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[1].Message, Equals, "Peer has ended the secure conversation.")
	c.Assert(hook.Entries[1].Data["peer"], DeepEquals, jid.Parse("someone@some.org/something"))

	c.Assert(hook.Entries[2].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[2].Message, Equals, "Unable to automatically tear down OTR conversation with peer")
	c.Assert(hook.Entries[2].Data["peer"], DeepEquals, jid.Parse("someone@some.org/something"))
	c.Assert(hook.Entries[2].Data["error"], ErrorMatches, "another marker error")
}

func (s *SessionSuite) Test_receiveClientMessage_handlesConversationEnded_withAutoTearDownSucceeding(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	eh := &otrclient.EventHandler{
		Log: l,
	}
	mcm := &mockConvManager{}
	ac := &config.Account{OTRAutoTearDown: true}
	sess := &session{
		log:           l,
		connStatus:    CONNECTED,
		convManager:   mcm,
		config:        &config.ApplicationConfig{},
		accountConfig: ac,
	}

	mc := &mockConv{}
	mc.eh = eh
	mc.receive = func(s []byte) ([]byte, error) { return nil, nil }
	mc.isEncrypted = func() bool { return true }
	mcm.ensureConversationWith = func(jid.Any, []byte) (otrclient.Conversation, bool) {
		return mc, false
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	eh.HandleSecurityEvent(otr3.GoneInsecure)
	sess.receiveClientMessage(jid.R("someone@some.org/something"), time.Now(), "hello")

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Peer)
		if !ok {
			return false
		}

		c.Assert(t.From, DeepEquals, jid.Parse("someone@some.org/something"))
		c.Assert(t.Type, Equals, events.OTREnded)
		return true
	})

	c.Assert(hook.Entries, HasLen, 3)

	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "HandleSecurityEvent()")
	c.Assert(hook.Entries[0].Data["event"], Equals, otr3.GoneInsecure)

	c.Assert(hook.Entries[1].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[1].Message, Equals, "Peer has ended the secure conversation.")
	c.Assert(hook.Entries[1].Data["peer"], DeepEquals, jid.Parse("someone@some.org/something"))

	c.Assert(hook.Entries[2].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[2].Message, Equals, "Secure session with peer has been automatically ended. Messages will be sent in the clear until another OTR session is established.")
	c.Assert(hook.Entries[2].Data["peer"], DeepEquals, jid.Parse("someone@some.org/something"))
}

func (s *SessionSuite) Test_receiveClientMessage_handlesSMPSecretNeeded(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)
	eh := &otrclient.EventHandler{
		Log: l,
	}
	mcm := &mockConvManager{}
	sess := &session{
		log:         l,
		connStatus:  CONNECTED,
		convManager: mcm,
		config:      &config.ApplicationConfig{},
	}

	mc := &mockConv{}
	mc.eh = eh
	mc.receive = func(s []byte) ([]byte, error) { return nil, nil }
	mc.isEncrypted = func() bool { return true }
	mcm.ensureConversationWith = func(jid.Any, []byte) (otrclient.Conversation, bool) {
		return mc, false
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	eh.HandleSMPEvent(otr3.SMPEventAskForSecret, 42, "the life, universe and everything")
	sess.receiveClientMessage(jid.R("someone@some.org/something"), time.Now(), "hello")

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.SMP)
		if !ok {
			return false
		}

		c.Assert(t.From, DeepEquals, jid.Parse("someone@some.org/something"))
		c.Assert(t.Type, Equals, events.SecretNeeded)
		c.Assert(t.Body, Equals, "the life, universe and everything")
		return true
	})

	c.Assert(hook.Entries, HasLen, 1)
}

func (s *SessionSuite) Test_receiveClientMessage_handlesSMPFailed(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)
	eh := &otrclient.EventHandler{
		Log: l,
	}
	mcm := &mockConvManager{}
	sess := &session{
		log:         l,
		connStatus:  CONNECTED,
		convManager: mcm,
		config:      &config.ApplicationConfig{},
	}

	mc := &mockConv{}
	mc.eh = eh
	mc.receive = func(s []byte) ([]byte, error) { return nil, nil }
	mc.isEncrypted = func() bool { return true }
	mcm.ensureConversationWith = func(jid.Any, []byte) (otrclient.Conversation, bool) {
		return mc, false
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	eh.HandleSMPEvent(otr3.SMPEventCheated, 42, "the life, universe and everything")
	sess.receiveClientMessage(jid.R("someone@some.org/something"), time.Now(), "hello")

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.SMP)
		if !ok {
			return false
		}

		c.Assert(t.From, DeepEquals, jid.Parse("someone@some.org/something"))
		c.Assert(t.Type, Equals, events.Failure)
		c.Assert(t.Body, Equals, "")
		return true
	})

	c.Assert(hook.Entries, HasLen, 1)
}

func (s *SessionSuite) Test_receiveClientMessage_handlesSMPComplete(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)
	eh := &otrclient.EventHandler{
		Log: l,
	}
	mcm := &mockConvManager{}

	var commands []interface{}
	mcmdm := &mockCommandManager{
		exec: func(v interface{}) {
			commands = append(commands, v)
		},
	}

	sess := &session{
		log:         l,
		connStatus:  CONNECTED,
		convManager: mcm,
		cmdManager:  mcmdm,
		config:      &config.ApplicationConfig{},
	}

	mc := &mockConv{}
	mc.eh = eh
	mc.receive = func(s []byte) ([]byte, error) { return nil, nil }
	mc.isEncrypted = func() bool { return true }
	mcm.ensureConversationWith = func(jid.Any, []byte) (otrclient.Conversation, bool) {
		return mc, false
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	eh.HandleSMPEvent(otr3.SMPEventSuccess, 100, "the life, universe and everything")
	sess.receiveClientMessage(jid.R("someone@some.org/something"), time.Now(), "hello")

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.SMP)
		if !ok {
			return false
		}

		c.Assert(t.From, DeepEquals, jid.Parse("someone@some.org/something"))
		c.Assert(t.Type, Equals, events.Success)
		c.Assert(t.Body, Equals, "")
		return true
	})

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(commands, HasLen, 1)
	c.Assert(commands[0], DeepEquals, otrclient.AuthorizeFingerprintCmd{
		Account:     nil,
		Session:     sess,
		Peer:        jid.NR("someone@some.org"),
		Fingerprint: nil,
		Tag:         "SMP",
	})
}

func (s *SessionSuite) Test_session_Timeout(c *C) {
	sess := &session{
		timeouts: make(map[data.Cookie]time.Time),
	}

	tt := time.Now()

	sess.Timeout(42, tt)

	c.Assert(sess.timeouts[data.Cookie(42)], DeepEquals, tt)
}

func (s *SessionSuite) Test_waitForNextRosterRequest(c *C) {
	orgRosterRequestDelay := rosterRequestDelay
	defer func() {
		rosterRequestDelay = orgRosterRequestDelay
	}()

	rosterRequestDelay = 1 * time.Millisecond

	waitForNextRosterRequest()
}

func (s *SessionSuite) Test_session_SendPing_failsIfConnectionFails(c *C) {
	mockIn := &mockConnIOReaderWriter{
		err: errors.New("another marker"),
	}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log:        l,
		connStatus: CONNECTED,
	}
	sess.conn = conn

	sess.SendPing()

	c.Assert(hook.Entries, HasLen, 1)

	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Failure to ping server")
	c.Assert(hook.Entries[0].Data["error"], ErrorMatches, "another marker")
}

func (s *SessionSuite) Test_session_SendPing_timesOut(c *C) {
	orgPingTimeout := pingTimeout
	defer func() {
		pingTimeout = orgPingTimeout
	}()

	pingTimeout = 1 * time.Millisecond

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
		connStatus: CONNECTED,
	}
	sess.conn = conn

	sess.SendPing()

	c.Assert(sess.connStatus, Equals, DISCONNECTED)

	c.Assert(hook.Entries, HasLen, 1)

	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Ping timeout. Disconnecting...")
}

type sendPingMockConn struct {
	*mock.Conn
	f func() (<-chan data.Stanza, data.Cookie, error)
}

func (m *sendPingMockConn) SendPing() (<-chan data.Stanza, data.Cookie, error) {
	return m.f()
}

func (s *SessionSuite) Test_session_SendPing_works(c *C) {
	ch := make(chan data.Stanza, 1)
	conn := &sendPingMockConn{
		f: func() (<-chan data.Stanza, data.Cookie, error) {
			return ch, 0, nil
		},
	}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		conn:       conn,
		log:        l,
		connStatus: CONNECTED,
	}

	ch <- data.Stanza{
		Value: &data.ClientIQ{
			Type: "result",
		},
	}

	sess.SendPing()

	c.Assert(sess.connStatus, Equals, CONNECTED)
	c.Assert(hook.Entries, HasLen, 0)
}

func (s *SessionSuite) Test_session_SendPing_failsWithWeirdStanzaMixup(c *C) {
	ch := make(chan data.Stanza, 1)
	conn := &sendPingMockConn{
		f: func() (<-chan data.Stanza, data.Cookie, error) {
			return ch, 0, nil
		},
	}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		conn:       conn,
		log:        l,
		connStatus: CONNECTED,
	}

	ch <- data.Stanza{
		Value: "hello",
	}

	sess.SendPing()

	c.Assert(sess.connStatus, Equals, CONNECTED)
	c.Assert(hook.Entries, HasLen, 1)

	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Server returned weird result")
	c.Assert(hook.Entries[0].Data["stanza"], DeepEquals, "hello")
}

func (s *SessionSuite) Test_session_SendPing_failsIfServerDoesntSupportPing(c *C) {
	ch := make(chan data.Stanza, 1)
	conn := &sendPingMockConn{
		f: func() (<-chan data.Stanza, data.Cookie, error) {
			return ch, 0, nil
		},
	}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		conn:       conn,
		log:        l,
		connStatus: CONNECTED,
	}

	ch <- data.Stanza{
		Value: &data.ClientIQ{
			Type: "error",
		},
	}

	sess.SendPing()

	c.Assert(sess.connStatus, Equals, CONNECTED)
	c.Assert(hook.Entries, HasLen, 1)

	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Server does not support Ping")
}

type mockRunnable struct {
	ret error
}

func (m *mockRunnable) Run() error {
	return m.ret
}

func (s *SessionSuite) Test_session_maybeNotify_works(c *C) {
	orgExecCommand := execCommand
	defer func() {
		execCommand = orgExecCommand
	}()

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
		config: &config.ApplicationConfig{
			IdleSecondsBeforeNotification: 1,
			NotifyCommand: []string{
				"hello",
				"goodbye",
				"somewhere",
			},
		},
		lastActionTime: time.Now().Add(-1000 * time.Second),
	}

	res := orgExecCommand("foo", "bar")
	c.Assert(res, Not(IsNil))

	called := []string{}

	execCommand = func(name string, arg ...string) runnable {
		called = append(called, name)
		called = append(called, arg...)
		return &mockRunnable{}
	}

	sess.maybeNotify()

	c.Assert(called, DeepEquals, []string{"hello", "goodbye", "somewhere"})
	c.Assert(hook.Entries, HasLen, 0)
}

func (s *SessionSuite) Test_session_maybeNotify_logsErrorWhenFailingCommand(c *C) {
	orgExecCommand := execCommand
	defer func() {
		execCommand = orgExecCommand
	}()

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log: l,
		config: &config.ApplicationConfig{
			IdleSecondsBeforeNotification: 1,
			NotifyCommand: []string{
				"hello",
				"goodbye",
				"somewhere",
			},
		},
		lastActionTime: time.Now().Add(-1000 * time.Second),
	}

	execCommand = func(string, ...string) runnable {
		return &mockRunnable{errors.New("a nice marker error")}
	}

	sess.maybeNotify()

	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.ErrorLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Failed to run notify command")
	c.Assert(hook.Entries[0].Data["error"], ErrorMatches, "a nice marker error")
}

type requestVCardMockConn struct {
	*mock.Conn
	f func() (<-chan data.Stanza, data.Cookie, error)
}

func (m *requestVCardMockConn) RequestVCard() (<-chan data.Stanza, data.Cookie, error) {
	return m.f()
}

func (s *SessionSuite) Test_session_getVCard_works(c *C) {
	mc := &requestVCardMockConn{}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		conn:       mc,
		log:        l,
		connStatus: CONNECTED,
	}

	ch := make(chan data.Stanza, 1)

	mc.f = func() (<-chan data.Stanza, data.Cookie, error) {
		return ch, 0, nil
	}

	stz := data.Stanza{
		Value: &data.ClientIQ{
			Query: []byte(`<vCard xmlns="vcard-temp">
<FN>Hello</FN>
<NICKNAME>Again</NICKNAME>
</vCard>`),
		},
	}

	ch <- stz

	sess.getVCard()

	c.Assert(sess.nicknames, DeepEquals, []string{"Again", "Hello"})
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Fetching VCard")
}

func (s *SessionSuite) Test_session_getVCard_reportsErrorWhenParsingTheXML(c *C) {
	mc := &requestVCardMockConn{}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		conn:       mc,
		log:        l,
		connStatus: CONNECTED,
	}

	ch := make(chan data.Stanza, 1)

	mc.f = func() (<-chan data.Stanza, data.Cookie, error) {
		return ch, 0, nil
	}

	stz := data.Stanza{
		Value: &data.ClientIQ{
			Query: []byte(`<vCard xmlns="vcard-temp">
<FN>Hel`),
		},
	}

	ch <- stz

	sess.getVCard()

	c.Assert(sess.nicknames, IsNil)
	c.Assert(hook.Entries, HasLen, 2)
	c.Assert(hook.Entries[1].Level, Equals, log.ErrorLevel)
	c.Assert(hook.Entries[1].Message, Equals, "Failed to parse vcard")
	c.Assert(hook.Entries[1].Data["error"], ErrorMatches, "XML syntax error.*")
}

func (s *SessionSuite) Test_session_getVCard_failsWhenChannelIsClosed(c *C) {
	mc := &requestVCardMockConn{}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		conn:       mc,
		log:        l,
		connStatus: CONNECTED,
	}

	ch := make(chan data.Stanza, 1)

	mc.f = func() (<-chan data.Stanza, data.Cookie, error) {
		return ch, 0, nil
	}

	close(ch)

	sess.getVCard()

	c.Assert(sess.nicknames, IsNil)
	c.Assert(hook.Entries, HasLen, 2)
	c.Assert(hook.Entries[1].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[1].Message, Equals, "session: vcard request cancelled or timed out")
}

func (s *SessionSuite) Test_session_getVCard_failsWhenRequestFails(c *C) {
	mc := &requestVCardMockConn{}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		conn:       mc,
		log:        l,
		connStatus: CONNECTED,
	}

	mc.f = func() (<-chan data.Stanza, data.Cookie, error) {
		return nil, 0, errors.New("another marker error")
	}

	sess.getVCard()

	c.Assert(sess.nicknames, IsNil)
	c.Assert(hook.Entries, HasLen, 2)
	c.Assert(hook.Entries[1].Level, Equals, log.ErrorLevel)
	c.Assert(hook.Entries[1].Message, Equals, "Failed to request vcard")
	c.Assert(hook.Entries[1].Data["error"], ErrorMatches, "another marker error")
}

func (s *SessionSuite) Test_session_getVCard_doesntDoAnythingIfNotConnected(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log:        l,
		connStatus: DISCONNECTED,
	}

	sess.getVCard()

	c.Assert(sess.nicknames, IsNil)
	c.Assert(hook.Entries, HasLen, 0)
}

func (s *SessionSuite) Test_session_DisplayName_works(c *C) {
	c.Assert((&session{accountConfig: &config.Account{}, nicknames: []string{}}).DisplayName(), Equals, "")
	c.Assert((&session{accountConfig: &config.Account{Nickname: "name 1", Account: "name 4"}, nicknames: []string{"", "name 2", "name 3"}}).DisplayName(), Equals, "name 1")
	c.Assert((&session{accountConfig: &config.Account{Account: "name 4"}, nicknames: []string{"", "name 2", "name 3"}}).DisplayName(), Equals, "name 2")
	c.Assert((&session{accountConfig: &config.Account{Account: "name 4"}, nicknames: []string{""}}).DisplayName(), Equals, "name 4")
}

type requestRosterMockConn struct {
	*mock.Conn
	frr func() (<-chan data.Stanza, data.Cookie, error)
	frd func() (string, error)
}

func (m *requestRosterMockConn) RequestRoster() (<-chan data.Stanza, data.Cookie, error) {
	return m.frr()
}

func (m *requestRosterMockConn) GetRosterDelimiter() (string, error) {
	return m.frd()
}

func (s *SessionSuite) Test_session_requestRoster_works(c *C) {
	mc := &requestRosterMockConn{}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		conn:       mc,
		log:        l,
		connStatus: CONNECTED,
		config:     &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		r: roster.New(),
	}

	mc.frd = func() (string, error) {
		return ":X:", nil
	}

	ch := make(chan data.Stanza, 1)

	mc.frr = func() (<-chan data.Stanza, data.Cookie, error) {
		return ch, 0, nil
	}

	stz := data.Stanza{
		Value: &data.ClientIQ{
			Query: []byte(`<query xmlns='jabber:iq:roster'>
  <item jid='nurse@example.com'/>
  <item jid='romeo@example.net'/>
  <item jid='foo@somewhere.com'/>
  <item jid='abc@example.org'/>
</query>`),
		},
	}

	ch <- stz

	res := sess.requestRoster()

	c.Assert(res, Equals, true)
	c.Assert(sess.groupDelimiter, Equals, ":X:")
	c.Assert(sess.r.ToSlice(), HasLen, 4)
	c.Assert(hook.Entries, HasLen, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Fetching roster")
	c.Assert(hook.Entries[0].Data, HasLen, 0)
	c.Assert(hook.Entries[1].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[1].Message, Equals, "Roster received")
	c.Assert(hook.Entries[1].Data, HasLen, 0)
}

func (s *SessionSuite) Test_session_requestRoster_parsingXMLFails(c *C) {
	mc := &requestRosterMockConn{}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		conn:       mc,
		log:        l,
		connStatus: CONNECTED,
		config:     &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		r: roster.New(),
	}

	mc.frd = func() (string, error) {
		return ":X:", nil
	}

	ch := make(chan data.Stanza, 1)

	mc.frr = func() (<-chan data.Stanza, data.Cookie, error) {
		return ch, 0, nil
	}

	stz := data.Stanza{
		Value: &data.ClientIQ{
			Query: []byte(`<query xmlns='jabber:iq:roster'`),
		},
	}

	ch <- stz

	res := sess.requestRoster()

	c.Assert(res, Equals, true)
	c.Assert(sess.groupDelimiter, Equals, ":X:")
	c.Assert(sess.r.ToSlice(), HasLen, 0)
	c.Assert(hook.Entries, HasLen, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Fetching roster")
	c.Assert(hook.Entries[0].Data, HasLen, 0)
	c.Assert(hook.Entries[1].Level, Equals, log.ErrorLevel)
	c.Assert(hook.Entries[1].Message, Equals, "Failed to parse roster")
	c.Assert(hook.Entries[1].Data, HasLen, 1)
	c.Assert(hook.Entries[1].Data["error"], ErrorMatches, "XML syntax error.*")
}

func (s *SessionSuite) Test_session_requestRoster_failsIfChannelClosed(c *C) {
	mc := &requestRosterMockConn{}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		conn:       mc,
		log:        l,
		connStatus: CONNECTED,
		config:     &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		r: roster.New(),
	}

	mc.frd = func() (string, error) {
		return ":X:", nil
	}

	ch := make(chan data.Stanza, 1)

	mc.frr = func() (<-chan data.Stanza, data.Cookie, error) {
		return ch, 0, nil
	}

	close(ch)

	res := sess.requestRoster()

	c.Assert(res, Equals, true)
	c.Assert(sess.groupDelimiter, Equals, ":X:")
	c.Assert(sess.r.ToSlice(), HasLen, 0)
	c.Assert(hook.Entries, HasLen, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Fetching roster")
	c.Assert(hook.Entries[0].Data, HasLen, 0)
	c.Assert(hook.Entries[1].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[1].Message, Equals, "session: roster request cancelled or timed out")
	c.Assert(hook.Entries[1].Data, HasLen, 0)
}

func (s *SessionSuite) Test_session_requestRoster_failsIfRequestingTheRosterFails(c *C) {
	mc := &requestRosterMockConn{}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		conn:       mc,
		log:        l,
		connStatus: CONNECTED,
		config:     &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		r: roster.New(),
	}

	mc.frd = func() (string, error) {
		return ":X:", nil
	}

	mc.frr = func() (<-chan data.Stanza, data.Cookie, error) {
		return nil, 0, errors.New("this is also a marker")
	}

	res := sess.requestRoster()

	c.Assert(res, Equals, true)
	c.Assert(sess.groupDelimiter, Equals, ":X:")
	c.Assert(sess.r.ToSlice(), HasLen, 0)
	c.Assert(hook.Entries, HasLen, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Fetching roster")
	c.Assert(hook.Entries[0].Data, HasLen, 0)
	c.Assert(hook.Entries[1].Level, Equals, log.ErrorLevel)
	c.Assert(hook.Entries[1].Message, Equals, "Failed to request roster")
	c.Assert(hook.Entries[1].Data, HasLen, 1)
	c.Assert(hook.Entries[1].Data["error"], ErrorMatches, "this is also a marker")
}

func (s *SessionSuite) Test_session_requestRoster_usesDefaultDelimiterWhenFailingToRequestIt(c *C) {
	mc := &requestRosterMockConn{}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		conn:       mc,
		log:        l,
		connStatus: CONNECTED,
		config:     &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		r: roster.New(),
	}

	mc.frd = func() (string, error) {
		return "foo", errors.New("humma")
	}

	mc.frr = func() (<-chan data.Stanza, data.Cookie, error) {
		return nil, 0, errors.New("this is also a marker")
	}

	res := sess.requestRoster()

	c.Assert(res, Equals, true)
	c.Assert(sess.groupDelimiter, Equals, "::")
	c.Assert(sess.r.ToSlice(), HasLen, 0)
	c.Assert(hook.Entries, HasLen, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Fetching roster")
	c.Assert(hook.Entries[0].Data, HasLen, 0)
	c.Assert(hook.Entries[1].Level, Equals, log.ErrorLevel)
	c.Assert(hook.Entries[1].Message, Equals, "Failed to request roster")
	c.Assert(hook.Entries[1].Data, HasLen, 1)
	c.Assert(hook.Entries[1].Data["error"], ErrorMatches, "this is also a marker")
}

func (s *SessionSuite) Test_session_requestRoster_usesDefaultDelimiterWhenEmptyReturned(c *C) {
	mc := &requestRosterMockConn{}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		conn:       mc,
		log:        l,
		connStatus: CONNECTED,
		config:     &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		r: roster.New(),
	}

	mc.frd = func() (string, error) {
		return "", nil
	}

	mc.frr = func() (<-chan data.Stanza, data.Cookie, error) {
		return nil, 0, errors.New("this is also a marker")
	}

	res := sess.requestRoster()

	c.Assert(res, Equals, true)
	c.Assert(sess.groupDelimiter, Equals, "::")
	c.Assert(sess.r.ToSlice(), HasLen, 0)
	c.Assert(hook.Entries, HasLen, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Fetching roster")
	c.Assert(hook.Entries[0].Data, HasLen, 0)
	c.Assert(hook.Entries[1].Level, Equals, log.ErrorLevel)
	c.Assert(hook.Entries[1].Message, Equals, "Failed to request roster")
	c.Assert(hook.Entries[1].Data, HasLen, 1)
	c.Assert(hook.Entries[1].Data["error"], ErrorMatches, "this is also a marker")
}

func (s *SessionSuite) Test_session_requestRoster_returnsWhenNotDisconnected(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log:        l,
		connStatus: DISCONNECTED,
	}

	res := sess.requestRoster()

	c.Assert(res, Equals, false)
	c.Assert(sess.groupDelimiter, Equals, "")
	c.Assert(hook.Entries, HasLen, 0)
}

func (s *SessionSuite) Test_session_watchRoster_works(c *C) {
	orgWaitForNextRosterRequest := waitForNextRosterRequest
	defer func() {
		waitForNextRosterRequest = orgWaitForNextRosterRequest
	}()

	mc := &requestRosterMockConn{}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		conn:       mc,
		log:        l,
		connStatus: CONNECTED,
		config:     &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		r: roster.New(),
	}

	waitForNextRosterRequest = func() {
		sess.connStatus = DISCONNECTED
	}

	mc.frd = func() (string, error) {
		return ":X:", nil
	}

	mc.frr = func() (<-chan data.Stanza, data.Cookie, error) {
		return nil, 0, errors.New("this is also a marker")
	}

	sess.watchRoster()

	c.Assert(sess.connStatus, Equals, DISCONNECTED)
	c.Assert(sess.groupDelimiter, Equals, ":X:")
	c.Assert(sess.r.ToSlice(), HasLen, 0)
	c.Assert(hook.Entries, HasLen, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Fetching roster")
	c.Assert(hook.Entries[0].Data, HasLen, 0)
	c.Assert(hook.Entries[1].Level, Equals, log.ErrorLevel)
	c.Assert(hook.Entries[1].Message, Equals, "Failed to request roster")
	c.Assert(hook.Entries[1].Data, HasLen, 1)
	c.Assert(hook.Entries[1].Data["error"], ErrorMatches, "this is also a marker")
}

type mockDialer struct {
	*mock.Dialer

	f func() (interfaces.Conn, error)
}

func (m *mockDialer) Dial() (interfaces.Conn, error) {
	if m.f != nil {
		return m.f()
	}
	return nil, nil
}

type mockConnectConn struct {
	*mock.Conn

	getResource      func() string
	serverHasFeature func(string) bool
}

func (m *mockConnectConn) GetJIDResource() string {
	if m.getResource != nil {
		return m.getResource()
	}
	return ""
}

func (m *mockConnectConn) ServerHasFeature(v string) bool {
	if m.serverHasFeature != nil {
		return m.serverHasFeature(v)
	}
	return false
}

func (s *SessionSuite) Test_session_Connect_works(c *C) {
	md := &mockDialer{}
	mc := &mockConnectConn{
		getResource:      func() string { return "hoho" },
		serverHasFeature: func(string) bool { return true },
	}
	md.f = func() (interfaces.Conn, error) { return mc, nil }

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		dialerFactory: func(tls.Verifier, tls.Factory) interfaces.Dialer { return md },
		log:           l,
		connStatus:    DISCONNECTED,
		config:        &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		wantToBeOnline: true,
		resource:       "somewhere",
		r:              roster.New(),
	}

	res := sess.Connect("one", nil)

	c.Assert(res, IsNil)
	c.Assert(sess.resource, Equals, "hoho")
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Connect()")
	c.Assert(hook.Entries[0].Data, HasLen, 2)
	c.Assert(hook.Entries[0].Data["resource"], Equals, "somewhere")
	c.Assert(hook.Entries[0].Data["wantToBeOnline"], Equals, true)
}

func (s *SessionSuite) Test_session_Connect_worksWithoutVCard(c *C) {
	md := &mockDialer{}
	mc := &mockConnectConn{
		getResource:      func() string { return "hoho" },
		serverHasFeature: func(string) bool { return false },
	}
	md.f = func() (interfaces.Conn, error) { return mc, nil }

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		dialerFactory: func(tls.Verifier, tls.Factory) interfaces.Dialer { return md },
		log:           l,
		connStatus:    DISCONNECTED,
		config:        &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		wantToBeOnline: false,
		resource:       "somewhere",
		r:              roster.New(),
	}

	res := sess.Connect("one", nil)

	c.Assert(res, IsNil)
	c.Assert(sess.resource, Equals, "hoho")
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Connect()")
	c.Assert(hook.Entries[0].Data, HasLen, 2)
	c.Assert(hook.Entries[0].Data["resource"], Equals, "")
	c.Assert(hook.Entries[0].Data["wantToBeOnline"], Equals, false)
}

func (s *SessionSuite) Test_session_Connect_closesIfWeChangeConnStatus(c *C) {
	md := &mockDialer{}
	mc := &mockConnectConn{
		getResource:      func() string { return "hoho" },
		serverHasFeature: func(string) bool { return true },
	}
	md.f = func() (interfaces.Conn, error) { return mc, nil }

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	var sess *session
	sess = &session{
		conn: mc,
		dialerFactory: func(tls.Verifier, tls.Factory) interfaces.Dialer {
			sess.connStatus = DISCONNECTED
			return md
		},
		log:        l,
		connStatus: DISCONNECTED,
		config:     &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		wantToBeOnline: true,
		resource:       "somewhere",
		r:              roster.New(),
	}

	res := sess.Connect("one", nil)

	c.Assert(res, IsNil)
	c.Assert(sess.resource, Equals, "somewhere")
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Connect()")
	c.Assert(hook.Entries[0].Data, HasLen, 2)
	c.Assert(hook.Entries[0].Data["resource"], Equals, "somewhere")
	c.Assert(hook.Entries[0].Data["wantToBeOnline"], Equals, true)
}

func (s *SessionSuite) Test_session_Connect_failsOnConnectionFailure(c *C) {
	md := &mockDialer{}
	md.f = func() (interfaces.Conn, error) { return nil, errors.New("dialer marker failure") }

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		dialerFactory: func(tls.Verifier, tls.Factory) interfaces.Dialer { return md },
		log:           l,
		connStatus:    DISCONNECTED,
		config:        &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		wantToBeOnline: true,
		resource:       "somewhere",
		r:              roster.New(),
	}

	res := sess.Connect("one", nil)

	c.Assert(res, ErrorMatches, "dialer marker failure")
	c.Assert(sess.connStatus, Equals, DISCONNECTED)
	c.Assert(hook.Entries, HasLen, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Connect()")
	c.Assert(hook.Entries[0].Data, HasLen, 2)
	c.Assert(hook.Entries[0].Data["resource"], Equals, "somewhere")
	c.Assert(hook.Entries[0].Data["wantToBeOnline"], Equals, true)
	c.Assert(hook.Entries[1].Level, Equals, log.ErrorLevel)
	c.Assert(hook.Entries[1].Message, Equals, "failed to connect")
	c.Assert(hook.Entries[1].Data, HasLen, 1)
	c.Assert(hook.Entries[1].Data["error"], ErrorMatches, "dialer marker failure")
}

func (s *SessionSuite) Test_session_Connect_doesntDoAnythingWhenConnected(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	sess := &session{
		log:        l,
		connStatus: CONNECTED,
		config:     &config.ApplicationConfig{},
		accountConfig: &config.Account{
			Account: "some@one.org",
		},
		wantToBeOnline: true,
		resource:       "somewhere",
		r:              roster.New(),
	}

	res := sess.Connect("one", nil)

	c.Assert(res, IsNil)
	c.Assert(hook.Entries, HasLen, 0)
}

func (s *SessionSuite) Test_session_Close_doesntDoAnythingIfDisconnected(c *C) {
	sess := &session{
		connStatus: DISCONNECTED,
	}

	sess.Close()
}

func (s *SessionSuite) Test_session_Close_setsTheConnectionDisconnectedIfNoConnectionExists(c *C) {
	sess := &session{
		conn:       nil,
		connStatus: CONNECTED,
	}

	sess.Close()
	c.Assert(sess.connStatus, Equals, DISCONNECTED)
}

type closeMockConn struct {
	*mock.Conn
	f func() error
}

func (m *closeMockConn) Close() error {
	if m.f != nil {
		return m.f()
	}
	return nil
}

func (s *SessionSuite) Test_session_Close_closesTheConnection(c *C) {
	called := false
	cn := &closeMockConn{
		f: func() error {
			called = true
			return nil
		},
	}

	sess := &session{
		conn:           cn,
		connStatus:     CONNECTED,
		wantToBeOnline: true,
	}

	sess.Close()
	c.Assert(sess.conn, IsNil)
	c.Assert(sess.connStatus, Equals, DISCONNECTED)
	c.Assert(called, Equals, true)
}

func (s *SessionSuite) Test_session_Close_closesTheConnectionAndTerminatesConnections(c *C) {
	cn := &closeMockConn{}
	cm := &mockConvManager{}

	sess := &session{
		conn:           cn,
		connStatus:     CONNECTED,
		wantToBeOnline: false,
		convManager:    cm,
	}
	called := false
	cm.terminateAll = func() {
		called = true
	}

	sess.Close()
	c.Assert(sess.conn, IsNil)
	c.Assert(sess.connStatus, Equals, DISCONNECTED)
	c.Assert(called, Equals, true)
}

func (s *SessionSuite) Test_session_connectionLost_closes(c *C) {
	called := false
	cn := &closeMockConn{
		f: func() error {
			called = true
			return nil
		},
	}

	sess := &session{
		conn:           cn,
		connStatus:     CONNECTED,
		wantToBeOnline: true,
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.connectionLost()
	c.Assert(sess.conn, IsNil)
	c.Assert(sess.connStatus, Equals, DISCONNECTED)
	c.Assert(called, Equals, true)

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Event)
		if !ok || t.Type != events.ConnectionLost {
			return false
		}

		c.Assert(t.Type, Equals, events.ConnectionLost)
		return true
	})
}

func (s *SessionSuite) Test_session_EncryptAndSendTo_returnsErrorWhenOffline(c *C) {
	sess := &session{
		connStatus: DISCONNECTED,
	}

	trace, delay, e := sess.EncryptAndSendTo(jid.Parse("someone@foo.org"), "hello")
	c.Assert(trace, Equals, 0)
	c.Assert(delay, Equals, false)
	c.Assert(e, ErrorMatches, "Couldn't send message since we are not connected")
}

func (s *SessionSuite) Test_session_EncryptAndSendTo_sendsMessageThroughOTR(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	eh := &otrclient.EventHandler{
		Log: l,
	}
	mcm := &mockConvManager{}

	sess := &session{
		connStatus:  CONNECTED,
		convManager: mcm,
		config:      &config.ApplicationConfig{},
	}

	mc := &mockConv{}
	mc.eh = eh
	vals := []byte{}
	mc.send = func(v []byte) (int, error) {
		vals = v
		return 42, errors.New("marker error")
	}
	mcm.ensureConversationWith = func(jid.Any, []byte) (otrclient.Conversation, bool) {
		return mc, false
	}

	trace, delayed, e := sess.EncryptAndSendTo(jid.Parse("some@two.org/bla"), "allo over there")
	c.Assert(string(vals), Equals, "allo over there")
	c.Assert(trace, Equals, 42)
	c.Assert(delayed, Equals, false)
	c.Assert(e, ErrorMatches, "marker error")
	c.Assert(hook.Entries, HasLen, 0)
}
