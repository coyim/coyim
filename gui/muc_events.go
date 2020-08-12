package gui

import (
	"github.com/coyim/coyim/session/events"
	log "github.com/sirupsen/logrus"
)

func (u *gtkUI) handleOneMUCRoomEvent(ev events.MUC, a *account) {
	switch ev.EventType {
	case events.MUCOccupantUpdate:
		t := ev.EventInfo.(events.MUCOccupantUpdated)
		u.handleMUCUpdatedEvent(t, a)
	case events.MUCOccupantJoin:
		t := ev.EventInfo.(events.MUCOccupantJoined)
		u.handleMUCJoinedEvent(t, a)
	default:
		u.log.WithField("event", ev).Warn("unsupported event")
	}
}

func (u *gtkUI) handleMUCUpdatedEvent(ev events.MUCOccupantUpdated, a *account) {
	a.log.WithField("Event", ev).Debug("handleMUCUpdatedEvent")
}

func (u *gtkUI) handleMUCJoinedEvent(ev events.MUCOccupantJoined, a *account) {
	a.log.WithFields(log.Fields{
		"from":        ev.From,
		"nickname":    ev.Nickname,
		"affiliation": ev.Affiliation,
		"role":        ev.Role,
	}).Debug("Room Joined event received")

	u.roomOcuppantJoinedOn(a, ev)
}
