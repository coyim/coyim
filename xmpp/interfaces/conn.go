package interfaces

import (
	"encoding/xml"
	"io"
	"reflect"
	"sync"

	"github.com/twstrike/coyim/xmpp/data"
)

type Conn interface {
	SetInOut(*xml.Decoder, io.Writer)
	SendInitialStreamHeader() error
	SetRawOut(io.WriteCloser)
	SetKeepaliveOut(io.Writer)
	Features() data.StreamFeatures
	RegisterAccount(user, password string) (bool, error)
	Authenticate(user, password string) error
	BindResource() error
	Config() *data.Config
	In() *xml.Decoder
	Out() io.Writer
	OriginDomain() string
	Lock() *sync.Mutex
	CustomStorage() map[xml.Name]reflect.Type
	SendPresence(to, typ, id string) error
	Send(to, msg string) error
	SendIQ(to, typ string, value interface{}) (reply chan data.Stanza, cookie data.Cookie, err error)
	SendIQReply(to, typ, id string, value interface{}) error
	Cancel(cookie data.Cookie) bool
	RequestRoster() (<-chan data.Stanza, data.Cookie, error)
	ReadStanzas(stanzaChan chan<- data.Stanza) error
	Close() error
	GetRosterDelimiter() (string, error)
	SignalPresence(state string) error
	Next() (stanza data.Stanza, err error)
}
