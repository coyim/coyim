package config

import (
	"crypto/x509"
	"errors"

	. "gopkg.in/check.v1"
)

type CertsSuite struct{}

var _ = Suite(&CertsSuite{})

func (s *CertsSuite) Test_rootCAFor_jabbercccde_withFailingParse(c *C) {
	origX509ParseCertificate := x509ParseCertificate
	defer func() {
		x509ParseCertificate = origX509ParseCertificate
	}()
	x509ParseCertificate = func([]byte) (*x509.Certificate, error) {
		return nil, errors.New("oh nooooooooooooooo")
	}

	res, e := rootCAFor("jabber.ccc.de")
	c.Assert(e, ErrorMatches, "oh no+")
	c.Assert(res, IsNil)
}
