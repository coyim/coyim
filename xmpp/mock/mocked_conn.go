package mock

import (
	"encoding/xml"
	"io"
	"reflect"
	"sync"

	"github.com/coyim/coyim/cache"
	"github.com/coyim/coyim/xmpp/data"

	mck "github.com/stretchr/testify/mock"
)

// MockedConn is a mock of the XMPP connection interface
type MockedConn struct {
	mck.Mock
}

// Authenticate implements the XMPP connection interface
func (m *MockedConn) Authenticate(v1, v2 string) error {
	args := m.Called(v1, v2)
	return args.Error(0)
}

// AuthenticationFailure implements the XMPP connection interface
func (m *MockedConn) AuthenticationFailure() error {
	args := m.Called()
	return args.Error(0)
}

// BindResource implements the XMPP connection interface
func (m *MockedConn) BindResource() error {
	args := m.Called()
	return args.Error(0)
}

// Cancel implements the XMPP connection interface
func (m *MockedConn) Cancel(v1 data.Cookie) bool {
	args := m.Called(v1)
	return args.Bool(0)
}

// Cache implements the XMPP connection interface
func (m *MockedConn) Cache() cache.WithExpiry {
	args := m.Called()
	res := args.Get(0)
	if res == nil {
		return nil
	}
	return res.(cache.WithExpiry)
}

// ChangePassword implements the XMPP connection interface
func (m *MockedConn) ChangePassword(v1, v2, v3 string) error {
	args := m.Called(v1, v2, v3)
	return args.Error(0)
}

