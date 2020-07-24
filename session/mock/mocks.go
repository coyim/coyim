package mocks

import (
	"bytes"
	"time"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/otrclient"
	"github.com/coyim/coyim/roster"
	"github.com/coyim/coyim/session/access"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/tls"
	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"

	"github.com/coyim/otr3"
)

// SessionMock is a mock of the Session interface
type SessionMock struct{}

// ApprovePresenceSubscription is the implementation for Session interface
func (*SessionMock) ApprovePresenceSubscription(jid.WithoutResource, string) error {
	return nil
}

// AwaitVersionReply is the implementation for Session interface
func (*SessionMock) AwaitVersionReply(<-chan data.Stanza, string) {}

// Close is the implementation for Session interface
func (*SessionMock) Close() {}

// AutoApprove is the implementation for Session interface
func (*SessionMock) AutoApprove(string) {}

// CommandManager is the implementation for Session interface
func (*SessionMock) CommandManager() otrclient.CommandManager {
	return nil
}

// Config is the implementation for Session interface
func (*SessionMock) Config() *config.ApplicationConfig {
	return nil
}

// Conn is the implementation for Session interface
func (*SessionMock) Conn() xi.Conn {
	return nil
}

// Connect is the implementation for Session interface
func (*SessionMock) Connect(string, tls.Verifier) error {
	return nil
}

// ConversationManager is the implementation for Session interface
func (*SessionMock) ConversationManager() otrclient.ConversationManager {
	return nil
}

// DenyPresenceSubscription is the implementation for Session interface
func (*SessionMock) DenyPresenceSubscription(jid.WithoutResource, string) error {
	return nil
}

// DisplayName is the implementation for Session interface
func (*SessionMock) DisplayName() string {
	return ""
}

// EncryptAndSendTo is the implementation for Session interface
func (*SessionMock) EncryptAndSendTo(jid.Any, string) (int, bool, error) {
	return 0, false, nil
}

// GetConfig is the implementation for Session interface
func (*SessionMock) GetConfig() *config.Account {
	return nil
}

// GetInMemoryLog is the implementation for Session interface
func (*SessionMock) GetInMemoryLog() *bytes.Buffer {
	return nil
}

// GroupDelimiter is the implementation for Session interface
func (*SessionMock) GroupDelimiter() string {
	return ""
}

// HandleConfirmOrDeny is the implementation for Session interface
func (*SessionMock) HandleConfirmOrDeny(jid.WithoutResource, bool) {}

// IsConnected is the implementation for Session interface
func (*SessionMock) IsConnected() bool {
	return false
}

// IsDisconnected is the implementation for Session interface
func (*SessionMock) IsDisconnected() bool {
	return false
}

// ManuallyEndEncryptedChat is the implementation for Session interface
func (*SessionMock) ManuallyEndEncryptedChat(jid.Any) error {
	return nil
}

// PrivateKeys is the implementation for Session interface
func (*SessionMock) PrivateKeys() []otr3.PrivateKey {
	return nil
}

// R is the implementation for Session interface
func (*SessionMock) R() *roster.List {
	return nil
}

// ReloadKeys is the implementation for Session interface
func (*SessionMock) ReloadKeys() {}

// RemoveContact is the implementation for Session interface
func (*SessionMock) RemoveContact(string) {}

// RequestPresenceSubscription is the implementation for Session interface
func (*SessionMock) RequestPresenceSubscription(jid.WithoutResource, string) error {
	return nil
}

// Send is the implementation for Session interface
func (*SessionMock) Send(jid.Any, string, bool) error {
	return nil
}

// SendPing is the implementation for Session interface
func (*SessionMock) SendPing() {}

// SetCommandManager is the implementation for Session interface
func (*SessionMock) SetCommandManager(otrclient.CommandManager) {}

// SetConnector is the implementation for Session interface
func (*SessionMock) SetConnector(access.Connector) {}

// SetLastActionTime is the implementation for Session interface
func (*SessionMock) SetLastActionTime(time.Time) {}

// SetSessionEventHandler is the implementation for Session interface
func (*SessionMock) SetSessionEventHandler(access.EventHandler) {}

// SetWantToBeOnline is the implementation for Session interface
func (*SessionMock) SetWantToBeOnline(bool) {}

// Subscribe is the implementation for Session interface
func (*SessionMock) Subscribe(chan<- interface{}) {}

// Timeout is the implementation for Session interface
func (*SessionMock) Timeout(data.Cookie, time.Time) {}

// Warn is the implementation for Session interface
func (*SessionMock) Warn(string) {}

// Info is the implementation for Session interface
func (*SessionMock) Info(string) {}

// StartSMP is the implementation for Session interface
func (*SessionMock) StartSMP(jid.WithResource, string, string) {}

// FinishSMP is the implementation for Session interface
func (*SessionMock) FinishSMP(jid.WithResource, string) {}

// AbortSMP is the implementation for Session interface
func (*SessionMock) AbortSMP(jid.WithResource) {}

// PublishEvent is the implementation for Session interface
func (*SessionMock) PublishEvent(interface{}) {}

// SendIQError is the implementation for Session interface
func (*SessionMock) SendIQError(*data.ClientIQ, interface{}) {}

// SendIQResult is the implementation for Session interface
func (*SessionMock) SendIQResult(*data.ClientIQ, interface{}) {}

// SendFileTo is the implementation for Session interface
func (*SessionMock) SendFileTo(jid.Any, string, func() bool, func(bool)) *sdata.FileTransferControl {
	return nil
}

// SendDirTo is the implementation for Session interface
func (*SessionMock) SendDirTo(jid.Any, string, func() bool, func(bool)) *sdata.FileTransferControl {
	return nil
}

// CreateSymmetricKeyFor is the implementation for Session interface
func (*SessionMock) CreateSymmetricKeyFor(jid.Any) []byte {
	return nil
}

// GetAndWipeSymmetricKeyFor is the implementation for Session interface
func (*SessionMock) GetAndWipeSymmetricKeyFor(jid.Any) []byte {
	return nil
}

// GetRooms is the implementation for Session interface
func (*SessionMock) GetRooms(jid.Domain) (<-chan *muc.RoomListing, <-chan *muc.ServiceListing, <-chan error) {
	return nil, nil, nil
}
