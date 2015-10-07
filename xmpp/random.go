package xmpp

import (
	"crypto/rand"
	"io"
)

func (c *Conn) rand() io.Reader {
	if c.Rand != nil {
		return c.Rand
	}
	return rand.Reader
}
