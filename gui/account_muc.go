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

func (a *account) enrollNewOccupantRoomEvent(ev events.MUCOccupantJoined) {
	rid := ev.From
	rv, room, err := a.getRoomView(rid)
	if err != nil {
		a.log.WithError(err).Debug()
	}
	// Updating the room occupant in the room manager
	fjid := ev.From.WithResource(jid.Resource(ev.Nickname))
	realjid := ev.Jid
	room.Roster().UpdatePresence(fjid, "", ev.Affiliation, ev.Role, "", ev.Status, "Room Joined", realjid)
	rv.processOccupantJoinedEvent(err)
}

func (a *account) updateOccupantRoomEvent(ev events.MUCOccupantUpdated) {
	//TODO: Implements the actions to do when a Occupant presence is received
	a.log.Debug("updateOccupantRoomEvent")
}

func (a *account) errorNewOccupantRoomEvent(ev events.MUCError) {
	ridwr, nickname := ev.EventInfo.From.PotentialSplit()
	rid := ridwr.(jid.Bare)
	rv, _, err := a.getRoomView(rid)
	if err != nil {
		a.log.WithError(err).Debug()
	}
	err = muc.NewNicknameConflictError(nickname).New()
	rv.processOccupantJoinedEvent(err)
}
