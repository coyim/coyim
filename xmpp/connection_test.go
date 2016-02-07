package xmpp

import (
	"crypto/tls"
	"encoding/xml"
	"io"
	"runtime"
	"strings"

	"github.com/twstrike/coyim/xmpp/data"

	. "gopkg.in/check.v1"
)

type ConnectionXmppSuite struct{}

var _ = Suite(&ConnectionXmppSuite{})

func (s *ConnectionXmppSuite) Test_Next_returnsErrorIfOneIsEncountered(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:foo xmlns:stream='http://etherx.jabber.org/streams' to='hello'></stream:foo>")}
	conn := Conn{
		in: xml.NewDecoder(mockIn),
	}

	_, err := conn.Next()
	c.Assert(err.Error(), Equals, "unexpected XMPP message http://etherx.jabber.org/streams <foo/>")
}

func (s *ConnectionXmppSuite) Test_Next_returnsErrorIfFailingToParseIQID(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='abczzzz'></client:iq>")}
	conn := Conn{
		in: xml.NewDecoder(mockIn),
	}

	_, err := conn.Next()
	c.Assert(err.Error(), Equals, "xmpp: failed to parse id from iq: strconv.ParseUint: parsing \"abczzzz\": invalid syntax")
}

func (s *ConnectionXmppSuite) Test_Next_returnsNothingIfThereIsNoInflightMatching(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000'></client:iq>")}
	conn := Conn{
		in: xml.NewDecoder(mockIn),
	}

	_, err := conn.Next()
	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXmppSuite) Test_Next_returnsNothingIfTheInflightIsToAnotherReceiver(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000' from='bar@somewhere.com'></client:iq>")}
	conn := Conn{
		in:        xml.NewDecoder(mockIn),
		inflights: make(map[Cookie]inflight),
	}
	cookie := Cookie(1048576)
	conn.inflights[cookie] = inflight{to: "foo@somewhere.com"}
	_, err := conn.Next()
	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXmppSuite) Test_Next_removesInflightIfItMatches(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000' from='foo@somewhere.com'></client:iq>")}
	inflights := make(map[Cookie]inflight)
	conn := Conn{
		in:        xml.NewDecoder(mockIn),
		inflights: inflights,
	}
	cookie := Cookie(1048576)
	reply := make(chan data.Stanza, 1)
	conn.inflights[cookie] =
		inflight{
			to:        "foo@somewhere.com",
			replyChan: reply,
		}

	go func() {
		<-reply
	}()

	_, err := conn.Next()
	c.Assert(err, Equals, io.EOF)
	_, ok := conn.inflights[cookie]
	c.Assert(ok, Equals, false)
}

func (s *ConnectionXmppSuite) Test_Next_continuesIfIqFromIsNotSimilarToJid(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000' from='foo@somewhere.com'></client:iq>")}
	inflights := make(map[Cookie]inflight)
	conn := Conn{
		in:        xml.NewDecoder(mockIn),
		inflights: inflights,
		jid:       "foo@myjid.com/blah",
	}
	cookie := Cookie(1048576)
	conn.inflights[cookie] = inflight{}
	_, err := conn.Next()
	c.Assert(err, Equals, io.EOF)
	_, ok := conn.inflights[cookie]
	c.Assert(ok, Equals, true)
}

func (s *ConnectionXmppSuite) Test_Next_removesIfThereIsNoFrom(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000'></client:iq>")}
	inflights := make(map[Cookie]inflight)
	conn := Conn{
		in:        xml.NewDecoder(mockIn),
		inflights: inflights,
	}
	cookie := Cookie(1048576)
	reply := make(chan data.Stanza, 1)
	conn.inflights[cookie] =
		inflight{
			replyChan: reply,
		}

	go func() {
		<-reply
	}()

	_, err := conn.Next()
	c.Assert(err, Equals, io.EOF)
	_, ok := conn.inflights[cookie]
	c.Assert(ok, Equals, false)
}

func (s *ConnectionXmppSuite) Test_Next_removesIfThereIsTheFromIsSameAsJid(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000' from='some@one.org/foo'></client:iq>")}
	inflights := make(map[Cookie]inflight)
	conn := Conn{
		in:        xml.NewDecoder(mockIn),
		inflights: inflights,
		jid:       "some@one.org/foo",
	}
	cookie := Cookie(1048576)
	reply := make(chan data.Stanza, 1)
	conn.inflights[cookie] =
		inflight{
			replyChan: reply,
		}

	go func() {
		<-reply
	}()

	_, err := conn.Next()
	c.Assert(err, Equals, io.EOF)
	_, ok := conn.inflights[cookie]
	c.Assert(ok, Equals, false)
}

func (s *ConnectionXmppSuite) Test_Next_removesIfThereIsTheFromIsSameAsJidWithoutResource(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000' from='some@one.org'></client:iq>")}
	inflights := make(map[Cookie]inflight)
	conn := Conn{
		in:        xml.NewDecoder(mockIn),
		inflights: inflights,
		jid:       "some@one.org/foo",
	}
	cookie := Cookie(1048576)
	reply := make(chan data.Stanza, 1)
	conn.inflights[cookie] =
		inflight{
			replyChan: reply,
		}

	go func() {
		<-reply
	}()

	_, err := conn.Next()
	c.Assert(err, Equals, io.EOF)
	_, ok := conn.inflights[cookie]
	c.Assert(ok, Equals, false)
}

func (s *ConnectionXmppSuite) Test_Next_removesIfThereIsTheFromIsSameAsJidDomain(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000' from='one.org'></client:iq>")}
	inflights := make(map[Cookie]inflight)
	conn := Conn{
		in:        xml.NewDecoder(mockIn),
		inflights: inflights,
		jid:       "some@one.org/foo",
	}
	cookie := Cookie(1048576)
	reply := make(chan data.Stanza, 1)
	conn.inflights[cookie] =
		inflight{
			replyChan: reply,
		}

	go func() {
		<-reply
	}()

	_, err := conn.Next()
	c.Assert(err, Equals, io.EOF)
	_, ok := conn.inflights[cookie]
	c.Assert(ok, Equals, false)
}

func (s *ConnectionXmppSuite) Test_Next_returnsNonIQMessage(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:message xmlns:client='jabber:client' to='fo@bar.com' from='bar@foo.com' type='chat'><client:body>something</client:body></client:message>")}
	conn := Conn{
		in:  xml.NewDecoder(mockIn),
		jid: "some@one.org/foo",
	}
	v, err := conn.Next()
	c.Assert(err, IsNil)
	c.Assert(v.Value.(*data.ClientMessage).From, Equals, "bar@foo.com")
	c.Assert(v.Value.(*data.ClientMessage).To, Equals, "fo@bar.com")
	c.Assert(v.Value.(*data.ClientMessage).Type, Equals, "chat")
	c.Assert(v.Value.(*data.ClientMessage).Body, Equals, "something")
}

