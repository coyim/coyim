// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// Conn represents a connection to an XMPP server.
type Conn struct {
	config *Config

	out     io.Writer
	rawOut  io.WriteCloser // doesn't log. Used for <auth>
	in      *xml.Decoder
	jid     string
	archive bool
	Rand    io.Reader

	lock          sync.Mutex
	inflights     map[Cookie]inflight
	customStorage map[xml.Name]reflect.Type

	lastPingRequest  time.Time
	lastPongResponse time.Time
}

// NewConn creates a new connection
func NewConn(in *xml.Decoder, out io.Writer, jid string) *Conn {
	return &Conn{
		in:  in,
		out: out,
		jid: jid,
	}
}

// Close closes the underlying connection
func (c *Conn) Close() error {
	// RFC 6120, Section 4.4 and 9.1.5
	fmt.Fprint(c.out, "</stream:stream>")
	//TODO: Wait for a </stream:stream> from server

	return c.rawOut.Close()
}

// Next reads stanzas from the server. If the stanza is a reply, it dispatches
// it to the correct channel and reads the next message. Otherwise it returns
// the stanza for processing.
func (c *Conn) Next() (stanza Stanza, err error) {
	for {
		if stanza.Name, stanza.Value, err = next(c); err != nil {
			return
		}

		if iq, ok := stanza.Value.(*ClientIQ); ok && (iq.Type == "result" || iq.Type == "error") {
			var cookieValue uint64
			if cookieValue, err = strconv.ParseUint(iq.ID, 16, 64); err != nil {
				err = errors.New("xmpp: failed to parse id from iq: " + err.Error())
				return
			}
			cookie := Cookie(cookieValue)

			c.lock.Lock()
			inflight, ok := c.inflights[cookie]
			c.lock.Unlock()

			if !ok {
				continue
			}
			if len(inflight.to) > 0 {
				// The reply must come from the address to
				// which we sent the request.
				if inflight.to != iq.From {
					continue
				}
			} else {
				// If there was no destination on the request
				// then the matching is more complex because
				// servers differ in how they construct the
				// reply.
				if len(iq.From) > 0 && iq.From != c.jid && iq.From != RemoveResourceFromJid(c.jid) && iq.From != domainFromJid(c.jid) {
					continue
				}
			}

			c.lock.Lock()
			delete(c.inflights, cookie)
			c.lock.Unlock()

			inflight.replyChan <- stanza
			continue
		}

		return
	}
}

// Cancel cancels and outstanding request. The request's channel is closed.
func (c *Conn) Cancel(cookie Cookie) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	inflight, ok := c.inflights[cookie]
	if !ok {
		return false
	}

	delete(c.inflights, cookie)
	close(inflight.replyChan)
	return true
}

// Send sends an IM message to the given user.
func (c *Conn) Send(to, msg string) error {
	archive := ""
	if !c.archive {
		// The first part of archive is from google:
		// See https://developers.google.com/talk/jep_extensions/otr
		// The second part of the stanza is from XEP-0136
		// http://xmpp.org/extensions/xep-0136.html#pref-syntax-item-otr
		// http://xmpp.org/extensions/xep-0136.html#otr-nego
		archive = "<nos:x xmlns:nos='google:nosave' value='enabled'/><arc:record xmlns:arc='http://jabber.org/protocol/archive' otr='require'/>"
	}
	_, err := fmt.Fprintf(c.out, "<message to='%s' from='%s' type='chat'><body>%s</body>%s</message>", xmlEscape(to), xmlEscape(c.jid), xmlEscape(msg), archive)
	return err
}
