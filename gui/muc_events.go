package gui

import (
	"github.com/coyim/coyim/session/events"
)

func (u *gtkUI) handleOneMUCRoomEvent(ev interface{}, rv *roomView) {
	switch t := ev.(type) {
	case events.MUCOccupantJoined:
		doInUIThread(func() {
			u.handleMUCJoinedEvent(t, rv)
		})
	case events.MUCOccupantUpdated:
		doInUIThread(func() {
			u.handleMUCUpdatedEvent(t, rv)
		})
	default:
		u.log.WithField("event", t).Warn("unsupported event")
	}
}

func (u *gtkUI) handleMUCJoinedEvent(ev events.MUCOccupantJoined, rv *roomView) {
	u.log.WithField("Event", ev).Info("handleMUCJoinedEvent")
}

func (u *gtkUI) handleMUCUpdatedEvent(ev events.MUCOccupantUpdated, rv *roomView) {
	u.log.WithField("Event", ev).Info("handleMUCUpdatedEvent")
}
