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

	joined, _, err := room.Roster().UpdatePresence(occupant, "")
	if err != nil {
		l.WithError(err).Error("An error occurred trying to update the occupant status in the roster")
		return
	}

	if joined {
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

	// TODO: Check this functionality, create two functions (join and remove) from roster
	occupant.UpdateStatus("unavailable", "Occupant left the room")
	_, left, err := r.Roster().UpdatePresence(occupant, "unavailable")
	if err != nil {
		m.log.WithError(err).Error("An error occurred trying to remove the occupant from the roster")
		return
	}

	if left {
		m.occupantLeft(roomID, occupant)
	}
}
