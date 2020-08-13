package gui

import (
	"errors"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (u *gtkUI) getRoomView(rid jid.Bare, account *account) (*roomView, *muc.Room, error) {
	room, exists := account.roomManager.GetRoom(rid)
	if !exists {
		return nil, nil, errors.New("The rooms doesn't exists")
	}

	rv := room.Opaque.(*roomView)

	return rv, room, nil
}

func (u *gtkUI) roomOcuppantJoinedOn(a *account, ev events.MUCOccupantJoined) {
	rid := ev.From
	rv, room, err := u.getRoomView(rid, a)
	if err != nil {
		a.log.WithError(err).Debug()
	}
	// Updating the room occupant in the room manager
	fjid := ev.From.WithResource(jid.Resource(ev.Nickname))
	realjid := ev.Jid
	room.Roster().UpdatePresence(fjid, "", ev.Affiliation, ev.Role, "", ev.Status, "Room Joined", realjid)
	rv.roomOcuppantJoinedOn(err)
}

func (u *gtkUI) roomOccupantUpdatedOn(a *account, ev events.MUCOccupantUpdated) {
	//TODO: Implements the actions to do when a Occupant presence is received
	a.log.Debug("roomOccupantUpdatedOn")
}

func (u *gtkUI) roomOcuppantJoinFailedOn(a *account, ev events.MUCError) {
	ridwr, nickname := ev.EventInfo.From.PotentialSplit()
	rid := ridwr.(jid.Bare)
	rv, _, err := u.getRoomView(rid, a)
	if err != nil {
		a.log.WithError(err).Debug()
	}
	err = muc.NewNicknameConflictError(nickname).New()
	rv.roomOcuppantJoinedOn(err)
}
