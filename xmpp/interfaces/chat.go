package interfaces

import "github.com/coyim/coyim/xmpp/data"

type Chat interface {
	CheckForSupport(service string) bool
	QueryRoomInformation(room string) (*data.DiscoveryInfoQuery, error)
}
