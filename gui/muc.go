package gui

import (
	"fmt"
	"sync"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/coyim/xmpp/jid"
)

func initMUC() {
	initMUCSupportedErrors()
	initMUCConfigUpdateMessages()
}

var supportedCreateMUCErrors map[error]string

// We should put here all MUC-related errors that we want to support
// and have a custom and useful user message for each one
func initMUCSupportedErrors() {
	supportedCreateMUCErrors = map[error]string{
		session.ErrInvalidInformationQueryRequest: i18n.Local("Couldn't send the information query to the server, please try again."),
		session.ErrUnexpectedResponse:             i18n.Local("The connection to the server can't be established."),
		session.ErrInformationQueryResponse:       i18n.Local("You don't have the permissions to create a room."),
	}
}

func newNicknameConflictError(n jid.Resource) error {
	return fmt.Errorf("the nickname \"%s\" is already being used", n)
}

func newRegistrationRequiredError(roomID jid.Bare) error {
	return fmt.Errorf("the room \"%s\" only allows registered members", roomID)
}

type callbacksSet struct {
	callbacks []func()
	lock      sync.RWMutex
}

func newCallbacksSet() *callbacksSet {
	return &callbacksSet{}
}

func (s *callbacksSet) add(cb func()) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.callbacks = append(s.callbacks, cb)
}

func (s *callbacksSet) invokeAll() {
	s.lock.Lock()
	callbacks := s.callbacks
	s.lock.Unlock()

	for _, cb := range callbacks {
		cb()
	}
}
