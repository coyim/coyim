package xmpp

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"testing"

	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/xmpp/data"
	"github.com/twstrike/gotk3adapter/glib_mock"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	log.SetOutput(ioutil.Discard)
	i18n.InitLocalization(&glib_mock.Mock{})
}

type XMPPSuite struct{}

var _ = Suite(&XMPPSuite{})

func (s *XMPPSuite) TestDiscoReplyVerSimple(c *C) {
	expect := "QgayPKawpkPSDYmwT/WM94uAlu0="
	input := []byte(`
  <query xmlns='http://jabber.org/protocol/disco#info'
         node='http://code.google.com/p/exodus#QgayPKawpkPSDYmwT/WM94uAlu0='>
    <identity category='client' name='Exodus 0.9.1' type='pc'/>
    <feature var='http://jabber.org/protocol/caps'/>
    <feature var='http://jabber.org/protocol/disco#info'/>
    <feature var='http://jabber.org/protocol/disco#items'/>
    <feature var='http://jabber.org/protocol/muc'/>
  </query>
  `)
	var dr data.DiscoveryReply
	c.Assert(xml.Unmarshal(input, &dr), IsNil)
	hash, err := VerificationString(&dr)
	c.Assert(err, IsNil)
	c.Assert(hash, Equals, expect)
}

func (s *XMPPSuite) TestDiscoReplyVerComplex(c *C) {
	expect := "q07IKJEyjvHSyhy//CH0CxmKi8w="
	input := []byte(`
  <query xmlns='http://jabber.org/protocol/disco#info'
         node='http://psi-im.org#q07IKJEyjvHSyhy//CH0CxmKi8w='>
    <identity xml:lang='en' category='client' name='Psi 0.11' type='pc'/>
    <identity xml:lang='el' category='client' name='Ψ 0.11' type='pc'/>
    <feature var='http://jabber.org/protocol/caps'/>
    <feature var='http://jabber.org/protocol/disco#info'/>
    <feature var='http://jabber.org/protocol/disco#items'/>
    <feature var='http://jabber.org/protocol/muc'/>
    <x xmlns='jabber:x:data' type='result'>
      <field var='FORM_TYPE' type='hidden'>
        <value>urn:xmpp:dataforms:softwareinfo</value>
      </field>
      <field var='ip_version'>
        <value>ipv4</value>
        <value>ipv6</value>
      </field>
      <field var='os'>
        <value>Mac</value>
      </field>
      <field var='os_version'>
        <value>10.5.1</value>
      </field>
      <field var='software'>
        <value>Psi</value>
      </field>
      <field var='software_version'>
        <value>0.11</value>
      </field>
    </x>
  </query>
`)
	var dr data.DiscoveryReply
	c.Assert(xml.Unmarshal(input, &dr), IsNil)
	hash, err := VerificationString(&dr)
	c.Assert(err, IsNil)
	c.Assert(hash, Equals, expect)
}

func (s *XMPPSuite) TestConnClose_sendsAStreamCloseTagWhenWeCloseFirst(c *C) {
	mockIn := &mockConnIOReaderWriter{
		read: []byte("<?xml version='1.0'?><str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'></str:stream>"),
	}
	mockCloser := &mockConnIOReaderWriter{}

	conn := newConn()
	conn.in = xml.NewDecoder(mockIn)
	conn.out = mockCloser
	conn.rawOut = mockCloser

	nextElement(conn.in) // Reads the opening tag and make the unmarshaller happy

	done := make(chan bool)
	go func() {
		stanza, err := conn.Next() // Reads the closing tag

		c.Assert(err, IsNil)
		c.Assert(stanza, DeepEquals, data.Stanza{
			Name:  xml.Name{Space: "http://etherx.jabber.org/streams", Local: "stream"},
			Value: &data.StreamClose{},
		})

		done <- true
	}()

	// blocks until it receives the </stream> or timeouts
	c.Assert(conn.Close(), IsNil)
	c.Assert(mockCloser.calledClose, Equals, 1)
	c.Assert(mockCloser.write, DeepEquals, []byte("</stream:stream>"))

	<-done
}

func (s *XMPPSuite) TestConnNext_replyWithAStreamCloseTagWhenTheyCloseFirst(c *C) {
	mockIn := &mockConnIOReaderWriter{
		read: []byte("<?xml version='1.0'?><str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'></str:stream>"),
	}
	mockCloser := &mockConnIOReaderWriter{}

	conn := newConn()
	conn.in = xml.NewDecoder(mockIn)
	conn.out = mockCloser
	conn.rawOut = mockCloser

	nextElement(conn.in)       // Reads the opening tag and make the unmarshaller happy
	stanza, err := conn.Next() // Reads the closing tag

	c.Assert(err, IsNil)
	c.Assert(stanza, DeepEquals, data.Stanza{
		Name:  xml.Name{Space: "http://etherx.jabber.org/streams", Local: "stream"},
		Value: &data.StreamClose{},
	})

	c.Assert(mockCloser.calledClose, Equals, 1)
	c.Assert(mockCloser.write, DeepEquals, []byte("</stream:stream>"))

	err = conn.Close()
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "xmpp: the connection is already closed")
}

func (s *XMPPSuite) TestConnNextEOF(c *C) {
	mockIn := &mockConnIOReaderWriter{err: io.EOF}
	conn := conn{
		in: xml.NewDecoder(mockIn),
	}
	stanza, err := conn.Next()
	c.Assert(stanza.Name, Equals, xml.Name{})
	c.Assert(stanza.Value, IsNil)
	c.Assert(err, Equals, io.EOF)
}

