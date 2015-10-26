package gui

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/session"
	"github.com/twstrike/coyim/xmpp"
)

func (u *gtkUI) observeAccountEvents() {
	//TODO: check for channel close
	for ev := range u.events {
		switch t := ev.(type) {
		case session.Event:
			u.handleSessionEvent(t)
		case session.PeerEvent:
			u.handlePeerEvent(t)
		case session.PresenceEvent:
			u.handlePresenceEvent(t)
		default:
			log.Printf("unsupported event %#v\n", t)
		}
	}
}

func (u *gtkUI) handleSessionEvent(ev session.Event) {
	account := u.findAccountForSession(ev.Session)

	switch ev.Type {
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
	}
}

func (u *gtkUI) handlePresenceEvent(ev session.PresenceEvent) {
	if ev.Session.CurrentAccount.HideStatusUpdates {
		return
	}

	u.Debug(fmt.Sprintf("[%s] Presence from %s: show: %s status: %s gone: %v\n", ev.To, ev.From, ev.Show, ev.Status, ev.Gone))
	u.rosterUpdated()

	account := u.findAccountForSession(ev.Session)
	if account == nil {
		u.Warn("couldn't find account for " + ev.To)
		return
	}

	u.roster.presenceUpdated(
		account,
		xmpp.RemoveResourceFromJid(ev.From),
		ev.Show,
		ev.Status,
		ev.Gone,
	)
}

func (u *gtkUI) handlePeerEvent(ev session.PeerEvent) {
	switch ev.Type {
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
	case session.Subscribed:
		jid := ev.Session.CurrentAccount.Account
		u.Debug(fmt.Sprintf("[%s] Subscribed to %s\n", jid, ev.From))
		u.rosterUpdated()
	case session.Unsubscribe:
		jid := ev.Session.CurrentAccount.Account
		u.Debug(fmt.Sprintf("[%s] Unsubscribed from %s\n", jid, ev.From))
		u.rosterUpdated()
	}
}
