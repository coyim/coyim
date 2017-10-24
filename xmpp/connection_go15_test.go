// +build !go1.6

package xmpp

import (
	"crypto/tls"
	"io"

	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
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
	tlsC := &tls.Config{
		SessionTicketKey: [32]byte{1},
		Rand:             fixedRand([]string{"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F"}),
	}

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
			TLSConfig: tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
	if isVersionOldish() {
		c.Assert(string(rw.write), Equals, ""+
			"<?xml version='1.0'?><stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>\x16\x03\x01\x00\x84\x01\x00\x00\x80\x03\x03\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
			"\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x00\x00\x1a\xc0/\xc0+\xc0\x11\xc0\a\xc0\x13\xc0\t\xc0\x14\xc0\n"+
			"\x00\x05\x00/\x005\xc0\x12\x00\n"+
			"\x01\x00\x00=\x00\x00\x00\v\x00\t\x00\x00\x06domain\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\n"+
			"\x00\b\x00\x06\x00\x17\x00\x18\x00\x19\x00\v\x00\x02\x01\x00\x00\r\x00\n"+
			"\x00\b\x04\x01\x04\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00",
		)
	} else {
		c.Assert(string(rw.write), Equals, ""+
			"<?xml version='1.0'?><stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>\x16\x03\x01\x00\x8a\x01\x00\x00\x86\x03\x03\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
			"\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x00\x00\x18\xc0/\xc0+\xc00\xc0,\xc0\x13\xc0\t\xc0\x14\xc0\n"+
			"\x00/\x005\xc0\x12\x00\n"+
			"\x01\x00\x00E\x00\x00\x00\v\x00\t\x00\x00\x06domain\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\n"+
			"\x00\b\x00\x06\x00\x17\x00\x18\x00\x19\x00\v\x00\x02\x01\x00\x00\r\x00\x0e\x00\f\x04\x01\x04\x03\x05\x01\x05\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00\x00\x12\x00\x00",
		)
	}
}

