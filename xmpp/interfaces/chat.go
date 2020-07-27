package interfaces

import (
	"github.com/coyim/coyim/xmpp/data"
)

// LegacyOldDoNotUseChat contains the basic functionality to do a chat in a room
type LegacyOldDoNotUseChat interface {
	Events() chan interface{}
	CheckForSupport(service string) bool
	QueryRooms(entity string) ([]data.DiscoveryItem, error)
	LegacyOldDoNotUseQueryRoomInformation(room string) (*data.LegacyOldDoNotUseRoomInfo, error)
	LegacyOldDoNotUseCreateRoom(*data.LegacyOldDoNotUseRoom) error
	LegacyOldDoNotUseEnterRoom(*data.LegacyOldDoNotUseOccupant) error
	LegacyOldDoNotUseLeaveRoom(*data.LegacyOldDoNotUseOccupant) error
	LegacyOldDoNotUseSendChatMessage(msg string, to *data.LegacyOldDoNotUseRoom) error
	LegacyOldDoNotUseRoomConfigForm(*data.LegacyOldDoNotUseRoom, data.FormCallback) error

	//TODO: Remove me
	LegacyOldDoNotUseRequestRoomConfigForm(*data.LegacyOldDoNotUseRoom) (*data.Form, error)
	LegacyOldDoNotUseUpdateRoomConfig(*data.LegacyOldDoNotUseRoom, *data.Form) error
}
