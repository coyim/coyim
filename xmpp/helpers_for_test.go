package xmpp

import (
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"gopkg.in/check.v1"
)

type mockConn struct {
	calledClose int
	net.TCPConn
}

func (c *mockConn) Close() error {
	c.calledClose++
	return nil
}

type mockConnIOReaderWriter struct {
	read      []byte
	readIndex int
	write     []byte
	errCount  int
	err       error

	calledClose int
}

func (iom *mockConnIOReaderWriter) Read(p []byte) (n int, err error) {
	if iom.readIndex >= len(iom.read) {
		return 0, io.EOF
	}
	i := copy(p, iom.read[iom.readIndex:])
	iom.readIndex += i
	var e error
	if iom.errCount == 0 {
		e = iom.err
	}
	iom.errCount--
	return i, e
}

func (iom *mockConnIOReaderWriter) Write(p []byte) (n int, err error) {
	iom.write = append(iom.write, p...)
	var e error
	if iom.errCount == 0 {
		e = iom.err
	}
	iom.errCount--
	return len(p), e
}

func (iom *mockConnIOReaderWriter) Close() error {
	iom.calledClose++
	return nil
}

type mockMultiConnIOReaderWriter struct {
	read      [][]byte
	readIndex int
	write     []byte
}

func (iom *mockMultiConnIOReaderWriter) Read(p []byte) (n int, err error) {
	if iom.readIndex >= len(iom.read) {
		return 0, io.EOF
	}
	i := copy(p, iom.read[iom.readIndex])
	iom.readIndex++
	return i, nil
}

func (iom *mockMultiConnIOReaderWriter) Write(p []byte) (n int, err error) {
	iom.write = append(iom.write, p...)
	return len(p), nil
}

type fullMockedConn struct {
	rw io.ReadWriter
}

func (c *fullMockedConn) Read(b []byte) (n int, err error) {
	return c.rw.Read(b)
}

func (c *fullMockedConn) Write(b []byte) (n int, err error) {
	return c.rw.Write(b)
}

func (c *fullMockedConn) Close() error {
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

type fixedRandReader struct {
	data []string
	at   int
}

func fixedRand(data []string) io.Reader {
	return &fixedRandReader{data, 0}
}

func bytesFromHex(s string) []byte {
	val, _ := hex.DecodeString(s)
	return val
}

func byteStringFromHex(s string) string {
	val, _ := hex.DecodeString(s)
	return string(val)
}

func (frr *fixedRandReader) Read(p []byte) (n int, err error) {
	if frr.at < len(frr.data) {
		plainBytes := bytesFromHex(frr.data[frr.at])
		frr.at++
		n = copy(p, plainBytes)
		return
	}
	return 0, io.EOF
}

func createTeeConn(c net.Conn, w io.Writer) net.Conn {
	return &teeConn{c, w}
}

type teeConn struct {
	c net.Conn
	w io.Writer
}

func (c *teeConn) Read(b []byte) (n int, err error) {
	n, err = c.c.Read(b)
	if n > 0 {
		fmt.Fprintf(c.w, "READ: %x\n", b[:n])
	}
	return
}

func (c *teeConn) Write(b []byte) (n int, err error) {
	n, err = c.c.Write(b)
	if n > 0 {
		fmt.Fprintf(c.w, "WRITE: %x\n", b[:n])
	}
	return n, err
}

func (c *teeConn) Close() error {
	return c.c.Close()
}

func (c *teeConn) LocalAddr() net.Addr {
	return c.c.LocalAddr()
}

func (c *teeConn) RemoteAddr() net.Addr {
	return c.c.RemoteAddr()
}

func (c *teeConn) SetDeadline(t time.Time) error {
	return c.c.SetDeadline(t)
}

func (c *teeConn) SetReadDeadline(t time.Time) error {
	return c.c.SetReadDeadline(t)
}

func (c *teeConn) SetWriteDeadline(t time.Time) error {
	return c.c.SetWriteDeadline(t)
}

type dialCall func(string, string) (c net.Conn, e error)
type dialCallExp struct {
	f      dialCall
	called bool
}

type mockProxy struct {
	called int
	calls  []dialCallExp
	sync.Mutex
}

func (p *mockProxy) Dial(network, addr string) (net.Conn, error) {
	if len(p.calls)-1 < p.called {
		return nil, fmt.Errorf("unexpected call to Dial: %s, %s", network, addr)
	}

	p.Lock()
	defer p.Unlock()

	fn := p.calls[p.called]
	p.called = p.called + 1

	fn.called = true
	return fn.f(network, addr)
}

func (p *mockProxy) Expects(f dialCall) {
	p.Lock()
	defer p.Unlock()

	if p.calls == nil {
		p.calls = []dialCallExp{}
	}

	p.calls = append(p.calls, dialCallExp{f: f})
}

var MatchesExpectations check.Checker = &allExpectations{
	&check.CheckerInfo{Name: "IsNil", Params: []string{"value"}},
}

type allExpectations struct {
	*check.CheckerInfo
}

func (checker *allExpectations) Check(params []interface{}, names []string) (result bool, error string) {
	p := params[0].(*mockProxy)

	if p.called != len(p.calls) {
		return false, fmt.Sprintf("expected: %d calls, got: %d", len(p.calls), p.called)
	}

	return true, ""
}
