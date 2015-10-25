package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/twstrike/coyim/session"
)

func (u *gtkUI) observeAccountEvents() {
	//TODO: check for channel close
	for ev := range u.events {
		account := u.findAccountForSession(ev.Session)

		switch ev.EventType {
		case session.Connected:
			if account == nil {
				return
			}

			u.window.Emit(account.connectedSignal.String())
		case session.Disconnected:
			if account != nil {
				u.window.Emit(account.disconnectedSignal.String())
			}

			for _, acc := range u.accounts {
				if acc.session.ConnStatus == session.CONNECTED {
					return
				}
			}

			u.roster.disconnected()
		case session.RosterReceived:
			if account == nil {
				return
			}

			u.roster.update(account, ev.Session.R)

			glib.IdleAdd(func() bool {
				u.roster.redraw()
				return false
			})
		case session.IQReceived:
			//TODO
			u.Debug(fmt.Sprintf("received iq: %v\n", ev.From))
		}
	}
}
