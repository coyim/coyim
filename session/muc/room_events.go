package muc

import "sync"

type roomSubscribers struct {
	sync.RWMutex
	subscribers []chan<- MUC
}

func newRoomSubscribers() *roomSubscribers {
	return &roomSubscribers{
		subscribers: make([]chan<- MUC, 0),
	}
}

func (s *roomSubscribers) subscribe(c chan<- MUC) {
	s.Lock()
	defer s.Unlock()
	s.subscribers = append(s.subscribers, c)
}

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
	defer s.RUnlock()

	for _, c := range s.subscribers {
		go s.publishEventTo(c, e)
	}
}
