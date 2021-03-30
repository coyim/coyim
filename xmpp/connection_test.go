package xmpp

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"io"
	"time"

	goerr "errors"

	"github.com/coyim/coyim/cache"
	"github.com/coyim/coyim/digests"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/errors"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	. "gopkg.in/check.v1"
)

type ConnectionXMPPSuite struct{}

var _ = Suite(&ConnectionXMPPSuite{})

func tlsConfigForRecordedHandshake() (*tls.Config, *basicTLSVerifier) {
	//This is the certificate in the recorded handshake
	peerCertificatePEM := []byte(`
-----BEGIN CERTIFICATE-----
MIIF5TCCA82gAwIBAgIQJkO7MqFmSHrhnWx5xD/iZjANBgkqhkiG9w0BAQsFADB9
MQswCQYDVQQGEwJJTDEWMBQGA1UEChMNU3RhcnRDb20gTHRkLjErMCkGA1UECxMi
U2VjdXJlIERpZ2l0YWwgQ2VydGlmaWNhdGUgU2lnbmluZzEpMCcGA1UEAxMgU3Rh
cnRDb20gQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkwHhcNMTUxMjE2MDEwMDA1WhcN
MzAxMjE2MDEwMDA1WjB4MQswCQYDVQQGEwJJTDEWMBQGA1UEChMNU3RhcnRDb20g
THRkLjEpMCcGA1UECxMgU3RhcnRDb20gQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkx
JjAkBgNVBAMTHVN0YXJ0Q29tIENsYXNzIDIgSVYgU2VydmVyIENBMIIBIjANBgkq
hkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnL29gjx6E467y4OsHo42TCn1rC7JXUnv
epzPE9KLbJiQi63JSLTr/QVGjhWFQBhqwXKlyTyBNGoOuV+yRoimqkPDdV6ZdnIn
RwmKAnVhvMVd2WXeqSJtq5STa2nuOnLTwYBnyVsOIo9YdnvFhDXAGjQ3hXWQIq00
f43XE8Fik+9EUG/oF7VLlIACAJnhotAj2dR2TvQmyBbEEN2PhLH3WANZklMbao2c
sASqSwyOmAB5+35nSagpMYuuVa4ZSnm2EaF8emLxiiFK5InCBZjRG4u+YLrEv7+m
KrnHOMVWkOE7mzKxtuHFYW2LRB++eJGLUdn1KiviZDS/ofOhIhfstwIDAQABo4IB
ZDCCAWAwDgYDVR0PAQH/BAQDAgEGMB0GA1UdJQQWMBQGCCsGAQUFBwMCBggrBgEF
BQcDATASBgNVHRMBAf8ECDAGAQH/AgEAMDIGA1UdHwQrMCkwJ6AloCOGIWh0dHA6
Ly9jcmwuc3RhcnRzc2wuY29tL3Nmc2NhLmNybDBmBggrBgEFBQcBAQRaMFgwJAYI
KwYBBQUHMAGGGGh0dHA6Ly9vY3NwLnN0YXJ0c3NsLmNvbTAwBggrBgEFBQcwAoYk
aHR0cDovL2FpYS5zdGFydHNzbC5jb20vY2VydHMvY2EuY3J0MB0GA1UdDgQWBBSU
3oVBKqXZRfZgLC5MkwmmLCN+PjAfBgNVHSMEGDAWgBROC+8apEBbpRdphzDKNGhD
0EGu8jA/BgNVHSAEODA2MDQGBFUdIAAwLDAqBggrBgEFBQcCARYeaHR0cDovL3d3
dy5zdGFydHNzbC5jb20vcG9saWN5MA0GCSqGSIb3DQEBCwUAA4ICAQC16kMuZh8h
lVsgzybaIix2qySQFU+rPgqSqeyrDSmJwpDbaKjwakm6LJ2DLX5MRFjNPCh+ArQf
CU1UUJa65n7UaQWt6q8kUwifHcIn+fFJdNV3N4zdvlKxwveqBSQZiXeIUO/hHr1U
i7Gw6s0On+K0fD9oNcgCRR3vPicB2frK7BhOFje6xowsWexxPfJHI69lCq73O7Ke
xXqp/V8f8uGF8L4KU3xW6RDG57RrXh5+LNxUQmZ2tIAaPyHTND5zbxff8Z/ZbgGG
HKbsuPkAUIG+bHpq5b6bf2x2NxMhqYSMI+GJJ9FmmiCV+P3+0ywBYGNhJkcFUYvo
SUduHz+/RXd6G/ejrvKp58rbZ9iCISLZjpo5gYEfLIl6IQJcZPM8FIWKLKhtIoKX
5ctNL3epV4DzIDZxLaSruEBQFeDQj6p/74pUYLQBP523anf6StXBtYgbfImRoIh4
I8L85aB/TUyLOJA/sKx/WFrXOxE9K4q+Pf5tq3gzZEchM/btMYn1cw1GPUt4nHya
zS52LrP0+Q77ao1Gza9svd8HE1NZ9NIVJO71QskqjxvGiTt048r4gLSXaM1zP2w9
nMsIw1IpxXE8h9UHAllgh8oNHno5I9nLfynbEhXxGy9RlfcLN/J8iOqyagfgxrUy
DPKMh5xGeLKMQSzjyQ1bV0WGC1JmJp+QDQ==
-----END CERTIFICATE-----
	`)

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(peerCertificatePEM) {
		panic("Bad PEM")
	}

	c := &tls.Config{
		RootCAs: pool,
		CipherSuites: []uint16{
			0xc02f, 0xc02b, 0xc030, 0xc02c, 0xc013, 0xc009, 0xc014,
			0xc00a, 0x009c, 0x009d, 0x002f, 0x0035, 0xc012, 0x000a,
		},
		//0x1d appears as preffered since go1.8
		//Manually inform what was used when the handshake for 1.6 was recorded
		CurvePreferences: []tls.CurveID{0x17, 0x18, 0x19},
	}

	v := &basicTLSVerifier{
		shaSum: []byte{
			0x82, 0x45, 0x44, 0x18, 0xcb, 0x04, 0x85, 0x4a,
			0xa7, 0x21, 0xbb, 0x05, 0x96, 0x52, 0x8f, 0xf8,
			0x02, 0xb1, 0xe1, 0x8a, 0x4e, 0x3a, 0x77, 0x67,
			0x41, 0x2a, 0xc9, 0xf1, 0x08, 0xc9, 0xd3, 0xa7,
		},
	}

	return c, v
}

