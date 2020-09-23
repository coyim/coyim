package session

import (
	"fmt"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func newMUCRoomOccupant(nickname string, affiliation muc.Affiliation, role muc.Role, realJid jid.Full) *muc.Occupant {
	return &muc.Occupant{
		Nick:        nickname,
		Affiliation: affiliation,
		Role:        role,
		Jid:         realJid,
	}
}

func (m *mucManager) handleOccupantUpdate(roomID jid.Bare, occupant *muc.Occupant) {
	// TODO: why does this ignore the second argument?
	room, _ := m.roomManager.GetRoom(roomID)
	joined, _, err := room.Roster().UpdatePresence(occupant, "")

	if err != nil {
		// TODO: We should probably add fields for room and occupant here
		m.log.WithError(err).Error("Error on trying to update the occupant status in the roster")
		return
	}

	// TODO: It is a bit confusing that we publish events in both the muc_events.go and this file.
	// TODO: Joining should probably happen before updating
	m.occupantUpdate(roomID, occupant)

	if joined {
		m.occupantJoined(roomID, occupant)
	}
}

func (m *mucManager) handleOccupantLeft(roomID jid.Bare, occupant *muc.Occupant) {
	// TODO: Same, why ignoring second result?
	r, _ := m.roomManager.GetRoom(roomID)
	// TODO: We should check whether the person ACTUALLY left or not.
	occupant.UpdateStatus("unavailable", "Occupant left the room")
	_, _, err := r.Roster().UpdatePresence(occupant, "unavailable")

	if err != nil {
		m.log.WithError(err).Error("An error occurred trying to remove the occupant from the roster")
		return
	}

	m.occupantLeft(roomID, occupant)
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
