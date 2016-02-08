package xmpp

import (
	"crypto/rand"
	"io"
)

func (c *conn) Rand() io.Reader {
	if c.rand != nil {
		return c.rand
	}
	return rand.Reader
}
