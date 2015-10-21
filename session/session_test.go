package session

import (
	"encoding/xml"
	"errors"
	"testing"
	"time"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/event"
	"github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/otr3"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type SessionXmppSuite struct{}

var _ = Suite(&SessionXmppSuite{})

func (s *SessionXmppSuite) Test_NewSession_returnsANewSession(c *C) {
	sess := NewSession(&config.Config{})
	c.Assert(sess, Not(IsNil))
}

func (s *SessionXmppSuite) Test_info_callsSessionHandlerInfo(c *C) {
	called := 0
	calledWith := ""

	sess := &Session{
		SessionEventHandler: &mockSessionEventHandler{
			info: func(v string) {
				called++
				calledWith = v
			},
		},
	}

	sess.info("hello world")

	c.Assert(called, Equals, 1)
	c.Assert(calledWith, Equals, "hello world")
}

func (s *SessionXmppSuite) Test_warn_callsSessionHandlerWarn(c *C) {
	called := 0
	calledWith := ""

	sess := &Session{
		SessionEventHandler: &mockSessionEventHandler{
			warn: func(v string) {
				called++
				calledWith = v
			},
		},
	}

	sess.warn("hello world2")

	c.Assert(called, Equals, 1)
	c.Assert(calledWith, Equals, "hello world2")
}

func (s *SessionXmppSuite) Test_alert_callsSessionHandlerAlert(c *C) {
	called := 0
	calledWith := ""

	sess := &Session{
		SessionEventHandler: &mockSessionEventHandler{
			alert: func(v string) {
				called++
				calledWith = v
			},
		},
	}

	sess.alert("hello world3")

	c.Assert(called, Equals, 1)
	c.Assert(calledWith, Equals, "hello world3")
}

func (s *SessionXmppSuite) Test_iqReceived_passesToSessionHandler(c *C) {
	called := 0
	calledWith := ""

	sess := &Session{
		SessionEventHandler: &mockSessionEventHandler{
			iqReceived: func(v string) {
				called++
				calledWith = v
			},
		},
	}

	sess.iqReceived("someone@somewhere")

	c.Assert(called, Equals, 1)
	c.Assert(calledWith, Equals, "someone@somewhere")
}

func (s *SessionXmppSuite) Test_readMessages_passesStanzaToChannel(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:message xmlns:client='jabber:client' to='fo@bar.com' from='bar@foo.com' type='chat'><client:body>something</client:body></client:message>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		nil,
		"some@one.org/foo",
	)

	sess := &Session{
		SessionEventHandler: &mockSessionEventHandler{},
	}
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
	called := 0
	calledWith := ""

	sess := &Session{
		SessionEventHandler: &mockSessionEventHandler{
			alert: func(v string) {
				called++
				calledWith = v
			},
		},
	}
	sess.Conn = conn

	stanzaChan := make(chan xmpp.Stanza)
	go sess.readMessages(stanzaChan)

	select {
	case _, ok := <-stanzaChan:
		c.Assert(ok, Equals, false)
	}

	c.Assert(called, Equals, 1)
	c.Assert(calledWith, Equals, "unexpected XMPP message clientx <message/>")
}

