package gui

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func (u *gtkUI) handleOneMUCEvent(ev events.MUC, a *account) {
	from := ev.From

	switch t := ev.Info.(type) {
	case events.MUCOccupantUpdated:
		u.handlMUCOccupantUpdatedEvent(from, t, a)
	case events.MUCOccupantJoined:
		u.handleMUCOccupantJoinedEvent(from, t, a)
	case events.MUCOccupantLeft:
		u.handleMUCOccupantLeftEvent(from, t, a)
	case events.MUCError:
		u.handleOneMUCErrorEvent(from, t, a)
	default:
		u.log.WithFields(log.Fields{
			"type": t,
			"from": ev.From,
		}).Warn("Unsupported MUC event")
	}
}

func (u *gtkUI) handleMUCOccupantJoinedEvent(from jid.Full, ev events.MUCOccupantJoined, a *account) {
	a.log.WithFields(log.Fields{
		"from":        ev.From,
		"nickname":    ev.Nickname,
		"affiliation": ev.Affiliation,
		"role":        ev.Role,
	}).Debug("Room occupant joined event received")

	a.onRoomOccupantJoined(from, ev.Jid, ev.Affiliation, ev.Role, ev.Status)
}

func (u *gtkUI) handlMUCOccupantUpdatedEvent(from jid.Full, ev events.MUCOccupantUpdated, a *account) {
	a.log.WithFields(log.Fields{
		"from":        ev.From,
		"nickname":    ev.Nickname,
		"affiliation": ev.Affiliation,
		"role":        ev.Role,
	}).Debug("Room occupant presence updated event received")

	a.onRoomOccupantUpdated(from, ev.Jid, ev.Affiliation, ev.Role)
}

func (u *gtkUI) handleMUCOccupantLeftEvent(from jid.Full, ev events.MUCOccupantLeft, a *account) {
	a.log.WithFields(log.Fields{
		"from":        ev.From,
		"nickname":    ev.Nickname,
		"affiliation": ev.Affiliation,
		"role":        ev.Role,
	}).Debug("Occupant left the room event received")

	a.onRoomOccupantLeftTheRoom(from, ev.Jid, ev.Affiliation, ev.Role)
}
