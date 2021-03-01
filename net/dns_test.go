package net

import (
	"errors"
	"net"
	"sort"
	"time"

	"github.com/miekg/dns"

	. "gopkg.in/check.v1"
)

type DNSXmppSuite struct{}

var _ = Suite(&DNSXmppSuite{})

func (s *DNSXmppSuite) Test_createCName_createsAValidCnameForAService(c *C) {
	ret := createCName("foo", "bar", "bax.com")
	c.Assert(ret, Equals, "_foo._bar.bax.com.")
}

func (s *DNSXmppSuite) Test_convertAnswerToSRV_returnsNilForNonSRVEntry(c *C) {
	cn := new(dns.CNAME)
	res := convertAnswerToSRV(cn)
	c.Assert(res, IsNil)
}

func (s *DNSXmppSuite) Test_convertAnswerToSRV_returnsAValidNetSRV(c *C) {
	srv := new(dns.SRV)
	srv.Target = "foo.com"
	srv.Port = 123
	srv.Priority = 5
	srv.Weight = 42
	res := convertAnswerToSRV(srv)
	c.Assert(res, Not(IsNil))
	c.Assert(res.Target, Equals, "foo.com")
	c.Assert(res.Port, Equals, uint16(123))
	c.Assert(res.Priority, Equals, uint16(5))
	c.Assert(res.Weight, Equals, uint16(42))
}

func (s *DNSXmppSuite) Test_convertAnswersToSRV_convertsAnswers(c *C) {
	cn := new(dns.CNAME)
	srv := new(dns.SRV)
	srv.Target = "foo2.com"

	in := make([]dns.RR, 2)
	in[0] = cn
	in[1] = srv
	res := convertAnswersToSRV(in)

	c.Assert(res, HasLen, 1)
	c.Assert(res[0].Target, Equals, "foo2.com")
}

func (s *DNSXmppSuite) Test_msgSRV_createsMessage(c *C) {
	res := msgSRV("foo.com")
	c.Assert(res.Question[0].Name, Equals, "foo.com")
	c.Assert(res.Question[0].Qtype, Equals, dns.TypeSRV)
}

func (s *DNSXmppSuite) Test_convertAnswersToSRV_sortsByPriority(c *C) {
	srv1 := &dns.SRV{
		Target:   "foo1.com",
		Priority: 5,
		Weight:   1,
	}
	srv2 := &dns.SRV{
		Target:   "foo2.com",
		Priority: 3,
		Weight:   1,
	}
	srv3 := &dns.SRV{
		Target:   "foo3.com",
		Priority: 6,
		Weight:   1,
	}
	srv4 := &dns.SRV{
		Target:   "foo4.com",
		Priority: 1,
		Weight:   1,
	}

	in := []dns.RR{
		srv1,
		srv2,
		srv3,
		srv4,
	}
	res := convertAnswersToSRV(in)
	c.Assert(res[0].Target, Equals, "foo4.com")
	c.Assert(res[1].Target, Equals, "foo2.com")
	c.Assert(res[2].Target, Equals, "foo1.com")
	c.Assert(res[3].Target, Equals, "foo3.com")
}

func (s *DNSXmppSuite) Test_convertAnswersToSRV_sortsByWeightIfPriotityIsTheSame(c *C) {
	srv1 := &dns.SRV{
		Target:   "foo1.com",
		Priority: 1,
		Weight:   5,
	}
	srv2 := &dns.SRV{
		Target:   "foo2.com",
		Priority: 1,
		Weight:   3,
	}
	srv3 := &dns.SRV{
		Target:   "foo3.com",
		Priority: 1,
		Weight:   6,
	}
	srv4 := &dns.SRV{
		Target:   "foo4.com",
		Priority: 1,
		Weight:   1,
	}

	in := []dns.RR{
		srv1,
		srv2,
		srv3,
		srv4,
	}
	res := convertAnswersToSRV(in)
	c.Assert(res[0].Target, Equals, "foo3.com")
	c.Assert(res[1].Target, Equals, "foo1.com")
	c.Assert(res[2].Target, Equals, "foo2.com")
	c.Assert(res[3].Target, Equals, "foo4.com")
}

