// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/twstrike/coyim/sasl"
	"github.com/twstrike/coyim/sasl/digestmd5"
	"github.com/twstrike/coyim/sasl/plain"
	"github.com/twstrike/coyim/sasl/scram"
)

var (
	// ErrAuthenticationFailed indicates a failure to authenticate to the server with the user and password provided.
	ErrAuthenticationFailed = errors.New("could not authenticate to the XMPP server")

	errUnsupportedSASLMechanism = errors.New("xmpp: server does not support any of the prefered SASL mechanism")
)

func init() {
	plain.Register()
	digestmd5.Register()
	scram.Register()
}

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

func (c *Conn) authenticate(user, password string) error {
	// TODO: Google accounts with 2-step auth MUST use app-specific passwords:
	// https://security.google.com/settings/security/apppasswords
	// An alternative would be implementing the Google authentication mechanisms
	// - X-OAUTH2: requires app registration on Google, and a server to receive the oauth callback
	// https://developers.google.com/talk/jep_extensions/oauth?hl=en
	// - X-GOOGLE-TOKEN: seems to be this https://developers.google.com/identity/protocols/AuthForInstalledApps

	return c.authenticateWithPreferedMethod(user, password)
}

func (c *Conn) authenticateWithPreferedMethod(user, password string) error {
	//TODO: this should be configurable by the client
	preferedMechanisms := []string{"SCRAM-SHA-1", "DIGEST-MD5", "PLAIN"}

	log.Println("sasl: server supports mechanisms", c.features.Mechanisms.Mechanism)

	for _, prefered := range preferedMechanisms {
		if !sasl.ClientSupport(prefered) {
			continue
		}

		for _, m := range c.features.Mechanisms.Mechanism {
			if prefered == m {
				log.Println("sasl: authenticating via", m)
				return c.authenticatWith(prefered, user, password)
			}
		}
	}

	return errUnsupportedSASLMechanism
}

func clientNonce(r io.Reader) (string, error) {
	//TODO: what is the appropriate size for this?
	//TODO: what is the appropriate way to generate a cnonce?
	conceRand := make([]byte, 7)
	_, err := io.ReadFull(r, conceRand)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(conceRand), nil
}

func (c *Conn) authenticatWith(mechanism string, user string, password string) error {
	clientAuth, err := sasl.NewClient(mechanism)
	if err != nil {
		return err
	}

	clientNonce, err := clientNonce(c.rand())
	if err != nil {
		return err
	}

	clientAuth.SetProperty(sasl.AuthID, user)
	clientAuth.SetProperty(sasl.Password, password)

	clientAuth.SetProperty(sasl.Service, "xmpp")
	clientAuth.SetProperty(sasl.QOP, "auth")

	//TODO: this should come from username if it were a full JID
	//clientAuth.SetProperty(sasl.Realm, "")

	clientAuth.SetProperty(sasl.ClientNonce, clientNonce)

	t, err := clientAuth.Step(nil)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.rawOut, "<auth xmlns='%s' mechanism='%s'>%s</auth>\n", NsSASL, mechanism, t.Encode())

	return c.challengeLoop(clientAuth)
}

func (c *Conn) challengeLoop(clientAuth sasl.Session) error {
	for {
		t, success, err := c.receiveChallenge()
		if err != nil {
			return err
		}

		t, err = clientAuth.Step(t)
		if err != nil {
			return err
		}

		if success {
			return nil
		}

		if !clientAuth.NeedsMore() {
			break
		}

		fmt.Fprintf(c.rawOut, "<response xmlns='%s'>%s</response>\n", NsSASL, t.Encode())
	}

	return ErrAuthenticationFailed
}

func (c *Conn) receiveChallenge() (t sasl.Token, success bool, err error) {
	var encodedChallenge []byte

	name, val, _ := next(c)
	switch v := val.(type) {
	case *saslFailure:
		err = errors.New("xmpp: authentication failure: " + v.String())
		return
	case *saslSuccess:
		encodedChallenge = v.Content
		success = true
	case *string:
		if name.Local != "challenge" || name.Space != NsSASL {
			err = errors.New("xmpp: unexpected <" + name.Local + "> in " + name.Space)
			return
		}

		encodedChallenge = []byte(*v)
	}

	t, err = sasl.DecodeToken(encodedChallenge)
	return
}

//TODO: decide how to get the OAuth token from the OAuth code
func (c *Conn) xOAuth(user, token string) error {
	encoding := base64.StdEncoding

	raw := "\x00" + user + "\x00" + token
	enc := make([]byte, encoding.EncodedLen(len(raw)))
	encoding.Encode(enc, []byte(raw))
	fmt.Fprintf(c.rawOut, "<auth xmlns='%s' mechanism='X-OAUTH2' auth:service='oauth2' xmlns:auth='%s'>%s</auth>\n", NsSASL, NsXOAuth, enc)

	_, err := c.verifyAuthenticationSuccess()
	return err
}

func (c *Conn) verifyAuthenticationSuccess() (*saslSuccess, error) {
	// Next message should be either success or failure.
	name, val, _ := next(c)
	switch v := val.(type) {
	case *saslSuccess:
		return v, nil
	case *saslFailure:
		return nil, errors.New("xmpp: authentication failure: " + v.String())
	}

	return nil, errors.New("expected <success> or <failure>, got <" + name.Local + "> in " + name.Space)
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
	Content []byte   `xml:",innerxml"`
}

type saslFailure struct {
	XMLName          xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl failure"`
	Text             string   `xml:"text,omitempty"`
	DefinedCondition Any      `xml:",any"`
}

// Condition returns a SASL-related error condition
func (f saslFailure) Condition() SASLErrorCondition {
	return SASLErrorCondition(f.DefinedCondition.XMLName.Local)
}

func (f saslFailure) String() string {
	if f.Text != "" {
		return fmt.Sprintf("%s: %q", f.Condition(), f.Text)
	}

	return string(f.Condition())
}

// SASLErrorCondition represents a defined SASL-related error conditions as defined in RFC 6120, section 6.5
type SASLErrorCondition string

// SASL error conditions as defined in RFC 6120, section 6.5
const (
	SASLAborted              SASLErrorCondition = "aborted"
	SASLAccountDisabled                         = "account-disabled"
	SASLCredentialsExpired                      = "credentials-expired"
	SASLEncryptionRequired                      = "encryption-required"
	SASLIncorrectEncoding                       = "incorrect-encoding"
	SASLInvalidAuthzid                          = "invalid-authzid"
	SASLInvalidMechanism                        = "invalid-mechanism"
	SASLMalformedRequest                        = "malformed-request"
	SASLMechanismTooWeak                        = "mechanism-too-weak"
	SASLNotAuthorized                           = "not-authorized"
	SASLTemporaryAuthFailure                    = "temporary-auth-failure"
)
