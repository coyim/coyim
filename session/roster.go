package session

import "github.com/coyim/coyim/xmpp/data"

// RemoveContact is used to remove a contact
//TODO: double check how the RFC calls this
func (s *session) RemoveContact(jid string) {
	s.conn.SendIQ("" /* to the server */, "set", data.RosterRequest{
		Item: data.RosterRequestItem{
			Jid:          jid,
			Subscription: "remove",
		},
	})
}
