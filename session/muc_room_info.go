package session

import (
	"time"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

const maxTimeForRoomDiscoInfoRequest = time.Second * 10

func (m *mucManager) requestRoomDiscoInfo(roomID jid.Bare) {
	result := make(chan *muc.RoomListing)
	go m.getRoomListing(roomID, result)

	select {
	case rl := <-result:
		m.onRoomDiscoInfoReceived(roomID, rl)
	case <-time.After(maxTimeForRoomDiscoInfoRequest):
		m.roomDiscoInfoRequestTimeout(roomID)
	}
}

func (m *mucManager) onRoomDiscoInfoReceived(roomID jid.Bare, rl *muc.RoomListing) {
	m.addRoomInfo(roomID, rl)
	m.roomDiscoInfoReceived(roomID, rl.GetDiscoInfo())
}

func (m *mucManager) discoInfoForRoom(roomID jid.Bare) *muc.RoomListing {
	rl, ok := m.getRoomInfo(roomID)
	if !ok {
		rl = m.newRoomListing(roomID)
		m.addRoomInfo(roomID, rl)
	}
	return rl
}
