package mock

import (
	"encoding/xml"
	"io"
	"reflect"
	"sync"

	"github.com/coyim/coyim/cache"
	"github.com/coyim/coyim/xmpp/data"
)

// Conn is a mock of the XMPP connection interface
type Conn struct{}

// Authenticate implements the XMPP connection interface
func (*Conn) Authenticate(string, string) error {
	return nil
}

// AuthenticationFailure implements the XMPP connection interface
func (*Conn) AuthenticationFailure() error {
	return nil
}

// BindResource implements the XMPP connection interface
func (*Conn) BindResource() error {
	return nil
}

// Cancel implements the XMPP connection interface
func (*Conn) Cancel(data.Cookie) bool {
	return false
}

// Cache implements the XMPP connection interface
func (*Conn) Cache() cache.WithExpiry {
	return nil
}

// ChangePassword implements the XMPP connection interface
func (*Conn) ChangePassword(string, string, string) error {
	return nil
}

// Close implements the XMPP connection interface
func (*Conn) Close() error {
	return nil
}

// Config implements the XMPP connection interface
func (*Conn) Config() *data.Config {
	return nil
}

// CustomStorage implements the XMPP connection interface
func (*Conn) CustomStorage() map[xml.Name]reflect.Type {
	return nil
}

// DiscoveryFeatures implements the XMPP connection interface
func (*Conn) DiscoveryFeatures(string) ([]string, bool) {
	return nil, false
}

// DiscoveryFeaturesAndIdentities implements the XMPP connection interface
func (*Conn) DiscoveryFeaturesAndIdentities(string) ([]data.DiscoveryIdentity, []string, bool) {
	return nil, nil, false
}

// Features implements the XMPP connection interface
func (*Conn) Features() data.StreamFeatures {
	return data.StreamFeatures{}
}

// GetRosterDelimiter implements the XMPP connection interface
func (*Conn) GetRosterDelimiter() (string, error) {
	return "", nil
}

// In implements the XMPP connection interface
func (*Conn) In() *xml.Decoder {
	return nil
}

// Lock implements the XMPP connection interface
func (*Conn) Lock() *sync.Mutex {
	return nil
}

// Next implements the XMPP connection interface
func (*Conn) Next() (data.Stanza, error) {
	return data.Stanza{}, nil
}

// OriginDomain implements the XMPP connection interface
func (*Conn) OriginDomain() string {
	return ""
}

// Out implements the XMPP connection interface
func (*Conn) Out() io.Writer {
	return nil
}

// Rand implements the XMPP connection interface
func (*Conn) Rand() io.Reader {
	return nil
}

// RawOut implements the XMPP connection interface
func (*Conn) RawOut() io.WriteCloser {
	return nil
}

// ReadStanzas implements the XMPP connection interface
func (*Conn) ReadStanzas(chan<- data.Stanza) error {
	return nil
}

// RegisterAccount implements the XMPP connection interface
func (*Conn) RegisterAccount(string, string) (bool, error) {
	return false, nil
}

// RequestRoster implements the XMPP connection interface
func (*Conn) RequestRoster() (<-chan data.Stanza, data.Cookie, error) {
	return nil, 0, nil
}

// RequestVCard implements the XMPP connection interface
func (*Conn) RequestVCard() (<-chan data.Stanza, data.Cookie, error) {
	return nil, 0, nil
}

// Send implements the XMPP connection interface
func (*Conn) Send(string, string, bool) error {
	return nil
}

// SendMessage implements the XMPP connection interface
func (*Conn) SendMessage(*data.Message) error {
	return nil
}

// SendIQ implements the XMPP connection interface
func (*Conn) SendIQ(string, string, interface{}) (<-chan data.Stanza, data.Cookie, error) {
	return nil, 0, nil
}

// SendIQReply implements the XMPP connection interface
func (*Conn) SendIQReply(string, string, string, interface{}) error {
	return nil
}

// SendInitialStreamHeader implements the XMPP connection interface
func (*Conn) SendInitialStreamHeader() error {
	return nil
}

// SendPing implements the XMPP connection interface
func (*Conn) SendPing() (reply <-chan data.Stanza, cookie data.Cookie, err error) {
	return nil, 0, nil
}

// SendPresence implements the XMPP connection interface
func (*Conn) SendPresence(string, string, string, string) error {
	return nil
}

// SendMUCPresence implements the XMPP connection interface
func (*Conn) SendMUCPresence(string, *data.MUC) error {
	return nil
}

// ServerAddress implements the XMPP connection interface
func (*Conn) ServerAddress() string {
	return ""
}

// SetInOut implements the XMPP connection interface
func (*Conn) SetInOut(*xml.Decoder, io.Writer) {}

// SetKeepaliveOut implements the XMPP connection interface
func (*Conn) SetKeepaliveOut(io.Writer) {}

// SetRawOut implements the XMPP connection interface
func (*Conn) SetRawOut(io.WriteCloser) {}

// SetServerAddress implements the XMPP connection interface
func (*Conn) SetServerAddress(string) {}

// SignalPresence implements the XMPP connection interface
func (*Conn) SignalPresence(string) error {
	return nil
}

// SetChannelBinding implements the XMPP connection interface
func (*Conn) SetChannelBinding([]byte) {}

// GetChannelBinding implements the XMPP connection interface
func (*Conn) GetChannelBinding() []byte {
	return nil
}

// GetJIDResource implements the XMPP connection interface
func (*Conn) GetJIDResource() string {
	return ""
}

// SetJIDResource implements the XMPP connection interface
func (*Conn) SetJIDResource(string) {}

// HasSupportTo implements the XMPP connection interface
func (*Conn) HasSupportTo(entity string, features ...string) bool {
	return false
}

// QueryServiceInformation implements the XMPP connection interface
func (*Conn) QueryServiceInformation(entity string) (*data.DiscoveryInfoQuery, error) {
	return nil, nil
}

// QueryServiceItems implements the XMPP connection interface
func (*Conn) QueryServiceItems(entity string) (*data.DiscoveryItemsQuery, error) {
	return nil, nil
}

// EntityExists implements the XMPP connection interface
func (*Conn) EntityExists(entity string) (bool, error) {
	return false, nil
}

// ServerHasFeature implements the XMPP connection interface
func (*Conn) ServerHasFeature(ns string) bool {
	return false
}
