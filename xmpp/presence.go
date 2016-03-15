// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"fmt"
	"strconv"
)

// SendPresence sends a presence stanza. If id is empty, a unique id is
// generated.
func (c *conn) SendPresence(to, typ, id, status string) error {
	if len(id) == 0 {
		id = strconv.FormatUint(uint64(c.getCookie()), 10)
	}
	end := "/>"
	if typ == "subscribe" && status != "" {
		end = fmt.Sprintf("><status>%s</status></presence>", xmlEscape(status))
	}
	_, err := fmt.Fprintf(c.out, "<presence id='%s' to='%s' type='%s'%s", xmlEscape(id), xmlEscape(to), xmlEscape(typ), end)
	return err
}

// SignalPresence will signal the current presence
func (c *conn) SignalPresence(state string) error {
	_, err := fmt.Fprintf(c.out, "<presence><show>%s</show></presence>", xmlEscape(state))
	return err
}
