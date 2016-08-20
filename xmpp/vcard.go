// Package xmpp implements the XMPP IM protocol, as specified in RFC 6120 and
// 6121.
package xmpp

import (
	"fmt"

	"github.com/twstrike/coyim/xmpp/data"
)

// RequestVCard requests the user's vcard from the server. It returns a
// channel on which the reply can be read when received and a Cookie that can
// be used to cancel the request.
func (c *conn) RequestVCard() (<-chan data.Stanza, data.Cookie, error) {
	cookie := c.getCookie()
	if _, err := fmt.Fprintf(c.out, "<iq type='get' id='%x'><vCard xmlns='vcard-temp'/></iq>", cookie); err != nil {
		return nil, 0, err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	ch := make(chan data.Stanza, 1)
	c.inflights[cookie] = inflight{ch, ""}
	return ch, cookie, nil
}
