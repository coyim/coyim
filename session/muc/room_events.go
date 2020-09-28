package muc

import "sync"

type roomSubscribers struct {
	sync.RWMutex
	subscribers []chan<- MUC
}

func newRoomSubscribers() *roomSubscribers {
	return &roomSubscribers{}
}

func (s *roomSubscribers) subscribe(c chan<- MUC) {
	s.Lock()
	defer s.Unlock()
	s.subscribers = append(s.subscribers, c)
}

// TODO: We will remove this when we remove the Room.Unsubscribe functionality

func (s *roomSubscribers) unsubscribe(c chan<- MUC) {
	s.Lock()
	defer s.Unlock()

	for i, subs := range s.subscribers {
		if subs == c {
			s.subscribers = append(
				s.subscribers[:i], s.subscribers[i+1:]...,
			)
			return
		}
	}
}

func (s *roomSubscribers) publishEventTo(subscriber chan<- MUC, e MUC) {
	defer func() {
		if r := recover(); r != nil {
			// published to a closed channel
			s.unsubscribe(subscriber)
		}
	}()

	subscriber <- e
}

func (s *roomSubscribers) publishEvent(e MUC) {
	s.RLock()
	subs := append([]chan<- MUC{}, s.subscribers...)
	s.RUnlock()

	for _, c := range subs {
		go s.publishEventTo(c, e)
	}
}
