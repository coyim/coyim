package session

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/otrclient"
	"github.com/coyim/coyim/roster"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/glib_mock"

	. "gopkg.in/check.v1"
)

func init() {
	log.SetOutput(ioutil.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
}

func Test(t *testing.T) { TestingT(t) }

type SessionSuite struct{}

var _ = Suite(&SessionSuite{})

func (s *SessionSuite) Test_NewSession_returnsANewSession(c *C) {
	sess := Factory(&config.ApplicationConfig{}, &config.Account{}, xmpp.DialerFactory)
	c.Assert(sess, Not(IsNil))
}

const testTimeout = time.Duration(5) * time.Second

func (s *SessionSuite) Test_info_publishesInfoEvent(c *C) {
	sess := &session{}

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 1)
	sess.eventsReachedZero = eventsDone

	sess.info("hello world")

	select {
	case <-eventsDone:
		close(observer)
		select {
		case ev := <-observer:
			t := ev.(events.Log)
			c.Assert(t.Level, Equals, events.Info)
			c.Assert(t.Message, Equals, "hello world")
		default:
			c.Errorf("did not receive event")
		}
	case <-time.After(testTimeout):
		c.Errorf("test timed out")
	}
}

func (s *SessionSuite) Test_warn_publishesWarnEvent(c *C) {
	sess := &session{}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 1)
	sess.eventsReachedZero = eventsDone

	sess.warn("hello world2")

	select {
	case <-eventsDone:
		close(observer)
		select {
		case ev := <-observer:
			t := ev.(events.Log)
			c.Assert(t.Level, Equals, events.Warn)
			c.Assert(t.Message, Equals, "hello world2")
		default:
			c.Errorf("did not receive event")
		}
	case <-time.After(testTimeout):
		c.Errorf("test timed out")
	}
}

func (s *SessionSuite) Test_alert_publishedAlertEvent(c *C) {
	sess := &session{}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 1)
	sess.eventsReachedZero = eventsDone

	sess.alert("hello world3")

	select {
	case <-eventsDone:
		close(observer)
		select {
		case ev := <-observer:
			t := ev.(events.Log)
			c.Assert(t.Level, Equals, events.Alert)
			c.Assert(t.Message, Equals, "hello world3")
		default:
			c.Errorf("did not receive event")
		}
	case <-time.After(testTimeout):
		c.Errorf("test timed out")
	}
}

func (s *SessionSuite) Test_iqReceived_publishesIQReceivedEvent(c *C) {
	sess := &session{}

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
				Session: sess,
				Type:    events.IQReceived,
				From:    jid.NR("someone@somewhere"),
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

	sess := &session{
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 1)
	sess.eventsReachedZero = eventsDone

	sess.watchStanzas()

	select {
	case <-eventsDone:
		close(observer)
		select {
		case ev := <-observer:
			t := ev.(events.Log)
			c.Assert(t.Message, Equals, "error reading XMPP message: unexpected XMPP message clientx <message/>")
		default:
			c.Errorf("did not receive event")
		}
	case <-time.After(testTimeout):
		c.Errorf("test timed out")
	}
}

func (s *SessionSuite) Test_WatchStanzas_handlesUnknownMessage(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<bind:bind xmlns:bind='urn:ietf:params:xml:ns:xmpp-bind'></bind:bind>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.watchStanzas()

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Log)
		if !ok || t.Level != events.Info {
			return false
		}
		c.Assert(t.Message, Equals, "unhandled stanza: {urn:ietf:params:xml:ns:xmpp-bind bind} &{{urn:ietf:params:xml:ns:xmpp-bind bind}  }")
		return true
	})
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
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)

	sess.watchStanzas()

	assertLogContains(c, observer, events.Log{
		Level:   events.Alert,
		Message: "Exiting in response to fatal error from server: bad horse showed up",
	})
}

