package session

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/data"
)

func (s *session) handleMUCError(stanza *data.ClientPresence) {
	switch {
	case stanza.Error.MUCNotAuthorized != nil:
		s.publishEvent(events.MUCNotAuthorized)
	case stanza.Error.MUCForbidden != nil:
		s.publishEvent(events.MUCForbidden)
	case stanza.Error.MUCItemNotFound != nil:
		s.publishEvent(events.MUCItemNotFound)
	case stanza.Error.MUCNotAllowed != nil:
		s.publishEvent(events.MUCNotAllowed)
	case stanza.Error.MUCNotAceptable != nil:
		s.publishEvent(events.MUCNotAceptable)
	case stanza.Error.MUCRegistrationRequired != nil:
		s.publishEvent(events.MUCRegistrationRequired)
	case stanza.Error.MUCConflict != nil:
		s.publishEvent(events.MUCConflict)
	case stanza.Error.MUCServiceUnavailable != nil:
		s.publishEvent(events.MUCServiceUnavailable)
	}
}
