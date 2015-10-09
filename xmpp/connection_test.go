package xmpp

import (
	"crypto/tls"
	"encoding/xml"
	"io"

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
	reply := make(chan Stanza, 1)
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
	reply := make(chan Stanza, 1)
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
	reply := make(chan Stanza, 1)
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
	reply := make(chan Stanza, 1)
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
	reply := make(chan Stanza, 1)
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

func (s *ConnectionXmppSuite) Test_makeInOut_returnsANewDecoderAndOriginalWriterWhenNoConfigIsGiven(c *C) {
	mockBoth := &mockConnIOReaderWriter{}
	_, rout := makeInOut(mockBoth, nil)
	c.Assert(rout, Equals, mockBoth)
}

func (s *ConnectionXmppSuite) Test_makeInOut_returnsANewDecoderAndWrappedWriterWhenConfigIsGiven(c *C) {
	mockBoth := &mockConnIOReaderWriter{}
	mockInLog := &mockConnIOReaderWriter{}
	config := &Config{InLog: mockInLog, OutLog: mockInLog}
	_, rout := makeInOut(mockBoth, config)
	c.Assert(rout, Not(Equals), mockBoth)
}

func (s *ConnectionXmppSuite) Test_Dial_returnsErrorFromGetFeatures(c *C) {
	rw := &mockConnIOReaderWriter{}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXmppSuite) Test_Dial_returnsErrorFromAuthenticateIfSkipTLS(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'><str:features></str:features>")}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "xmpp: PLAIN authentication is not an option")
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
	config := &Config{Conn: conn, SkipTLS: true}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err, Equals, io.EOF)
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
	config := &Config{Conn: conn, SkipTLS: true}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "unmarshal <iq>: EOF")
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
	config := &Config{Conn: conn, SkipTLS: true}
	_, err := Dial("addr", "user", "domain", "pass", config)
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
	config := &Config{Conn: conn, SkipTLS: false}
	_, err := Dial("addr", "user", "domain", "pass", config)
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
	config := &Config{Conn: conn, SkipTLS: false}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err, Equals, io.EOF)
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
	config := &Config{Conn: conn, SkipTLS: false}
	_, err := Dial("addr", "user", "domain", "pass", config)
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
	config := &Config{Conn: conn, SkipTLS: false}
	_, err := Dial("addr", "user", "domain", "pass", config)
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
	config := &Config{Conn: conn, SkipTLS: false}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err, Equals, io.EOF)
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
	config := &Config{Conn: conn, SkipTLS: false, TLSConfig: &tlsC}
	_, err := Dial("addr", "user", "domain", "pass", config)
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
			"</str:features>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true, CreateCallback: func(title, instructions string, fields []interface{}) error {
		return nil
	}}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "unmarshal <iq>: EOF")
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
			"</str:features>" +
			"<iq xmlns='jabber:client' type='something'></iq>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true, CreateCallback: func(title, instructions string, fields []interface{}) error {
		return nil
	}}
	_, err := Dial("addr", "user", "domain", "pass", config)
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
			"</str:features>" +
			"<iq xmlns='jabber:client' type='result'></iq>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true, CreateCallback: func(title, instructions string, fields []interface{}) error {
		return nil
	}}
	_, err := Dial("addr", "user", "domain", "pass", config)
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
			"</str:features>" +
			"<iq xmlns='jabber:client' type='result'>" +
			"<query xmlns='jabber:iq:register'></query>" +
			"</iq>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true, CreateCallback: func(title, instructions string, fields []interface{}) error {
		return nil
	}}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "unmarshal <iq>: EOF")
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
			"</str:features>" +
			"<iq xmlns='jabber:client' type='result'>" +
			"<query xmlns='jabber:iq:register'></query>" +
			"</iq>" +
			"<iq xmlns='jabber:client' type='error'></iq>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true, CreateCallback: func(title, instructions string, fields []interface{}) error {
		return nil
	}}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "xmpp: account creation failed")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_continuesWithAuthenticationAfterRegistering(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<iq xmlns='jabber:client' type='result'>" +
			"<query xmlns='jabber:iq:register'><username/></query>" +
			"</iq>" +
			"<iq xmlns='jabber:client' type='result'></iq>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true, CreateCallback: func(title, instructions string, fields []interface{}) error {
		return nil
	}}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "expected <success> or <failure>, got <> in ")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_continuesWithAuthenticationAfterRegistering2(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<iq xmlns='jabber:client' type='result'>" +
			"<query xmlns='jabber:iq:register'><password/></query>" +
			"</iq>" +
			"<iq xmlns='jabber:client' type='result'></iq>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true, CreateCallback: func(title, instructions string, fields []interface{}) error {
		return nil
	}}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "expected <success> or <failure>, got <> in ")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n",
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
			"</str:features>" +
			"<iq xmlns='jabber:client' type='result'>" +
			"<query xmlns='jabber:iq:register'><username/><password/></query>" +
			"</iq>" +
			"<iq xmlns='jabber:client' type='result'></iq>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true, CreateCallback: func(title, instructions string, fields []interface{}) error {
		return nil
	}}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "expected <success> or <failure>, got <> in ")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>"+
		"<iq type='set' id='create_2'><query xmlns='jabber:iq:register'><username>user</username><password>pass</password></query></iq>"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n",
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
	config := &Config{Conn: conn, SkipTLS: true, CreateCallback: func(title, instructions string, fields []interface{}) error {
		return nil
	}}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "expected <success> or <failure>, got <> in ")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>"+
		"<iq type='set' id='create_2'><query xmlns='jabber:iq:register'><x xmlns=\"jabber:x:data\" type=\"submit\"><field var=\"FORM_TYPE\"><value>jabber:iq:register</value></field><field var=\"first\"><value></value></field></x></query></iq>"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n",
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
			"</str:features>",
	)}
	conn := &fullMockedConn{rw: rw}
	config := &Config{Conn: conn, SkipTLS: true, Log: l, CreateCallback: func(title, instructions string, fields []interface{}) error {
		return nil
	}}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "unmarshal <iq>: EOF")
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
	config := &Config{Conn: conn, SkipTLS: true}
	_, err := Dial("addr", "user", "domain", "pass", config)
	c.Assert(err.Error(), Equals, "xmpp: unmarshal <iq>: EOF")
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
	config := &Config{Conn: conn, SkipTLS: true}
	_, err := Dial("addr", "user", "domain", "pass", config)
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
	config := &Config{Conn: conn, SkipTLS: true}
	_, err := Dial("addr", "user", "domain", "pass", config)
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
