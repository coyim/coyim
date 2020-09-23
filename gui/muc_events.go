package gui

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/muc"
)

func (u *gtkUI) handleMUCEvent(ev muc.MUC, a *account) {
	switch t := ev.(type) {
	case events.MUCError:
		a.handleMUCErrorEvent(t)
	}
}
