package socks5

import (
	"errors"
	"net"
	"time"

	"golang.org/x/net/proxy"
	. "gopkg.in/check.v1"
)

type Socks5Suite struct{}

var _ = Suite(&Socks5Suite{})

func (s *Socks5Suite) Test_XMPP_createsWithoutAuthorizationData(c *C) {
	v, e := XMPP("tcp", "foo:84", nil, nil)
	c.Assert(e, IsNil)

	vv := v.(*socks5)
	c.Assert(vv.user, Equals, "")
	c.Assert(vv.password, Equals, "")
	c.Assert(vv.network, Equals, "tcp")
	c.Assert(vv.addr, Equals, "foo:84")
	c.Assert(vv.forward, IsNil)
}

func (s *Socks5Suite) Test_XMPP_createsWithAuthorizationData(c *C) {
	v, e := XMPP("tcp", "foo:84", &proxy.Auth{User: "hello", Password: "goodbye"}, nil)
	c.Assert(e, IsNil)

	vv := v.(*socks5)
	c.Assert(vv.user, Equals, "hello")
	c.Assert(vv.password, Equals, "goodbye")
	c.Assert(vv.network, Equals, "tcp")
	c.Assert(vv.addr, Equals, "foo:84")
	c.Assert(vv.forward, IsNil)
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

func (s *Socks5Suite) Test_socks5_Dial_failsOnUnknownNetwork(c *C) {
	v := &socks5{}
	_, e := v.Dial("something", "")
	c.Assert(e, ErrorMatches, ".*no support for SOCKS5 proxy connections of type .*")
}

func (s *Socks5Suite) Test_socks5_Dial_failsIfDialingRealConnectionFails(c *C) {
	v := &socks5{
		forward: &mockDialer{returnErr: errors.New("something is bad")},
	}
	_, e := v.Dial("tcp6", "")
	c.Assert(e, ErrorMatches, "something is bad")
}

func (s *Socks5Suite) Test_socks5_Dial_failsIfConnectFails(c *C) {
	conn := &fullMockedConn{}
	v := &socks5{
		forward: &mockDialer{
			returnConn: conn,
		},
	}
	_, e := v.Dial("tcp4", "blarg%$%$56")
	c.Assert(e, ErrorMatches, ".*missing port in address.*")
	c.Assert(conn.closeCalled, Equals, 1)
}

func (s *Socks5Suite) Test_socks5_Dial_succeedsAndReturnsConnection(c *C) {
	conn := &fullMockedConn{
		readRetBuf: [][]byte{
			[]byte{5, 0},
			[]byte{0, 0, 0, socks5Domain},
			[]byte{5},
			[]byte{0, 0, 0, 0, 0},
			[]byte{2, 0},
		},
		readRetErr: []error{
			nil,
			nil,
			nil,
			nil,
			nil,
		},
	}
	v := &socks5{
		forward: &mockDialer{
			returnConn: conn,
		},
	}
	cc, e := v.Dial("tcp", "localhost:42")
	c.Assert(e, IsNil)
	c.Assert(cc, Equals, conn)
}

func (s *Socks5Suite) Test_socks5_connect_failsOnUnparseablePortNumber(c *C) {
	v := &socks5{}
	e := v.connect(nil, "localhost:something")
	c.Assert(e, ErrorMatches, ".*proxy: failed to parse port number.*")
}

func (s *Socks5Suite) Test_socks5_connect_failsOnTooLargePortNumber(c *C) {
	v := &socks5{}
	e := v.connect(nil, "localhost:123243243")
	c.Assert(e, ErrorMatches, ".*proxy: port number out of range:.*")
}

func (s *Socks5Suite) Test_socks5_connect_failsOnWritingGreeting(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0},
		writeRetErr: []error{errors.New("something wrongy")},
	}
	v := &socks5{}
	e := v.connect(conn, "localhost:123")

	c.Assert(e, ErrorMatches, "proxy: failed to write greeting to SOCKS5 proxy at .*")
	c.Assert(conn.writeArgs[0], DeepEquals, []byte{0x05, 0x01, 0x00})
}

func (s *Socks5Suite) Test_socks5_connect_failsOnReadingGreeting(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0},
		writeRetErr: []error{nil},

		readRetBuf: [][]byte{
			[]byte{0},
		},
		readRetErr: []error{
			errors.New("reading is bad"),
		},
	}
	v := &socks5{}
	e := v.connect(conn, "localhost:123")

	c.Assert(e, ErrorMatches, "proxy: failed to read greeting from SOCKS5 proxy at .*")
}

