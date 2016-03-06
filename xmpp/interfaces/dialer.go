package interfaces

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/golang.org/x/net/proxy"
	"github.com/twstrike/coyim/xmpp/data"
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
