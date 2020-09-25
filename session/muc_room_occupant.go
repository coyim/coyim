package session

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
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
	l := m.log.WithFields(log.Fields{
		"room":     roomID,
		"occupant": occupant.Nick,
		"method":   "handleOccupantUpdate",
	})

	room, ok := m.roomManager.GetRoom(roomID)
	if !ok {
		l.Error("Trying to get a room that is not in the room manager")
		return
	}

	updated := room.Roster().UpdateOrAddOccupant(occupant)
	if !updated {
		m.occupantJoined(roomID, occupant)
	}
	m.occupantUpdate(roomID, occupant)
}

func (m *mucManager) handleOccupantLeft(roomID jid.Bare, occupant *muc.Occupant) {
	l := m.log.WithFields(log.Fields{
		"room":     roomID,
		"occupant": occupant.Nick,
		"method":   "handleOccupantLeft",
	})

	r, ok := m.roomManager.GetRoom(roomID)
	if !ok {
		l.Error("Trying to get a room that is not in the room manager")
		return
	}

	occupant.UpdateStatus("unavailable", "Occupant left the room")

	err := r.Roster().RemoveOccupant(occupant)
	if err != nil {
		m.log.WithError(err).Error("An error occurred trying to remove the occupant from the roster")
		return
	}

	m.occupantLeft(roomID, occupant)
}