func (s *Socks5Suite) Test_socks5_connect_failsOnUnexpectedVersion(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0},
		writeRetErr: []error{nil},

		readRetBuf: [][]byte{
			[]byte{6, 0},
		},
		readRetErr: []error{
			nil,
		},
	}
	v := &socks5{}
	e := v.connect(conn, "localhost:123")

	c.Assert(e, ErrorMatches, ".*has unexpected version.*")
}

func (s *Socks5Suite) Test_socks5_connect_failsOnAuthenticationRequirement(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0},
		writeRetErr: []error{nil},

		readRetBuf: [][]byte{
			[]byte{5, 0xFF},
		},
		readRetErr: []error{
			nil,
		},
	}
	v := &socks5{}
	e := v.connect(conn, "localhost:123")

	c.Assert(e, ErrorMatches, ".*requires authentication.*")
}

func (s *Socks5Suite) Test_socks5_connect_failsOnTooLongHostname(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0},
		writeRetErr: []error{nil},

		readRetBuf: [][]byte{
			[]byte{5, 0},
		},
		readRetErr: []error{
			nil,
		},
	}
	v := &socks5{}
	e := v.connect(conn, "localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost.localhost:123")

	c.Assert(e, ErrorMatches, ".*destination hostname too long.*")
}

func (s *Socks5Suite) Test_socks5_connect_failsOnWritingConnectionRequest(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0, 0},
		writeRetErr: []error{nil, errors.New("don't want to")},

		readRetBuf: [][]byte{
			[]byte{5, 0},
		},
		readRetErr: []error{
			nil,
		},
	}
	v := &socks5{}
	e := v.connect(conn, "localhost:123")

	c.Assert(e, ErrorMatches, ".*failed to write connect request to SOCKS5 proxy at.*")
	c.Assert(conn.writeArgs[1], DeepEquals, []byte{socks5Version, socks5Connect, 0x0, socks5Domain, 0x9, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, 0x0, 0x7b})
}

func (s *Socks5Suite) Test_socks5_connect_failsOnReadingConnectionRequest(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0, 0},
		writeRetErr: []error{nil, nil},

		readRetBuf: [][]byte{
			[]byte{5, 0},
			[]byte{0},
		},
		readRetErr: []error{
			nil,
			errors.New("really don't want to connect"),
		},
	}
	v := &socks5{}
	e := v.connect(conn, "localhost:123")

	c.Assert(e, ErrorMatches, ".*failed to read connect reply from SOCKS5 proxy.*")
}

func (s *Socks5Suite) Test_socks5_connect_reportsConnectionRequestFailure(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0, 0},
		writeRetErr: []error{nil, nil},

		readRetBuf: [][]byte{
			[]byte{5, 0},
			[]byte{0, 3, 0, socks5Domain},
		},
		readRetErr: []error{
			nil,
			nil,
		},
	}
	v := &socks5{}
	e := v.connect(conn, "localhost:123")

	c.Assert(e, ErrorMatches, ".*failed to connect: network unreachable.*")
}

func (s *Socks5Suite) Test_socks5_connect_reportsUnknownFailure(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0, 0},
		writeRetErr: []error{nil, nil},

		readRetBuf: [][]byte{
			[]byte{5, 0},
			[]byte{0, 30, 0, socks5Domain},
		},
		readRetErr: []error{
			nil,
			nil,
		},
	}
	v := &socks5{}
	e := v.connect(conn, "localhost:123")

	c.Assert(e, ErrorMatches, ".*failed to connect: unknown error.*")
}

func (s *Socks5Suite) Test_socks5_connect_failsOnUnknownDomainType(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0, 0},
		writeRetErr: []error{nil, nil},

		readRetBuf: [][]byte{
			[]byte{5, 0},
			[]byte{0, 0, 0, 0xFF},
		},
		readRetErr: []error{
			nil,
			nil,
		},
	}
	v := &socks5{}
	e := v.connect(conn, "localhost:123")

	c.Assert(e, ErrorMatches, "proxy: got unknown address type.*")
}

func (s *Socks5Suite) Test_socks5_connect_failsIfCantReadDomainLength(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0, 0},
		writeRetErr: []error{nil, nil},

		readRetBuf: [][]byte{
			[]byte{5, 0},
			[]byte{0, 0, 0, socks5Domain},
			[]byte{},
		},
		readRetErr: []error{
			nil,
			nil,
			errors.New("no read for you"),
		},
	}
	v := &socks5{}
	e := v.connect(conn, "localhost:123")

	c.Assert(e, ErrorMatches, "proxy: failed to read domain length from SOCKS5.*")
}

