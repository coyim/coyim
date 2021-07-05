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

// RoomOccupantItemList represents a list of room occupant items
type RoomOccupantItemList []*RoomOccupantItem

// IncludesJid returns a boolean that indicates if the given account ID (jid) is in the list
func (l RoomOccupantItemList) IncludesJid(id string) bool {
	for _, itm := range l {
		if itm.Jid.String() == id {
			return true
		}
	}
	return false
}
