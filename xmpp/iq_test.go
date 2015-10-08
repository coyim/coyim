package xmpp

import (
	"errors"

	. "gopkg.in/check.v1"
)

type IqXmppSuite struct{}

var _ = Suite(&IqXmppSuite{})

func (s *IqXmppSuite) Test_SendIQReply_returnsErrorIfOneIsEncounteredWhenWriting(c *C) {
	mockIn := &mockConnIOReaderWriter{err: errors.New("some error")}
	conn := Conn{
		out: mockIn,
		jid: "somewhat@foo.com/somewhere",
	}

	err := conn.SendIQReply("fo", "bar", "baz", nil)
	c.Assert(err.Error(), Equals, "some error")
}

func (s *IqXmppSuite) Test_SendIQReply_writesAnEmptyReplyIfEmptyIsGiven(c *C) {
	mockIn := &mockConnIOReaderWriter{}
	conn := Conn{
		out: mockIn,
		jid: "som'ewhat@foo.com/somewhere",
	}

	err := conn.SendIQReply("f&o", "b\"ar", "b<az", EmptyReply{})
	c.Assert(err, IsNil)
	c.Assert(string(mockIn.write), Equals, "<iq to='f&amp;o' from='som&apos;ewhat@foo.com/somewhere' type='b&quot;ar' id='b&lt;az'></iq>")
}

type testNonXmlValue struct {
	x int
}

func (s *IqXmppSuite) Test_SendIQReply_returnsErrorIfAnUnXMLableEntryIsGiven(c *C) {
	mockIn := &mockConnIOReaderWriter{}
	conn := Conn{
		out: mockIn,
		jid: "som'ewhat@foo.com/somewhere",
	}
	err := conn.SendIQReply("f&o", "b\"ar", "b<az", func() int { return 42 })
	c.Assert(err.Error(), Equals, "xml: unsupported type: func() int")
}

func (s *IqXmppSuite) Test_SendIQ_returnsErrorIfWritingDataFails(c *C) {
	mockIn := &mockConnIOReaderWriter{err: errors.New("this also fails...")}
	conn := Conn{
		out: mockIn,
		jid: "som'ewhat@foo.com/somewhere",
	}
	_, _, err := conn.SendIQ("", "", nil)
	c.Assert(err.Error(), Equals, "this also fails...")
}

func (s *IqXmppSuite) Test_Send_returnsErrorIfAnUnXMLableEntryIsGiven(c *C) {
	mockIn := &mockConnIOReaderWriter{}
	conn := Conn{
		out: mockIn,
		jid: "som'ewhat@foo.com/somewhere",
	}
	_, _, err := conn.SendIQ("", "", func() int { return 42 })
	c.Assert(err.Error(), Equals, "xml: unsupported type: func() int")
}

func (s *IqXmppSuite) Test_SendIQ_returnsErrorIfWritingDataFailsTheSecondTime(c *C) {
	mockIn := &mockConnIOReaderWriter{err: errors.New("this also fails again..."), errCount: 1}
	conn := Conn{
		out: mockIn,
		jid: "som'ewhat@foo.com/somewhere",
	}
	_, _, err := conn.SendIQ("", "", nil)
	c.Assert(err.Error(), Equals, "this also fails again...")
	c.Assert(string(mockIn.write), Matches, "<iq  from='som&apos;ewhat@foo.com/somewhere' type='' id='.*?'></iq>")
}
