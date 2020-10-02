package gui

import (
	"github.com/coyim/coyim/session/events"
)

// For now, this looks pretty useless, but later on we will
// have other events coming, such as for example for invites and other
// MUC functionality, so we retain this method for those purposes.

func (u *gtkUI) handleMUCEvent(ev events.MUC, a *account) {
	switch t := ev.(type) {
	case events.MUCError:
		a.handleMUCErrorEvent(t)
	}
}
