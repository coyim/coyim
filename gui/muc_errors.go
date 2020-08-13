package gui

import (
	"github.com/coyim/coyim/session/events"
	log "github.com/sirupsen/logrus"
)

func (u *gtkUI) handleOneMUCErrorEvent(ev events.MUCError, a *account) {
	switch ev.EventType {
	case events.MUCNotAuthorized:
		a.log.Debug("MUC Error NotAuthorized received")
	case events.MUCForbidden:
		a.log.Debug("MUC Error MUCForbidden received")
	case events.MUCItemNotFound:
		a.log.Debug("MUC Error MUCItemNotFound received")
	case events.MUCNotAllowed:
		a.log.Debug("MUC Error MUCNotAllowed received")
	case events.MUCNotAceptable:
		a.log.Debug("MUC Error MUCNotAceptable received")
	case events.MUCRegistrationRequired:
		a.log.Debug("MUC Error MUCRegistrationRequired received")
	case events.MUCConflict:
		u.handleErrorMUCConflictEvent(a, ev)
	case events.MUCServiceUnavailable:
		a.log.Debug("MUC Error MUCServiceUnavailable received")
	default:
		a.log.WithField("event", ev).Warn("unsupported event")
	}
}

func (u *gtkUI) handleErrorMUCConflictEvent(a *account, ev events.MUCError) {
	a.log.WithField("Event", ev).Debug("handleErrorMUCConflictEvent")
	a.log.WithFields(log.Fields{
		"from": ev.EventInfo.From,
	}).Info("Nickname conflict received")

	u.roomOcuppantJoinFailedOn(a, ev)
}
