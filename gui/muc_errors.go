package gui

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (u *gtkUI) handleOneMUCErrorEvent(ev events.MUCError, a *account) {
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
		u.handleErrorMUCRegistrationRequiredEvent(ev.Room, ev.Nickname, a)
	case events.MUCConflict:
		u.handleErrorMUCNicknameConflictEvent(ev.Room, ev.Nickname, a)
	case events.MUCServiceUnavailable:
		a.log.Debug("MUC Error MUCServiceUnavailable received")
	default:
		a.log.WithField("event", ev).Warn("unsupported muc error event")
	}
}

func (u *gtkUI) handleErrorMUCNicknameConflictEvent(room jid.Bare, nickname string, a *account) {
	a.log.WithFields(log.Fields{
		"room":     room,
		"nickname": nickname,
	}).Error("Room nickname conflict event received")

	a.onRoomNicknameConflict(room, nickname)
}

func (u *gtkUI) handleErrorMUCRegistrationRequiredEvent(room jid.Bare, nickname string, a *account) {
	a.log.WithFields(log.Fields{
		"room":     room,
		"nickname": nickname,
	}).Error("Room registration required event received")

	a.onRoomRegistrationRequired(room, nickname)
}
