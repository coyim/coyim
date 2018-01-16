package gui

import (
	"log"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/ui"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
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
	case events.DelayedMessageSent:
		doInUIThread(func() {
			u.handleDelayedMessageSentEvent(t)
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
	case events.FileTransfer:
		doInUIThread(func() {
			u.handleFileTransfer(t)
		})
	case events.SMP:
		doInUIThread(func() {
			u.handleSMPEvent(t)
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

	from, resource := ev.From.PotentialSplit()
	timestamp := ev.When
	encrypted := ev.Encrypted
	message := ev.Body

	p, ok := u.getPeer(account, from)
	if ok {
		p.LastResource(resource)
	}

	doInUIThread(func() {
		//TODO: here we dont want to open, only findOrCreate
		//this is what the false means
		conv := u.openConversationView(account, from, false, "")

		sent := sentMessage{
			from:            u.displayNameFor(account, from),
			timestamp:       timestamp,
			isEncrypted:     encrypted,
			isOutgoing:      false,
			strippedMessage: ui.StripSomeHTML(message),
		}
		conv.appendMessage(sent)

		if !conv.isVisible() {
			u.maybeNotify(timestamp, account, from, string(ui.StripSomeHTML(message)))
		}
	})

}

func (u *gtkUI) handleSessionEvent(ev events.Event) {
	account := u.findAccountForSession(ev.Session)
	if account == nil {
		return
	}

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

	u.rosterUpdated()
}

func (u *gtkUI) handlePresenceEvent(ev events.Presence) {
	if !ev.Session.GetConfig().HideStatusUpdates {
		u.presenceUpdated(ev)
	}

	log.Printf("[%s] Presence from %v: show: %v status: %v gone: %v\n", ev.To, ev.From, ev.Show, ev.Status, ev.Gone)
	u.rosterUpdated()
}

func convWindowNowOrLater(account *account, peer jid.Any, ui *gtkUI, f func(conversationView)) {
	// TODO: this is clearly wrong
	fullJID := peer.NoResource().String()
	convWin, ok := ui.getConversationView(account, fullJID)
	if !ok {
		account.afterConversationWindowCreated(peer.String(), f)
	} else {
		f(convWin)
	}
}

func (u *gtkUI) handlePeerEvent(ev events.Peer) {
	identityWarning := func(cv conversationView) {
		cv.updateSecurityWarning()
		cv.removeIdentityVerificationWarning()
		cv.showIdentityVerificationWarning(u)
	}

	switch ev.Type {
	case events.IQReceived:
		//TODO
		log.Printf("received iq: %v\n", ev.From)
	case events.OTREnded:
		account := u.findAccountForSession(ev.Session)
		convWindowNowOrLater(account, ev.From, u, func(cv conversationView) {
			cv.displayNotification(i18n.Local("Private conversation has ended."))
			cv.updateSecurityWarning()
			cv.removeIdentityVerificationWarning()
			cv.haveShownPrivateEndedNotification()
		})

	case events.OTRNewKeys:
		account := u.findAccountForSession(ev.Session)
		convWindowNowOrLater(account, ev.From, u, func(cv conversationView) {
			cv.displayNotificationVerifiedOrNot(u, i18n.Local("Private conversation started."), i18n.Local("Unverified conversation started."))
			cv.appendPendingDelayed()
			identityWarning(cv)
			cv.haveShownPrivateNotification()
		})

	case events.OTRRenewedKeys:
		account := u.findAccountForSession(ev.Session)
		convWindowNowOrLater(account, ev.From, u, func(cv conversationView) {
			cv.displayNotificationVerifiedOrNot(u, i18n.Local("Successfully refreshed the private conversation."), i18n.Local("Successfully refreshed the unverified private conversation."))
			identityWarning(cv)
		})

	case events.SubscriptionRequest:
		confirmDialog := authorizePresenceSubscriptionDialog(u.window, ev.From.String())

		doInUIThread(func() {
			responseType := gtki.ResponseType(confirmDialog.Run())
			switch responseType {
			case gtki.RESPONSE_YES:
				ev.Session.HandleConfirmOrDeny(ev.From.NoResource(), true)
			case gtki.RESPONSE_NO:
				ev.Session.HandleConfirmOrDeny(ev.From.NoResource(), false)
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
	account := u.findAccountForSession(ev.Session)
	convWin := u.openConversationView(account, ev.Peer.NoResource(), false, jid.Resource(""))

	convWin.displayNotification(i18n.Local(ev.Notification))
}

func (u *gtkUI) handleDelayedMessageSentEvent(ev events.DelayedMessageSent) {
	account := u.findAccountForSession(ev.Session)
	convWin := u.openConversationView(account, ev.Peer.NoResource(), false, jid.Resource(""))

	convWin.delayedMessageSent(ev.Tracer)
}

func (u *gtkUI) handleSMPEvent(ev events.SMP) {
	account := u.findAccountForSession(ev.Session)
	convWin := u.openConversationView(account, ev.From.NoResource(), false, jid.Resource(""))

	switch ev.Type {
	case events.SecretNeeded:
		convWin.showSMPRequestForSecret(ev.Body)
	case events.Success:
		convWin.showSMPSuccess()
	case events.Failure:
		convWin.showSMPFailure()
	}
}
