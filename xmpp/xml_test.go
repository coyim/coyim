package xmpp

import (
	"encoding/xml"
	"reflect"

	. "gopkg.in/check.v1"
)

type XMLXmppSuite struct{}

var _ = Suite(&XMLXmppSuite{})

func (s *XMLXmppSuite) Test_xmlEscape_escapesSpecialCharactersButNotRegularOnes(c *C) {
	res := xmlEscape("abc\"<foo>bar>bar\\and &her'e;")
	c.Assert(res, Equals, "abc&quot;&lt;foo&gt;bar&gt;bar\\and &amp;her&apos;e;")
}

type testFoo struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams foo"`
	To      string   `xml:"to,attr"`
}

func (s *XMLXmppSuite) Test_next_usesCustomStorageIfAvailable(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:foo xmlns:stream='http://etherx.jabber.org/streams' to='hello'></stream:foo>")}
	conn := Conn{
		in: xml.NewDecoder(mockIn),
		customStorage: map[xml.Name]reflect.Type{
			xml.Name{Space: NsStream, Local: "foo"}: reflect.TypeOf(testFoo{}),
		},
	}

	nm, i, e := next(&conn)
	c.Assert(e, IsNil)
	c.Assert(nm, Equals, xml.Name{Space: NsStream, Local: "foo"})
	val, _ := i.(*testFoo)
	c.Assert(*val, Equals, testFoo{XMLName: nm, To: "hello"})
}

func (s *XMLXmppSuite) Test_next_causesErrorWhenTryingToDecodeWrong(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:foo xmlns:stream='http://etherx.jabber.org/streams'><something></foo></stream:foo>")}
	conn := Conn{
		in: xml.NewDecoder(mockIn),
		customStorage: map[xml.Name]reflect.Type{
			xml.Name{Space: NsStream, Local: "foo"}: reflect.TypeOf(testFoo{}),
		},
	}

	_, _, e := next(&conn)
	c.Assert(e.Error(), Equals, "XML syntax error on line 1: element <something> closed by </foo>")
}
