package net

import . "gopkg.in/check.v1"

type ProxiesSuite struct{}

var _ = Suite(&ProxiesSuite{})

func (s *ProxiesSuite) Test_FindProxyTypeFor_findsTheIndexForAProxyType(c *C) {
	c.Check(FindProxyTypeFor("tor-auto"), Equals, 0)
	c.Check(FindProxyTypeFor("socks5"), Equals, 1)
	c.Check(FindProxyTypeFor("something-weird"), Equals, -1)
}

func (s *ProxiesSuite) Test_GetProxyTypeNames_yieldsAllProxyTypeNames(c *C) {
	result := []string{}
	GetProxyTypeNames(func(s string) {
		result = append(result, s)
	})

	c.Assert(result, DeepEquals, []string{"Automatic Tor", "SOCKS5", "SOCKS5 over Unix Domain Socket"})
}

func (s *ProxiesSuite) Test_GetProxyTypeFor_returnsTheCorrectProxyType(c *C) {
	c.Check(GetProxyTypeFor("Automatic Tor"), Equals, "tor-auto")
	c.Check(GetProxyTypeFor("SOCKS5"), Equals, "socks5")
	c.Check(GetProxyTypeFor("something else"), Equals, "")
}

func (s *ProxiesSuite) Test_ParseProxy_returnsTheSchemeOfTheProxySpecification(c *C) {
	c.Check(ParseProxy("socks4://localhost").Scheme, Equals, "socks4")
	c.Check(ParseProxy("socks5://127.1.1.2").Scheme, Equals, "socks5")
	c.Check(ParseProxy("socks4://abc:foo@127.1.1.3:http/foo/bar").Scheme, Equals, "socks4")
}

func (s *ProxiesSuite) Test_ParseProxy_returnsTheHostOfTheProxySpecification(c *C) {
	c.Check(*ParseProxy("socks4://localhost").Host, Equals, "localhost")
	c.Check(*ParseProxy("socks5://127.1.1.2").Host, Equals, "127.1.1.2")
	c.Check(*ParseProxy("socks4://abc:foo@127.1.1.3:http/foo/bar").Host, Equals, "127.1.1.3")
}

func (s *ProxiesSuite) Test_ParseProxy_returnsThePortIfSet(c *C) {
	c.Check(ParseProxy("socks5://localhost").Port, IsNil)
	c.Check(ParseProxy("socks5://127.1.1.2").Port, IsNil)
	c.Check(*ParseProxy("socks5://abc:foo@127.1.1.3:http/foo/bar").Port, Equals, "http")
}

func (s *ProxiesSuite) Test_ParseProxy_returnsThePathIfSet(c *C) {
	c.Check(*ParseProxy("socks5+unix://abc:foo@/foo/bar").Path, Equals, "/foo/bar")
}

func (s *ProxiesSuite) Test_ParseProxy_returnsTheUserIfSet(c *C) {
	c.Check(ParseProxy("socks5://localhost").User, IsNil)
	c.Check(*ParseProxy("socks5://ooo@127.1.1.2").User, Equals, "ooo")
	c.Check(*ParseProxy("socks5://abc:foo@127.1.1.3:http/foo/bar").User, Equals, "abc")
}

func (s *ProxiesSuite) Test_ParseProxy_returnsThePasswordIfSet(c *C) {
	c.Check(ParseProxy("socks5://localhost").Pass, IsNil)
	c.Check(ParseProxy("socks5://ooo@127.1.1.2").Pass, IsNil)
	c.Check(*ParseProxy("socks5://abc:foo@127.1.1.3:http/foo/bar").Pass, Equals, "foo")
}

func (s *ProxiesSuite) Test_Proxy_ForPresentation_returnsAStringSuitableForUserPresentation(c *C) {
	c.Check(ParseProxy("socks5://localhost").ForPresentation(), Equals, "socks5://localhost")
	c.Check(ParseProxy("socks5://ooo@127.1.1.2").ForPresentation(), Equals, "socks5://ooo@127.1.1.2")
	c.Check(ParseProxy("socks5://abc:foo@127.1.1.3:http/foo/bar").ForPresentation(), Equals, "socks5://abc:*****@127.1.1.3:http/foo/bar")
	c.Check(ParseProxy("socks5+unix://abc:foo@/foo/bar").ForPresentation(), Equals, "socks5+unix://abc:*****@/foo/bar")
}

func (s *ProxiesSuite) Test_Proxy_ForProcessing_returnsAStringSuitableForProcessing(c *C) {
	c.Check(ParseProxy("socks5://localhost").ForProcessing(), Equals, "socks5://localhost")
	c.Check(ParseProxy("socks5://ooo@127.1.1.2").ForProcessing(), Equals, "socks5://ooo@127.1.1.2")
	c.Check(ParseProxy("socks5://abc:foo@127.1.1.3:http/foo/bar").ForProcessing(), Equals, "socks5://abc:foo@127.1.1.3:http/foo/bar")
	c.Check(ParseProxy("socks5+unix:///var/run/tor.fd").ForProcessing(), Equals, "socks5+unix:///var/run/tor.fd")
}
