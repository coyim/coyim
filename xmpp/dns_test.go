package xmpp

import (
	"encoding/hex"
	"net"

	. "gopkg.in/check.v1"
)

type DNSXMPPSuite struct{}

var _ = Suite(&DNSXMPPSuite{})

func fakeTCPConnToDNS(answer []byte) (net.Conn, error) {
	fakeResolver, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		conn, _ := fakeResolver.Accept()

		var dest [46]byte
		conn.Read(dest[:])
		conn.Write(answer)

		conn.Close()
	}()

	return net.Dial("tcp", fakeResolver.Addr().String())
}

func (s *DNSXMPPSuite) Test_Resolve_resolvesCorrectly(c *C) {
	dec, _ := hex.DecodeString("00511eea818000010001000000000c5f786d70702d636c69656e74045f746370076f6c6162696e690273650000210001c00c0021000100000258001700000005146604786d7070076f6c6162696e6902736500")

	p := &mockProxy{}
	p.Expects(func(network, addr string) (net.Conn, error) {
		c.Check(network, Equals, "tcp")
		c.Check(addr, Equals, "208.67.222.222:53")

		return fakeTCPConnToDNS(dec)
	})

	hostport, err := ResolveSRVWithProxy(p, "olabini.se")
	c.Assert(err, IsNil)
	c.Assert(hostport[0], Equals, "xmpp.olabini.se:5222")
	c.Check(p, MatchesExpectations)
}

// WARNING: this test requires a real live connection to the Internet. Not so good...
func (s *DNSXMPPSuite) Test_Resolve_handlesErrors(c *C) {
	_, err := Resolve("doesntexist.olabini.se")

	//It only happens when using golang resolver
	//ResolveSRVWithProxy will not return an error
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Matches, "lookup _xmpp-client._tcp.doesntexist.olabini.se.*?")
}
