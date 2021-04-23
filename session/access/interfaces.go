package access

import (
	"bytes"
	"time"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/otrclient"
	"github.com/coyim/coyim/roster"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/session/muc"
	mdata "github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/tls"
	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"

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

// CanSendIQ represents anything that can send IQ results or errors
type CanSendIQ interface {
	SendIQResult(*data.ClientIQ, interface{})
	SendIQError(*data.ClientIQ, interface{})
}

// HasSymmetricKey represents anything that can manage symmetric OTR keys
type HasSymmetricKey interface {
	CreateSymmetricKeyFor(jid.Any) []byte
	GetAndWipeSymmetricKeyFor(jid.Any) []byte
}

// Publisher is anything that can publish events
type Publisher interface {
	PublishEvent(interface{})
}

// SMP contains functionality related to SMP
type SMP interface {
	StartSMP(jid.WithResource, string, string)
	FinishSMP(jid.WithResource, string)
	AbortSMP(jid.WithResource)
}

// Rooms contains all the functionality for MUC
type Rooms interface {
	HasRoom(jid.Bare, chan<- *muc.RoomListing) (<-chan bool, <-chan error)
	GetRooms(jid.Domain, string) (<-chan *muc.RoomListing, <-chan *muc.ServiceListing, <-chan error)
	JoinRoom(jid.Bare, string, string) error
	CreateInstantRoom(jid.Bare) (<-chan bool, <-chan error)
	CreateReservedRoom(jid.Bare) (<-chan *muc.RoomConfigForm, <-chan error)
	SubmitRoomConfigurationForm(jid.Bare, *muc.RoomConfigForm) (<-chan bool, <-chan error)
	CancelRoomConfiguration(jid.Bare) <-chan error
	GetChatServices(jid.Domain) (<-chan jid.Domain, <-chan error, func())
	GetRoomListing(jid.Bare, chan<- *muc.RoomListing)
	GetRoomInformation(jid.Bare)
	LeaveRoom(room jid.Bare, nickname string) (<-chan bool, <-chan error)
	DestroyRoom(room jid.Bare, reason string, alternativeRoom jid.Bare, password string) (<-chan bool, <-chan error)
	UpdateOccupantAffiliation(roomID jid.Bare, occupantNickname string, occupantRealJID jid.Full, affiliation mdata.Affiliation, reason string) (<-chan bool, <-chan error)
	UpdateOccupantRole(roomID jid.Bare, occupantNickname string, role mdata.Role, reason string) (<-chan bool, <-chan error)
	NewRoom(jid.Bare) *muc.Room
	SendMUCMessage(to, from, body string) error
	GetRoomBanList(jid.Bare) (<-chan []*muc.RoomBanListItem, <-chan error)
	ModifyRoomBanList(jid.Bare, []*muc.RoomBanListItem) (<-chan bool, <-chan error)
}

// Connection contains the connection related functionality
type Connection interface {
	Close()
	IsConnected() bool
	IsDisconnected() bool
	Connect(string, tls.Verifier) error
	SetConnector(Connector)
}

// Subscription contains functionality related to subscriptions
type Subscription interface {
	ApprovePresenceSubscription(jid.WithoutResource, string) error
	AutoApprove(string)
	DenyPresenceSubscription(jid.WithoutResource, string) error
	HandleConfirmOrDeny(jid.WithoutResource, bool)
	RequestPresenceSubscription(jid.WithoutResource, string) error
	RemoveContact(string)
}

// Transfer gives access to file and directory transfer functionality
type Transfer interface {
	SendFileTo(jid.Any, string, func() bool, func(bool)) *sdata.FileTransferControl
	SendDirTo(jid.Any, string, func() bool, func(bool)) *sdata.FileTransferControl
}

// EncryptedChat allows control over OTR functionality
type EncryptedChat interface {
	CommandManager() otrclient.CommandManager
	ConversationManager() otrclient.ConversationManager
	PrivateKeys() []otr3.PrivateKey
	SetCommandManager(otrclient.CommandManager)
	ManuallyEndEncryptedChat(jid.Any) error
	ReloadKeys()
}

// Roster exposes roster functionality
type Roster interface {
	GroupDelimiter() string
	R() *roster.List
}

// Sending gives access to sending functionality
type Sending interface {
	EncryptAndSendTo(jid.Any, string) (int, bool, error)
	Send(jid.Any, string, bool) error
	SendPing()
}

// ConnectionData gives access to information about the connection and session
type ConnectionData interface {
	DisplayName() string
	SetLastActionTime(time.Time)
	SetWantToBeOnline(bool)
	Timeout(data.Cookie, time.Time)
}

// Logging gives access to the in memory log
type Logging interface {
	GetInMemoryLog() *bytes.Buffer
}

// Events allow you to subscribe to events
type Events interface {
	Subscribe(chan<- interface{})
}

// Version gives access to version functionality
type Version interface {
	AwaitVersionReply(<-chan data.Stanza, string)
}

// Session is an interface that defines the functionality of a Session
type Session interface {
	CanSendIQ
	HasSymmetricKey
	Publisher
	coylog.Has
	config.Has
	config.HasApplication
	xi.Has
	SMP
	Rooms
	Connection
	Subscription
	Transfer
	EncryptedChat
	Roster
	Sending
	ConnectionData
	Logging
	Events
	Version
}

// Factory is a function that can create new Sessions
type Factory func(*config.ApplicationConfig, *config.Account, func(tls.Verifier, tls.Factory) xi.Dialer) Session
