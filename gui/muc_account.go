package gui

import (
	"errors"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (a *account) roomViewFor(rid jid.Bare) (*roomView, *muc.Room, error) {
	room, exists := a.roomManager.GetRoom(rid)
	if !exists {
		return nil, nil, errors.New("The room doesn't exist")
	}

	rv := room.Opaque.(*roomView)

	return rv, room, nil
}

func (a *account) addOccupantToRoster(from jid.Full, userj jid.Full, affiliation, role, status string) {
	roomj := jid.NewBare(from.Local(), from.Host())
	rv, room, err := a.roomViewFor(roomj)
	if err != nil {
		rv.lastErrorMessage = err.Error()
		a.log.WithError(err).Error("An error occurred trying to get the room view")
		rv.onJoin <- false
		return
	}

	joined, _, err := room.Roster().UpdatePresence(from, "", affiliation, role, "", status, "Room Joined", userj)
	if err != nil {
		rv.lastErrorMessage = err.Error()
		a.log.WithError(err).Error("An error occurred trying to add the occupant to the roster")
	}
	rv.onJoin <- joined
}

func (a *account) updateOccupantRoomEvent(ev events.MUCOccupantUpdated) {
	//TODO: Implements the actions to do when a Occupant presence is received
	a.log.Debug("updateOccupantRoomEvent")
}

func (a *account) generateNicknameConflictError(from jid.Full) {
	ridwr, nickname := from.Split()
	rid := ridwr.(jid.Bare)
	rv, _, err := a.roomViewFor(rid)
	if err != nil {
		rv.lastErrorMessage = err.Error()
		a.log.WithError(err).Error("An error occurred trying to get the room view")
		rv.onJoin <- false
		return
	}
	// Generating a custom error for the nickname conflict event received
	err = muc.NewNicknameConflictError(nickname)
	a.log.WithError(err).Error("Nickname conflict event received")
	rv.lastErrorMessage = err.Error()
	rv.onJoin <- false
}
