package session

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

func (m *mucManager) publishMUCError(from jid.Full, e *data.StanzaError) {
	ev := events.MUCError{}
	ev.ErrorType = getEventErrorTypeBasedOnStanzaError(e)
	ev.Room = from.Bare()
	ev.Nickname = from.Resource().String()

	m.publishEvent(ev)
}

func getEventErrorTypeBasedOnStanzaError(e *data.StanzaError) events.MUCErrorType {
	var t events.MUCErrorType
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
