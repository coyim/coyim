package xmpp

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"

	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
)

type TLSXmppSuite struct{}

var _ = Suite(&TLSXmppSuite{})

func (s *TLSXmppSuite) Test_certName_returnsEmptyInformation(c *C) {
	cert := &x509.Certificate{}
	cert.Subject = pkix.Name{}
	res := certName(cert)
	c.Assert(res, Equals, "")
}

func (s *TLSXmppSuite) Test_certName_usesNameInformation(c *C) {
	cert := &x509.Certificate{}
	cert.Subject = pkix.Name{}
	cert.Subject.Organization = []string{"Foo", "Bar.com"}
	cert.Subject.OrganizationalUnit = []string{"Somewhere", "Else", "Above", "Beyond"}
	cert.Subject.CommonName = "test.coyim"
	res := certName(cert)
	c.Assert(res, Equals, "O=Foo/O=Bar.com/OU=Somewhere/OU=Else/OU=Above/OU=Beyond/CN=test.coyim/")
}

func (s *TLSXmppSuite) Test_printTLSDetails_printsUnknownVersions(c *C) {
	state := tls.ConnectionState{
		Version: 0x0200,
	}
	ll := log.New()
	buf := new(bytes.Buffer)
	ll.SetOutput(buf)

	printTLSDetails(ll, state)

	c.Assert(buf.String(), Matches, "(?s).*?version=unknown.*?")
	c.Assert(buf.String(), Matches, "(?s).*?cipherSuite=unknown.*?")
}

func (s *TLSXmppSuite) Test_printTLSDetails_printsCorrectVersions(c *C) {
	state := tls.ConnectionState{
		Version:     tls.VersionTLS12,
		CipherSuite: 0x1303,
	}
	ll := log.New()
	buf := new(bytes.Buffer)
	ll.SetOutput(buf)

	printTLSDetails(ll, state)

	c.Assert(buf.String(), Matches, "(?s).*?version=\"TLS 1\\.2\".*?")
	c.Assert(buf.String(), Matches, "(?s).*?cipherSuite=TLS_CHACHA20_POLY1305_SHA256.*?")
}
