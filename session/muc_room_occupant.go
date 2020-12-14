package session

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func newMUCRoomOccupant(nickname string, affiliation data.Affiliation, role data.Role, realJid jid.Full) *muc.Occupant {
	return &muc.Occupant{
		Nickname:    nickname,
		Affiliation: affiliation,
		Role:        role,
		RealJid:     realJid,
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
	// Added IsSelfOccupantInTheRoom validation to avoid publishing the events of
	// other occupants until receive the selfPresence.
	// This validation is temporally while 'state machine' pattern is implemented.
	if room.IsSelfOccupantInTheRoom() {
		if !updated {
			m.occupantJoined(roomID, op)
		}
		m.occupantUpdate(roomID, op)
	}
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
		l.WithError(err).Error("An error occurred trying to remove the occupant from the roster")
		return
	}

	m.occupantLeft(roomID, op)
}

func (m *mucManager) handleOccupantUnavailable(roomID jid.Bare, op *muc.OccupantPresenceInfo, u *xmppData.MUCUser) {
	if u == nil || u.Destroy == nil {
		return
	}

	m.handleRoomDestroyed(roomID, u.Destroy)
}

func (m *mucManager) handleRoomDestroyed(roomID jid.Bare, d *xmppData.MUCRoomDestroy) {
	j, ok := jid.TryParseBare(d.Jid)
	if d.Jid != "" && !ok {
		m.log.WithFields(log.Fields{
			"room":            roomID,
			"alternativeRoom": d.Jid,
			"method":          "handleRoomDestroyed",
		}).Warn("Invalid alternative room ID")
	}

	m.roomDestroyed(roomID, d.Reason, j, d.Password)
}

func (m *mucManager) handleNonMembersRemoved(roomID jid.Bare, op *muc.OccupantPresenceInfo) {
	l := m.log.WithFields(log.Fields{
		"room":     roomID,
		"occupant": op.Nickname,
		"method":   "handleNonMembersRemoved",
	})

	r, ok := m.roomManager.GetRoom(roomID)
	if !ok {
		l.Error("Trying to get a room that is not in the room manager")
		return
	}

	err := r.Roster().RemoveOccupant(op.Nickname)
	if err != nil {
		l.WithError(err).Error("An error occurred trying to remove the occupant from the roster")
	}

	if r.SelfOccupant().Nickname == op.Nickname {
		m.removeSelfOccupant(roomID)
		_ = m.roomManager.LeaveRoom(roomID)
		return
	}
	m.occupantRemoved(roomID, op.Nickname)
}
