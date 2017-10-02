// +build go1.8

package xmpp

import (
	"crypto/tls"
	"io"

	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

func (s *ConnectionXMPPSuite) Test_Dial_failsWhenStartingAHandshake(c *C) {
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

	expectedXmppHeader := "<?xml version='1.0'?>" +
		"<stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n" +
		"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>"
	expected := []byte{
		0x16, 0x3, 0x1, 0x0, 0x94, 0x1, 0x0, 0x0,
		0x90, 0x3, 0x3, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc,
		0xd, 0xe, 0xf, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c,
		0x1d, 0x1e, 0x1f, 0x0, 0x0, 0x20, 0xcc, 0xa8, 0xcc, 0xa9, 0xc0, 0x2f, 0xc0, 0x30, 0xc0, 0x2b,
		0xc0, 0x2c, 0xc0, 0x13, 0xc0, 0x9, 0xc0, 0x14, 0xc0, 0xa, 0x0, 0x9c, 0x0, 0x9d, 0x0, 0x2f,
		0x0, 0x35, 0xc0, 0x12, 0x0, 0xa, 0x1, 0x0, 0x0, 0x47, 0x0, 0x0, 0x0, 0xb, 0x0, 0x9,
		0x0, 0x0, 0x6, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x0, 0x5, 0x0, 0x5, 0x1, 0x0, 0x0,
		0x0, 0x0, 0x0, 0xa, 0x0, 0xa, 0x0, 0x8, 0x0, 0x1d, 0x0, 0x17, 0x0, 0x18, 0x0, 0x19,
		0x0, 0xb, 0x0, 0x2, 0x1, 0x0, 0x0, 0xd, 0x0, 0xe, 0x0, 0xc, 0x4, 0x1, 0x4, 0x3,
		0x5, 0x1, 0x5, 0x3, 0x2, 0x1, 0x2, 0x3, 0xff, 0x1, 0x0, 0x1, 0x0, 0x0, 0x12, 0x0, 0x0,
	}

	c.Assert(len(rw.write), Equals, len(expectedXmppHeader)+len(expected))
	c.Assert(string(rw.write[:len(expectedXmppHeader)]), Equals, expectedXmppHeader)
	c.Assert(rw.write[len(expectedXmppHeader):], DeepEquals, expected)
}

func (s *ConnectionXMPPSuite) Test_Dial_worksIfTheHandshakeSucceeds(c *C) {
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

	c.Assert(err.Error(), Equals, "tls: server's Finished message was incorrect")
	c.Assert(string(rw.write), Equals,
		"<?xml version='1.0'?><stream:stream to='www.olabini.se' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>\x16\x03\x01\x00\x9c\x01\x00\x00\x98\x03\x03\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
			"\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x00\x00 ̨̩\xc0/\xc00\xc0+\xc0,\xc0\x13\xc0\t\xc0\x14\xc0\n"+
			"\x00\x9c\x00\x9d\x00/\x005\xc0\x12\x00\n"+
			"\x01\x00\x00O\x00\x00\x00\x13\x00\x11\x00\x00\x0ewww.olabini.se\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\n"+
			"\x00\n"+
			"\x00\b\x00\x1d\x00\x17\x00\x18\x00\x19\x00\v\x00\x02\x01\x00\x00\r\x00\x0e\x00\f\x04\x01\x04\x03\x05\x01\x05\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00\x00\x12\x00\x00\x16\x03\x03\x00F\x10\x00\x00BA\x04\b\xda\xf8\xb0\xab\xae5\t\xf3\\\xe1\xd31\x04\xcb\x01\xb9Qb̹\x18\xba\x1f\x81o8\xd3\x13\x0f\xb8\u007f\x92\xa3\b7\xf8o\x9e\xef\x19\u007fCy\xa5\n"+
			"b\x06\x82fy]\xb9\xf83\xea6\x1d\x03\xafT[\xe7\x92\x14\x03\x03\x00\x01\x01\x16\x03\x03\x00(\x00\x00\x00\x00\x00\x00\x00\x00j\xc9\xd4\xef\v.~\xb5?\xdaD#q\xed\xf9?\a\xb8\xcfF\x9c\xe8\xe8\xdc\xfaE\xe1\xc80\xf4S5\x15\x03\x03\x00\x1a\x00\x00\x00\x00\x00\x00\x00\x01i01i\x98I\x1a\x9e\x13NL\x9b3\a\x10\xb9\xe0\xa5",
	)
}

func (s *ConnectionXMPPSuite) Test_Dial_worksIfTheHandshakeSucceedsButFailsOnInvalidCertHash(c *C) {
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

	c.Assert(err.Error(), Equals, "tls: server's Finished message was incorrect")
}

func (s *ConnectionXMPPSuite) Test_Dial_worksIfTheHandshakeSucceedsButSucceedsOnValidCertHash(c *C) {
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
		verifier:      &basicTLSVerifier{bytesFromHex("82454418cb04854aa721bb0596528ff802b1e18a4e3a7767412ac9f108c9d3a7")},

		config: data.Config{
			TLSConfig: &tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "tls: server's Finished message was incorrect")
}
