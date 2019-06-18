package xmpp

import (
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
	conn.out = mockIn
	conn.jid = "user@xmpp.org"

	_, _, err := conn.CancelRegistration()
	c.Assert(err, IsNil)
	c.Assert(string(mockIn.write), Matches, expectedoOut)
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

// Test for ChangePassword, when a successful request is made.
// XXX: Stubbed implementation
func (s *RegisterSuite) Test_ChangePassword_SuccessfulChangePasswordRequest(c *C) {
	expectedoOut := ""

	mockIn := &mockConnIOReaderWriter{}
	conn := newConn()
	conn.out = mockIn
	conn.jid = "user@xmpp.org"
	newPassword := "123"

	status, err := conn.ChangePassword(conn.jid, newPassword)

	c.Assert(err, IsNil)
	c.Assert(status, Equals, true)
	c.Assert(string(mockIn.write), Matches, expectedoOut)
}
