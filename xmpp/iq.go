// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"encoding/xml"
	"fmt"
)

// SendIQ sends an info/query message to the given user. It returns a channel
// on which the reply can be read when received and a Cookie that can be used
// to cancel the request.
func (c *Conn) SendIQ(to, typ string, value interface{}) (reply chan Stanza, cookie Cookie, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	cookie = c.getCookie()
	reply = make(chan Stanza, 1)

	toAttr := ""
	if len(to) > 0 {
		toAttr = "to='" + xmlEscape(to) + "'"
	}
	if _, err = fmt.Fprintf(c.out, "<iq %s from='%s' type='%s' id='%x'>", toAttr, xmlEscape(c.jid), xmlEscape(typ), cookie); err != nil {
		return
	}
	if _, ok := value.(EmptyReply); !ok {
		if err = xml.NewEncoder(c.out).Encode(value); err != nil {
			return
		}
	}
	if _, err = fmt.Fprintf(c.out, "</iq>"); err != nil {
		return
	}

	c.inflights[cookie] = inflight{reply, to}
	return
}

// SendIQReply sends a reply to an IQ query.
func (c *Conn) SendIQReply(to, typ, id string, value interface{}) error {
	if _, err := fmt.Fprintf(c.out, "<iq to='%s' from='%s' type='%s' id='%s'>", xmlEscape(to), xmlEscape(c.jid), xmlEscape(typ), xmlEscape(id)); err != nil {
		return err
	}
	if _, ok := value.(EmptyReply); !ok {
		if err := xml.NewEncoder(c.out).Encode(value); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(c.out, "</iq>")
	return err
}

// ClientIQ contains a specific information query request
type ClientIQ struct { // info/query
	XMLName xml.Name    `xml:"jabber:client iq"`
	From    string      `xml:"from,attr"`
	ID      string      `xml:"id,attr"`
	To      string      `xml:"to,attr"`
	Type    string      `xml:"type,attr"` // error, get, result, set
	Error   ClientError `xml:"error"`
	Bind    bindBind    `xml:"bind"`
	Query   []byte      `xml:",innerxml"`
}

// An EmptyReply results in in no XML.
type EmptyReply struct {
}
