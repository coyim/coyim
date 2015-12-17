package gui

import (
	. "gopkg.in/check.v1"
)

type JidSuite struct{}

var _ = Suite(&JidSuite{})

func (s *JidSuite) Test_verify_addressWithSimpleDomainPart(c *C) {
	address := "local@domain.com"
	valid, err := verify(address)
	c.Assert(valid, Equals, true)
	c.Assert(err, Equals, "")
}

func (s *JidSuite) Test_verify_addressWithIncompleteDomainPart(c *C) {
	address := "local@domain."
	valid, err := verify(address)
	c.Assert(valid, Equals, false)
	c.Assert(err, NotNil)
}

func (s *JidSuite) Test_verify_addressWithCompoundDomainPart(c *C) {
	address := "local@domain.com.ec"
	valid, _ := verify(address)
	c.Assert(valid, Equals, true)
}

func (s *JidSuite) Test_verify_addressWithoutServerInDomainPart(c *C) {
	address := "local@.com"
	valid, err := verify(address)
	c.Assert(valid, Equals, false)
	c.Assert(err, NotNil)
}

func (s *JidSuite) Test_verify_addressWithoutLocalPart(c *C) {
	address := "@domain.com"
	valid, err := verify(address)
	c.Assert(valid, Equals, false)
	c.Assert(err, NotNil)
}
