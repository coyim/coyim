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

type withCallbacks struct {
	callbacks []func()
	lock      sync.RWMutex
}

func newWithCallbacks() *withCallbacks {
	return &withCallbacks{}
}

func (wc *withCallbacks) add(cb func()) {
	wc.lock.Lock()
	defer wc.lock.Unlock()

	wc.callbacks = append(wc.callbacks, cb)
}

func (wc *withCallbacks) invokeAll() {
	wc.lock.Lock()
	callbacks := wc.callbacks
	wc.lock.Unlock()

	for _, cb := range callbacks {
		cb()
	}
}
