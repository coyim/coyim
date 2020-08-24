package gui

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (a *account) roomForIdentity(ident jid.Bare) (*muc.Room, error) {
	room, exists := a.roomManager.GetRoom(ident)
	if !exists {
		return nil, errors.New("the room doesn't exist")
	}

	return room, nil
}

func (a *account) roomViewFor(ident jid.Bare) (*roomView, error) {
	room, err := a.roomForIdentity(ident)
	if err != nil {
		a.log.WithError(err).Error("An error occurred trying to get the room")
		return nil, err
	}

	return getViewFromRoom(room), nil
}

func (a *account) addOccupantToRoomRoster(from jid.Full, occupant jid.Full, affiliation, role, status string) {
	room, err := a.roomForIdentity(from.Bare())
	if err != nil {
		a.log.WithField("from", from).WithError(err).Error("An error occurred trying to get the room")
		return
	}

	view := getViewFromRoom(room)

	joined, _, err := room.Roster().UpdatePresence(from, "", affiliation, role, "", status, "Room Joined", occupant)
	if err != nil {
		a.log.WithFields(log.Fields{
			"occupant":    occupant,
			"affiliation": affiliation,
			"role":        role,
			"status":      status,
		}).WithError(err).Error("An error occurred trying to add the occupant to the roster")
		view.lastErrorMessage = err.Error()
		view.onJoin <- false
		return
	}

	view.onJoin <- joined
}

func (a *account) updateOccupantRoomEvent(ev events.MUCOccupantUpdated) {
	//TODO: Implements the actions to do when a Occupant presence is received
	a.log.Debug("updateOccupantRoomEvent")
}

func (a *account) onRoomNicknameConflict(from jid.Full, message string) {
	view, err := a.roomViewFor(from.Bare())
	if err != nil {
		a.log.WithField("from", from).WithError(err).Error("An error occurred trying to get the room view")
		return
	}

	view.lastErrorMessage = message
	view.onJoin <- false
}
