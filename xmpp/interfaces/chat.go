package interfaces

import (
	"github.com/coyim/coyim/xmpp/data"
)

type Chat interface {
	CheckForSupport(service string) bool
	QueryRooms(entity string) ([]data.DiscoveryItem, error)
	QueryRoomInformation(room string) (*data.RoomInfo, error)
	EnterRoom(*data.Occupant) error
	LeaveRoom(*data.Occupant) error
	SendChatMessage(msg string, to *data.Room) error
	RequestRoomConfigForm(*data.Room) (*data.Form, error)
}
