package muc

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
)

type nicknameError struct {
	message  string
	nickname jid.Resource
}

// Error returns the error message
func (e *nicknameError) Error() string {
	return i18n.Localf("Can't join the room using \"%s\". Nickname is already being used", e.nickname)
}

// NewNicknameConflictError creates a new nickname conflict error
func NewNicknameConflictError(n jid.Resource) error {
	return &nicknameError{
		nickname: n,
	}
}
