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

func (s *session) handleMUCPresence(stanza *data.ClientPresence) {
	from := jid.Parse(stanza.From)
	rid, nickname := from.PotentialSplit()

	switch {
	case stanza.MUCUser != nil:
		if stanza.MUCUser.Item != nil {
			s.mucOccupantUpdate(rid.String(), string(nickname), stanza.MUCUser.Item.Affiliation, stanza.MUCUser.Item.Role)
		}

		if len(stanza.MUCUser.Status) > 0 {
			affiliation := stanza.MUCUser.Item.Affiliation
			jid := stanza.MUCUser.Item.Jid
			role := stanza.MUCUser.Item.Role
			for _, status := range stanza.MUCUser.Status {
				switch status.Code {
				case MUCStatusPresenceJoined:
					s.mucOccupantJoined(rid.String(), string(nickname), affiliation, jid, role, status.Code, true)
				}
			}
		}
	}
}

func (s *session) mucOccupantUpdate(rid, nickname, affiliation, role string) {
	ev := events.MUCOccupantUpdated{}
	ev.From = rid
	ev.Nickname = nickname
	ev.Affiliation = affiliation
	ev.Role = role

	s.publishMUCEvent(ev, events.MUCOccupantUpdate)
}

func (s *session) mucOccupantJoined(rid, nickname, affiliation, jid, role, status string, v bool) {
	ev := events.MUCOccupantJoined{}
	ev.From = rid
	ev.Nickname = nickname
	ev.Affiliation = affiliation
	ev.Jid = jid
	ev.Role = role
	ev.Status = status
	ev.Joined = v

	s.publishMUCEvent(ev, events.MUCOccupantJoin)
}

func (s *session) publishMUCEvent(e interface{}, t events.MUCEventType) {
	ev := events.MUC{}
	ev.EventInfo = e
	ev.EventType = t

	s.publishEvent(ev)
}

func (s *session) hasSomeConferenceService(identities []data.DiscoveryIdentity) bool {
	for _, identity := range identities {
		if identity.Category == "conference" && identity.Type == "text" {
			return true
		}
	}
	return false
}

func (s *session) hasSomeChatService(di data.DiscoveryItem) bool {
	iq, err := s.conn.QueryServiceInformation(di.Jid)
	if err != nil {
		s.log.WithError(err).Error("Error getting the information query for the service:", di.Jid)
		return false
	}
	return s.hasSomeConferenceService(iq.Identities)
}

func (s *session) filterOnlyChatServices(items *data.DiscoveryItemsQuery) []jid.Domain {
	var chatServices []jid.Domain
	for _, item := range items.DiscoveryItems {
		if s.hasSomeChatService(item) {
			chatServices = append(chatServices, jid.Parse(item.Jid).Host())
		}
	}
	return chatServices
}

//GetChatServices offers the chat services from a xmpp server.
func (s *session) GetChatServices(server jid.Domain) ([]jid.Domain, error) {
	items, err := s.conn.QueryServiceItems(server.String())
	if err != nil {
		return nil, err
	}
	return s.filterOnlyChatServices(items), nil
}
