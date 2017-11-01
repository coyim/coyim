package xmpp

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

type DiscoveryXMPPSuite struct{}

var _ = Suite(&DiscoveryXMPPSuite{})

func (s *DiscoveryXMPPSuite) Test_SendDiscoveryInfoRequest(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		out: mockOut,
		jid: "juliet@example.com/chamber",
	}
	conn.inflights = make(map[data.Cookie]inflight)

	reply, cookie, err := conn.sendDiscoveryInfo("example.com")
	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Matches, "<iq xmlns='jabber:client' to='example.com' from='juliet@example.com/chamber' type='get' id='.+'><query xmlns=\"http://jabber.org/protocol/disco#info\"></query></iq>")
	c.Assert(reply, NotNil)
	c.Assert(cookie, NotNil)
}

func (s *DiscoveryXMPPSuite) Test_ReceiveDiscoveryResult(c *C) {
	// See: XEP-0030, Section: 3.1 Basic Protocol, Example: 2
	fromServer := `
<iq xmlns='jabber:client' type='result'
    from='plays.shakespeare.lit'
    to='romeo@montague.net/orchard'
    id='100000'>
  <query xmlns='http://jabber.org/protocol/disco#info'>
    <identity
        category='conference'
        type='text'
        name='Play-Specific Chatrooms'/>
    <identity
        category='directory'
        type='chatroom'
        name='Play-Specific Chatrooms'/>
    <feature var='http://jabber.org/protocol/disco#info'/>
    <feature var='http://jabber.org/protocol/disco#items'/>
  </query>
</iq>
`
	reply := make(chan data.Stanza, 1)
	mockIn := &mockConnIOReaderWriter{read: []byte(fromServer)}
	conn := conn{
		in: xml.NewDecoder(mockIn),
		inflights: map[data.Cookie]inflight{
			0x100000: inflight{
				to:        "plays.shakespeare.lit",
				replyChan: reply,
			},
		},
	}

	_, err := conn.Next()
	c.Assert(err, Equals, io.EOF)

	iq, ok := (<-reply).Value.(*data.ClientIQ)
	c.Assert(ok, Equals, true)

	discoveryReply, err := parseDiscoveryReply(iq)
	c.Assert(err, IsNil)

	c.Assert(discoveryReply, DeepEquals, data.DiscoveryReply{
		XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "query"},
		Identities: []data.DiscoveryIdentity{
			data.DiscoveryIdentity{
				XMLName:  xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "identity"},
				Category: "conference",
				Type:     "text",
				Name:     "Play-Specific Chatrooms"},
			data.DiscoveryIdentity{
				XMLName:  xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "identity"},
				Category: "directory",
				Type:     "chatroom",
				Name:     "Play-Specific Chatrooms"},
		},
		Features: []data.DiscoveryFeature{
			data.DiscoveryFeature{
				XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
				Var:     "http://jabber.org/protocol/disco#info",
			},
			data.DiscoveryFeature{
				XMLName: xml.Name{Space: "http://jabber.org/protocol/disco#info", Local: "feature"},
				Var:     "http://jabber.org/protocol/disco#items",
			},
		},
	})

	_, ok = conn.inflights[0x100000]
	c.Assert(ok, Equals, false)
}

func (s *DiscoveryXMPPSuite) Test_HasSupportTo(c *C) {
	// See: XEP-0030, Section: 3.1 Basic Protocol, Example: 2
	fromServer := `
<iq xmlns='jabber:client' type='result'
    from='plays.shakespeare.lit'
    to='romeo@montague.net/orchard'
    id='%s'>
  <query xmlns='http://jabber.org/protocol/disco#info'>
    <feature var='jabber:iq:privacy'/>
    <feature var='http://jabber.org/protocol/disco#items'/>
  </query>
</iq>
`
	mockOut := &mockConnIOReaderWriter{}

	conn := conn{
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan bool, 1)
	go func() {
		ok := conn.HasSupportTo("plays.shakespeare.lit", "jabber:iq:privacy")
		c.Assert(ok, Equals, true)
		done <- true
	}()

	<-time.After(1 * time.Millisecond)

	var iq data.ClientIQ
	c.Assert(string(mockOut.Written()), Matches, "<iq xmlns='jabber:client' to='plays.shakespeare.lit' from='romeo@montague.net/orchard' type='get' id='.+'><query xmlns=\"http://jabber.org/protocol/disco#info\"></query></iq>")

	err := xml.Unmarshal(mockOut.Written(), &iq)
	c.Assert(err, IsNil)

	mockIn := &mockConnIOReaderWriter{read: []byte(fmt.Sprintf(fromServer, iq.ID))}
	conn.in = xml.NewDecoder(mockIn)

	_, err = conn.Next()
	c.Assert(err, Equals, io.EOF)
	<-done
}

