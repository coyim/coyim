package xmpp

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"sync"

	"github.com/coyim/coyim/xmpp/data"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

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
		log: testLogger(),
		in:  xml.NewDecoder(mockIn),
		customStorage: map[xml.Name]reflect.Type{
			xml.Name{Space: NsStream, Local: "foo"}: reflect.TypeOf(testFoo{}),
		},
	}

	nm, i, e := next(&conn, conn.log)
	c.Assert(e, IsNil)
	c.Assert(nm, Equals, xml.Name{Space: NsStream, Local: "foo"})
	val, _ := i.(*testFoo)
	c.Assert(*val, Equals, testFoo{XMLName: nm, To: "hello"})
}

func (s *XMLXMPPSuite) Test_next_causesErrorWhenTryingToDecodeWrong(c *C) {
	mockIn := &mockConnIOReaderWriter{read: []byte("<stream:foo xmlns:stream='http://etherx.jabber.org/streams'><something></foo></stream:foo>")}
	conn := conn{
		log: testLogger(),
		in:  xml.NewDecoder(mockIn),
		customStorage: map[xml.Name]reflect.Type{
			xml.Name{Space: NsStream, Local: "foo"}: reflect.TypeOf(testFoo{}),
		},
	}

	_, _, e := next(&conn, conn.log)
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

func (s *XMLXMPPSuite) Test_decodeEndElement_returnsStreamEnd(c *C) {
	name, val, err := decodeEndElement(xml.EndElement{Name: xml.Name{Space: "http://etherx.jabber.org/streams", Local: "stream"}})
	c.Assert(err, IsNil)
	c.Assert(val, DeepEquals, &data.StreamClose{})
	c.Assert(name, Equals, xml.Name{Space: "http://etherx.jabber.org/streams", Local: "stream"})
}

func (s *XMLXMPPSuite) Test_decodeEndElement_returnsOtherElement(c *C) {
	name, val, err := decodeEndElement(xml.EndElement{Name: xml.Name{Space: "foo", Local: "bar"}})
	c.Assert(err, IsNil)
	c.Assert(val, IsNil)
	c.Assert(name, Equals, xml.Name{Space: "foo", Local: "bar"})
}

type xmlMockConn struct {
	sync.Mutex

	dec *xml.Decoder
}

func (m *xmlMockConn) In() *xml.Decoder {
	return m.dec
}

func (m *xmlMockConn) Lock() *sync.Mutex {
	return &m.Mutex
}

func (m *xmlMockConn) CustomStorage() map[xml.Name]reflect.Type {
	return nil
}

func (s *XMLXMPPSuite) Test_next_ignoresOtherElements(c *C) {
	l, hook := test.NewNullLogger()
	l.SetLevel(log.DebugLevel)

	np := []byte(`<?xml version="1.0"?><?xml something?><?foo hello?><!-- a comment -->
<foobarium>
</foobarium>
`)
	m := &xmlMockConn{
		dec: xml.NewDecoder(bytes.NewBuffer(np)),
	}

	name, val, err := next(m, l)
	c.Assert(name, DeepEquals, xml.Name{})
	c.Assert(val, IsNil)
	c.Assert(err, ErrorMatches, "unexpected XMPP message  <foobarium/>")

	c.Assert(len(hook.Entries), Equals, 4)
	c.Assert(hook.Entries[0].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[0].Message, Equals, "xmpp: received unhandled ProcInst element")
	c.Assert(hook.Entries[0].Data, HasLen, 2)
	c.Assert(string(hook.Entries[0].Data["inst"].([]byte)), Equals, "foobarium")
	c.Assert(hook.Entries[0].Data["target"], Equals, "xml")

	c.Assert(hook.Entries[1].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[1].Message, Equals, "xmpp: received unhandled ProcInst element")
	c.Assert(hook.Entries[1].Data, HasLen, 2)
	c.Assert(string(hook.Entries[1].Data["inst"].([]byte)), Equals, "fooba")
	c.Assert(hook.Entries[1].Data["target"], Equals, "foo")

	c.Assert(hook.Entries[2].Level, Equals, log.WarnLevel)
	c.Assert(hook.Entries[2].Message, Equals, "xmpp: received unhandled element")
	c.Assert(hook.Entries[2].Data, HasLen, 1)

	c.Assert(hook.Entries[3].Level, Equals, log.InfoLevel)
	c.Assert(hook.Entries[3].Message, Equals, "xmpp: received whitespace ping")
	c.Assert(hook.Entries[3].Data, HasLen, 0)
}
