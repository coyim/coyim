package gui

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/jid"
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
	case events.MUCNotAcceptable:
		a.log.Debug("MUC Error MUCNotAcceptable received")
	case events.MUCRegistrationRequired:
		a.log.Debug("MUC Error MUCRegistrationRequired received")
	case events.MUCConflict:
		u.handleErrorMUCConflictEvent(from, a)
	case events.MUCServiceUnavailable:
		a.log.Debug("MUC Error MUCServiceUnavailable received")
	default:
		a.log.WithField("event", ev).Warn("unsupported muc error event")
	}
}

func (u *gtkUI) handleErrorMUCConflictEvent(from jid.Full, a *account) {
	view, err := a.roomViewFor(from.Bare())
	if err != nil {
		a.log.WithField("from", from).WithError(err).Error("An error occurred trying to get the room view")
		return
	}

	err = newNicknameConflictError(from.Resource())
	a.log.WithField("from", from).WithError(err).Error("Nickname conflict event received")
	view.lastErrorMessage = err.Error()
	view.onJoin <- false
}