func (s *ConnectionXmppSuite) Test_makeInOut_returnsANewDecoderAndOriginalWriterWhenNoConfigIsGiven(c *C) {
	mockBoth := &mockConnIOReaderWriter{}
	_, rout := makeInOut(mockBoth, Config{})
	c.Assert(rout, Equals, mockBoth)
}

func (s *ConnectionXmppSuite) Test_makeInOut_returnsANewDecoderAndWrappedWriterWhenConfigIsGiven(c *C) {
	mockBoth := &mockConnIOReaderWriter{}
	mockInLog := &mockConnIOReaderWriter{}
	config := Config{InLog: mockInLog, OutLog: mockInLog}
	_, rout := makeInOut(mockBoth, config)
	c.Assert(rout, Not(Equals), mockBoth)
}

func (s *ConnectionXmppSuite) Test_Dial_returnsErrorFromGetFeatures(c *C) {
	rw := &mockConnIOReaderWriter{}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXmppSuite) Test_Dial_returnsErrorFromAuthenticateIfSkipTLS(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'><str:features></str:features>")}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config:   Config{SkipTLS: true},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, ErrAuthenticationFailed)
}

func (s *ConnectionXmppSuite) Test_Dial_returnsErrorFromSecondFeatureCheck(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>")}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config:   Config{SkipTLS: true},
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Matches, "(XML syntax error on line 1: unexpected )?EOF")

	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n")
}

func (s *ConnectionXmppSuite) Test_Dial_returnsErrorFromIQReturn(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>" +
			"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"</str:features>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config:   Config{SkipTLS: true},
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Matches, "unmarshal <iq>:( XML syntax error on line 1: unexpected)? EOF")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_returnsWorkingConnIfEverythingPasses(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>" +
			"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"</str:features>" +
			"<client:iq xmlns:client='jabber:client'></client:iq>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config:   Config{SkipTLS: true},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, IsNil)
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_failsIfTheServerDoesntSupportTLS(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "xmpp: server doesn't support TLS")
}

func (s *ConnectionXmppSuite) Test_Dial_failsIfReceivingEOFAfterStartingTLS(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'>" +
			"</starttls>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Matches, "(XML syntax error on line 1: unexpected )?EOF")
}

func (s *ConnectionXmppSuite) Test_Dial_failsIfReceivingTheWrongNamespaceAfterStarttls(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'>" +
			"</starttls>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<str:proceed>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "xmpp: expected <proceed> after <starttls> but got <proceed> in http://etherx.jabber.org/streams")
}

func (s *ConnectionXmppSuite) Test_Dial_failsIfReceivingTheWrongTagName(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'>" +
			"</starttls>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<things xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "xmpp: expected <proceed> after <starttls> but got <things> in urn:ietf:params:xml:ns:xmpp-tls")
}

func (s *ConnectionXmppSuite) Test_Dial_failsWhenStartingAHandshake(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'>" +
			"</starttls>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<proceed xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>",
	)}
	conn := &fullMockedConn{rw: rw}
	var tlsC tls.Config
	tlsC.Rand = fixedRand([]string{"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F"})

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config: Config{
			TLSConfig: &tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
	if isVersionOldish() {
		c.Assert(string(rw.write), Equals, ""+
			"<?xml version='1.0'?><stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>\x16\x03\x01\x00\x84\x01\x00\x00\x80\x03\x03\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
			"\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x00\x00\x1a\xc0/\xc0+\xc0\x11\xc0\a\xc0\x13\xc0\t\xc0\x14\xc0\n"+
			"\x00\x05\x00/\x005\xc0\x12\x00\n"+
			"\x01\x00\x00=\x00\x00\x00\v\x00\t\x00\x00\x06domain\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\n"+
			"\x00\b\x00\x06\x00\x17\x00\x18\x00\x19\x00\v\x00\x02\x01\x00\x00\r\x00\n"+
			"\x00\b\x04\x01\x04\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00",
		)
	} else {
		c.Assert(string(rw.write), Equals, ""+
			"<?xml version='1.0'?><stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>\x16\x03\x01\x00\x8a\x01\x00\x00\x86\x03\x03\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
			"\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x00\x00\x18\xc0/\xc0+\xc00\xc0,\xc0\x13\xc0\t\xc0\x14\xc0\n"+
			"\x00/\x005\xc0\x12\x00\n"+
			"\x01\x00\x00E\x00\x00\x00\v\x00\t\x00\x00\x06domain\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\n"+
			"\x00\b\x00\x06\x00\x17\x00\x18\x00\x19\x00\v\x00\x02\x01\x00\x00\r\x00\x0e\x00\f\x04\x01\x04\x03\x05\x01\x05\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00\x00\x12\x00\x00",
		)
	}
}

