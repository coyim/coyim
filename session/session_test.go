package session

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/event"
	"github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/otr3"

	. "gopkg.in/check.v1"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func Test(t *testing.T) { TestingT(t) }

type SessionXmppSuite struct{}

var _ = Suite(&SessionXmppSuite{})

func (s *SessionXmppSuite) Test_NewSession_returnsANewSession(c *C) {
	sess := NewSession(&config.Accounts{}, &config.Account{})
	c.Assert(sess, Not(IsNil))
}

func (s *SessionXmppSuite) Test_info_publishesInfoEvent(c *C) {
	sess := &Session{}

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.info("hello world")

	select {
	case ev := <-observer:
		t := ev.(LogEvent)
		c.Assert(t.Level, Equals, Info)
		c.Assert(t.Message, Equals, "hello world")
	case <-time.After(1 * time.Millisecond):
		c.Errorf("did not receive event")
	}
}

func (s *SessionXmppSuite) Test_warn_publishesWarnEvent(c *C) {
	sess := &Session{}

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.warn("hello world2")

	select {
	case ev := <-observer:
		t := ev.(LogEvent)
		c.Assert(t.Level, Equals, Warn)
		c.Assert(t.Message, Equals, "hello world2")
	case <-time.After(1 * time.Millisecond):
		c.Errorf("did not receive event")
	}
}

func (s *SessionXmppSuite) Test_alert_publishedAlertEvent(c *C) {
	sess := &Session{}

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.alert("hello world3")

	select {
	case ev := <-observer:
		t := ev.(LogEvent)
		c.Assert(t.Level, Equals, Alert)
		c.Assert(t.Message, Equals, "hello world3")
	case <-time.After(1 * time.Millisecond):
		c.Errorf("did not receive event")
	}
}

func (s *SessionXmppSuite) Test_iqReceived_publishesIQReceivedEvent(c *C) {
	sess := &Session{}

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.iqReceived("someone@somewhere")

	select {
	case ev := <-observer:
		c.Assert(ev, Equals, PeerEvent{
			Session: sess,
			Type:    IQReceived,
			From:    "someone@somewhere",
		})
	case <-time.After(1 * time.Millisecond):
		c.Error("did not receive event")
	}
}

func (s *SessionXmppSuite) Test_readMessages_passesStanzaToChannel(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:message xmlns:client='jabber:client' to='fo@bar.com' from='bar@foo.com' type='chat'><client:body>something</client:body></client:message>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		nil,
		"some@one.org/foo",
	)

	sess := &Session{}
	sess.Conn = conn

	stanzaChan := make(chan xmpp.Stanza)
	go sess.readMessages(stanzaChan)

	select {
	case rawStanza, ok := <-stanzaChan:
		c.Assert(ok, Equals, true)
		c.Assert(rawStanza.Name.Local, Equals, "message")
		c.Assert(rawStanza.Value.(*xmpp.ClientMessage).Body, Equals, "something")
	}
}

func (s *SessionXmppSuite) Test_readMessages_alertsOnError(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<clientx:message xmlns:client='jabber:client' to='fo@bar.com' from='bar@foo.com' type='chat'><client:body>something</client:body></client:message>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		nil,
		"some@one.org/foo",
	)

	sess := &Session{}
	sess.Conn = conn

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	stanzaChan := make(chan xmpp.Stanza)
	go sess.readMessages(stanzaChan)

	select {
	case _, ok := <-stanzaChan:
		c.Assert(ok, Equals, false)
	}

	select {
	case ev := <-observer:
		t := ev.(LogEvent)
		c.Assert(t.Level, Equals, Alert)
		c.Assert(t.Message, Equals, "unexpected XMPP message clientx <message/>")
	case <-time.After(1 * time.Millisecond):
		c.Errorf("did not receive event")
	}
}

