// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import "encoding/xml"

// Extensions implements generic XEPs.
type Extensions []*Extension

// Extension represents any XML node not included in the Stanza definition
type Extension struct {
	XMLName xml.Name
	Body    string `xml:",innerxml"`
}

// StanzaError implements RFC 3920, section 9.3.
type StanzaError struct {
	// cancel -- do not retry (the error is unrecoverable)
	// continue -- proceed (the condition was only a warning)
	// modify -- retry after changing the data sent
	// auth -- retry after providing credentials
	// wait -- retry after waiting (the error is temporary)
	Type string `xml:"type,attr"`

	Condition struct {
		XMLName xml.Name
		Body    string `xml:",innerxml"`
	} `xml:",any"`
}

// ClientMessage implements RFC 3921  B.1  jabber:client
type ClientMessage struct {
	XMLName xml.Name `xml:"jabber:client message"`
	From    string   `xml:"from,attr"`
	ID      string   `xml:"id,attr"`
	To      string   `xml:"to,attr"`
	Type    string   `xml:"type,attr"` // chat, error, groupchat, headline, or normal

	// These should technically be []clientText,
	// but string is much more convenient.
	Subject string `xml:"subject"`
	Body    string `xml:"body"`
	Thread  string `xml:"thread"`
	Delay   *Delay `xml:"delay,omitempty"`

	Error *StanzaError `xml:"error"`

	Extensions `xml:",any,omitempty"`
}

// ClientCaps contains information about client capabilities
type ClientCaps struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/caps c"`
	Ext     string   `xml:"ext,attr"`
	Hash    string   `xml:"hash,attr"`
	Node    string   `xml:"node,attr"`
	Ver     string   `xml:"ver,attr"`
}

// ClientError represents a client error
type ClientError struct {
	XMLName xml.Name `xml:"jabber:client error"`
	Code    string   `xml:"code,attr"`
	Type    string   `xml:"type,attr"`
	Any     xml.Name `xml:",any"`
	Text    string   `xml:"text"`
}
