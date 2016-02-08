// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"fmt"

	"github.com/twstrike/coyim/xmpp/data"
)

// RequestRoster requests the user's roster from the server. It returns a
// channel on which the reply can be read when received and a Cookie that can
// be used to cancel the request.
func (c *conn) RequestRoster() (<-chan data.Stanza, data.Cookie, error) {
	cookie := c.getCookie()
	if _, err := fmt.Fprintf(c.out, "<iq type='get' id='%x'><query xmlns='jabber:iq:roster'/></iq>", cookie); err != nil {
		return nil, 0, err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	ch := make(chan data.Stanza, 1)
	c.inflights[cookie] = inflight{ch, ""}
	return ch, cookie, nil
}