func (s *ConnectionXmppSuite) Test_Dial_setsServerNameOnTLSContext(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'>" +
			"</starttls>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<proceed xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>",
	)}
	var tlsC tls.Config
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config: Config{
			TLSConfig: &tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXmppSuite) Test_Dial_failsIfDecodingFallbackFails(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"<register xmlns='http://jabber.org/features/iq-register'/>" +
			"</str:features>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config: Config{
			SkipTLS: true,
			CreateCallback: func(title, instructions string, fields []interface{}) error {
				return nil
			},
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Matches, "unmarshal <iq>:( XML syntax error on line 1: unexpected)? EOF")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_failsIfAccountCreationFails(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"<register xmlns='http://jabber.org/features/iq-register'/>" +
			"</str:features>" +
			"<iq xmlns='jabber:client' type='something'></iq>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config: Config{
			SkipTLS: true,
			CreateCallback: func(title, instructions string, fields []interface{}) error {
				return nil
			},
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "xmpp: account creation failed")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_failsIfTheIQQueryHasNoContent(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"<register xmlns='http://jabber.org/features/iq-register'/>" +
			"</str:features>" +
			"<iq xmlns='jabber:client' type='result'></iq>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config: Config{
			SkipTLS: true,
			CreateCallback: func(title, instructions string, fields []interface{}) error {
				return nil
			},
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_ifRegisterQueryDoesntContainDataFailsAtNextIQ(c *C) {
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
			"<query xmlns='jabber:iq:register'></query>" +
			"</iq>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config: Config{
			SkipTLS: true,
			CreateCallback: func(title, instructions string, fields []interface{}) error {
				return nil
			},
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Matches, "unmarshal <iq>:( XML syntax error on line 1: unexpected)? EOF")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_afterRegisterFailsIfReceivesAnErrorElement(c *C) {
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
			"<query xmlns='jabber:iq:register'></query>" +
			"</iq>" +
			"<iq xmlns='jabber:client' type='error'></iq>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config: Config{
			SkipTLS: true,
			CreateCallback: func(title, instructions string, fields []interface{}) error {
				return nil
			},
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "xmpp: account creation failed")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_sendsBackUsernameAndPassword(c *C) {
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
			"<query xmlns='jabber:iq:register'><username/><password/></query>" +
			"</iq>" +
			"<iq xmlns='jabber:client' type='result'></iq>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config: Config{
			SkipTLS: true,
			CreateCallback: func(title, instructions string, fields []interface{}) error {
				return nil
			},
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, IsNil)
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>"+
		"<iq type='set' id='create_2'><query xmlns='jabber:iq:register'><username>user</username><password>pass</password></query></iq>"+
		"</stream:stream>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_runsForm(c *C) {
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
			"<query xmlns='jabber:iq:register'>" +
			"<x xmlns='jabber:x:data' type='form'>" +
			"<title>Contest Registration</title>" +
			"<field type='hidden' var='FORM_TYPE'>" +
			"<value>jabber:iq:register</value>" +
			"</field>" +
			"<field type='text-single' label='Given Name' var='first'>" +
			"<required/>" +
			"</field>" +
			"</x>" +
			"</query>" +
			"</iq>" +
			"<iq xmlns='jabber:client' type='result'></iq>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config: Config{
			SkipTLS: true,
			CreateCallback: func(title, instructions string, fields []interface{}) error {
				return nil
			},
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, IsNil)
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>"+
		"<iq type='set' id='create_2'><query xmlns='jabber:iq:register'><x xmlns=\"jabber:x:data\" type=\"submit\"><field var=\"FORM_TYPE\"><value>jabber:iq:register</value></field><field var=\"first\"><value></value></field></x></query></iq>"+
		"</stream:stream>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_setsLog(c *C) {
	l := &mockConnIOReaderWriter{}
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"<register xmlns='http://jabber.org/features/iq-register'/>" +
			"</str:features>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config: Config{
			SkipTLS: true,
			Log:     l,
			CreateCallback: func(title, instructions string, fields []interface{}) error {
				return nil
			},
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Matches, "unmarshal <iq>:( XML syntax error on line 1: unexpected)? EOF")
	c.Assert(string(l.write), Equals, "Attempting to create account\n")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_failsWhenTryingToEstablishSession(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>" +
			"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<str:session>foobar</str:session>" +
			"</str:features>" +
			"<client:iq xmlns:client='jabber:client'></client:iq>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config: Config{
			SkipTLS: true,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Matches, "xmpp: unmarshal <iq>:( XML syntax error on line 1: unexpected)? EOF")

	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq>"+
		"<iq to='domain' type='set' id='sess_1'><session xmlns='urn:ietf:params:xml:ns:xmpp-session'/></iq>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_failsWhenTryingToEstablishSessionAndGetsTheWrongIQBack(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>" +
			"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<str:session>foobar</str:session>" +
			"</str:features>" +
			"<client:iq xmlns:client='jabber:client'></client:iq>" +
			"<client:iq xmlns:client='jabber:client' type='foo'></client:iq>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config: Config{
			SkipTLS: true,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "xmpp: session establishment failed")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq>"+
		"<iq to='domain' type='set' id='sess_1'><session xmlns='urn:ietf:params:xml:ns:xmpp-session'/></iq>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_succeedsEstablishingASession(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>" +
			"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<str:session>foobar</str:session>" +
			"</str:features>" +
			"<client:iq xmlns:client='jabber:client'></client:iq>" +
			"<client:iq xmlns:client='jabber:client' type='result'></client:iq>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &Dialer{
		JID:      "user@domain",
		Password: "pass",
		Config: Config{
			SkipTLS: true,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, IsNil)
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq>"+
		"<iq to='domain' type='set' id='sess_1'><session xmlns='urn:ietf:params:xml:ns:xmpp-session'/></iq>",
	)
}

// func (s *ConnectionXmppSuite) Test_blaData(c *C) {
// 	println("Trying!")
// 	var tlsC tls.Config
// 	tlsC.ServerName = "www.olabini.se"
// 	tlsC.InsecureSkipVerify = true
// 	tlsC.Rand = fixedRand([]string{
// 		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
// 		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
// 		"000102030405060708090A0B0C0D0E0F",
// 		"000102030405060708090A0B0C0D0E0F",
// 	})
// 	conn, _ := net.setupStream("tcp", "www.olabini.se:443")
// 	tee := createTeeConn(conn, os.Stdout)
// 	tlsConn := tls.Client(tee, &tlsC)
// 	err := tlsConn.Handshake()
// 	if err != nil {
// 		println("Error: ", err.Error())
// 	}
// }

