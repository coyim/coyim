package session

import (
	"sync"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/jid"
)

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

func (s *session) publishEventTo(subscriber chan<- interface{}, e interface{}, subWait *sync.WaitGroup) {
	defer func() {
		subWait.Done()
		if r := recover(); r != nil {
			//published to a closed channel
			s.unsubscribe(subscriber)
		}
	}()

	subscriber <- e
}

func (s *session) publish(e events.EventType) {
	s.publishEvent(events.Event{
		Type: e,
	})
}

func (s *session) publishPeerEvent(e events.PeerType, peer jid.Any) {
	s.publishEvent(events.Peer{
		Type: e,
		From: peer,
	})
}

func (s *session) publishEvent(e interface{}) {
	s.pendingEventsLock.Lock()
	s.pendingEvents++
	s.pendingEventsLock.Unlock()

	s.subscribers.RLock()
	defer s.subscribers.RUnlock()

	var allSubs sync.WaitGroup

	allSubs.Add(len(s.subscribers.subs))
	for _, c := range s.subscribers.subs {
		go s.publishEventTo(c, e, &allSubs)
	}

	go func() {
		allSubs.Wait()

		s.pendingEventsLock.Lock()

		s.pendingEvents--
		if s.pendingEvents == 0 && s.eventsReachedZero != nil {
			s.eventsReachedZero <- true
		}

		s.pendingEventsLock.Unlock()
	}()
}

func (s *session) PublishEvent(e interface{}) {
	s.publishEvent(e)
}

func (s *session) publishSMPEvent(t events.SMPType, peer jid.WithResource, body string) {
	s.publishEvent(events.SMP{
		Type: t,
		From: peer,
		Body: body,
	})
}
