package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (a *account) getRoomView(ident jid.Bare) (*roomView, bool) {
	room, ok := a.roomManager.GetRoom(ident)
	if !ok {
		a.log.WithField("room", ident).Error("Trying to get a room that is not in the manager")
		return nil, false
	}
	return getViewFromRoom(room), true
}

func (a *account) onRoomNicknameConflict(room jid.Bare, nickname string) {
	view, ok := a.getRoomView(room)
	if !ok {
		a.log.WithField("room", room).Error("Room view not available when a room nickname conflict event was received")
		return
	}

	view.onNicknameConflictReceived(room, nickname)
}

func (a *account) handleMUCLoggingEnabled(room jid.Bare) {
	view, ok := a.getRoomView(room)
	if !ok {
		a.log.WithField("room", room).Error("Not possible to get room view when handling Logging Enabled event")
		return
	}

	view.loggingIsEnabled()
}

func (a *account) handleMUCLoggingDisabled(room jid.Bare) {
	view, ok := a.getRoomView(room)
	if !ok {
		a.log.WithField("room", room).Error("Not possible to get room view when handling Logging Disabled event")
		return
	}

	view.loggingIsDisabled()
}

func (a *account) onRoomRegistrationRequired(room jid.Bare, nickname string) {
	view, ok := a.getRoomView(room)
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

	room, ok := a.roomManager.GetRoom(roomName)
	if !ok {
		l.Error("Room view not available when a room occupant joined event was received")
		return
	}

	view := getViewFromRoom(room)

	roster := room.Roster()
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

	view.onRoomOccupantJoinedReceived(nickname, roster.AllOccupants())
}

func (a *account) onRoomOccupantUpdated(roomName jid.Bare, nickname string, occupant jid.Full, affiliation muc.Affiliation, role muc.Role) {
	l := a.log.WithFields(log.Fields{
		"room":        roomName,
		"nickname":    nickname,
		"occupant":    occupant,
		"affiliation": affiliation,
		"role":        role,
	})

	room, ok := a.roomManager.GetRoom(roomName)
	if !ok {
		l.Error("Room view not available when a room occupant updated event was received")
		return
	}

	roster := room.Roster()
	_, _, err := roster.UpdatePresence(roomName.WithResource(jid.NewResource(nickname)), "", affiliation, role, "", "", "Occupant updated", occupant)
	if err != nil {
		l.WithError(err).Error("Error on trying to update the occupant status in the roster")
		return
	}

	view := getViewFromRoom(room)
	view.onRoomOccupantUpdateReceived(roster.AllOccupants())
}

func (a *account) onRoomOccupantLeftTheRoom(roomName jid.Bare, nickname string, ident jid.Full, affiliation muc.Affiliation, role muc.Role) {
	l := a.log.WithFields(log.Fields{
		"room":        roomName,
		"nickname":    nickname,
		"occupant":    ident,
		"affiliation": affiliation,
		"role":        role,
	})

	room, ok := a.roomManager.GetRoom(roomName)
	if !ok {
		l.Error("Room view not available when an occupant left the room event was received")
		return
	}

	roster := room.Roster()
	_, left, err := roster.UpdatePresence(roomName.WithResource(jid.NewResource(nickname)), "unavailable", affiliation, role, "", "unavailable", "Occupant left the room", ident)
	if err != nil {
		l.WithError(err).Error("An error occurred trying to remove the occupant from the roster")
		return
	}

	if !left {
		l.Error("The occupant can't left the room roster")
		return
	}

	view := getViewFromRoom(room)
	view.onRoomOccupantLeftTheRoomReceived(ident.Resource(), roster.AllOccupants())
}

func (a *account) onRoomMessageReceived(roomName jid.Bare, nickname jid.Resource, message string) {
	l := a.log.WithFields(log.Fields{
		"roomName": roomName,
		"nickname": nickname,
		"message":  message,
	})

	room, ok := a.roomManager.GetRoom(roomName)
	if !ok {
		l.Error("Room view not available when a live message was received")
		return
	}

	view := getViewFromRoom(room)
	view.onRoomMessageToTheRoomReceived(nickname, message)
}
