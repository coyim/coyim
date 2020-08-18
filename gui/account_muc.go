package gui

import (
	"errors"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

// TODO[OB]-MUC: I think this filename should be muc_account instead of account_muc

func (a *account) roomViewFor(rid jid.Bare) (*roomView, *muc.Room, error) {
	room, exists := a.roomManager.GetRoom(rid)
	if !exists {
		return nil, nil, errors.New("The room doesn't exist")
	}

	rv := room.Opaque.(*roomView)

	return rv, room, nil
}

// TODO[OB]-MUC: I'm not a huge fan of this method name

func (a *account) enrollNewOccupantRoomEvent(from jid.Bare, ev events.MUCOccupantJoined) {
	rv, room, err := a.roomViewFor(from)
	if err != nil {
		a.log.WithError(err).Error("An error occurred while trying to change the room occupant status.")
		return
	}

	fjid := from.WithResource(jid.Resource(ev.Nickname))
	realjid := ev.Jid

	// TODO[OB]-MUC: You should not ignore the results of this call
	room.Roster().UpdatePresence(fjid, "", ev.Affiliation, ev.Role, "", ev.Status, "Room Joined", realjid)

	rv.processOccupantJoinedEvent(err)
}

func (a *account) updateOccupantRoomEvent(ev events.MUCOccupantUpdated) {
	//TODO: Implements the actions to do when a Occupant presence is received
	a.log.Debug("updateOccupantRoomEvent")
}

// TODO[OB]-MUC: I'm not a fan of this method name either

func (a *account) errorNewOccupantRoomEvent(ev events.MUC) {
	// TODO[OB]-MUC: It doesn't make sense to use PotentialSplit here, since you are assuming the nickname exists, you
	// should make sure the From is a full JID and use Split() instead
	ridwr, nickname := ev.From.PotentialSplit()
	rid := ridwr.(jid.Bare)
	rv, _, err := a.roomViewFor(rid)
	if err != nil {
		// TODO[OB]-MUC: It might be better to say something about WHEN the error occurred... This is not Java, we don't have backtraces on exceptions
		a.log.WithError(err).Error("An error occurred ")
		return
	}

	err = muc.NewNicknameConflictError(nickname)

	// TODO[OB]-MUC: I don't understand this at all. Why are you sending in the err here instead of processing it here?
	// The processOccupantJoinedEvent doesn't seem like it should take an error as argument. So it seems like the method is badly named?
	rv.processOccupantJoinedEvent(err)
}