type basicTLSVerifier struct {
	shaSum []byte
}

func (v *basicTLSVerifier) verifyCert(state tls.ConnectionState, conf *tls.Config) ([][]*x509.Certificate, error) {
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

func (v *basicTLSVerifier) Verify(state tls.ConnectionState, conf *tls.Config, originDomain string) error {
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
		log: testLogger(),
		in:  xml.NewDecoder(mockIn),
	}

	_, err := conn.Next()
	c.Assert(err.Error(), Equals, "unexpected XMPP message http://etherx.jabber.org/streams <foo/>")
}

func (s *ConnectionXMPPSuite) Test_Next_returnsErrorIfFailingToParseIQID(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='abczzzz'></client:iq>")}
	conn := conn{
		log: testLogger(),
		in:  xml.NewDecoder(mockIn),
	}

	_, err := conn.Next()
	c.Assert(err.Error(), Equals, "xmpp: failed to parse id from iq: strconv.ParseUint: parsing \"abczzzz\": invalid syntax")
}

func (s *ConnectionXMPPSuite) Test_Next_returnsNothingIfThereIsNoInflightMatching(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000'></client:iq>")}
	conn := conn{
		log: testLogger(),
		in:  xml.NewDecoder(mockIn),
	}

	_, err := conn.Next()
	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXMPPSuite) Test_Next_returnsNothingIfTheInflightIsToAnotherReceiver(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000' from='bar@somewhere.com'></client:iq>")}
	conn := conn{
		log:       testLogger(),
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
		log:       testLogger(),
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
		log:       testLogger(),
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
		log:       testLogger(),
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
		log:       testLogger(),
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
		log:       testLogger(),
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
		log:       testLogger(),
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
		log: testLogger(),
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
		log:      testLogger(),
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
		log:      testLogger(),
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
		log:      testLogger(),
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Matches, "unmarshal <iq>:( XML syntax error on line 1: unexpected)? EOF")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'></bind></iq>",
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
		log:      testLogger(),
	}
	_, err := d.setupStream(conn)

	c.Assert(err, IsNil)
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'></bind></iq>",
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
		log: testLogger(),
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
		log: testLogger(),
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
		log: testLogger(),
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
		log: testLogger(),
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
		log: testLogger(),
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
		log: testLogger(),
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
		log: testLogger(),
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
		log: testLogger(),
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Matches, "unmarshal <iq>:( XML syntax error on line 1: unexpected)? EOF")
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
		log: testLogger(),
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Matches, "xmpp: unmarshal <iq>:( XML syntax error on line 1: unexpected)? EOF")

	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'></bind></iq>"+
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
		log: testLogger(),
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "xmpp: session establishment failed")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'></bind></iq>"+
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
		log: testLogger(),
	}
	_, err := d.setupStream(conn)

	c.Assert(err, IsNil)
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'></bind></iq>"+
		"<iq to='domain' type='set' id='sess_1'><session xmlns='urn:ietf:params:xml:ns:xmpp-session'/></iq>",
	)
}

