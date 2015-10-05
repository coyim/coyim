// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import "encoding/xml"

// bobData is a data element from http://xmpp.org/extensions/xep-0231.html.
type bobData struct {
	XMLName  xml.Name `xml:"urn:xmpp:bob data"`
	CID      string   `xml:"cid,attr"`
	MIMEType string   `xml:"type,attr"`
	Base64   string   `xml:",chardata"`
}
