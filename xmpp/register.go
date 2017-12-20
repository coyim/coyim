// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
)

// Various errors signalled by the registration component
var (
	ErrUsernameConflict                = errors.New("xmpp: the username is not available for registration")
	ErrMissingRequiredRegistrationInfo = errors.New("xmpp: missing required registration information")
	ErrRegistrationFailed              = errors.New("xmpp: account creation failed")
	ErrWrongCaptcha                    = errors.New("xmpp: the captcha entered is wrong")
	ErrResourceConstraint              = errors.New("xmpp: already reached the configured number of allowable resources")
)

// XEP-0077
func (d *dialer) negotiateInBandRegistration(c interfaces.Conn) (bool, error) {
	if c.Features().InBandRegistration == nil {
		return false, nil
	}

	user := d.getJIDLocalpart()
	return c.RegisterAccount(user, d.password)
}

func (c *conn) RegisterAccount(user, password string) (bool, error) {
	if c.config.CreateCallback == nil {
		return false, nil
	}

	err := c.createAccount(user, password)
	if err != nil {
		return false, err
	}

	return true, c.closeImmediately()
}

func (c *conn) createAccount(user, password string) error {
	io.WriteString(c.config.GetLog(), "Attempting to create account\n")
	fmt.Fprintf(c.out, "<iq type='get' id='create_1'><query xmlns='jabber:iq:register'/></iq>")
	var iq data.ClientIQ
	if err := c.in.DecodeElement(&iq, nil); err != nil {
		return errors.New("unmarshal <iq>: " + err.Error())
	}

	if iq.Type != "result" {
		return ErrRegistrationFailed
	}

	var register data.RegisterQuery

	if err := xml.NewDecoder(bytes.NewBuffer(iq.Query)).Decode(&register); err != nil {
		return err
	}

	if len(register.Form.Type) > 0 {
		reply, err := processForm(&register.Form, register.Datas, c.config.CreateCallback)
		fmt.Fprintf(c.rawOut, "<iq type='set' id='create_2'><query xmlns='jabber:iq:register'>")
		if err = xml.NewEncoder(c.rawOut).Encode(reply); err != nil {
			return err
		}

		fmt.Fprintf(c.rawOut, "</query></iq>")
	} else if register.Username != nil && register.Password != nil {
		//TODO: make sure this only happens via SSL
		//TODO: should generate form asking for username and password,
		//and call processForm for consistency

		// Try the old-style registration.
		fmt.Fprintf(c.rawOut, "<iq type='set' id='create_2'><query xmlns='jabber:iq:register'><username>%s</username><password>%s</password></query></iq>", user, password)
	}

	iq2 := &data.ClientIQ{}
	if err := c.in.DecodeElement(iq2, nil); err != nil {
		return errors.New("unmarshal <iq>: " + err.Error())
	}

	if iq2.Type == "error" {
		switch iq2.Error.Condition.XMLName.Local {
		case "conflict":
			// <conflict xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>
			return ErrUsernameConflict
		case "not-acceptable":
			// <not-acceptable xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>
			return ErrMissingRequiredRegistrationInfo
		// TODO: this case shouldn't happen
		case "bad-request":
			//<bad-request xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>
			return ErrRegistrationFailed
		case "not-allowed":
			//<not-allowed xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>
			return ErrWrongCaptcha
		case "resource-constraint":
			//<resource-constraint xmlns='urn:ietf:params:xml:ns:xmpp-stanzas'/>
			return ErrResourceConstraint
		default:
			return ErrRegistrationFailed
		}
	}

	return nil
}

// CancelRegistration cancels the account registration with the server
func (c *conn) CancelRegistration() (reply <-chan data.Stanza, cookie data.Cookie, err error) {
	// https://xmpp.org/extensions/xep-0077.html#usecases-cancel
	registrationCancel := rawXML(`
	<query xmlns='jabber:iq:register'>
		<remove/>
	</query>
	`)

	return c.SendIQ("", "set", registrationCancel)
}
