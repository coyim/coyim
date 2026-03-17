package session

import (
	"github.com/coyim/coyim/session/events"
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

func (m *mucManager) publishMUCError(roomID jid.Bare, nickname string, e *xmppData.StanzaError) {
	ev := events.MUCError{}
	ev.ErrorType = getEventErrorTypeBasedOnStanzaError(e)
	ev.Room = roomID
	ev.Nickname = nickname

	m.publishEvent(ev)
}

func (m *mucManager) publishMUCMessageError(roomID jid.Bare, e *xmppData.StanzaError) {
	ev := events.MUCError{}
	ev.ErrorType = getEventErrorTypeBasedOnMessageError(e)
	ev.Room = roomID

	m.publishEvent(ev)
}

func isMUCErrorPresence(e *xmppData.StanzaError) bool {
	return e != nil && (e.MUCNotAuthorized != nil ||
		e.MUCForbidden != nil ||
		e.MUCItemNotFound != nil ||
		e.MUCNotAllowed != nil ||
		e.MUCNotAcceptable != nil ||
		e.MUCRegistrationRequired != nil ||
		e.MUCConflict != nil ||
		e.MUCServiceUnavailable != nil)
}

func getEventErrorTypeBasedOnMessageError(e *xmppData.StanzaError) events.MUCErrorType {
	t := events.MUCNoError
	switch {
	case e.MUCForbidden != nil:
		t = events.MUCMessageForbidden
	case e.MUCNotAcceptable != nil:
		t = events.MUCMessageNotAcceptable
	}
	return t
}

func getEventErrorTypeBasedOnStanzaError(e *xmppData.StanzaError) events.MUCErrorType {
	t := events.MUCNoError
	switch {
	case e.MUCNotAuthorized != nil:
		t = events.MUCNotAuthorized
	case e.MUCForbidden != nil:
		t = events.MUCForbidden
	case e.MUCItemNotFound != nil:
		t = events.MUCItemNotFound
	case e.MUCNotAllowed != nil:
		t = events.MUCNotAllowed
	case e.MUCNotAcceptable != nil:
		t = events.MUCNotAcceptable
	case e.MUCRegistrationRequired != nil:
		t = events.MUCRegistrationRequired
	case e.MUCConflict != nil:
		t = events.MUCConflict
	case e.MUCServiceUnavailable != nil:
		t = events.MUCServiceUnavailable
	}
	return t
}

func isMUCError(e *xmppData.StanzaError) bool {
	return getEventErrorTypeBasedOnStanzaError(e) != events.MUCNoError
}
