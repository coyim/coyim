package gui

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session"
	"github.com/twstrike/coyim/xmpp"
)

func (u *gtkUI) handleOneAccountEvent(ev interface{}) {
	switch t := ev.(type) {
	case session.Event:
		doInUIThread(func() {
			u.handleSessionEvent(t)
		})
	case session.PeerEvent:
		doInUIThread(func() {
			u.handlePeerEvent(t)
		})
	case session.NotificationEvent:
		doInUIThread(func() {
			u.handleNotificationEvent(t)
		})
	case session.PresenceEvent:
		doInUIThread(func() {
			u.handlePresenceEvent(t)
		})
	case session.MessageEvent:
		doInUIThread(func() {
			u.handleMessageEvent(t)
		})
	case session.LogEvent:
		doInUIThread(func() {
			u.handleLogEvent(t)
		})
	default:
		log.Printf("unsupported event %#v\n", t)
	}
}

func (u *gtkUI) observeAccountEvents() {
	for ev := range u.events {
		u.handleOneAccountEvent(ev)
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

	if account != nil {
		switch ev.Type {
		case session.Connected:
			account.enableExistingConversationWindows(true)
		case session.Disconnected:
			account.enableExistingConversationWindows(false)
		case session.ConnectionLost:
			u.notifyConnectionFailure(account)
			go u.connectWithRandomDelay(account)
		case session.RosterReceived:
			u.roster.update(account, ev.Session.R())
		}
	}

	u.rosterUpdated()
}

func (u *gtkUI) handlePresenceEvent(ev session.PresenceEvent) {
	if ev.Session.GetConfig().HideStatusUpdates {
		return
	}

	log.Printf("[%s] Presence from %v: show: %v status: %v gone: %v\n", ev.To, ev.From, ev.Show, ev.Status, ev.Gone)
	u.rosterUpdated()

	account := u.findAccountForSession(ev.Session)
	if account == nil {
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

func convWindowNowOrLater(account *account, peer string, f func(*conversationWindow)) {
	convWin, ok := account.getConversationWith(peer)
	if !ok {
		account.afterConversationWindowCreated(peer, f)
	} else {
		f(convWin)
	}
}

func (u *gtkUI) handlePeerEvent(ev session.PeerEvent) {
	identityWarning := func(cw *conversationWindow) {
		cw.updateSecurityWarning()
		cw.showIdentityVerificationWarning(u)
	}

	switch ev.Type {
	case session.IQReceived:
		//TODO
		log.Printf("received iq: %v\n", ev.From)
	case session.OTREnded:
		peer := ev.From
		account := u.findAccountForSession(ev.Session)
		convWindowNowOrLater(account, peer, func(cw *conversationWindow) {
			cw.displayNotification(i18n.Local("Private conversation lost."))
			cw.updateSecurityWarning()
		})

	case session.OTRNewKeys:
		peer := ev.From
		account := u.findAccountForSession(ev.Session)
		convWindowNowOrLater(account, peer, func(cw *conversationWindow) {
			cw.displayNotificationVerifiedOrNot(i18n.Local("Private conversation started."), i18n.Local("Unverified conversation started."))
			identityWarning(cw)
		})

	case session.OTRRenewedKeys:
		peer := ev.From
		account := u.findAccountForSession(ev.Session)
		convWindowNowOrLater(account, peer, func(cw *conversationWindow) {
			cw.displayNotificationVerifiedOrNot(i18n.Local("Successfully refreshed the private conversation."), i18n.Local("Successfully refreshed the unverified private conversation."))
			identityWarning(cw)
		})

	case session.SubscriptionRequest:
		confirmDialog := authorizePresenceSubscriptionDialog(&u.window.Window, ev.From)

		doInUIThread(func() {
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
		})
	case session.Subscribed:
		jid := ev.Session.GetConfig().Account
		log.Printf("[%s] Subscribed to %s\n", jid, ev.From)
		u.rosterUpdated()
	case session.Unsubscribe:
		jid := ev.Session.GetConfig().Account
		log.Printf("[%s] Unsubscribed from %s\n", jid, ev.From)
		u.rosterUpdated()
	}
}

func (u *gtkUI) handleNotificationEvent(ev session.NotificationEvent) {
	peer := ev.Peer
	account := u.findAccountForSession(ev.Session)
	convWin, ok := account.getConversationWith(peer)
	if !ok {
		account.afterConversationWindowCreated(peer, func(cw *conversationWindow) {
			cw.displayNotification(i18n.Local(ev.Notification))
		})
		return
	}

	convWin.displayNotification(ev.Notification)
}
