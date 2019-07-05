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

// TODO: move all of these tests
// Tests ChangePassword() for a successful password change.
func (s *RegisterSuite) Test_ChangePassword_SuccessfulChangePasswordRequest(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"<register xmlns='http://jabber.org/features/iq-register'/>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>" +
			"<?xml version='1.0'?>" +
			"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>" +
			"<stream:features />" +
			"<iq type='set' to='domain' id='change1'>" +
			"<query xmlns='jabber:iq:register'>" +
			"<username>user</username>" +
			"<password>newpass</password>" +
			"</query>" +
			"</iq>" +
			"<iq xmlns='jabber:client' type='result'></iq>",
	)}

	conn := &fullMockedConn{rw: rw}

	d := &dialer{
		JID:         "user@domain",
		password:    "pass",
		newPassword: "newPass",
		config: data.Config{
			SkipTLS: true,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, IsNil)
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq><iq type='set' to='domain' id='change1'><query xmlns='jabber:iq:register'><username>user</username><password>newPass</password></query></iq>"+
		"</stream:stream>",
	)
}

// Tests ChangePassword() with an XML that has an error of type "bad-request"
func (s *RegisterSuite) Test_ChangePassword_FailedBadRequest(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"<register xmlns='http://jabber.org/features/iq-register'/>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>" +
			"<?xml version='1.0'?>" +
			"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>" +
			"<stream:features />" +
			"<iq type='set' to='domain' id='change1'>" +
			"<query xmlns='jabber:iq:register'>" +
			"<username>user</username>" +
			"<password>newpass</password>" +
			"</query>" +
			"</iq>" +
			"<iq type='error' from='domain' to='user@domain' id='change1'>" +
			"<error code='400' type='modify'>" +
			"<bad-request xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>" +
			"</error></iq>" +
			"<iq type='set' to='domain' id='change1'>",
	)}

	conn := &fullMockedConn{rw: rw}

	d := &dialer{
		JID:         "user@domain",
		password:    "pass",
		newPassword: "newPass",
		config: data.Config{
			SkipTLS: true,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "xmpp: password change request was malformed")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq><iq type='set' to='domain' id='change1'><query xmlns='jabber:iq:register'><username>user</username><password>newPass</password></query></iq>",
	)
}

// Tests ChangePassword() with an XML that has an error of type "not-authorized"
func (s *RegisterSuite) Test_ChangePassword_FailedNotAuthorized(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"<register xmlns='http://jabber.org/features/iq-register'/>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>" +
			"<?xml version='1.0'?>" +
			"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>" +
			"<stream:features />" +
			"<iq type='set' to='domain' id='change1'>" +
			"<query xmlns='jabber:iq:register'>" +
			"<username>user</username>" +
			"<password>newpass</password>" +
			"</query>" +
			"</iq>" +
			"<iq type='error' from='domain' to='user@domain' id='change1'>" +
			"<error code='401' type='modify'>" +
			"<not-authorized xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>" +
			"</error></iq>",
	)}

	conn := &fullMockedConn{rw: rw}

	d := &dialer{
		JID:         "user@domain",
		password:    "pass",
		newPassword: "newPass",
		config: data.Config{
			SkipTLS: true,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "xmpp: password change not authorized")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq><iq type='set' to='domain' id='change1'><query xmlns='jabber:iq:register'><username>user</username><password>newPass</password></query></iq>",
	)
}

// Tests ChangePassword() with an XML that has an error of type "not-allowed"
func (s *RegisterSuite) Test_ChangePassword_FailedNotAllowed(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"<register xmlns='http://jabber.org/features/iq-register'/>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>" +
			"<?xml version='1.0'?>" +
			"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>" +
			"<stream:features />" +
			"<iq type='set' to='domain' id='change1'>" +
			"<query xmlns='jabber:iq:register'>" +
			"<username>user</username>" +
			"<password>newpass</password>" +
			"</query>" +
			"</iq>" +
			"<iq type='error' from='domain' to='user@domain' id='change1'>" +
			"<error code='405' type='modify'>" +
			"<not-allowed xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>" +
			"</error></iq>",
	)}

	conn := &fullMockedConn{rw: rw}

	d := &dialer{
		JID:         "user@domain",
		password:    "pass",
		newPassword: "newPass",
		config: data.Config{
			SkipTLS: true,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "xmpp: server does not allow password changes")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq><iq type='set' to='domain' id='change1'><query xmlns='jabber:iq:register'><username>user</username><password>newPass</password></query></iq>",
	)
}

