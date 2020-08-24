package gui

import (
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

// Error returns the error message
func (e *nicknameError) Error() string {
	return i18n.Localf("Can't join the room using \"%s\" because the nickname is already being used.", e.nickname)
}

func newNicknameConflictError(n jid.Resource) error {
	return &nicknameError{
		nickname: n,
	}
}
