// +build !go1.6

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
			"\v\f\r\x0e\x0fҟ\xb0>ʋ,\xcbI\x9e\x1c\x94\b\x1an\x12cT\x81\xac'\xf4{\rV\xb1V\xad]\xc5\b\xfe\xf4rh˱\b?\x10&ե\x89\"\xbf\x8a0\x17\x03\x03\x00\xc0\x00\x01\x02\x03\x04\x05\x06\a\b\t\n"+
			"\v\f\r\x0e\x0f7ѱ\n"+
			"9[\x9d\\!\xcd4\xfa\x18\xd6y\xb0y\xf2\x94T\x96\xb3\x9e\x16\xfdŷ\x9a\xf3Ʒ\xa9n@\xa5\tO\x8cE\x84\x85\x15#\x04\x97\xce\xfc\x12\x93;\xd8\xcdl>Q\xe52\x9f\xc1\x84\xf2cj\x81_U\x86\xcf6\xadڦC\xe6\x13\xfa\xc8%-\x15\xd5E\xe2h\x91\xfc\xa0+\x94\x13\xe3gG\xe6\xff\xe5`'و\x9fM\xea\xc780N\xd5\u9fdc\xf5)\xfdݫ3\xf0\x8ee/`8@\x88t\xdfc\xd3d-\xa9S\x80\x1a\x95rV\x98\x1e\xb2\xde\xdc\xd8\xc8\n"+
			"\x9aE\x12\xd1U\x95\xd37\xfe\xccm\u007f \x1a\xb5\x91\xd80;\\S\xb0\x04\x85Av+\xefǗ")
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

	c.Assert(err.Error(), Equals, "tls: server certificate does not match expected hash (got: 82454418cb04854aa721bb0596528ff802b1e18a4e3a7767412ac9f108c9d3a7, want: 6161616161)")
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
		verifier:      &basicTLSVerifier{bytesFromHex("82454418cb04854aa721bb0596528ff802b1e18a4e3a7767412ac9f108c9d3a7")},

		config: data.Config{
			TLSConfig: &tlsC,
		},
	}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
}
