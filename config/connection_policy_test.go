package config

import (
	"net"
	"net/url"
	"time"

	"github.com/twstrike/coyim/net/nettest"
	"golang.org/x/net/proxy"
	. "gopkg.in/check.v1"
)

type ConnectionPolicySuite struct{}

var _ = Suite(&ConnectionPolicySuite{})

func (s *ConnectionPolicySuite) Test_buildDialerFor_ValidatesJid(c *C) {
	account := &Account{
		Account: "invalid.com",
	}

	policy := ConnectionPolicy{}

	_, err := policy.buildDialerFor(account)

	c.Check(err.Error(), Equals, "invalid username (want user@domain): invalid.com")
}

func (s *ConnectionPolicySuite) Test_buildDialerFor_UsesCustomRootCAForJabberDotCCCDotDe(c *C) {
	account := &Account{
		Account: "coyim@jabber.ccc.de",
	}

	policy := ConnectionPolicy{}

	expectedRootCA, _ := rootCAFor("jabber.ccc.de")
	dialer, err := policy.buildDialerFor(account)

	c.Check(err, IsNil)
	c.Check(dialer.Config.TLSConfig.RootCAs.Subjects(),
		DeepEquals,
		expectedRootCA.Subjects(),
	)
}

func (s *ConnectionPolicySuite) Test_buildDialerFor_UsesConfiguredServerAddressAndPortAndMakesSRVLookup(c *C) {
	policy := ConnectionPolicy{}

	dialer, err := policy.buildDialerFor(&Account{
		Account: "coyim@coy.im",
		Server:  "xmpp.coy.im",
		Port:    5234,
	})

	c.Check(err, IsNil)
	c.Check(dialer.ServerAddress, Equals, "xmpp.coy.im:5234")

	dialer, err = policy.buildDialerFor(&Account{
		Account: "coyim@coy.im",
		Server:  "coy.im",
		Port:    5234,
	})

	c.Check(err, IsNil)
	c.Check(dialer.Config.SkipSRVLookup, Equals, false)
	c.Check(dialer.ServerAddress, Equals, "coy.im:5234")
}

func (s *ConnectionPolicySuite) Test_buildDialerFor_UsesAssociatedHiddenServiceIfFoundAndSkipsSRVLookup(c *C) {
	account := &Account{
		Account: "coyim@riseup.net",
	}

	policy := ConnectionPolicy{
		UseHiddenService: true,
	}

	dialer, err := policy.buildDialerFor(account)

	c.Check(err, IsNil)
	c.Check(dialer.Config.SkipSRVLookup, Equals, true)
	c.Check(dialer.ServerAddress, Equals, "4cjw6cwpeaeppfqz.onion:5222")
}

func (s *ConnectionPolicySuite) Test_buildDialerFor_IgnoresAssociatedHiddenService(c *C) {
	account := &Account{
		Account: "coyim@riseup.net",
	}

	policy := ConnectionPolicy{}

	dialer, err := policy.buildDialerFor(account)

	c.Check(err, IsNil)
	c.Check(dialer.ServerAddress, Equals, "")
}

func (s *ConnectionPolicySuite) Test_buildDialerFor_ErrorsIfTorIsRequiredButNotFound(c *C) {
	account := &Account{
		Account:    "coyim@riseup.net",
		RequireTor: true,
	}

	policy := ConnectionPolicy{
		torState: nettest.MockTorState(""),
	}

	_, err := policy.buildDialerFor(account)

	c.Check(err, Equals, ErrTorNotRunning)
}

func (s *ConnectionPolicySuite) Test_buildDialerFor_EnforcesTorProxyOnAccountIfRequired(c *C) {
	account := &Account{
		Account:    "coyim@riseup.net",
		RequireTor: true,
	}

	policy := ConnectionPolicy{
		torState: nettest.MockTorState("127.0.0.1:9999"),
	}

	c.Check(account.Proxies, HasLen, 0)

	_, err := policy.buildDialerFor(account)

	c.Check(err, IsNil)
	c.Check(account.Proxies, HasLen, 1)
}

func (s *ConnectionPolicySuite) Test_buildProxyChain_ErrorsIfProxyIsMalformed(c *C) {
	proxies := []string{
		"%gh&%ij",
	}

	_, err := buildProxyChain(proxies)
	c.Check(err.Error(), Equals, `Failed to parse %gh&%ij as a URL: parse %gh&%ij: invalid URL escape "%gh"`)
}

func (s *ConnectionPolicySuite) Test_buildProxyChain_ErrorsIfProxyIsNotCompatible(c *C) {
	proxies := []string{
		"socks4://proxy.local",
	}

	_, err := buildProxyChain(proxies)
	c.Check(err.Error(), Equals, "Failed to parse socks4://proxy.local as a proxy: proxy: unknown scheme: socks4")
}

func (s *ConnectionPolicySuite) Test_buildProxyChain_Returns(c *C) {
	proxies := []string{
		"socks5://proxy.local",
		"socks5://proxy.remote",
	}

	direct := &net.Dialer{Timeout: 30 * time.Second}
	p1, _ := proxy.FromURL(
		&url.URL{
			Scheme: "socks5",
			Host:   "proxy.remote",
		}, direct)

	expectedProxy, _ := proxy.FromURL(&url.URL{
		Scheme: "socks5",
		Host:   "proxy.local",
	}, p1)

	chain, err := buildProxyChain(proxies)
	c.Check(err, IsNil)
	c.Check(chain, DeepEquals, expectedProxy)

}
