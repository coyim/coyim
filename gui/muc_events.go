package gui

import (
	"github.com/coyim/coyim/session/events"
	log "github.com/sirupsen/logrus"
)

func (u *gtkUI) handleOneMUCRoomEvent(ev interface{}, a *account) {
	switch t := ev.(type) {
	case events.MUCOccupantJoined:
		doInUIThread(func() {
			u.handleMUCJoinedEvent(t, a)
		})
	case events.MUCOccupantUpdated:
		doInUIThread(func() {
			u.handleMUCUpdatedEvent(t, a)
		})
	// Handling Errors
	case events.MUCErrorEvent:
		doInUIThread(func() {
			u.handleMUCErrorEvent(t, a)
		})
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
	}).Info("Room Joined event received")

	u.roomOcuppantJoinedOn(a, ev)
}

func (u *gtkUI) handleMUCUpdatedEvent(ev events.MUCOccupantUpdated, a *account) {
	u.log.WithField("Event", ev).Debug("handleMUCUpdatedEvent")
}

func (u *gtkUI) handleMUCErrorEvent(ev events.MUCErrorEvent, a *account) {
	switch ev.Event {
	case events.MUCNotAuthorized:
		u.log.WithField("account", a).Debug("MUC Error NotAuthorized received")
	case events.MUCForbidden:
		u.log.WithField("account", a).Debug("MUC Error MUCForbidden received")
	case events.MUCItemNotFound:
		u.log.WithField("account", a).Debug("MUC Error MUCItemNotFound received")
	case events.MUCNotAllowed:
		u.log.WithField("account", a).Debug("MUC Error MUCNotAllowed received")
	case events.MUCNotAceptable:
		u.log.WithField("account", a).Debug("MUC Error MUCNotAceptable received")
	case events.MUCRegistrationRequired:
		u.log.WithField("account", a).Debug("MUC Error MUCRegistrationRequired received")
	case events.MUCConflict:
		u.handleErrorMUCConflictEvent(a, ev)
	case events.MUCServiceUnavailable:
		u.log.WithField("account", a).Debug("MUC Error MUCServiceUnavailable received")
	}
}
