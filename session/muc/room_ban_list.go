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

// IsDifferentFrom returns a boolean indicating if the ban list item is different in any
// of its values from the given information
func (itm *RoomBanListItem) IsDifferentFrom(itm2 *RoomBanListItem) bool {
	return itm.Jid.String() != itm2.Jid.String() ||
		itm.Affiliation.IsDifferentFrom(itm2.Affiliation) ||
		itm.Reason != itm2.Reason
}

// RoomBanList represents a list of banned users items
type RoomBanList []*RoomBanListItem

// IncludesJid returns a boolean that indicates if the given account ID (jid) is in the list
func (l RoomBanList) IncludesJid(id string) bool {
	for _, itm := range l {
		if itm.Jid.String() == id {
			return true
		}
	}
	return false
}
