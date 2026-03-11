package session

import (
	"github.com/coyim/coyim/internal/util"
	"github.com/coyim/coyim/xmpp/data"
)

// RemoveContact is used to remove a contact
func (s *session) RemoveContact(jid string) {
	_, _, e := s.conn.SendIQ("" /* to the server */, "set", data.RosterRequest{
		Item: data.RosterRequestItem{
			Jid:          jid,
			Subscription: "remove",
		},
	})

	util.LogIgnoredError(e, s.log, "removing contact")
}
