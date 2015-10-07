package xmpp

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"

	. "gopkg.in/check.v1"
)

type TlsXmppSuite struct{}

var _ = Suite(&TlsXmppSuite{})

func (s *TlsXmppSuite) Test_certName_returnsEmptyInformation(c *C) {
	cert := &x509.Certificate{}
	cert.Subject = pkix.Name{}
	res := certName(cert)
	c.Assert(res, Equals, "")
}

func (s *TlsXmppSuite) Test_certName_usesNameInformation(c *C) {
	cert := &x509.Certificate{}
	cert.Subject = pkix.Name{}
	cert.Subject.Organization = []string{"Foo", "Bar.com"}
	cert.Subject.OrganizationalUnit = []string{"Somewhere", "Else", "Above", "Beyond"}
	cert.Subject.CommonName = "test.coyim"
	res := certName(cert)
	c.Assert(res, Equals, "O=Foo/O=Bar.com/OU=Somewhere/OU=Else/OU=Above/OU=Beyond/CN=test.coyim/")
}

func (s *TlsXmppSuite) Test_printTLSDetails_printsUnknownVersions(c *C) {
	mockWriter := mockConnIOReaderWriter{}
	state := tls.ConnectionState{
		Version: 0x0200,
	}

	printTLSDetails(&mockWriter, state)

	c.Assert(string(mockWriter.write), Equals, ""+
		"  SSL/TLS version: unknown\n"+
		"  Cipher suite: unknown\n",
	)
}

func (s *TlsXmppSuite) Test_printTLSDetails_printsCorrectVersions(c *C) {
	mockWriter := mockConnIOReaderWriter{}
	state := tls.ConnectionState{
		Version:     tls.VersionTLS11,
		CipherSuite: 0xc00a,
	}

	printTLSDetails(&mockWriter, state)

	c.Assert(string(mockWriter.write), Equals, ""+
		"  SSL/TLS version: TLS 1.1\n"+
		"  Cipher suite: TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA\n",
	)
}
