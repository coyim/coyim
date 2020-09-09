package gui

import (
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

func (a *account) onRoomNicknameConflict(from jid.Full) {
	view, ok := a.getRoomView(from.Bare())
	if !ok {
		a.log.WithField("from", from).Error("Room view not available when a room nickname conflict event was received")
		return
	}

	view.onNicknameConflictReceived(from)
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

func (a *account) onRoomRegistrationRequired(from jid.Full) {
	view, ok := a.getRoomView(from.Bare())
	if !ok {
		a.log.WithField("from", from).Error("Room view not available when a room registration required event was received")
		return
	}

	view.onRegistrationRequiredReceived(from)
}

func (a *account) onRoomOccupantJoined(from jid.Full, ident jid.Full, affiliation, role, status string) {
	l := a.log.WithFields(log.Fields{
		"from":        from,
		"occupant":    ident,
		"affiliation": affiliation,
		"role":        role,
		"status":      status,
	})

	room, ok := a.roomManager.GetRoom(from.Bare())
	if !ok {
		l.Error("Room view not available when a room occupant joined event was received")
		return
	}

	view := getViewFromRoom(room)

	roster := room.Roster()
	joined, _, err := roster.UpdatePresence(from, "", affiliation, role, "", status, "Occupant joined", ident)
	if err != nil {
		l.WithError(err).Error("An error occurred trying to add the occupant to the roster")
		view.onRoomOccupantErrorReceived(from)
		return
	}

	if !joined {
		l.Error("The occupant can't join the room roster")
		return
	}

	view.onRoomOccupantJoinedReceived(ident.Resource(), roster.AllOccupants())
}

func (a *account) onRoomOccupantUpdated(from jid.Full, occupant jid.Full, affiliation, role string) {
	l := a.log.WithFields(log.Fields{
		"from":        from,
		"occupant":    occupant,
		"affiliation": affiliation,
		"role":        role,
	})

	room, ok := a.roomManager.GetRoom(from.Bare())
	if !ok {
		l.Error("Room view not available when a room occupant updated event was received")
		return
	}

	roster := room.Roster()
	_, _, err := roster.UpdatePresence(from, "", affiliation, role, "", "", "Occupant updated", occupant)
	if err != nil {
		l.WithError(err).Error("Error on trying to update the occupant status in the roster")
		return
	}

	view := getViewFromRoom(room)
	view.onRoomOccupantUpdateReceived(roster.AllOccupants())
}

func (a *account) onRoomOccupantLeftTheRoom(from jid.Full, ident jid.Full, affiliation, role string) {
	l := a.log.WithFields(log.Fields{
		"from":        from,
		"occupant":    ident,
		"affiliation": affiliation,
		"role":        role,
	})

	room, ok := a.roomManager.GetRoom(from.Bare())
	if !ok {
		l.Error("Room view not available when an occupant left the room event was received")
		return
	}

	roster := room.Roster()
	_, left, err := roster.UpdatePresence(from, "unavailable", affiliation, role, "", "unavailable", "Occupant left the room", ident)
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
