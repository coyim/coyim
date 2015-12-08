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
		u.roster.enableExistingConversationWindows(account, true)
	case session.Disconnected:
		u.roster.enableExistingConversationWindows(account, false)
	case session.ConnectionLost:
		log.Println("Connection lost. Should we reconnect?")

		u.notifyConnectionFailure(account)

		if ev.Session.CurrentAccount.ConnectAutomatically {
			go u.connectWithRandomDelay(account)
			//u.ExecuteCmd(connectAccountCmd(account))
		}
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

	log.Printf("[%s] Presence from %v: show: %v status: %v gone: %v\n", ev.To, ev.From, ev.Show, ev.Status, ev.Gone)
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
		log.Println("New OTR keys from", ev.From)

		peer := ev.From
		conversation := ev.Session.GetConversationWith(peer)

		theirKey := conversation.GetTheirKey()
		if theirKey == nil {
			//Something is VERY wrong
			log.Println("Conversation has no theirKey")
			return
		}

		fingerprint := theirKey.Fingerprint()
		conf := ev.Session.CurrentAccount

		//TODO: this only returns the userID if the fingerprint matches AND is not
		//untrusted. What if this fingerprint is associated with another (untrusted)
		//userID and we trust it for a different userID? Is this a problem?
		userID := conf.UserIDForVerifiedFingerprint(fingerprint)

		switch userID {
		case "":
			//TODO: Unknown fingerprint. User must verify.
		case ev.From:
			//TODO: Already verifyed. Should we notify?
			log.Println("Fingerprint already verified")
			return
		default:
			//TODO: The fingerprint is associated with someone else. Warn!!!
			log.Println("Fingerprint verified with another userID")
			return
		}

		convWin, ok := u.roster.conversations[peer]
		if !ok {
			log.Println("Could not find a conversation window")
			return
		}

		account := u.findAccountForSession(ev.Session)
		infoBar := buildVerifyIdentityNotification(peer)
		infoBar.Connect("response", func(info *gtk.InfoBar, response gtk.ResponseType) {
			if response != gtk.RESPONSE_ACCEPT {
				log.Println("Got response", response)
				return
			}

			glib.IdleAdd(func() {
				verifyFingerprintDialog(account, peer, convWin.win)
			})

			//TODO: should only hide the notification when the identity is verified
			info.Hide()
			info.Destroy()
		})

		log.Println("Notification as added")

		convWin.addNotification(infoBar)
		infoBar.ShowAll()
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
