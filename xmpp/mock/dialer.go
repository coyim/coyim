package mock

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/servers"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
	"golang.org/x/net/proxy"
)

// Dialer is a mock of the XMPP dialer interface
type Dialer struct{}

// Config is an implementation of the Dialer interface
func (*Dialer) Config() data.Config {
	return data.Config{}
}

// Dial is an implementation of the Dialer interface
func (*Dialer) Dial() (interfaces.Conn, error) {
	return nil, nil
}

// GetServer is an implementation of the Dialer interface
func (*Dialer) GetServer() string {
	return ""
}

// RegisterAccount is an implementation of the Dialer interface
func (*Dialer) RegisterAccount(data.FormCallback) (interfaces.Conn, error) {
	return nil, nil
}

// ServerAddress is an implementation of the Dialer interface
func (*Dialer) ServerAddress() string {
	return ""
}

// SetConfig is an implementation of the Dialer interface
func (*Dialer) SetConfig(data.Config) {}

// SetJID is an implementation of the Dialer interface
func (*Dialer) SetJID(string) {}

// SetPassword is an implementation of the Dialer interface
func (*Dialer) SetPassword(string) {}

// SetProxy is an implementation of the Dialer interface
func (*Dialer) SetProxy(proxy.Dialer) {}

// SetResource is an implementation of the Dialer interface
func (*Dialer) SetResource(string) {}

// SetServerAddress is an implementation of the Dialer interface
func (*Dialer) SetServerAddress(string) {}

// SetShouldConnectTLS is an implementation of the Dialer interface
func (*Dialer) SetShouldConnectTLS(bool) {}

// SetShouldSendALPN is an implementation of the Dialer interface
func (*Dialer) SetShouldSendALPN(bool) {}

// SetLogger is an implementation of the Dialer interface
func (*Dialer) SetLogger(coylog.Logger) {}

// SetKnown is an implementation of the Dialer interface
func (*Dialer) SetKnown(*servers.Server) {}
