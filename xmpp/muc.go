package xmpp

import "fmt"

const (
	mucSupport = "<x xmlns='http://jabber.org/protocol/muc'/>"
)

type Occupant struct {
	Room, Service, Nick string
}

func (o *Occupant) JID() string {
	return fmt.Sprintf("%s@%s/%s", o.Room, o.Service, o.Nick)
}

func (c *conn) enterRoom(roomID, service, nickname string) error {
	occupant := Occupant{Room: roomID, Service: service, Nick: nickname}
	return c.sendPresenceWithChildren(occupant.JID(), "", "", mucSupport)
}

func (c *conn) leaveRoom(roomID, service, nickname string) error {
	occupant := Occupant{Room: roomID, Service: service, Nick: nickname}
	return c.sendPresenceWithChildren(occupant.JID(), "unavailable", "", mucSupport)
}
