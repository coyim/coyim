package xmpp

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"io"

	goerr "errors"

	"github.com/twstrike/coyim/digests"
	"github.com/twstrike/coyim/xmpp/data"
	"github.com/twstrike/coyim/xmpp/errors"

	. "gopkg.in/check.v1"
)

type ConnectionXMPPSuite struct{}

var _ = Suite(&ConnectionXMPPSuite{})

type basicTLSVerifier struct {
	shaSum []byte
}

func (v *basicTLSVerifier) verifyCert(state tls.ConnectionState, conf tls.Config) ([][]*x509.Certificate, error) {
	opts := x509.VerifyOptions{
		Intermediates: x509.NewCertPool(),
		Roots:         conf.RootCAs,
	}

	for _, cert := range state.PeerCertificates[1:] {
		opts.Intermediates.AddCert(cert)
	}

	return state.PeerCertificates[0].Verify(opts)
}

func (v *basicTLSVerifier) verifyFailure(err error) error {
	return goerr.New("xmpp: failed to verify TLS certificate: " + err.Error())
}

func (v *basicTLSVerifier) verifyHostnameFailure(err error) error {
	return goerr.New("xmpp: failed to match TLS certificate to name: " + err.Error())
}

func (v *basicTLSVerifier) verifyHostName(leafCert *x509.Certificate, originDomain string) error {
	return leafCert.VerifyHostname(originDomain)
}

func (v *basicTLSVerifier) hasPinned(certs []*x509.Certificate) error {
	savedHash := v.shaSum
	if len(savedHash) == 0 {
		return nil
	}

	if digest := digests.Sha256(certs[0].Raw); !bytes.Equal(digest, savedHash) {
		return fmt.Errorf("tls: server certificate does not match expected hash (got: %x, want: %x)", digest, savedHash)
	}

	return nil
}

func (v *basicTLSVerifier) Verify(state tls.ConnectionState, conf tls.Config, originDomain string) error {
	if len(state.PeerCertificates) == 0 {
		return goerr.New("tls: server has no certificates")
	}

	if err := v.hasPinned(state.PeerCertificates); err != nil {
		return err
	}

	chains, err := v.verifyCert(state, conf)
	if err != nil {
		return v.verifyFailure(err)
	}

	if err = v.verifyHostName(chains[0][0], originDomain); err != nil {
		return v.verifyHostnameFailure(err)
	}

	return nil
}

func (s *ConnectionXMPPSuite) Test_Next_returnsErrorIfOneIsEncountered(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:foo xmlns:stream='http://etherx.jabber.org/streams' to='hello'></stream:foo>")}
	conn := conn{
		in: xml.NewDecoder(mockIn),
	}

	_, err := conn.Next()
	c.Assert(err.Error(), Equals, "unexpected XMPP message http://etherx.jabber.org/streams <foo/>")
}

func (s *ConnectionXMPPSuite) Test_Next_returnsErrorIfFailingToParseIQID(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='abczzzz'></client:iq>")}
	conn := conn{
		in: xml.NewDecoder(mockIn),
	}

	_, err := conn.Next()
	c.Assert(err.Error(), Equals, "xmpp: failed to parse id from iq: strconv.ParseUint: parsing \"abczzzz\": invalid syntax")
}

