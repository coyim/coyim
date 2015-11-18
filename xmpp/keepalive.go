package xmpp

import (
	"encoding/xml"
	"log"
	"net"
	"time"
)

var (
	keepaliveInterval = 10 * time.Second
	keepaliveTimeout  = 30 * time.Second
)

// Manage whitespace keepalives as specified in RFC 6210, section 4.6.1
func (c *Conn) watchKeepAlive(conn net.Conn) {
	for !c.closed {
		<-time.After(keepaliveInterval)

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
	return err == nil
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
