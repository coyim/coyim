package xmpp

import (
	"encoding/xml"
	"io"
	"time"

	"github.com/coyim/coyim/xmpp/data"
)

var (
	keepaliveInterval = 10 * time.Second
	keepaliveTimeout  = 30 * time.Second
)

// Manage whitespace keepalives as specified in RFC 6120, section 4.6.1
func (c *conn) watchKeepAlive() {
	tick := time.NewTicker(keepaliveInterval)
	defer tick.Stop()
	defer c.log.Info("xmpp: no more watching keepalives")

	for range tick.C {
		if c.closed {
			return
		}

		if c.sendKeepalive() {
			continue
		}

		c.log.Info("xmpp: keepalive failed")

		go func() {
			_ = c.sendStreamError(data.StreamError{
				DefinedCondition: data.ConnectionTimeout,
			})
		}()

		return
	}
}

func (c *conn) sendKeepalive() bool {
	_, err := c.keepaliveOut.Write([]byte{0x20})
	return c.closed || err == nil || err == io.EOF
}

func (c *conn) sendStreamError(streamError data.StreamError) error {
	enc, _ := xml.Marshal(streamError)

	//This is expected to error since the connection may be unreliable at this moment
	_, _ = c.safeWrite(enc)

	return c.Close()
}
