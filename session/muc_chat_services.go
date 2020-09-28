package session

import (
	"sync"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// GetChatServices offers the chat services from a xmpp server.
func (s *session) GetChatServices(server jid.Domain) (<-chan jid.Domain, <-chan error, func()) {
	r := s.createChatServicesReceiver()
	go r.fetchChatServices(server)
	return r.resultsChannel, r.errorChannel, r.finish
}

func (s *session) hasSomeChatService(di data.DiscoveryItem) bool {
	iq, err := s.conn.QueryServiceInformation(di.Jid)
	if err != nil {
		s.log.WithField("jid", di.Jid).WithError(err).Error("Error getting the information query for the service")
		return false
	}
	return hasSomeConferenceService(iq.Identities)
}

type chatServicesReceiver struct {
	sync.RWMutex

	resultsChannel chan jid.Domain
	errorChannel   chan error

	s *session
}

func (s *session) createChatServicesReceiver() *chatServicesReceiver {
	result := &chatServicesReceiver{}

	result.resultsChannel = make(chan jid.Domain)
	result.errorChannel = make(chan error)
	result.s = s

	return result
}

func (r *chatServicesReceiver) fetchChatServices(server jid.Domain) {
	defer r.finish()
	items, err := r.s.conn.QueryServiceItems(server.String())
	if err != nil {
		r.RLock()
		defer r.RUnlock()
		if r.errorChannel != nil {
			r.errorChannel <- err
		}
		return
	}
	for _, item := range items.DiscoveryItems {
		if r.s.hasSomeChatService(item) {
			r.RLock()
			defer r.RUnlock()
			if r.resultsChannel == nil {
				return
			}
			r.resultsChannel <- jid.Parse(item.Jid).Host()
		}
	}
}

func (r *chatServicesReceiver) finish() {
	r.Lock()
	defer r.Unlock()
	if r.resultsChannel != nil {
		close(r.resultsChannel)
		close(r.errorChannel)
		r.resultsChannel = nil
		r.errorChannel = nil
	}
}

func hasSomeConferenceService(identities []data.DiscoveryIdentity) bool {
	for _, identity := range identities {
		if identity.Category == "conference" && identity.Type == "text" {
			return true
		}
	}
	return false
}
