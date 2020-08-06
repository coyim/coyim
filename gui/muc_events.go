package gui

import (
	"github.com/coyim/coyim/session/events"
)

func (u *gtkUI) handleOneMUCRoomEvent(ev interface{}, r *roomView) {
	switch t := ev.(type) {
	case events.MUCPresence:
		doInUIThread(func() {
			u.handleMUCPresenceEvent(t, r)
		})
	default:
		u.log.WithField("event", t).Warn("unsupported event")
	}
}

func (u *gtkUI) observeMUCRoomEvents(r *roomView) {
	for ev := range r.events {
		u.handleOneMUCRoomEvent(ev, r)
	}
}

func (u *gtkUI) handleMUCPresenceEvent(ev events.MUCPresence, r *roomView) {
	for _, f := range r.connectionEventHandlers {
		f()
	}
}
