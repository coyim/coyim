package xmpp

import (
	"errors"
	"io"
	"reflect"

	"github.com/coyim/coyim/xmpp/data"

	. "gopkg.in/check.v1"
)

type IqXMPPSuite struct{}

var _ = Suite(&IqXMPPSuite{})

func (s *IqXMPPSuite) Test_SendIQReply_returnsErrorIfOneIsEncounteredWhenWriting(c *C) {
	mockIn := &mockConnIOReaderWriter{err: errors.New("some error")}
	conn := conn{
		log: testLogger(),
		out: mockIn,
		jid: "somewhat@foo.com/somewhere",
	}

	err := conn.SendIQReply("fo", "bar", "baz", nil)
	c.Assert(err.Error(), Equals, "some error")
}

func (s *IqXMPPSuite) Test_SendIQReply_writesAnEmptyReplyIfEmptyIsGiven(c *C) {
	mockIn := &mockConnIOReaderWriter{}
	conn := conn{
		log: testLogger(),
		out: mockIn,
		jid: "som'ewhat@foo.com/somewhere",
	}

	err := conn.SendIQReply("f&o", "b\"ar", "b<az", data.EmptyReply{})
	c.Assert(err, IsNil)
	c.Assert(string(mockIn.write), Equals, "<iq xmlns='jabber:client' to='f&amp;o' from='som&apos;ewhat@foo.com/somewhere' type='b&quot;ar' id='b&lt;az'></iq>")
}

func (s *IqXMPPSuite) Test_SendIQReply_returnsErrorIfAnUnXMLableEntryIsGiven(c *C) {
	mockIn := &mockConnIOReaderWriter{}
	conn := conn{
		log: testLogger(),
		out: mockIn,
		jid: "som'ewhat@foo.com/somewhere",
	}
	err := conn.SendIQReply("f&o", "b\"ar", "b<az", func() int { return 42 })
	c.Assert(err.Error(), Equals, "xml: unsupported type: func() int")
}

func (s *IqXMPPSuite) Test_SendIQ_returnsErrorIfWritingDataFails(c *C) {
	mockIn := &mockConnIOReaderWriter{err: errors.New("this also fails")}
	conn := conn{
		log: testLogger(),
		out: mockIn,
		jid: "som'ewhat@foo.com/somewhere",
	}
	_, _, err := conn.SendIQ("", "", nil)
	c.Assert(err.Error(), Equals, "this also fails")
}

func (s *IqXMPPSuite) Test_Send_returnsErrorIfAnUnXMLableEntryIsGiven(c *C) {
	mockIn := &mockConnIOReaderWriter{}
	conn := conn{
		log: testLogger(),
		out: mockIn,
		jid: "som'ewhat@foo.com/somewhere",
	}
	_, _, err := conn.SendIQ("", "", func() int { return 42 })
	c.Assert(err.Error(), Equals, "xml: unsupported type: func() int")
}

func (s *IqXMPPSuite) TestConnSendIQReplyAndTyp(c *C) {
	mockOut := mockConnIOReaderWriter{}
	conn := conn{
		log: testLogger(),
		out: &mockOut,
		jid: "jid",
	}
	conn.inflights = make(map[data.Cookie]inflight)
	reply, cookie, err := conn.SendIQ("example@xmpp.com", "typ", nil)
	c.Assert(string(mockOut.write), Matches, "<iq xmlns='jabber:client' to='example@xmpp.com' from='jid' type='typ' id='.+'></iq>")
	c.Assert(reply, NotNil)
	c.Assert(cookie, NotNil)
	c.Assert(err, IsNil)
}

func (s *IqXMPPSuite) TestConnSendIQRaw(c *C) {
	mockOut := mockConnIOReaderWriter{}
	conn := conn{
		log: testLogger(),
		out: &mockOut,
		jid: "jid",
	}

	conn.inflights = make(map[data.Cookie]inflight)
	reply, cookie, err := conn.SendIQ("example@xmpp.com", "typ", rawXML("<foo param='bar' />"))
	c.Assert(string(mockOut.write), Matches, "<iq xmlns='jabber:client' to='example@xmpp.com' from='jid' type='typ' id='.*'><foo param='bar' /></iq>")
	c.Assert(reply, NotNil)
	c.Assert(cookie, NotNil)
	c.Assert(err, IsNil)
}

func (s *IqXMPPSuite) TestConnSendIQErr(c *C) {
	mockOut := mockConnIOReaderWriter{err: io.EOF}
	conn := conn{
		log: testLogger(),
		out: &mockOut,
		jid: "jid",
	}
	reply, cookie, err := conn.SendIQ("example@xmpp.com", "typ", nil)
	c.Assert(string(mockOut.write), Matches, "<iq xmlns='jabber:client' to='example@xmpp.com' from='jid' type='typ' id='.*'></iq>$")
	c.Assert(reply, IsNil)
	c.Assert(cookie, Equals, data.Cookie(0))
	c.Assert(err, Equals, io.EOF)
}

func (s *IqXMPPSuite) TestConnSendIQEmptyReply(c *C) {
	mockOut := mockConnIOReaderWriter{}
	conn := conn{
		log: testLogger(),
		out: &mockOut,
		jid: "jid",
	}
	conn.inflights = make(map[data.Cookie]inflight)
	reply, cookie, err := conn.SendIQ("example@xmpp.com", "typ", reflect.ValueOf(data.EmptyReply{}))
	c.Assert(string(mockOut.write), Matches, "<iq xmlns='jabber:client' to='example@xmpp.com' from='jid' type='typ' id='.+'><Value><flag>153</flag></Value></iq>")
	c.Assert(reply, NotNil)
	c.Assert(cookie, NotNil)
	c.Assert(err, IsNil)
}

func (s *IqXMPPSuite) TestConnSendIQReply(c *C) {
	mockOut := mockConnIOReaderWriter{}
	conn := conn{
		log: testLogger(),
		out: &mockOut,
		jid: "jid",
	}
	err := conn.SendIQReply("example@xmpp.com", "typ", "id", nil)
	c.Assert(string(mockOut.write), Matches, "<iq xmlns='jabber:client' to='example@xmpp.com' from='jid' type='typ' id='id'></iq>")
	c.Assert(err, IsNil)
}
