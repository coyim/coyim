package session

import (
	"fmt"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (m *mucManager) roomCreated(from jid.Full, room jid.Bare) {
	ev := events.MUCRoomCreated{}
	ev.Room = room

	m.publishEvent(ev)
}

func (m *mucManager) roomRenamed(from jid.Full, room jid.Bare) {
	ev := events.MUCRoomRenamed{}
	ev.Room = room

	m.publishEvent(ev)
}

func parseAffiliationAndReport(a string) muc.Affiliation {
	aa, e := muc.AffiliationFromString(a)
	if e != nil {
		fmt.Printf("error when parsing affiliation: %v\n", e)
	}
	return aa
}

func parseRoleAndReport(a string) muc.Role {
	aa, e := muc.RoleFromString(a)
	if e != nil {
		fmt.Printf("error when parsing role: %v\n", e)
	}
	return aa
}

func (m *mucManager) occupantLeft(from jid.Full, room jid.Bare, occupant jid.Resource, affiliation, role string) {
	ev := events.MUCOccupantLeft{}
	ev.Room = room
	ev.Nickname = occupant.String()
	ev.RealJid = from
	ev.Affiliation = parseAffiliationAndReport(affiliation)
	ev.Role = parseRoleAndReport(role)

	m.publishEvent(ev)
}

func (m *mucManager) occupantUpdate(from jid.Full, room jid.Bare, occupant jid.Resource, affiliation, role string) {
	ev := events.MUCOccupantUpdated{}
	ev.Room = room
	ev.Nickname = occupant.String()
	ev.RealJid = from
	ev.Affiliation = parseAffiliationAndReport(affiliation)
	ev.Role = parseRoleAndReport(role)

	m.publishEvent(ev)
}

func (m *mucManager) publishLoggingEnabled(room jid.Bare) {
	ev := events.MUCLoggingEnabled{}
	ev.Room = room

	m.publishEvent(ev)
}

func (m *mucManager) publishLoggingDisabled(room jid.Bare) {
	ev := events.MUCLoggingDisabled{}
	ev.Room = room

	m.publishEvent(ev)
}

// selfOccupantUpdated can happen several times - every time a status code update is changed, or role or affiliation
// is updated, this can lead to the method being called. For now, it will generate a event about joining, but this
// should be cleaned up and fixed
func (m *mucManager) selfOccupantUpdated(from jid.Full, room jid.Bare, occupant jid.Resource, ident jid.Full, affiliation, role string, status mucUserStatuses) {
	ev := events.MUCOccupantJoined{}
	ev.Room = room
	ev.Nickname = occupant.String()
	ev.RealJid = ident
	ev.Affiliation = parseAffiliationAndReport(affiliation)
	ev.Role = parseRoleAndReport(role)

	m.publishEvent(ev)

	if status.contains(MUCStatusRoomLoggingEnabled) {
		m.publishLoggingEnabled(room)
	}

	if status.contains(MUCStatusRoomLoggingDisabled) {
		m.publishLoggingDisabled(room)
	}
}

func (m *mucManager) messageReceived(room jid.Bare, nickname, message, subject string) {
	ev := events.MUCMessageReceived{}
	ev.Room = room
	ev.Nickname = nickname
	ev.BodyMessage = message
	ev.Subject = subject

	m.publishEvent(ev)
}