func (s *SessionXmppSuite) Test_WatchStanzas_warnsAndExitsOnBadStanza(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<clientx:message xmlns:client='jabber:client' to='fo@bar.com' from='bar@foo.com' type='chat'><client:body>something</client:body></client:message>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		nil,
		"some@one.org/foo",
	)

	sess := &Session{
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	for i := 0; i < 2; i++ {
		select {
		case ev := <-observer:
			t := ev.(LogEvent)
			switch t.Level {
			case Warn:
				c.Assert(t.Message, Equals, "Exiting because channel to server closed")
			case Alert:
				c.Assert(t.Message, Equals, "unexpected XMPP message clientx <message/>")
			}
		case <-time.After(1 * time.Millisecond):
			c.Errorf("did not receive event")
		}
	}
}

func (s *SessionXmppSuite) Test_WatchStanzas_handlesUnknownMessage(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<bind:bind xmlns:bind='urn:ietf:params:xml:ns:xmpp-bind'></bind:bind>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		nil,
		"some@one.org/foo",
	)

	sess := &Session{
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	for {
		select {
		case ev := <-observer:
			t := ev.(LogEvent)
			if t.Level != Info {
				continue
			}

			c.Assert(t.Message, Equals, "{urn:ietf:params:xml:ns:xmpp-bind bind} &{{urn:ietf:params:xml:ns:xmpp-bind bind}  }")
			return

		case <-time.After(1 * time.Millisecond):
			c.Errorf("did not receive event")
			return
		}
	}
}

func (s *SessionXmppSuite) Test_WatchStanzas_handlesStreamError_withText(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:error xmlns:stream='http://etherx.jabber.org/streams'><stream:text>bad horse showed up</stream:text></stream:error>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		nil,
		"some@one.org/foo",
	)

	sess := &Session{
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	assertLogContains(c, observer, LogEvent{
		Level:   Alert,
		Message: "Exiting in response to fatal error from server: bad horse showed up",
	})
}

func (s *SessionXmppSuite) Test_WatchStanzas_handlesStreamError_withEmbeddedTag(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:error xmlns:stream='http://etherx.jabber.org/streams'><not-well-formed xmlns='urn:ietf:params:xml:ns:xmpp-streams'/></stream:error>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		nil,
		"some@one.org/foo",
	)

	sess := &Session{
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	for i := 0; i < 2; i++ {
		select {
		case ev := <-observer:
			t := ev.(LogEvent)
			if i < 1 {
				continue
			}

			c.Assert(t.Level, Equals, Alert)
			c.Assert(t.Message, Equals, "Exiting in response to fatal error from server: {urn:ietf:params:xml:ns:xmpp-streams not-well-formed}")
			return

		case <-time.After(1 * time.Millisecond):
			c.Errorf("did not receive event")
			return
		}
	}
}

func (s *SessionXmppSuite) Test_WatchStanzas_receivesAMessage(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:message xmlns:client='jabber:client' type='chat' to='some@one.org/foo' from='bla@hmm.org/somewhere'><client:body>well, hello there</client:body></client:message>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		Config:          &config.Accounts{},
		CurrentAccount:  &config.Account{InstanceTag: uint32(42)},
		Conversations:   make(map[string]*otr3.Conversation),
		OtrEventHandler: make(map[string]*event.OtrEventHandler),
		ConnStatus:      DISCONNECTED,
	}
	sess.Conn = conn

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	for {
		select {
		case ev := <-observer:
			switch t := ev.(type) {
			case MessageEvent:
				c.Assert(t.Session, Equals, sess)
				c.Assert(t.Encrypted, Equals, false)
				c.Assert(t.From, Equals, "bla@hmm.org")
				c.Assert(string(t.Body), Equals, "well, hello there")
				return
			default:
				//ignore
			}
		case <-time.After(1 * time.Millisecond):
			c.Errorf("did not receive event")
			return
		}
	}
}

