package muc

import (
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// RoomOccupantItem contains information related with occupants to be configured in a room with a specific affiliation
type RoomOccupantItem struct {
	Jid         jid.Any
	Affiliation data.Affiliation
	Reason      string
}
