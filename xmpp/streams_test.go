package xmpp

import (
	"encoding/xml"
	"errors"

	. "gopkg.in/check.v1"
)

type StreamsXmppSuite struct{}

var _ = Suite(&StreamsXmppSuite{})

func (s *StreamsXmppSuite) Test_sendInitialStreamHeader_returnsErrorIfSomethingGoesWrongWithFmtPrintf(c *C) {
	conn := Conn{
		out: &mockConnIOReaderWriter{err: errors.New("Hello")},
	}
	_, err := conn.sendInitialStreamHeader("foo.com")
	c.Assert(err, Not(IsNil))
}

func (s *StreamsXmppSuite) Test_sendInitialStreamHeader_returnsErrorIfSomethingGoesWrongWithReadingAStream(c *C) {
	mockIn := &mockConnIOReaderWriter{err: errors.New("Hello")}
	conn := Conn{
		out: &mockConnIOReaderWriter{},
		in:  xml.NewDecoder(mockIn),
	}
	_, err := conn.sendInitialStreamHeader("foo.com")
	c.Assert(err, Not(IsNil))
}

func (s *StreamsXmppSuite) Test_sendInitialStreamHeader_sendsInitialStreamHeaderToOutput(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{err: errors.New("Hello")}
	conn := Conn{
		out: mockOut,
		in:  xml.NewDecoder(mockIn),
	}
	conn.sendInitialStreamHeader("somewhere.org")
	c.Assert(string(mockOut.write), Equals, "<?xml version='1.0'?><stream:stream to='somewhere.org' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n")
}

func (s *StreamsXmppSuite) Test_sendInitialStreamHeader_expectsResponseStreamHeaderInReturn(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><stream:stream xmlns:stream='http://etherx.jabber.org/streams' version='1.0'></stream:stream>")}
	conn := Conn{
		out: mockOut,
		in:  xml.NewDecoder(mockIn),
	}
	_, err := conn.sendInitialStreamHeader("somewhereElse.org")
	c.Assert(err.Error(), Equals, "unmarshal <features>: EOF")
}

func (s *StreamsXmppSuite) Test_sendInitialStreamHeader_failsIfReturnedStreamIsNotCorrectNamespace(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><str:stream xmlns:str='http://etherx.jabber.org/streams2' version='1.0'>")}
	conn := Conn{
		out: mockOut,
		in:  xml.NewDecoder(mockIn),
	}
	_, err := conn.sendInitialStreamHeader("somewhereElse.org")
	c.Assert(err.Error(), Equals, "xmpp: expected <stream> but got <stream> in http://etherx.jabber.org/streams2")
}

func (s *StreamsXmppSuite) Test_sendInitialStreamHeader_failsIfReturnedElementIsNotStream(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><str:feature xmlns:str='http://etherx.jabber.org/streams' version='1.0'>")}
	conn := Conn{
		out: mockOut,
		in:  xml.NewDecoder(mockIn),
	}
	_, err := conn.sendInitialStreamHeader("somewhereElse.org")
	c.Assert(err.Error(), Equals, "xmpp: expected <stream> but got <feature> in http://etherx.jabber.org/streams")
}

func (s *StreamsXmppSuite) Test_sendInitialStreamHeader_expectsFeaturesInReturn(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	mockIn := &mockConnIOReaderWriter{read: []byte("<?xml version='1.0'?><str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'><str:features></str:features>")}
	conn := Conn{
		out: mockOut,
		in:  xml.NewDecoder(mockIn),
	}
	feat, err := conn.sendInitialStreamHeader("somewhereElse.org")
	c.Assert(err, IsNil)
	expected := streamFeatures{}
	expected.XMLName = xml.Name{Space: "http://etherx.jabber.org/streams", Local: "features"}
	c.Assert(feat, DeepEquals, expected)
}
