package gui

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/muc"
)

func (u *gtkUI) handleMUCEvent(ev muc.MUC, a *account) {
	view, ok := a.getRoomView(ev.WhichRoom())
	if !ok {
		a.log.WithField("room", ev.WhichRoom()).Error("Not possible to get room view when handling multi user chat event")
		return
	}

	switch t := ev.(type) {
	case events.MUCError:
		a.handleMUCErrorEvent(t, view)
	}
}
