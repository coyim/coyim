package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

type mucJoinRoomContext struct {
	a      *account
	v      *mucJoinRoomView
	roomID jid.Bare
	done   func()
}

func (v *mucJoinRoomView) newJoinRoomContext(ca *account, roomID jid.Bare, done func()) *mucJoinRoomContext {
	return &mucJoinRoomContext{
		a:      ca,
		v:      v,
		roomID: roomID,
		done:   done,
	}
}

func (c *mucJoinRoomContext) onFinishWithError(err error, errorReceived bool) {
	if errorReceived {
		c.v.onJoinError(c.a, c.roomID, err)
		return
	}
	c.v.onServiceUnavailable(c.a, c.roomID)
}

func (c *mucJoinRoomContext) waitToFinish(result <-chan bool, errors <-chan error, roomInfo <-chan *muc.RoomListing) {
	defer doInUIThread(c.done)

	select {
	case value, ok := <-result:
		if !ok {
			c.v.onServiceUnavailable(c.a, c.roomID)
			return
		}

		if !value {
			c.v.onJoinFails(c.a, c.roomID)
			return
		}

		ri := <-roomInfo

		c.v.onJoinSuccess(c.a, c.roomID, ri)
	case err, ok := <-errors:
		c.onFinishWithError(err, !ok)
	}
}

func (c *mucJoinRoomContext) joinRoom() {
	c.v.beforeJoiningRoom()
	roomInfo := make(chan *muc.RoomListing)
	result, errors := c.a.session.HasRoom(c.roomID, roomInfo)
	go c.waitToFinish(result, errors, roomInfo)
}