func (s *DNSXmppSuite) Test_byPriorityWeight_sortsWithNils(c *C) {
	v := []*net.SRV{
		&net.SRV{Priority: 3, Weight: 1},
		nil,
		&net.SRV{Priority: 5, Weight: 1},
		&net.SRV{Priority: 1, Weight: 1},
		nil,
		&net.SRV{Priority: 7, Weight: 1},
		nil,
		&net.SRV{Priority: 8, Weight: 1},
	}
	sort.Sort(byPriorityWeight(v))
	c.Assert(v[0], IsNil)
	c.Assert(v[1], IsNil)
	c.Assert(v[2], IsNil)
	c.Assert(v[3].Priority, Equals, uint16(1))
	c.Assert(v[4].Priority, Equals, uint16(3))
	c.Assert(v[5].Priority, Equals, uint16(5))
	c.Assert(v[6].Priority, Equals, uint16(7))
	c.Assert(v[7].Priority, Equals, uint16(8))
}

type mockDialer struct {
	argNetwork string
	argAddress string
	returnConn net.Conn
	returnErr  error
}

func (m *mockDialer) Dial(network, address string) (net.Conn, error) {
	m.argNetwork = network
	m.argAddress = address
	return m.returnConn, m.returnErr
}

type fullMockedConn struct {
	closeCalled int

	readRetBuf [][]byte
	readRetErr []error

	writeArgs   [][]byte
	writeRetInt []int
	writeRetErr []error
}

func (c *fullMockedConn) Read(b []byte) (n int, err error) {
	if len(c.readRetBuf) > 0 {
		buf := c.readRetBuf[0]
		c.readRetBuf = c.readRetBuf[1:]
		copy(b, buf)
		e := c.readRetErr[0]
		c.readRetErr = c.readRetErr[1:]
		return len(buf), e
	}
	return 0, nil
}

func (c *fullMockedConn) Write(b []byte) (n int, err error) {
	newBuf := make([]byte, len(b))
	copy(newBuf, b)
	c.writeArgs = append(c.writeArgs, newBuf)
	if len(c.writeRetInt) > 0 {
		ret := c.writeRetInt[0]
		c.writeRetInt = c.writeRetInt[1:]
		e := c.writeRetErr[0]
		c.writeRetErr = c.writeRetErr[1:]
		return ret, e
	}
	return 0, nil
}

func (c *fullMockedConn) Close() error {
	c.closeCalled++
	return nil
}

func (c *fullMockedConn) LocalAddr() net.Addr {
	return nil
}

func (c *fullMockedConn) RemoteAddr() net.Addr {
	return nil
}

