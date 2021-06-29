package gui

import (
	"github.com/coyim/coyim/session/events"
)

func (a *account) handleMUCErrorEvent(ev events.MUCError) {
	if view, ok := a.getRoomView(ev.Room); ok {
		switch ev.ErrorType {
		case events.MUCNotAuthorized:
			view.notAuthorized()
		case events.MUCForbidden:
			view.occupantForbidden()
		case events.MUCMessageForbidden:
			view.messageForbidden()
		case events.MUCItemNotFound:
			a.log.Debug("MUC Error MUCItemNotFound received")
		case events.MUCNotAllowed:
			a.log.Debug("MUC Error MUCNotAllowed received")
		case events.MUCNotAcceptable:
			a.log.Debug("MUC Error MUCNotAcceptable received")
		case events.MUCMessageNotAcceptable:
			view.messageNotAccepted()
		case events.MUCRegistrationRequired:
			view.registrationRequired(ev.Nickname)
		case events.MUCConflict:
			view.nicknameConflict(ev.Nickname)
		case events.MUCServiceUnavailable:
			view.serviceUnavailable()
		default:
			a.log.WithField("event", ev).Warn("Unsupported MUC error event")
			view.unknownError()
		}
	}
}