func (s *SessionSuite) Test_WatchStanzas_handlesStreamError_withEmbeddedTag(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:error xmlns:stream='http://etherx.jabber.org/streams'><not-well-formed xmlns='urn:ietf:params:xml:ns:xmpp-streams'/></stream:error>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)

	sess.watchStanzas()

	assertLogContains(c, observer, events.Log{
		Level:   events.Alert,
		Message: "Exiting in response to fatal error from server: {urn:ietf:params:xml:ns:xmpp-streams not-well-formed}",
	})
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
		c.Assert(t.Session, Equals, sess)
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

	sess := &session{
		connStatus: DISCONNECTED,
	}
	sess.conn = conn

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.watchStanzas()

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Log)
		if !ok || t.Level != events.Info {
			return false
		}

		c.Assert(t.Message, Equals, "unrecognized iq: &data.ClientIQ{XMLName:xml.Name{Space:\"jabber:client\", Local:\"iq\"}, From:\"\", ID:\"\", To:\"\", Type:\"something\", Error:data.StanzaError{By:\"\", Code:\"\", Type:\"\", Text:\"\", Condition:struct { XMLName xml.Name; Body string \"xml:\\\",innerxml\\\"\" }{XMLName:xml.Name{Space:\"\", Local:\"\"}, Body:\"\"}, ApplicationCondition:(*data.Any)(nil)}, Bind:data.BindBind{XMLName:xml.Name{Space:\"\", Local:\"\"}, Resource:\"\", Jid:\"\"}, Query:[]uint8{}}")
		return true
	})
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
		"<iq to='abc' from='some@one.org/foo' type='result' id=''>"+
		"<query xmlns=\"http://jabber.org/protocol/disco#info\">"+
		"<identity xmlns=\"http://jabber.org/protocol/disco#info\" category=\"client\" type=\"pc\" name=\"foo.bar@somewhere.org\"></identity>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/disco#info\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"urn:xmpp:bob\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"urn:xmpp:ping\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/caps\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"jabber:iq:version\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"vcard-temp\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"jabber:x:data\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/si\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/si/profile/file-transfer\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/si/profile/directory-transfer\"></feature>"+
		//		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/si/profile/encrypted-data-transfer\"></feature>"+
		"<feature xmlns=\"http://jabber.org/protocol/disco#info\" var=\"http://jabber.org/protocol/bytestreams\"></feature>"+
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
		"<iq to='abc' from='some@one.org/foo' type='result' id=''>"+
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

	sess := &session{
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

	for {
		select {
		case ev := <-observer:
			t := ev.(events.Log)
			if t.Level != events.Info {
				continue
			}

			c.Assert(t.Message, Equals, "Unknown IQ: <query xmlns='jabber:iq:somethingStrange'/>")
			return
		case <-time.After(time.Duration(10) * time.Second):
			c.Errorf("did not receive event")
			return
		}
	}
}

func (s *SessionSuite) Test_WatchStanzas_iq_set_roster_withBadFrom(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='set' from='some2@one.org' to='cde'><query xmlns='jabber:iq:roster'/></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
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

	assertLogContains(c, observer, events.Log{
		Level:   events.Warn,
		Message: "Ignoring roster IQ from bad address: some2@one.org",
	})

	c.Assert(string(mockIn.write), Equals, "")
}

func (s *SessionSuite) Test_WatchStanzas_iq_set_roster_withFromContainingJid(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='set' from='some@one.org/foo' to='cde'><query xmlns='jabber:iq:roster'/></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &session{
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

	assertLogContains(c, observer, events.Log{
		Level:   events.Warn,
		Message: "Failed to parse roster push IQ",
	})
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
		config:        &config.ApplicationConfig{},
		accountConfig: &config.Account{},
		r:             roster.New(),
		connStatus:    DISCONNECTED,
	}
	sess.conn = conn
	sess.r.AddOrReplace(roster.PeerWithState(jid.NR("some2@one.org"), "somewhere", "", "", "balcony"))

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
		config:        &config.ApplicationConfig{},
		accountConfig: &config.Account{},
		r:             roster.New(),
		connStatus:    DISCONNECTED,
	}
	sess.conn = conn
	sess.r.AddOrReplace(roster.PeerWithState(jid.NR("some2@one.org"), "dnd", "", "", "main"))

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
	sess := &session{
		r: roster.New(),
	}

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.HandleConfirmOrDeny(jid.NR("foo@bar.com"), true)

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Log)
		if !ok {
			return false
		}

		c.Assert(t.Level, Equals, events.Warn)
		return true
	})
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

	sess := &session{
		r: roster.New(),
	}
	sess.conn = conn
	sess.r.SubscribeRequest(jid.NR("foo@bar.com"), "123", "")

	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 2)
	sess.eventsReachedZero = eventsDone

	sess.HandleConfirmOrDeny(jid.NR("foo@bar.com"), true)

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Log)
		if !ok || t.Level != events.Warn {
			return false
		}

		c.Assert(t.Message, Equals, "Error sending presence stanza: foo bar")
		return true
	})
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
	ensureConversationWith func(jid.Any) (otrclient.Conversation, bool)
	terminateAll           func()
}

