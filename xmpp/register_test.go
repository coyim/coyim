package xmpp

import (
	"bytes"
	"encoding/xml"
	"errors"

	"github.com/coyim/coyim/xmpp/data"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"
)

type RegisterSuite struct{}

var _ = Suite(&RegisterSuite{})

func (s *RegisterSuite) Test_CancelRegistration_SendCancelationRequest(c *C) {
	expectedoOut := "<iq xmlns='jabber:client' from='user@xmpp.org' type='set' id='.+'>\n" +
		"\t<query xmlns='jabber:iq:register'>\n" +
		"\t\t<remove/>\n" +
		"\t</query>\n" +
		"\t</iq>"

	mockIn := &mockConnIOReaderWriter{}
	conn := newConn()
	conn.log = testLogger()
	conn.out = mockIn
	conn.jid = "user@xmpp.org"

	_, _, err := conn.CancelRegistration()
	c.Assert(err, IsNil)
	c.Assert(string(mockIn.write), Matches, expectedoOut)
}

func (s *RegisterSuite) Test_SendChangePasswordInfo(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		log:  testLogger(),
		out:  mockOut,
		jid:  "crone1@shakespeare.lit",
		rand: bytes.NewBuffer([]byte{1, 0, 0, 0, 0, 0, 0, 0}),
	}

	conn.inflights = make(map[data.Cookie]inflight)

	reply, cookie, err := conn.sendChangePasswordInfo("crone1", "shakespeare.lit", "pass")

	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Matches, "<iq xmlns='jabber:client' to='shakespeare.lit' from='crone1@shakespeare.lit' type='set' id='1'><query xmlns='jabber:iq:register'><username>crone1</username><password>pass</password></query></iq>")
	c.Assert(reply, NotNil)
	c.Assert(cookie, NotNil)
}

func (s *RegisterSuite) Test_setupStream_registerWithoutAuthenticating(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"<register xmlns='http://jabber.org/features/iq-register'/>" +
			"</str:features>" +
			"<iq xmlns='jabber:client' type='result'>" +
			"<query xmlns='jabber:iq:register'><username/></query>" +
			"</iq>" +
			"<iq xmlns='jabber:client' type='result'></iq>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
			SkipTLS: true,
			CreateCallback: func(title, instructions string, fields []interface{}) error {
				return nil
			},
		},
		log: testLogger(),
	}
	_, err := d.setupStream(conn)

	c.Assert(err, IsNil)
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>"+
		"</stream:stream>",
	)
}

func (s *RegisterSuite) Test_conn_RegisterAccount_withNoCallback(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	cn := &conn{
		log:    testLogger(),
		out:    mockOut,
		jid:    "crone1@shakespeare.lit",
		config: data.Config{},
	}

	ok, e := cn.RegisterAccount("foo", "bar")
	c.Assert(ok, Equals, false)
	c.Assert(e, IsNil)
}

