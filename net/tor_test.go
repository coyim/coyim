package net

import (
	"fmt"
	"net"

	. "gopkg.in/check.v1"
)

type TorSuite struct{}

var _ = Suite(&TorSuite{})

func (s *TorSuite) TestDetectTor(c *C) {
	host := "127.0.0.1"
	if isTails() {
		host = getLocalIP()
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:0", host))
	c.Assert(err, IsNil)

	_, port, err := net.SplitHostPort(ln.Addr().String())
	c.Assert(err, IsNil)

	tor := &defaultTorManager{
		torPorts: []string{port},
		torHost:  host,
	}

	torAddress := ln.Addr().String()
	c.Assert(tor.Address(), Equals, torAddress)

	_ = ln.Close()

	c.Assert(tor.Address(), Equals, torAddress)

	c.Assert(tor.Detect(), Equals, false)
	c.Assert(tor.Address(), Equals, "")
}

func (s *TorSuite) TestDetectTorConnectionRefused(c *C) {
	host := "127.0.0.1"
	if isTails() {
		host = getLocalIP()
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:0", host))
	c.Assert(err, IsNil)

	_, port, err := net.SplitHostPort(ln.Addr().String())
	c.Assert(err, IsNil)

	_ = ln.Close()

	tor := &defaultTorManager{
		torPorts: []string{port},
		torHost:  host,
	}

	c.Assert(tor.Detect(), Equals, false)
	c.Assert(tor.Address(), Equals, "")
}
