package access

import (
	"io"
	"time"

	"github.com/twstrike/coyim/client"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/event"
	"github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/tls"
	"github.com/twstrike/coyim/xmpp/data"
	xi "github.com/twstrike/coyim/xmpp/interfaces"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/otr3"
)

// EventHandler represents the main notifications that the session can emit
// It's really more an observer than an even handler
type EventHandler interface {
	RegisterCallback(title, instructions string, fields []interface{}) error
}

// Connector represents something that connect
type Connector interface {
	Connect()
}

// Session is an interface that defines the functionality of a Session
type Session interface {
	ApprovePresenceSubscription(string, string) error
	AutoApprove(string)
	AwaitVersionReply(<-chan data.Stanza, string)
	Close()
	CommandManager() client.CommandManager
	Config() *config.ApplicationConfig
	Conn() xi.Conn
	Connect(string, tls.Verifier) error
	ConversationManager() client.ConversationManager
	DenyPresenceSubscription(string, string) error
	DisplayName() string
	EncryptAndSendTo(string, string, string) (int, bool, error)
	GetConfig() *config.Account
	GroupDelimiter() string
	HandleConfirmOrDeny(string, bool)
	IsConnected() bool
	IsDisconnected() bool
	ManuallyEndEncryptedChat(string, string) error
	OtrEventHandler() map[string]*event.OtrEventHandler
	PrivateKeys() []otr3.PrivateKey
	R() *roster.List
	ReloadKeys()
	RemoveContact(string)
	RequestPresenceSubscription(string, string) error
	Send(string, string, string) error
	SendPing()
	SetCommandManager(client.CommandManager)
	SetConnectionLogger(io.Writer)
	SetConnector(Connector)
	SetLastActionTime(time.Time)
	SetSessionEventHandler(EventHandler)
	SetWantToBeOnline(bool)
	Subscribe(chan<- interface{})
	Timeout(data.Cookie, time.Time)
}

// Factory is a function that can create new Sessions
type Factory func(*config.ApplicationConfig, *config.Account, func(tls.Verifier) xi.Dialer) Session
