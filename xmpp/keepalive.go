package xmpp

import (
	"encoding/xml"
	"io"
	"log"
	"net"
	"time"
)

var (
	keepaliveInterval = 10 * time.Second
	keepaliveTimeout  = 30 * time.Second
)

// Manage whitespace keepalives as specified in RFC 6120, section 4.6.1
func (c *Conn) watchKeepAlive(conn net.Conn) {
	tick := time.NewTicker(keepaliveInterval)
	defer tick.Stop()
	defer log.Println("xmpp: no more watching keepalives")

	for range tick.C {
		if c.closed {
			return
		}

		if c.sendKeepalive() {
			log.Println("xmpp: keepalive sent")
			continue
		}

		log.Println("xmpp: failed to send keepalive")

		go c.sendStreamError(StreamError{
			DefinedCondition: ConnectionTimeout,
		})

		return
	}
}

func (c *Conn) sendKeepalive() bool {
	_, err := c.keepaliveOut.Write([]byte{0x20})
	return c.closed || err == nil || err == io.EOF
}

func (c *Conn) sendStreamError(streamError StreamError) error {
	enc, err := xml.Marshal(streamError)
	if err != nil {
		return err
	}

	//This is expected to error since the connection may be unreliable at this moment
	c.out.Write(enc)

	return c.Close()
}
