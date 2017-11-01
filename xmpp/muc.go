package xmpp

import "fmt"

const (
	mucSupport = "<x xmlns='http://jabber.org/protocol/muc'/>"
)

func (c *conn) enterRoom(roomID, service, nickname string) error {
	to := fmt.Sprintf("%s@%s/%s", roomID, service, nickname)
	return c.sendPresenceWithChildren(to, "", "", mucSupport)
}