func (s *ConnectionXMPPSuite) Test_Next_returnsNothingIfThereIsNoInflightMatching(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000'></client:iq>")}
	conn := conn{
		in: xml.NewDecoder(mockIn),
	}

	_, err := conn.Next()
	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXMPPSuite) Test_Next_returnsNothingIfTheInflightIsToAnotherReceiver(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000' from='bar@somewhere.com'></client:iq>")}
	conn := conn{
		in:        xml.NewDecoder(mockIn),
		inflights: make(map[data.Cookie]inflight),
	}
	cookie := data.Cookie(1048576)
	conn.inflights[cookie] = inflight{to: "foo@somewhere.com"}
	_, err := conn.Next()
	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXMPPSuite) Test_Next_removesInflightIfItMatches(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000' from='foo@somewhere.com'></client:iq>")}
	inflights := make(map[data.Cookie]inflight)
	conn := conn{
		in:        xml.NewDecoder(mockIn),
		inflights: inflights,
	}
	cookie := data.Cookie(1048576)
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

func (s *ConnectionXMPPSuite) Test_Next_continuesIfIqFromIsNotSimilarToJid(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000' from='foo@somewhere.com'></client:iq>")}
	inflights := make(map[data.Cookie]inflight)
	conn := conn{
		in:        xml.NewDecoder(mockIn),
		inflights: inflights,
		jid:       "foo@myjid.com/blah",
	}
	cookie := data.Cookie(1048576)
	conn.inflights[cookie] = inflight{}
	_, err := conn.Next()
	c.Assert(err, Equals, io.EOF)
	_, ok := conn.inflights[cookie]
	c.Assert(ok, Equals, true)
}

func (s *ConnectionXMPPSuite) Test_Next_removesIfThereIsNoFrom(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000'></client:iq>")}
	inflights := make(map[data.Cookie]inflight)
	conn := conn{
		in:        xml.NewDecoder(mockIn),
		inflights: inflights,
	}
	cookie := data.Cookie(1048576)
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

func (s *ConnectionXMPPSuite) Test_Next_removesIfThereIsTheFromIsSameAsJid(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000' from='some@one.org/foo'></client:iq>")}
	inflights := make(map[data.Cookie]inflight)
	conn := conn{
		in:        xml.NewDecoder(mockIn),
		inflights: inflights,
		jid:       "some@one.org/foo",
	}
	cookie := data.Cookie(1048576)
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

func (s *ConnectionXMPPSuite) Test_Next_removesIfThereIsTheFromIsSameAsJidWithoutResource(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000' from='some@one.org'></client:iq>")}
	inflights := make(map[data.Cookie]inflight)
	conn := conn{
		in:        xml.NewDecoder(mockIn),
		inflights: inflights,
		jid:       "some@one.org/foo",
	}
	cookie := data.Cookie(1048576)
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

func (s *ConnectionXMPPSuite) Test_Next_removesIfThereIsTheFromIsSameAsJidDomain(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000' from='one.org'></client:iq>")}
	inflights := make(map[data.Cookie]inflight)
	conn := conn{
		in:        xml.NewDecoder(mockIn),
		inflights: inflights,
		jid:       "some@one.org/foo",
	}
	cookie := data.Cookie(1048576)
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

func (s *ConnectionXMPPSuite) Test_Next_returnsNonIQMessage(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:message xmlns:client='jabber:client' to='fo@bar.com' from='bar@foo.com' type='chat'><client:body>something</client:body></client:message>")}
	conn := conn{
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

func (s *ConnectionXMPPSuite) Test_makeInOut_returnsANewDecoderAndOriginalWriterWhenNoConfigIsGiven(c *C) {
	mockBoth := &mockConnIOReaderWriter{}
	_, rout := makeInOut(mockBoth, data.Config{})
	c.Assert(rout, Equals, mockBoth)
}

func (s *ConnectionXMPPSuite) Test_makeInOut_returnsANewDecoderAndWrappedWriterWhenConfigIsGiven(c *C) {
	mockBoth := &mockConnIOReaderWriter{}
	mockInLog := &mockConnIOReaderWriter{}
	config := data.Config{InLog: mockInLog, OutLog: mockInLog}
	_, rout := makeInOut(mockBoth, config)
	c.Assert(rout, Not(Equals), mockBoth)
}

func (s *ConnectionXMPPSuite) Test_Dial_returnsErrorFromGetFeatures(c *C) {
	rw := &mockConnIOReaderWriter{}
	conn := &fullMockedConn{rw: rw}

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXMPPSuite) Test_Dial_returnsErrorFromAuthenticateIfSkipTLS(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'><str:features></str:features>")}
	conn := &fullMockedConn{rw: rw}

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config:   data.Config{SkipTLS: true},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, errors.ErrAuthenticationFailed)
}

func (s *ConnectionXMPPSuite) Test_Dial_returnsErrorFromSecondFeatureCheck(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config:   data.Config{SkipTLS: true},
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

func (s *ConnectionXMPPSuite) Test_Dial_returnsErrorFromIQReturn(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config:   data.Config{SkipTLS: true},
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

func (s *ConnectionXMPPSuite) Test_Dial_returnsWorkingConnIfEverythingPasses(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config:   data.Config{SkipTLS: true},
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

func (s *ConnectionXMPPSuite) Test_Dial_failsIfTheServerDoesntSupportTLS(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "xmpp: server doesn't support TLS")
}

func (s *ConnectionXMPPSuite) Test_Dial_failsIfReceivingEOFAfterStartingTLS(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Matches, "(XML syntax error on line 1: unexpected )?EOF")
}

func (s *ConnectionXMPPSuite) Test_Dial_failsIfReceivingTheWrongNamespaceAfterStarttls(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "xmpp: expected <proceed> after <starttls> but got <proceed> in http://etherx.jabber.org/streams")
}

func (s *ConnectionXMPPSuite) Test_Dial_failsIfReceivingTheWrongTagName(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "xmpp: expected <proceed> after <starttls> but got <things> in urn:ietf:params:xml:ns:xmpp-tls")
}

func (s *ConnectionXMPPSuite) Test_Dial_setsServerNameOnTLSContext(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
			TLSConfig: &tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXMPPSuite) Test_Dial_failsIfDecodingFallbackFails(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
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

func (s *ConnectionXMPPSuite) Test_Dial_failsIfAccountCreationFails(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
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

func (s *ConnectionXMPPSuite) Test_Dial_failsIfTheIQQueryHasNoContent(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
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

func (s *ConnectionXMPPSuite) Test_Dial_ifRegisterQueryDoesntContainDataFailsAtNextIQ(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
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

func (s *ConnectionXMPPSuite) Test_Dial_afterRegisterFailsIfReceivesAnErrorElement(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
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

func (s *ConnectionXMPPSuite) Test_Dial_sendsBackUsernameAndPassword(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
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

func (s *ConnectionXMPPSuite) Test_Dial_runsForm(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
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

func (s *ConnectionXMPPSuite) Test_Dial_setsLog(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
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

func (s *ConnectionXMPPSuite) Test_Dial_failsWhenTryingToEstablishSession(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
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

func (s *ConnectionXMPPSuite) Test_Dial_failsWhenTryingToEstablishSessionAndGetsTheWrongIQBack(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
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

func (s *ConnectionXMPPSuite) Test_Dial_succeedsEstablishingASession(c *C) {
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

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
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

// func (s *ConnectionXMPPSuite) Test_blaData(c *C) {
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
// 	conn, _ := net.Dial("tcp", "www.olabini.se:443")
// 	tee := createTeeConn(conn, os.Stdout)
// 	tlsConn := tls.Client(tee, &tlsC)
// 	err := tlsConn.Handshake()
// 	if err != nil {
// 		println("Error: ", err.Error())
// 	}
// }

func (s *ConnectionXMPPSuite) Test_readMessages_passesStanzaToChannel(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:message xmlns:client='jabber:client' to='fo@bar.com' from='bar@foo.com' type='chat'><client:body>something</client:body></client:message>")}

	conn := &conn{
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

func (s *ConnectionXMPPSuite) Test_readMessages_alertsOnError(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<clientx:message xmlns:client='jabber:client' to='fo@bar.com' from='bar@foo.com' type='chat'><client:body>something</client:body></client:message>")}

	conn := &conn{
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