var validTLSExchange = [][]byte{
	[]byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'>" +
			"</starttls>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<proceed xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>"),
	bytesFromHex("160303005d020000590303561976181b6028026857f7c44a881c5b76ecf47207a0aa4c5137160d3d9edb6520eb5c272ac25b34b590bd1a2972bebe34be767cfe80cbdfb6d8f2b1a6764f3d07c01400001100000000ff01000100000b00040300010216030310240b00102000101d0009df308209db308208c3a0030201020203022008300d06092a864886f70d010105050030818c310b300906035504061302494c31163014060355040a130d5374617274436f6d204c74642e312b3029060355040b1322536563757265204469676974616c204365727469666963617465205369676e696e67313830360603550403132f5374617274436f6d20436c6173732032205072696d61727920496e7465726d65646961746520536572766572204341301e170d3134303432313130313531335a170d3136303432313131313131395a30819d31193017060355040d13104e6a4d4833335634776757646b34614a310b30090603550406130255533111300f06035504081308496c6c696e6f69733110300e060355040713074368696361676f3111300f060355040a13084f6c612042696e69311530130603550403140c2a2e6f6c6162696e692e73653124302206092a864886f70d0109011615706f73746d6173746572406f6c6162696e692e736530820222300d06092a864886f70d01010105000382020f003082020a0282020100c30142b1459535ab3365827e3f4a654a87f11b96e307998ba6f799ab490d81e890f1121ce092bb0aae969d801510493e7e280356a5a0f888cb3ccb5abab5104c5630c7c1d88a0525948a6ba283b85f5a88365cb7dd1eba93b21d83a40a8fe0b9848677d21367e01b6dc75959e6a50e4ca5e0c2f9a2f93ddf2db9de7ef70df4be7ede241167d02208fbc5e8924f8ced4ce71bff966e80a9f8159a1cc466814f037c1434c7ba7c663c4101ecb1186c6f5828aaf5d562ee3e4c2a97e1a92f64a76fc25332d5e5a71a1d49eb0291b0faca2644333ca22cac55c5c7f7b3275c0801a22f0fc95a02ac13e947ed91f542ad6cbe580cceb0862940d8b558403614cd4a5d094e51db04eec944809a0b4d8da74cae3cbff9af3017bc87a66a84b9688d871474d89e7dbecedfcd246dbd58471272618ec1959a4923958639625476aad24832b692f88b4cedddf03b6aaf6be83ee0ca7e9f4dfbc55d589643f3113a86ec5e4bbb0cad357ea6d9634888870503da4e3d842dd05c5d9328393570660551361847b6dde4a54bb8452a4a6746fca4fd83f7556f39fea8a5d6339fc8ba1f967f687f3c6ae58e7be4a5459eb42ccaa0090f075202c8bb37d7505002a0bf66e3361f73697e825d20351e4c70bbaa1f8c4142f1b4b7d1b0c8adacf03ec69af7e8f8ab9d576f436f6ea51d60681d858ae777ee2eaa71fa3a53475737aaea71b1c15d"),
	bytesFromHex("04110203010001a38205313082052d30090603551d1304023000300b0603551d0f0404030203a8301d0603551d250416301406082b0601050507030206082b06010505070301301d0603551d0e0416041443d7a7c90643825a3150e3f331d4e1c9ba96174b301f0603551d2304183016801411db2345fd54cc6a716f848a03d7bef7012f26863082026b0603551d11048202623082025e820c2a2e6f6c6162696e692e7365820a6f6c6162696e692e736582092a2e62696e692e65638208696f6b652e6f7267820a2a2e696f6b652e6f726782106d616c696e73616e64656c6c2e636f6d82122a2e6d616c696e73616e64656c6c2e636f6d820f6d616c696e73616e64656c6c2e6e7582112a2e6d616c696e73616e64656c6c2e6e75820f6d616c696e73616e64656c6c2e736582112a2e6d616c696e73616e64656c6c2e7365820b6f6c6162696e692e636f6d820d2a2e6f6c6162696e692e636f6d820a6f6c6162696e692e6563820c2a2e6f6c6162696e692e6563820c6f6c6162696e692e696e666f820e2a2e6f6c6162696e692e696e666f820a6f6c6162696e692e6d65820c2a2e6f6c6162"),
	bytesFromHex("696e692e6d65820b6f6c6162696e692e6e6574820d2a2e6f6c6162696e692e6e6574820a6f6c6162696e692e6e75820c2a2e6f6c6162696e692e6e75820762696e692e6563820a6f6c6f6769782e6e6574820c2a2e6f6c6f6769782e6e657482096f6c6f6769782e7365820b2a2e6f6c6f6769782e73658207726561702e656382092a2e726561702e6563820a73616e64656c6c2e6563820c2a2e73616e64656c6c2e6563820d736570682d6c616e672e6f7267820f2a2e736570682d6c616e672e6f726782117374656c6c6173616e64656c6c2e636f6d82132a2e7374656c6c6173616e64656c6c2e636f6d82117374656c6c6173616e64656c6c2e6e657482132a2e7374656c6c6173616e64656c6c2e6e657482107374656c6c6173616e64656c6c2e736582122a2e7374656c6c6173616e64656c6c2e736582067477732e656382082a2e7477732e6563308201560603551d200482014d308201493008060667810c0102023082013b060b2b0601040181b5370102033082012a302e06082b060105050702011622687474703a2f2f7777772e737461727473736c2e636f6d2f706f6c6963792e7064663081f706082b060105050702023081ea302716205374617274436f6d2043657274696669636174696f6e20417574686f7269747930030201011a81be546869732063657274696669636174652077617320697373756564206163636f7264696e6720746f2074686520436c61737320322056616c69646174696f6e20726571756972656d656e7473206f6620746865205374617274436f6d20434120706f6c6963792c2072656c69616e6365206f6e6c7920666f722074686520696e74656e64656420707572706f736520696e20636f6d706c69616e6365206f66207468652072656c79696e67207061727479206f626c69676174696f6e732e30350603551d1f042e302c302aa028a0268624687474703a2f2f63726c2e737461727473736c2e636f6d2f637274322d63726c2e63726c30818e06082b06010505070101048181307f303906082b06010505073001862d687474703a2f2f6f6373702e737461727473736c2e636f6d2f7375622f636c617373322f7365727665722f6361304206082b060105050730028636687474703a2f2f6169612e737461727473736c2e636f6d2f63657274732f7375622e636c617373322e7365727665722e63612e63727430230603551d12041c301a8618687474703a2f2f7777772e737461727473736c2e636f6d2f300d06092a864886f70d01010505000382010100debc4553a1ee7e34e0d2aa9100e08865153aa3a9055bca798feb006b429b7985a2e89451f65171c704fc9576c20818da71cccd1eb987cb4e5da4af1a5a0097c1e64a995941ee612865ef1c9b12f70c9ecf54e4fb35bac20b5f317f9f878f964779361adb66638a53c6c72f1ebfac6b083e5f5329f89f060ebd7671e26caa53597e9cfc5511ad1e5ab4a75b7b3ed1d01f6167eefca0708c5fdf18d8df0fc6ad794e774475c3205d02eaa3d88347181a5bad96284e2e70028ef84f21b6dc5c71aa635c949a574cef75eba1cd37b4c6af09564f047a573b1dae41574821cb1484d362e6250cc8f912f2a8684213f6016f0335011140986b390d98821697e492dc20000638308206343082041ca00302010202011a300d06092a864886f70d0101050500307d310b300906035504061302494c31163014060355040a130d5374617274436f6d204c74642e312b3029060355040b1322536563757265204469676974616c204365727469666963617465205369676e696e6731293027060355040313205374617274436f6d2043657274696669636174696f6e20417574686f72697479301e170d3037313032343230353730395a170d3137313032343230353730395a30818c310b300906035504061302494c31163014060355040a130d5374617274436f6d204c74642e312b3029060355040b1322536563757265204469676974616c204365727469666963617465205369676e696e67313830360603550403132f5374617274436f6d20436c6173732032205072696d61727920496e7465726d6564696174652053657276657220434130820122300d06092a864886f70d01010105000382010f003082010a0282010100e24f392fa18c9a85ad080e083e57f28801211b94a96ce2b8dbaa1918463a52a1f50ff46e8cea968c9687791340512f22f20c8b870f65df7174344355b135099bd9bc1ffaeb42d0974072b743963dba969d5d50021c9b918d9cc0acd7bb2f17d7cb3e829d73eb074292b2cd64b374551bb44b86212cf7788732e016e4dabd4c95eaa40a7eb60a0d2e8acf55abc3e5dd418a4ee66f656cb240cf175db9c36a0b2711847761f6c27cedc08d7814189981997563b7e853d3ba61e90efaa230f346a2b9c91f6c805a40ac27ed484733b054c6461af33561c1022990547e644dc430520282d7dfce216e1891d7b8ab8c2717b5f0a3012f8ed22e873a3db429678ac4030203010001a38201ad308201a9300f0603551d130101ff040530030101ff300e0603551d0f0101ff040403020106301d0603551d0e0416041411db2345fd54cc6a716f848a03d7bef7012f2686301f0603551d230418301680144e0bef1aa4405ba517698730ca346843d041aef2306606082b06010505070101045a3058302706082b06010505073001861b687474703a2f2f6f6373702e737461727473736c2e636f6d2f6361302d06082b060105050730028621687474703a2f2f7777772e737461727473736c2e636f6d2f73667363612e637274305b0603551d1f045430523027a025a0238621687474703a2f2f7777772e737461727473736c2e636f6d2f73667363612e63726c3027a025a0238621687474703a2f2f63726c2e737461727473736c2e636f6d2f73667363612e63726c3081800603551d20047930773075060b2b0601040181b5370102013066302e06082b060105050702011622687474703a2f2f7777772e737461727473736c2e636f6d2f706f6c6963792e706466303406082b060105050702011628687474703a2f2f7777772e737461727473736c2e636f6d2f696e7465726d6564696174652e706466300d06092a864886f70d010105050003820201009d07e1ee907631671645708ccb848b4b576844a589c1f27ecb288bf5e77077d5b6f40b2160a5a17473242280d6d8ba8da2625d09354229fb3963450ba4b0381a68f49513cce04394eceb391aec5729d9996df584cd8e73aec9dc6afa9e9d16649308c71cc289549e778090f6b92976eb13674859f82e3a31b8c9d388e55f4ed2193d438ed792ffcf38b6e15b8a531dceacb4762fd8f74063d5ee69f3457da062c161c375edb27b4dac2127304e59466a9317cac8392d0173655be9419b11179cc8c84aefa176602dae93ff0cd533139f4f13cedd86f1fcf8355415a85be7857efa3709ff8bb831499e0d6edeb4d2122db8edc8c3f1b642a04c9779dffec3a39fa1f46d2c8477a4a205e117ff31dd9af3b87ac352c21111b750318a7fcce75a89ccf7869a61924f2f94b698c778e0624b437d3cded69ab410a1409c4b2adcb8d0d49efdf184781b0e578f695442687beaa0ef750f07a28c7399ab55f50709d2af38036a90030c2f8fe2e843c231e96fad87e58dbd4e2c894b51e69c4c5476c01281539beca0fc2c9cda18956e1e38264227786008df7f6d32e8d8c06f1feb26759f93fc7b1bfe3590dc53a307a63f83550a2b4e628225ce66305d2ce0f9191b75b99d9856a683277ad18f8d5993fc3f73d72eb42c95d88bf7c97ec7fc9dac72041fd2cc17f4ed34609b9e4a9704fedd720e57545106704defaa1ca482e033c7f4160303024d0c0002490300174104880c0195c2c45bbfc4e34bd0295dade14ce6d3b1d4630dee0dda7183affb9f84dd96dd672b8cefa666f0faba9c823217dcefa99d2ce2f66b03de4bf6b998682802010200271ec5f8fb61c62cbcd4ebfae7fc58f32b3db921275c308b9c67bc2a61fa669be124c3fc99005e5bc9f9601e1da32a36825ab0651857c0c80a7f87adfdc631103b4ab85bd60a6f6828c5fcb146f41621abb9bdd2b0a8743b6de852369d5c2624b8d7cd0c7dd64199d37bec1dd0146fd58c0496e8aa1bd1348eed35ce411a72d01faf2c917b4435340452497962ae66c2b87306c053b3e99f5dd798b2ec3c8c4c4b81fdf06a47b4905e46c9401f66f915ec4fe803779f8ca9c4a4a1af31b33a2982b180d80d542a5f8ca3a6a5fc3cc62196d4467d6483abe0b7fa24c266172b3fd31565a377f0546f055084b6b1d5fc82df31cac7f5114fe7ce7c802d7a880119a0a41421893a002b89ad9cbb81902b8b3bd636a7208a975e27c773e9146fb48a582b0981f4bda469e9e5f0d0b0660790971ac90d92638180094a68a448f6737622e9dbb717b72c876accd1e6d36123e6a6854a43597294e10d949a06886ba275ffdde903576fb770c1f9460c8010874870020eee832ebdb1cf4d87aa1b9f5aedddc152a611d946832e426aa58f0d9049a39c238f9444863fc48dc7ef69cd3e047560860223a5fff7d735169aba8fe6551b97f20fd075ac87f3a1d199698fdf605414da4b13f0dea8855eb52e539f8b986f5b33bfaa11241af212dd39c428441834939c4729267fd33ef66aaa67fcfffbb3946c7db2805715518980098a42568516030300040e000000"),
	bytesFromHex("1403030001011603030040672edb3cd5374ab5f61c9dc5f350e3222ab87a61edb8f177409dee1cd974472ce78defd749972cc86bf7a6a4917aab3639f0129427e8b5cd05c9a10ec8a88675"),
}

