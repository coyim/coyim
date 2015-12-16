package session

import "../xmpp"

//TODO: double check how the RFC calls this
func (s *Session) RemoveContact(jid string) {
	s.Conn.SendIQ("" /* to the server */, "set", xmpp.RosterRequest{
		Item: xmpp.RosterRequestItem{
			Jid:          jid,
			Subscription: "remove",
		},
	})
}
