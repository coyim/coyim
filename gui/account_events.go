package gui

import (
	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/ui"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

func (u *gtkUI) handleOneAccountEvent(ev interface{}, a *account) {
	switch t := ev.(type) {
	case events.Event:
		doInUIThread(func() {
			u.handleSessionEvent(t, a)
		})
	case events.Peer:
		doInUIThread(func() {
			u.handlePeerEvent(t, a)
		})
	case events.Notification:
		doInUIThread(func() {
			u.handleNotificationEvent(t, a)
		})
	case events.DelayedMessageSent:
		doInUIThread(func() {
			u.handleDelayedMessageSentEvent(t, a)
		})
	case events.Presence:
		doInUIThread(func() {
			u.handlePresenceEvent(t, a)
		})
	case events.Message:
		doInUIThread(func() {
			u.handleMessageEvent(t, a)
		})
	case events.FileTransfer:
		doInUIThread(func() {
			u.handleFileTransfer(t, a)
		})
	case events.SMP:
		doInUIThread(func() {
			u.handleSMPEvent(t, a)
		})
	case events.MUC:
		doInUIThread(func() {
			u.handleOneMUCRoomEvent(t, a)
		})
	case events.MUCError:
		doInUIThread(func() {
			u.handleOneMUCErrorEvent(t, a)
		})
	default:
		a.log.WithField("event", t).Warn("unsupported event")
	}
}

func (u *gtkUI) observeAccountEvents(a *account) {
	for ev := range a.events {
		u.handleOneAccountEvent(ev, a)
	}
}

func (u *gtkUI) handleMessageEvent(ev events.Message, a *account) {
	timestamp := ev.When
	encrypted := ev.Encrypted
	message := ev.Body

	p, ok := u.getPeer(a, ev.From.NoResource())
	if ok {
		p.LastSeen(ev.From)
	}

	doInUIThread(func() {
		conv := u.openConversationView(a, ev.From, false)

		sent := sentMessage{
			from:            u.displayNameFor(a, ev.From.NoResource()),
			timestamp:       timestamp,
			isEncrypted:     encrypted,
			isOutgoing:      false,
			strippedMessage: ui.StripSomeHTML(message),
		}
		conv.appendMessage(sent)

		if !conv.isVisible() {
			u.maybeNotify(timestamp, a, ev.From.NoResource(), string(ui.StripSomeHTML(message)))
		}
	})
}

func (u *gtkUI) handleSessionEvent(ev events.Event, a *account) {
	switch ev.Type {
	case events.Connected:
		a.enableExistingConversationWindows(true)
		u.notifyChangeOfConnectedAccounts()
	case events.Disconnected:
		a.enableExistingConversationWindows(false)
		u.notifyChangeOfConnectedAccounts()
	case events.ConnectionLost:
		u.notifyConnectionFailure(a, u.connectionFailureMoreInfoConnectionLost)
		u.notifyChangeOfConnectedAccounts()
		go u.connectWithRandomDelay(a)
	case events.RosterReceived:
		u.roster.update(a, a.session.R())
	}

	u.rosterUpdated()
}

func (u *gtkUI) handlePresenceEvent(ev events.Presence, a *account) {
	peer := jid.R(ev.From)
	// if p, ok := u.getPeer(account, peer.NoResource()); ok {
	// 	p.Presence(peer, ev.Show, ev.Status, ev.Gone)
	// }

	if !a.session.GetConfig().HideStatusUpdates {
		u.presenceUpdated(a, peer, ev)
	}

	a.log.WithFields(log.Fields{
		"from":   peer,
		"show":   ev.Show,
		"status": ev.Status,
		"gone":   ev.Gone,
	}).Info("Presence received")
	u.rosterUpdated()
}

func convWindowNowOrLater(account *account, peer jid.Any, ui *gtkUI, f func(conversationView)) {
	ui.NewConversationViewFactory(account, peer, false).IfConversationView(f, func() {
		account.afterConversationWindowCreated(peer, f)
	})
}

func (u *gtkUI) handlePeerEvent(ev events.Peer, a *account) {
	switch ev.Type {
	case events.IQReceived:
		//TODO
		// TODO WHAT?
		a.log.WithFields(log.Fields{"from": ev.From, "event": ev}).Info("received iq")
	case events.OTREnded:
		convWindowNowOrLater(a, ev.From, u, func(cv conversationView) {
			cv.updateSecurityStatus()

			cv.removeOtrLock()
			cv.displayNotification(i18n.Local("Private conversation has ended."))
			cv.haveShownPrivateEndedNotification()
		})

	case events.OTRNewKeys:
		convWindowNowOrLater(a, ev.From, u, func(cv conversationView) {
			cv.setOtrLock(ev.From.(jid.WithResource))
			cv.calculateNewKeyStatus()
			cv.savePeerFingerprint(u)
			cv.updateSecurityStatus()

			cv.displayNotificationVerifiedOrNot(i18n.Local("Private conversation started."), i18n.Local("Private conversation started (tagged: '%s')."), i18n.Local("Unverified conversation started."))
			cv.appendPendingDelayed()
			cv.haveShownPrivateNotification()
		})

	case events.OTRRenewedKeys:
		convWindowNowOrLater(a, ev.From, u, func(cv conversationView) {
			cv.updateSecurityStatus()

			cv.displayNotificationVerifiedOrNot(i18n.Local("Successfully refreshed the private conversation."), i18n.Local("Successfully refreshed the private conversation (tagged: '%s')."), i18n.Local("Successfully refreshed the unverified private conversation."))
		})

	case events.SubscriptionRequest:
		authorizePresenceSubscriptionDialog(u.window, ev.From.NoResource(), func(r gtki.ResponseType) {
			switch r {
			case gtki.RESPONSE_YES:
				a.session.HandleConfirmOrDeny(ev.From.NoResource(), true)
			case gtki.RESPONSE_NO:
				a.session.HandleConfirmOrDeny(ev.From.NoResource(), false)
			default:
				// We got a different response, such as a close of the window. In this case we want
				// to keep the subscription request open
			}
		})
	case events.Subscribed:
		a.log.WithField("to", ev.From).Info("Subscribed to peer")
		u.rosterUpdated()
	case events.Unsubscribe:
		a.log.WithField("from", ev.From).Info("Unsubscribed from peer")
		u.rosterUpdated()
	}
}

func (u *gtkUI) handleNotificationEvent(ev events.Notification, a *account) {
	convWin := u.openConversationView(a, ev.Peer, false)
	convWin.displayNotification(i18n.Local(ev.Notification))
}

func (u *gtkUI) handleDelayedMessageSentEvent(ev events.DelayedMessageSent, a *account) {
	convWin := u.openConversationView(a, ev.Peer, false)
	convWin.delayedMessageSent(ev.Tracer)
}

func (u *gtkUI) handleSMPEvent(ev events.SMP, a *account) {
	convWin := u.openConversationView(a, ev.From, false)

	switch ev.Type {
	case events.SecretNeeded:
		convWin.showSMPRequestForSecret(ev.Body)
	case events.Success:
		convWin.showSMPSuccess()
	case events.Failure:
		convWin.showSMPFailure()
	}
}
