// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"encoding/xml"
	"fmt"

	"github.com/coyim/coyim/xmpp/data"
)

type rawXML []byte

// SendIQ sends an info/query message to the given user. It returns a channel
// on which the reply can be read when received and a Cookie that can be used
// to cancel the request.
func (c *conn) SendIQ(to, typ string, value interface{}) (reply chan data.Stanza, cookie data.Cookie, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	cookie = c.getCookie()
	reply = make(chan data.Stanza, 1)

	toAttr := ""
	if len(to) > 0 {
		toAttr = "to='" + xmlEscape(to) + "'"
	}

	if _, err = fmt.Fprintf(c.out, "<iq xmlns='jabber:client' %s from='%s' type='%s' id='%x'>", toAttr, xmlEscape(c.jid), xmlEscape(typ), cookie); err != nil {
		return
	}

	switch v := value.(type) {
	case data.EmptyReply:
		//nothing
	case rawXML:
		_, err = c.out.Write(v)
	default:
		err = xml.NewEncoder(c.out).Encode(value)
	}

	if err != nil {
		return
	}

	if _, err = fmt.Fprintf(c.out, "</iq>"); err != nil {
		return
	}

	c.inflights[cookie] = inflight{reply, to}
	return
}

// SendIQReply sends a reply to an IQ query.
func (c *conn) SendIQReply(to, typ, id string, value interface{}) error {
	if _, err := fmt.Fprintf(c.out, "<iq to='%s' from='%s' type='%s' id='%s'>", xmlEscape(to), xmlEscape(c.jid), xmlEscape(typ), xmlEscape(id)); err != nil {
		return err
	}
	if _, ok := value.(data.EmptyReply); !ok {
		if err := xml.NewEncoder(c.out).Encode(value); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(c.out, "</iq>")
	return err
}