func (s *ConnectionXMPPSuite) Test_readMessages_passesStanzaToChannel(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:message xmlns:client='jabber:client' to='fo@bar.com' from='bar@foo.com' type='chat'><client:body>something</client:body></client:message>")}

	conn := &conn{
		log:    testLogger(),
		in:     xml.NewDecoder(mockIn),
		closed: true, //This avoids trying to close the connection after the EOF
	}
	stanzaChan := make(chan data.Stanza)
	go func() {
		_ = conn.ReadStanzas(stanzaChan)
	}()

	rawStanza, ok := <-stanzaChan
	c.Assert(ok, Equals, true)
	c.Assert(rawStanza.Name.Local, Equals, "message")
	c.Assert(rawStanza.Value.(*data.ClientMessage).Body, Equals, "something")
}

func (s *ConnectionXMPPSuite) Test_readMessages_alertsOnError(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<clientx:message xmlns:client='jabber:client' to='fo@bar.com' from='bar@foo.com' type='chat'><client:body>something</client:body></client:message>")}

	conn := &conn{
		log:    testLogger(),
		in:     xml.NewDecoder(mockIn),
		closed: true, //This avoids trying to close the connection after the EOF
	}

	stanzaChan := make(chan data.Stanza, 1)
	err := conn.ReadStanzas(stanzaChan)

	_, ok := <-stanzaChan
	c.Assert(ok, Equals, false)
	c.Assert(err.Error(), Equals, "unexpected XMPP message clientx <message/>")
}

func (s *ConnectionXMPPSuite) Test_Dial_failsWhenStartingAHandshake(c *C) {
	tlsC, v := tlsConfigForRecordedHandshake()

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
	t := &tlsMock1{returnFromHandshake: io.EOF}
	d := &dialer{
		log:      testLogger(),
		JID:      "user@domain",
		password: "pass",

		verifier: v,
		config: data.Config{
			TLSConfig: tlsC,
		},
		tlsConnFactory: fixedTLSFactory(t),
	}

	conn := &fullMockedConn{rw: rw}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXMPPSuite) Test_Dial_worksIfTheHandshakeSucceeds(c *C) {
	tlsC, _ := tlsConfigForRecordedHandshake()

	rw := &mockMultiConnIOReaderWriter{read: validTLSExchange}
	conn := &fullMockedConn{rw: rw}
	connState := tls.ConnectionState{
		Version:           tls.VersionTLS12,
		HandshakeComplete: true,
		CipherSuite:       tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	}
	t := &tlsMock1{
		returnFromHandshake: nil,
		returnFromConnState: connState,
		returnFromRead1:     0,
		returnFromRead2:     io.EOF,
	}
	v := &mockTLSVerifier{
		toReturn: nil,
	}
	d := &dialer{
		log:           testLogger(),
		JID:           "user@www.olabini.se",
		password:      "pass",
		serverAddress: "www.olabini.se:443",
		verifier:      v,
		config: data.Config{
			TLSConfig: tlsC,
		},
		tlsConnFactory: fixedTLSFactory(t),
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
	c.Assert(v.verifyCalled, Equals, 1)
	c.Assert(v.originDomain, Equals, "www.olabini.se")
}

func (s *ConnectionXMPPSuite) Test_Dial_worksIfTheHandshakeSucceedsButFailsOnInvalidCertHash(c *C) {
	tlsC, _ := tlsConfigForRecordedHandshake()

	rw := &mockMultiConnIOReaderWriter{read: validTLSExchange}
	conn := &fullMockedConn{rw: rw}
	t := &tlsMock1{
		returnFromHandshake: nil,
		returnFromConnState: tls.ConnectionState{},
	}
	v := &mockTLSVerifier{
		toReturn: goerr.New("tls: server certificate does not match expected hash (got: 82454418cb04854aa721bb0596528ff802b1e18a4e3a7767412ac9f108c9d3a7, want: 6161616161)"),
	}
	d := &dialer{
		log:           testLogger(),
		JID:           "user@www.olabini.se",
		password:      "pass",
		serverAddress: "www.olabini.se:443",
		verifier:      v,

		config: data.Config{
			TLSConfig: tlsC,
		},
		tlsConnFactory: fixedTLSFactory(t),
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "tls: server certificate does not match expected hash (got: 82454418cb04854aa721bb0596528ff802b1e18a4e3a7767412ac9f108c9d3a7, want: 6161616161)")
}

func (s *ConnectionXMPPSuite) Test_Dial_worksIfTheHandshakeSucceedsButSucceedsOnValidCertHash(c *C) {
	tlsC, _ := tlsConfigForRecordedHandshake()
	tlsC.Rand = fixedRand([]string{
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F",
		"000102030405060708090A0B0C0D0E0F",
	})

	rw := &mockMultiConnIOReaderWriter{read: validTLSExchange}
	conn := &fullMockedConn{rw: rw}
	connState := tls.ConnectionState{
		Version:           tls.VersionTLS12,
		HandshakeComplete: true,
		CipherSuite:       tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	}
	t := &tlsMock1{
		returnFromHandshake: nil,
		returnFromConnState: connState,
		returnFromRead1:     0,
		returnFromRead2:     io.EOF,
	}
	v := &mockTLSVerifier{
		toReturn: nil,
	}
	d := &dialer{
		log:           testLogger(),
		JID:           "user@www.olabini.se",
		password:      "pass",
		serverAddress: "www.olabini.se:443",
		verifier:      v,

		config: data.Config{
			TLSConfig: tlsC,
		},
		tlsConnFactory: fixedTLSFactory(t),
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
	c.Assert(v.verifyCalled, Equals, 1)
	c.Assert(v.originDomain, Equals, "www.olabini.se")
}

func (s *ConnectionXMPPSuite) Test_conn_Cache(c *C) {
	one := cache.NewWithExpiry()
	two := cache.NewWithExpiry()

	cc := NewConn(nil, nil, "").(*conn)
	cc.c = one

	c.Assert(cc.Cache(), Equals, one)
	c.Assert(cc.Cache(), Not(Equals), two)
}

func (s *ConnectionXMPPSuite) Test_conn_RawOut(c *C) {
	one := &mockConnIOReaderWriter{}
	cc := &conn{rawOut: one}
	c.Assert(cc.RawOut(), Equals, one)
}

func (s *ConnectionXMPPSuite) Test_conn_ServerAddress(c *C) {
	cc := &conn{serverAddress: "something"}
	c.Assert(cc.ServerAddress(), Equals, "something")
}

func (s *ConnectionXMPPSuite) Test_conn_Resource(c *C) {
	cc := &conn{}
	cc.SetJIDResource("bla blu")

	c.Assert(cc.GetJIDResource(), Equals, "bla blu")
}

func (s *ConnectionXMPPSuite) Test_conn_Close(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)
	out := &mockConnIOReaderWriter{
		err: goerr.New("haha"),
	}
	cc := &conn{
		closed: false,
		log:    l,
		out:    out,
		rawOut: out,
	}

	e := cc.Close()
	c.Assert(e, IsNil)
	c.Assert(len(hook.Entries), Equals, 2)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "xmpp: sending closing stream tag")
	c.Assert(hook.Entries[1].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[1].Message, Equals, "xmpp: TCP closed")
}

