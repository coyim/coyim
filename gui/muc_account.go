package gui

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/i18n"
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

	view.updateOccupantsInModel(room.Roster().AllOccupants())
	view.onJoin <- joined
}

func (a *account) updateOccupantRoomEvent(from jid.Full, occupant jid.Full, affiliation, role string) {
	room, err := a.roomForIdentity(from.Bare())
	if err != nil {
		a.log.WithField("from", from).WithError(err).Error("An error occurred trying to get the room")
		return
	}

	view := getViewFromRoom(room)
	_, _, err = room.Roster().UpdatePresence(from, "", affiliation, role, "", "", "Room Joined", occupant)

	if err != nil {
		a.log.WithFields(log.Fields{
			"from":        from,
			"occupant":    occupant,
			"affiliation": affiliation,
		}).WithError(err).Error("Error on trying to join a new occupant")
		return
	}

	view.updateOccupantsInModel(room.Roster().AllOccupants())
}

func (a *account) onRoomNicknameConflict(from jid.Full) {
	view, err := a.roomViewFor(from.Bare())
	if err != nil {
		a.log.WithField("from", from).WithError(err).Error("An error occurred trying to get the room view")
		return
	}

	view.lastErrorMessage = i18n.Localf("Can't join the room using \"%s\" because the nickname is already being used.", from.Resource())
	view.onJoin <- false
}

func (a *account) onErrorRegistrationRequired(from jid.Full) {
	view, err := a.roomViewFor(from.Bare())
	if err != nil {
		a.log.WithError(err).Error("Error getting the room view")
		return
	}
	view.lastErrorMessage = i18n.Local("Sorry, this room only allows registered members")
	view.onJoin <- false
}

func (a *account) removeOccupantFromRoomRoster(from jid.Full, occupant jid.Full, affiliation, role string) {
	room, err := a.roomForIdentity(from.Bare())
	if err != nil {
		a.log.WithField("from", from).WithError(err).Error("An error occurred trying to get the room")
		return
	}

	view := getViewFromRoom(room)

	_, left, err := room.Roster().UpdatePresence(from, "unavailable", affiliation, role, "", "unavailable", "Room left", occupant)
	if err != nil {
		a.log.WithFields(log.Fields{
			"occupant":    occupant,
			"affiliation": affiliation,
			"role":        role,
		}).WithError(err).Error("An error occurred trying to remove the occupant from the roster")
		return
	}

	if left {
		view.showOccupantLeftRoom(occupant.Resource())
		view.updateOccupantsInModel(room.Roster().AllOccupants())
	}

}