var validTLSExchange2 = [][]byte{
	[]byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'>" +
			"</starttls>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<proceed xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>"),
	bytesFromHex("160303005d020000590303561986f1db1afdfdc2458be105ffb6ea5ce74deabece1a6aafa93ec9e300b0ce208f60d1038ddd8930ea0e5a7d3bc52a91de884e0e368a2539c453a6f33b1d79acc03000001100000000ff01000100000b00040300010216030310240b00102000101d0009df308209db308208c3a0030201020203022008300d06092a864886f70d010105050030818c310b300906035504061302494c31163014060355040a130d5374617274436f6d204c74642e312b3029060355040b1322536563757265204469676974616c204365727469666963617465205369676e696e67313830360603550403132f5374617274436f6d20436c6173732032205072696d61727920496e7465726d65646961746520536572766572204341301e170d3134303432313130313531335a170d3136303432313131313131395a30819d31193017060355040d13104e6a4d4833335634776757646b34614a310b30090603550406130255533111300f06035504081308496c6c696e6f69733110300e060355040713074368696361676f3111300f060355040a13084f6c612042696e69311530130603550403140c2a2e6f6c6162696e692e73653124302206092a864886f70d0109011615706f73746d6173746572406f6c6162696e692e736530820222300d06092a864886f70d01010105000382020f003082020a0282020100c30142b1459535ab3365827e3f4a654a87f11b96e307998ba6f799ab490d81e890f1121ce092bb0aae969d801510493e7e280356a5a0f888cb3ccb5abab5104c5630c7c1d88a0525948a6ba283b85f5a88365cb7dd1eba93b21d83a40a8fe0b9848677d21367e01b6dc75959e6a50e4ca5e0c2f9a2f93ddf2db9de7ef70df4be7ede241167d02208fbc5e8924f8ced4ce71bff966e80a9f8159a1cc466814f037c1434c7ba7c663c4101ecb1186c6f5828aaf5d562ee3e4c2a97e1a92f64a76fc25332d5e5a71a1d49eb0291b0faca2644333ca22cac55c5c7f7b3275c0801a22f0fc95a02ac13e947ed91f542ad6cbe580cceb0862940d8b558403614cd4a5d094e51db04eec944809a0b4d8da74cae3cbff9af3017bc87a66a84b9688d871474d89e7dbecedfcd246dbd58471272618ec1959a4923958639625476aad24832b692f88b4cedddf03b6aaf6be83ee0ca7e9f4dfbc55d589643f3113a86ec5e4bbb0cad357ea6d9634888870503da4e3d842dd05c5d9328393570660551361847b6dde4a54bb8452a4a6746fca4fd83f7556f39fea8a5d6339fc8ba1f967f687f3c6ae58e7be4a5459eb42ccaa0090f075202c8bb37d7505002a0bf66e3361f73697e825d20351e4c70bbaa1f8c4142f1b4b7d1b0c8adacf03ec69af7e8f8ab9d576f436f6ea51d60681d858ae777ee2eaa71fa3a53475737aaea71b1c15d"),
	bytesFromHex("04110203010001a38205313082052d30090603551d1304023000300b0603551d0f0404030203a8301d0603551d250416301406082b0601050507030206082b06010505070301301d0603551d0e0416041443d7a7c90643825a3150e3f331d4e1c9ba96174b301f0603551d2304183016801411db2345fd54cc6a716f848a03d7bef7012f26863082026b0603551d11048202623082025e820c2a2e6f6c6162696e692e7365820a6f6c6162696e692e736582092a2e62696e692e65638208696f6b652e6f7267820a2a2e696f6b652e6f726782106d616c696e73616e64656c6c2e636f6d82122a2e6d616c696e73616e64656c6c2e636f6d820f6d616c696e73616e64656c6c2e6e7582112a2e6d616c696e73616e64656c6c2e6e75820f6d616c696e73616e64656c6c2e736582112a2e6d616c696e73616e64656c6c2e7365820b6f6c6162696e692e636f6d820d2a2e6f6c6162696e692e636f6d820a6f6c6162696e692e6563820c2a2e6f6c6162696e692e6563820c6f6c6162696e692e696e666f820e2a2e6f6c6162696e692e696e666f820a6f6c6162696e692e6d65820c2a2e6f6c6162696e692e6d65820b6f6c6162696e692e6e6574820d2a2e6f6c6162696e692e6e6574820a6f6c6162696e692e6e75820c2a2e6f6c6162696e692e6e75820762696e692e6563820a6f6c6f6769782e6e6574820c2a2e6f6c6f6769782e6e657482096f6c6f6769782e7365820b2a2e6f6c6f6769782e73658207726561702e656382092a2e726561702e6563820a73616e64656c6c2e6563820c2a2e73616e64656c6c2e6563820d736570682d6c616e672e6f7267820f2a2e736570682d6c616e672e6f726782117374656c6c6173616e64656c6c2e636f6d82132a2e7374656c6c6173616e64656c6c2e636f6d82117374656c6c6173616e64656c6c2e6e657482132a2e7374656c6c6173616e64656c6c2e6e657482107374656c6c6173616e64656c6c2e736582122a2e7374656c6c6173616e64656c6c2e736582067477732e656382082a2e7477732e6563308201560603551d200482014d308201493008060667810c0102023082013b060b2b0601040181b5370102033082012a302e06082b060105050702011622687474703a2f2f7777772e737461727473736c2e636f6d2f706f6c6963792e7064663081f706082b060105050702023081ea302716205374617274436f6d2043657274696669636174696f6e20417574686f7269747930030201011a81be546869732063657274696669636174652077617320697373756564206163636f7264696e6720746f2074686520436c61737320322056616c69646174696f6e20726571756972656d656e7473206f6620746865205374617274436f6d20434120706f6c6963792c2072656c69616e6365206f6e6c7920666f722074686520696e74656e64656420707572706f736520696e20636f6d706c69616e6365206f66207468652072656c79696e67207061727479206f626c69676174696f6e732e30350603551d1f042e302c302aa028a0268624687474703a2f2f63726c2e737461727473736c2e636f6d2f637274322d63726c2e63726c30818e06082b06010505070101048181307f303906082b06010505073001862d687474703a2f2f6f6373702e737461727473736c2e636f6d2f7375622f636c617373322f7365727665722f6361304206082b060105050730028636687474703a2f2f6169612e737461727473736c2e636f6d2f63657274732f7375622e636c617373322e7365727665722e63612e63727430230603551d12041c301a8618687474703a2f2f7777772e737461727473736c2e636f6d2f300d06092a864886f70d01010505000382010100debc4553a1ee7e34e0d2aa9100e08865153aa3a9055bca798feb006b429b7985a2e89451f65171c704fc9576c20818da71cccd1eb987cb4e5da4af1a5a0097c1e64a995941ee612865ef1c9b12f70c9ecf54e4fb35bac20b5f317f9f878f964779361adb66638a53c6c72f1ebfac6b083e5f5329f89f060ebd7671e26caa53597e9cfc5511ad1e5ab4a75b7b3ed1d01f6167eefca0708c5fdf18d8df0fc6ad794e774475c3205d02eaa3d88347181a5bad96284e2e70028ef84f21b6dc5c71aa635c949a574cef75eba1cd37b4c6af09564f047a573b1dae41574821cb1484d362e6250cc8f912f2a8684213f6016f0335011140986b390d98821697e492dc20000638308206343082041ca00302010202011a300d06092a864886f70d0101050500307d310b300906035504061302494c31163014060355040a130d5374617274436f6d204c74642e312b3029060355040b1322536563757265204469676974616c204365727469666963617465205369676e696e6731293027060355040313205374617274436f6d2043657274696669636174696f6e20417574686f72697479301e170d3037313032343230353730395a170d3137313032343230353730395a30818c310b300906035504061302494c31163014060355040a130d5374617274436f6d204c74642e312b3029060355040b1322536563757265204469676974616c204365727469666963617465205369676e696e67313830360603550403132f5374617274436f6d20436c6173732032205072696d61727920496e7465726d6564696174652053657276657220434130820122300d06092a864886f70d01010105000382010f003082010a0282010100e24f392fa18c9a85ad080e083e57f28801211b94a96ce2b8dbaa1918463a52a1f50ff46e8cea968c9687791340512f22f20c8b870f65df7174344355b135099bd9bc1ffaeb42d0974072b743963dba969d5d50021c9b918d9cc0acd7bb2f17d7cb3e829d73eb074292b2cd64b374551bb44b86212cf7788732e016e4dabd4c95eaa40a7eb60a0d2e8acf55abc3e5dd418a4ee66f656cb240cf175db9c36a0b2711847761f6c27cedc08d7814189981997563b7e853d3ba61e90efaa230f346a2b9c91f6c805a40ac27ed484733b054c6461af33561c1022990547e644dc430520282d7dfce216e1891d7b8ab8c2717b5f0a3012f8ed22e873a3db429678ac4030203010001a38201ad308201a9300f0603551d130101ff040530030101ff300e0603551d0f0101ff040403020106301d0603551d0e0416041411db2345fd54cc6a716f848a03d7bef7012f2686301f0603551d230418301680144e0bef1aa4405ba517698730ca346843d041aef2306606082b06010505070101045a3058302706082b06010505073001861b687474703a2f2f6f6373702e737461727473736c2e636f6d2f6361302d06082b060105050730028621687474703a2f2f7777772e737461727473736c2e636f6d2f73667363612e637274305b0603551d1f045430523027a025a0238621687474703a2f2f7777772e737461727473736c2e636f6d2f73667363612e63726c3027a025a0238621687474703a2f2f63726c2e737461727473736c2e636f6d2f73667363612e63726c3081800603551d20047930773075060b2b0601040181b5370102013066302e06082b060105050702011622687474703a2f2f7777772e737461727473736c2e636f6d2f706f6c6963792e706466303406082b060105050702011628687474703a2f2f7777772e737461727473736c2e636f6d2f696e7465726d6564696174652e706466300d06092a864886f70d010105050003820201009d07e1ee907631671645708ccb848b4b576844a589c1f27ecb288bf5e77077d5b6f40b2160a5a17473242280d6d8ba8da2625d09354229fb3963450ba4b0381a68f49513cce04394eceb391aec5729d9996df584cd8e73aec9dc6afa9e9d16649308c71cc289549e778090f6b92976eb13674859f82e3a31b8c9d388e55f4ed2193d438ed792ffcf38b6e15b8a531dceacb4762fd8f74063d5ee69f3457da062c161c375edb27b4dac2127304e59466a9317cac8392d0173655be9419b11179cc8c84aefa176602dae93ff0cd533139f4f13cedd86f1fcf8355415a85be7857efa3709ff8bb831499e0d6edeb4d2122db8edc8c3f1b642a04c9779dffec3a39fa1f46d2c8477a4a205e117ff31dd9af3b87ac352c21111b750318a7fcce75a89ccf7869a61924f2f94b698c778e0624b437d3cded69ab410a1409c4b2adcb8d0d49efdf184781b0e578f695442687beaa0ef750f07a28c7399ab55f50709d2af38036a90030c2f8fe2e843c231e96fad87e58dbd4e2c894b51e69c4c5476c01281539beca0fc2c9cda18956e1e38264227786008df7f6d32e8d8c06f1feb26759f93fc7b1bfe3590dc53a307a63f83550a2b4e628225ce66305d2ce0f9191b75b99d9856a683277ad18f8d5993fc3f73d72eb42c95d88bf7c97ec7fc9dac72041fd2cc17f4ed34609b9e4a9704fedd720e57545106704defaa1ca482e033c7f4160303024d0c000249030017410458fb4cf2e793e1da9e80cb9495bc8b8e699f4a0f6e7c361f01616f90cc47a7960bb190775f9725ea45eab1a7b8d982fa9433faea7fd731524b3cafa27870cd0c0201020015e77bb80a24b5a267e2cc3683c0b5477ad572d9d48fc16518e53a33eb9654f18b84d3b1ed74353606130be3d7707c0413922b2d989a2fcfac9f8ce63a48dfa4255e78181914f3d1588e79d13f3cfef213f8329fe8ad170a2d89bb93b247947fe361aac987782a47adfa06afe23b11475f309cfd0690daf86fc89bc305831b107b1d299e01b4b38e626cf126389532e838b7f57168e93ba655672367e0d5b6a4170dec1fa0f9ec26765c28976fa83422d5a28c3e99e0db500fe16d1eeb285f7f97c2f5f3c1cef5a34622e9a3b29552826c5e99889c5f2e7b323aabdc76eb0a54cc0755be63d6c952b6c68d04dda38a0a17ad9e891d1bac8d9bba4caec8d9c5e66ff53c789bbf3caa05fa928fead38c2271f71adca6706ae9de46aae8a31be97af6618d5e9223fcf9141651b53e99d2d142c2559f102b9ef28f948aec94b7dba85095f517b26481b8c75014b1315fae8ded43a9af117d87de375e6f1df736973aee1dfce71f5392b00c22139daa21eb0b29d2b424cacf251484bbfd1c05c9bb10571947292d19bc7c657250f3f1c937c3219ac79ad51caa472f4ef61a355fe299b199995c4b6476816821ffb0554783d9200eb12b891b4483c8a15c8011fa173c2d0ec3a4ffcd414f36e99ac82379bda42fa8ceae269ac008724f8d4b6e9ed8b330d597fad33ec1862ba6f700c5506cddefa1ac66fe262897f3fda961e5c140d216030300040e000000"),
	bytesFromHex("14030300010116030300283867a1b1ef4eaac9acda2164189ece901ff133bf9f98480f8a3aac09d26b5bd1b5ac929e10dc3209"),
}

