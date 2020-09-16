package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (a *account) getRoomView(roomIdentity string) (*roomView, bool) {
	room, ok := a.rooms[roomIdentity]
	return room, ok
}

func (a *account) addRoomView(v *roomView) {
	a.rooms[v.identity.String()] = v
}

func (a *account) newRoomModel(identity jid.Bare) *muc.Room {
	return a.session.NewRoom(identity)
}

func (a *account) onRoomNicknameConflict(room jid.Bare, nickname string) {
	view, ok := a.getRoomView(room.String())
	if !ok {
		a.log.WithField("room", room).Error("Room view not available when a room nickname conflict event was received")
		return
	}

	view.onNicknameConflictReceived(room, nickname)
}

func (a *account) handleMUCLoggingEnabled(room jid.Bare) {
	view, ok := a.getRoomView(room.String())
	if !ok {
		a.log.WithField("room", room).Error("Not possible to get room view when handling Logging Enabled event")
		return
	}

	view.loggingIsEnabled()
}

func (a *account) handleMUCLoggingDisabled(room jid.Bare) {
	view, ok := a.getRoomView(room.String())
	if !ok {
		a.log.WithField("room", room).Error("Not possible to get room view when handling Logging Disabled event")
		return
	}

	view.loggingIsDisabled()
}

func (a *account) onRoomRegistrationRequired(room jid.Bare, nickname string) {
	view, ok := a.getRoomView(room.String())
	if !ok {
		a.log.WithField("room", room).Error("Room view not available when a room registration required event was received")
		return
	}

	view.onRegistrationRequiredReceived(room, nickname)
}

func (a *account) onRoomOccupantJoined(roomName jid.Bare, nickname string, ident jid.Full, affiliation muc.Affiliation, role muc.Role, status string) {
	l := a.log.WithFields(log.Fields{
		"room":        roomName,
		"nickname":    nickname,
		"occupant":    ident,
		"affiliation": affiliation,
		"role":        role,
		"status":      status,
	})

	view, ok := a.getRoomView(roomName.String())
	if !ok {
		l.Error("Room view not available when a room occupant joined event was received")
		return
	}

	roster := view.room.Roster()
	joined, _, err := roster.UpdatePresence(roomName.WithResource(jid.NewResource(nickname)), "", affiliation, role, "", status, "Occupant joined", ident)
	if err != nil {
		l.WithError(err).Error("An error occurred trying to add the occupant to the roster")
		view.onRoomOccupantErrorReceived(roomName, nickname)
		return
	}

	if !joined {
		l.Error("The occupant can't join the room roster")
		return
	}

	view.onRoomOccupantJoinedReceived(nickname)
}

func (a *account) onRoomOccupantUpdated(roomName jid.Bare, nickname string, occupant jid.Full, affiliation muc.Affiliation, role muc.Role) {
	l := a.log.WithFields(log.Fields{
		"room":        roomName,
		"nickname":    nickname,
		"occupant":    occupant,
		"affiliation": affiliation,
		"role":        role,
	})

	view, ok := a.getRoomView(roomName.String())
	if !ok {
		l.Error("Room view not available when a room occupant updated event was received")
		return
	}

	roster := view.room.Roster()
	joined, _, err := roster.UpdatePresence(roomName.WithResource(jid.NewResource(nickname)), "", affiliation, role, "", "", "Occupant updated", occupant)
	if err != nil {
		l.WithError(err).Error("Error on trying to update the occupant status in the roster")
		return
	}

	if joined {
		view.someoneJoinedTheRoom(nickname)
	} else {
		view.onRoomOccupantUpdateReceived()
	}
}

func (a *account) onRoomOccupantLeftTheRoom(roomName jid.Bare, nickname string, ident jid.Full, affiliation muc.Affiliation, role muc.Role) {
	l := a.log.WithFields(log.Fields{
		"room":        roomName,
		"nickname":    nickname,
		"occupant":    ident,
		"affiliation": affiliation,
		"role":        role,
	})

	view, ok := a.getRoomView(roomName.String())
	if !ok {
		l.Error("Room view not available when an occupant left the room event was received")
		return
	}

	roster := view.room.Roster()
	_, left, err := roster.UpdatePresence(roomName.WithResource(jid.NewResource(nickname)), "unavailable", affiliation, role, "", "unavailable", "Occupant left the room", ident)
	if err != nil {
		l.WithError(err).Error("An error occurred trying to remove the occupant from the roster")
		return
	}

	if !left {
		l.Error("The occupant can't left the room roster")
		return
	}

	view.onRoomOccupantLeftTheRoomReceived(nickname)
}

func (a *account) onRoomMessageReceived(roomName jid.Bare, nickname, subject, message string) {
	l := a.log.WithFields(log.Fields{
		"room":     roomName,
		"nickname": nickname,
		"subject":  subject,
		"message":  message,
	})

	view, ok := a.getRoomView(roomName.String())
	if !ok {
		l.Error("Room view not available when a live message was received")
		return
	}

	view.onRoomMessageToTheRoomReceived(nickname, subject, message)
}