func (s *DiscoveryXMPPSuite) Test_VerificationString_failsIfThereAreDuplicateIdentities(c *C) {
	reply := &data.DiscoveryReply{
		Identities: []data.DiscoveryIdentity{
			data.DiscoveryIdentity{
				Lang:     "en",
				Category: "stuff",
				Type:     "thing",
				Name:     "something",
			},
			data.DiscoveryIdentity{
				Lang:     "en",
				Category: "stuff",
				Type:     "thing",
				Name:     "something",
			},
		},
	}

	_, err := VerificationString(reply)
	c.Assert(err.Error(), Equals, "duplicate discovery identity")
}

func (s *DiscoveryXMPPSuite) Test_VerificationString_failsIfThereAreDuplicateFeatures(c *C) {
	reply := &data.DiscoveryReply{
		Features: []data.DiscoveryFeature{
			data.DiscoveryFeature{
				Var: "foo",
			},
			data.DiscoveryFeature{
				Var: "foo",
			},
		},
	}

	_, err := VerificationString(reply)
	c.Assert(err.Error(), Equals, "duplicate discovery feature")
}

func (s *DiscoveryXMPPSuite) Test_VerificationString_failsIfThereAreDuplicateFormTypes(c *C) {
	reply := &data.DiscoveryReply{
		Forms: []data.Form{
			data.Form{
				Type: "foo",
				Fields: []data.FormFieldX{
					data.FormFieldX{
						Var:    "FORM_TYPE",
						Type:   "Foo",
						Values: []string{"foo"},
					},
				},
			},
			data.Form{
				Type: "foo",
				Fields: []data.FormFieldX{
					data.FormFieldX{
						Var:    "FORM_TYPE",
						Type:   "Foo",
						Values: []string{"foo"},
					},
				},
			},
		},
	}

	_, err := VerificationString(reply)
	c.Assert(err.Error(), Equals, "multiple forms of the same type")
}

func (s *DiscoveryXMPPSuite) Test_VerificationString_failsIfThereAreNoValues(c *C) {
	reply := &data.DiscoveryReply{
		Forms: []data.Form{
			data.Form{
				Type: "foo",
			},
			data.Form{
				Type: "foo",
				Fields: []data.FormFieldX{
					data.FormFieldX{
						Var:  "FORM_TYPE2",
						Type: "Foo",
					},
				},
			},
			data.Form{
				Type: "foo",
				Fields: []data.FormFieldX{
					data.FormFieldX{
						Var:  "FORM_TYPE",
						Type: "Foo",
					},
				},
			},
			data.Form{
				Type: "foo",
				Fields: []data.FormFieldX{
					data.FormFieldX{
						Var:  "FORM_TYPE",
						Type: "Foo",
					},
				},
			},
		},
	}

	_, err := VerificationString(reply)
	c.Assert(err.Error(), Equals, "form does not have a single FORM_TYPE value")
}

func (s *DiscoveryXMPPSuite) Test_DiscoveryReply_returnsSupportedValues(c *C) {
	rep := DiscoveryReply("foo@bar.com")
	c.Assert(rep, DeepEquals,
		data.DiscoveryReply{
			XMLName: xml.Name{Space: "", Local: ""},
			Node:    "",
			Identities: []data.DiscoveryIdentity{
				data.DiscoveryIdentity{
					XMLName:  xml.Name{Space: "", Local: ""},
					Lang:     "",
					Category: "client",
					Type:     "pc",
					Name:     "foo@bar.com"}},
			Features: []data.DiscoveryFeature{
				{Var: "http://jabber.org/protocol/disco#info"},
				{Var: "urn:xmpp:bob"},
				{Var: "urn:xmpp:ping"},
				{Var: "http://jabber.org/protocol/caps"},
				{Var: "jabber:iq:version"},
				{Var: "vcard-temp"},
				{Var: "jabber:x:data"},
				{Var: "http://jabber.org/protocol/si"},
				{Var: "http://jabber.org/protocol/si/profile/file-transfer"},
				{Var: "http://jabber.org/protocol/si/profile/directory-transfer"},
				{Var: "http://jabber.org/protocol/bytestreams"},
			},
			Forms: []data.Form(nil)})
}
