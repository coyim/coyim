package xmpp

import (
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"net"

	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/xmpp/data"
	"github.com/twstrike/coyim/xmpp/errors"
	"github.com/twstrike/gotk3adapter/glib_mock"

	. "gopkg.in/check.v1"
)

type TCPSuite struct{}

func (*TCPSuite) SetUpSuite(c *C) {
	log.SetOutput(ioutil.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
}

var _ = Suite(&TCPSuite{})

func (s *TCPSuite) Test_newTCPConn_SkipsSRVAndConnectsToOriginDomain(c *C) {
	p := &mockProxy{}
	d := &dialer{
		JID: "foo@jabber.com",

		proxy: p,
		config: data.Config{
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
	d := &dialer{
		JID:           "foo@jabber.com",
		serverAddress: "jabber.im:5333",

		proxy: p,
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
	c.Check(d.config.SkipSRVLookup, Equals, true)

	c.Check(p, MatchesExpectations)
}

func (s *TCPSuite) Test_newTCPConn_ErrorsIfServiceIsNotAvailable(c *C) {
	p := &mockProxy{}
	d := &dialer{
		JID: "foo@jabber.com",

		proxy: p,
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
	d := &dialer{
		JID: "foo@jabber.com",

		proxy: p,
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

func (s *TCPSuite) Test_newTCPConn_ErrorsWhenTCPBindingFails(c *C) {
	p := &mockProxy{}
	d := &dialer{
		JID: "foo@jabber.com",

		proxy: p,
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
	c.Check(err, Equals, errors.ErrTCPBindingFailed)

	c.Check(p, MatchesExpectations)
}

func (s *TCPSuite) Test_newTCPConn_ErrorsWhenTCPBindingSucceedsButConnectionFails(c *C) {
	dec, _ := hex.DecodeString("00511eea818000010001000000000c5f786d70702d636c69656e74045f746370076f6c6162696e690273650000210001c00c0021000100000258001700000005146604786d7070076f6c6162696e6902736500")

	p := &mockProxy{}
	d := &dialer{
		JID: "foo@olabini.se",

		proxy: p,
	}

	p.Expects(func(network, addr string) (net.Conn, error) {
		c.Check(network, Equals, "tcp")
		c.Check(addr, Equals, "208.67.222.222:53")

		return fakeTCPConnToDNS(dec)
	})

	p.Expects(func(network, addr string) (net.Conn, error) {
		c.Check(network, Equals, "tcp")
		c.Check(addr, Equals, "xmpp.olabini.se:5222")

		return nil, io.EOF
	})

	_, err := d.newTCPConn()
	c.Check(err, Equals, errors.ErrConnectionFailed)
	c.Check(p, MatchesExpectations)
}
