//+build go1.8

package xmpp

import (
	"crypto/tls"
	"io"

	"github.com/coyim/coyim/xmpp/data"
	. "gopkg.in/check.v1"
)

func recordTLSHandshake(rw *mockConnIOReaderWriter, conf *tls.Config) ([]byte, error) {
	conn := &fullMockedConn{rw: rw}
	return rw.write, tls.Client(conn, conf).Handshake()
}

func (s *ConnectionXMPPSuite) Test_Dial_failsWhenStartingAHandshake_Proposed(c *C) {
	r := []string{"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F"}
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

	expectedXMPPStartTLS := "<?xml version='1.0'?><stream:stream to='domain' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n" +
		"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>"

	//The server did not send its part of the handshake (see rw above), so there's nothing to be read
	recordedHandshake, recordedErr := recordTLSHandshake(&mockConnIOReaderWriter{}, &tls.Config{
		ServerName: "domain",
		Rand:       fixedRand(r),
	})

	tlsC := &tls.Config{
		Rand: fixedRand(r),
	}

	d := &dialer{
		JID:      "user@domain",
		password: "pass",
		config: data.Config{
			TLSConfig: tlsC,
		},
	}
	conn := &fullMockedConn{rw: rw}
	_, err := d.setupStream(conn)

	c.Assert(err, Equals, io.EOF)
	c.Assert(err, Equals, recordedErr)
	c.Assert(string(rw.write)[:len(expectedXMPPStartTLS)], Equals, expectedXMPPStartTLS)
	c.Assert(rw.write[len(expectedXMPPStartTLS):], DeepEquals, recordedHandshake)
}

func (s *ConnectionXMPPSuite) Test_Dial_worksIfTheHandshakeSucceeds_Proposed(c *C) {
	r := []string{
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F",
		"000102030405060708090A0B0C0D0E0F",
		"000102030405060708090A0B0C0D0E0F",
	}
	expectedXMPPStartTLS := "<?xml version='1.0'?>" +
		"<stream:stream to='www.olabini.se' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams' version='1.0'>\n" +
		"<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/>"

	rw := &mockMultiConnIOReaderWriter{read: validTLSExchange}

	//Skips the XMPP-specific reply
	reply := []byte{}
	for _, part := range validTLSExchange[1:] {
		reply = append(reply, part...)
	}
	recordedServerReply := &mockConnIOReaderWriter{read: reply}
	recordedHandshake, recordedErr := recordTLSHandshake(recordedServerReply, &tls.Config{
		ServerName: "www.olabini.se",
		Rand:       fixedRand(r),
	})

	tlsC := &tls.Config{
		Rand: fixedRand(r),
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

	conn := &fullMockedConn{rw: rw}
	_, err := d.setupStream(conn)

	//This is different on 1.6, 1.7 ("EOF") but does the error string really matter?
	//c.Assert(err.Error(), Equals, "tls: server's Finished message was incorrect")
	c.Assert(err.Error(), Equals, recordedErr.Error())
	c.Assert(string(rw.write)[:len(expectedXMPPStartTLS)], Equals, expectedXMPPStartTLS)
	c.Assert(rw.write[len(expectedXMPPStartTLS):], DeepEquals, recordedHandshake)
}
