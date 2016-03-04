package mocks

import (
	"io"
	"time"

	"github.com/twstrike/coyim/client"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/event"
	"github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/session/access"
	"github.com/twstrike/coyim/xmpp/data"
	xi "github.com/twstrike/coyim/xmpp/interfaces"

	"github.com/twstrike/otr3"
)

type SessionMock struct{}

func (*SessionMock) ApprovePresenceSubscription(string, string) error {
	return nil
}

func (*SessionMock) AwaitVersionReply(<-chan data.Stanza, string) {}
func (*SessionMock) Close()                                       {}
func (*SessionMock) CommandManager() client.CommandManager {
	return nil
}
func (*SessionMock) Config() *config.ApplicationConfig {
	return nil
}
func (*SessionMock) Conn() xi.Conn {
	return nil
}
func (*SessionMock) Connect(string) error {
	return nil
}
func (*SessionMock) ConversationManager() client.ConversationManager {
	return nil
}
func (*SessionMock) DenyPresenceSubscription(string, string) error {
	return nil
}
func (*SessionMock) EncryptAndSendTo(string, string) error {
	return nil
}
func (*SessionMock) GetConfig() *config.Account {
	return nil
}
func (*SessionMock) GroupDelimiter() string {
	return ""
}
func (*SessionMock) HandleConfirmOrDeny(string, bool) {}
func (*SessionMock) IsConnected() bool {
	return false
}
func (*SessionMock) IsDisconnected() bool {
	return false
}
func (*SessionMock) OtrEventHandler() map[string]*event.OtrEventHandler {
	return nil
}
func (*SessionMock) PrivateKeys() []otr3.PrivateKey {
	return nil
}
func (*SessionMock) R() *roster.List {
	return nil
}
func (*SessionMock) ReloadKeys()          {}
func (*SessionMock) RemoveContact(string) {}
func (*SessionMock) RequestPresenceSubscription(string) error {
	return nil
}
func (*SessionMock) Send(string, string) error {
	return nil
}
func (*SessionMock) SetCommandManager(client.CommandManager)    {}
func (*SessionMock) SetConnectionLogger(io.Writer)              {}
func (*SessionMock) SetConnector(access.Connector)              {}
func (*SessionMock) SetLastActionTime(time.Time)                {}
func (*SessionMock) SetSessionEventHandler(access.EventHandler) {}
func (*SessionMock) SetWantToBeOnline(bool)                     {}
func (*SessionMock) Subscribe(chan<- interface{})               {}
func (*SessionMock) Timeout(data.Cookie, time.Time)             {}
