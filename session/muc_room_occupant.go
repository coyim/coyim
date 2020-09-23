package session

import (
	"fmt"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

type mucRoomOccupant struct {
	nickname    string
	realJid     jid.Full
	affiliation muc.Affiliation
	role        muc.Role
}

// TODO: this should return a pointer instead. No point in copying structs all over the place

// TODO: we need to think about this struct, it is a bit weird to have this and muc.Occupant at the same time
// The problem is that UpdatePresence on the RoomRoster is responsible for creating the Occupant if it does not
// exist. We need to think about whether this should be refactored or what makes sense.

func newMUCRoomOccupant(nickname jid.Resource, affiliation, role string, realJid jid.Full) mucRoomOccupant {
	return mucRoomOccupant{
		nickname:    nickname.String(),
		affiliation: parseAffiliationAndReport(affiliation),
		role:        parseRoleAndReport(role),
		realJid:     realJid,
	}
}

func (o *mucRoomOccupant) sameFrom(from jid.Full) bool {
	if o.realJid != nil {
		// TODO: Maybe better to implement an Equals method on the jid objects
		return o.realJid.String() == from.String()
	}
	return false
}

// TODO: This method should not take both "from" and "room" and "occupant"

func (m *mucManager) occupantUpdate(from jid.Full, room jid.Bare, occupant mucRoomOccupant) {
	// TODO: why does this ignore the second argument?
	r, _ := m.roomManager.GetRoom(room)
	joined, _, err := r.Roster().UpdatePresence(
		room.WithResource(jid.NewResource(occupant.nickname)),
		"",
		occupant.affiliation,
		occupant.role,
		"",
		"",
		"Occupant updated",
		occupant.realJid,
	)

	if err != nil {
		// TODO: We should probably add fields for room and occupant here
		m.log.WithError(err).Error("Error on trying to update the occupant status in the roster")
		return
	}

	// TODO: It is a bit confusing that we publish events in both the muc_events.go and this file.
	// TODO: Joining should probably happen before updating
	m.publishOccupantUpdate(from, room, occupant)

	if joined {
		m.publishOccupantJoined(from, room, occupant)
	}
}

func (m *mucManager) occupantLeft(from jid.Full, room jid.Bare, occupant mucRoomOccupant) {
	// TODO: Same, why ignoring second result?
	r, _ := m.roomManager.GetRoom(room)
	// TODO: We should check whether the person ACTUALLY left or not.
	_, _, err := r.Roster().UpdatePresence(
		room.WithResource(jid.NewResource(occupant.nickname)),
		"unavailable",
		occupant.affiliation,
		occupant.role,
		"",
		"unavailable",
		"Occupant left the room",
		occupant.realJid,
	)

	if err != nil {
		m.log.WithError(err).Error("An error occurred trying to remove the occupant from the roster")
		return
	}

	m.publishOccupantLeft(from, room, occupant)
}

// TODO: Maybe move these and other utility methods to the muc.go file, and move
// muc_manager and other more focused functions and methods from muc.go into new files.

func parseAffiliationAndReport(a string) muc.Affiliation {
	aa, e := muc.AffiliationFromString(a)
	if e != nil {
		// We use printf here because this is a programmer error, and should not happen in
		// production code.
		fmt.Printf("error when parsing affiliation: %v\n", e)
	}
	return aa
}

func parseRoleAndReport(a string) muc.Role {
	aa, e := muc.RoleFromString(a)
	if e != nil {
		// We use printf here because this is a programmer error, and should not happen in
		// production code.
		fmt.Printf("error when parsing role: %v\n", e)
	}
	return aa
}
