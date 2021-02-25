package config

import (
	"bytes"
	"crypto/x509"
	"errors"
	"fmt"

	"golang.org/x/net/proxy"

	"github.com/coyim/coyim/coylog"
	ournet "github.com/coyim/coyim/net"
	"github.com/coyim/coyim/servers"
	ourtls "github.com/coyim/coyim/tls"
	"github.com/coyim/coyim/xmpp"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
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

func (s *ConnectionPolicySuite) Test_buildDialerFor_failsIfCantBuildCAForJabberCCCDE(c *C) {
	account := &Account{
		Account: "coyim@jabber.ccc.de",
	}

	policy := ConnectionPolicy{DialerFactory: xmpp.DialerFactory, torState: mockTorState("", false)}

	origX509ParseCertificate := x509ParseCertificate
	defer func() {
		x509ParseCertificate = origX509ParseCertificate
	}()
	x509ParseCertificate = func([]byte) (*x509.Certificate, error) {
		return nil, errors.New("oh nooooooooooooooo")
	}

	dialer, err := policy.buildDialerFor(account, nil)

	c.Assert(err, ErrorMatches, "oh no+")
	c.Assert(dialer, IsNil)
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

func (s *ConnectionPolicySuite) Test_buildDialerFor_FailsIfItCantDoProxyCreation(c *C) {
	policy := ConnectionPolicy{DialerFactory: xmpp.DialerFactory, torState: mockTorState("", false)}

	dialer, err := buildDialerFor(&policy, &Account{
		Account: "coyim@coy.im",
		Server:  "xmpp.coy.im",
		Port:    5234,
		Proxies: []string{
			"%gh&%ij",
		},
	}, nil)

	c.Assert(err, ErrorMatches, "Failed to parse.*")
	c.Assert(dialer, IsNil)
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
	c.Check(err, ErrorMatches, `Failed to parse %gh&%ij as a URL: parse "?%gh&%ij"?: invalid URL escape "%gh"`)
}

func (s *ConnectionPolicySuite) Test_buildProxyChain_ErrorsIfProxyIsNotCompatible(c *C) {
	proxies := []string{
		"socks4://proxy.local",
	}

	_, err := buildProxyChain(proxies)
	c.Check(err.Error(), Equals, "Failed to parse socks4://proxy.local as a proxy: proxy: unknown scheme: socks4")
}

func (s *ConnectionPolicySuite) Test_buildProxyChain_Returns(c *C) {
	chain, err := buildProxyChain([]string{
		"socks5://proxy.local",
		"socks5://proxy.remote",
	})
	c.Check(err, IsNil)
	c.Check(fmt.Sprintf("%#v", chain), Matches, ".*socks.*?proxy\\.local.*")

	chain, err = buildProxyChain([]string{
		"socks5://proxy.remote",
	})
	c.Check(err, IsNil)
	c.Check(fmt.Sprintf("%#v", chain), Matches, ".*socks.*?proxy\\.remote.*")

	chain, err = buildProxyChain([]string{})
	c.Check(err, IsNil)
	c.Check(chain, IsNil)
}

func (s *ConnectionPolicySuite) Test_Account_CreateTorProxy_doesntUseTorAuto(c *C) {
	a := &Account{}
	chain, e := a.CreateTorProxy()
	c.Assert(chain, IsNil)
	c.Assert(e, IsNil)
}

func (s *ConnectionPolicySuite) Test_Account_CreateTorProxy_withTorAuto_triesToDetect(c *C) {
	origTorDetect := torDetect
	defer func() {
		torDetect = origTorDetect
	}()

	origTor := ournet.Tor
	defer func() {
		ournet.Tor = origTor
	}()

	a := &Account{
		Proxies: []string{"tor-auto://"},
	}

	called := false
	torDetect = func() bool {
		ournet.Tor = mockTorState("127.0.0.1:9999", true)
		called = true
		return true
	}

	chain, e := a.CreateTorProxy()
	c.Assert(chain, Not(IsNil))
	c.Assert(e, IsNil)
	c.Assert(called, Equals, true)
}

func (s *ConnectionPolicySuite) Test_Account_CreateTorProxy_withTorAuto_failsToDetect(c *C) {
	origTorDetect := torDetect
	defer func() {
		torDetect = origTorDetect
	}()

	a := &Account{
		Proxies: []string{"tor-auto://"},
	}

	torDetect = func() bool {
		return false
	}

	chain, e := a.CreateTorProxy()
	c.Assert(chain, IsNil)
	c.Assert(e, ErrorMatches, "Tor is not running")
}

func (s *ConnectionPolicySuite) Test_ConnectionPolicy_initTorState_usesOurnetTorIfNoneSet(c *C) {
	cp := &ConnectionPolicy{}
	cp.initTorState()

	c.Assert(cp.torState, Equals, ournet.Tor)
}

func (s *ConnectionPolicySuite) Test_torDetect_detects(c *C) {
	origTor := ournet.Tor
	defer func() {
		ournet.Tor = origTor
	}()

	ournet.Tor = mockTorState("127.0.0.1:9999", true)

	c.Assert(torDetect(), Equals, true)
}

func (s *ConnectionPolicySuite) Test_buildInOutLogs_createsLogger(c *C) {
	raw := bytes.Buffer{}
	in, out := buildInOutLogs(&raw)

	inr := in.(*rawLogger)
	outr := out.(*rawLogger)

	c.Assert(inr.out, Equals, &raw)
	c.Assert(outr.out, Equals, &raw)

	c.Assert(inr.prefix, DeepEquals, []byte("<- "))
	c.Assert(outr.prefix, DeepEquals, []byte("-> "))

	c.Assert(inr.lock, Equals, outr.lock)

	c.Assert(inr.other, Equals, outr)
	c.Assert(outr.other, Equals, inr)
}

func (s *ConnectionPolicySuite) Test_ConnectionPolicy_Connect_failsIfBuildingDialerFails(c *C) {
	origBuildDialerForFunc := buildDialerFor
	defer func() {
		buildDialerFor = origBuildDialerForFunc
	}()

	buildDialerFor = func(p *ConnectionPolicy, conf *Account, verifier ourtls.Verifier) (interfaces.Dialer, error) {
		return nil, errors.New("foooo")
	}

	cp := &ConnectionPolicy{}
	conn, e := cp.Connect("", "", nil, nil)
	c.Assert(e, ErrorMatches, "foooo")
	c.Assert(conn, IsNil)
}

type mockDialer struct {
	argPassword         string
	argResource         string
	argShouldConnectTLS bool
	argShouldSendALPN   bool

	returnDialConn interfaces.Conn
	returnDialErr  error

	returnRegisterConn interfaces.Conn
	returnRegisterErr  error
}

func (md *mockDialer) Config() data.Config { return data.Config{} }
func (md *mockDialer) Dial() (interfaces.Conn, error) {
	return md.returnDialConn, md.returnDialErr
}
func (md *mockDialer) GetServer() string { return "" }
func (md *mockDialer) RegisterAccount(data.FormCallback) (interfaces.Conn, error) {
	return md.returnRegisterConn, md.returnRegisterErr
}
func (md *mockDialer) ServerAddress() string { return "" }
func (md *mockDialer) SetConfig(data.Config) {}
func (md *mockDialer) SetJID(string)         {}
func (md *mockDialer) SetPassword(v string) {
	md.argPassword = v
}
func (md *mockDialer) SetProxy(proxy.Dialer) {}
func (md *mockDialer) SetResource(v string) {
	md.argResource = v
}
func (md *mockDialer) SetServerAddress(string) {}
func (md *mockDialer) SetShouldConnectTLS(v bool) {
	md.argShouldConnectTLS = v
}
func (md *mockDialer) SetShouldSendALPN(v bool) {
	md.argShouldSendALPN = v
}
func (md *mockDialer) SetLogger(coylog.Logger)  {}
func (md *mockDialer) SetKnown(*servers.Server) {}

func (s *ConnectionPolicySuite) Test_ConnectionPolicy_Connect_succeedsAndDials(c *C) {
	origBuildDialerForFunc := buildDialerFor
	defer func() {
		buildDialerFor = origBuildDialerForFunc
	}()

	expConn := xmpp.NewConn(nil, nil, "")

	dialer := &mockDialer{
		returnDialConn: expConn,
	}
	buildDialerFor = func(p *ConnectionPolicy, conf *Account, verifier ourtls.Verifier) (interfaces.Dialer, error) {
		return dialer, nil
	}

	cp := &ConnectionPolicy{}
	a := &Account{
		ConnectTLS: true,
		SetALPN:    true,
	}
	conn, e := cp.Connect("p1", "r1", a, nil)
	c.Assert(e, IsNil)
	c.Assert(conn, Equals, expConn)
	c.Assert(dialer.argShouldConnectTLS, Equals, true)
	c.Assert(dialer.argShouldSendALPN, Equals, true)
	c.Assert(dialer.argPassword, Equals, "p1")
	c.Assert(dialer.argResource, Equals, "r1")
}

func (s *ConnectionPolicySuite) Test_ConnectionPolicy_RegisterAccount_failsIfBuildingDialerFails(c *C) {
	origBuildDialerForFunc := buildDialerFor
	defer func() {
		buildDialerFor = origBuildDialerForFunc
	}()

	buildDialerFor = func(p *ConnectionPolicy, conf *Account, verifier ourtls.Verifier) (interfaces.Dialer, error) {
		return nil, errors.New("foooo")
	}

	cp := &ConnectionPolicy{}
	conn, e := cp.RegisterAccount(nil, nil, nil)
	c.Assert(e, ErrorMatches, "foooo")
	c.Assert(conn, IsNil)
}

func (s *ConnectionPolicySuite) Test_ConnectionPolicy_RegisterAccount_succeedsAndRegisters(c *C) {
	origBuildDialerForFunc := buildDialerFor
	defer func() {
		buildDialerFor = origBuildDialerForFunc
	}()

	expConn := xmpp.NewConn(nil, nil, "")

	dialer := &mockDialer{
		returnRegisterConn: expConn,
	}
	buildDialerFor = func(p *ConnectionPolicy, conf *Account, verifier ourtls.Verifier) (interfaces.Dialer, error) {
		return dialer, nil
	}

	cp := &ConnectionPolicy{}
	conn, e := cp.RegisterAccount(nil, nil, nil)
	c.Assert(e, IsNil)
	c.Assert(conn, Equals, expConn)
}

func (s *ConnectionPolicySuite) Test_ConnectionPolicy_RegisterAccount_returnsErrorFromRegistration(c *C) {
	origBuildDialerForFunc := buildDialerFor
	defer func() {
		buildDialerFor = origBuildDialerForFunc
	}()

	dialer := &mockDialer{
		returnRegisterErr: errors.New("reg fail"),
	}
	buildDialerFor = func(p *ConnectionPolicy, conf *Account, verifier ourtls.Verifier) (interfaces.Dialer, error) {
		return dialer, nil
	}

	cp := &ConnectionPolicy{}
	conn, e := cp.RegisterAccount(nil, nil, nil)
	c.Assert(e, ErrorMatches, "reg fail")
	c.Assert(conn, IsNil)
}
