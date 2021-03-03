package mock

import (
	"encoding/xml"
	"io"
	"reflect"
	"sync"

	"github.com/coyim/coyim/cache"
	"github.com/coyim/coyim/xmpp/data"
)

type Conn struct{}

func (*Conn) Authenticate(string, string) error {
	return nil
}
func (*Conn) AuthenticationFailure() error {
	return nil
}
func (*Conn) BindResource() error {
	return nil
}
func (*Conn) Cancel(data.Cookie) bool {
	return false
}
func (*Conn) Cache() cache.WithExpiry {
	return nil
}
func (*Conn) ChangePassword(string, string, string) error {
	return nil
}
func (*Conn) Close() error {
	return nil
}
func (*Conn) Config() *data.Config {
	return nil
}
func (*Conn) CustomStorage() map[xml.Name]reflect.Type {
	return nil
}
func (*Conn) DiscoveryFeatures(string) ([]string, bool) {
	return nil, false
}
func (*Conn) DiscoveryFeaturesAndIdentities(string) ([]data.DiscoveryIdentity, []string, bool) {
	return nil, nil, false
}
func (*Conn) Features() data.StreamFeatures {
	return data.StreamFeatures{}
}
func (*Conn) GetRosterDelimiter() (string, error) {
	return "", nil
}
func (*Conn) In() *xml.Decoder {
	return nil
}
func (*Conn) Lock() *sync.Mutex {
	return nil
}
func (*Conn) Next() (data.Stanza, error) {
	return data.Stanza{}, nil
}
func (*Conn) OriginDomain() string {
	return ""
}
func (*Conn) Out() io.Writer {
	return nil
}
func (*Conn) Rand() io.Reader {
	return nil
}
func (*Conn) RawOut() io.WriteCloser {
	return nil
}
func (*Conn) ReadStanzas(chan<- data.Stanza) error {
	return nil
}
func (*Conn) RegisterAccount(string, string) (bool, error) {
	return false, nil
}
func (*Conn) RequestRoster() (<-chan data.Stanza, data.Cookie, error) {
	return nil, 0, nil
}
func (*Conn) RequestVCard() (<-chan data.Stanza, data.Cookie, error) {
	return nil, 0, nil
}
func (*Conn) Send(string, string, bool) error {
	return nil
}
func (*Conn) SendMessage(*data.Message) error {
	return nil
}
func (*Conn) SendIQ(string, string, interface{}) (<-chan data.Stanza, data.Cookie, error) {
	return nil, 0, nil
}
func (*Conn) SendIQReply(string, string, string, interface{}) error {
	return nil
}
func (*Conn) SendInitialStreamHeader() error {
	return nil
}
func (*Conn) SendPing() (reply <-chan data.Stanza, cookie data.Cookie, err error) {
	return nil, 0, nil
}
func (*Conn) SendPresence(string, string, string, string) error {
	return nil
}
func (*Conn) SendMUCPresence(string, *data.MUC) error {
	return nil
}
func (*Conn) ServerAddress() string {
	return ""
}
func (*Conn) SetInOut(*xml.Decoder, io.Writer) {}
func (*Conn) SetKeepaliveOut(io.Writer)        {}
func (*Conn) SetRawOut(io.WriteCloser)         {}
func (*Conn) SetServerAddress(string)          {}
func (*Conn) SignalPresence(string) error {
	return nil
}

func (*Conn) SetChannelBinding([]byte) {}
func (*Conn) GetChannelBinding() []byte {
	return nil
}

func (*Conn) GetJIDResource() string {
	return ""
}

func (*Conn) SetJIDResource(string) {}

func (*Conn) HasSupportTo(entity string, features ...string) bool {
	return false
}

func (*Conn) QueryServiceInformation(entity string) (*data.DiscoveryInfoQuery, error) {
	return nil, nil
}
func (*Conn) QueryServiceItems(entity string) (*data.DiscoveryItemsQuery, error) {
	return nil, nil
}
func (*Conn) EntityExists(entity string) (bool, error) {
	return false, nil
}

func (*Conn) ServerHasFeature(ns string) bool {
	return false
}
