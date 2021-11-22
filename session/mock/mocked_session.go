package mock

import (
	"bytes"
	"time"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/otrclient"
	"github.com/coyim/coyim/roster"
	"github.com/coyim/coyim/session/access"
	sdata "github.com/coyim/coyim/session/data"
	"github.com/coyim/coyim/session/muc"
	mdata "github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/tls"
	"github.com/coyim/coyim/xmpp/data"
	xi "github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/otr3"
	mck "github.com/stretchr/testify/mock"
)

// MockedSession is a mock of the Session interface
type MockedSession struct {
	mck.Mock
}

// ApprovePresenceSubscription is the implementation for Session interface
func (m *MockedSession) ApprovePresenceSubscription(v1 jid.WithoutResource, v2 string) error {
	return m.Called(v1, v2).Error(0)
}

// AwaitVersionReply is the implementation for Session interface
func (m *MockedSession) AwaitVersionReply(v1 <-chan data.Stanza, v2 string) {
	m.Called(v1, v2)
}

// Close is the implementation for Session interface
func (m *MockedSession) Close() {
	m.Called()
}

// AutoApprove is the implementation for Session interface
func (m *MockedSession) AutoApprove(v1 string) {
	m.Called(v1)
}

// CommandManager is the implementation for Session interface
func (m *MockedSession) CommandManager() otrclient.CommandManager {
	return m.Called().Get(0).(otrclient.CommandManager)
}

// Config is the implementation for Session interface
func (m *MockedSession) Config() *config.ApplicationConfig {
	return m.Called().Get(0).(*config.ApplicationConfig)
}

// Conn is the implementation for Session interface
func (m *MockedSession) Conn() xi.Conn {
	return m.Called().Get(0).(xi.Conn)
}

// Connect is the implementation for Session interface
func (m *MockedSession) Connect(v1 string, v2 tls.Verifier) error {
	return m.Called(v1, v2).Error(0)
}

// ConversationManager is the implementation for Session interface
func (m *MockedSession) ConversationManager() otrclient.ConversationManager {
	return m.Called().Get(0).(otrclient.ConversationManager)
}

// DenyPresenceSubscription is the implementation for Session interface
func (m *MockedSession) DenyPresenceSubscription(v1 jid.WithoutResource, v2 string) error {
	return m.Called(v1, v2).Error(0)
}

// DisplayName is the implementation for Session interface
func (m *MockedSession) DisplayName() string {
	return m.Called().String(0)
}

// EncryptAndSendTo is the implementation for Session interface
func (m *MockedSession) EncryptAndSendTo(v1 jid.Any, v2 string) (int, bool, error) {
	args := m.Called(v1, v2)
	return args.Int(0), args.Bool(1), args.Error(2)
}

// GetConfig is the implementation for Session interface
func (m *MockedSession) GetConfig() *config.Account {
	return m.Called().Get(0).(*config.Account)
}

// GetInMemoryLog is the implementation for Session interface
func (m *MockedSession) GetInMemoryLog() *bytes.Buffer {
	return m.Called().Get(0).(*bytes.Buffer)
}

// GroupDelimiter is the implementation for Session interface
func (m *MockedSession) GroupDelimiter() string {
	return m.Called().String(0)
}

// HandleConfirmOrDeny is the implementation for Session interface
func (m *MockedSession) HandleConfirmOrDeny(v1 jid.WithoutResource, v2 bool) {
	m.Called(v1, v2)
}

// IsConnected is the implementation for Session interface
func (m *MockedSession) IsConnected() bool {
	return m.Called().Bool(0)
}

// IsDisconnected is the implementation for Session interface
func (m *MockedSession) IsDisconnected() bool {
	return m.Called().Bool(0)
}

// ManuallyEndEncryptedChat is the implementation for Session interface
func (m *MockedSession) ManuallyEndEncryptedChat(v1 jid.Any) error {
	return m.Called(v1).Error(0)
}

// PrivateKeys is the implementation for Session interface
func (m *MockedSession) PrivateKeys() []otr3.PrivateKey {
	return m.Called().Get(0).([]otr3.PrivateKey)
}

// R is the implementation for Session interface
func (m *MockedSession) R() *roster.List {
	return m.Called().Get(0).(*roster.List)
}

// ReloadKeys is the implementation for Session interface
func (m *MockedSession) ReloadKeys() {
	m.Called()
}

// RemoveContact is the implementation for Session interface
func (m *MockedSession) RemoveContact(v1 string) {
	m.Called(v1)
}

// RequestPresenceSubscription is the implementation for Session interface
func (m *MockedSession) RequestPresenceSubscription(v1 jid.WithoutResource, v2 string) error {
	return m.Called(v1, v2).Error(0)
}

// Send is the implementation for Session interface
func (m *MockedSession) Send(v1 jid.Any, v2 string, v3 bool) error {
	return m.Called(v1, v2, v3).Error(0)
}

