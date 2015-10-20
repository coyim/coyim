// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import "encoding/xml"

// Delay represents XEP-0203: Delayed Delivery of <message/> and <presence/> stanzas.
type Delay struct {
	XMLName xml.Name `xml:"urn:xmpp:delay delay"`
	From    string   `xml:"from,attr,omitempty"`
	Stamp   string   `xml:"stamp,attr"`

	Body string `xml:",chardata"`
}
