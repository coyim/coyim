package tls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"

	. "gopkg.in/check.v1"
)

type VerifierSuite struct{}

var _ = Suite(&VerifierSuite{})

func (s *VerifierSuite) Test_verifyHostName(c *C) {
	e := verifyHostName(&x509.Certificate{}, "foo")
	c.Assert(e, ErrorMatches, "x509: certificate is not valid for any names.*")
}

func (s *VerifierSuite) Test_BasicVerifier_verifyCert(c *C) {
	certs := []*x509.Certificate{
		&x509.Certificate{},
		&x509.Certificate{},
		&x509.Certificate{},
	}

	v := &BasicVerifier{}
	res, e := v.verifyCert(certs, &tls.Config{})
	c.Assert(e, ErrorMatches, "x509: missing ASN.1 contents; use ParseCertificate")
	c.Assert(res, IsNil)
}

func (s *VerifierSuite) Test_BasicVerifier_verifyCert_verifiesCorrectly(c *C) {
	roots := x509.NewCertPool()
	_ = roots.AppendCertsFromPEM([]byte(testCAcert))

	block, _ := pem.Decode([]byte(testCert))
	cert, _ := x509.ParseCertificate(block.Bytes)

	conf := &tls.Config{
		RootCAs: roots,
	}

	certs := []*x509.Certificate{cert}

	v := &BasicVerifier{}
	res, e := v.verifyCert(certs, conf)
	c.Assert(e, IsNil)
	c.Assert(res, HasLen, 1)
	c.Assert(res[0], HasLen, 2)
	c.Assert(res[0][0], Equals, cert)
}

func (s *VerifierSuite) Test_BasicVerifier_canConnectInPresenceOfPins_withPinningPolicy_none(c *C) {
	v := &BasicVerifier{
		PinningPolicy: "none",
	}

	c.Assert(v.canConnectInPresenceOfPins(nil), IsNil)
}

func (s *VerifierSuite) Test_BasicVerifier_canConnectInPresenceOfPins_withPinningPolicy_deny(c *C) {
	called := false
	v := &BasicVerifier{
		PinningPolicy: "deny",
		OnPinDeny: func() {
			called = true
		},
	}

	c.Assert(v.canConnectInPresenceOfPins(nil), ErrorMatches, "tls: you have a pinning policy that stops us from connecting using this certificate")
	c.Assert(called, Equals, true)
}

func (s *VerifierSuite) Test_BasicVerifier_canConnectInPresenceOfPins_withPinningPolicy_add(c *C) {
	certs := []*x509.Certificate{
		&x509.Certificate{},
		&x509.Certificate{},
		&x509.Certificate{},
	}
	called := false
	v := &BasicVerifier{
		PinningPolicy: "add",
		AddCert: func(*x509.Certificate) {
			called = true
		},
	}

	c.Assert(v.canConnectInPresenceOfPins(certs), IsNil)
	c.Assert(called, Equals, true)
}

func (s *VerifierSuite) Test_BasicVerifier_canConnectInPresenceOfPins_withPinningPolicy_default_withoutCerts(c *C) {
	certs := []*x509.Certificate{
		&x509.Certificate{},
		&x509.Certificate{},
		&x509.Certificate{},
	}
	called := false
	v := &BasicVerifier{
		PinningPolicy: "",
		HasCertificates: func() bool {
			return false
		},
		AddCert: func(*x509.Certificate) {
			called = true
		},
	}

	c.Assert(v.canConnectInPresenceOfPins(certs), IsNil)
	c.Assert(called, Equals, true)
}

func (s *VerifierSuite) Test_BasicVerifier_canConnectInPresenceOfPins_withPinningPolicy_default_withCerts(c *C) {
	certs := []*x509.Certificate{
		&x509.Certificate{},
		&x509.Certificate{},
		&x509.Certificate{},
	}
	v := &BasicVerifier{
		PinningPolicy: "add-first-ask-rest",
		HasCertificates: func() bool {
			return true
		},
		AskPinning: func([]*x509.Certificate) error {
			return errors.New("marker return")
		},
	}

	c.Assert(v.canConnectInPresenceOfPins(certs), ErrorMatches, "marker return")
}

func (s *VerifierSuite) Test_BasicVerifier_canConnectInPresenceOfPins_withPinningPolicy_addFirstDenyRest_withoutCertificates(c *C) {
	certs := []*x509.Certificate{
		&x509.Certificate{},
		&x509.Certificate{},
		&x509.Certificate{},
	}
	called := false
	v := &BasicVerifier{
		PinningPolicy: "add-first-deny-rest",
		HasCertificates: func() bool {
			return false
		},
		AddCert: func(*x509.Certificate) {
			called = true
		},
	}

	c.Assert(v.canConnectInPresenceOfPins(certs), IsNil)
	c.Assert(called, Equals, true)
}

func (s *VerifierSuite) Test_BasicVerifier_canConnectInPresenceOfPins_withPinningPolicy_addFirstDenyRest_withCertificates(c *C) {
	certs := []*x509.Certificate{
		&x509.Certificate{},
		&x509.Certificate{},
		&x509.Certificate{},
	}
	called := false
	v := &BasicVerifier{
		PinningPolicy: "add-first-deny-rest",
		HasCertificates: func() bool {
			return true
		},
		OnPinDeny: func() {
			called = true
		},
	}

	c.Assert(v.canConnectInPresenceOfPins(certs), ErrorMatches, "tls: you have a pinning policy that stops us from connecting using this certificate")
	c.Assert(called, Equals, true)
}