// SendMUCMessage is the implementation for Session interface
func (m *MockedSession) SendMUCMessage(to, from, body string) error {
	return m.Called(to, from, body).Error(0)
}

// UpdateRoomSubject is the implementation for Session interface
func (m *MockedSession) UpdateRoomSubject(roomID jid.Bare, from, subject string) error {
	return m.Called(roomID, from, subject).Error(0)
}

// UpdateOccupantAffiliations is the implementation for Session interface
func (m *MockedSession) UpdateOccupantAffiliations(v1 jid.Bare, v2 []*muc.RoomOccupantItem) (<-chan bool, <-chan error) {
	args := m.Called(v1, v2)
	return args.Get(0).(<-chan bool), args.Get(1).(<-chan error)
}

// GetRoomOccupantsByAffiliation is the implementation for Session interface
func (m *MockedSession) GetRoomOccupantsByAffiliation(roomID jid.Bare, a mdata.Affiliation) (<-chan []*muc.RoomOccupantItem, <-chan error) {
	args := m.Called(roomID, a)
	return args.Get(0).(<-chan []*muc.RoomOccupantItem), args.Get(1).(<-chan error)
}

// SendPing is the implementation for Session interface
func (m *MockedSession) SendPing() {
	m.Called()
}

// SetCommandManager is the implementation for Session interface
func (m *MockedSession) SetCommandManager(v1 otrclient.CommandManager) {
	m.Called(v1)
}

// SetConnector is the implementation for Session interface
func (m *MockedSession) SetConnector(v1 access.Connector) {
	m.Called(v1)
}

// SetLastActionTime is the implementation for Session interface
func (m *MockedSession) SetLastActionTime(v1 time.Time) {
	m.Called(v1)
}

// SetWantToBeOnline is the implementation for Session interface
func (m *MockedSession) SetWantToBeOnline(v1 bool) {
	m.Called(v1)
}

// Subscribe is the implementation for Session interface
func (m *MockedSession) Subscribe(v1 chan<- interface{}) {
	m.Called(v1)
}

// Timeout is the implementation for Session interface
func (m *MockedSession) Timeout(v1 data.Cookie, v2 time.Time) {
	m.Called(v1, v2)
}

// StartSMP is the implementation for Session interface
func (m *MockedSession) StartSMP(v1 jid.WithResource, v2 string, v3 string) {
	m.Called(v1, v2, v3)
}

// FinishSMP is the implementation for Session interface
func (m *MockedSession) FinishSMP(v1 jid.WithResource, v2 string) {
	m.Called(v1, v2)
}

// AbortSMP is the implementation for Session interface
func (m *MockedSession) AbortSMP(v1 jid.WithResource) {
	m.Called(v1)
}

// PublishEvent is the implementation for Session interface
func (m *MockedSession) PublishEvent(v1 interface{}) {
	m.Called(v1)
}

// SendIQError is the implementation for Session interface
func (m *MockedSession) SendIQError(v1 *data.ClientIQ, v2 interface{}) {
	m.Called(v1, v2)
}

// SendIQResult is the implementation for Session interface
func (m *MockedSession) SendIQResult(v1 *data.ClientIQ, v2 interface{}) {
	m.Called(v1, v2)
}

// SendFileTo is the implementation for Session interface
func (m *MockedSession) SendFileTo(v1 jid.Any, v2 string, v3 func() bool, v4 func(bool)) *sdata.FileTransferControl {
	return m.Called(v1, v2, v3, v4).Get(0).(*sdata.FileTransferControl)
}

// SendDirTo is the implementation for Session interface
func (m *MockedSession) SendDirTo(v1 jid.Any, v2 string, v3 func() bool, v4 func(bool)) *sdata.FileTransferControl {
	return m.Called(v1, v2, v3, v4).Get(0).(*sdata.FileTransferControl)
}

// CreateSymmetricKeyFor is the implementation for Session interface
func (m *MockedSession) CreateSymmetricKeyFor(v1 jid.Any) []byte {
	return m.Called(v1).Get(0).([]byte)
}

// GetAndWipeSymmetricKeyFor is the implementation for Session interface
func (m *MockedSession) GetAndWipeSymmetricKeyFor(v1 jid.Any) []byte {
	return m.Called(v1).Get(0).([]byte)
}

// HasRoom is the implementation for Session interface
func (m *MockedSession) HasRoom(v1 jid.Bare, v2 chan<- *muc.RoomListing) (<-chan bool, <-chan error) {
	args := m.Called(v1, v2)
	return args.Get(0).(<-chan bool), args.Get(1).(<-chan error)
}

// GetRoomListing is the implementation for Session interface
func (m *MockedSession) GetRoomListing(v1 jid.Bare, v2 chan<- *muc.RoomListing) {
	m.Called(v1, v2)
}

