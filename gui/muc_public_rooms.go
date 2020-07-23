package gui

import (
	"fmt"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (u *gtkUI) updatedRoomListing(rl *muc.RoomListing, data interface{}) {
	fmt.Printf("We have an updated room listing! - %#v\n", rl)
}

func (u *gtkUI) mucShowPublicRooms() {
	for _, a := range u.accounts {
		if a.session.IsConnected() {
			res, ec := a.session.GetRooms(jid.Parse(a.session.GetConfig().Account).Host())
			go func() {
				for {
					select {
					case rl, ok := <-res:
						if !ok {
							return
						}
						rl.OnUpdate(u.updatedRoomListing, nil)
						fmt.Printf("We have a room listing! - %#v\n", rl)
					case e, ok := <-ec:
						if !ok {
							return
						}
						if e != nil {
							fmt.Printf("Had an error: %v\n", e)
						}
						return
					}
				}
			}()
		}
	}
}
