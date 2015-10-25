package session

import "sync"

// Event represents a XMPP event
type Event struct {
	EventType
	*Session
	From string
}

// EventType represents the type of XMPP event
type EventType int

// XMPP event types
const (
	Disconnected EventType = iota
	Connected
	RosterReceived
	IQReceived

	OTREnded
	OTRNewKeys

	SubscriptionRequest
	Subscribed
	Unsubscribe
)

var subscribers = struct {
	sync.RWMutex
	subs []chan<- Event
}{
	subs: make([]chan<- Event, 0),
}

// Subscribe subscribes the observer to XMPP events
func (s *Session) Subscribe(c chan<- Event) {
	subscribers.Lock()
	defer subscribers.Unlock()

	subscribers.subs = append(subscribers.subs, c)
}

func publishEvent(c chan<- Event, e Event) {
	//prevents from blocking the publisher if any subscriber is not listening to the channel
	go func(subscriber chan<- Event) {
		subscriber <- e
	}(c)
}

func (s *Session) publish(e EventType) {
	s.publishEvent(Event{
		EventType: e,
	})
}

func (s *Session) publishEvent(e Event) {
	subscribers.RLock()
	defer subscribers.RUnlock()

	for _, c := range subscribers.subs {
		e.Session = s
		publishEvent(c, e)
	}
}
