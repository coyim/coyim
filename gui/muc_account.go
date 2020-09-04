package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (a *account) getRoomByIdentity(ident jid.Bare) (*muc.Room, error) {
	room, exists := a.roomManager.GetRoom(ident)
	if !exists {
		return nil, newRoomNotExistsError(ident)
	}
	return room, nil
}

// getRoom will return a specific room based on it's bare.
// If the room doesn't exists this method will return false
func (a *account) getRoom(ident jid.Bare) (*muc.Room, bool) {
	room, err := a.getRoomByIdentity(ident)
	if err != nil {
		a.log.WithField("room", ident).WithError(err).Error("Trying to get a room that is not in the manager")
		return nil, false
	}
	return room, true
}

// getRoomView will return a specific room view based on it's bare.
// If the view doesn't exists this method will return false
func (a *account) getRoomView(ident jid.Bare) (*roomView, bool) {
	room, ok := a.getRoom(ident)
	if !ok {
		a.log.WithField("room", ident).Error("Trying to get a room that is not in the manager")
		return nil, false
	}
	return getViewFromRoom(room), true
}

func (a *account) onRoomNicknameConflict(from jid.Full) {
	view, ok := a.getRoomView(from.Bare())
	if ok {
		view.onNicknameConflictReceived(from)
	}
}

func (a *account) onRoomRegistrationRequired(from jid.Full) {
	view, ok := a.getRoomView(from.Bare())
	if ok {
		view.onRegistrationRequiredReceived(from)
	}
}

func (a *account) onRoomOccupantJoined(from jid.Full, ident jid.Full, affiliation, role, status string) {
	room, ok := a.getRoom(from.Bare())
	if !ok {
		return
	}

	view := getViewFromRoom(room)

	roster := room.Roster()
	joined, _, err := roster.UpdatePresence(from, "", affiliation, role, "", status, "Occupant joined", ident)
	if err != nil {
		a.log.WithFields(log.Fields{
			"from":        from,
			"occupant":    ident,
			"affiliation": affiliation,
			"role":        role,
			"status":      status,
		}).WithError(err).Error("An error occurred trying to add the occupant to the roster")

		view.onRoomOccupantErrorReceived(from)
		return
	}

	// TODO: we are receiving a `join the room` event, so,
	// if !joined we should log it
	if joined {
		view.onRoomOccupantJoinedReceived(ident.Resource(), roster.AllOccupants())
	}
}

func (a *account) onRoomOccupantUpdated(from jid.Full, occupant jid.Full, affiliation, role string) {
	room, ok := a.getRoom(from.Bare())
	if !ok {
		return
	}

	roster := room.Roster()
	_, _, err := roster.UpdatePresence(from, "", affiliation, role, "", "", "Occupant updated", occupant)
	if err != nil {
		a.log.WithFields(log.Fields{
			"from":        from,
			"occupant":    occupant,
			"affiliation": affiliation,
		}).WithError(err).Error("Error on trying to update the occupant status in the roster")
		return
	}

	view := getViewFromRoom(room)
	view.onRoomOccupantUpdateReceived(roster.AllOccupants())
}

func (a *account) onRoomOccupantLeftTheRoom(from jid.Full, ident jid.Full, affiliation, role string) {
	room, ok := a.getRoom(from.Bare())
	if !ok {
		return
	}

	roster := room.Roster()
	_, left, err := roster.UpdatePresence(from, "unavailable", affiliation, role, "", "unavailable", "Occupant left the room", ident)
	if err != nil {
		a.log.WithFields(log.Fields{
			"from":        from,
			"occupant":    ident,
			"affiliation": affiliation,
			"role":        role,
		}).WithError(err).Error("An error occurred trying to remove the occupant from the roster")
		return
	}

	// TODO: we are receiving a `left the room` event, so,
	// if !left we should log it
	if left {
		view := getViewFromRoom(room)
		view.onRoomOccupantLeftTheRoomReceived(ident.Resource(), roster.AllOccupants())
	}
}
