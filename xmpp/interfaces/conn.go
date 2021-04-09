package interfaces

import (
	"encoding/xml"
	"io"
	"reflect"
	"sync"

	"github.com/coyim/coyim/cache"
	"github.com/coyim/coyim/xmpp/data"
)

// Has represents any object that has a connection and can give access to it
type Has interface {
	Conn() Conn
}

// Authentication contains authentication related functionality for the XMPP connection
type Authentication interface {
	Authenticate(string, string) error
	AuthenticationFailure() error
}

// Registration contains registration related functionality for the XMPP connection
type Registration interface {
	ChangePassword(string, string, string) error
	RegisterAccount(string, string) (bool, error)
}

// Resource contains resource related functionality for the XMPP connection
type Resource interface {
	GetJIDResource() string
	SetJIDResource(string)
	BindResource() error
	OriginDomain() string
}

// ConnectionData contains connection information for the XMPP connection
type ConnectionData interface {
	ServerAddress() string
	SetInOut(*xml.Decoder, io.Writer)
	SetKeepaliveOut(io.Writer)
	SetRawOut(io.WriteCloser)
	SetServerAddress(string)
}

// Presence contains functionality related to signaling presence for the XMPP connection
type Presence interface {
	SendPresence(string, string, string, string) error
	SendMUCPresence(string, *data.MUC) error
	SignalPresence(string) error
}

// Features contains functionality for managing discovery of features and identities for the XMPP connection
type Features interface {
	DiscoveryFeatures(string) ([]string, bool)
	DiscoveryFeaturesAndIdentities(string) ([]data.DiscoveryIdentity, []string, bool)
	Features() data.StreamFeatures
	HasSupportTo(entity string, features ...string) bool
	QueryServiceInformation(entity string) (*data.DiscoveryInfoQuery, error)
	QueryServiceItems(entity string) (*data.DiscoveryItemsQuery, error)
	EntityExists(entity string) (bool, error)
	ServerHasFeature(ns string) bool
}

// Transport contains functionality for managing the transport of information over the XMPP connection
type Transport interface {
	SetChannelBinding([]byte)
	GetChannelBinding() []byte
	Close() error
}

// InputOutput manages the low level input and output functionality of the XMPP connection
type InputOutput interface {
	In() *xml.Decoder
	Next() (data.Stanza, error)
	Out() io.Writer
	RawOut() io.WriteCloser
}

// Configuration contains functionality related to the configuration of the XMPP connection
type Configuration interface {
	Cache() cache.WithExpiry
	Config() *data.Config
	CustomStorage() map[xml.Name]reflect.Type
	Lock() *sync.Mutex
	Rand() io.Reader
}

// Requests manages the different request functionalities for the XMPP connection
type Requests interface {
	GetRosterDelimiter() (string, error)
	RequestRoster() (<-chan data.Stanza, data.Cookie, error)
	RequestVCard() (<-chan data.Stanza, data.Cookie, error)
}

// Reading gives access to reading of stanzas over the XMPP connection
type Reading interface {
	ReadStanzas(chan<- data.Stanza) error
}

// Sending allows for sending of different kinds of data over the XMPP connection
type Sending interface {
	Send(string, string, bool) error
	SendMessage(*data.Message) error
	SendIQ(string, string, interface{}) (<-chan data.Stanza, data.Cookie, error)
	SendIQReply(string, string, string, interface{}) error
	SendInitialStreamHeader() error
	SendPing() (reply <-chan data.Stanza, cookie data.Cookie, err error)
	Cancel(data.Cookie) bool
}

// Conn represents a connection to an XMPP server.
type Conn interface {
	ConnectionData
	Authentication
	Registration
	Resource
	Presence
	Features
	Transport
	InputOutput
	Configuration
	Requests
	Reading
	Sending
}
