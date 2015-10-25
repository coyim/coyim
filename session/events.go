package session

import "sync"

type Event struct {
	EventType
	*Session
}

type EventType int

const (
	Disconnected EventType = iota
	Connected
)

var subscribers = struct {
	sync.RWMutex
	subs []chan<- Event
}{
	subs: make([]chan<- Event, 0),
}

func (s *Session) Subscribe(c chan<- Event) {
	subscribers.Lock()
	defer subscribers.Unlock()

	subscribers.subs = append(subscribers.subs, c)
}

func (s *Session) publish(e EventType) {
	subscribers.RLock()
	defer subscribers.RUnlock()

	for _, c := range subscribers.subs {
		//prevents from blocking the publisher if any subscriber is not listening to the channel
		go func(subscriber chan<- Event) {
			subscriber <- Event{e, s}
		}(c)
	}
}
