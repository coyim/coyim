// +build go1.8

package xmpp

import (
	"crypto/tls"
	"io"

	"github.com/twstrike/coyim/xmpp/data"
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
	c.Assert(string(rw.write), Equals,
		"<?xml version='1.0'?><stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>\x16\x03\x01\x00\x94\x01\x00\x00\x90\x03\x03\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
			"\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x00\x00 \xc0/\xc00\xc0+\xc0,̨̩\xc0\x13\xc0\t\xc0\x14\xc0\n"+
			"\x00\x9c\x00\x9d\x00/\x005\xc0\x12\x00\n"+
			"\x01\x00\x00G\x00\x00\x00\v\x00\t\x00\x00\x06domain\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\n"+
			"\x00\n"+
			"\x00\b\x00\x1d\x00\x17\x00\x18\x00\x19\x00\v\x00\x02\x01\x00\x00\r\x00\x0e\x00\f\x04\x01\x04\x03\x05\x01\x05\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00\x00\x12\x00\x00",
	)
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
			"\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x00\x00 \xc0/\xc00\xc0+\xc0,̨̩\xc0\x13\xc0\t\xc0\x14\xc0\n"+
			"\x00\x9c\x00\x9d\x00/\x005\xc0\x12\x00\n"+
			"\x01\x00\x00O\x00\x00\x00\x13\x00\x11\x00\x00\x0ewww.olabini.se\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\n"+
			"\x00\n"+
			"\x00\b\x00\x1d\x00\x17\x00\x18\x00\x19\x00\v\x00\x02\x01\x00\x00\r\x00\x0e\x00\f\x04\x01\x04\x03\x05\x01\x05\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00\x00\x12\x00\x00\x16\x03\x03\x00F\x10\x00\x00BA\x04\b\xda\xf8\xb0\xab\xae5\t\xf3\\\xe1\xd31\x04\xcb\x01\xb9Qb̹\x18\xba\x1f\x81o8\xd3\x13\x0f\xb8\u007f\x92\xa3\b7\xf8o\x9e\xef\x19\u007fCy\xa5\n"+
			"b\x06\x82fy]\xb9\xf83\xea6\x1d\x03\xafT[\xe7\x92\x14\x03\x03\x00\x01\x01\x16\x03\x03\x00(\x00\x00\x00\x00\x00\x00\x00\x00j\xc9\xd4\uf8b3\xe9h\x83P+\x8c\xe2\b\"\x9d\xcc\xe6Ar\xbe\xa1\x9f\x96\xfe\xb63X'\x8a5\x10\x15\x03\x03\x00\x1a\x00\x00\x00\x00\x00\x00\x00\x01i01i\x98I\x1a\x9e\x13NL\x9b3\a\x10\xb9\xe0\xa5",
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
