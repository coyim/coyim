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
	"reflect"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/cache"
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
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

	rand io.Reader
	lock sync.Mutex

	inflightsMutex sync.Mutex
	inflights      map[data.Cookie]inflight
	customStorage  map[xml.Name]reflect.Type

	lastPingRequest  time.Time
	lastPongResponse time.Time

	streamCloseReceived chan bool
	closed              bool
	closedLock          sync.Mutex

	c cache.WithExpiry

	statusUpdates chan<- string

	channelBinding []byte

	log coylog.Logger
}

func (c *conn) Cache() cache.WithExpiry {
	return c.c
}

func (c *conn) isClosed() bool {
	c.closedLock.Lock()
	defer c.closedLock.Unlock()
	return c.closed
}

func (c *conn) setClosed(v bool) {
	c.closedLock.Lock()
	defer c.closedLock.Unlock()
	c.closed = v
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
	l := log.New()
	l.SetOutput(ioutil.Discard)
	conn := &conn{
		in:     in,
		out:    out,
		rawOut: out,
		jid:    jid,

		inflights:           make(map[data.Cookie]inflight),
		streamCloseReceived: make(chan bool),

		rand: rand.Reader,
		c:    cache.NewWithExpiry(),
		log:  l,
	}

	conn.setClosed(true)
	close(conn.streamCloseReceived) // closes immediately
	return conn
}

func newConn() *conn {
	return &conn{
		inflights:           make(map[data.Cookie]inflight),
		streamCloseReceived: make(chan bool),
		c:                   cache.NewWithExpiry(),

		rand: rand.Reader,
	}
}

// Close closes the underlying connection
func (c *conn) Close() error {
	if c.isClosed() {
		return errors.New("xmpp: the connection is already closed")
	}

	// RFC 6120, Section 4.4 and 9.1.5
	c.log.Info("xmpp: sending closing stream tag")

	c.closedLock.Lock()
	_, err := fmt.Fprint(c.out, "</stream:stream>")

	//TODO: find a better way to prevent sending message.
	c.out = ioutil.Discard
	c.closedLock.Unlock()

	if err != nil {
		return c.closeTCP()
	}

	if c.statusUpdates != nil {
		c.statusUpdates <- "waitingForClose"
	}

	// Wait for </stream:stream>
	select {
	// Since no-one ever writes to the streamCloseReceived
	// channel, this select will wait not for a value, but for
	// the channel to be closed. In golang, reading from a closed
	// channel will immediately give you the zero-value for that
	// channel...
	case <-c.streamCloseReceived:
	case <-time.After(30 * time.Second):
		c.log.Info("xmpp: timed out waiting for closing stream")
	}

	return c.closeTCP()
}

func (c *conn) receivedStreamClose() error {
	c.log.Info("xmpp: received closing stream tag")
	return c.closeImmediately()
}

func (c *conn) closeImmediately() error {
	if c.isClosed() {
		return nil
	}

	close(c.streamCloseReceived)
	return c.Close()
}

func (c *conn) closeTCP() error {
	if c.isClosed() {
		return nil
	}

	//Close all pending requests at this moment. It will include pending pings
	c.cancelInflights()

	c.log.Info("xmpp: TCP closed")
	c.setClosed(true)
	return c.rawOut.Close()
}

func (c *conn) createInflight(cookie data.Cookie, to string) (<-chan data.Stanza, data.Cookie, error) {
	c.inflightsMutex.Lock()
	defer c.inflightsMutex.Unlock()

	ch := make(chan data.Stanza, 1)
	c.inflights[cookie] = inflight{ch, to}
	return ch, cookie, nil
}

func (c *conn) asyncReturnIQResponse(stanza data.Stanza) error {
	iq := stanza.Value.(*data.ClientIQ)
	cookieValue, err := strconv.ParseUint(iq.ID, 16, 64)
	if err != nil {
		return errors.New("xmpp: failed to parse id from iq: " + err.Error())
	}

	cookie := data.Cookie(cookieValue)
	c.inflightsMutex.Lock()
	inflight, ok := c.inflights[cookie]
	c.inflightsMutex.Unlock()

	if !ok {
		c.log.WithField("iq", iq).Warn("xmpp: received reply to unknown iq")
		return nil
	}

	if len(inflight.to) > 0 {
		// The reply must come from the address to
		// which we sent the request.
		if inflight.to != iq.From {
			return nil
		}
	} else {
		// If there was no destination on the request
		// then the matching is more complex because
		// servers differ in how they construct the
		// reply.
		bare := jid.Parse(c.jid).NoResource().String()
		dm := string(jid.Parse(c.jid).Host())
		if len(iq.From) > 0 && iq.From != c.jid && iq.From != bare && iq.From != dm {
			return nil
		}
	}

	c.inflightsMutex.Lock()
	delete(c.inflights, cookie)
	c.inflightsMutex.Unlock()

	//replyChan is buffered with size 1
	inflight.replyChan <- stanza
	return nil
}

// Next reads stanzas from the server. If the stanza is a reply, it dispatches
// it to the correct channel and reads the next message. Otherwise it returns
// the stanza for processing.
func (c *conn) Next() (stanza data.Stanza, err error) {
	for {
		if stanza.Name, stanza.Value, err = next(c, c.log); err != nil {
			return
		}

		if _, ok := stanza.Value.(*data.StreamClose); ok {
			err = c.receivedStreamClose()
			return
		}

		switch v := stanza.Value.(type) {
		case *data.ClientIQ:
			switch v.Type {
			case "result", "error":
				if err = c.asyncReturnIQResponse(stanza); err != nil {
					return
				}
			default:
				return
			}

		default:
			return
		}

	}
}

// Cancel cancels and outstanding request. The request's channel is closed.
func (c *conn) Cancel(cookie data.Cookie) bool {
	c.inflightsMutex.Lock()
	defer c.inflightsMutex.Unlock()

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
			c.log.WithError(err).Warn("xmpp: error receiving stanza")
			return err
		}

		//The receiving entity has closed the channel
		if _, quit := stanza.Value.(*data.StreamClose); quit {
			return nil
		}

		stanzaChan <- stanza
	}
}

// SetChannelBinding sets the specific value to use for channel binding with authentication modes that
// support this
func (c *conn) SetChannelBinding(v []byte) {
	c.channelBinding = v
}

// GetChannelBinding gets the specific value to use for channel binding with authentication modes that
// support this
func (c *conn) GetChannelBinding() []byte {
	return c.channelBinding
}
