package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (a *account) getRoom(ident jid.Bare) (*muc.Room, error) {
	room, ok := a.roomManager.GetRoom(ident)
	if !ok {
		a.log.WithField("room", ident).Error("Trying to get a room that is not in the manager")
		return nil, newRoomNotExistsError(ident)
	}
	return room, nil
}

func (a *account) getRoomView(ident jid.Bare) (*roomView, error) {
	room, err := a.getRoom(ident)
	if err != nil {
		a.log.WithField("room", ident).WithError(err).Error("An error occurred while trying to get a room view")
		return nil, err
	}
	return getViewFromRoom(room), nil
}

func (a *account) onRoomNicknameConflict(from jid.Full) {
	view, err := a.getRoomView(from.Bare())
	if err != nil {
		a.log.WithField("from", from).WithError(err).Error("An error occurred when a room nickname conflict event was received")
		return
	}

	view.onNicknameConflictReceived(from)
}

func (a *account) onRoomRegistrationRequired(from jid.Full) {
	view, err := a.getRoomView(from.Bare())
	if err != nil {
		a.log.WithField("from", from).WithError(err).Error("An error occurred when a room registration required event was received")
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

	room, err := a.getRoom(from.Bare())
	if err != nil {
		l.WithError(err).Error("An error occurred when a room occupant joined event was received")
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

	room, err := a.getRoom(from.Bare())
	if err != nil {
		l.WithError(err).Error("An error occurred when a room occupant updated event was received")
		return
	}

	roster := room.Roster()
	_, _, err = roster.UpdatePresence(from, "", affiliation, role, "", "", "Occupant updated", occupant)
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

	room, err := a.getRoom(from.Bare())
	if err != nil {
		l.WithError(err).Error("An error occurred when an occupant left the room event was received")
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