func (mcm *mockConvManager) GetConversationWith(peer jid.Any) (otrclient.Conversation, bool) {
	return mcm.getConversationWith(peer)
}

func (mcm *mockConvManager) EnsureConversationWith(peer jid.Any) (otrclient.Conversation, bool) {
	return mcm.ensureConversationWith(peer)
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

// func otrEventHandlerWith(s string, eh *otrclient.EventHandler) *otrclient.EventHandlers {
// 	ehs := otrclient.NewEventHandlers("one", func(jid.Any, *otrclient.EventHandler, chan string, chan int) {})
// 	ehs.Add(jid.Parse(s), eh)
// 	return ehs
// }

func (s *SessionSuite) Test_receiveClientMessage_willNotProcessBRTagsWhenNotEncrypted(c *C) {
	mcm := &mockConvManager{}
	sess := &session{
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

	mcm.ensureConversationWith = func(jid.Any) (otrclient.Conversation, bool) {
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
		connStatus:  CONNECTED,
		convManager: mcm,
		config:      &config.ApplicationConfig{},
	}

	mc := &mockConv{}
	mc.receive = func(s []byte) ([]byte, error) { return s, nil }
	mc.isEncrypted = func() bool { return true }
	mcm.ensureConversationWith = func(jid.Any) (otrclient.Conversation, bool) {
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

func (ncm *convManagerWithoutConversation) EnsureConversationWith(peer jid.Any) (otrclient.Conversation, bool) {
	return nil, false
}

func (ncm *convManagerWithoutConversation) TerminateAll() {
}

func sessionWithConvMngrWithoutConvs() *session {
	return &session{
		connStatus:  CONNECTED,
		convManager: &convManagerWithoutConversation{},
		config:      &config.ApplicationConfig{},
	}
}

func (s *SessionSuite) Test_logsError_whenWeStartSMPWithoutAConversation(c *C) {
	sess := sessionWithConvMngrWithoutConvs()
	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 1)
	sess.eventsReachedZero = eventsDone

	go sess.StartSMP(jid.R("someone's peer/resource"), "Im a question", "im an answer")

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Log)
		if !ok || t.Level != events.Alert {
			return false
		}

		c.Assert(t.Level, Equals, events.Alert)
		c.Assert(t.Message, Equals, "error: tried to start SMP when a conversation does not exist")
		return true
	})
}

func (s *SessionSuite) Test_logsError_whenWeFinishSMPWithoutAConversation(c *C) {
	sess := sessionWithConvMngrWithoutConvs()
	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 1)
	sess.eventsReachedZero = eventsDone

	go sess.FinishSMP(jid.R("someone's peer/resource"), "im an answer")

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Log)
		if !ok || t.Level != events.Alert {
			return false
		}

		c.Assert(t.Level, Equals, events.Alert)
		c.Assert(t.Message, Equals, "error: tried to finish SMP when a conversation does not exist")
		return true
	})
}

func (s *SessionSuite) Test_logsError_whenWeAbortSMPWithoutAConversation(c *C) {
	sess := sessionWithConvMngrWithoutConvs()
	observer := make(chan interface{}, 1000)
	sess.Subscribe(observer)
	eventsDone := make(chan bool, 1)
	sess.eventsReachedZero = eventsDone

	go sess.AbortSMP(jid.R("someone's peer/resource"))

	assertReceivesEvent(c, eventsDone, observer, func(ev interface{}) bool {
		t, ok := ev.(events.Log)
		if !ok || t.Level != events.Alert {
			return false
		}

		c.Assert(t.Level, Equals, events.Alert)
		c.Assert(t.Message, Equals, "error: tried to abort SMP when a conversation does not exist")
		return true
	})
}
