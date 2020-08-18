package gui

import (
	"errors"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (a *account) getRoomView(rid jid.Bare) (*roomView, *muc.Room, error) {
	room, exists := a.roomManager.GetRoom(rid)
	if !exists {
		return nil, nil, errors.New("The rooms doesn't exists")
	}

	rv := room.Opaque.(*roomView)

	return rv, room, nil
}

func (a *account) enrollNewOccupantRoomEvent(from jid.Bare, ev events.MUCOccupantJoined) {
	rv, room, err := a.getRoomView(from)
	if err != nil {
		a.log.WithError(err).Error("An error ocurred while trying to change the room occupant status.")
		return
	}

	fjid := from.WithResource(jid.Resource(ev.Nickname))
	realjid := ev.Jid

	room.Roster().UpdatePresence(fjid, "", ev.Affiliation, ev.Role, "", ev.Status, "Room Joined", realjid)

	rv.processOccupantJoinedEvent(err)
}

func (a *account) updateOccupantRoomEvent(ev events.MUCOccupantUpdated) {
	//TODO: Implements the actions to do when a Occupant presence is received
	a.log.Debug("updateOccupantRoomEvent")
}

func (a *account) errorNewOccupantRoomEvent(ev events.MUC) {
	ridwr, nickname := ev.From.PotentialSplit()
	rid := ridwr.(jid.Bare)
	rv, _, err := a.getRoomView(rid)
	if err != nil {
		a.log.WithError(err).Error("An error occurred ")
		return
	}

	err = muc.NewNicknameConflictError(nickname)
	rv.processOccupantJoinedEvent(err)
}
