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
	// TODO[OB]-MUC: It's an antipattern to send a custom message to Sprintf
	return fmt.Sprintf(e.message, e.nickname)
}

// NewNicknameConflictError creates a new nickname conflict error
func NewNicknameConflictError(n jid.Resource) error {
	e := &nicknameError{}
	e.message = i18n.Local("Nickname conflict, can't join to the room using \"%s\"")
	e.nickname = n
	return e
}
