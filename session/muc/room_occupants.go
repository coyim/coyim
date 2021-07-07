package muc

import (
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// RoomOccupantItem contains information related with occupants to be configured in a room with a specific affiliation
type RoomOccupantItem struct {
	Jid           jid.Any
	Affiliation   data.Affiliation
	Reason        string
	MustBeUpdated bool
}

// ChangeAffiliationToNone changes an occupant's affiliation to none
func (roi *RoomOccupantItem) ChangeAffiliationToNone() {
	roi.Affiliation = &data.NoneAffiliation{}
}

// RoomOccupantItemList represents a list of room occupant items
type RoomOccupantItemList []*RoomOccupantItem

// IncludesJid returns a boolean that indicates if the given account ID (jid) is in the list
func (l RoomOccupantItemList) IncludesJid(id jid.Any) bool {
	for _, itm := range l {
		if itm.Jid.String() == id.String() {
			return true
		}
	}
	return false
}

// retrieveOccupantsToUpdate extracts occupants from list when MustBeUpdated attribute is true
func (l RoomOccupantItemList) retrieveOccupantsToUpdate() RoomOccupantItemList {
	extracted := RoomOccupantItemList{}
	for _, o := range l {
		if o.MustBeUpdated {
			extracted = append(extracted, o)
		}
	}
	return extracted
}
