// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
// Ping implements the XMPP extension Ping, as specified in xep-0199
package xmpp

import (
	"errors"
	"log"
	"time"

	"github.com/twstrike/coyim/xmpp/data"
)

// SendPing sends a Ping request.
func (c *conn) SendPing() (reply chan data.Stanza, cookie data.Cookie, err error) {
	c.lastPingRequest = time.Now() //TODO: this seems should not belong to Conn
	return c.SendIQ("", "get", data.PingRequest{})
}

// SendPingReply sends a reply to a Ping request.
func (c *conn) SendPingReply(id string) error {
	return c.SendIQReply("", "result", id, data.EmptyReply{})
}

// ReceivePong update the timestamp for lastPongResponse,
func (c *conn) ReceivePong() {
	c.lastPongResponse = time.Now() //TODO: this seems should not belong to Conn
}

// ParsePong parse a reply of a Pong response.
func ParsePong(reply data.Stanza) error {
	iq, ok := reply.Value.(*data.ClientIQ)
	if !ok {
		return errors.New("xmpp: ping request resulted in tag of type " + reply.Name.Local)
	}
	switch iq.Type {
	case "result":
		return nil
	case "error":
		return errors.New("xmpp: ping request resulted in a error: " + iq.Error.Text)
	default:
		return errors.New("xmpp: ping request resulted in a unexpected type")
	}
}

var (
	pingIterval     = 10 * time.Second //should be 5 minutes at least, per spec
	pingTimeout     = 30 * time.Second
	maxPingFailures = 2
)

func (c *conn) watchPings() {
	tick := time.NewTicker(pingIterval)
	defer tick.Stop()
	failures := 0

	for _ = range tick.C {
		if c.closed {
			return
		}

		pongReply, _, err := c.SendPing()
		if err != nil {
			return
		}

		select {
		case <-time.After(pingTimeout):
			// ping timed out
			failures = failures + 1
		case pongStanza, ok := <-pongReply:
			if !ok {
				// pong cancelled
				continue
			}

			failures = 0
			iq, ok := pongStanza.Value.(*data.ClientIQ)
			if !ok {
				//not the expected IQ
				return
			}

			//TODO: check for <service-unavailable/>
			if iq.Type == "error" {
				//server does not support Ping
				return
			}
		}

		if failures < maxPingFailures {
			continue
		}

		log.Println("xmpp: ping failures reached threshold of", maxPingFailures)
		go c.sendStreamError(data.StreamError{
			DefinedCondition: data.ConnectionTimeout,
		})

		return
	}
}
