package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
)

func initMUC() {
	initMUCSupportedErrors()
}

var supportedCreateMUCErrors map[error]string

// We should put here all MUC-related errors that we want to support
// and have a custom a useful user message for each one
func initMUCSupportedErrors() {
	supportedCreateMUCErrors = map[error]string{
		session.ErrInvalidInformationQueryRequest: i18n.Local("Couldn't send the information query to the server, please try again."),
		session.ErrUnexpectedResponse:             i18n.Local("The connection to the server can't be established."),
		session.ErrInformationQueryResponse:       i18n.Local("You don't have the permissions to create a room on the server or the room is already created."),
	}
}