func isVersionOldish() bool {
	v := runtime.Version()
	return strings.HasPrefix(v, "go1.3") ||
		strings.HasPrefix(v, "go1.4")
}

func decideTLSExchangeFromVersion() [][]byte {
	if isVersionOldish() {
		return validTLSExchange
	}
	return validTLSExchange2
}

func (s *ConnectionXmppSuite) Test_Dial_worksIfTheHandshakeSucceeds(c *C) {
	rw := &mockMultiConnIOReaderWriter{read: decideTLSExchangeFromVersion()}
	conn := &fullMockedConn{rw: rw}
	var tlsC tls.Config
	tlsC.Rand = fixedRand([]string{
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F",
		"000102030405060708090A0B0C0D0E0F",
	})

	d := &Dialer{
		JID:           "user@www.olabini.se",
		Password:      "pass",
		ServerAddress: "www.olabini.se:443",

		Config: Config{
			TLSConfig: &tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
	if isVersionOldish() {
		c.Assert(string(rw.write), Equals, ""+
			"<?xml version='1.0'?><stream:stream to='www.olabini.se' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>\x16\x03\x01\x00\x8c\x01\x00\x00\x88\x03\x03\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
			"\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x00\x00\x1a\xc0/\xc0+\xc0\x11\xc0\a\xc0\x13\xc0\t\xc0\x14\xc0\n"+
			"\x00\x05\x00/\x005\xc0\x12\x00\n"+
			"\x01\x00\x00E\x00\x00\x00\x13\x00\x11\x00\x00\x0ewww.olabini.se\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\n"+
			"\x00\b\x00\x06\x00\x17\x00\x18\x00\x19\x00\v\x00\x02\x01\x00\x00\r\x00\n"+
			"\x00\b\x04\x01\x04\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00\x16\x03\x03\x00F\x10\x00\x00BA\x04\b\xda\xf8\xb0\xab\xae5\t\xf3\\\xe1\xd31\x04\xcb\x01\xb9Qb̹\x18\xba\x1f\x81o8\xd3\x13\x0f\xb8\u007f\x92\xa3\b7\xf8o\x9e\xef\x19\u007fCy\xa5\n"+
			"b\x06\x82fy]\xb9\xf83\xea6\x1d\x03\xafT[\xe7\x92\x14\x03\x03\x00\x01\x01\x16\x03\x03\x00@\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
			"\v\f\r\x0e\x0fҟ\xb0>ʋ,\xcbI\x9e\x1c\x94\b\x1an\x12cT\x81\xac'\xf4{\rV\xb1V\xad]\xc5\b\xfe\xf4rh˱\b?\x10&ե\x89\"\xbf\x8a0\x17\x03\x03\x00\xc0\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
			"\v\f\r\x0e\x0f7ѱ\n"+
			"9[\x9d\\!\xcd4\xfa\x18\xd6y\xb0y\xf2\x94T\x96\xb3\x9e\x16\xfdŷ\x9a\xf3Ʒ\xa9n@\xa5\tO\x8cE\x84\x85\x15#\x04\x97\xce\xfc\x12\x93;\xd8\xcdl>Q\xe52\x9f\xc1\x84\xf2cj\x81_U\x86\xcf6\xadڦC\xe6\x13\xfa\xc8%-\x15\xd5E\xe2h\x91\xfc\xa0+\x94\x13\xe3gG\xe6\xff\xe5`'و\x9fM\xea\xc780N\xd5\u9fdc\xf5)\xfdݫ3\xf0\x8ee/`8@\x88t\xdfc\xd3d-\xa9S\x80\x1a\x95rV\x98\x1e\xb2\xde\xdc\xd8\xc8\n"+
			"\x9aE\x12\xd1U\x95\xd37\xfe\xccm\u007f \x1a\xb5\x91\xd80;\\S\xb0\x04\x85Av+\xefǗ")
	} else {
		c.Assert(string(rw.write), Equals, ""+
			"<?xml version='1.0'?><stream:stream to='www.olabini.se' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>\x16\x03\x01\x00\x92\x01\x00\x00\x8e\x03\x03\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
			"\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x00\x00\x18\xc0/\xc0+\xc00\xc0,\xc0\x13\xc0\t\xc0\x14\xc0\n"+
			"\x00/\x005\xc0\x12\x00\n"+
			"\x01\x00\x00M\x00\x00\x00\x13\x00\x11\x00\x00\x0ewww.olabini.se\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\n"+
			"\x00\b\x00\x06\x00\x17\x00\x18\x00\x19\x00\v\x00\x02\x01\x00\x00\r\x00\x0e\x00\f\x04\x01\x04\x03\x05\x01\x05\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00\x00\x12\x00\x00\x16\x03\x03\x00F\x10\x00\x00BA\x04\b\xda\xf8\xb0\xab\xae5\t\xf3\\\xe1\xd31\x04\xcb\x01\xb9Qb̹\x18\xba\x1f\x81o8\xd3\x13\x0f\xb8\u007f\x92\xa3\b7\xf8o\x9e\xef\x19\u007fCy\xa5\n"+
			"b\x06\x82fy]\xb9\xf83\xea6\x1d\x03\xafT[\xe7\x92\x14\x03\x03\x00\x01\x01\x16\x03\x03\x00(\x00\x00\x00\x00\x00\x00\x00\x00漡noB\xe7z\x0f\xdb\xf7TxAдeϳ(\xc9\xe1m\xdd\xe6\x05\xa3~\x9f\tQ<\x17\x03\x03\x00\xa5\x00\x00\x00\x00\x00\x00\x00\x01Xa\x06\x8c;+\xbf\xdc\x12d<\x13'h\x9fk\xff\xea\xf60$~\x9e\x17ߕ\xb0J\xe6&\x90r\x1e\xb0\x18\xc3\r\xbbl\xe3\x12]\xb6\xa7X\xfc\xa2N̚`v-\xaeʐ\xbaT\xa7N\u007f\n"+
			"\xbc\xf7\x045\xc5\x03^\x9e\x92Fؼ!2\xe2.\xb2\xcbD\xfc\x91\x05\x86Ӽ\x16\xf4\x122ǼX\x8e\x8b\xb6t\x15\x9e<\xca\"\xd1\x19\xda\xd6\x06x\x98yx\x8b4ǚz\x99\xf3\xc1k\x89\xb0\xc0\xf6\xaf\xae\x92!9:hK\xf9j\a͕\xcc˱\xc1F@F{\x00\xe0{\xb5\x92y`\x98\bֽ")
	}
}

func (s *ConnectionXmppSuite) Test_Dial_worksIfTheHandshakeSucceedsButFailsOnInvalidCertHash(c *C) {
	rw := &mockMultiConnIOReaderWriter{read: decideTLSExchangeFromVersion()}
	conn := &fullMockedConn{rw: rw}
	var tlsC tls.Config
	tlsC.Rand = fixedRand([]string{
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F",
		"000102030405060708090A0B0C0D0E0F",
	})

	d := &Dialer{
		JID:           "user@www.olabini.se",
		Password:      "pass",
		ServerAddress: "www.olabini.se:443",

		Config: Config{
			TLSConfig:               &tlsC,
			ServerCertificateSHA256: []byte("aaaaa"),
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "xmpp: server certificate does not match expected hash (got: 2300818fdc977ce5eb357694d421e47869a952990bc3230ef6aca2bb6ee6f00b, want: 6161616161)")
}

func (s *ConnectionXmppSuite) Test_Dial_worksIfTheHandshakeSucceedsButSucceedsOnValidCertHash(c *C) {
	rw := &mockMultiConnIOReaderWriter{read: decideTLSExchangeFromVersion()}
	conn := &fullMockedConn{rw: rw}
	var tlsC tls.Config
	tlsC.Rand = fixedRand([]string{
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F",
		"000102030405060708090A0B0C0D0E0F",
	})

	d := &Dialer{
		JID:           "user@www.olabini.se",
		Password:      "pass",
		ServerAddress: "www.olabini.se:443",

		Config: Config{
			TLSConfig:               &tlsC,
			ServerCertificateSHA256: bytesFromHex("2300818fdc977ce5eb357694d421e47869a952990bc3230ef6aca2bb6ee6f00b"),
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXmppSuite) Test_readMessages_passesStanzaToChannel(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:message xmlns:client='jabber:client' to='fo@bar.com' from='bar@foo.com' type='chat'><client:body>something</client:body></client:message>")}

	conn := &Conn{
		in:     xml.NewDecoder(mockIn),
		closed: true, //This avoids trying to close the connection after the EOF
	}
	stanzaChan := make(chan data.Stanza)
	go conn.ReadStanzas(stanzaChan)

	select {
	case rawStanza, ok := <-stanzaChan:
		c.Assert(ok, Equals, true)
		c.Assert(rawStanza.Name.Local, Equals, "message")
		c.Assert(rawStanza.Value.(*data.ClientMessage).Body, Equals, "something")
	}
}

func (s *ConnectionXmppSuite) Test_readMessages_alertsOnError(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<clientx:message xmlns:client='jabber:client' to='fo@bar.com' from='bar@foo.com' type='chat'><client:body>something</client:body></client:message>")}

	conn := &Conn{
		in:     xml.NewDecoder(mockIn),
		closed: true, //This avoids trying to close the connection after the EOF
	}

	stanzaChan := make(chan data.Stanza, 1)
	err := conn.ReadStanzas(stanzaChan)

	select {
	case _, ok := <-stanzaChan:
		c.Assert(ok, Equals, false)
	}

	c.Assert(err.Error(), Equals, "unexpected XMPP message clientx <message/>")
}
