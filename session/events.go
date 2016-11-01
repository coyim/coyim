package session

import "github.com/twstrike/coyim/session/events"

// Subscribe subscribes the observer to XMPP events
func (s *session) Subscribe(c chan<- interface{}) {
	s.subscribers.Lock()
	defer s.subscribers.Unlock()

	if s.subscribers.subs == nil {
		s.subscribers.subs = make([]chan<- interface{}, 0)
	}

	s.subscribers.subs = append(s.subscribers.subs, c)
}

// Unsubscribe unsubscribes the observer to XMPP events
func (s *session) unsubscribe(c chan<- interface{}) {
	s.subscribers.Lock()
	defer s.subscribers.Unlock()

	for i, subs := range s.subscribers.subs {
		if subs == c {
			s.subscribers.subs = append(
				s.subscribers.subs[:i], s.subscribers.subs[i+1:]...,
			)
			return
		}
	}
}

func (s *session) publishEventTo(subscriber chan<- interface{}, e interface{}) {
	defer func() {
		if r := recover(); r != nil {
			//published to a closed channel
			s.unsubscribe(subscriber)
		}
	}()

	subscriber <- e
}

func (s *session) publish(e events.EventType) {
	s.publishEvent(events.Event{
		Session: s,
		Type:    e,
	})
}

func (s *session) publishPeerEvent(e events.PeerType, peer string) {
	s.publishEvent(events.Peer{
		Session: s,
		Type:    e,
		From:    peer,
	})
}

func (s *session) publishEvent(e interface{}) {
	s.subscribers.RLock()
	defer s.subscribers.RUnlock()

	for _, c := range s.subscribers.subs {
		go s.publishEventTo(c, e)
	}
}

func (s *session) publishSMPEvent(t events.SMPType, from, r, body string) {
	s.publishEvent(events.SMP{
		Type:     t,
		Session:  s,
		From:     from,
		Resource: r,
		Body:     body,
	})
}
