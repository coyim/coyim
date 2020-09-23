package gui

import (
	"github.com/coyim/coyim/session/events"
)

func (a *account) handleMUCErrorEvent(ev events.MUCError) {
	view, ok := a.getRoomView(ev.WhichRoom())
	if !ok {
		a.log.WithField("room", ev.WhichRoom()).Error("Not possible to get room view when handling multi user chat event")
		return
	}

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
		view.registrationRequired(ev.WhichNickname())
	case events.MUCConflict:
		view.nicknameConflict(ev.WhichNickname())
	case events.MUCServiceUnavailable:
		a.log.Debug("MUC Error MUCServiceUnavailable received")
	default:
		a.log.WithField("event", ev).Warn("unsupported muc error event")
	}
}
