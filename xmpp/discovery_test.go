package xmpp

import (
	"encoding/xml"
	"errors"
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
		log: testLogger(),
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
		log: testLogger(),
		in:  xml.NewDecoder(mockIn),
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

	discoveryReply, err := parseDiscoveryInfoReply(iq)
	c.Assert(err, IsNil)

	c.Assert(discoveryReply, DeepEquals, &data.DiscoveryInfoQuery{
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

func waitForInflightTo(c *conn, to string) {
	done := make(chan bool)

	go func() {
		for {
			time.Sleep(3 * time.Millisecond)
			for _, v := range c.inflights {
				if v.to == to {
					done <- true
					return
				}
			}
		}
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
	}
}

func (s *DiscoveryXMPPSuite) Test_conn_HasSupportTo(c *C) {
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
		log:       testLogger(),
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

	waitForInflightTo(&conn, "plays.shakespeare.lit")

	c.Assert(string(mockOut.Written()), Matches, "<iq xmlns='jabber:client' to='plays.shakespeare.lit' from='romeo@montague.net/orchard' type='get' id='.+'><query xmlns=\"http://jabber.org/protocol/disco#info\"></query></iq>")

	var iq data.ClientIQ
	err := xml.Unmarshal(mockOut.Written(), &iq)
	c.Assert(err, IsNil)

	mockIn := &mockConnIOReaderWriter{read: []byte(fmt.Sprintf(fromServer, iq.ID))}
	conn.in = xml.NewDecoder(mockIn)

	_, err = conn.Next()
	c.Assert(err, Equals, io.EOF)
	<-done
}

func (s *DiscoveryXMPPSuite) Test_conn_HasSupportTo_fails(c *C) {
	// See: XEP-0030, Section: 3.1 Basic Protocol, Example: 2
	fromServer := `
<iq xmlns='jabber:client' type='error'
    from='plays.shakespeare.lit'
    to='romeo@montague.net/orchard'
    id='%s'>
</iq>
`
	mockOut := &mockConnIOReaderWriter{}

	conn := conn{
		log:       testLogger(),
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan bool, 1)
	go func() {
		ok := conn.HasSupportTo("plays.shakespeare.lit", "jabber:iq:privacy")
		c.Assert(ok, Equals, false)
		done <- true
	}()

	waitForInflightTo(&conn, "plays.shakespeare.lit")

	var iq data.ClientIQ
	err := xml.Unmarshal(mockOut.Written(), &iq)
	c.Assert(err, IsNil)

	mockIn := &mockConnIOReaderWriter{read: []byte(fmt.Sprintf(fromServer, iq.ID))}
	conn.in = xml.NewDecoder(mockIn)

	_, err = conn.Next()
	c.Assert(err, Equals, io.EOF)
	<-done
}

func (s *DiscoveryXMPPSuite) Test_conn_HasSupportTo_doesntSupport(c *C) {
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
		log:       testLogger(),
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan bool, 1)
	go func() {
		ok := conn.HasSupportTo("plays.shakespeare.lit", "jabber:iq:privacy", "jabber:iq:something:else")
		c.Assert(ok, Equals, false)
		done <- true
	}()

	waitForInflightTo(&conn, "plays.shakespeare.lit")

	var iq data.ClientIQ
	err := xml.Unmarshal(mockOut.Written(), &iq)
	c.Assert(err, IsNil)

	mockIn := &mockConnIOReaderWriter{read: []byte(fmt.Sprintf(fromServer, iq.ID))}
	conn.in = xml.NewDecoder(mockIn)

	_, err = conn.Next()
	c.Assert(err, Equals, io.EOF)
	<-done
}

func (s *DiscoveryXMPPSuite) Test_VerificationString_failsIfThereAreDuplicateIdentities(c *C) {
	reply := &data.DiscoveryInfoQuery{
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
	reply := &data.DiscoveryInfoQuery{
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
	reply := &data.DiscoveryInfoQuery{
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
	reply := &data.DiscoveryInfoQuery{
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
	rep := DiscoveryReply("foo@bar.com", "")
	c.Assert(rep, DeepEquals,
		data.DiscoveryInfoQuery{
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
				{Var: "http://jabber.org/protocol/disco#items"},
				{Var: "urn:xmpp:bob"},
				{Var: "urn:xmpp:ping"},
				{Var: "http://jabber.org/protocol/caps"},
				{Var: "jabber:iq:version"},
				{Var: "vcard-temp"},
				{Var: "jabber:x:data"},
				{Var: "http://jabber.org/protocol/si"},
				{Var: "http://jabber.org/protocol/si/profile/file-transfer"},
				{Var: "http://jabber.org/protocol/si/profile/directory-transfer"},
				{Var: "http://jabber.org/protocol/si/profile/encrypted-data-transfer"},
				{Var: "http://jabber.org/protocol/bytestreams"},
				{Var: "urn:xmpp:eme:0"},
				{Var: "http://jabber.org/protocol/muc"},
			},
			Forms: []data.Form(nil)})
}

func (s *DiscoveryXMPPSuite) Test_conn_QueryServiceItems(c *C) {
	// See: XEP-0030, Section: 3.1 Basic Protocol, Example: 2
	fromServer := `
<iq xmlns='jabber:client' type='result'
    from='plays.shakespeare.lit'
    to='romeo@montague.net/orchard'
    id='%s'>
  <query xmlns='http://jabber.org/protocol/disco#items'>
    <item jid="foo@somewhere.com/bla"/>
    <item jid="bla@else.org" name="Something"/>
  </query>
</iq>
`
	mockOut := &mockConnIOReaderWriter{}

	conn := conn{
		log:       testLogger(),
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan bool, 1)
	go func() {
		q, e := conn.QueryServiceItems("plays.shakespeare.lit")
		c.Assert(e, IsNil)
		c.Assert(q, Not(IsNil))
		c.Assert(q.DiscoveryItems, HasLen, 2)
		done <- true
	}()

	waitForInflightTo(&conn, "plays.shakespeare.lit")

	c.Assert(string(mockOut.Written()), Matches, "<iq xmlns='jabber:client' to='plays.shakespeare.lit' from='romeo@montague.net/orchard' type='get' id='.+'><query xmlns=\"http://jabber.org/protocol/disco#items\"></query></iq>")

	var iq data.ClientIQ
	err := xml.Unmarshal(mockOut.Written(), &iq)
	c.Assert(err, IsNil)

	mockIn := &mockConnIOReaderWriter{read: []byte(fmt.Sprintf(fromServer, iq.ID))}
	conn.in = xml.NewDecoder(mockIn)

	_, err = conn.Next()
	c.Assert(err, Equals, io.EOF)
	<-done
}

func (s *DiscoveryXMPPSuite) Test_conn_QueryServiceItems_failsWriting(c *C) {
	mockOut := &mockConnIOReaderWriter{
		err: errors.New("an IO marker"),
	}

	conn := conn{
		log:       testLogger(),
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	q, e := conn.QueryServiceItems("plays.shakespeare.lit")
	c.Assert(e, ErrorMatches, "an IO marker")
	c.Assert(q, IsNil)
}

type resultData struct {
	data interface{}
	more *resultData
}

func (s *DiscoveryXMPPSuite) Test_conn_QueryServiceItems_noResponse(c *C) {
	mockOut := &mockConnIOReaderWriter{}

	conn := conn{
		log:       testLogger(),
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan *resultData, 1)
	go func() {
		q, e := conn.QueryServiceItems("plays.shakespeare.lit")
		done <- &resultData{q, &resultData{e, nil}}
	}()

	waitForInflightTo(&conn, "plays.shakespeare.lit")

	for _, inf := range conn.inflights {
		close(inf.replyChan)
	}

	rd := <-done

	q := rd.data.(*data.DiscoveryItemsQuery)
	e := rd.more.data.(error)

	c.Assert(e, ErrorMatches, "xmpp: failed to receive response")
	c.Assert(q, IsNil)
}

func (s *DiscoveryXMPPSuite) Test_conn_QueryServiceItems_badResponse(c *C) {
	mockOut := &mockConnIOReaderWriter{}

	conn := conn{
		log:       testLogger(),
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan *resultData, 1)
	go func() {
		q, e := conn.QueryServiceItems("plays.shakespeare.lit")
		done <- &resultData{q, &resultData{e, nil}}
	}()

	waitForInflightTo(&conn, "plays.shakespeare.lit")

	for _, inf := range conn.inflights {
		inf.replyChan <- data.Stanza{
			Value: "hello",
		}
	}

	rd := <-done

	q := rd.data.(*data.DiscoveryItemsQuery)
	e := rd.more.data.(error)

	c.Assert(e, ErrorMatches, "xmpp: failed to parse response")
	c.Assert(q, IsNil)
}

func (s *DiscoveryXMPPSuite) Test_conn_QueryServiceInformation_failsWriting(c *C) {
	mockOut := &mockConnIOReaderWriter{
		err: errors.New("an IO marker"),
	}

	conn := conn{
		log:       testLogger(),
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	q, e := conn.QueryServiceInformation("plays.shakespeare.lit")
	c.Assert(e, ErrorMatches, "an IO marker")
	c.Assert(q, IsNil)
}

func (s *DiscoveryXMPPSuite) Test_conn_QueryServiceInformation_noResponse(c *C) {
	mockOut := &mockConnIOReaderWriter{}

	conn := conn{
		log:       testLogger(),
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan *resultData, 1)
	go func() {
		q, e := conn.QueryServiceInformation("plays.shakespeare.lit")
		done <- &resultData{q, &resultData{e, nil}}
	}()

	waitForInflightTo(&conn, "plays.shakespeare.lit")

	for _, inf := range conn.inflights {
		close(inf.replyChan)
	}

	rd := <-done

	q := rd.data.(*data.DiscoveryInfoQuery)
	e := rd.more.data.(error)

	c.Assert(e, ErrorMatches, "xmpp: failed to receive response")
	c.Assert(q, IsNil)
}

func (s *DiscoveryXMPPSuite) Test_conn_QueryServiceInformation_badResponse(c *C) {
	mockOut := &mockConnIOReaderWriter{}

	conn := conn{
		log:       testLogger(),
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan *resultData, 1)
	go func() {
		q, e := conn.QueryServiceInformation("plays.shakespeare.lit")
		done <- &resultData{q, &resultData{e, nil}}
	}()

	waitForInflightTo(&conn, "plays.shakespeare.lit")

	for _, inf := range conn.inflights {
		inf.replyChan <- data.Stanza{
			Value: "hello",
		}
	}

	rd := <-done

	q := rd.data.(*data.DiscoveryInfoQuery)
	e := rd.more.data.(error)

	c.Assert(e, ErrorMatches, "xmpp: failed to parse response")
	c.Assert(q, IsNil)
}

func (s *DiscoveryXMPPSuite) Test_conn_EntityExists(c *C) {
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
		log:       testLogger(),
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan *resultData, 1)
	go func() {
		q, e := conn.EntityExists("plays.shakespeare.lit")
		done <- &resultData{q, &resultData{e, nil}}
	}()

	waitForInflightTo(&conn, "plays.shakespeare.lit")

	c.Assert(string(mockOut.Written()), Matches, "<iq xmlns='jabber:client' to='plays.shakespeare.lit' from='romeo@montague.net/orchard' type='get' id='.+'><query xmlns=\"http://jabber.org/protocol/disco#info\"></query></iq>")

	var iq data.ClientIQ
	err := xml.Unmarshal(mockOut.Written(), &iq)
	c.Assert(err, IsNil)

	mockIn := &mockConnIOReaderWriter{read: []byte(fmt.Sprintf(fromServer, iq.ID))}
	conn.in = xml.NewDecoder(mockIn)

	_, err = conn.Next()
	c.Assert(err, Equals, io.EOF)
	rd := <-done

	ok := rd.data.(bool)
	e := rd.more.data

	c.Assert(e, IsNil)
	c.Assert(ok, Equals, true)
}

func (s *DiscoveryXMPPSuite) Test_conn_EntityExists_fails(c *C) {
	mockOut := &mockConnIOReaderWriter{
		err: errors.New("an IO marker"),
	}

	conn := conn{
		log:       testLogger(),
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	ok, e := conn.EntityExists("plays.shakespeare.lit")

	c.Assert(e, ErrorMatches, "an IO marker")
	c.Assert(ok, Equals, false)
}

func (s *DiscoveryXMPPSuite) Test_conn_EntityExists_mucItemDoesntExist(c *C) {
	fromServer := `
<iq xmlns='jabber:client' type='error'
    from='plays.shakespeare.lit'
    to='romeo@montague.net/orchard'
    id='%s'>
  <error>
    <item-not-found xmlns="urn:ietf:params:xml:ns:xmpp-stanzas"/>
  </error>
</iq>
`
	mockOut := &mockConnIOReaderWriter{}

	conn := conn{
		log:       testLogger(),
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan *resultData, 1)
	go func() {
		q, e := conn.EntityExists("plays.shakespeare.lit")
		done <- &resultData{q, &resultData{e, nil}}
	}()

	waitForInflightTo(&conn, "plays.shakespeare.lit")

	c.Assert(string(mockOut.Written()), Matches, "<iq xmlns='jabber:client' to='plays.shakespeare.lit' from='romeo@montague.net/orchard' type='get' id='.+'><query xmlns=\"http://jabber.org/protocol/disco#info\"></query></iq>")

	var iq data.ClientIQ
	err := xml.Unmarshal(mockOut.Written(), &iq)
	c.Assert(err, IsNil)

	mockIn := &mockConnIOReaderWriter{read: []byte(fmt.Sprintf(fromServer, iq.ID))}
	conn.in = xml.NewDecoder(mockIn)

	_, err = conn.Next()
	c.Assert(err, Equals, io.EOF)
	rd := <-done

	ok := rd.data.(bool)
	e := rd.more.data

	c.Assert(e, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *DiscoveryXMPPSuite) Test_conn_DiscoveryFeaturesAndIdentities(c *C) {
	fromServer := `
<iq xmlns='jabber:client' type='result'
    from='plays.shakespeare.lit'
    to='romeo@montague.net/orchard'
    id='%s'>
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
    <feature var='http://jabber.org/protocol/muc'/>
    <feature var='jabber:iq:register'/>
    <feature var='jabber:iq:search'/>
    <feature var='jabber:iq:time'/>
    <feature var='jabber:iq:version'/>
  </query>
</iq>
`
	mockOut := &mockConnIOReaderWriter{}

	conn := conn{
		log:       testLogger(),
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan *resultData, 1)
	go func() {
		ids, feats, ok := conn.DiscoveryFeaturesAndIdentities("plays.shakespeare.lit")
		done <- &resultData{ids, &resultData{feats, &resultData{ok, nil}}}
	}()

	waitForInflightTo(&conn, "plays.shakespeare.lit")

	c.Assert(string(mockOut.Written()), Matches, "<iq xmlns='jabber:client' to='plays.shakespeare.lit' from='romeo@montague.net/orchard' type='get' id='.+'><query xmlns=\"http://jabber.org/protocol/disco#info\"></query></iq>")

	var iq data.ClientIQ
	err := xml.Unmarshal(mockOut.Written(), &iq)
	c.Assert(err, IsNil)

	mockIn := &mockConnIOReaderWriter{read: []byte(fmt.Sprintf(fromServer, iq.ID))}
	conn.in = xml.NewDecoder(mockIn)

	_, err = conn.Next()
	c.Assert(err, Equals, io.EOF)
	rd := <-done

	ids := rd.data.([]data.DiscoveryIdentity)
	feats := rd.more.data.([]string)
	ok := rd.more.more.data.(bool)

	c.Assert(ok, Equals, true)
	c.Assert(ids, HasLen, 2)
	c.Assert(feats, HasLen, 7)
}

func (s *DiscoveryXMPPSuite) Test_conn_DiscoveryFeaturesAndIdentities_fails(c *C) {
	mockOut := &mockConnIOReaderWriter{
		err: errors.New("an IO marker"),
	}

	conn := conn{
		log:       testLogger(),
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	_, _, ok := conn.DiscoveryFeaturesAndIdentities("plays.shakespeare.lit")

	c.Assert(ok, Equals, false)
}

func (s *DiscoveryXMPPSuite) Test_DiscoveryReply_forCustomNode(c *C) {
	res := DiscoveryReply("some@foo.bar", "chat")
	c.Assert(res, DeepEquals, data.ErrorReply{
		Type:  "cancel",
		Error: data.ErrorServiceUnavailable{},
	})
}

func (s *DiscoveryXMPPSuite) Test_conn_ServerHasFeature(c *C) {
	fromServer := `
<iq xmlns='jabber:client' type='result'
    from='montague.net'
    to='romeo@montague.net/orchard'
    id='%s'>
  <query xmlns='http://jabber.org/protocol/disco#info'>
    <feature var='http://jabber.org/protocol/disco#info'/>
    <feature var='http://jabber.org/protocol/disco#items'/>
    <feature var='http://jabber.org/protocol/muc'/>
    <feature var='jabber:iq:register'/>
    <feature var='jabber:iq:search'/>
    <feature var='jabber:iq:time'/>
    <feature var='jabber:iq:version'/>
  </query>
</iq>
`
	mockOut := &mockConnIOReaderWriter{}

	conn := conn{
		log:       testLogger(),
		in:        xml.NewDecoder(nil),
		out:       mockOut,
		inflights: make(map[data.Cookie]inflight),

		jid: "romeo@montague.net/orchard",
	}

	done := make(chan bool, 1)
	go func() {
		res := conn.ServerHasFeature("trombones")
		done <- res
	}()

	waitForInflightTo(&conn, "plays.shakespeare.lit")

	c.Assert(string(mockOut.Written()), Matches, "<iq xmlns='jabber:client' from='romeo@montague.net/orchard' type='get' id='.+'><query xmlns=\"http://jabber.org/protocol/disco#info\"></query></iq>")

	var iq data.ClientIQ
	err := xml.Unmarshal(mockOut.Written(), &iq)
	c.Assert(err, IsNil)

	mockIn := &mockConnIOReaderWriter{read: []byte(fmt.Sprintf(fromServer, iq.ID))}
	conn.in = xml.NewDecoder(mockIn)

	_, err = conn.Next()
	c.Assert(err, Equals, io.EOF)
	res := <-done

	c.Assert(res, Equals, false)
	c.Assert(conn.ServerHasFeature("jabber:iq:time"), Equals, true)
}
