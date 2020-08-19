package session

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

func (m *mucManager) publishMUCError(stanza *data.ClientPresence) {
	e := events.MUC{}
	e.From = jid.ParseFull(stanza.From)

	var t events.MUCErrorType
	switch {
	case stanza.Error.MUCNotAuthorized != nil:
		t = events.MUCNotAuthorized
	case stanza.Error.MUCForbidden != nil:
		t = events.MUCForbidden
	case stanza.Error.MUCItemNotFound != nil:
		t = events.MUCItemNotFound
	case stanza.Error.MUCNotAllowed != nil:
		t = events.MUCNotAllowed
	case stanza.Error.MUCNotAcceptable != nil:
		t = events.MUCNotAceptable
	case stanza.Error.MUCRegistrationRequired != nil:
		t = events.MUCRegistrationRequired
	case stanza.Error.MUCConflict != nil:
		t = events.MUCConflict
	case stanza.Error.MUCServiceUnavailable != nil:
		t = events.MUCServiceUnavailable
	}

	errorInfo := events.MUCError{}
	errorInfo.ErrorType = t
	e.Info = errorInfo

	m.publishEvent(e)
}
