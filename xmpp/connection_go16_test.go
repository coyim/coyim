// +build go1.6

package xmpp

import (
	"crypto/tls"
	"io"

	. "github.com/twstrike/coyim/Godeps/_workspace/src/gopkg.in/check.v1"
	"github.com/twstrike/coyim/xmpp/data"
)

func (s *ConnectionXmppSuite) Test_Dial_failsWhenStartingAHandshake(c *C) {
	rw := &mockConnIOReaderWriter{read: []byte(
		"<?xml version='1.0'?>" +
			"<str:stream xmlns:str='http://etherx.jabber.org/streams' version='1.0'>" +
			"<str:features>" +
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'>" +
			"</starttls>" +
			"<mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" +
			"<mechanism>PLAIN</mechanism>" +
			"</mechanisms>" +
			"</str:features>" +
			"<proceed xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>",
	)}
	conn := &fullMockedConn{rw: rw}
	var tlsC tls.Config
	tlsC.Rand = fixedRand([]string{"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F"})

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
			TLSConfig: &tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?><stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>\x16\x03\x01\x00\x8e\x01\x00\x00\x8a\x03\x03\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
		"\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x00\x00\x1c\xc0/\xc0+\xc00\xc0,\xc0\x13\xc0\t\xc0\x14\xc0\n"+
		"\x00\x9c\x00\x9d\x00/\x005\xc0\x12\x00\n"+
		"\x01\x00\x00E\x00\x00\x00\v\x00\t\x00\x00\x06domain\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\n"+
		"\x00\b\x00\x06\x00\x17\x00\x18\x00\x19\x00\v\x00\x02\x01\x00\x00\r\x00\x0e\x00\f\x04\x01\x04\x03\x05\x01\x05\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00\x00\x12\x00\x00",
	)
}

func (s *ConnectionXmppSuite) Test_Dial_worksIfTheHandshakeSucceeds(c *C) {
	rw := &mockMultiConnIOReaderWriter{read: validTLSExchange}
	conn := &fullMockedConn{rw: rw}
	var tlsC tls.Config
	tlsC.Rand = fixedRand([]string{
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F",
		"000102030405060708090A0B0C0D0E0F",
	})

	d := &dialer{
		JID:           "user@www.olabini.se",
		password:      "pass",
		serverAddress: "www.olabini.se:443",
		verifier:      &basicTLSVerifier{},

		config: data.Config{
			TLSConfig: &tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
	c.Assert(string(rw.write), Equals, ""+
		"<?xml version='1.0'?><stream:stream to='www.olabini.se' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
		"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>\x16\x03\x01\x00\x96\x01\x00\x00\x92\x03\x03\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
		"\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x00\x00\x1c\xc0/\xc0+\xc00\xc0,\xc0\x13\xc0\t\xc0\x14\xc0\n"+
		"\x00\x9c\x00\x9d\x00/\x005\xc0\x12\x00\n"+
		"\x01\x00\x00M\x00\x00\x00\x13\x00\x11\x00\x00\x0ewww.olabini.se\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\n"+
		"\x00\b\x00\x06\x00\x17\x00\x18\x00\x19\x00\v\x00\x02\x01\x00\x00\r\x00\x0e\x00\f\x04\x01\x04\x03\x05\x01\x05\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00\x00\x12\x00\x00\x16\x03\x03\x00F\x10\x00\x00BA\x04\b\xda\xf8\xb0\xab\xae5\t\xf3\\\xe1\xd31\x04\xcb\x01\xb9Qb̹\x18\xba\x1f\x81o8\xd3\x13\x0f\xb8\u007f\x92\xa3\b7\xf8o\x9e\xef\x19\u007fCy\xa5\n"+
		"b\x06\x82fy]\xb9\xf83\xea6\x1d\x03\xafT[\xe7\x92\x14\x03\x03\x00\x01\x01\x16\x03\x03\x00(\x00\x00\x00\x00\x00\x00\x00\x00K\xe8\x9fz\xf8d\xf8&<=\xe4\xae\xd3\\A\xbb\u007f\x8f<\xfc-\x9f\x11\x9f\xd5\x01\xe5S\x86:\xed\x99\x17\x03\x03\x00\xa5\x00\x00\x00\x00\x00\x00\x00\x01v@*\x1bʱ\x05\xa7\x92+@}\xb8\x02p\xaeF\x17\x99\xb5\x01\x19\x05\xc4r\xc6v1]\xf7\x9b\xf6\x1c\xf6\xb2*\xe4\f\x93\xbd4-~z\x8a\xb9S\x14\xacnm\xe5c\xcff\xaf\x83s~\xfcJ9j5*!\x9e\xbe\x00\x1d2\x80\xf1(\xf4\xcff\xd0tq\xb2U\x9a(\x8bAl\xac\xc1\xf0\x9a\xd4\u007f\x0fR\\,\xf9Ρע\xc2\f\xe6\xa8г\x98O\x05\xdf\xd2(\b\xf8v\x8b\xfdݲ\xea\x18kȄ\xbd6 \xbe\xd9\xe1\xc5&\xd2Ň\xde\xe1ği\xa6\xee\xf4r\x8c\x13\xbd\x97֮\x87\xbd\x92F\xf2")

}

func (s *ConnectionXmppSuite) Test_Dial_worksIfTheHandshakeSucceedsButFailsOnInvalidCertHash(c *C) {
	rw := &mockMultiConnIOReaderWriter{read: validTLSExchange}
	conn := &fullMockedConn{rw: rw}
	var tlsC tls.Config
	tlsC.Rand = fixedRand([]string{
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F",
		"000102030405060708090A0B0C0D0E0F",
	})

	d := &dialer{
		JID:           "user@www.olabini.se",
		password:      "pass",
		serverAddress: "www.olabini.se:443",
		verifier:      &basicTLSVerifier{[]byte("aaaaa")},

		config: data.Config{
			TLSConfig: &tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "tls: server certificate does not match expected hash (got: 2300818fdc977ce5eb357694d421e47869a952990bc3230ef6aca2bb6ee6f00b, want: 6161616161)")
}

func (s *ConnectionXmppSuite) Test_Dial_worksIfTheHandshakeSucceedsButSucceedsOnValidCertHash(c *C) {
	rw := &mockMultiConnIOReaderWriter{read: validTLSExchange}
	conn := &fullMockedConn{rw: rw}
	var tlsC tls.Config
	tlsC.Rand = fixedRand([]string{
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F",
		"000102030405060708090A0B0C0D0E0F",
	})

	d := &dialer{
		JID:           "user@www.olabini.se",
		password:      "pass",
		serverAddress: "www.olabini.se:443",
		verifier:      &basicTLSVerifier{bytesFromHex("2300818fdc977ce5eb357694d421e47869a952990bc3230ef6aca2bb6ee6f00b")},

		config: data.Config{
			TLSConfig: &tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
}
