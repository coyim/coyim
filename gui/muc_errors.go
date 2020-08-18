package gui

import (
	"github.com/coyim/coyim/session/events"
	log "github.com/sirupsen/logrus"
)

func (u *gtkUI) handleOneMUCErrorEvent(ev events.MUC, a *account) {
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
		u.handleErrorMUCConflictEvent(ev, a)
	case events.MUCServiceUnavailable:
		a.log.Debug("MUC Error MUCServiceUnavailable received")
	default:
		a.log.WithField("event", ev).Warn("unsupported muc error event")
	}
}

func (u *gtkUI) handleErrorMUCConflictEvent(ev events.MUC, a *account) {
	// TODO[OB]-MUC: Is debug level the right level for this one?
	// TODO[OB]-MUC: When it's only one field, you should use WithField(), not WithFields()

	a.log.WithFields(log.Fields{
		"from": ev.From,
	}).Debug("Nickname conflict event received")

	a.errorNewOccupantRoomEvent(ev)
}
