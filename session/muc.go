package session

import (
	"sync"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

const (
	// MUCStatusPresenceJoined inform user that presence refers to one of its own room occupants
	MUCStatusPresenceJoined = "110"
)

func isMUCPresence(stanza *data.ClientPresence) bool {
	return stanza.MUC != nil
}

func isMUCUserPresence(stanza *data.ClientPresence) bool {
	return stanza.MUCUser != nil
}

func (s *session) handleMUCPresence(stanza *data.ClientPresence) {
	from := jid.Parse(stanza.From)
	ridwr, nr := from.PotentialSplit()
	rid := ridwr.(jid.Bare)
	nickname := nr.String()

	switch {
	case stanza.MUCUser != nil:
		if stanza.MUCUser.Item != nil {
			s.mucOccupantUpdate(rid, nickname, stanza.MUCUser.Item.Affiliation, stanza.MUCUser.Item.Role)
		}

		if len(stanza.MUCUser.Status) > 0 {
			affiliation := stanza.MUCUser.Item.Affiliation
			realjid := jid.Parse(stanza.MUCUser.Item.Jid).(jid.WithResource)
			role := stanza.MUCUser.Item.Role
			for _, status := range stanza.MUCUser.Status {
				switch status.Code {
				case MUCStatusPresenceJoined:
					s.mucOccupantJoined(rid, realjid, nickname, affiliation, role, status.Code, true)
				}
			}
		}
	}
}

func (s *session) mucOccupantUpdate(rid jid.Bare, nickname, affiliation, role string) {
	ev := events.MUCOccupantUpdated{}
	ev.From = rid
	ev.Nickname = nickname
	ev.Affiliation = affiliation
	ev.Role = role

	s.publishMUCEvent(ev, events.MUCOccupantUpdate)
}

func (s *session) mucOccupantJoined(rid jid.Bare, realjid jid.WithResource, nickname, affiliation, role, status string, v bool) {
	ev := events.MUCOccupantJoined{}
	ev.From = rid
	ev.Nickname = nickname
	ev.Affiliation = affiliation
	ev.Jid = realjid
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

type chatServiceReceivalContext struct {
	sync.RWMutex

	resultsChannel chan jid.Domain
	errorChannel   chan error

	s *session
}

func (c *chatServiceReceivalContext) end() {
	c.Lock()
	defer c.Unlock()
	if c.resultsChannel != nil {
		close(c.resultsChannel)
		close(c.errorChannel)
		c.resultsChannel = nil
		c.errorChannel = nil
	}
}

func (s *session) createChatServiceReceivalContext() *chatServiceReceivalContext {
	result := &chatServiceReceivalContext{}

	result.resultsChannel = make(chan jid.Domain)
	result.errorChannel = make(chan error)
	result.s = s

	return result
}

func (c *chatServiceReceivalContext) fetchChatServices(server jid.Domain) {
	defer c.end()
	items, err := c.s.conn.QueryServiceItems(server.String())
	if err != nil {
		c.RLock()
		defer c.RUnlock()
		if c.errorChannel != nil {
			c.errorChannel <- err
		}
		return
	}
	for _, item := range items.DiscoveryItems {
		if c.s.hasSomeChatService(item) {
			c.RLock()
			defer c.RUnlock()
			if c.resultsChannel == nil {
				return
			}
			c.resultsChannel <- jid.Parse(item.Jid).Host()
		}
	}
}

//GetChatServices offers the chat services from a xmpp server.
func (s *session) GetChatServices(server jid.Domain) (<-chan jid.Domain, <-chan error, func()) {
	ctx := s.createChatServiceReceivalContext()
	go ctx.fetchChatServices(server)
	return ctx.resultsChannel, ctx.errorChannel, ctx.end
}
