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

// handleOccupantAffiliationUpdate will lock until the update process of the given occupant finishes
func (m *mucManager) handleOccupantAffiliationUpdate(roomID jid.Bare, presence *muc.OccupantPresenceInfo, isOwnPresence bool) {
	m.occupantAffiliationUpdateLock.Lock()
	defer m.occupantAffiliationUpdateLock.Unlock()

	room, ok := m.roomManager.GetRoom(roomID)
	if !ok {
		m.log.WithFields(log.Fields{
			"room":     roomID,
			"occupant": presence.Nickname,
			"method":   "handleOccupantAffiliationUpdate",
		}).Error("Trying to get a room that is not in the room manager")
		return
	}

	occupant, exist := room.Roster().GetOccupant(presence.Nickname)
	if exist {
		actorOccupant := &data.OccupantUpdateActor{
			Nickname: presence.AffiliationRole.Actor,
		}

		actor, ok := room.Roster().GetOccupant(presence.AffiliationRole.Actor)
		if ok {
			actorOccupant.Affiliation = actor.Affiliation
			actorOccupant.Role = actor.Role
		}

		commonUpdateInfo := data.OccupantUpdateAffiliationRole{
			Nickname: presence.Nickname,
			Actor:    actorOccupant,
			Reason:   presence.AffiliationRole.Reason,
		}

		switch {
		case occupant.Affiliation.Name() != presence.AffiliationRole.Affiliation.Name():
			affiliationUpate := data.AffiliationUpdate{
				OccupantUpdateAffiliationRole: commonUpdateInfo,
				New:                           presence.AffiliationRole.Affiliation,
				Previous:                      occupant.Affiliation,
			}

			occupant.UpdateAffiliation(presence.AffiliationRole.Affiliation)

			if isOwnPresence {
				m.selfOccupantAffiliationUpdated(roomID, affiliationUpate)
				return
			}

			m.occupantAffiliationUpdated(roomID, affiliationUpate)

			break
		case occupant.Role.Name() != presence.AffiliationRole.Role.Name():
			roleUpdate := data.RoleUpdate{
				OccupantUpdateAffiliationRole: commonUpdateInfo,
				New:                           presence.AffiliationRole.Role,
				Previous:                      occupant.Role,
			}

			occupant.UpdateRole(presence.AffiliationRole.Role)

			if isOwnPresence {
				m.selfOccupantRoleUpdated(roomID, roleUpdate)
				return
			}

			m.occupantRoleUpdated(roomID, roleUpdate)

			break
		}
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
