package access

import (
	"bytes"
	"io"
	"time"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/event"
	"github.com/coyim/coyim/otr_client"
	"github.com/coyim/coyim/roster"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/tls"
	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"

	"github.com/coyim/otr3"
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
	ApprovePresenceSubscription(data.JIDWithoutResource, string) error
	AutoApprove(string)
	AwaitVersionReply(<-chan data.Stanza, string)
	Close()
	CommandManager() otr_client.CommandManager
	Config() *config.ApplicationConfig
	Conn() xi.Conn
	Connect(string, tls.Verifier) error
	ConversationManager() otr_client.ConversationManager
	CreateSymmetricKeyFor(data.JID) []byte
	DenyPresenceSubscription(data.JIDWithoutResource, string) error
	DisplayName() string
	EncryptAndSendTo(data.JIDWithoutResource, data.JIDResource, string) (int, bool, error)
	GetConfig() *config.Account
	GetInMemoryLog() *bytes.Buffer
	GroupDelimiter() string
	HandleConfirmOrDeny(data.JIDWithoutResource, bool)
	IsConnected() bool
	IsDisconnected() bool
	ManuallyEndEncryptedChat(data.JIDWithoutResource, data.JIDResource) error
	OtrEventHandler() map[string]*event.OtrEventHandler
	PrivateKeys() []otr3.PrivateKey
	R() *roster.List
	ReloadKeys()
	RemoveContact(string)
	RequestPresenceSubscription(data.JIDWithoutResource, string) error
	Send(data.JIDWithoutResource, data.JIDResource, string) error
	SendPing()
	SetCommandManager(otr_client.CommandManager)
	SetConnectionLogger(io.Writer)
	SetConnector(Connector)
	SetLastActionTime(time.Time)
	SetSessionEventHandler(EventHandler)
	SetWantToBeOnline(bool)
	Subscribe(chan<- interface{})
	Timeout(data.Cookie, time.Time)
	Info(string)
	Warn(string)
	SendIQError(*data.ClientIQ, interface{})
	SendIQResult(*data.ClientIQ, interface{})
	PublishEvent(interface{})
	SendFileTo(data.JID, string, bool) *sdata.FileTransferControl
	SendDirTo(data.JID, string, bool) *sdata.FileTransferControl
	StartSMP(data.JIDWithoutResource, data.JIDResource, string, string)
	FinishSMP(data.JIDWithoutResource, data.JIDResource, string)
	AbortSMP(data.JIDWithoutResource, data.JIDResource)
}

// Factory is a function that can create new Sessions
type Factory func(*config.ApplicationConfig, *config.Account, func(tls.Verifier) xi.Dialer) Session
