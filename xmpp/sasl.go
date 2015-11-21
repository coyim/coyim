// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
)

var (
	// ErrAuthenticationFailed indicates a failure to authenticate to the server with the user and password provided.
	ErrAuthenticationFailed = errors.New("could not authenticate to the XMPP server")
)

// SASL negotiation. RFC 6120, section 6
func (d *Dialer) negotiateSASL(c *Conn) error {
	originDomain := d.getJIDDomainpart()
	user := d.getJIDLocalpart()
	password := d.Password

	if err := c.authenticate(user, password); err != nil {
		return ErrAuthenticationFailed
	}

	return c.sendInitialStreamHeader(originDomain)
}

func (c *Conn) authenticate(user, password string) (err error) {
	// TODO: Section 13.8 mandates clients and servers to implement SCRAM, DIGEST-MD5 and PLAIN.
	// SCRAM replaces DIGEST-MD5 as mandatory, but for compatibility reasons we should keep DIGEST-MD5
	// - SCRAM-SHA-1 and SCRAM-SHA-1-PLUS: RFC 5802
	// - DIGEST-MD5: RFC 2831
	// See: https://prosody.im/doc/plain_or_hashed

	// TODO: Google accounts with 2-step auth MUST use app-specific passwords:
	// https://security.google.com/settings/security/apppasswords
	// An alternative would be implementing the Google authentication mechanisms
	// - X-OAUTH2: requires app registration on Google, and a server to receive the oauth callback
	// https://developers.google.com/talk/jep_extensions/oauth?hl=en
	// - X-GOOGLE-TOKEN: seems to be this https://developers.google.com/identity/protocols/AuthForInstalledApps

	// TODO: this is dangerous, specially with skipTLS
	havePlain := false
	for _, m := range c.features.Mechanisms.Mechanism {
		if m == "PLAIN" {
			havePlain = true
			break
		}
	}

	if !havePlain {
		return errors.New("xmpp: PLAIN authentication is not an option")
	}

	// Plain authentication: send base64-encoded \x00 user \x00 password.
	raw := "\x00" + user + "\x00" + password
	enc := make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
	base64.StdEncoding.Encode(enc, []byte(raw))
	fmt.Fprintf(c.rawOut, "<auth xmlns='%s' mechanism='PLAIN'>%s</auth>\n", NsSASL, enc)

	// Next message should be either success or failure.
	name, val, err := next(c)
	switch v := val.(type) {
	case *saslSuccess:
	case *saslFailure:
		// v.Any is type of sub-element in failure,
		// which gives a description of what failed.
		return errors.New("xmpp: authentication failure: " + v.Any.Local)
	default:
		return errors.New("expected <success> or <failure>, got <" + name.Local + "> in " + name.Space)
	}

	return nil
}

// RFC 3920  C.4  SASL name space
//TODO RFC 6120 obsoletes RFC 3920
type saslMechanisms struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	Mechanism []string `xml:"mechanism"`
}

type saslAuth struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl auth"`
	Mechanism string   `xml:"mechanism,attr"`
}

type saslChallenge string

type saslResponse string

type saslAbort struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl abort"`
}

type saslSuccess struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl success"`
}

type saslFailure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl failure"`
	Any     xml.Name `xml:",any"`
}
