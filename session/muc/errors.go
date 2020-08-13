package muc

import (
	"errors"
	"fmt"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
)

// CustomError base structure to define custom errors
type CustomError interface {
	Error() string
	New() error
}

type customError struct {
	message string
}

func newCustomError() CustomError {
	return &customError{}
}

func (e *customError) Error() string {
	return e.message
}

func (e *customError) New() error {
	return errors.New(e.Error())
}

// NicknameError description
type NicknameError interface {
	CustomError
}

type nicknameError struct {
	customError
	nickname jid.Resource
}

func (ne *nicknameError) Error() string {
	return fmt.Sprintf(ne.message, ne.nickname)
}

func (ne *nicknameError) New() error {
	return errors.New(ne.Error())
}

// NewNicknameConflictError description
func NewNicknameConflictError(n jid.Resource) NicknameError {
	e := &nicknameError{}
	e.nickname = n
	e.message = i18n.Local("Nickname conflict, can't join to the room using \"%s\"")
	return e
}
