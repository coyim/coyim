package gui

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
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
		case session.OTREnded:
			//TODO
			log.Println("OTR conversation ended with", ev.From)
		case session.OTRNewKeys:
			//TODO
			u.Info(fmt.Sprintf("TODO: notify new keys from %s", ev.From))
		case session.SubscriptionRequest:
			confirmDialog := authorizePresenceSubscriptionDialog(u.window, ev.From)

			glib.IdleAdd(func() bool {
				responseType := gtk.ResponseType(confirmDialog.Run())
				switch responseType {
				case gtk.RESPONSE_YES:
					ev.Session.HandleConfirmOrDeny(ev.From, true)
				case gtk.RESPONSE_NO:
					ev.Session.HandleConfirmOrDeny(ev.From, false)
				default:
					// We got a different response, such as a close of the window. In this case we want
					// to keep the subscription request open
				}
				confirmDialog.Destroy()

				return false
			})
		}
	}
}
