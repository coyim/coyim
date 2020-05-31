package gui

import "github.com/coyim/coyim/coylog"

type withLog interface {
	Log() coylog.Logger
}

func (u *gtkUI) Log() coylog.Logger {
	return u.log
}

func (m *accountManager) Log() coylog.Logger {
	return m.log
}

func (a *account) Log() coylog.Logger {
	return a.log
}

func (c *conversationPane) Log() coylog.Logger {
	return c.account.log
}
