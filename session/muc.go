package session

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

const (
	// MUCStatusPresenceJoined inform user that presence refers to one of its own room occupants
	MUCStatusPresenceJoined = "110"
)

func (s *session) isMUCPresence(stanza *data.ClientPresence) bool {
	return stanza.MUC != nil
}

func (s *session) receivedMUCPresence(stanza *data.ClientPresence) bool {
	from := jid.Parse(stanza.From)
	rid, nickname := from.PotentialSplit()

	if stanza.MUCUser != nil {
		if stanza.MUCUser.Item != nil {
			s.mucOccupantUpdate(rid.String(), string(nickname), stanza.MUCUser.Item.Affiliation, stanza.MUCUser.Item.Role)
		}

		if len(stanza.MUCUser.Status) > 0 {
			for _, status := range stanza.MUCUser.Status {
				switch status.Code {
				case MUCStatusPresenceJoined:
					s.mucOccupantJoined(rid.String(), string(nickname), true)
				}
			}
		}

		return true
	}

	return false
}

func (s *session) mucOccupantUpdate(rid, nickname, affiliation, role string) {
	s.publishEvent(events.MUCOccupantUpdatedType{
		MUCOccupantType: &events.MUCOccupantType{
			MUCType: &events.MUCType{
				From: rid,
			},
			Nickname: nickname,
		},
		Affiliation: affiliation,
		Role:        role,
	})

	s.mucRosterUpdated()
}

func (s *session) mucRosterUpdated() {
	s.publishEvent(events.MUCOccupantUpdate)
}

func (s *session) mucOccupantJoined(rid, nickname string, v bool) {
	s.publishEvent(events.MUCOccupantJoinedType{
		MUCOccupantType: &events.MUCOccupantType{
			MUCType: &events.MUCType{
				From: rid,
			},
			Nickname: nickname,
		},
		Joined: v,
	})
}
