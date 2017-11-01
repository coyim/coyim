package xmpp

import "fmt"

const (
	mucSupport = "<x xmlns='http://jabber.org/protocol/muc'/>"
)

func (c *conn) enterRoom(roomID, service, nickname string) error {
	occupantJID := fmt.Sprintf("%s@%s/%s", roomID, service, nickname)
	return c.sendPresenceWithChildren(occupantJID, "", "", mucSupport)
}

func (c *conn) leaveRoom(roomID, service, nickname string) error {
	occupantJID := fmt.Sprintf("%s@%s/%s", roomID, service, nickname)
	return c.sendPresenceWithChildren(occupantJID, "unavailable", "", mucSupport)
}
