package xmpp

import (
	"encoding/hex"
	"fmt"
	"net"

	. "gopkg.in/check.v1"
)

type DNSXMPPSuite struct{}

var _ = Suite(&DNSXMPPSuite{})

func fakeTCPConnToDNS(answer []byte) (net.Conn, error) {
	host := "127.0.0.1"
	if isTails() {
		host = getLocalIP()
	}

	fakeResolver, _ := net.Listen("tcp", fmt.Sprintf("%s:0", host))
	go func() {
		conn, _ := fakeResolver.Accept()

		var dest [46]byte
		_, _ = conn.Read(dest[:])
		_, _ = conn.Write(answer)

		_ = conn.Close()
	}()

	return net.Dial("tcp", fakeResolver.Addr().String())
}

func (s *DNSXMPPSuite) Test_resolve_resolvesCorrectly(c *C) {
	dec, _ := hex.DecodeString("00511eea818000010001000000000c5f786d70702d636c69656e74045f746370076f6c6162696e690273650000210001c00c0021000100000258001700000005146604786d7070076f6c6162696e6902736500")

	p := &mockProxy{}
	p.Expects(func(network, addr string) (net.Conn, error) {
		c.Check(network, Equals, "tcp")
		c.Check(addr, Equals, "208.67.222.222:53")

		return fakeTCPConnToDNS(dec)
	})
	p.Expects(func(network, addr string) (net.Conn, error) {
		c.Check(network, Equals, "tcp")
		c.Check(addr, Equals, "208.67.222.222:53")

		return fakeTCPConnToDNS(dec)
	})

	hostport, err := resolveSRVWithProxy(p, "olabini.se")
	c.Assert(err, IsNil)
	c.Assert(hostport[0], DeepEquals, &connectEntry{host: "xmpp.olabini.se", port: 5222, priority: 0, weight: 5, tls: true})
	c.Check(p, MatchesExpectations)
}

// WARNING: this test requires a real live connection to the Internet. Not so good...
func (s *DNSXMPPSuite) Test_resolve_handlesErrors(c *C) {
	_, err := resolve("doesntexist.olabini.se")

	//It only happens when using golang resolver
	//resolveSRVWithProxy will not return an error
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Matches, "lookup _xmpps-client._tcp.doesntexist.olabini.se.*?")
}
