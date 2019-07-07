package interfaces

import (
	"github.com/coyim/coyim/tls"
	"github.com/coyim/coyim/xmpp/data"
	"golang.org/x/net/proxy"
)

// Dialer connects and authenticates to an XMPP server
type Dialer interface {
	Config() data.Config
	Dial() (Conn, error)
	GetServer() string
	RegisterAccount(data.FormCallback) (Conn, error)
	ServerAddress() string
	SetConfig(data.Config)
	SetJID(string)
	SetPassword(string)
	SetProxy(proxy.Dialer)
	SetServerAddress(v string)
}

// DialerFactory represents a function that can create a Dialer
type DialerFactory func(tls.Verifier, tls.Factory) Dialer
