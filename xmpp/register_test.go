package xmpp

import (
	"bytes"

	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

type RegisterSuite struct{}

var _ = Suite(&RegisterSuite{})

func (s *RegisterSuite) Test_CancelRegistration_SendCancelationRequest(c *C) {
	expectedoOut := "<iq xmlns='jabber:client'  from='user@xmpp.org' type='set' id='.+'>\n" +
		"\t<query xmlns='jabber:iq:register'>\n" +
		"\t\t<remove/>\n" +
		"\t</query>\n" +
		"\t</iq>"

	mockIn := &mockConnIOReaderWriter{}
	conn := newConn()
	conn.log = testLogger()
	conn.out = mockIn
	conn.jid = "user@xmpp.org"

	_, _, err := conn.CancelRegistration()
	c.Assert(err, IsNil)
	c.Assert(string(mockIn.write), Matches, expectedoOut)
}

func (s *RegisterSuite) Test_SendChangePasswordInfo(c *C) {
	mockOut := &mockConnIOReaderWriter{}
	conn := conn{
		log:  testLogger(),
		out:  mockOut,
		jid:  "crone1@shakespeare.lit",
		rand: bytes.NewBuffer([]byte{1, 0, 0, 0, 0, 0, 0, 0}),
	}

	conn.inflights = make(map[data.Cookie]inflight)

	reply, cookie, err := conn.sendChangePasswordInfo("crone1", "shakespeare.lit", "pass")

	c.Assert(err, IsNil)
	c.Assert(string(mockOut.write), Matches, "<iq xmlns='jabber:client' to='shakespeare.lit' from='crone1@shakespeare.lit' type='set' id='1'><query xmlns='jabber:iq:register'><username>crone1</username><password>pass</password></query></iq>")
	c.Assert(reply, NotNil)
	c.Assert(cookie, NotNil)
}

func (s *RegisterSuite) Test_setupStream_registerWithoutAuthenticating(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"<register xmlns='http://jabber.org/features/iq-register'/>" +
			"</str:features>" +
			"<iq xmlns='jabber:client' type='result'>" +
			"<query xmlns='jabber:iq:register'><username/></query>" +
			"</iq>" +
			"<iq xmlns='jabber:client' type='result'></iq>",
	)}
	conn := &fullMockedConn{rw: rw}

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
			SkipTLS: true,
			CreateCallback: func(title, instructions string, fields []interface{}) error {
				return nil
			},
		},
		log: testLogger(),
	}
	_, err := d.setupStream(conn)

	c.Assert(err, IsNil)
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>"+
		"</stream:stream>",
	)
}
