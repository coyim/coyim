package muc

import (
	"fmt"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
)

type nicknameError struct {
	message  string
	nickname jid.Resource
}

// Error returns the error message
func (e *nicknameError) Error() string {
	return fmt.Sprintf(i18n.Local("Nickname conflict, can't join to the room using \"%s\""), e.nickname)
}

// NewNicknameConflictError creates a new nickname conflict error
func NewNicknameConflictError(n jid.Resource) error {
	e := &nicknameError{
		nickname: n,
	}
	return e
}
