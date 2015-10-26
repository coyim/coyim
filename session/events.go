package session

import (
	"sync"
	"time"

	"github.com/twstrike/coyim/xmpp"
)

// Event represents a Session event
type Event struct {
	Type EventType
	*Session
}

// EventType represents the type of Session event
type EventType int

// Session event types
const (
	Disconnected EventType = iota
	Connected
	RosterReceived
)

// PeerEvent represents an event associated to a peer
type PeerEvent struct {
	*Session
	Type PeerEventType
	From string
}

// PeerEventType represents the type of Peer event
type PeerEventType int

// PeerEvent types
const (
	IQReceived PeerEventType = iota

	OTREnded
	OTRNewKeys

	SubscriptionRequest
	Subscribed
	Unsubscribe
)

type PresenceEvent struct {
	*Session
	*xmpp.ClientPresence
	Gone bool
}

type MessageEvent struct {
	*Session
	From      string
	When      time.Time
	Body      []byte
	Encrypted bool
}

var subscribers = struct {
	sync.RWMutex
	subs []chan<- interface{}
}{
	subs: make([]chan<- interface{}, 0),
}

// Subscribe subscribes the observer to XMPP events
func (s *Session) Subscribe(c chan<- interface{}) {
	subscribers.Lock()
	defer subscribers.Unlock()

	subscribers.subs = append(subscribers.subs, c)
}

func publishEvent(c chan<- interface{}, e interface{}) {
	//prevents from blocking the publisher if any subscriber is not listening to the channel
	go func(subscriber chan<- interface{}) {
		subscriber <- e
	}(c)
}

func (s *Session) publish(e EventType) {
	s.publishEvent(Event{
		Session: s,
		Type:    e,
	})
}

func (s *Session) publishPeerEvent(e PeerEventType, peer string) {
	s.publishEvent(PeerEvent{
		Session: s,
		Type:    e,
		From:    peer,
	})
}

func (s *Session) publishEvent(e interface{}) {
	subscribers.RLock()
	defer subscribers.RUnlock()

	for _, c := range subscribers.subs {
		publishEvent(c, e)
	}
}
