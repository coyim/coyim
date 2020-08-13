package errors

import (
	"errors"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
)

// NicknameError custom error for nicknames
type NicknameError struct {
	err      string
	nickname jid.Resource
}

// Error returns the error message
func (e *NicknameError) Error() string {
	return i18n.Localf(e.err, e.nickname)
}

// GetError returns the error interface with the error information
func (e *NicknameError) GetError() error {
	return errors.New(e.Error())
}

// NewNicknameConflictError creates a new nickname conflict error
func NewNicknameConflictError(nickname jid.Resource) *NicknameError {
	return &NicknameError{
		err:      "Nickname conflict, can't join to the room using \"%s\"",
		nickname: nickname,
	}
}
