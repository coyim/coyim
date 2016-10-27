package net

import (
	"net"

	. "gopkg.in/check.v1"
)

type TorSuite struct{}

var _ = Suite(&TorSuite{})

func (s *TorSuite) TestDetectTor(c *C) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	c.Assert(err, IsNil)

	_, port, err := net.SplitHostPort(ln.Addr().String())
	c.Assert(err, IsNil)

	tor := &defaultTorManager{
		torPorts: []string{port},
	}

	torAddress := ln.Addr().String()
	c.Assert(tor.Address(), Equals, torAddress)

	ln.Close()

	c.Assert(tor.Address(), Equals, torAddress)

	c.Assert(tor.Detect(), Equals, false)
	c.Assert(tor.Address(), Equals, "")
}

func (s *TorSuite) TestDetectTorConnectionRefused(c *C) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	c.Assert(err, IsNil)

	_, port, err := net.SplitHostPort(ln.Addr().String())
	c.Assert(err, IsNil)

	ln.Close()

	tor := &defaultTorManager{
		torPorts: []string{port},
	}

	c.Assert(tor.Detect(), Equals, false)
	c.Assert(tor.Address(), Equals, "")
}
