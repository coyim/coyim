package interfaces

import (
	"encoding/xml"
	"io"
	"reflect"
	"sync"

	"github.com/coyim/coyim/cache"
	"github.com/coyim/coyim/xmpp/data"
)

// Conn represents a connection to an XMPP server.
type Conn interface {
	Authenticate(string, string) error
	AuthenticationFailure() error
	BindResource() error
	Cancel(data.Cookie) bool
	Cache() cache.WithExpiry
	ChangePassword(string, string, string) error
	Close() error
	Config() *data.Config
	CustomStorage() map[xml.Name]reflect.Type
	DiscoveryFeatures(string) ([]string, bool)
	DiscoveryFeaturesAndIdentities(string) ([]data.DiscoveryIdentity, []string, bool)
	Features() data.StreamFeatures
	GetRosterDelimiter() (string, error)
	In() *xml.Decoder
	Lock() *sync.Mutex
	Next() (data.Stanza, error)
	OriginDomain() string
	Out() io.Writer
	Rand() io.Reader
	RawOut() io.WriteCloser
	ReadStanzas(chan<- data.Stanza) error
	RegisterAccount(string, string) (bool, error)
	RequestRoster() (<-chan data.Stanza, data.Cookie, error)
	RequestVCard() (<-chan data.Stanza, data.Cookie, error)
	Send(string, string, bool) error
	SendMessage(*data.Message) error
	SendIQ(string, string, interface{}) (<-chan data.Stanza, data.Cookie, error)
	SendIQReply(string, string, string, interface{}) error
	SendInitialStreamHeader() error
	SendPing() (reply <-chan data.Stanza, cookie data.Cookie, err error)
	SendPresence(string, string, string, string) error
	SendMUCPresence(string, *data.MUC) error
	ServerAddress() string
	SetInOut(*xml.Decoder, io.Writer)
	SetKeepaliveOut(io.Writer)
	SetRawOut(io.WriteCloser)
	SetServerAddress(string)
	SignalPresence(string) error

	SetChannelBinding([]byte)
	GetChannelBinding() []byte

	GetJIDResource() string
	SetJIDResource(string)

	HasSupportTo(entity string, features ...string) bool

	QueryServiceInformation(entity string) (*data.DiscoveryInfoQuery, error)
	QueryServiceItems(entity string) (*data.DiscoveryItemsQuery, error)
	EntityExists(entity string) (bool, error)

	ServerHasFeature(ns string) bool
}
