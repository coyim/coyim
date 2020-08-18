package gui

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (u *gtkUI) handleOneMUCErrorEvent(from jid.Full, ev events.MUCError, a *account) {
	switch ev.ErrorType {
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
		u.handleErrorMUCConflictEvent(from, ev, a)
	case events.MUCServiceUnavailable:
		a.log.Debug("MUC Error MUCServiceUnavailable received")
	default:
		a.log.WithField("event", ev).Warn("unsupported muc error event")
	}
}

func (u *gtkUI) handleErrorMUCConflictEvent(from jid.Full, ev events.MUCError, a *account) {
	// TODO[OB]-MUC: Is debug level the right level for this one?
	// TODO[OB]-MUC: When it's only one field, you should use WithField(), not WithFields()

	a.log.WithFields(log.Fields{
		"from": from,
	}).Debug("Nickname conflict event received")

	a.errorNewOccupantRoomEvent(from, ev)
}