func (s *ConnectionXmppSuite) Test_Dial_worksIfTheHandshakeSucceeds(c *C) {
	rw := &mockMultiConnIOReaderWriter{read: validTLSExchange}
	conn := &fullMockedConn{rw: rw}
	tlsC := &tls.Config{
		SessionTicketKey: [32]byte{1},
		Rand: fixedRand([]string{
			"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
			"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
			"000102030405060708090A0B0C0D0E0F",
			"000102030405060708090A0B0C0D0E0F",
		}),
	}

	d := &dialer{
		JID:           "user@www.olabini.se",
		password:      "pass",
		serverAddress: "www.olabini.se:443",
		verifier:      &basicTLSVerifier{},

		config: data.Config{
			TLSConfig: tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
	if isVersionOldish() {
		c.Assert(string(rw.write), Equals, ""+
			"<?xml version='1.0'?><stream:stream to='www.olabini.se' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>\x16\x03\x01\x00\x8c\x01\x00\x00\x88\x03\x03\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
			"\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x00\x00\x1a\xc0/\xc0+\xc0\x11\xc0\a\xc0\x13\xc0\t\xc0\x14\xc0\n"+
			"\x00\x05\x00/\x005\xc0\x12\x00\n"+
			"\x01\x00\x00E\x00\x00\x00\x13\x00\x11\x00\x00\x0ewww.olabini.se\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\n"+
			"\x00\b\x00\x06\x00\x17\x00\x18\x00\x19\x00\v\x00\x02\x01\x00\x00\r\x00\n"+
			"\x00\b\x04\x01\x04\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00\x16\x03\x03\x00F\x10\x00\x00BA\x04\b\xda\xf8\xb0\xab\xae5\t\xf3\\\xe1\xd31\x04\xcb\x01\xb9Qb̹\x18\xba\x1f\x81o8\xd3\x13\x0f\xb8\u007f\x92\xa3\b7\xf8o\x9e\xef\x19\u007fCy\xa5\n"+
			"b\x06\x82fy]\xb9\xf83\xea6\x1d\x03\xafT[\xe7\x92\x14\x03\x03\x00\x01\x01\x16\x03\x03\x00@\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
			"\v\f\r\x0e\x0fL\x86\xb5⌍\xfe\xc77\x8b\x01\x18\xces\xf4\x01S\xcbI9\t \x9e1\xe6U\n"+
			"\xff\xa67\xe4,Z\x05\x9e\xcfՈ\xfd-ڰ\x9dn\xac[Ud\x17\x03\x03\x00\xc0\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
			"\v\f\r\x0e\x0f\xd6|\xd9h1\xb4}ȹ\x80\xe7]\xbf\b\xe5Ю~N\x11n(Ӯ`\xcd\xfb\xc4E1Ѧ\xa0d\x83#?\b\xde\x1bp5@\x94\xe8\x89\xc0\xdb.\x03pd\xa6&|\xa2\xd8\f6\x91\a\xeb6G\xe5\xe1@\x02a\x89\x95@\x81\x0e\x161Dy\xf1N\xbf4\xf1\x93\xa4\xd2\xfc\xb6JZi\xb5\b\x13\xb7{n\\\xd6\x15M!\xe5\x10\xdc\x15\xbb\xab\xc0'\x11\xf6\xb8\xfa\x82I\x18\x96\xe1>\xa1\xa6GT\x1cy\x94X\xa7\xfa\x9c\xf5\xa8\xe8\t\xc8\xf5I\xce\xd9\xd8<\xf8\x9d\x93\xf8\x9b\xdd8\xbbIdVP\x0f\x88\x05{\x9b\x84Z\x82\x91\x9e4\x91}\x97nZ\xa2n\xfe.#?")
	} else {
		c.Assert(string(rw.write), Equals, ""+
			"<?xml version='1.0'?><stream:stream to='www.olabini.se' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n"+
			"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>\x16\x03\x01\x00\x92\x01\x00\x00\x8e\x03\x03\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
			"\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x00\x00\x18\xc0/\xc0+\xc00\xc0,\xc0\x13\xc0\t\xc0\x14\xc0\n"+
			"\x00/\x005\xc0\x12\x00\n"+
			"\x01\x00\x00M\x00\x00\x00\x13\x00\x11\x00\x00\x0ewww.olabini.se\x00\x05\x00\x05\x01\x00\x00\x00\x00\x00\n"+
			"\x00\b\x00\x06\x00\x17\x00\x18\x00\x19\x00\v\x00\x02\x01\x00\x00\r\x00\x0e\x00\f\x04\x01\x04\x03\x05\x01\x05\x03\x02\x01\x02\x03\xff\x01\x00\x01\x00\x00\x12\x00\x00\x16\x03\x03\x00F\x10\x00\x00BA\x04\b\xda\xf8\xb0\xab\xae5\t\xf3\\\xe1\xd31\x04\xcb\x01\xb9Qb̹\x18\xba\x1f\x81o8\xd3\x13\x0f\xb8\u007f\x92\xa3\b7\xf8o\x9e\xef\x19\u007fCy\xa5\n"+
			"b\x06\x82fy]\xb9\xf83\xea6\x1d\x03\xafT[\xe7\x92\x14\x03\x03\x00\x01\x01\x16\x03\x03\x00(\x00\x00\x00\x00\x00\x00\x00\x00~nq\xd7\x04\\\x97\xdfR\f\xd5Aѥ\xf3\t>\x81P\xc0\xffO\xc2\xd9`\xf3\u0094\xb3\x14\x04\x90\x17\x03\x03\x00\xa5\x00\x00\x00\x00\x00\x00\x00\x01\x1c3\xb8T0zp\x01\xe2P\xffI\xec\x1d a^w\x94\xf2+\xd6\xc7\xd8\x10\x108`\x9f\x98h(\x1f\xf0\xd5\xd0\U000163fb\x9d\x98\xd5s\xc36\xe6\x9bR\x88\x9d\xb7/#\x18^!1\x06\x90;#\x04\u03a2\x9djf\xd5\xd8:\x05+\x10T\x1a\x00\xd1\xc1\xe6\xd9V\xeb\xb6p+\xfa\xa3`\x92\xc0Z\xd3\xe6}K;\x02\xfc\xc9\u0530\x0f\xbd\xeeX\xf8\xf7\xdf\x16\x05q9\xecat\x80\x97\x80\xa0\x1fŊ\xa7\n"+
			"\x9e\xfaڣ\xae\xeb\xc0k\n"+
			"\x14\"\xb1N\x89{\x90z\xaf\xbeCP\xfd\x98\x9b\x13v\x14\xe9\xd1\x1f\x9e\x96")
	}
}

func (s *ConnectionXmppSuite) Test_Dial_worksIfTheHandshakeSucceedsButFailsOnInvalidCertHash(c *C) {
	rw := &mockMultiConnIOReaderWriter{read: validTLSExchange}
	conn := &fullMockedConn{rw: rw}
	tlsC := &tls.Config{
		SessionTicketKey: [32]byte{1},
		Rand: fixedRand([]string{
			"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
			"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
			"000102030405060708090A0B0C0D0E0F",
			"000102030405060708090A0B0C0D0E0F",
		}),
	}

	d := &dialer{
		JID:           "user@www.olabini.se",
		password:      "pass",
		serverAddress: "www.olabini.se:443",
		verifier:      &basicTLSVerifier{[]byte("aaaaa")},

		config: data.Config{
			TLSConfig: tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err.Error(), Equals, "tls: server certificate does not match expected hash (got: 82454418cb04854aa721bb0596528ff802b1e18a4e3a7767412ac9f108c9d3a7, want: 6161616161)")
}

func (s *ConnectionXmppSuite) Test_Dial_worksIfTheHandshakeSucceedsButSucceedsOnValidCertHash(c *C) {
	rw := &mockMultiConnIOReaderWriter{read: validTLSExchange}
	conn := &fullMockedConn{rw: rw}
	tlsC := &tls.Config{
		SessionTicketKey: [32]byte{1},
		Rand: fixedRand([]string{
			"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
			"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
			"000102030405060708090A0B0C0D0E0F",
			"000102030405060708090A0B0C0D0E0F",
		}),
	}

	d := &dialer{
		JID:           "user@www.olabini.se",
		password:      "pass",
		serverAddress: "www.olabini.se:443",
		verifier:      &basicTLSVerifier{bytesFromHex("82454418cb04854aa721bb0596528ff802b1e18a4e3a7767412ac9f108c9d3a7")},

		config: data.Config{
			TLSConfig: tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
}
