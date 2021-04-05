package xmpp

import (
	gotls "crypto/tls"
	"errors"
	"net"
	"time"

	"github.com/coyim/coyim/servers"
	"github.com/coyim/coyim/tls"
	"github.com/coyim/coyim/xmpp/data"
	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
)

type DialerSuite struct{}

var _ = Suite(&DialerSuite{})

func (s *DialerSuite) Test_DialerFactory(c *C) {
	dd := DialerFactory(nil, nil)
	c.Assert(dd, Not(IsNil))
}

func (s *DialerSuite) Test_dialer_SetJID(c *C) {
	dd := &dialer{}
	dd.SetJID("helo@fmm")
	c.Assert(dd.JID, Equals, "helo@fmm")
	c.Assert(dd.getJIDLocalpart(), Equals, "helo")
	c.Assert(dd.getJIDDomainpart(), Equals, "fmm")
}

func (s *DialerSuite) Test_dialer_SetKnown(c *C) {
	dd := &dialer{}
	kn := &servers.Server{}
	dd.SetKnown(kn)
	c.Assert(dd.known, Equals, kn)
}

func (s *DialerSuite) Test_dialer_ServerAddress(c *C) {
	dd := &dialer{JID: "hmm@haha.com"}
	c.Assert(dd.hasCustomServer(), Equals, false)
	c.Assert(dd.GetServer(), Equals, "haha.com:5222")
	dd.SetServerAddress("hmm@bla")
	c.Assert(dd.ServerAddress(), Equals, "hmm@bla")
	c.Assert(dd.hasCustomServer(), Equals, true)
	c.Assert(dd.GetServer(), Equals, "hmm@bla")
}

func (s *DialerSuite) Test_dialer_settingProperties(c *C) {
	dd := &dialer{}

	dd.SetShouldConnectTLS(true)
	c.Assert(dd.connectTLS, Equals, true)

	dd.SetShouldSendALPN(true)
	c.Assert(dd.sendALPN, Equals, true)

	dd.SetPassword("hello pass")
	c.Assert(dd.password, Equals, "hello pass")

	dd.SetResource("hello resource")
	c.Assert(dd.resource, Equals, "hello resource")

	pp := &net.Dialer{Timeout: 60 * time.Second}
	dd.SetProxy(pp)
	c.Assert(dd.proxy, Equals, pp)

	conf := data.Config{Archive: true}
	dd.SetConfig(conf)
	c.Assert(dd.Config(), DeepEquals, conf)

	ll := log.New()
	dd.SetLogger(ll)
	c.Assert(dd.log, Equals, ll)
}

func (s *DialerSuite) Test_dialer_Dial(c *C) {
	p := &mockProxy{}
	d := &dialer{
		JID: "foo@jabber.com",

		proxy: p,
		config: data.Config{
			SkipSRVLookup: true,
			SkipTLS:       true,
		},
		log: testLogger(),
	}

	expectedConn := &fullMockedConn{
		&mockConnIOReaderWriter{
			read: []byte(`<stream xmlns="http://etherx.jabber.org/streams">
  <features>
    <mechanisms xmlns="urn:ietf:params:xml:ns:xmpp-sasl">
      <mechanism>PLAIN</mechanism>
    </mechanisms>
  </features>
  <success xmlns="urn:ietf:params:xml:ns:xmpp-sasl"></success>
  <streams:stream xmlns:streams="http://etherx.jabber.org/streams" xmlns="jabber:client">
  <streams:features>
    <mechanisms xmlns="urn:ietf:params:xml:ns:xmpp-sasl">
      <mechanism>PLAIN</mechanism>
    </mechanisms>
  </streams:features>
  <iq type="result" id="bind_1"></iq>
  <iq type="result" id="sess_1"></iq>
`),
		},
	}
	p.Expects(func(network, addr string) (net.Conn, error) {
		return expectedConn, nil
	})

	cn, e := d.Dial()
	c.Assert(e, IsNil)
	c.Assert(cn, Not(IsNil))
	cc := cn.(*conn)
	c.Assert(cc.rawOut, Equals, expectedConn)
}

