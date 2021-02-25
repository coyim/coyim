package config

import (
	"net/url"

	"golang.org/x/net/proxy"
	. "gopkg.in/check.v1"
)

type SOCKSSuite struct{}

var _ = Suite(&SOCKSSuite{})

func (s *SOCKSSuite) Test_socks5UnixProxy_returnsASimpleProxyWithoutAuth(c *C) {
	unixURL := &url.URL{
		Path: "/hello/goodbye",
	}

	res, e := socks5UnixProxy(unixURL, nil)
	c.Assert(e, IsNil)
	expected, _ := proxy.SOCKS5("unix", "/hello/goodbye", nil, nil)
	c.Assert(res, DeepEquals, expected)
}

func (s *SOCKSSuite) Test_socks5UnixProxy_returnsASimpleProxyWithUsername(c *C) {
	unixURL := &url.URL{
		Path: "/hello/goodbye",
		User: url.User("someone"),
	}

	_, e := socks5UnixProxy(unixURL, nil)
	c.Assert(e, IsNil)
}

func (s *SOCKSSuite) Test_socks5UnixProxy_returnsASimpleProxyWithUsernameAndPassword(c *C) {
	unixURL := &url.URL{
		Path: "/hello/goodbye",
		User: url.UserPassword("someone", "something"),
	}

	_, e := socks5UnixProxy(unixURL, nil)
	c.Assert(e, IsNil)
}

type TORSuite struct{}

var _ = Suite(&TORSuite{})

func (s *TORSuite) Test_genTorAutoAuth_returnsAuthWithUserOnly(c *C) {
	u := &url.URL{
		Path: "/hello/goodbye",
		User: url.User("someone"),
	}

	a := genTorAutoAuth(u)
	c.Assert(a.User, Equals, "someone")
	c.Assert(a.Password, Equals, "")
}

func (s *TORSuite) Test_genTorAutoAuth_returnsAuthWithUserAndPassword(c *C) {
	u := &url.URL{
		Path: "/hello/goodbye",
		User: url.UserPassword("someone", "well"),
	}

	a := genTorAutoAuth(u)
	c.Assert(a.User, Equals, "someone")
	c.Assert(a.Password, Equals, "well")
}

func (s *TORSuite) Test_genTorAutoAddr_returnsDetectedTorAddressIfNoneGiven(c *C) {
	origAddressFunc := tornetAddress
	defer func() {
		tornetAddress = origAddressFunc
	}()

	tornetAddress = func() string {
		return "hello.tor.world"
	}

	u := &url.URL{
		Host: "",
	}

	h := genTorAutoAddr(u)
	c.Assert(h, Equals, "hello.tor.world")
}

func (s *TORSuite) Test_genTorAutoAddr_returnsTheGivenHostIfOneExists(c *C) {
	u := &url.URL{
		Host: "some.where.com",
	}

	h := genTorAutoAddr(u)
	c.Assert(h, Equals, "some.where.com")
}

func (s *TORSuite) Test_torAutoProxy_returnsErrorIfTorCantBeFound(c *C) {
	origAddressFunc := tornetAddress
	defer func() {
		tornetAddress = origAddressFunc
	}()

	tornetAddress = func() string {
		return ""
	}

	u := &url.URL{
		Host: "",
	}

	d, e := torAutoProxy(u, nil)
	c.Assert(d, IsNil)
	c.Assert(e, Equals, ErrTorNotRunning)
}
