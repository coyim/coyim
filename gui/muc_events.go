package gui

import (
	"github.com/coyim/coyim/session/events"
	log "github.com/sirupsen/logrus"
)

func (u *gtkUI) handleOneMUCRoomEvent(ev interface{}, a *account) {
	switch t := ev.(type) {
	case events.MUCOccupantJoined:
		doInUIThread(func() {
			u.handleMUCJoinedEvent(t, a)
		})
	case events.MUCOccupantUpdated:
		doInUIThread(func() {
			u.handleMUCUpdatedEvent(t, a)
		})
	default:
		u.log.WithField("event", t).Warn("unsupported event")
	}
}

func (u *gtkUI) handleMUCJoinedEvent(ev events.MUCOccupantJoined, a *account) {
	u.log.WithField("Event", ev).Debug("handleMUCJoinedEvent")
	a.log.WithFields(log.Fields{
		"from":        ev.From,
		"nickname":    ev.Nickname,
		"affiliation": ev.Affiliation,
		"role":        ev.Role,
	}).Info("Room Joined event received")
}

func (u *gtkUI) handleMUCUpdatedEvent(ev events.MUCOccupantUpdated, a *account) {
	u.log.WithField("Event", ev).Debug("handleMUCUpdatedEvent")
}
