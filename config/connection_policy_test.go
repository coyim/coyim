package config

import . "gopkg.in/check.v1"

type ConnectionPolicySuite struct{}

var _ = Suite(&ConnectionPolicySuite{})

func (s *ConnectionPolicySuite) TestBuildDialerFor_ValidatesJid(c *C) {
	account := &Account{
		Account: "invalid.com",
	}

	policy := ConnectionPolicy{}

	_, err := policy.buildDialerFor(account)

	c.Check(err.Error(), Equals, "invalid username (want user@domain): invalid.com")
}

func (s *ConnectionPolicySuite) TestBuildDialerFor_UsesCustomRootCAForJabberDotCCCDotDe(c *C) {
	account := &Account{
		Account: "coyim@jabber.ccc.de",
	}

	policy := ConnectionPolicy{}

	expectedRootCA, _ := rootCAFor("jabber.ccc.de")
	dialer, err := policy.buildDialerFor(account)

	c.Check(err, Equals, nil)
	c.Check(dialer.Config.TLSConfig.RootCAs.Subjects(),
		DeepEquals,
		expectedRootCA.Subjects(),
	)
}

func (s *ConnectionPolicySuite) TestBuildDialerFor_UsesConfiguredServerAddressAndPortAndMakesSRVLookup(c *C) {
	policy := ConnectionPolicy{}

	dialer, err := policy.buildDialerFor(&Account{
		Account: "coyim@coy.im",
		Server:  "xmpp.coy.im",
		Port:    5234,
	})

	c.Check(err, Equals, nil)
	c.Check(dialer.ServerAddress, Equals, "xmpp.coy.im:5234")

	dialer, err = policy.buildDialerFor(&Account{
		Account: "coyim@coy.im",
		Server:  "coy.im",
		Port:    5234,
	})

	c.Check(err, Equals, nil)
	c.Check(dialer.Config.SkipSRVLookup, Equals, false)
	c.Check(dialer.ServerAddress, Equals, "coy.im:5234")
}

func (s *ConnectionPolicySuite) TestBuildDialerFor_UsesAssociatedHiddenServiceIfFoundAndSkipsSRVLookup(c *C) {
	account := &Account{
		Account: "coyim@riseup.net",
	}

	policy := ConnectionPolicy{
		UseHiddenService: true,
	}

	dialer, err := policy.buildDialerFor(account)

	c.Check(err, Equals, nil)
	c.Check(dialer.Config.SkipSRVLookup, Equals, true)
	c.Check(dialer.ServerAddress, Equals, "4cjw6cwpeaeppfqz.onion:5222")
}

func (s *ConnectionPolicySuite) TestBuildDialerFor_IgnoresAssociatedHiddenService(c *C) {
	account := &Account{
		Account: "coyim@riseup.net",
	}

	policy := ConnectionPolicy{}

	dialer, err := policy.buildDialerFor(account)

	c.Check(err, Equals, nil)
	c.Check(dialer.ServerAddress, Equals, "")
}
