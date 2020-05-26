package xmpp

import (
	gotls "crypto/tls"
	"net"
	"time"

	"github.com/coyim/coyim/tls"
)

type tlsMock1 struct {
	returnFromHandshake error
	returnFromConnState gotls.ConnectionState
	returnFromRead1     int
	returnFromRead2     error
}

func (t *tlsMock1) Handshake() error {
	return t.returnFromHandshake
}

func (t *tlsMock1) ConnectionState() gotls.ConnectionState {
	return t.returnFromConnState
}

func (t *tlsMock1) Read(b []byte) (n int, err error) {
	return t.returnFromRead1, t.returnFromRead2
}

func (t *tlsMock1) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (t *tlsMock1) Close() error {
	return nil
}

func (t *tlsMock1) LocalAddr() net.Addr {
	return nil
}

func (t *tlsMock1) RemoteAddr() net.Addr {
	return nil
}

func (t *tlsMock1) SetDeadline(t2 time.Time) error {
	return nil
}

func (t *tlsMock1) SetReadDeadline(t2 time.Time) error {
	return nil
}

func (t *tlsMock1) SetWriteDeadline(t2 time.Time) error {
	return nil
}

func fixedTLSFactory(t tls.Conn) tls.Factory {
	return func(net.Conn, *gotls.Config) tls.Conn {
		return t
	}
}

type mockTLSVerifier struct {
	verifyCalled int
	state        gotls.ConnectionState
	conf         *gotls.Config
	originDomain string

	toReturn error
}

func (v *mockTLSVerifier) Verify(state gotls.ConnectionState, conf *gotls.Config, originDomain string) error {
	v.verifyCalled++
	v.state = state
	v.conf = conf
	v.originDomain = originDomain

	return v.toReturn
}
