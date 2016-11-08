package gui

import (
	"strings"

	. "gopkg.in/check.v1"
)

type JidSuite struct{}

var _ = Suite(&JidSuite{})

func (s *JidSuite) Test_verify_addressWithSimpleDomainPart(c *C) {
	address := "local@domain.com"
	valid, err := verifyXMPPAddress(address)
	c.Assert(valid, Equals, true)
	c.Assert(err, Equals, "")
}

func (s *JidSuite) Test_verify_addressWithIncompleteDomainPart(c *C) {
	address := "local@domain."
	valid, err := verifyXMPPAddress(address)
	c.Assert(valid, Equals, false)
	c.Assert(strings.Contains(err, "domain"), Equals, true)
}

func (s *JidSuite) Test_verify_addressWithCompoundDomainPart(c *C) {
	address := "local@domain.com.ec"
	valid, _ := verifyXMPPAddress(address)
	c.Assert(valid, Equals, true)
}

func (s *JidSuite) Test_verify_addressWithoutServerInDomainPart(c *C) {
	address := "local@.com"
	valid, err := verifyXMPPAddress(address)
	c.Assert(valid, Equals, false)
	c.Assert(strings.Contains(err, "domain"), Equals, true)
}

func (s *JidSuite) Test_verify_addressWithoutLocalPart(c *C) {
	address := "@domain.com"
	valid, err := verifyXMPPAddress(address)
	c.Assert(valid, Equals, false)
	c.Assert(strings.Contains(err, "part"), Equals, true)
}

func (s *JidSuite) Test_verify_addressWithLocalPartIncludingDot(c *C) {
	address := "local.foo@domain.com"
	valid, err := verifyXMPPAddress(address)
	c.Assert(valid, Equals, true)
	c.Assert(err, Equals, "")
}

func (s *JidSuite) Test_verify_addressWithSeveralErrors(c *C) {
	address := "@.com"
	valid, err := verifyXMPPAddress(address)
	c.Assert(valid, Equals, false)
	c.Assert(strings.Contains(err, "domain"), Equals, true)
	c.Assert(strings.Contains(err, "part"), Equals, true)
	c.Assert(strings.Contains(err, "local@domain.com"), Equals, true)
}
