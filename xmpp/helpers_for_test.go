package xmpp

import (
	"io"
	"net"
	"time"
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
}

func (in *mockConnIOReaderWriter) Read(p []byte) (n int, err error) {
	if in.readIndex >= len(in.read) {
		return 0, io.EOF
	}
	i := copy(p, in.read[in.readIndex:])
	in.readIndex += i
	var e error
	if in.errCount == 0 {
		e = in.err
	}
	in.errCount--
	return i, e
}

func (out *mockConnIOReaderWriter) Write(p []byte) (n int, err error) {
	out.write = append(out.write, p...)
	var e error
	if out.errCount == 0 {
		e = out.err
	}
	out.errCount--
	return len(p), e
}

type fullMockedConn struct {
	rw *mockConnIOReaderWriter
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
