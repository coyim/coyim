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

func (m *mucManager) handleOccupantUpdate(roomID jid.Bare, op *muc.OccupantPresenceInfo) {
	l := m.log.WithFields(log.Fields{
		"room":     roomID,
		"occupant": op.Nickname,
		"method":   "handleOccupantUpdate",
	})

	room, ok := m.roomManager.GetRoom(roomID)
	if !ok {
		l.Error("Trying to get a room that is not in the room manager")
		return
	}

	updated := room.Roster().UpdateOrAddOccupant(op)
	if !updated {
		m.occupantJoined(roomID, op)
	}
	m.occupantUpdate(roomID, op)
}

func (m *mucManager) handleOccupantLeft(roomID jid.Bare, op *muc.OccupantPresenceInfo) {
	l := m.log.WithFields(log.Fields{
		"room":     roomID,
		"occupant": op.Nickname,
		"method":   "handleOccupantLeft",
	})

	r, ok := m.roomManager.GetRoom(roomID)
	if !ok {
		l.Error("Trying to get a room that is not in the room manager")
		return
	}

	err := r.Roster().RemoveOccupant(op.Nickname)
	if err != nil {
		m.log.WithError(err).Error("An error occurred trying to remove the occupant from the roster")
		return
	}

	m.occupantLeft(roomID, op)
}
