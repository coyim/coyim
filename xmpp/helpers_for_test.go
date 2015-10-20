package xmpp

import (
	"encoding/hex"
	"fmt"
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
