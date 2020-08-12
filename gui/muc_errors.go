package gui

import (
	"github.com/coyim/coyim/session/events"
	log "github.com/sirupsen/logrus"
)

func (u *gtkUI) handleOneMUCErrorEvent(ev events.MUCError, a *account) {
	switch ev.EventType {
	case events.MUCNotAuthorized:
		a.log.WithField("account", a).Debug("MUC Error NotAuthorized received")
	case events.MUCForbidden:
		a.log.WithField("account", a).Debug("MUC Error MUCForbidden received")
	case events.MUCItemNotFound:
		a.log.WithField("account", a).Debug("MUC Error MUCItemNotFound received")
	case events.MUCNotAllowed:
		a.log.WithField("account", a).Debug("MUC Error MUCNotAllowed received")
	case events.MUCNotAceptable:
		a.log.WithField("account", a).Debug("MUC Error MUCNotAceptable received")
	case events.MUCRegistrationRequired:
		a.log.WithField("account", a).Debug("MUC Error MUCRegistrationRequired received")
	case events.MUCConflict:
		u.handleErrorMUCConflictEvent(a, ev)
	case events.MUCServiceUnavailable:
		a.log.WithField("account", a).Debug("MUC Error MUCServiceUnavailable received")
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
