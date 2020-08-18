package gui

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (u *gtkUI) handleOneMUCEvent(ev events.MUC, a *account) {
	from := ev.From

	switch t := ev.Info.(type) {
	case events.MUCOccupantUpdated:
		u.handleMUCUpdatedEvent(from, t, a)
	case events.MUCOccupantJoined:
		u.handleMUCJoinedEvent(from, t, a)
	case events.MUCError:
		u.handleOneMUCErrorEvent(ev, a)
	default:
		u.log.WithFields(log.Fields{
			"Type":      t,
			"From":      ev.From,
			"EventType": ev.EventType,
		}).Warn("Unsupported received MUC event")
	}
}

func (u *gtkUI) handleMUCUpdatedEvent(from jid.Bare, ev events.MUCOccupantUpdated, a *account) {
	a.log.WithField("Event", ev).Debug("handleMUCUpdatedEvent")
}

func (u *gtkUI) handleMUCJoinedEvent(from jid.Bare, ev events.MUCOccupantJoined, a *account) {
	a.log.WithFields(log.Fields{
		"from":        ev.From,
		"nickname":    ev.Nickname,
		"affiliation": ev.Affiliation,
		"role":        ev.Role,
	}).Debug("Room Joined event received")

	a.enrollNewOccupantRoomEvent(from, ev)
}
