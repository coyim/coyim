package xmpp

import (
	"encoding/xml"
	"reflect"

	"github.com/twstrike/coyim/xmpp/data"

	. "gopkg.in/check.v1"
)

type XMLXMPPSuite struct{}

var _ = Suite(&XMLXMPPSuite{})

func (s *XMLXMPPSuite) Test_xmlEscape_escapesSpecialCharactersButNotRegularOnes(c *C) {
	res := xmlEscape("abc\"<foo>bar>bar\\and &her'e;")
	c.Assert(res, Equals, "abc&quot;&lt;foo&gt;bar&gt;bar\\and &amp;her&apos;e;")
}

type testFoo struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams foo"`
	To      string   `xml:"to,attr"`
}

func (s *XMLXMPPSuite) Test_next_usesCustomStorageIfAvailable(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:foo xmlns:stream='http://etherx.jabber.org/streams' to='hello'></stream:foo>")}
	conn := conn{
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

func (s *XMLXMPPSuite) Test_next_causesErrorWhenTryingToDecodeWrong(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:foo xmlns:stream='http://etherx.jabber.org/streams'><something></foo></stream:foo>")}
	conn := conn{
		in: xml.NewDecoder(mockIn),
		customStorage: map[xml.Name]reflect.Type{
			xml.Name{Space: NsStream, Local: "foo"}: reflect.TypeOf(testFoo{}),
		},
	}

	_, _, e := next(&conn)
	c.Assert(e.Error(), Equals, "XML syntax error on line 1: element <something> closed by </foo>")
}

func (s *XMLXMPPSuite) Test_ClientMesage_unmarshalsXMPPExtensions(c *C) {
	datax := `
	<message
    xmlns='jabber:client'
		xmlns:stream='http://etherx.jabber.org/streams'
		id='coyim1234'
    from='bernardo@shakespeare.lit/pda'
		to='francisco@shakespeare.lit/elsinore'
		type='chat'>
		<composing xmlns='http://jabber.org/protocol/chatstates'/>
		<x xmlns='jabber:x:event'>
	    <offline/>
			<delivered/>
			<composing/>
		</x>
	</message>
	`

	v := &data.ClientMessage{}
	err := xml.Unmarshal([]byte(datax), v)
	c.Assert(err, Equals, nil)

	c.Assert(v.ID, Equals, "coyim1234")
	c.Assert(v.From, Equals, "bernardo@shakespeare.lit/pda")
	c.Assert(v.To, Equals, "francisco@shakespeare.lit/elsinore")
	c.Assert(v.Type, Equals, "chat")
	c.Assert(v.Extensions, DeepEquals, data.Extensions{
		&data.Extension{XMLName: xml.Name{Space: "http://jabber.org/protocol/chatstates", Local: "composing"}},
		&data.Extension{
			XMLName: xml.Name{Space: "jabber:x:event", Local: "x"},
			Body:    "\n\t    <offline/>\n\t\t\t<delivered/>\n\t\t\t<composing/>\n\t\t",
		},
	})
}
