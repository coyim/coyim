// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"crypto/rand"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/twstrike/coyim/xmpp/data"
	"github.com/twstrike/coyim/xmpp/interfaces"
	"github.com/twstrike/coyim/xmpp/utils"
)

// conn represents a connection to an XMPP server.
type conn struct {
	config data.Config

	in           *xml.Decoder
	out          io.Writer
	rawOut       io.WriteCloser // doesn't log. Used for <auth>
	keepaliveOut io.Writer

	jid           string
	originDomain  string
	features      data.StreamFeatures
	serverAddress string

	rand          io.Reader
	lock          sync.Mutex
	inflights     map[data.Cookie]inflight
	customStorage map[xml.Name]reflect.Type

	lastPingRequest  time.Time
	lastPongResponse time.Time

	streamCloseReceived chan bool
	closed              bool
}

func (c *conn) CustomStorage() map[xml.Name]reflect.Type {
	return c.customStorage
}

func (c *conn) OriginDomain() string {
	return c.originDomain
}

func (c *conn) Lock() *sync.Mutex {
	return &c.lock
}

func (c *conn) SetInOut(i *xml.Decoder, o io.Writer) {
	c.in = i
	c.out = o
}

func (c *conn) SetRawOut(o io.WriteCloser) {
	c.rawOut = o
}

func (c *conn) SetKeepaliveOut(o io.Writer) {
	c.keepaliveOut = o
}

func (c *conn) Features() data.StreamFeatures {
	return c.features
}

func (c *conn) Config() *data.Config {
	return &c.config
}

func (c *conn) In() *xml.Decoder {
	return c.in
}

func (c *conn) RawOut() io.WriteCloser {
	return c.rawOut
}

func (c *conn) SetServerAddress(s1 string) {
	c.serverAddress = s1
}

func (c *conn) ServerAddress() string {
	return c.serverAddress
}

func (c *conn) Out() io.Writer {
	return c.out
}

// NewConn creates a new connection
//TODO: this is used only for testing. Remove when we have a Conn interface
func NewConn(in *xml.Decoder, out io.WriteCloser, jid string) interfaces.Conn {
	conn := &conn{
		in:     in,
		out:    out,
		rawOut: out,
		jid:    jid,

		inflights:           make(map[data.Cookie]inflight),
		streamCloseReceived: make(chan bool),

		rand: rand.Reader,
	}

	conn.closed = true
	close(conn.streamCloseReceived) // closes immediately
	return conn
}

func newConn() *conn {
	return &conn{
		inflights:           make(map[data.Cookie]inflight),
		streamCloseReceived: make(chan bool),

		rand: rand.Reader,
	}
}

// Close closes the underlying connection
func (c *conn) Close() error {
	if c.closed {
		return errors.New("xmpp: the connection is already closed")
	}

	// RFC 6120, Section 4.4 and 9.1.5
	log.Println("xmpp: sending closing stream tag")
	_, err := fmt.Fprint(c.out, "</stream:stream>")

	//TODO: find a better way to prevent sending message.
	c.out = ioutil.Discard

	if err != nil {
		return c.closeTCP()
	}

	// Wait for </stream:stream>
	select {
	case <-c.streamCloseReceived:
	case <-time.After(30 * time.Second):
		log.Println("xmpp: timed out waiting for closing stream")
	}

	return c.closeTCP()
}

func (c *conn) receivedStreamClose() error {
	log.Println("xmpp: received closing stream tag")
	return c.closeImmediately()
}

func (c *conn) closeImmediately() error {
	if c.closed {
		return nil
	}

	close(c.streamCloseReceived)
	return c.Close()
}

func (c *conn) closeTCP() error {
	if c.closed {
		return nil
	}

	//Close all pending requests at this moment. It will include pending pings
	c.cancelInflights()

	log.Println("xmpp: TCP closed")
	c.closed = true
	return c.rawOut.Close()
}

// Next reads stanzas from the server. If the stanza is a reply, it dispatches
// it to the correct channel and reads the next message. Otherwise it returns
// the stanza for processing.
func (c *conn) Next() (stanza data.Stanza, err error) {
	for {
		if stanza.Name, stanza.Value, err = next(c); err != nil {
			return
		}

		if _, ok := stanza.Value.(*data.StreamClose); ok {
			err = c.receivedStreamClose()
			return
		}

		if iq, ok := stanza.Value.(*data.ClientIQ); ok && (iq.Type == "result" || iq.Type == "error") {
			var cookieValue uint64
			if cookieValue, err = strconv.ParseUint(iq.ID, 16, 64); err != nil {
				err = errors.New("xmpp: failed to parse id from iq: " + err.Error())
				return
			}
			cookie := data.Cookie(cookieValue)

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
				if len(iq.From) > 0 && iq.From != c.jid && iq.From != utils.RemoveResourceFromJid(c.jid) && iq.From != utils.DomainFromJid(c.jid) {
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
func (c *conn) Cancel(cookie data.Cookie) bool {
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
func (c *conn) Send(to, msg string) error {
	archive := ""
	if !c.config.Archive {
		// The first part of archive is from google:
		// See https://developers.google.com/talk/jep_extensions/otr
		// The second part of the stanza is from XEP-0136
		// http://xmpp.org/extensions/xep-0136.html#pref-syntax-item-otr
		// http://xmpp.org/extensions/xep-0136.html#otr-nego
		archive = "<nos:x xmlns:nos='google:nosave' value='enabled'/><arc:record xmlns:arc='http://jabber.org/protocol/archive' otr='require'/>"
	}
	nocopy := "<no-copy xmlns='urn:xmpp:hints'/><no-permanent-store xmlns='urn:xmpp:hints'/><private xmlns='urn:xmpp:carbons:2'/>"
	_, err := fmt.Fprintf(c.out, "<message to='%s' from='%s' type='chat'><body>%s</body>%s%s</message>", xmlEscape(to), xmlEscape(c.jid), xmlEscape(msg), archive, nocopy)
	return err
}

// ReadStanzas reads XMPP stanzas
func (c *conn) ReadStanzas(stanzaChan chan<- data.Stanza) error {
	defer close(stanzaChan)
	defer c.closeImmediately()

	for {
		stanza, err := c.Next()
		if err != nil {
			log.Printf("xmpp: error receiving stanza. %s\n", err)
			return err
		}

		//The receiving entity has closed the channel
		if _, quit := stanza.Value.(*data.StreamClose); quit {
			return nil
		}

		stanzaChan <- stanza
	}
}
