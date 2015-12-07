package session

import (
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
	Connecting
	Connected
	ConnectionLost

	RosterReceived
	Ping
	PongReceived
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

// PresenceEvent represents a presence event
type PresenceEvent struct {
	*Session
	*xmpp.ClientPresence
	Gone bool
}

// MessageEvent represents a message event
type MessageEvent struct {
	*Session
	From      string
	When      time.Time
	Body      []byte
	Encrypted bool
}

// LogLevel is the current log level
type LogLevel int

// The different available log levels
const (
	Info LogLevel = iota
	Warn
	Alert
)

// LogEvent contains information one specific log event
type LogEvent struct {
	Level   LogLevel
	Message string
}

// Subscribe subscribes the observer to XMPP events
func (s *Session) Subscribe(c chan<- interface{}) {
	s.subscribers.Lock()
	defer s.subscribers.Unlock()

	if s.subscribers.subs == nil {
		s.subscribers.subs = make([]chan<- interface{}, 0)
	}

	s.subscribers.subs = append(s.subscribers.subs, c)
}

// Unsubscribe unsubscribes the observer to XMPP events
func (s *Session) Unsubscribe(c chan<- interface{}) {
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

func (s *Session) publishEventTo(subscriber chan<- interface{}, e interface{}) {
	defer func() {
		if r := recover(); r != nil {
			//published to a closed channel
			s.Unsubscribe(subscriber)
		}
	}()

	subscriber <- e
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
	s.subscribers.RLock()
	defer s.subscribers.RUnlock()

	for _, c := range s.subscribers.subs {
		go s.publishEventTo(c, e)
	}
}
