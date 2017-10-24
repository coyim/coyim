package events

import (
	"time"

	"github.com/coyim/coyim/session/access"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/xmpp/data"
)

// Event represents a Session event
type Event struct {
	Type    EventType
	Session access.Session
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

// Peer represents an event associated to a peer
type Peer struct {
	Session access.Session
	Type    PeerType
	From    string
}

// Notification represents a notification event
type Notification struct {
	Session      access.Session
	Peer         string
	Notification string
}

// DelayedMessageSent represents the event that a delayed message is sent
type DelayedMessageSent struct {
	Session access.Session
	Peer    string
	Tracer  int
}

// PeerType represents the type of Peer event
type PeerType int

// Peer types
const (
	IQReceived PeerType = iota

	OTREnded
	OTRNewKeys
	OTRRenewedKeys

	SubscriptionRequest
	Subscribed
	Unsubscribe
)

// Presence represents a presence event
type Presence struct {
	Session access.Session
	*data.ClientPresence
	Gone bool
}

// Message represents a message event
type Message struct {
	Session   access.Session
	From      string
	Resource  string
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

// Log contains information one specific log event
type Log struct {
	Level   LogLevel
	Message string
}

// FileTransfer represents an event associated with file transfers
type FileTransfer struct {
	Session access.Session
	Peer    string

	Mime             string
	DateLastModified string
	Name             string
	Size             int64
	Description      string
	IsDirectory      bool

	Answer  chan<- *string // one time use
	Control *sdata.FileTransferControl
}

// SMP is an event related to SMP
type SMP struct {
	Type     SMPType
	Session  access.Session
	From     string
	Resource string
	Body     string
}

// SMPType denotes the type of an SMP event
type SMPType int

// SMP types
const (
	SecretNeeded SMPType = iota
	Success
	Failure
)
