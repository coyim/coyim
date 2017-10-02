package config

import (
	"net"
	"net/url"
	"time"

	"golang.org/x/net/proxy"

	ournet "github.com/coyim/coyim/net"
	"github.com/coyim/coyim/xmpp"
	. "gopkg.in/check.v1"
)

type ConnectionPolicySuite struct{}

var _ = Suite(&ConnectionPolicySuite{})

func mockTorState(addr string, overTor bool) ournet.TorState {
	return &torStateMock{addr, overTor}
}

type torStateMock struct {
	addr    string
	overTor bool
}

func (s *torStateMock) Address() string {
	return s.addr
}

func (s *torStateMock) Detect() bool {
	return len(s.addr) > 0
}

func (s *torStateMock) IsConnectionOverTor(proxy.Dialer) bool {
	return s.overTor
}

func (s *ConnectionPolicySuite) Test_buildDialerFor_ValidatesJid(c *C) {
	account := &Account{
		Account: "invalid.com",
	}

	policy := ConnectionPolicy{DialerFactory: xmpp.DialerFactory}

	_, err := policy.buildDialerFor(account, nil)

	c.Check(err.Error(), Equals, "invalid username (want user@domain): invalid.com")
}

func (s *ConnectionPolicySuite) Test_buildDialerFor_UsesCustomRootCAForJabberDotCCCDotDe(c *C) {
	account := &Account{
		Account: "coyim@jabber.ccc.de",
	}

	policy := ConnectionPolicy{DialerFactory: xmpp.DialerFactory, torState: mockTorState("", false)}

	expectedRootCA, _ := rootCAFor("jabber.ccc.de")
	dialer, err := policy.buildDialerFor(account, nil)

	c.Check(err, IsNil)
	c.Check(dialer.Config().TLSConfig.RootCAs.Subjects(),
		DeepEquals,
		expectedRootCA.Subjects(),
	)
}

func (s *ConnectionPolicySuite) Test_buildDialerFor_UsesConfiguredServerAddressAndPortAndMakesSRVLookup(c *C) {
	policy := ConnectionPolicy{DialerFactory: xmpp.DialerFactory, torState: mockTorState("", false)}

	dialer, err := policy.buildDialerFor(&Account{
		Account: "coyim@coy.im",
		Server:  "xmpp.coy.im",
		Port:    5234,
	}, nil)

	c.Check(err, IsNil)
	c.Check(dialer.ServerAddress(), Equals, "xmpp.coy.im:5234")

	dialer, err = policy.buildDialerFor(&Account{
		Account: "coyim@coy.im",
		Server:  "coy.im",
		Port:    5234,
	}, nil)

	c.Check(err, IsNil)
	c.Check(dialer.Config().SkipSRVLookup, Equals, false)
	c.Check(dialer.ServerAddress(), Equals, "coy.im:5234")
}

func (s *ConnectionPolicySuite) Test_buildDialerFor_UsesAssociatedHiddenServiceIfFound(c *C) {
	account := &Account{
		Account: "coyim@riseup.net",
		Proxies: []string{
			"tor-auto://",
		},
	}

	currentTor := ournet.Tor

	ournet.Tor = mockTorState("127.0.0.1:9999", true)
	policy := ConnectionPolicy{
		DialerFactory: xmpp.DialerFactory,
		torState:      ournet.Tor,
	}
	dialer, err := policy.buildDialerFor(account, nil)
	ournet.Tor = currentTor

	c.Check(err, IsNil)
	c.Check(dialer.ServerAddress(), Equals, "4cjw6cwpeaeppfqz.onion:5222")
}

func (s *ConnectionPolicySuite) Test_buildDialerFor_IgnoresAssociatedHiddenService(c *C) {
	account := &Account{
		Account: "coyim@riseup.net",
	}

	policy := ConnectionPolicy{DialerFactory: xmpp.DialerFactory, torState: mockTorState("", false)}

	dialer, err := policy.buildDialerFor(account, nil)

	c.Check(err, IsNil)
	c.Check(dialer.ServerAddress(), Equals, "")
}

func (s *ConnectionPolicySuite) Test_buildDialerFor_ErrorsIfTorIsRequiredButNotFound(c *C) {
	account := &Account{
		Account: "coyim@riseup.net",
		Proxies: []string{"tor-auto://"},
	}

	policy := ConnectionPolicy{
		DialerFactory: xmpp.DialerFactory,
		torState:      mockTorState("", false),
	}

	_, err := policy.buildDialerFor(account, nil)

	c.Check(err, Equals, ErrTorNotRunning)
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

	direct := &net.Dialer{Timeout: 60 * time.Second}
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
