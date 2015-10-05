// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import "encoding/xml"

type RegisterQuery struct {
	XMLName  xml.Name  `xml:"jabber:iq:register query"`
	Username *xml.Name `xml:"username"`
	Password *xml.Name `xml:"password"`
	Form     Form      `xml:"x"`
	Datas    []bobData `xml:"data"`
}