// Tests ChangePassword() with an XML that has an error of type "unexpected-request"
func (s *RegisterSuite) Test_ChangePassword_FailedUnexpectedRequest(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"<register xmlns='http://jabber.org/features/iq-register'/>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>" +
			"<?xml version='1.0'?>" +
			"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>" +
			"<stream:features />" +
			"<iq type='set' to='domain' id='change1'>" +
			"<query xmlns='jabber:iq:register'>" +
			"<username>user</username>" +
			"<password>newpass</password>" +
			"</query>" +
			"</iq>" +
			"<iq type='error' from='domain' to='user@domain' id='change1'>" +
			"<error code='400' type='modify'>" +
			"<unexpected-request xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>" +
			"</error></iq>",
	)}

	conn := &fullMockedConn{rw: rw}

	d := &dialer{
		JID:         "user@domain",
		password:    "pass",
		newPassword: "newPass",
		config: data.Config{
			SkipTLS: true,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "xmpp: user is not registered with server")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq><iq type='set' to='domain' id='change1'><query xmlns='jabber:iq:register'><username>user</username><password>newPass</password></query></iq>",
	)
}

// Tests ChangePassword() with an XML that has an error not specificed in XEP-0077
func (s *RegisterSuite) Test_ChangePassword_FailedGeneric(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"<register xmlns='http://jabber.org/features/iq-register'/>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>" +
			"<?xml version='1.0'?>" +
			"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>" +
			"<stream:features />" +
			"<iq type='set' to='domain' id='change1'>" +
			"<query xmlns='jabber:iq:register'>" +
			"<username>user</username>" +
			"<password>newpass</password>" +
			"</query>" +
			"</iq>" +
			"<iq type='error' from='domain' to='user@domain' id='change1'>" +
			"<error code='400' type='modify'>" +
			"<this-will-not-match xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>" +
			"</error></iq>",
	)}

	conn := &fullMockedConn{rw: rw}

	d := &dialer{
		JID:         "user@domain",
		password:    "pass",
		newPassword: "newPass",
		config: data.Config{
			SkipTLS: true,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "xmpp: password change failed")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq><iq type='set' to='domain' id='change1'><query xmlns='jabber:iq:register'><username>user</username><password>newPass</password></query></iq>",
	)
}

// Tests ChangePassword() with an XML that faces Unmarshalling error
func (s *RegisterSuite) Test_ChangePassword_FailedUnmarshal(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"<register xmlns='http://jabber.org/features/iq-register'/>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>" +
			"<?xml version='1.0'?>" +
			"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>" +
			"<stream:features />" +
			"<iq type='set' to='domain' id='change1'>" +
			"<query xmlns='jabber:iq:register'>" +
			"<username>user</username>" +
			"<password>newpass</password>" +
			"</query>" +
			"</iq>",
	)}

	conn := &fullMockedConn{rw: rw}

	d := &dialer{
		JID:         "user@domain",
		password:    "pass",
		newPassword: "newPass",
		config: data.Config{
			SkipTLS: true,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "unmarshal <iq>: XML syntax error on line 1: unexpected EOF")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq><iq type='set' to='domain' id='change1'><query xmlns='jabber:iq:register'><username>user</username><password>newPass</password></query></iq>",
	)
}

// Tests ChangePassword() with an XML that does not have any matching type
func (s *RegisterSuite) Test_ChangePassword_FailedNoMatch(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"<register xmlns='http://jabber.org/features/iq-register'/>" +
			"</str:features>" +
			"<sasl:success xmlns:sasl='urn:ietf:params:xml:ns:xmpp-sasl'></sasl:success>" +
			"<?xml version='1.0'?>" +
			"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>" +
			"<stream:features />" +
			"<iq type='set' to='domain' id='change1'>" +
			"<query xmlns='jabber:iq:register'>" +
			"<username>user</username>" +
			"<password>newpass</password>" +
			"</query>" +
			"</iq>" +
			"<iq xmlns='jabber:client' type='no-match'></iq>",
	)}

	conn := &fullMockedConn{rw: rw}

	d := &dialer{
		JID:         "user@domain",
		password:    "pass",
		newPassword: "newPass",
		config: data.Config{
			SkipTLS: true,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "xmpp: password change failed")
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<auth xmlns='urn:ietf:params:xml:ns:xmpp-sasl' mechanism='PLAIN'>AHVzZXIAcGFzcw==</auth>\n"+
		"<?xml version='1.0'?>"+
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<iq type='set' id='bind_1'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'/></iq><iq type='set' to='domain' id='change1'><query xmlns='jabber:iq:register'><username>user</username><password>newPass</password></query></iq>",
	)
}