func (s *Socks5Suite) Test_socks5_connect_failsIfCantReadTheDomainToDiscard(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0, 0},
		writeRetErr: []error{nil, nil},

		readRetBuf: [][]byte{
			[]byte{5, 0},
			[]byte{0, 0, 0, socks5Domain},
			[]byte{0xFF},
			[]byte{0},
		},
		readRetErr: []error{
			nil,
			nil,
			nil,
			errors.New("no read for you"),
		},
	}
	v := &socks5{}
	e := v.connect(conn, "localhost:123")

	c.Assert(e, ErrorMatches, "proxy: failed to read address from SOCKS5.*")
}

func (s *Socks5Suite) Test_socks5_connect_failsIfCantReadThePortToDiscard(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0, 0},
		writeRetErr: []error{nil, nil},

		readRetBuf: [][]byte{
			[]byte{5, 0},
			[]byte{0, 0, 0, socks5Domain},
			[]byte{0x5},
			[]byte{0, 0, 0, 0, 0},
			[]byte{},
		},
		readRetErr: []error{
			nil,
			nil,
			nil,
			nil,
			errors.New("no read for you"),
		},
	}
	v := &socks5{}
	e := v.connect(conn, "localhost:123")

	c.Assert(e, ErrorMatches, "proxy: failed to read port from SOCKS5.*")
}

func (s *Socks5Suite) Test_socks5_connect_succeedsWithUsernameAndPassword(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0, 0, 0},
		writeRetErr: []error{nil, nil, nil},

		readRetBuf: [][]byte{
			[]byte{5, socks5AuthPassword},
			[]byte{0, 0},
			[]byte{0, 0, 0, socks5Domain},
			[]byte{5},
			[]byte{0, 0, 0, 0, 0},
			[]byte{2, 0},
		},
		readRetErr: []error{
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		},
	}
	v := &socks5{
		user:     "foo",
		password: "bar",
	}
	e := v.connect(conn, "localhost:123")
	c.Assert(e, IsNil)
	c.Assert(conn.writeArgs[0], DeepEquals, []byte{0x05, 0x02, 0x00, 0x2})
	c.Assert(conn.writeArgs[1], DeepEquals, []byte{0x1, 0x3, 0x66, 0x6f, 0x6f, 0x3, 0x62, 0x61, 0x72})
}

func (s *Socks5Suite) Test_socks5_connect_failsOnWritingAuthenticationRequest(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0, 0, 0},
		writeRetErr: []error{nil, errors.New("auth fail"), nil},

		readRetBuf: [][]byte{
			[]byte{5, socks5AuthPassword},
			[]byte{0, 0},
			[]byte{0, 0, 0, socks5Domain},
			[]byte{5},
			[]byte{0, 0, 0, 0, 0},
			[]byte{2, 0},
		},
		readRetErr: []error{
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		},
	}
	v := &socks5{
		user:     "foo",
		password: "bar",
	}
	e := v.connect(conn, "localhost:123")
	c.Assert(e, ErrorMatches, ".*failed to write authentication request to SOCKS5 proxy.*")
	c.Assert(conn.writeArgs[1], DeepEquals, []byte{0x1, 0x3, 0x66, 0x6f, 0x6f, 0x3, 0x62, 0x61, 0x72})
}

func (s *Socks5Suite) Test_socks5_connect_failsOnReadingAuthenticationResponse(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0, 0, 0},
		writeRetErr: []error{nil, nil, nil},

		readRetBuf: [][]byte{
			[]byte{5, socks5AuthPassword},
			[]byte{0},
			[]byte{0, 0, 0, socks5Domain},
			[]byte{5},
			[]byte{0, 0, 0, 0, 0},
			[]byte{2, 0},
		},
		readRetErr: []error{
			nil,
			errors.New("auth fail"),
			nil,
			nil,
			nil,
			nil,
		},
	}
	v := &socks5{
		user:     "foo",
		password: "bar",
	}
	e := v.connect(conn, "localhost:123")
	c.Assert(e, ErrorMatches, ".*failed to read authentication reply from SOCKS5 proxy.*")
}

func (s *Socks5Suite) Test_socks5_connect_handlesAuthFailureCorrectly(c *C) {
	conn := &fullMockedConn{
		writeRetInt: []int{0, 0, 0},
		writeRetErr: []error{nil, nil, nil},

		readRetBuf: [][]byte{
			[]byte{5, socks5AuthPassword},
			[]byte{0, 1},
			[]byte{0, 0, 0, socks5Domain},
			[]byte{5},
			[]byte{0, 0, 0, 0, 0},
			[]byte{2, 0},
		},
		readRetErr: []error{
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		},
	}
	v := &socks5{
		user:     "foo",
		password: "bar",
	}
	e := v.connect(conn, "localhost:123")
	c.Assert(e, ErrorMatches, ".*rejected username/password.*")
}
