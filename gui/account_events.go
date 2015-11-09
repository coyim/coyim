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
	for ev := range u.events {
		switch t := ev.(type) {
		case session.Event:
			glib.IdleAdd(func() bool {
				u.handleSessionEvent(t)
				return false
			})
		case session.PeerEvent:
			glib.IdleAdd(func() bool {
				u.handlePeerEvent(t)
				return false
			})
		case session.PresenceEvent:
			glib.IdleAdd(func() bool {
				u.handlePresenceEvent(t)
				return false
			})
		case session.MessageEvent:
			glib.IdleAdd(func() bool {
				u.handleMessageEvent(t)
				return false
			})
		case session.LogEvent:
			glib.IdleAdd(func() bool {
				u.handleLogEvent(t)
				return false
			})
		default:
			log.Printf("unsupported event %#v\n", t)
		}
	}
}

func (u *gtkUI) handleLogEvent(ev session.LogEvent) {
	m := ev.Message

	switch ev.Level {
	case session.Info:
		fmt.Println(">>> INFO", m)
	case session.Warn:
		fmt.Println(">>> WARN", m)
	case session.Alert:
		fmt.Println(">>> ALERT", m)
	}
}

func (u *gtkUI) handleMessageEvent(ev session.MessageEvent) {
	account := u.findAccountForSession(ev.Session)
	if account == nil {
		//TODO error
		return
	}

	u.roster.messageReceived(
		account,
		xmpp.RemoveResourceFromJid(ev.From),
		ev.When,
		ev.Encrypted,
		ev.Body,
	)
}

func (u *gtkUI) handleSessionEvent(ev session.Event) {
	account := u.findAccountForSession(ev.Session)

	switch ev.Type {
	case session.Connected:
	case session.Disconnected:
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

	log.Printf("[%s] Presence from %s: show: %s status: %s gone: %v\n", ev.To, ev.From, ev.Show, ev.Status, ev.Gone)
	u.rosterUpdated()

	account := u.findAccountForSession(ev.Session)
	if account == nil {
		//u.Warn("couldn't find account for " + ev.To)
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
		log.Printf("received iq: %v\n", ev.From)
	case session.OTREnded:
		//TODO
		log.Println("OTR conversation ended with", ev.From)
	case session.OTRNewKeys:
		//TODO
		log.Printf("TODO: notify new keys from %s", ev.From)
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
		log.Printf("[%s] Subscribed to %s\n", jid, ev.From)
		u.rosterUpdated()
	case session.Unsubscribe:
		jid := ev.Session.CurrentAccount.Account
		log.Printf("[%s] Unsubscribed from %s\n", jid, ev.From)
		u.rosterUpdated()
	}
}
