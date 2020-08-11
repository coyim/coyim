package gui

import (
	"github.com/coyim/coyim/session/events"
)

func (u *gtkUI) handleOneMUCRoomEvent(ev interface{}, rv *roomView) {
	switch t := ev.(type) {
	case events.MUCOccupantJoinedType:
		doInUIThread(func() {
			u.handleMUCJoinedEvent(t, rv)
		})
	case events.MUCOccupantUpdatedType:
		doInUIThread(func() {
			u.handleMUCUpdatedEvent(t, rv)
		})
	default:
		u.log.WithField("event", t).Warn("unsupported event")
	}
}

func (u *gtkUI) observeMUCRoomEvents(rv *roomView) {
	for ev := range rv.events {
		u.handleOneMUCRoomEvent(ev, rv)
	}
}

func (u *gtkUI) handleMUCJoinedEvent(ev events.MUCOccupantJoinedType, rv *roomView) {
	u.log.WithField("Event", ev).Info("handleMUCJoinedEvent")
}

func (u *gtkUI) handleMUCUpdatedEvent(ev events.MUCOccupantUpdatedType, rv *roomView) {
	u.log.WithField("Event", ev).Info("handleMUCUpdatedEvent")
}
