package xmpp

import (
	"encoding/hex"
	"fmt"
	"net"
	"sort"

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

func (s *DNSXMPPSuite) Test_intoConnectEntry_returnsNilOnFailure(c *C) {
	res := intoConnectEntry("bla")
	c.Assert(res, IsNil)
}

func (s *DNSXMPPSuite) Test_byPriorityWeight_sortsConnectEntries(c *C) {
	res := []*connectEntry{
		&connectEntry{host: "a", priority: 1, weight: 1},
		&connectEntry{host: "b", priority: 1, weight: 42},
		&connectEntry{host: "c", priority: 10, weight: 1},
		&connectEntry{host: "d", priority: 1, weight: 1},
		&connectEntry{host: "e", priority: 1, weight: 3},
		&connectEntry{host: "f", priority: 1, weight: 1},
		&connectEntry{host: "g", priority: 6, weight: 1},
		&connectEntry{host: "h", priority: 1, weight: 1},
	}

	sort.Sort(byPriorityWeight(res))
	c.Assert(res[0].host, Equals, "b")
	c.Assert(res[1].host, Equals, "e")
	c.Assert(res[2].host, Equals, "a")
	c.Assert(res[3].host, Equals, "d")
	c.Assert(res[4].host, Equals, "f")
	c.Assert(res[5].host, Equals, "h")
	c.Assert(res[6].host, Equals, "g")
	c.Assert(res[7].host, Equals, "c")
}

func (s *DNSXMPPSuite) Test_resolveWithCustom(c *C) {
	resv := func(part, tp, domain string) (string, []*net.SRV, error) {
		if part == "xmpps-client" {
			return "", nil, nil
		}

		return "", []*net.SRV{
			&net.SRV{Target: "."},
		}, nil
	}

	res, e := resolveWithCustom("foobar.com", resv)
	c.Assert(res, HasLen, 0)
	c.Assert(e, ErrorMatches, "service not available")
}
