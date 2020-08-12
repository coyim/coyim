package gui

import (
	"github.com/coyim/coyim/session/events"
	log "github.com/sirupsen/logrus"
)

func (u *gtkUI) handleOneMUCRoomEvent(ev events.MUC, a *account) {
	switch t := ev.EventInfo.(type) {
	case events.MUCOccupantJoined:
		u.handleMUCJoinedEvent(t, a)
	case events.MUCOccupantUpdated:
		u.handleMUCUpdatedEvent(t, a)
	default:
		u.log.WithField("event", t).Warn("unsupported event")
	}
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

func (u *gtkUI) handleMUCUpdatedEvent(ev events.MUCOccupantUpdated, a *account) {
	u.log.WithField("Event", ev).Debug("handleMUCUpdatedEvent")
}
