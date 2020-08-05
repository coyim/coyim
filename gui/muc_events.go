package gui

import (
	"github.com/coyim/coyim/session/events"
)

func (m *roomViewsManager) handleOneMUCRoomEvent(ev interface{}, r *roomView) {
	switch t := ev.(type) {
	case events.MUCPresence:
		doInUIThread(func() {
			m.handlePresenceEvent(t, r)
		})
	default:
		r.log.WithField("event", t).Warn("unsupported event")
	}
}

func (m *roomViewsManager) observeRoomEvents(r *roomView) {
	for ev := range r.events {
		m.handleOneMUCRoomEvent(ev, r)
	}
}

func (m *roomViewsManager) handlePresenceEvent(ev events.MUCPresence, r *roomView) {
	for _, f := range r.connectionEventHandlers {
		f()
	}
}