func (s *SessionXmppSuite) Test_WatchStanzas_warnsAndExitsOnBadStanza(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<clientx:message xmlns:client='jabber:client' to='fo@bar.com' from='bar@foo.com' type='chat'><client:body>something</client:body></client:message>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		nil,
		"some@one.org/foo",
	)

	warnCalled := 0
	warnCalledWith := ""

	alertCalled := 0
	alertCalledWith := ""

	sess := &Session{
		SessionEventHandler: &mockSessionEventHandler{
			warn: func(v string) {
				warnCalled++
				warnCalledWith = v
			},
			alert: func(v string) {
				alertCalled++
				alertCalledWith = v
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(warnCalled, Equals, 1)
	c.Assert(warnCalledWith, Equals, "Exiting because channel to server closed")

	c.Assert(alertCalled, Equals, 1)
	c.Assert(alertCalledWith, Equals, "unexpected XMPP message clientx <message/>")
}

func (s *SessionXmppSuite) Test_WatchStanzas_handlesUnknownMessage(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<bind:bind xmlns:bind='urn:ietf:params:xml:ns:xmpp-bind'></bind:bind>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		nil,
		"some@one.org/foo",
	)

	infoCalled := 0
	infoCalledWith := ""

	sess := &Session{
		SessionEventHandler: &mockSessionEventHandler{
			info: func(v string) {
				infoCalled++
				infoCalledWith = v
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(infoCalled, Equals, 1)
	c.Assert(infoCalledWith, Equals, "{urn:ietf:params:xml:ns:xmpp-bind bind} &{{urn:ietf:params:xml:ns:xmpp-bind bind}  }")
}

func (s *SessionXmppSuite) Test_WatchStanzas_handlesStreamError_withText(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:error xmlns:stream='http://etherx.jabber.org/streams'><stream:text>bad horse showed up</stream:text></stream:error>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		nil,
		"some@one.org/foo",
	)

	alertCalled := 0
	alertCalledWith := ""

	sess := &Session{
		SessionEventHandler: &mockSessionEventHandler{
			alert: func(v string) {
				alertCalled++
				alertCalledWith = v
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(alertCalled, Equals, 2)
	c.Assert(alertCalledWith, Equals, "Exiting in response to fatal error from server: bad horse showed up")
}

func (s *SessionXmppSuite) Test_WatchStanzas_handlesStreamError_withEmbeddedTag(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:error xmlns:stream='http://etherx.jabber.org/streams'><not-well-formed xmlns='urn:ietf:params:xml:ns:xmpp-streams'/></stream:error>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		nil,
		"some@one.org/foo",
	)

	alertCalled := 0
	alertCalledWith := ""

	sess := &Session{
		SessionEventHandler: &mockSessionEventHandler{
			alert: func(v string) {
				alertCalled++
				alertCalledWith = v
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(alertCalled, Equals, 2)
	c.Assert(alertCalledWith, Equals, "Exiting in response to fatal error from server: {urn:ietf:params:xml:ns:xmpp-streams not-well-formed}")
}

func (s *SessionXmppSuite) Test_WatchStanzas_receivesAMessage(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:message xmlns:client='jabber:client' type='chat' to='some@one.org/foo' from='bla@hmm.org/somewhere'><client:body>well, hello there</client:body></client:message>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	called := 0

	var sess *Session
	sess = &Session{
		Config:          &config.Config{},
		Conversations:   make(map[string]*otr3.Conversation),
		OtrEventHandler: make(map[string]*event.OtrEventHandler),
		SessionEventHandler: &mockSessionEventHandler{
			messageReceived: func(s *Session, from string, timestamp time.Time, encrypted bool, message []byte) {
				called++
				c.Assert(s, Equals, sess)
				c.Assert(encrypted, Equals, false)
				c.Assert(from, Equals, "bla@hmm.org")
				c.Assert(string(message), Equals, "well, hello there")
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(called, Equals, 1)
}

func (s *SessionXmppSuite) Test_WatchStanzas_failsOnUnrecognizedIQ(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='something'></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	called := 0

	var sess *Session
	sess = &Session{
		SessionEventHandler: &mockSessionEventHandler{
			info: func(v string) {
				called++
				c.Assert(v, Equals, "unrecognized iq: &xmpp.ClientIQ{XMLName:xml.Name{Space:\"jabber:client\", Local:\"iq\"}, From:\"\", ID:\"\", To:\"\", Type:\"something\", Error:xmpp.ClientError{XMLName:xml.Name{Space:\"\", Local:\"\"}, Code:\"\", Type:\"\", Any:xml.Name{Space:\"\", Local:\"\"}, Text:\"\"}, Bind:xmpp.bindBind{XMLName:xml.Name{Space:\"\", Local:\"\"}, Resource:\"\", Jid:\"\"}, Query:[]uint8{}}")
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(called, Equals, 1)
}

func (s *SessionXmppSuite) Test_WatchStanzas_getsDiscoInfoIQ(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='get' from='abc' to='cde'><query xmlns='http://jabber.org/protocol/disco#info'/></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	sess := &Session{
		Config: &config.Config{
			Account: "foo.bar@somewhere.org",
		},
		SessionEventHandler: &mockSessionEventHandler{},
		ConnStatus:          DISCONNECTED,
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
		Config: &config.Config{
			Account: "foo.bar@somewhere.org",
		},
		SessionEventHandler: &mockSessionEventHandler{},
		ConnStatus:          DISCONNECTED,
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

	called := 0

	sess := &Session{
		Config: &config.Config{
			Account: "foo.bar@somewhere.org",
		},
		SessionEventHandler: &mockSessionEventHandler{
			info: func(v string) {
				called++
				c.Assert(v, Equals, "Unknown IQ: jabber:iq:somethingStrange query")
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(called, Equals, 1)
}

func (s *SessionXmppSuite) Test_WatchStanzas_iq_set_roster_withBadFrom(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='get' from='some2@one.org' to='cde'><query xmlns='jabber:iq:roster'/></client:iq>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	called := 0

	sess := &Session{
		Config: &config.Config{Account: "some@one.org"},
		SessionEventHandler: &mockSessionEventHandler{
			warn: func(v string) {
				called++
				if called == 1 {
					c.Assert(v, Equals, "Ignoring roster IQ from bad address: some2@one.org")
				}
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(called, Equals, 2)
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

	called := 0

	sess := &Session{
		Config: &config.Config{Account: "some@one.org"},
		SessionEventHandler: &mockSessionEventHandler{
			warn: func(v string) {
				called++
				if called == 1 {
					c.Assert(v, Equals, "Failed to parse roster push IQ")
				}
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(called, Equals, 2)
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

	called := 0

	sess := &Session{
		Config: &config.Config{Account: "some@one.org"},
		R:      roster.New(),
		SessionEventHandler: &mockSessionEventHandler{
			iqReceived: func(v string) {
				called++
				c.Assert(v, Equals, "romeo@example.net")
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(called, Equals, 1)
	c.Assert(sess.R.ToSlice(), DeepEquals, []*roster.Peer{
		roster.PeerFrom(xmpp.RosterEntry{Jid: "romeo@example.net", Subscription: "both", Name: "Romeo", Group: []string{"Friends"}})})
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
		Config: &config.Config{Account: "some@one.org"},
		R:      roster.New(),
		SessionEventHandler: &mockSessionEventHandler{
			iqReceived: func(v string) {
				called++
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.R.AddOrReplace(roster.PeerFrom(xmpp.RosterEntry{Jid: "jill@example.net", Subscription: "both", Name: "Jill", Group: []string{"Foes"}}))
	sess.R.AddOrReplace(roster.PeerFrom(xmpp.RosterEntry{Jid: "romeo@example.net", Subscription: "both", Name: "Mo", Group: []string{"Foes"}}))

	sess.WatchStanzas()

	c.Assert(called, Equals, 0)
	c.Assert(sess.R.ToSlice(), DeepEquals, []*roster.Peer{
		roster.PeerFrom(xmpp.RosterEntry{Jid: "jill@example.net", Subscription: "both", Name: "Jill", Group: []string{"Foes"}}),
		roster.PeerFrom(xmpp.RosterEntry{Jid: "romeo@example.net", Subscription: "both", Name: "Romeo", Group: []string{"Friends"}}),
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

	called := 0

	sess := &Session{
		Config: &config.Config{Account: "some@one.org"},
		R:      roster.New(),
		SessionEventHandler: &mockSessionEventHandler{
			iqReceived: func(v string) {
				called++
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.R.AddOrReplace(roster.PeerFrom(xmpp.RosterEntry{Jid: "romeo@example.net", Subscription: "both", Name: "Mo", Group: []string{"Foes"}}))
	sess.R.AddOrReplace(roster.PeerFrom(xmpp.RosterEntry{Jid: "jill@example.net", Subscription: "both", Name: "Jill", Group: []string{"Foes"}}))
	sess.R.AddOrReplace(roster.PeerFrom(xmpp.RosterEntry{Jid: "romeo@example.net", Subscription: "both", Name: "Mo", Group: []string{"Foes"}}))

	sess.WatchStanzas()

	c.Assert(called, Equals, 0)
	c.Assert(sess.R.ToSlice(), DeepEquals, []*roster.Peer{
		roster.PeerFrom(xmpp.RosterEntry{Jid: "jill@example.net", Subscription: "both", Name: "Jill", Group: []string{"Foes"}}),
	})
}

func (s *SessionXmppSuite) Test_WatchStanzas_presence_unavailable_forNoneKnownUser(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org' type='unavailable'><client:status>going on vacation</client:status></client:presence>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	called := 0

	sess := &Session{
		R: roster.New(),
		SessionEventHandler: &mockSessionEventHandler{
			processPresence: func(from, to, show, status string, gone bool) {
				called++
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(called, Equals, 0)
}

func (s *SessionXmppSuite) Test_WatchStanzas_presence_unavailable_forKnownUser(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org' type='unavailable'><client:status>going on vacation</client:status></client:presence>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	called := 0

	sess := &Session{
		Config: &config.Config{},
		R:      roster.New(),
		SessionEventHandler: &mockSessionEventHandler{
			processPresence: func(from, to, show, status string, gone bool) {
				called++
				c.Assert(gone, Equals, true)
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn
	sess.R.AddOrReplace(roster.PeerWithState("some2@one.org", "somewhere", ""))

	sess.WatchStanzas()

	c.Assert(called, Equals, 1)
	p, _ := sess.R.Get("some2@one.org")
	c.Assert(p.Online, Equals, false)
}

func (s *SessionXmppSuite) Test_WatchStanzas_presence_subscribe(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org' type='subscribe' id='adf12112'/>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	called := 0

	sess := &Session{
		Config: &config.Config{},
		R:      roster.New(),
		SessionEventHandler: &mockSessionEventHandler{
			subscriptionRequest: func(s *Session, uid string) {
				called++
				c.Assert(uid, Equals, "some2@one.org")
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(called, Equals, 1)
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

	called := 0

	sess := &Session{
		Config: &config.Config{},
		SessionEventHandler: &mockSessionEventHandler{
			subscriptionRequest: func(s *Session, uid string) {
				called++
			},
			processPresence: func(from, to, show, status string, gone bool) {
				called++
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(called, Equals, 0)
}

func (s *SessionXmppSuite) Test_WatchStanzas_presence_regularPresenceIsAdded(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org'><client:show>dnd</client:show></client:presence>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	called := 0

	sess := &Session{
		Config: &config.Config{},
		R:      roster.New(),
		SessionEventHandler: &mockSessionEventHandler{
			processPresence: func(from, to, show, status string, gone bool) {
				called++
				c.Assert(gone, Equals, false)
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(called, Equals, 1)
	st, _, _ := sess.R.StateOf("some2@one.org")
	c.Assert(st, Equals, "dnd")
}

func (s *SessionXmppSuite) Test_WatchStanzas_presence_ignoresInitialAway(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org'><client:show>away</client:show></client:presence>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	called := 0

	sess := &Session{
		Config: &config.Config{},
		R:      roster.New(),
		SessionEventHandler: &mockSessionEventHandler{
			processPresence: func(from, to, show, status string, gone bool) {
				called++
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn

	sess.WatchStanzas()

	c.Assert(called, Equals, 0)
	st, _, _ := sess.R.StateOf("some2@one.org")
	c.Assert(st, Equals, "")
}

func (s *SessionXmppSuite) Test_WatchStanzas_presence_ignoresSameState(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:presence xmlns:client='jabber:client' from='some2@one.org/balcony' to='some@one.org'><client:show>dnd</client:show></client:presence>")}
	conn := xmpp.NewConn(
		xml.NewDecoder(mockIn),
		mockIn,
		"some@one.org/foo",
	)

	called := 0

	sess := &Session{
		Config: &config.Config{},
		R:      roster.New(),
		SessionEventHandler: &mockSessionEventHandler{
			processPresence: func(from, to, show, status string, gone bool) {
				called++
			},
		},
		ConnStatus: DISCONNECTED,
	}
	sess.Conn = conn
	sess.R.AddOrReplace(roster.PeerWithState("some2@one.org", "dnd", ""))

	sess.WatchStanzas()

	c.Assert(called, Equals, 0)
	st, _, _ := sess.R.StateOf("some2@one.org")
	c.Assert(st, Equals, "dnd")
}

func (s *SessionXmppSuite) Test_HandleConfirmOrDeny_failsWhenNoPendingSubscribeIsWaiting(c *C) {
	called := 0

	sess := &Session{
		R: roster.New(),
		SessionEventHandler: &mockSessionEventHandler{
			warn: func(v string) {
				called++
				c.Assert(v, Equals, "No pending subscription from foo@bar.com")
			},
		},
	}

	sess.HandleConfirmOrDeny("foo@bar.com", true)
	c.Assert(called, Equals, 1)
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
		R: roster.New(),
		SessionEventHandler: &mockSessionEventHandler{
			warn: func(v string) {
				called++
			},
		},
	}
	sess.Conn = conn
	sess.R.SubscribeRequest("foo@bar.com", "123")

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
		R: roster.New(),
		SessionEventHandler: &mockSessionEventHandler{
			warn: func(v string) {
				called++
			},
		},
	}
	sess.Conn = conn
	sess.R.SubscribeRequest("foo@bar.com", "123")

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

	called := 0

	sess := &Session{
		R: roster.New(),
		SessionEventHandler: &mockSessionEventHandler{
			warn: func(v string) {
				called++
				c.Assert(v, Equals, "Error sending presence stanza: foo bar")
			},
		},
	}
	sess.Conn = conn
	sess.R.SubscribeRequest("foo@bar.com", "123")

	sess.HandleConfirmOrDeny("foo@bar.com", true)

	c.Assert(called, Equals, 1)
}
