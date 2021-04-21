package filetransfer

import (
	"errors"
	"net"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/xmpp/data"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"golang.org/x/net/proxy"
	. "gopkg.in/check.v1"
)

type BytestreamsSuite struct{}

var _ = Suite(&BytestreamsSuite{})

type mockConfigAndLog struct {
	c *config.Account
	l coylog.Logger
}

func (m *mockConfigAndLog) GetConfig() *config.Account {
	return m.c
}

func (m *mockConfigAndLog) Log() coylog.Logger {
	return m.l
}

type mockDialer struct {
	f func(network, addr string) (net.Conn, error)
}

func (m *mockDialer) Dial(network, addr string) (net.Conn, error) {
	return m.f(network, addr)
}

func (s *BytestreamsSuite) Test_createTorProxy(c *C) {
	a := &config.Account{}
	res, e := createTorProxy(a)
	c.Assert(res, IsNil)
	c.Assert(e, IsNil)
}

func (s *BytestreamsSuite) Test_tryStreamhost_succeeds(c *C) {
	orgCreateTorProxy := createTorProxy
	defer func() {
		createTorProxy = orgCreateTorProxy
	}()

	orgSocks5XMPP := socks5XMPP
	defer func() {
		socks5XMPP = orgSocks5XMPP
	}()

	createTorProxy = func(a *config.Account) (proxy.Dialer, error) {
		return nil, nil
	}

	md := &mockDialer{
		f: func(network, addr string) (net.Conn, error) {
			return nil, nil
		},
	}

	socks5XMPP = func(_, _ string, _ *proxy.Auth, _ proxy.Dialer) (proxy.Dialer, error) {
		return md, nil
	}

	mockSess := &mockConfigAndLog{
		c: nil,
		l: nil,
	}

	sh := data.BytestreamStreamhost{}

	called := false
	f := func(net.Conn) {
		called = true
	}

	res := tryStreamhost(mockSess, sh, "hello.com", f)
	c.Assert(res, Equals, true)
	c.Assert(called, Equals, true)
}

func (s *BytestreamsSuite) Test_tryStreamhost_failsIfGettingProxyFails(c *C) {
	orgCreateTorProxy := createTorProxy
	defer func() {
		createTorProxy = orgCreateTorProxy
	}()

	createTorProxy = func(a *config.Account) (proxy.Dialer, error) {
		return nil, errors.New("marker error 1")
	}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockSess := &mockConfigAndLog{
		c: nil,
		l: l,
	}

	sh := data.BytestreamStreamhost{}

	f := func(net.Conn) {}

	res := tryStreamhost(mockSess, sh, "hello.com", f)
	c.Assert(res, Equals, false)
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Had error when trying to connect")
	c.Assert(hook.Entries[0].Data, HasLen, 1)
	c.Assert(hook.Entries[0].Data["error"], ErrorMatches, "marker error 1")
}

func (s *BytestreamsSuite) Test_tryStreamhost_failsCreatingSocks(c *C) {
	orgCreateTorProxy := createTorProxy
	defer func() {
		createTorProxy = orgCreateTorProxy
	}()

	orgSocks5XMPP := socks5XMPP
	defer func() {
		socks5XMPP = orgSocks5XMPP
	}()

	createTorProxy = func(a *config.Account) (proxy.Dialer, error) {
		return nil, nil
	}

	socks5XMPP = func(_, _ string, _ *proxy.Auth, _ proxy.Dialer) (proxy.Dialer, error) {
		return nil, errors.New("marker error 22")
	}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockSess := &mockConfigAndLog{
		c: nil,
		l: l,
	}

	sh := data.BytestreamStreamhost{}

	f := func(net.Conn) {}

	res := tryStreamhost(mockSess, sh, "hello.com", f)
	c.Assert(res, Equals, false)
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Error setting up socks5")
	c.Assert(hook.Entries[0].Data["error"], ErrorMatches, "marker error 22")
}

func (s *BytestreamsSuite) Test_tryStreamhost_failsDialing(c *C) {
	orgCreateTorProxy := createTorProxy
	defer func() {
		createTorProxy = orgCreateTorProxy
	}()

	orgSocks5XMPP := socks5XMPP
	defer func() {
		socks5XMPP = orgSocks5XMPP
	}()

	createTorProxy = func(a *config.Account) (proxy.Dialer, error) {
		return nil, nil
	}

	md := &mockDialer{
		f: func(network, addr string) (net.Conn, error) {
			return nil, errors.New("marker error 3")
		},
	}

	socks5XMPP = func(_, _ string, _ *proxy.Auth, _ proxy.Dialer) (proxy.Dialer, error) {
		return md, nil
	}

	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	mockSess := &mockConfigAndLog{
		c: nil,
		l: l,
	}

	sh := data.BytestreamStreamhost{}

	f := func(net.Conn) {
	}

	res := tryStreamhost(mockSess, sh, "hello.com", f)
	c.Assert(res, Equals, false)
	c.Assert(hook.Entries, HasLen, 1)
	c.Assert(hook.Entries[0].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[0].Message, Equals, "Error connecting socks5")
	c.Assert(hook.Entries[0].Data["error"], ErrorMatches, "marker error 3")
}
