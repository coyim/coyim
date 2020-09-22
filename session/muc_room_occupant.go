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
		return o.realJid.String() == from.String()
	}
	return false
}

func (m *mucManager) occupantUpdate(from jid.Full, room jid.Bare, occupant mucRoomOccupant) {
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
		m.log.WithError(err).Error("Error on trying to update the occupant status in the roster")
		return
	}

	m.publishOccupantUpdate(from, room, occupant)

	if joined {
		m.publishOccupantJoined(from, room, occupant)
	}
}

func (m *mucManager) occupantLeft(from jid.Full, room jid.Bare, occupant mucRoomOccupant) {
	r, _ := m.roomManager.GetRoom(room)
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

func parseAffiliationAndReport(a string) muc.Affiliation {
	aa, e := muc.AffiliationFromString(a)
	if e != nil {
		fmt.Printf("error when parsing affiliation: %v\n", e)
	}
	return aa
}

func parseRoleAndReport(a string) muc.Role {
	aa, e := muc.RoleFromString(a)
	if e != nil {
		fmt.Printf("error when parsing role: %v\n", e)
	}
	return aa
}