func (c *fullMockedConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *fullMockedConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *fullMockedConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (s *DNSXmppSuite) Test_LookupSRV_works(c *C) {
	conn := &fullMockedConn{
		readRetBuf: [][]byte{
			[]byte{0x0, 0x78},
			[]byte{0x5e, 0x92, 0x81, 0x80, 0x0, 0x1, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0xc, 0x5f, 0x78, 0x6d, 0x70, 0x70, 0x2d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4, 0x5f, 0x74, 0x63, 0x70, 0x4, 0x72, 0x65, 0x61, 0x70, 0x2, 0x65, 0x63, 0x0, 0x0, 0x21, 0x0, 0x1, 0xc0, 0xc, 0x0, 0x21, 0x0, 0x1, 0x0, 0x0, 0x2, 0x58, 0x0, 0x17, 0x0, 0x0, 0x0, 0x5, 0x14, 0x66, 0x4, 0x78, 0x6d, 0x70, 0x70, 0x7, 0x6f, 0x6c, 0x61, 0x62, 0x69, 0x6e, 0x69, 0x2, 0x73, 0x65, 0x0, 0xc0, 0xc, 0x0, 0x21, 0x0, 0x1, 0x0, 0x0, 0x2, 0x58, 0x0, 0x1e, 0x0, 0x0, 0x0, 0xa, 0x14, 0x66, 0x10, 0x65, 0x70, 0x36, 0x76, 0x35, 0x61, 0x77, 0x76, 0x79, 0x76, 0x35, 0x76, 0x74, 0x73, 0x33, 0x75, 0x5, 0x6f, 0x6e, 0x69, 0x6f, 0x6e, 0x0},
		},
		readRetErr: []error{
			nil,
			nil,
		},
	}
	dialer := &mockDialer{
		returnConn: conn,
	}

	cname, addrs, e := LookupSRV(dialer, "xmpp-client", "tcp", "reap.ec")
	c.Assert(e, IsNil)
	c.Assert(cname, Equals, "_xmpp-client._tcp.reap.ec.")
	c.Assert(addrs, HasLen, 2)
	c.Assert(conn.writeArgs, HasLen, 2)
	c.Assert(conn.writeArgs[0], DeepEquals, []byte{0x0, 0x2b})
}

func (s *DNSXmppSuite) Test_timingOutLookup_timesOut(c *C) {
	_, _, e := timingOutLookup(func() (string, []*net.SRV, error) {
		time.Sleep(time.Duration(10) * time.Second)
		return "", nil, nil
	}, time.Duration(1)*time.Millisecond)

	c.Assert(e, ErrorMatches, "i/o timeout")
}

func (s *DNSXmppSuite) Test_LookupSRV_failsOnDial(c *C) {
	dialer := &mockDialer{
		returnConn: nil,
		returnErr:  errors.New("stuff"),
	}

	_, addrs, e := LookupSRV(dialer, "xmpp-client", "tcp", "reap.ec")
	c.Assert(e, ErrorMatches, "stuff")
	c.Assert(addrs, IsNil)
}

func (s *DNSXmppSuite) Test_LookupSRV_failsOnWriteOfExchange(c *C) {
	conn := &fullMockedConn{
		writeRetErr: []error{
			errors.New("haha"),
		},
		writeRetInt: []int{
			0,
		},
	}
	dialer := &mockDialer{
		returnConn: conn,
	}

	_, _, e := LookupSRV(dialer, "xmpp-client", "tcp", "reap.ec")
	c.Assert(e, ErrorMatches, "haha")
}

func (s *DNSXmppSuite) Test_LookupSRV_failsOnReadOfExchange(c *C) {
	conn := &fullMockedConn{
		readRetBuf: [][]byte{
			[]byte{},
			[]byte{},
		},
		readRetErr: []error{
			errors.New("oh no"),
			nil,
		},
	}
	dialer := &mockDialer{
		returnConn: conn,
	}

	_, _, e := LookupSRV(dialer, "xmpp-client", "tcp", "reap.ec")
	c.Assert(e, ErrorMatches, "oh no")
}

func (s *DNSXmppSuite) Test_LookupSRV_failsIfResultIsNotSuccess(c *C) {
	conn := &fullMockedConn{
		readRetBuf: [][]byte{
			[]byte{0x0, 0x78},
			[]byte{0x5e, 0x92, 0x81, 0x85, 0x0, 0x1, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0xc, 0x5f, 0x78, 0x6d, 0x70, 0x70, 0x2d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4, 0x5f, 0x74, 0x63, 0x70, 0x4, 0x72, 0x65, 0x61, 0x70, 0x2, 0x65, 0x63, 0x0, 0x0, 0x21, 0x0, 0x1, 0xc0, 0xc, 0x0, 0x21, 0x0, 0x1, 0x0, 0x0, 0x2, 0x58, 0x0, 0x17, 0x0, 0x0, 0x0, 0x5, 0x14, 0x66, 0x4, 0x78, 0x6d, 0x70, 0x70, 0x7, 0x6f, 0x6c, 0x61, 0x62, 0x69, 0x6e, 0x69, 0x2, 0x73, 0x65, 0x0, 0xc0, 0xc, 0x0, 0x21, 0x0, 0x1, 0x0, 0x0, 0x2, 0x58, 0x0, 0x1e, 0x0, 0x0, 0x0, 0xa, 0x14, 0x66, 0x10, 0x65, 0x70, 0x36, 0x76, 0x35, 0x61, 0x77, 0x76, 0x79, 0x76, 0x35, 0x76, 0x74, 0x73, 0x33, 0x75, 0x5, 0x6f, 0x6e, 0x69, 0x6f, 0x6e, 0x0},
		},
		readRetErr: []error{
			nil,
			nil,
		},
	}
	dialer := &mockDialer{
		returnConn: conn,
	}

	_, _, e := LookupSRV(dialer, "xmpp-client", "tcp", "reap.ec")
	c.Assert(e, ErrorMatches, "got return: 5")
}