func (s *RegisterSuite) Test_conn_createAccount_getsErrorAsFinalResult_conflict(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockIn := &mockConnIOReaderWriter{read: []byte(`
<iq xmlns="jabber:client" type="result">
  <query xmlns="jabber:iq:register">
  </query>
</iq>
<iq xmlns="jabber:client" type="error">
  <error code='409' type='cancel'>
    <conflict xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>
  </error>
</iq>
`)}

	mockOut := &mockConnIOReaderWriter{}
	cn := &conn{
		log:    l,
		in:     xml.NewDecoder(mockIn),
		out:    mockOut,
		jid:    "crone1@shakespeare.lit",
		config: data.Config{},
	}

	res := cn.createAccount("hello", "goodbye")

	c.Assert(res, Equals, ErrUsernameConflict)
	c.Assert(len(hook.Entries), Equals, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to create account")
	c.Assert(hook.Entries[1].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[1].Message, Equals, "createAccount() - received the registration form")
}

func (s *RegisterSuite) Test_conn_createAccount_getsErrorAsFinalResult_notAcceptable(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockIn := &mockConnIOReaderWriter{read: []byte(`
<iq xmlns="jabber:client" type="result">
  <query xmlns="jabber:iq:register">
  </query>
</iq>
<iq xmlns="jabber:client" type="error">
  <error code='409' type='cancel'>
    <not-acceptable xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>
  </error>
</iq>
`)}

	mockOut := &mockConnIOReaderWriter{}
	cn := &conn{
		log:    l,
		in:     xml.NewDecoder(mockIn),
		out:    mockOut,
		jid:    "crone1@shakespeare.lit",
		config: data.Config{},
	}

	res := cn.createAccount("hello", "goodbye")

	c.Assert(res, Equals, ErrMissingRequiredRegistrationInfo)
	c.Assert(len(hook.Entries), Equals, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to create account")
	c.Assert(hook.Entries[1].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[1].Message, Equals, "createAccount() - received the registration form")
}

func (s *RegisterSuite) Test_conn_createAccount_getsErrorAsFinalResult_notAllowed(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockIn := &mockConnIOReaderWriter{read: []byte(`
<iq xmlns="jabber:client" type="result">
  <query xmlns="jabber:iq:register">
  </query>
</iq>
<iq xmlns="jabber:client" type="error">
  <error code='409' type='cancel'>
    <not-allowed xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>
  </error>
</iq>
`)}

	mockOut := &mockConnIOReaderWriter{}
	cn := &conn{
		log:    l,
		in:     xml.NewDecoder(mockIn),
		out:    mockOut,
		jid:    "crone1@shakespeare.lit",
		config: data.Config{},
	}

	res := cn.createAccount("hello", "goodbye")

	c.Assert(res, Equals, ErrWrongCaptcha)
	c.Assert(len(hook.Entries), Equals, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to create account")
	c.Assert(hook.Entries[1].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[1].Message, Equals, "createAccount() - received the registration form")
}

func (s *RegisterSuite) Test_conn_createAccount_getsErrorAsFinalResult_badRequest(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockIn := &mockConnIOReaderWriter{read: []byte(`
<iq xmlns="jabber:client" type="result">
  <query xmlns="jabber:iq:register">
  </query>
</iq>
<iq xmlns="jabber:client" type="error">
  <error code='409' type='cancel'>
    <bad-request xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>
  </error>
</iq>
`)}

	mockOut := &mockConnIOReaderWriter{}
	cn := &conn{
		log:    l,
		in:     xml.NewDecoder(mockIn),
		out:    mockOut,
		jid:    "crone1@shakespeare.lit",
		config: data.Config{},
	}

	res := cn.createAccount("hello", "goodbye")

	c.Assert(res, Equals, ErrRegistrationFailed)
	c.Assert(len(hook.Entries), Equals, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to create account")
	c.Assert(hook.Entries[1].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[1].Message, Equals, "createAccount() - received the registration form")
}

func (s *RegisterSuite) Test_conn_createAccount_getsErrorAsFinalResult_resourceConstraint(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockIn := &mockConnIOReaderWriter{read: []byte(`
<iq xmlns="jabber:client" type="result">
  <query xmlns="jabber:iq:register">
  </query>
</iq>
<iq xmlns="jabber:client" type="error">
  <error code='409' type='cancel'>
    <resource-constraint xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>
  </error>
</iq>
`)}

	mockOut := &mockConnIOReaderWriter{}
	cn := &conn{
		log:    l,
		in:     xml.NewDecoder(mockIn),
		out:    mockOut,
		jid:    "crone1@shakespeare.lit",
		config: data.Config{},
	}

	res := cn.createAccount("hello", "goodbye")

	c.Assert(res, Equals, ErrResourceConstraint)
	c.Assert(len(hook.Entries), Equals, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to create account")
	c.Assert(hook.Entries[1].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[1].Message, Equals, "createAccount() - received the registration form")
}

func (s *RegisterSuite) Test_conn_ChangePassword(c *C) {
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

	done := make(chan error, 1)
	go func() {
		e := cn.ChangePassword("foo", "bar.com", "baz")
		done <- e
	}()

	inf := waitForInflightTo(cn, "bar.com")
	c.Assert(inf, Not(IsNil))

	inf.replyChan <- data.Stanza{
		Value: &data.ClientIQ{
			Type: "result",
		},
	}

	res := <-done

	c.Assert(res, IsNil)
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to change account's password")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["user"], Equals, "foo")
}

func (s *RegisterSuite) Test_conn_ChangePassword_failsWithError_badRequest(c *C) {
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

	done := make(chan error, 1)
	go func() {
		e := cn.ChangePassword("foo", "bar.com", "baz")
		done <- e
	}()

	inf := waitForInflightTo(cn, "bar.com")
	c.Assert(inf, Not(IsNil))

	ciq := &data.ClientIQ{}
	_ = xml.NewDecoder(bytes.NewBuffer([]byte(`
	<iq xmlns="jabber:client" type="error">
  <error code='409' type='cancel'>
    <bad-request xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>
  </error>
</iq>
`))).DecodeElement(ciq, nil)

	inf.replyChan <- data.Stanza{
		Value: ciq,
	}

	res := <-done

	c.Assert(res, Equals, ErrBadRequest)
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to change account's password")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["user"], Equals, "foo")
}

func (s *RegisterSuite) Test_conn_ChangePassword_failsWithError_notAuthorized(c *C) {
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

	done := make(chan error, 1)
	go func() {
		e := cn.ChangePassword("foo", "bar.com", "baz")
		done <- e
	}()

	inf := waitForInflightTo(cn, "bar.com")
	c.Assert(inf, Not(IsNil))

	ciq := &data.ClientIQ{}
	_ = xml.NewDecoder(bytes.NewBuffer([]byte(`
	<iq xmlns="jabber:client" type="error">
  <error code='409' type='cancel'>
    <not-authorized xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>
  </error>
</iq>
`))).DecodeElement(ciq, nil)

	inf.replyChan <- data.Stanza{
		Value: ciq,
	}

	res := <-done

	c.Assert(res, Equals, ErrNotAuthorized)
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to change account's password")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["user"], Equals, "foo")
}

func (s *RegisterSuite) Test_conn_ChangePassword_failsWithError_notAllowed(c *C) {
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

	done := make(chan error, 1)
	go func() {
		e := cn.ChangePassword("foo", "bar.com", "baz")
		done <- e
	}()

	inf := waitForInflightTo(cn, "bar.com")
	c.Assert(inf, Not(IsNil))

	ciq := &data.ClientIQ{}
	_ = xml.NewDecoder(bytes.NewBuffer([]byte(`
	<iq xmlns="jabber:client" type="error">
  <error code='409' type='cancel'>
    <not-allowed xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>
  </error>
</iq>
`))).DecodeElement(ciq, nil)

	inf.replyChan <- data.Stanza{
		Value: ciq,
	}

	res := <-done

	c.Assert(res, Equals, ErrNotAllowed)
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to change account's password")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["user"], Equals, "foo")
}

func (s *RegisterSuite) Test_conn_ChangePassword_failsWithError_unexpectedRequest(c *C) {
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

	done := make(chan error, 1)
	go func() {
		e := cn.ChangePassword("foo", "bar.com", "baz")
		done <- e
	}()

	inf := waitForInflightTo(cn, "bar.com")
	c.Assert(inf, Not(IsNil))

	ciq := &data.ClientIQ{}
	_ = xml.NewDecoder(bytes.NewBuffer([]byte(`
	<iq xmlns="jabber:client" type="error">
  <error code='409' type='cancel'>
    <unexpected-request xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>
  </error>
</iq>
`))).DecodeElement(ciq, nil)

	inf.replyChan <- data.Stanza{
		Value: ciq,
	}

	res := <-done

	c.Assert(res, Equals, ErrUnexpectedRequest)
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to change account's password")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["user"], Equals, "foo")
}

func (s *RegisterSuite) Test_conn_ChangePassword_failsWithError_other(c *C) {
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

	done := make(chan error, 1)
	go func() {
		e := cn.ChangePassword("foo", "bar.com", "baz")
		done <- e
	}()

	inf := waitForInflightTo(cn, "bar.com")
	c.Assert(inf, Not(IsNil))

	ciq := &data.ClientIQ{}
	_ = xml.NewDecoder(bytes.NewBuffer([]byte(`
	<iq xmlns="jabber:client" type="error">
  <error code='409' type='cancel'>
    <foobar xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>
  </error>
</iq>
`))).DecodeElement(ciq, nil)

	inf.replyChan <- data.Stanza{
		Value: ciq,
	}

	res := <-done

	c.Assert(res, Equals, ErrChangePasswordFailed)
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to change account's password")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["user"], Equals, "foo")
}

func (s *RegisterSuite) Test_conn_ChangePassword_badCiqType(c *C) {
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

	done := make(chan error, 1)
	go func() {
		e := cn.ChangePassword("foo", "bar.com", "baz")
		done <- e
	}()

	inf := waitForInflightTo(cn, "bar.com")
	c.Assert(inf, Not(IsNil))

	inf.replyChan <- data.Stanza{
		Value: &data.ClientIQ{
			Type: "set",
		},
	}

	res := <-done

	c.Assert(res, Equals, ErrChangePasswordFailed)
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to change account's password")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["user"], Equals, "foo")
}

func (s *RegisterSuite) Test_conn_ChangePassword_noClientIQ(c *C) {
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

	done := make(chan error, 1)
	go func() {
		e := cn.ChangePassword("foo", "bar.com", "baz")
		done <- e
	}()

	inf := waitForInflightTo(cn, "bar.com")
	c.Assert(inf, Not(IsNil))

	inf.replyChan <- data.Stanza{
		Value: "hmm",
	}

	res := <-done

	c.Assert(res, ErrorMatches, "xmpp: failed to parse response")
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to change account's password")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["user"], Equals, "foo")
}

func (s *RegisterSuite) Test_conn_ChangePassword_noResultChannel(c *C) {
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

	done := make(chan error, 1)
	go func() {
		e := cn.ChangePassword("foo", "bar.com", "baz")
		done <- e
	}()

	inf := waitForInflightTo(cn, "bar.com")
	c.Assert(inf, Not(IsNil))

	close(inf.replyChan)

	res := <-done

	c.Assert(res, ErrorMatches, "xmpp: failed to receive response")
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to change account's password")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["user"], Equals, "foo")
}

func (s *RegisterSuite) Test_conn_ChangePassword_failedToSend(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockOut := &mockConnIOReaderWriter{
		err: errors.New("an IO error"),
	}

	cn := &conn{
		closed:    false,
		log:       l,
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	e := cn.ChangePassword("foo", "bar.com", "baz")

	c.Assert(e, ErrorMatches, "xmpp: failed to send request")
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to change account's password")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["user"], Equals, "foo")
}

func (s *RegisterSuite) Test_conn_createAccount_failsProcessingForm(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockIn := &mockConnIOReaderWriter{read: []byte(`
<iq xmlns="jabber:client" type="result">
  <query xmlns="jabber:iq:register">
    <x xmlns="jabber:x:data" type="something">
    </x>
  </query>
</iq>
`)}

	mockOut := &mockConnIOReaderWriter{}
	cn := &conn{
		log: l,
		in:  xml.NewDecoder(mockIn),
		out: mockOut,
		jid: "crone1@shakespeare.lit",
		config: data.Config{
			CreateCallback: func(string, string, []interface{}) error {
				return errors.New("couldn't create form")
			},
		},
	}

	res := cn.createAccount("hello", "goodbye")

	c.Assert(res, ErrorMatches, "couldn't create form")
	c.Assert(len(hook.Entries), Equals, 4)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to create account")
	c.Assert(hook.Entries[1].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[1].Message, Equals, "createAccount() - received the registration form")
}

func (s *RegisterSuite) Test_conn_createAccount_failsOnWritingToRawOut(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockIn := &mockConnIOReaderWriter{read: []byte(`
<iq xmlns="jabber:client" type="result">
  <query xmlns="jabber:iq:register">
    <x xmlns="jabber:x:data" type="something">
    </x>
  </query>
</iq>
`)}

	mockOut := &mockConnIOReaderWriter{
		err:      errors.New("Oh noes"),
		errCount: 2,
	}
	cn := &conn{
		log:    l,
		in:     xml.NewDecoder(mockIn),
		out:    mockOut,
		rawOut: mockOut,
		jid:    "crone1@shakespeare.lit",
		config: data.Config{
			CreateCallback: func(string, string, []interface{}) error {
				return nil
			},
		},
	}

	res := cn.createAccount("hello", "goodbye")

	c.Assert(res, ErrorMatches, "Oh noes")
	c.Assert(len(hook.Entries), Equals, 3)
	c.Assert(hook.Entries[0].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Attempting to create account")
	c.Assert(hook.Entries[1].Level, Equals, log.DebugLevel)
	c.Assert(hook.Entries[1].Message, Equals, "createAccount() - received the registration form")
}