// Close implements the XMPP connection interface
func (m *MockedConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Config implements the XMPP connection interface
func (m *MockedConn) Config() *data.Config {
	args := m.Called()
	return args.Get(0).(*data.Config)
}

// CustomStorage implements the XMPP connection interface
func (m *MockedConn) CustomStorage() map[xml.Name]reflect.Type {
	args := m.Called()
	return args.Get(0).(map[xml.Name]reflect.Type)
}

// DiscoveryFeatures implements the XMPP connection interface
func (m *MockedConn) DiscoveryFeatures(entity string) ([]string, bool) {
	args := m.Called(entity)
	return args.Get(0).([]string), args.Bool(1)
}

// DiscoveryFeaturesAndIdentities implements the XMPP connection interface
func (m *MockedConn) DiscoveryFeaturesAndIdentities(entity string) ([]data.DiscoveryIdentity, []string, bool) {
	args := m.Called(entity)
	return args.Get(0).([]data.DiscoveryIdentity), args.Get(1).([]string), args.Bool(2)
}

// Features implements the XMPP connection interface
func (m *MockedConn) Features() data.StreamFeatures {
	args := m.Called()
	return args.Get(0).(data.StreamFeatures)
}

// GetRosterDelimiter implements the XMPP connection interface
func (m *MockedConn) GetRosterDelimiter() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// In implements the XMPP connection interface
func (m *MockedConn) In() *xml.Decoder {
	args := m.Called()
	return args.Get(0).(*xml.Decoder)
}

// Lock implements the XMPP connection interface
func (m *MockedConn) Lock() *sync.Mutex {
	args := m.Called()
	return args.Get(0).(*sync.Mutex)
}

// Next implements the XMPP connection interface
func (m *MockedConn) Next() (data.Stanza, error) {
	args := m.Called()
	return args.Get(0).(data.Stanza), args.Error(1)
}

// OriginDomain implements the XMPP connection interface
func (m *MockedConn) OriginDomain() string {
	args := m.Called()
	return args.String(0)
}

// Out implements the XMPP connection interface
func (m *MockedConn) Out() io.Writer {
	args := m.Called()
	return args.Get(0).(io.Writer)
}

// Rand implements the XMPP connection interface
func (m *MockedConn) Rand() io.Reader {
	args := m.Called()
	return args.Get(0).(io.Reader)
}

// RawOut implements the XMPP connection interface
func (m *MockedConn) RawOut() io.WriteCloser {
	args := m.Called()
	return args.Get(0).(io.WriteCloser)
}

// ReadStanzas implements the XMPP connection interface
func (m *MockedConn) ReadStanzas(v1 chan<- data.Stanza) error {
	args := m.Called(v1)
	return args.Error(0)
}

// RegisterAccount implements the XMPP connection interface
func (m *MockedConn) RegisterAccount(v1, v2 string) (bool, error) {
	args := m.Called(v1, v2)
	return args.Bool(0), args.Error(1)
}

// RequestRoster implements the XMPP connection interface
func (m *MockedConn) RequestRoster() (<-chan data.Stanza, data.Cookie, error) {
	args := m.Called()
	return args.Get(0).(<-chan data.Stanza), args.Get(1).(data.Cookie), args.Error(2)
}

// RequestVCard implements the XMPP connection interface
func (m *MockedConn) RequestVCard() (<-chan data.Stanza, data.Cookie, error) {
	args := m.Called()
	return args.Get(0).(<-chan data.Stanza), args.Get(1).(data.Cookie), args.Error(2)
}

// Send implements the XMPP connection interface
func (m *MockedConn) Send(v1, v2 string, v3 bool) error {
	args := m.Called(v1, v2, v3)
	return args.Error(0)
}

// SendMessage implements the XMPP connection interface
func (m *MockedConn) SendMessage(v1 *data.ClientMessage) error {
	args := m.Called(v1)
	return args.Error(0)
}

// SendIQ implements the XMPP connection interface
func (m *MockedConn) SendIQ(v1, v2 string, v3 interface{}) (<-chan data.Stanza, data.Cookie, error) {
	args := m.Called(v1, v2, v3)
	return args.Get(0).(<-chan data.Stanza), args.Get(1).(data.Cookie), args.Error(2)
}

// SendIQReply implements the XMPP connection interface
func (m *MockedConn) SendIQReply(v1, v2, v3 string, v4 interface{}) error {
	args := m.Called(v1, v2, v3, v4)
	return args.Error(0)
}

// SendInitialStreamHeader implements the XMPP connection interface
func (m *MockedConn) SendInitialStreamHeader() error {
	args := m.Called()
	return args.Error(0)
}

// SendPing implements the XMPP connection interface
func (m *MockedConn) SendPing() (reply <-chan data.Stanza, cookie data.Cookie, err error) {
	args := m.Called()
	return args.Get(0).(<-chan data.Stanza), args.Get(1).(data.Cookie), args.Error(2)
}

// SendPresence implements the XMPP connection interface
func (m *MockedConn) SendPresence(v1, v2, v3, v4 string) error {
	args := m.Called(v1, v2, v3, v4)
	return args.Error(0)
}

// SendMUCPresence implements the XMPP connection interface
func (m *MockedConn) SendMUCPresence(v1 string, v2 *data.MUC) (presenceID string, err error) {
	args := m.Called(v1, v2)
	return args.String(0), args.Error(1)
}

// ServerAddress implements the XMPP connection interface
func (m *MockedConn) ServerAddress() string {
	args := m.Called()
	return args.String(0)
}

// SetInOut implements the XMPP connection interface
func (m *MockedConn) SetInOut(v1 *xml.Decoder, v2 io.Writer) {
	m.Called(v1, v2)
}

// SetKeepaliveOut implements the XMPP connection interface
func (m *MockedConn) SetKeepaliveOut(v1 io.Writer) {
	m.Called(v1)
}

// SetRawOut implements the XMPP connection interface
func (m *MockedConn) SetRawOut(v1 io.WriteCloser) {
	m.Called(v1)
}

// SetServerAddress implements the XMPP connection interface
func (m *MockedConn) SetServerAddress(v1 string) {
	m.Called(v1)
}

// SignalPresence implements the XMPP connection interface
func (m *MockedConn) SignalPresence(v1 string) error {
	args := m.Called(v1)
	return args.Error(0)
}

// SetChannelBinding implements the XMPP connection interface
func (m *MockedConn) SetChannelBinding(binding []byte) {
	m.Called(binding)
}

// GetChannelBinding implements the XMPP connection interface
func (m *MockedConn) GetChannelBinding() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

// GetJIDResource implements the XMPP connection interface
func (m *MockedConn) GetJIDResource() string {
	args := m.Called()
	return args.String(0)
}

// SetJIDResource implements the XMPP connection interface
func (m *MockedConn) SetJIDResource(resource string) {
	m.Called(resource)
}

// HasSupportTo implements the XMPP connection interface
func (m *MockedConn) HasSupportTo(entity string, features ...string) bool {
	args := m.Called(entity, features)
	return args.Bool(0)
}

// QueryServiceInformation implements the XMPP connection interface
func (m *MockedConn) QueryServiceInformation(entity string) (*data.DiscoveryInfoQuery, error) {
	args := m.Called(entity)
	return args.Get(0).(*data.DiscoveryInfoQuery), args.Error(1)
}

// QueryServiceItems implements the XMPP connection interface
func (m *MockedConn) QueryServiceItems(entity string) (*data.DiscoveryItemsQuery, error) {
	args := m.Called(entity)
	return args.Get(0).(*data.DiscoveryItemsQuery), args.Error(1)
}

// EntityExists implements the XMPP connection interface
func (m *MockedConn) EntityExists(entity string) (bool, error) {
	args := m.Called(entity)
	return args.Bool(0), args.Error(1)
}

// ServerHasFeature implements the XMPP connection interface
func (m *MockedConn) ServerHasFeature(ns string) bool {
	args := m.Called(ns)
	return args.Bool(0)
}
