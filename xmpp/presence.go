// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/coyim/coyim/xmpp/data"
)

func (c *conn) sendPresence(presence *data.ClientPresence) error {
	if len(presence.ID) == 0 {
		presence.ID = strconv.FormatUint(uint64(c.getCookie()), 10)
	}

	var outb bytes.Buffer

	_ = xml.NewEncoder(&outb).Encode(presence)
	_, e := c.safeWrite(outb.Bytes())
	return e
}

// SendPresence sends a presence stanza. If id is empty, a unique id is
// generated.
func (c *conn) SendPresence(to, typ, id, status string) error {
	p := &data.ClientPresence{
		ID:   id,
		To:   to,
		Type: typ,
	}

	if typ == "subscribe" && status != "" {
		p.Status = status
	}

	return c.sendPresence(p)
}

// SendMUCPresence sends a presence as first step for create a new room
func (c *conn) SendMUCPresence(to string, m *data.MUC) error {
	p := &data.ClientPresence{
		To:  to,
		MUC: m,
	}

	return c.sendPresence(p)
}

// SignalPresence will signal the current presence
func (c *conn) SignalPresence(state string) error {
	var outb bytes.Buffer

	//We dont use c.sendPresence() because this presence does not have `id` (why?)
	_, _ = fmt.Fprintf(&outb, "<presence><show>%s</show></presence>", xmlEscape(state))
	_, err := c.safeWrite(outb.Bytes())
	return err
}
