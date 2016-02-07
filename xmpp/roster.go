// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"sort"

	"github.com/twstrike/coyim/xmpp/data"
)

// RequestRoster requests the user's roster from the server. It returns a
// channel on which the reply can be read when received and a Cookie that can
// be used to cancel the request.
func (c *Conn) RequestRoster() (<-chan data.Stanza, Cookie, error) {
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

type rosterEntries []data.RosterEntry

func (entries rosterEntries) Len() int {
	return len(entries)
}

func (entries rosterEntries) Less(i, j int) bool {
	return entries[i].Jid < entries[j].Jid
}

func (entries rosterEntries) Swap(i, j int) {
	entries[i], entries[j] = entries[j], entries[i]
}

// ParseRoster extracts roster information from the given Stanza.
func ParseRoster(reply data.Stanza) ([]data.RosterEntry, error) {
	iq, ok := reply.Value.(*data.ClientIQ)
	if !ok {
		return nil, errors.New("xmpp: roster request resulted in tag of type " + reply.Name.Local)
	}

	var roster data.Roster
	if err := xml.NewDecoder(bytes.NewBuffer(iq.Query)).Decode(&roster); err != nil {
		return nil, err
	}
	sort.Sort(rosterEntries(roster.Item))
	return roster.Item, nil
}
