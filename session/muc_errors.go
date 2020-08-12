package session

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/data"
)

func (s *session) publishMUCError(stanza *data.ClientPresence) {
	e := events.MUCErrorEvent{
		MUC: &events.MUC{
			From: stanza.From,
		},
	}

	switch {
	case stanza.Error.MUCNotAuthorized != nil:
		e.Event = events.MUCNotAuthorized
	case stanza.Error.MUCForbidden != nil:
		e.Event = events.MUCForbidden
	case stanza.Error.MUCItemNotFound != nil:
		e.Event = events.MUCItemNotFound
	case stanza.Error.MUCNotAllowed != nil:
		e.Event = events.MUCNotAllowed
	case stanza.Error.MUCNotAceptable != nil:
		e.Event = events.MUCNotAceptable
	case stanza.Error.MUCRegistrationRequired != nil:
		e.Event = events.MUCRegistrationRequired
	case stanza.Error.MUCConflict != nil:
		e.Event = events.MUCConflict
	case stanza.Error.MUCServiceUnavailable != nil:
		e.Event = events.MUCServiceUnavailable
	}
	s.publishEvent(e)
}
