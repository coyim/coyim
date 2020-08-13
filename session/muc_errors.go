package session

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

func (s *session) publishMUCError(stanza *data.ClientPresence) {
	e := events.MUCError{}
	e.EventInfo = events.MUCInfo{
		From: jid.Parse(stanza.From).(jid.Bare),
	}

	switch {
	case stanza.Error.MUCNotAuthorized != nil:
		e.EventType = events.MUCNotAuthorized
	case stanza.Error.MUCForbidden != nil:
		e.EventType = events.MUCForbidden
	case stanza.Error.MUCItemNotFound != nil:
		e.EventType = events.MUCItemNotFound
	case stanza.Error.MUCNotAllowed != nil:
		e.EventType = events.MUCNotAllowed
	case stanza.Error.MUCNotAceptable != nil:
		e.EventType = events.MUCNotAceptable
	case stanza.Error.MUCRegistrationRequired != nil:
		e.EventType = events.MUCRegistrationRequired
	case stanza.Error.MUCConflict != nil:
		e.EventType = events.MUCConflict
	case stanza.Error.MUCServiceUnavailable != nil:
		e.EventType = events.MUCServiceUnavailable
	}

	s.publishEvent(e)
}
