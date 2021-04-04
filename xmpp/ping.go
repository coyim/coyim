// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
// Ping implements the XMPP extension Ping, as specified in xep-0199
package xmpp

import (
	"errors"
	"time"

	"github.com/coyim/coyim/xmpp/data"
)

// SendPing sends a Ping request.
func (c *conn) SendPing() (reply <-chan data.Stanza, cookie data.Cookie, err error) {
	//TODO: should not this be set when we send any message? Why would I send a ping
	//just after sending a message?
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
		return errors.New("xmpp: ping request resulted in an error: " + iq.Error.Text)
	default:
		return errors.New("xmpp: ping request resulted in an unexpected type")
	}
}

var (
	pingInterval    = 10 * time.Second //should be 5 minutes at least, per spec
	pingTimeout     = 30 * time.Second
	maxPingFailures = 2
)

var newTicker = time.NewTicker

func (c *conn) watchPings() {
	tick := newTicker(pingInterval)
	defer tick.Stop()
	failures := 0

	for range tick.C {
		if c.closed {
			c.log.Info("xmpp: trying to send ping on closed connection")
			return
		}

		pongReply, _, err := c.SendPing()
		if err != nil {
			c.log.WithError(err).Warn("xmpp: error when sending ping")
			return
		}

		select {
		case <-time.After(pingTimeout):
			failures = failures + 1

			if failures >= maxPingFailures {
				c.log.WithField("threshold", maxPingFailures).Warn("xmpp: ping failures reached threshold")

				_ = c.sendStreamError(data.StreamError{
					DefinedCondition: data.ConnectionTimeout,
				})

				return
			}
		case pongStanza, ok := <-pongReply:
			if !ok {
				c.log.Info("xmpp: ping result channel closed")
				continue
			}

			failures = 0
			iq, ok := pongStanza.Value.(*data.ClientIQ)
			if !ok || iq.Type == "error" {
				c.log.WithField("value", pongStanza.Value).Info("xmpp: received invalid result to ping")
				return
			}
		}

	}
}
