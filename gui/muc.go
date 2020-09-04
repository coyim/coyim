package gui

import (
	"fmt"

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

type nicknameError struct {
	nickname jid.Resource
}

type registrationRequiredError struct {
	room jid.Bare
}

type roomNotExistsError struct {
	ident jid.Bare
}

func (e *nicknameError) Error() string {
	return fmt.Sprintf("the nickname \"%s\" is already being used", e.nickname)
}

func (e *registrationRequiredError) Error() string {
	return fmt.Sprintf("the room \"%s\" only allows registered members", e.room)
}

func newNicknameConflictError(n jid.Resource) error {
	return &nicknameError{
		nickname: n,
	}
}

func newRegistrationRequiredError(ident jid.Bare) error {
	return &registrationRequiredError{
		room: ident,
	}
}