func (s *XMPPSuite) TestConnNextErr(c *C) {
	mockIn := &mockConnIOReaderWriter{
		read: []byte(`
      <field var='os'>
        <value>Mac</value>
      </field>
		`),
	}
	conn := conn{
		in: xml.NewDecoder(mockIn),
	}
	stanza, err := conn.Next()
	c.Assert(stanza.Name, Equals, xml.Name{})
	c.Assert(stanza.Value, IsNil)
	c.Assert(err.Error(), Equals, "unexpected XMPP message  <field/>")
}

func (s *XMPPSuite) TestConnNextIQSet(c *C) {
	mockIn := &mockConnIOReaderWriter{
		read: []byte(`
<iq to='example.com'
    xmlns='jabber:client'
    type='set'
    id='sess_1'>
  <session xmlns='urn:ietf:params:xml:ns:xmpp-session'/>
</iq>
  `),
	}
	conn := conn{
		in: xml.NewDecoder(mockIn),
	}
	stanza, err := conn.Next()
	c.Assert(stanza.Name, Equals, xml.Name{Space: NsClient, Local: "iq"})
	iq, ok := stanza.Value.(*data.ClientIQ)
	c.Assert(ok, Equals, true)
	c.Assert(iq.To, Equals, "example.com")
	c.Assert(iq.Type, Equals, "set")
	c.Assert(err, IsNil)
}

func (s *XMPPSuite) TestConnNextIQResult(c *C) {
	mockIn := &mockConnIOReaderWriter{
		read: []byte(`
<iq from='example.com'
    xmlns='jabber:client'
    type='result'
    id='sess_1'/>
  `),
	}
	conn := conn{
		in: xml.NewDecoder(mockIn),
	}
	stanza, err := conn.Next()
	c.Assert(stanza.Name, Equals, xml.Name{Space: NsClient, Local: "iq"})
	iq, ok := stanza.Value.(*data.ClientIQ)
	c.Assert(ok, Equals, true)
	c.Assert(iq.From, Equals, "example.com")
	c.Assert(iq.Type, Equals, "result")
	c.Assert(err, ErrorMatches, "xmpp: failed to parse id from iq: .*")
}

func (s *XMPPSuite) TestConnCancelError(c *C) {
	conn := conn{}
	ok := conn.Cancel(conn.getCookie())
	c.Assert(ok, Equals, false)
}

func (s *XMPPSuite) TestConnCancelOK(c *C) {
	conn := conn{}
	cookie := conn.getCookie()
	ch := make(chan data.Stanza, 1)
	conn.inflights = make(map[data.Cookie]inflight)
	conn.inflights[cookie] = inflight{ch, ""}
	ok := conn.Cancel(cookie)
	c.Assert(ok, Equals, true)
	_, ok = conn.inflights[cookie]
	c.Assert(ok, Equals, false)
}

func (s *XMPPSuite) TestConnRequestRoster(c *C) {
	mockOut := mockConnIOReaderWriter{}
	conn := conn{
		out: &mockOut,
	}
	conn.inflights = make(map[data.Cookie]inflight)
	ch, cookie, err := conn.RequestRoster()
	c.Assert(string(mockOut.write), Matches, "<iq type='get' id='.*'><query xmlns='jabber:iq:roster'/></iq>")
	c.Assert(ch, NotNil)
	c.Assert(cookie, NotNil)
	c.Assert(err, IsNil)
}

func (s *XMPPSuite) TestConnRequestRosterErr(c *C) {
	mockOut := mockConnIOReaderWriter{err: io.EOF}
	conn := conn{
		out: &mockOut,
	}
	conn.inflights = make(map[data.Cookie]inflight)
	ch, cookie, err := conn.RequestRoster()
	c.Assert(string(mockOut.write), Matches, "<iq type='get' id='.*'><query xmlns='jabber:iq:roster'/></iq>")
	c.Assert(ch, IsNil)
	c.Assert(cookie, NotNil)
	c.Assert(err, Equals, io.EOF)
}

func (s *XMPPSuite) TestParseRoster(c *C) {
	iq := data.ClientIQ{}
	iq.Query = []byte(`
  <query xmlns='jabber:iq:roster'>
    <item jid='romeo@example.net'
          name='Romeo'
          subscription='both'>
      <group>Friends</group>
    </item>
    <item jid='mercutio@example.org'
          name='Mercutio'
          subscription='from'>
      <group>Friends</group>
    </item>
    <item jid='benvolio@example.org'
          name='Benvolio'
          subscription='both'>
      <group>Friends</group>
    </item>
  </query>
  `)
	reply := data.Stanza{
		Value: &iq,
	}
	rosterEntrys, err := data.ParseRoster(reply)
	c.Assert(rosterEntrys, NotNil)
	c.Assert(err, IsNil)
}

func (s *XMPPSuite) TestConnSend(c *C) {
	mockOut := mockConnIOReaderWriter{}
	conn := conn{
		out: &mockOut,
		jid: "jid",
	}
	err := conn.Send("example@xmpp.com", "message")
	c.Assert(string(mockOut.write), Matches, "<message to='example@xmpp.com' from='jid' type='chat'><body>message</body><nos:x xmlns:nos='google:nosave' value='enabled'/><arc:record xmlns:arc='http://jabber.org/protocol/archive' otr='require'/><no-copy xmlns='urn:xmpp:hints'/><no-permanent-store xmlns='urn:xmpp:hints'/><private xmlns='urn:xmpp:carbons:2'/></message>")
	c.Assert(err, IsNil)
}