// RefreshRoomProperties is the implementation for Session interface
func (m *MockedSession) RefreshRoomProperties(v1 jid.Bare) {
	m.Called(v1)
}

// GetRooms is the implementation for Session interface
func (m *MockedSession) GetRooms(v1 jid.Domain, v2 string) (<-chan *muc.RoomListing, <-chan *muc.ServiceListing, <-chan error) {
	args := m.Called(v1, v2)
	return args.Get(0).(<-chan *muc.RoomListing), args.Get(1).(<-chan *muc.ServiceListing), args.Get(2).(<-chan error)
}

// JoinRoom is the implementation for Session interface
func (m *MockedSession) JoinRoom(v1 jid.Bare, v2 string, v3 string) error {
	args := m.Called(v1, v2, v3)
	return args.Error(0)
}

// CreateInstantRoom is the implementation for session interface
func (m *MockedSession) CreateInstantRoom(v1 jid.Bare) (<-chan bool, <-chan error) {
	args := m.Called(v1)
	return args.Get(0).(<-chan bool), args.Get(1).(<-chan error)
}

// CreateReservedRoom is the implementation for session interface
func (m *MockedSession) CreateReservedRoom(v1 jid.Bare) (<-chan *muc.RoomConfigForm, <-chan error) {
	args := m.Called(v1)
	return args.Get(0).(<-chan *muc.RoomConfigForm), args.Get(1).(<-chan error)
}

// GetRoomConfigurationForm is the implementation for session interface
func (m *MockedSession) GetRoomConfigurationForm(v1 jid.Bare) (<-chan *muc.RoomConfigForm, <-chan error) {
	args := m.Called(v1)
	return args.Get(0).(<-chan *muc.RoomConfigForm), args.Get(1).(<-chan error)
}

// SubmitRoomConfigurationForm is the implementation for session interface
func (m *MockedSession) SubmitRoomConfigurationForm(v1 jid.Bare, v2 *muc.RoomConfigForm) (<-chan bool, <-chan *muc.SubmitFormError) {
	args := m.Called(v1, v2)
	return args.Get(0).(<-chan bool), args.Get(1).(<-chan *muc.SubmitFormError)
}

// CancelRoomConfiguration is the implementation for session interface
func (m *MockedSession) CancelRoomConfiguration(v1 jid.Bare) <-chan error {
	args := m.Called(v1)
	return args.Get(0).(<-chan error)
}

// GetChatServices is the implementation for session interface
func (m *MockedSession) GetChatServices(v1 jid.Domain) (<-chan jid.Domain, <-chan error, func()) {
	args := m.Called(v1)
	return args.Get(0).(<-chan jid.Domain), args.Get(1).(<-chan error), args.Get(2).(func())
}

// DestroyRoom is the implementation for session interface
func (m *MockedSession) DestroyRoom(v1 jid.Bare, v2 string, v3 jid.Bare, v4 string) (<-chan bool, <-chan error) {
	args := m.Called(v1, v2, v3, v4)
	return args.Get(0).(<-chan bool), args.Get(1).(<-chan error)
}

// UpdateOccupantAffiliation is the implementation for session interface
func (m *MockedSession) UpdateOccupantAffiliation(roomID jid.Bare, occupantNickname string, occupantRealJID jid.Full, affiliation mdata.Affiliation, reason string) (<-chan bool, <-chan error) {
	args := m.Called(roomID, occupantNickname, occupantRealJID, affiliation, reason)
	return args.Get(0).(<-chan bool), args.Get(1).(<-chan error)
}

// UpdateOccupantRole is the implementation for session interface
func (m *MockedSession) UpdateOccupantRole(v1 jid.Bare, v2 string, v3 mdata.Role, v4 string) (<-chan bool, <-chan error) {
	args := m.Called(v1, v2, v3, v4)
	return args.Get(0).(<-chan bool), args.Get(1).(<-chan error)
}

// Log is the implementation for session interface
func (m *MockedSession) Log() coylog.Logger {
	return m.Called().Get(0).(coylog.Logger)
}

// LeaveRoom is the implementation for session interface
func (m *MockedSession) LeaveRoom(room jid.Bare, nickname string) (<-chan bool, <-chan error) {
	args := m.Called(room, nickname)
	return args.Get(0).(<-chan bool), args.Get(1).(<-chan error)
}

// GetRoom is the implementation for session interface
func (m *MockedSession) GetRoom(v1 jid.Bare) (*muc.Room, bool) {
	args := m.Called(v1)
	return args.Get(0).(*muc.Room), args.Bool(1)
}

// NewRoom is the implementation for session interface
func (m *MockedSession) NewRoom(v1 jid.Bare) *muc.Room {
	args := m.Called(v1)
	return args.Get(0).(*muc.Room)
}