func (s *ConnectionXMPPSuite) Test_conn_waitForStreamClosed_withTimeout(c *C) {
	orgStreamClosedTimeout := streamClosedTimeout
	defer func() {
		streamClosedTimeout = orgStreamClosedTimeout
	}()

	streamClosedTimeout = 1 * time.Millisecond

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	cc := &conn{
		log: l,
	}
	cc.streamCloseReceived = make(chan bool)

	waitForStreamClosed(cc)

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "xmpp: timed out waiting for closing stream")
}

func (s *ConnectionXMPPSuite) Test_conn_waitForStreamClosed_withoutTimeout(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	cc := &conn{
		log: l,
	}
	ch := make(chan bool, 1)
	cc.streamCloseReceived = ch
	ch <- true

	waitForStreamClosed(cc)

	c.Assert(len(hook.Entries), Equals, 0)
}

func (s *ConnectionXMPPSuite) Test_conn_SendMessage_withoutID(c *C) {
	out := &mockConnIOReaderWriter{}
	cc := &conn{out: out}

	m := &data.Message{}

	e := cc.SendMessage(m)

	c.Assert(e, IsNil)
	c.Assert(string(out.Written()), Matches, `<message xmlns="jabber:client" from="" id=".+" to="" type=""><body></body></message>`)
}

func (s *ConnectionXMPPSuite) Test_conn_SendMessage_withID(c *C) {
	out := &mockConnIOReaderWriter{}
	cc := &conn{out: out}

	m := &data.Message{ID: "hello"}

	e := cc.SendMessage(m)

	c.Assert(e, IsNil)
	c.Assert(string(out.Written()), Matches, `<message xmlns="jabber:client" from="" id="hello" to="" type=""><body></body></message>`)
}

func (s *ConnectionXMPPSuite) Test_conn_ReadStanzas_returnsOnStreamClose(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte(`<stream:stream xmlns:stream="http://etherx.jabber.org/streams" version="1.0"></stream:stream>`)}
	dec := xml.NewDecoder(mockIn)
	_, _ = dec.Token()
	cc := conn{
		log:                 testLogger(),
		out:                 mockIn,
		rawOut:              mockIn,
		in:                  dec,
		streamCloseReceived: make(chan bool),
	}

	ch := make(chan data.Stanza)
	e := cc.ReadStanzas(ch)
	c.Assert(e, IsNil)
}
