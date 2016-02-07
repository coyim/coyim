package access

import (
	"io"
	"time"

	"github.com/twstrike/coyim/client"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/event"
	"github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/xmpp"

	"github.com/twstrike/otr3"
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
	AwaitVersionReply(ch <-chan xmpp.Stanza, user string)
	Close()
	CommandManager() client.CommandManager
	Config() *config.ApplicationConfig
	Conn() *xmpp.Conn
	Connect(password string) error
	ConversationManager() client.ConversationManager
	DenyPresenceSubscription(string, string) error
	EncryptAndSendTo(string, string) error
	GetConfig() *config.Account
	GroupDelimiter() string
	HandleConfirmOrDeny(jid string, isConfirm bool)
	IsConnected() bool
	IsDisconnected() bool
	OtrEventHandler() map[string]*event.OtrEventHandler
	PrivateKeys() []otr3.PrivateKey
	R() *roster.List
	ReloadKeys()
	RemoveContact(string)
	RequestPresenceSubscription(jid string) error
	Send(to string, msg string) error
	SetCommandManager(client.CommandManager)
	SetConnectionLogger(l io.Writer)
	SetConnector(Connector)
	SetLastActionTime(t time.Time)
	SetSessionEventHandler(eh EventHandler)
	SetWantToBeOnline(bool)
	Subscribe(c chan<- interface{})
	Timeout(c xmpp.Cookie, t time.Time)
}

// Factory is a function that can create new Sessions
type Factory func(*config.ApplicationConfig, *config.Account) Session
