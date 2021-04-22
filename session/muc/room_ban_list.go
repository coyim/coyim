package muc

import (
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// RoomBanListItem contains information about an item of the room's list of banned users
type RoomBanListItem struct {
	Jid         jid.Any
	Affiliation data.Affiliation
	Reason      string
}