func (s *DialerSuite) Test_dialer_Dial_failsOnTCPConn(c *C) {
	p := &mockProxy{}
	d := &dialer{
		JID: "foo@jabber.com",

		proxy: p,
		config: data.Config{
			SkipSRVLookup: true,
			SkipTLS:       true,
		},
		log: testLogger(),
	}

	p.Expects(func(network, addr string) (net.Conn, error) {
		return nil, errors.New("marker interface")
	})

	cn, e := d.Dial()
	c.Assert(e, ErrorMatches, ".*marker interface")
	c.Assert(cn, IsNil)
}

func (s *DialerSuite) Test_dialer_RegisterAccount(c *C) {
	p := &mockProxy{}
	d := &dialer{
		JID: "foo@jabber.com",

		proxy: p,
		config: data.Config{
			SkipSRVLookup: true,
			SkipTLS:       true,
		},
		log: testLogger(),
	}

	expectedConn := &fullMockedConn{
		&mockConnIOReaderWriter{
			read: []byte(`<stream xmlns="http://etherx.jabber.org/streams">
  <features>
    <mechanisms xmlns="urn:ietf:params:xml:ns:xmpp-sasl">
      <mechanism>PLAIN</mechanism>
    </mechanisms>
  </features>
  <success xmlns="urn:ietf:params:xml:ns:xmpp-sasl"></success>
  <streams:stream xmlns:streams="http://etherx.jabber.org/streams" xmlns="jabber:client">
  <streams:features>
    <mechanisms xmlns="urn:ietf:params:xml:ns:xmpp-sasl">
      <mechanism>PLAIN</mechanism>
    </mechanisms>
  </streams:features>
  <iq type="result" id="bind_1"></iq>
  <iq type="result" id="sess_1"></iq>
`),
		},
	}
	p.Expects(func(network, addr string) (net.Conn, error) {
		return expectedConn, nil
	})

	called := false
	ff := func(string, string, []interface{}) error {
		called = true
		return nil
	}

	cn, e := d.RegisterAccount(ff)
	c.Assert(e, IsNil)
	c.Assert(cn, Not(IsNil))
	cc := cn.(*conn)
	c.Assert(cc.rawOut, Equals, expectedConn)
	_ = d.config.CreateCallback("", "", nil)
	c.Assert(called, Equals, true)
}

type fullMockedTLSConn struct {
	*fullMockedConn

	e error
}

func (m *fullMockedTLSConn) Handshake() error {
	return m.e
}

func (m *fullMockedTLSConn) ConnectionState() gotls.ConnectionState {
	return gotls.ConnectionState{}
}

func (s *DialerSuite) Test_dialer_Dial_withOuterTLS(c *C) {
	p := &mockProxy{}

	mtls := &fullMockedTLSConn{
		e: errors.New("another marker"),
	}

	d := &dialer{
		JID: "foo@jabber.com",

		proxy: p,
		config: data.Config{
			SkipSRVLookup: true,
			SkipTLS:       true,
		},
		log:      testLogger(),
		sendALPN: true,

		connectTLS: true,

		tlsConnFactory: func(net.Conn, *gotls.Config) tls.Conn {
			return mtls
		},
	}

	expectedConn := &fullMockedConn{
		&mockConnIOReaderWriter{
			read: []byte(`<stream xmlns="http://etherx.jabber.org/streams">
  <features>
    <mechanisms xmlns="urn:ietf:params:xml:ns:xmpp-sasl">
      <mechanism>PLAIN</mechanism>
    </mechanisms>
  </features>
  <success xmlns="urn:ietf:params:xml:ns:xmpp-sasl"></success>
  <streams:stream xmlns:streams="http://etherx.jabber.org/streams" xmlns="jabber:client">
  <streams:features>
    <mechanisms xmlns="urn:ietf:params:xml:ns:xmpp-sasl">
      <mechanism>PLAIN</mechanism>
    </mechanisms>
  </streams:features>
  <iq type="result" id="bind_1"></iq>
  <iq type="result" id="sess_1"></iq>
`),
		},
	}
	p.Expects(func(network, addr string) (net.Conn, error) {
		return expectedConn, nil
	})

	cn, e := d.Dial()
	c.Assert(e, ErrorMatches, "another marker")
	c.Assert(cn, IsNil)
}
