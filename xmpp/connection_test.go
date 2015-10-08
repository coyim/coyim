package xmpp

import (
	"encoding/xml"
	"io"

	. "gopkg.in/check.v1"
)

type ConnectionXmppSuite struct{}

var _ = Suite(&ConnectionXmppSuite{})

func (s *ConnectionXmppSuite) Test_Next_returnsErrorIfOneIsEncountered(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:foo xmlns:stream='http://etherx.jabber.org/streams' to='hello'></stream:foo>")}
	conn := Conn{
		in: xml.NewDecoder(mockIn),
	}

	_, err := conn.Next()
	c.Assert(err.Error(), Equals, "unexpected XMPP message http://etherx.jabber.org/streams <foo/>")
}

func (s *ConnectionXmppSuite) Test_Next_returnsErrorIfFailingToParseIQID(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='abczzzz'></client:iq>")}
	conn := Conn{
		in: xml.NewDecoder(mockIn),
	}

	_, err := conn.Next()
	c.Assert(err.Error(), Equals, "xmpp: failed to parse id from iq: strconv.ParseUint: parsing \"abczzzz\": invalid syntax")
}

func (s *ConnectionXmppSuite) Test_Next_returnsNothingIfThereIsNoInflightMatching(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<client:iq xmlns:client='jabber:client' type='result' id='100000'></client:iq>")}
	conn := Conn{
		in: xml.NewDecoder(mockIn),
	}

	_, err := conn.Next()
	c.Assert(err, Equals, io.EOF)
}

func (s *ConnectionXmppSuite) Test_makeInOut_returnsANewDecoderAndOriginalWriterWhenNoConfigIsGiven(c *C) {
	mockBoth := &mockConnIOReaderWriter{}
	_, rout := makeInOut(mockBoth, nil)
	c.Assert(rout, Equals, mockBoth)
}

func (s *ConnectionXmppSuite) Test_makeInOut_returnsANewDecoderAndWrappedWriterWhenConfigIsGiven(c *C) {
	mockBoth := &mockConnIOReaderWriter{}
	mockInLog := &mockConnIOReaderWriter{}
	config := &Config{InLog: mockInLog, OutLog: mockInLog}
	_, rout := makeInOut(mockBoth, config)
	c.Assert(rout, Not(Equals), mockBoth)
}
