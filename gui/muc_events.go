package gui

import (
	"github.com/coyim/coyim/session/events"
	log "github.com/sirupsen/logrus"
)

func (u *gtkUI) handleOneMUCEvent(ev events.MUC, a *account) {
	// TODO: we should probably just look up the room here directly, instead of duplicating it all over the place

	switch t := ev.(type) {
	case events.MUCOccupantUpdated:
		u.handleMUCOccupantUpdatedEvent(t, a)
	case events.MUCOccupantJoined:
		u.handleMUCOccupantJoinedEvent(t, a)
	case events.MUCOccupantLeft:
		u.handleMUCOccupantLeftEvent(t, a)
	case events.MUCError:
		u.handleOneMUCErrorEvent(t, a)
	case events.MUCLoggingEnabled:
		a.handleMUCLoggingEnabled(t.Room)
	case events.MUCLoggingDisabled:
		a.handleMUCLoggingDisabled(t.Room)
	default:
		u.log.WithField("event", ev).Warn("Unsupported MUC event")
	}
}

func (u *gtkUI) handleMUCOccupantJoinedEvent(ev events.MUCOccupantJoined, a *account) {
	a.log.WithFields(log.Fields{
		"room":        ev.Room,
		"nickname":    ev.Nickname,
		"affiliation": ev.Affiliation,
		"role":        ev.Role,
	}).Debug("Room occupant joined event received")

	a.onRoomOccupantJoined(ev.Room, ev.Nickname, ev.RealJid, ev.Affiliation, ev.Role, ev.Status)
}

func (u *gtkUI) handleMUCOccupantUpdatedEvent(ev events.MUCOccupantUpdated, a *account) {
	a.log.WithFields(log.Fields{
		"room":        ev.Room,
		"nickname":    ev.Nickname,
		"affiliation": ev.Affiliation,
		"role":        ev.Role,
	}).Debug("Room occupant presence updated event received")

	a.onRoomOccupantUpdated(ev.Room, ev.Nickname, ev.RealJid, ev.Affiliation, ev.Role)
}

func (u *gtkUI) handleMUCOccupantLeftEvent(ev events.MUCOccupantLeft, a *account) {
	a.log.WithFields(log.Fields{
		"room":        ev.Room,
		"nickname":    ev.Nickname,
		"affiliation": ev.Affiliation,
		"role":        ev.Role,
	}).Debug("Occupant left the room event received")

	a.onRoomOccupantLeftTheRoom(ev.Room, ev.Nickname, ev.RealJid, ev.Affiliation, ev.Role)
}
