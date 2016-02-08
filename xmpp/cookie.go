// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"encoding/binary"

	"github.com/twstrike/coyim/xmpp/data"
)

func (c *conn) getCookie() data.Cookie {
	var buf [8]byte
	if _, err := c.Rand().Read(buf[:]); err != nil {
		panic("Failed to read random bytes: " + err.Error())
	}
	return data.Cookie(binary.LittleEndian.Uint64(buf[:]))
}

func (c *conn) cancelInflights() {
	for cookie := range c.inflights {
		c.Cancel(cookie)
	}
}
