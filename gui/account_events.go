package gui

import (
	"log"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/coyim/session/events"
	"github.com/twstrike/coyim/xmpp/utils"
)

func (u *gtkUI) handleOneAccountEvent(ev interface{}) {
	switch t := ev.(type) {
	case events.Event:
		doInUIThread(func() {
			u.handleSessionEvent(t)
		})
	case events.Peer:
		doInUIThread(func() {
			u.handlePeerEvent(t)
		})
	case events.Notification:
		doInUIThread(func() {
			u.handleNotificationEvent(t)
		})
	case events.Presence:
		doInUIThread(func() {
			u.handlePresenceEvent(t)
		})
	case events.Message:
		doInUIThread(func() {
			u.handleMessageEvent(t)
		})
	case events.Log:
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

func (u *gtkUI) handleLogEvent(ev events.Log) {
	m := ev.Message

	switch ev.Level {
	case events.Info:
		log.Println(">>> INFO", m)
	case events.Warn:
		log.Println(">>> WARN", m)
	case events.Alert:
		log.Println(">>> ALERT", m)
	}
}

func (u *gtkUI) handleMessageEvent(ev events.Message) {
	account := u.findAccountForSession(ev.Session)
	if account == nil {
		//TODO error
		return
	}

	u.roster.messageReceived(
		account,
		ev.From,
		ev.Resource,
		ev.When,
		ev.Encrypted,
		ev.Body,
	)
}

func (u *gtkUI) handleSessionEvent(ev events.Event) {
	account := u.findAccountForSession(ev.Session)

	if account != nil {
		switch ev.Type {
		case events.Connected:
			account.enableExistingConversationWindows(true)
		case events.Disconnected:
			account.enableExistingConversationWindows(false)
		case events.ConnectionLost:
			u.notifyConnectionFailure(account, u.connectionFailureMoreInfoConnectionLost)
			go u.connectWithRandomDelay(account)
		case events.RosterReceived:
			u.roster.update(account, ev.Session.R())
		}
	}

	u.rosterUpdated()
}

func (u *gtkUI) handlePresenceEvent(ev events.Presence) {
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
		utils.RemoveResourceFromJid(ev.From),
		ev.Show,
		ev.Status,
		ev.Gone,
	)
}

func convWindowNowOrLater(account *account, peer string, f func(conversationView)) {
	convWin, ok := account.getConversationWith(peer)
	if !ok {
		account.afterConversationWindowCreated(peer, f)
	} else {
		f(convWin)
	}
}

func (u *gtkUI) handlePeerEvent(ev events.Peer) {
	identityWarning := func(cv conversationView) {
		cv.updateSecurityWarning()
		cv.showIdentityVerificationWarning(u)
	}

	switch ev.Type {
	case events.IQReceived:
		//TODO
		log.Printf("received iq: %v\n", ev.From)
	case events.OTREnded:
		peer := ev.From
		account := u.findAccountForSession(ev.Session)
		convWindowNowOrLater(account, peer, func(cv conversationView) {
			cv.displayNotification(i18n.Local("Private conversation lost."))
			cv.updateSecurityWarning()
		})

	case events.OTRNewKeys:
		peer := ev.From
		account := u.findAccountForSession(ev.Session)
		convWindowNowOrLater(account, peer, func(cv conversationView) {
			cv.displayNotificationVerifiedOrNot(i18n.Local("Private conversation started."), i18n.Local("Unverified conversation started."))
			identityWarning(cv)
		})

	case events.OTRRenewedKeys:
		peer := ev.From
		account := u.findAccountForSession(ev.Session)
		convWindowNowOrLater(account, peer, func(cv conversationView) {
			cv.displayNotificationVerifiedOrNot(i18n.Local("Successfully refreshed the private conversation."), i18n.Local("Successfully refreshed the unverified private conversation."))
			identityWarning(cv)
		})

	case events.SubscriptionRequest:
		confirmDialog := authorizePresenceSubscriptionDialog(u.window, ev.From)

		doInUIThread(func() {
			responseType := gtki.ResponseType(confirmDialog.Run())
			switch responseType {
			case gtki.RESPONSE_YES:
				ev.Session.HandleConfirmOrDeny(ev.From, true)
			case gtki.RESPONSE_NO:
				ev.Session.HandleConfirmOrDeny(ev.From, false)
			default:
				// We got a different response, such as a close of the window. In this case we want
				// to keep the subscription request open
			}
			confirmDialog.Destroy()
		})
	case events.Subscribed:
		jid := ev.Session.GetConfig().Account
		log.Printf("[%s] Subscribed to %s\n", jid, ev.From)
		u.rosterUpdated()
	case events.Unsubscribe:
		jid := ev.Session.GetConfig().Account
		log.Printf("[%s] Unsubscribed from %s\n", jid, ev.From)
		u.rosterUpdated()
	}
}

func (u *gtkUI) handleNotificationEvent(ev events.Notification) {
	peer := ev.Peer
	account := u.findAccountForSession(ev.Session)
	convWin, ok := account.getConversationWith(peer)
	if !ok {
		account.afterConversationWindowCreated(peer, func(cv conversationView) {
			cv.displayNotification(i18n.Local(ev.Notification))
		})
		return
	}

	convWin.displayNotification(ev.Notification)
}
