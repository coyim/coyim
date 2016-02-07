package session

import (
	"github.com/twstrike/coyim/xmpp"
)

// RemoveContact is used to remove a contact
//TODO: double check how the RFC calls this
func (s *Session) RemoveContact(jid string) {
	s.conn.SendIQ("" /* to the server */, "set", xmpp.RosterRequest{
		Item: xmpp.RosterRequestItem{
			Jid:          jid,
			Subscription: "remove",
		},
	})
}
