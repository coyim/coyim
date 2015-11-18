// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
// Ping implements the XMPP extension Ping, as specified in xep-0199
package xmpp

import (
	"encoding/xml"
	"errors"
	"log"
	"time"
)

type PingRequest struct {
	XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}

// SendPing sends a Ping request.
func (c *Conn) SendPing() (reply chan Stanza, cookie Cookie, err error) {
	c.lastPingRequest = time.Now() //TODO: this seems should not belong to Conn
	return c.SendIQ("", "get", PingRequest{})
}

// SendPingReply sends a reply to a Ping request.
func (c *Conn) SendPingReply(id string) error {
	return c.SendIQReply("", "result", id, EmptyReply{})
}

// ReceivePong update the timestamp for lastPongResponse,
func (c *Conn) ReceivePong() {
	c.lastPongResponse = time.Now() //TODO: this seems should not belong to Conn
}

// ParsePong parse a reply of a Pong response.
func ParsePong(reply Stanza) error {
	iq, ok := reply.Value.(*ClientIQ)
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

func (c *Conn) watchPings() {
	failures := 0

	for !c.closed {
		<-time.After(pingIterval)

		log.Println("xmpp: send ping")
		pongReply, _, err := c.SendPing()
		if err != nil {
			return
		}

		select {
		case <-time.After(pingTimeout):
			log.Println("xmpp: pong not received")
			failures = failures + 1
		case pongStanza, ok := <-pongReply:
			if !ok {
				//somebody canceled the IQ
				continue
			}

			log.Println("xmpp: receive pong")
			failures = 0
			iq, ok := pongStanza.Value.(*ClientIQ)
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

		log.Println("xmpp: ping failures reached treshold of ", maxPingFailures)
		go c.sendStreamError(StreamError{
			DefinedCondition: ConnectionTimeout,
		})

		return
	}
}