func (s *VerifierSuite) Test_BasicVerifier_canConnectInPresenceOfPins_withPinningPolicy_ask(c *C) {
	certs := []*x509.Certificate{
		&x509.Certificate{},
		&x509.Certificate{},
		&x509.Certificate{},
	}
	v := &BasicVerifier{
		PinningPolicy: "ask",
		AskPinning: func([]*x509.Certificate) error {
			return errors.New("another marker")
		},
	}

	c.Assert(v.canConnectInPresenceOfPins(certs), ErrorMatches, "another marker")
}

func (s *VerifierSuite) Test_BasicVerifier_canConnectInPresenceOfPins_withPinningPolicy_other(c *C) {
	certs := []*x509.Certificate{
		&x509.Certificate{},
		&x509.Certificate{},
		&x509.Certificate{},
	}
	called := false
	v := &BasicVerifier{
		PinningPolicy: "bla",
		OnPinDeny: func() {
			called = true
		},
	}

	c.Assert(v.canConnectInPresenceOfPins(certs), ErrorMatches, "tls: you have a pinning policy that stops us from connecting using other certificates")
	c.Assert(called, Equals, true)
}

func (s *VerifierSuite) Test_BasicVerifier_Verify_withNoCerts(c *C) {
	called := false
	v := &BasicVerifier{
		OnNoPeerCertificates: func() {
			called = true
		},
	}
	e := v.Verify(tls.ConnectionState{PeerCertificates: []*x509.Certificate{}}, nil, "")
	c.Assert(e, ErrorMatches, "tls: server has no certificates")
	c.Assert(called, Equals, true)
}

func (s *VerifierSuite) Test_BasicVerifier_Verify_withAlreadyPinnedCerts(c *C) {
	called := false
	v := &BasicVerifier{
		HasPinned: func([]*x509.Certificate) bool {
			called = true
			return true
		},
	}
	e := v.Verify(tls.ConnectionState{PeerCertificates: []*x509.Certificate{&x509.Certificate{}}}, nil, "")
	c.Assert(e, IsNil)
	c.Assert(called, Equals, true)
}

func (s *VerifierSuite) Test_BasicVerifier_Verify_verifiesCorrectly(c *C) {
	roots := x509.NewCertPool()
	_ = roots.AppendCertsFromPEM([]byte(testCAcert))

	block, _ := pem.Decode([]byte(testCert))
	cert, _ := x509.ParseCertificate(block.Bytes)

	conf := &tls.Config{
		RootCAs: roots,
	}

	v := &BasicVerifier{
		HasPinned: func([]*x509.Certificate) bool {
			return false
		},
	}

	e := v.Verify(tls.ConnectionState{PeerCertificates: []*x509.Certificate{cert}}, conf, "foo.bar.com")
	c.Assert(e, IsNil)
}

func (s *VerifierSuite) Test_BasicVerifier_Verify_callsVerifyFailureWhenCantVerify(c *C) {
	block, _ := pem.Decode([]byte(testCert))
	cert, _ := x509.ParseCertificate(block.Bytes)

	conf := &tls.Config{}

	v := &BasicVerifier{
		HasPinned: func([]*x509.Certificate) bool {
			return false
		},
		VerifyFailure: func([]*x509.Certificate, error) error {
			return errors.New("marker error")
		},
	}

	e := v.Verify(tls.ConnectionState{PeerCertificates: []*x509.Certificate{cert}}, conf, "foo.bar.com")
	c.Assert(e, ErrorMatches, "marker error")
}

func (s *VerifierSuite) Test_BasicVerifier_Verify_failsOnBadHostname(c *C) {
	roots := x509.NewCertPool()
	_ = roots.AppendCertsFromPEM([]byte(testCAcert))

	block, _ := pem.Decode([]byte(testCert))
	cert, _ := x509.ParseCertificate(block.Bytes)

	conf := &tls.Config{
		RootCAs: roots,
	}

	v := &BasicVerifier{
		HasPinned: func([]*x509.Certificate) bool {
			return false
		},
		VerifyHostnameFailure: func([]*x509.Certificate, string, error) error {
			return errors.New("other marker error")
		},
	}

	e := v.Verify(tls.ConnectionState{PeerCertificates: []*x509.Certificate{cert}}, conf, "something.else.com")
	c.Assert(e, ErrorMatches, "other marker error")
}

func (s *VerifierSuite) Test_BasicVerifier_Verify_checksPinsForVerifiedDomain(c *C) {
	roots := x509.NewCertPool()
	_ = roots.AppendCertsFromPEM([]byte(testCAcert))

	block, _ := pem.Decode([]byte(testCert))
	cert, _ := x509.ParseCertificate(block.Bytes)

	conf := &tls.Config{
		RootCAs: roots,
	}

	called := false
	v := &BasicVerifier{
		HasPinned: func([]*x509.Certificate) bool {
			return false
		},
		NeedToCheckPins: true,
		PinningPolicy:   "deny",
		OnPinDeny: func() {
			called = true
		},
	}

	e := v.Verify(tls.ConnectionState{PeerCertificates: []*x509.Certificate{cert}}, conf, "foo.bar.com")
	c.Assert(e, ErrorMatches, "tls: you have a pinning policy that stops us from connecting using this certificate")
	c.Assert(called, Equals, true)
}
