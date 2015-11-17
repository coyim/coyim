package xmpp

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"

	. "gopkg.in/check.v1"
)

type TCPSuite struct{}

func (*TCPSuite) SetUpSuite(c *C) {
	log.SetOutput(ioutil.Discard)
}

func (*TCPSuite) TearDownSuite(c *C) {
	log.SetOutput(os.Stderr)
}

var _ = Suite(&TCPSuite{})

func (s *TCPSuite) Test_newTCPConn_SkipsSRVAndConnectsToOriginDomain(c *C) {
	p := &mockProxy{}
	d := &Dialer{
		JID: "foo@jabber.com",

		Proxy: p,
		Config: Config{
			SkipSRVLookup: true,
		},
	}

	expectedConn := &net.TCPConn{}
	p.Expects(func(network, addr string) (net.Conn, error) {
		c.Check(network, Equals, "tcp")
		c.Check(addr, Equals, "jabber.com:5222")

		return expectedConn, nil
	})

	conn, err := d.newTCPConn()
	c.Check(err, IsNil)
	c.Check(conn, Equals, expectedConn)

	c.Check(p, MatchesExpectations)
}

func (s *TCPSuite) Test_newTCPConn_SkipsSRVAndConnectsToConfiguredServerAddress(c *C) {
	p := &mockProxy{}
	d := &Dialer{
		JID:           "foo@jabber.com",
		ServerAddress: "jabber.im:5333",

		Proxy: p,
		Config: Config{
			SkipSRVLookup: true,
		},
	}

	expectedConn := &net.TCPConn{}
	p.Expects(func(network, addr string) (net.Conn, error) {
		c.Check(network, Equals, "tcp")
		c.Check(addr, Equals, "jabber.im:5333")

		return expectedConn, nil
	})

	conn, err := d.newTCPConn()
	c.Check(err, IsNil)
	c.Check(conn, Equals, expectedConn)

	c.Check(p, MatchesExpectations)
}

func (s *TCPSuite) Test_newTCPConn_ErrorsOnMalformedServerAddress(c *C) {
	p := &mockProxy{}
	d := &Dialer{
		JID:           "foo@jabber.com",
		ServerAddress: "jabber.im",

		Proxy: p,
	}

	_, err := d.newTCPConn()
	c.Check(err.Error(), Equals, "missing port in address jabber.im")

	c.Check(p, MatchesExpectations)
}

func (s *TCPSuite) Test_newTCPConn_ErrorsIfServiceIsNotAvailable(c *C) {
	p := &mockProxy{}
	d := &Dialer{
		JID: "foo@jabber.com",

		Proxy: p,
	}

	// We exploit ResolveSRVWithProxy forwarding conn errors
	// to fake an error it should generated.
	p.Expects(func(network, addr string) (net.Conn, error) {
		c.Check(network, Equals, "tcp")
		c.Check(addr, Equals, "208.67.222.222:53")

		return nil, ErrServiceNotAvailable
	})

	_, err := d.newTCPConn()
	c.Check(err, Equals, ErrServiceNotAvailable)

	c.Check(p, MatchesExpectations)
}

func (s *TCPSuite) Test_newTCPConn_DefaultsToOriginDomainAtDefaultPortAfterSRVFails(c *C) {
	p := &mockProxy{}
	d := &Dialer{
		JID: "foo@jabber.com",

		Proxy: p,
	}

	p.Expects(func(network, addr string) (net.Conn, error) {
		c.Check(network, Equals, "tcp")
		c.Check(addr, Equals, "208.67.222.222:53")

		return nil, io.EOF
	})

	expectedConn := &net.TCPConn{}
	p.Expects(func(network, addr string) (net.Conn, error) {
		c.Check(network, Equals, "tcp")
		c.Check(addr, Equals, "jabber.com:5222")

		return expectedConn, nil
	})

	conn, err := d.newTCPConn()
	c.Check(err, IsNil)
	c.Check(conn, Equals, expectedConn)

	c.Check(p, MatchesExpectations)
}

func (s *TCPSuite) Test_newTCPConn_IgnoresConfiguredServerWhenFallingBackAfterSRVFailure(c *C) {
	p := &mockProxy{}
	d := &Dialer{
		JID:           "foo@jabber.com",
		ServerAddress: "jabber.im:5333",

		Proxy: p,
	}

	p.Expects(func(network, addr string) (net.Conn, error) {
		c.Check(network, Equals, "tcp")
		c.Check(addr, Equals, "208.67.222.222:53")

		return nil, io.EOF
	})

	expectedConn := &net.TCPConn{}
	p.Expects(func(network, addr string) (net.Conn, error) {
		c.Check(network, Equals, "tcp")
		c.Check(addr, Equals, "jabber.com:5222")

		return expectedConn, nil
	})

	conn, err := d.newTCPConn()
	c.Check(err, IsNil)
	c.Check(conn, Equals, expectedConn)

	c.Check(p, MatchesExpectations)
}

func (s *TCPSuite) Test_newTCPConn_ErrorsWhenFallbackMethodFails(c *C) {
	p := &mockProxy{}
	d := &Dialer{
		JID: "foo@jabber.com",

		Proxy: p,
	}

	p.Expects(func(network, addr string) (net.Conn, error) {
		c.Check(network, Equals, "tcp")
		c.Check(addr, Equals, "208.67.222.222:53")

		return nil, io.EOF
	})

	p.Expects(func(network, addr string) (net.Conn, error) {
		c.Check(network, Equals, "tcp")
		c.Check(addr, Equals, "jabber.com:5222")

		return nil, io.EOF
	})

	_, err := d.newTCPConn()
	c.Check(err, Equals, ErrConnectionFailed)

	c.Check(p, MatchesExpectations)
}
