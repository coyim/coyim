package events

import (
	"time"

	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// Event represents a Session event
type Event struct {
	Type EventType
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
	Type PeerType
	// This can be either with or without Resource depending on the peer type
	From jid.Any
}

// Notification represents a notification event
type Notification struct {
	Peer         jid.Any
	Notification string
}

// DelayedMessageSent represents the event that a delayed message is sent
type DelayedMessageSent struct {
	Peer   jid.Any
	Tracer int
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
	*data.ClientPresence
	Gone bool
}

// Message represents a message event
type Message struct {
	From      jid.Any
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
	Peer jid.WithResource

	Mime             string
	DateLastModified string
	Name             string
	Size             int64
	Description      string
	IsDirectory      bool
	Encrypted        bool

	Answer  chan<- *string // one time use
	Control *sdata.FileTransferControl
}

// SMP is an event related to SMP
type SMP struct {
	Type     SMPType
	From     jid.WithResource
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