func (s *SessionXmppSuite) Test_WatchStanzas_failsOnUnrecognizedIQ(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='something'></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	var sess *Session
	sess = &Session{
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	for {
		select {
		case ev := <-observer:
			t := ev.(LogEvent)
			if t.Level != Info {
				continue
			}

			c.Assert(t.Message, Equals, "unrecognized iq: &xmpp.ClientIQ{XMLName:xml.Name{Space:\"jabber:client\", Local:\"iq\"}, From:\"\", ID:\"\", To:\"\", Type:\"something\", Error:xmpp.ClientError{XMLName:xml.Name{Space:\"\", Local:\"\"}, Code:\"\", Type:\"\", Any:xml.Name{Space:\"\", Local:\"\"}, Text:\"\"}, Bind:xmpp.bindBind{XMLName:xml.Name{Space:\"\", Local:\"\"}, Resource:\"\", Jid:\"\"}, Query:[]uint8{}}")
			return

		case <-time.After(1 * time.Millisecond):
			c.Errorf("did not receive event")
			return
		}
	}
}

func (s *SessionXmppSuite) Test_WatchStanzas_getsDiscoInfoIQ(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='get' from='abc' to='cde'><query xmlns='http://jabber.org/protocol/disco#info'/></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		Config: &config.Accounts{},
		CurrentAccount: &config.Account{
			Account: "foo.bar@somewhere.org",
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(string(mockIn.write), Equals, ""+
		"<iq to='abc' from='some@one.org/foo' type='result' id=''>"+
		"<query xmlns=\"http://jabber.org/protocol/disco#info\">"+
		"<node></node>"+
		"<identity xmlns=\"http://jabber.org/protocol/disco#info\" category=\"client\" type=\"pc\" name=\"foo.bar@somewhere.org\"></identity>"+
		"</query>"+
		"</iq>")
}

func (s *SessionXmppSuite) Test_WatchStanzas_getsVersionInfoIQ(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='get' from='abc' to='cde'><query xmlns='jabber:iq:version'/></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		Config: &config.Accounts{},
		CurrentAccount: &config.Account{
			Account: "foo.bar@somewhere.org",
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(string(mockIn.write), Equals, ""+
		"<iq to='abc' from='some@one.org/foo' type='result' id=''>"+
		"<query xmlns=\"jabber:iq:version\">"+
		"<name>testing</name>"+
		"<version>version</version>"+
		"<os>none</os>"+
		"</query>"+
		"</iq>")
}

func (s *SessionXmppSuite) Test_WatchStanzas_getsUnknown(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='get' from='abc' to='cde'><query xmlns='jabber:iq:somethingStrange'/></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		Config: &config.Accounts{},
		CurrentAccount: &config.Account{
			Account: "foo.bar@somewhere.org",
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	for {
		select {
		case ev := <-observer:
			t := ev.(LogEvent)
			if t.Level != Info {
				continue
			}

			c.Assert(t.Message, Equals, "Unknown IQ: jabber:iq:somethingStrange query")
			return

		case <-time.After(1 * time.Millisecond):
			c.Errorf("did not receive event")
			return
		}
	}
}

func (s *SessionXmppSuite) Test_WatchStanzas_iq_set_roster_withBadFrom(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='get' from='some2@one.org' to='cde'><query xmlns='jabber:iq:roster'/></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		Config: &config.Accounts{},
		CurrentAccount: &config.Account{
			Account: "some@one.org",
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	assertLogContains(c, observer, LogEvent{
		Level:   Warn,
		Message: "Ignoring roster IQ from bad address: some2@one.org",
	})

	// TODO: this is actually incorrect - something like this should be ignored, not responded to
	c.Assert(string(mockIn.write), Equals, "<iq to='some2@one.org' from='some@one.org/foo' type='result' id=''><error type=\"cancel\"><bad-request xmlns=\"urn:ietf:params:xml:ns:xmpp-stanzas\"></bad-request></error></iq>")
}

func (s *SessionXmppSuite) Test_WatchStanzas_iq_set_roster_withFromContainingJid(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='get' from='some@one.org/foo' to='cde'><query xmlns='jabber:iq:roster'/></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		Config: &config.Accounts{},
		CurrentAccount: &config.Account{
			Account: "some@one.org",
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	assertLogContains(c, observer, LogEvent{
		Level:   Warn,
		Message: "Failed to parse roster push IQ",
	})
}

func (s *SessionXmppSuite) Test_WatchStanzas_iq_set_roster_addsANewRosterItem(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='get' to='cde'><query xmlns='jabber:iq:roster'>" +
		"<item jid='romeo@example.net' name='Romeo' subscription='both'>" +
		"<group>Friends</group>" +
		"</item>" +
		"</query></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		Config: &config.Accounts{},
		CurrentAccount: &config.Account{
			Account: "some@one.org",
		},
		R:          roster.New(),
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(sess.R.ToSlice(), DeepEquals, []*roster.Peer{
		roster.PeerFrom(xmpp.RosterEntry{Jid: "romeo@example.net", Subscription: "both", Name: "Romeo", Group: []string{"Friends"}}, sess.CurrentAccount.ID())})
}

func (s *SessionXmppSuite) Test_WatchStanzas_iq_set_roster_setsExistingRosterItem(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='get' to='cde'><query xmlns='jabber:iq:roster'>" +
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

	sess := &Session{
		Config: &config.Accounts{},
		CurrentAccount: &config.Account{
			Account: "some@one.org",
		},
		R:          roster.New(),
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.R.AddOrReplace(roster.PeerFrom(xmpp.RosterEntry{Jid: "jill@example.net", Subscription: "both", Name: "Jill", Group: []string{"Foes"}}, sess.CurrentAccount.ID()))
	sess.R.AddOrReplace(roster.PeerFrom(xmpp.RosterEntry{Jid: "romeo@example.net", Subscription: "both", Name: "Mo", Group: []string{"Foes"}}, sess.CurrentAccount.ID()))

	sess.WatchStanzas()

	c.Assert(called, Equals, 0)
	c.Assert(sess.R.ToSlice(), DeepEquals, []*roster.Peer{
		roster.PeerFrom(xmpp.RosterEntry{Jid: "jill@example.net", Subscription: "both", Name: "Jill", Group: []string{"Foes"}}, sess.CurrentAccount.ID()),
		roster.PeerFrom(xmpp.RosterEntry{Jid: "romeo@example.net", Subscription: "both", Name: "Romeo", Group: []string{"Friends"}}, sess.CurrentAccount.ID()),
	})
}

func (s *SessionXmppSuite) Test_WatchStanzas_iq_set_roster_removesRosterItems(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='get' to='cde'><query xmlns='jabber:iq:roster'>" +
		"<item jid='romeo@example.net' name='Romeo' subscription='remove'>" +
		"<group>Friends</group>" +
		"</item>" +
		"</query></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		Config: &config.Accounts{},
		CurrentAccount: &config.Account{
			Account: "some@one.org",
		},
		R:          roster.New(),
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.R.AddOrReplace(roster.PeerFrom(xmpp.RosterEntry{Jid: "romeo@example.net", Subscription: "both", Name: "Mo", Group: []string{"Foes"}}, ""))
	sess.R.AddOrReplace(roster.PeerFrom(xmpp.RosterEntry{Jid: "jill@example.net", Subscription: "both", Name: "Jill", Group: []string{"Foes"}}, ""))
	sess.R.AddOrReplace(roster.PeerFrom(xmpp.RosterEntry{Jid: "romeo@example.net", Subscription: "both", Name: "Mo", Group: []string{"Foes"}}, ""))

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	c.Assert(sess.R.ToSlice(), DeepEquals, []*roster.Peer{
		roster.PeerFrom(xmpp.RosterEntry{Jid: "jill@example.net", Subscription: "both", Name: "Jill", Group: []string{"Foes"}}, ""),
	})

	select {
	case ev := <-observer:
		switch ev.(type) {
		case PeerEvent:
			c.Error("Received peer event")
			return
		default:
			// ignore
		}
	case <-time.After(1 * time.Millisecond):
		return
	}
}

func (s *SessionXmppSuite) Test_WatchStanzas_presence_unavailable_forNoneKnownUser(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org' type='unavailable'><client:status>going on vacation</client:status></client:presence>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		R:          roster.New(),
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	select {
	case ev := <-observer:
		switch ev.(type) {
		case PresenceEvent:
			c.Error("Received presence event")
			return
		default:
			// ignore
		}
	case <-time.After(1 * time.Millisecond):
		return
	}
}

func (s *SessionXmppSuite) Test_WatchStanzas_presence_unavailable_forKnownUser(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org' type='unavailable'><client:status>going on vacation</client:status></client:presence>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		Config:         &config.Accounts{},
		CurrentAccount: &config.Account{},
		R:              roster.New(),
		ConnStatus:     DISCONNECTED,
	}
	sess.Conn = conn
	sess.R.AddOrReplace(roster.PeerWithState("some2@one.org", "somewhere", "", ""))

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	p, _ := sess.R.Get("some2@one.org")
	c.Assert(p.Online, Equals, false)

	for {
		select {
		case ev := <-observer:
			switch t := ev.(type) {
			case PresenceEvent:
				c.Assert(t.Gone, Equals, true)
				return
			default:
				//ignore
			}
		case <-time.After(1 * time.Millisecond):
			c.Errorf("did not receive event")
			return
		}
	}

}

func (s *SessionXmppSuite) Test_WatchStanzas_presence_subscribe(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org' type='subscribe' id='adf12112'/>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		Config:         &config.Accounts{},
		CurrentAccount: &config.Account{},
		R:              roster.New(),
		ConnStatus:     DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	v, _ := sess.R.GetPendingSubscribe("some2@one.org")
	c.Assert(v, Equals, "adf12112")
}

func (s *SessionXmppSuite) Test_WatchStanzas_presence_unknown(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org' type='weird'/>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		Config:         &config.Accounts{},
		CurrentAccount: &config.Account{},
		ConnStatus:     DISCONNECTED,
	}
	sess.Conn = conn

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	select {
	case ev := <-observer:
		switch t := ev.(type) {
		case PresenceEvent:
			c.Error("Received presence event")
			return
		case PeerEvent:
			if t.Type == SubscriptionRequest {
				c.Error("Received subscription request event")
			}
			return
		default:
			// ignore
		}
	case <-time.After(1 * time.Millisecond):
		return
	}
}

func (s *SessionXmppSuite) Test_WatchStanzas_presence_regularPresenceIsAdded(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org'><client:show>dnd</client:show></client:presence>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		Config:         &config.Accounts{},
		CurrentAccount: &config.Account{},
		R:              roster.New(),
		ConnStatus:     DISCONNECTED,
	}
	sess.Conn = conn

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	st, _, _ := sess.R.StateOf("some2@one.org")
	c.Assert(st, Equals, "dnd")

	for {
		select {
		case ev := <-observer:
			switch t := ev.(type) {
			case PresenceEvent:
				c.Assert(t.Gone, Equals, false)
			default:
				//ignore
			}
			return
		case <-time.After(1 * time.Millisecond):
			c.Errorf("did not receive event")
			return
		}
	}
}

func (s *SessionXmppSuite) Test_WatchStanzas_presence_ignoresInitialAway(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org'><client:show>away</client:show></client:presence>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		Config:         &config.Accounts{},
		CurrentAccount: &config.Account{},
		R:              roster.New(),
		ConnStatus:     DISCONNECTED,
	}
	sess.Conn = conn

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	st, _, _ := sess.R.StateOf("some2@one.org")
	c.Assert(st, Equals, "")

	select {
	case ev := <-observer:
		switch ev.(type) {
		case PresenceEvent:
			c.Error("Received presence event")
			return
		default:
			// ignore
		}
	case <-time.After(1 * time.Millisecond):
		return
	}
}

func (s *SessionXmppSuite) Test_WatchStanzas_presence_ignoresSameState(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org'><client:show>dnd</client:show></client:presence>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		Config:         &config.Accounts{},
		CurrentAccount: &config.Account{},
		R:              roster.New(),
		ConnStatus:     DISCONNECTED,
	}
	sess.Conn = conn
	sess.R.AddOrReplace(roster.PeerWithState("some2@one.org", "dnd", "", ""))

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.WatchStanzas()

	st, _, _ := sess.R.StateOf("some2@one.org")
	c.Assert(st, Equals, "dnd")

	select {
	case ev := <-observer:
		switch ev.(type) {
		case PresenceEvent:
			c.Error("Received presence event")
			return
		default:
			// ignore
		}
	case <-time.After(1 * time.Millisecond):
		return
	}
}

func (s *SessionXmppSuite) Test_HandleConfirmOrDeny_failsWhenNoPendingSubscribeIsWaiting(c *C) {
	sess := &Session{
		R: roster.New(),
	}

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.HandleConfirmOrDeny("foo@bar.com", true)

	select {
	case ev := <-observer:
		t := ev.(LogEvent)
		c.Assert(t.Level, Equals, Warn)
	case <-time.After(1 * time.Millisecond):
		c.Errorf("did not receive event")
	}
}

func (s *SessionXmppSuite) Test_HandleConfirmOrDeny_succeedsOnNotAllowed(c *C) {
	mockIn := &mockConnIOReaderWriter{}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	called := 0

	sess := &Session{
		R:                   roster.New(),
		SessionEventHandler: &mockSessionEventHandler{
		//warn: func(v string) {
		//	called++
		//},
		},
	}
	sess.Conn = conn
	sess.R.SubscribeRequest("foo@bar.com", "123", "")

	sess.HandleConfirmOrDeny("foo@bar.com", false)

	c.Assert(called, Equals, 0)
	c.Assert(string(mockIn.write), Equals, "<presence id='123' to='foo@bar.com' type='unsubscribed'/>")
	_, inMap := sess.R.GetPendingSubscribe("foo@bar.com")
	c.Assert(inMap, Equals, false)
}

func (s *SessionXmppSuite) Test_HandleConfirmOrDeny_succeedsOnAllowed(c *C) {
	mockIn := &mockConnIOReaderWriter{}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	called := 0

	sess := &Session{
		R:                   roster.New(),
		SessionEventHandler: &mockSessionEventHandler{
		//warn: func(v string) {
		//	called++
		//},
		},
	}
	sess.Conn = conn
	sess.R.SubscribeRequest("foo@bar.com", "123", "")

	sess.HandleConfirmOrDeny("foo@bar.com", true)

	c.Assert(called, Equals, 0)
	c.Assert(string(mockIn.write), Equals, "<presence id='123' to='foo@bar.com' type='subscribed'/>")
	_, inMap := sess.R.GetPendingSubscribe("foo@bar.com")
	c.Assert(inMap, Equals, false)
}

func (s *SessionXmppSuite) Test_HandleConfirmOrDeny_handlesSendPresenceError(c *C) {
	mockIn := &mockConnIOReaderWriter{}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		&mockConnIOReaderWriter{err: errors.New("foo bar")},
		"some@one.org/foo",
	)

	sess := &Session{
		R: roster.New(),
	}
	sess.Conn = conn
	sess.R.SubscribeRequest("foo@bar.com", "123", "")

	observer := make(chan interface{}, 1)
	sess.Subscribe(observer)

	sess.HandleConfirmOrDeny("foo@bar.com", true)

	for {
		select {
		case ev := <-observer:
			t := ev.(LogEvent)
			if t.Level != Warn {
				continue
			}

			c.Assert(t.Message, Equals, "Error sending presence stanza: foo bar")
			return

		case <-time.After(1 * time.Millisecond):
			c.Errorf("did not receive event")
			return
		}
	}
}
