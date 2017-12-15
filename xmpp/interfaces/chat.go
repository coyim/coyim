package interfaces

import (
	"github.com/coyim/coyim/xmpp/data"
)

type Chat interface {
	Events() chan interface{}
	CheckForSupport(service string) bool
	QueryRooms(entity string) ([]data.DiscoveryItem, error)
	QueryRoomInformation(room string) (*data.RoomInfo, error)
	EnterRoom(*data.Occupant) error
	LeaveRoom(*data.Occupant) error
	SendChatMessage(msg string, to *data.Room) error
	RoomConfigForm(*data.Room, data.FormCallback) error

	//TODO: Remove me
	RequestRoomConfigForm(*data.Room) (*data.Form, error)
	UpdateRoomConfig(*data.Room, *data.Form) error
}
