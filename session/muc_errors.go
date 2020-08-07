package session

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/data"
)

func (s *session) receivedMUCError(stanza *data.ClientPresence) bool {
	if s.isMUCPresence(stanza) {
		functions := []func(*data.ClientPresence) bool{
			s.mucErrorNotAuthorized,
			s.mucErrorForbidden,
			s.mucErrorItemNotFound,
			s.mucErrorNotAllowed,
			s.mucErrorNotAceptable,
			s.mucErrorRegistrationRequired,
			s.mucErrorConflict,
			s.mucErrorServiceUnavailable,
		}

		for _, fn := range functions {
			if ok := fn(stanza); ok {
				return true
			}
		}
	}

	return false
}

func (s *session) isMUCErrorNotAuthorized(stanza *data.ClientPresence) bool {
	return stanza.Error.MUCNotAuthorized != nil
}

func (s *session) mucErrorNotAuthorized(stanza *data.ClientPresence) bool {
	if !s.isMUCErrorNotAuthorized(stanza) {
		return false
	}

	s.publishEvent(events.MUCNotAuthorized)

	return true
}

func (s *session) isMUCErrorForbidden(stanza *data.ClientPresence) bool {
	return stanza.Error.MUCForbidden != nil
}

func (s *session) mucErrorForbidden(stanza *data.ClientPresence) bool {
	if !s.isMUCErrorForbidden(stanza) {
		return false
	}

	s.publishEvent(events.MUCForbidden)

	return true
}

func (s *session) isMUCErrorItemNotFound(stanza *data.ClientPresence) bool {
	return stanza.Error.MUCItemNotFound != nil
}

func (s *session) mucErrorItemNotFound(stanza *data.ClientPresence) bool {
	if !s.isMUCErrorItemNotFound(stanza) {
		return false
	}

	s.publishEvent(events.MUCItemNotFound)

	return true
}

func (s *session) isMUCErrorNotAllowed(stanza *data.ClientPresence) bool {
	return stanza.Error.MUCNotAllowed != nil
}

func (s *session) mucErrorNotAllowed(stanza *data.ClientPresence) bool {
	if !s.isMUCErrorNotAllowed(stanza) {
		return false
	}

	s.publishEvent(events.MUCNotAllowed)

	return true
}

func (s *session) isMUCErrorNotAceptable(stanza *data.ClientPresence) bool {
	return stanza.Error.MUCNotAceptable != nil
}

func (s *session) mucErrorNotAceptable(stanza *data.ClientPresence) bool {
	if !s.isMUCErrorNotAceptable(stanza) {
		return false
	}

	s.publishEvent(events.MUCNotAceptable)

	return true
}

func (s *session) isMUCErrorRegistrationRequired(stanza *data.ClientPresence) bool {
	return stanza.Error.MUCRegistrationRequired != nil
}

func (s *session) mucErrorRegistrationRequired(stanza *data.ClientPresence) bool {
	if !s.isMUCErrorRegistrationRequired(stanza) {
		return false
	}

	s.publishEvent(events.MUCRegistrationRequired)

	return true
}

func (s *session) isMUCErrorConflict(stanza *data.ClientPresence) bool {
	return stanza.Error.MUCConflict != nil
}

func (s *session) mucErrorConflict(stanza *data.ClientPresence) bool {
	if !s.isMUCErrorConflict(stanza) {
		return false
	}

	s.publishEvent(events.MUCConflict)

	return true
}

func (s *session) isMUCErrorServiceUnavailable(stanza *data.ClientPresence) bool {
	return stanza.Error.MUCServiceUnavailable != nil
}

func (s *session) mucErrorServiceUnavailable(stanza *data.ClientPresence) bool {
	if !s.isMUCErrorServiceUnavailable(stanza) {
		return false
	}

	s.publishEvent(events.MUCServiceUnavailable)

	return true
}
