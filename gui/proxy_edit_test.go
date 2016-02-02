package gui

import . "gopkg.in/check.v1"

type ProxiesSuite struct{}

var _ = Suite(&ProxiesSuite{})

func (s *ProxiesSuite) Test_FindProxyTypeFor_findsTheIndexForAProxyType(c *C) {
	c.Check(findProxyTypeFor("tor-auto"), Equals, 0)
	c.Check(findProxyTypeFor("socks5"), Equals, 1)
	c.Check(findProxyTypeFor("something-weird"), Equals, -1)
}

func (s *ProxiesSuite) Test_GetProxyTypeNames_yieldsAllProxyTypeNames(c *C) {
	result := []string{}
	getProxyTypeNames(func(s string) {
		result = append(result, s)
	})

	c.Assert(result, DeepEquals, []string{"Automatic Tor", "SOCKS5"})
}

func (s *ProxiesSuite) Test_GetProxyTypeFor_returnsTheCorrectProxyType(c *C) {
	c.Check(getProxyTypeFor("Automatic Tor"), Equals, "tor-auto")
	c.Check(getProxyTypeFor("SOCKS5"), Equals, "socks5")
	c.Check(getProxyTypeFor("something else"), Equals, "")
}
