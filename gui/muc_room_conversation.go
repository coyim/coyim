package gui

import (
	"time"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

const bottomPositionValue = 3

type roomViewConversation struct {
	u                        *gtkUI
	tags                     *conversationTags
	roomID                   jid.Bare
	account                  *account
	roomView                 *roomView
	canSendMessages          bool
	selfOccupantNickname     func() string
	saveNotificationMessage  func(*data.DelayedMessage)
	selfOccupantJoined       chan bool
	historyPrinted           chan bool
	disconnectMessagePrinted bool
	currentAdjustment        float64
	maxAdjustment            float64

	view                  gtki.Box            `gtk-widget:"room-conversation"`
	chatScrolledWindow    gtki.ScrolledWindow `gtk-widget:"chat-scrolled-window"`
	chatTextView          gtki.TextView       `gtk-widget:"chat-text-view"`
	messageView           gtki.Box            `gtk-widget:"message-view"`
	messageScrolledWindow gtki.ScrolledWindow `gtk-widget:"message-scrolled-window"`
	messageTextView       gtki.TextView       `gtk-widget:"message-text-view"`
	sendButton            gtki.Button         `gtk-widget:"message-send-button"`
	sendButtonIcon        gtki.Image          `gtk-widget:"message-send-icon"`
	notificationBox       gtki.Box            `gtk-widget:"notification-box"`

	messageBoxNotification *roomMessageBoxNotification

	log coylog.Logger
}

func (v *roomView) newRoomViewConversation() *roomViewConversation {
	c := &roomViewConversation{
		u:                       v.u,
		roomID:                  v.room.ID,
		account:                 v.account,
		roomView:                v,
		selfOccupantNickname:    v.room.SelfOccupantNickname,
		saveNotificationMessage: v.room.AddMessage,
		selfOccupantJoined:      make(chan bool),
		historyPrinted:          make(chan bool),
	}

	c.log = c.account.log.WithFields(log.Fields{
		"who":  "roomViewConversation",
		"room": c.roomID,
	})

	c.initBuilder()
	c.initDefaults(v)
	c.initSubscribers(v)
	c.initTagsAndTextBuffers()

	return c
}

func (c *roomViewConversation) initBuilder() {
	builder := newBuilder("MUCRoomConversation")
	panicOnDevError(builder.bindObjects(c))

	builder.ConnectSignals(map[string]interface{}{
		"on_send_message": c.onSendMessage,
		"on_key_press":    c.onKeyPress,
		"on_edge_reached": c.onEdgeReached,
	})

	adj := c.chatScrolledWindow.GetVAdjustment()
	adj.Connect("changed", c.onAdjustmentChanged)
	adj.Connect("value-changed", c.updateCurrentAdjustmentValue)

	mucStyles.setScrolledWindowStyle(c.chatScrolledWindow)
	mucStyles.setScrolledWindowStyle(c.messageScrolledWindow)
	mucStyles.setMessageViewBoxStyle(c.messageView)
}

func (c *roomViewConversation) initDefaults(v *roomView) {
	c.sendButtonIcon.SetFromPixbuf(getMUCIconPixbuf("send"))

	c.messageBoxNotification = newRoomMessageBoxNotification()
	c.notificationBox.Add(c.messageBoxNotification.infoBar())

	c.enableSendCapabilitiesIfHasVoice(v.room.SelfOccupant().HasVoice())
}

func (c *roomViewConversation) initSubscribers(v *roomView) {
	v.subscribe("conversation", func(ev roomViewEvent) {
		switch t := ev.(type) {
		case occupantLeftEvent:
			c.occupantLeftEvent(t.nickname)
		case occupantJoinedEvent:
			c.occupantJoinedEvent(t.nickname)
		case selfOccupantReconnectedEvent:
			c.selfOccupantJoinedEvent(t.nickname, t.role)
		case selfOccupantJoinedEvent:
			c.selfOccupantJoinedEvent(t.nickname, t.role)
		case occupantUpdatedEvent:
			c.occupantUpdatedEvent(t.nickname, t.role)
		case messageEvent:
			c.messageEvent(t.tp, t.nickname, t.message, t.timestamp)
		case discussionHistoryEvent:
			c.discussionHistoryEvent(t.history)
		case messageForbidden:
			c.messageForbiddenEvent()
		case messageNotAcceptable:
			c.messageNotAcceptableEvent()
		case subjectUpdatedEvent:
			c.subjectUpdatedEvent(t.nickname, t.subject)
		case subjectReceivedEvent:
			c.subjectReceivedEvent(t.subject)
		case loggingEnabledEvent:
			c.loggingEnabledEvent()
		case loggingDisabledEvent:
			c.loggingDisabledEvent()
		case roomAnonymityEvent:
			c.roomAnonymityChangedEvent(t.anonymityLevel)
		case roomConfigChangedEvent:
			c.roomConfigChangedEvent(t.changes, t.discoInfo)
		case selfOccupantRemovedEvent:
			c.selfOccupantRemovedEvent(v.room.SelfOccupantNickname())
		case occupantRemovedEvent:
			c.occupantRemovedEvent(t.nickname)
		case roomDestroyedEvent:
			c.roomDestroyedEvent(t.reason, t.alternative, t.password)
		case occupantAffiliationRoleUpdatedEvent:
			c.occupantAffiliationRoleUpdatedEvent(t.affiliationRoleUpdate)
		case selfOccupantAffiliationRoleUpdatedEvent:
			c.selfOccupantAffiliationRoleUpdatedEvent(t.selfAffiliationRoleUpdate)
		case occupantAffiliationUpdatedEvent:
			c.occupantAffiliationEvent(t.affiliationUpdate)
		case selfOccupantAffiliationUpdatedEvent:
			c.selfOccupantAffiliationEvent(t.selfAffiliationUpdate.AffiliationUpdate)
		case occupantRoleUpdatedEvent:
			c.occupantRoleEvent(t.roleUpdate)
		case selfOccupantRoleUpdatedEvent:
			c.selfOccupantRoleEvent(t.selfRoleUpdate)
		case selfOccupantDisconnectedEvent:
			c.selfOccupantDisconnectedEvent()
		case accountAffiliationUpdated:
			c.occupantModifiedEvent(t.accountAddress, t.affiliation)
		case roomDisableEvent:
			c.roomDisableEvent()
		case roomEnableEvent:
			c.roomEnableEvent()
		case reopenRoomEvent:
			c.reopenRoomWindowEvent(t.history)
		case selfOccupantRemovedOnAffiliationChangeEvent:
			c.selfOccupantRemovedOnAffiliationChangeEvent(v.room.SelfOccupantNickname())
		case occupantRemovedOnAffiliationChangeEvent:
			c.occupantRemovedOnAffiliationChangeEvent(t.nickname)
		}
	})
}

func (c *roomViewConversation) roomDestroyedEvent(reason string, alternative jid.Bare, password string) {
	doInUIThread(func() {
		c.updateNotificationMessage(messageForRoomDestroyedEvent())
		c.displayNotificationWhenRoomDestroyed(reason, alternative, password)
		c.disableSendCapabilities()
	})
}

func (c *roomViewConversation) occupantAffiliationRoleUpdatedEvent(aru data.AffiliationRoleUpdate) {
	doInUIThread(func() {
		c.displayOccupantUpdateMessageFor(aru)
	})
}

func (c *roomViewConversation) selfOccupantAffiliationRoleUpdatedEvent(saru data.SelfAffiliationRoleUpdate) {
	doInUIThread(func() {
		c.displayOccupantUpdateMessageFor(saru.AffiliationRoleUpdate)
	})
}

func (c *roomViewConversation) occupantAffiliationEvent(au data.AffiliationUpdate) {
	doInUIThread(func() {
		c.displayOccupantUpdateMessageFor(au)
	})
}

func (c *roomViewConversation) occupantRoleEvent(ru data.RoleUpdate) {
	doInUIThread(func() {
		c.displayOccupantUpdateMessageFor(ru)
	})
}

func (c *roomViewConversation) selfOccupantAffiliationEvent(au data.AffiliationUpdate) {
	c.occupantAffiliationEvent(au)

	if au.New.IsBanned() {
		doInUIThread(c.onSelfOccupantBanned)
	}
}

func (c *roomViewConversation) selfOccupantRoleEvent(sru data.SelfRoleUpdate) {
	c.occupantRoleEvent(sru.RoleUpdate)

	switch {
	case sru.New.IsNone():
		doInUIThread(c.onSelfOccupantKicked)
	case sru.New.IsVisitor():
		doInUIThread(c.onSelfOccupantVoiceRevoked)
	}
}

func (c *roomViewConversation) selfOccupantDisconnectedEvent() {
	doInUIThread(func() {
		c.updateNotificationMessage(messageForSelfOccupantDisconnected())
		c.disableSendCapabilities()

		if !c.disconnectMessagePrinted {
			c.occupantLeftEvent(c.selfOccupantNickname())
			c.disconnectMessagePrinted = true
		}
	})
}

func (c *roomViewConversation) occupantModifiedEvent(accountAddress jid.Any, affiliation data.Affiliation) {
	doInUIThread(func() {
		message := ""

		switch {
		case affiliation.IsOwner():
			message = i18n.Localf("$nickname{%s} has been added as $affiliation{an owner}.", accountAddress.String())
		case affiliation.IsAdmin():
			message = i18n.Localf("$nickname{%s} has been added as $affiliation{an administrator}.", accountAddress.String())
		case affiliation.IsMember():
			message = i18n.Localf("$nickname{%s} has been added as $affiliation{a member}.", accountAddress.String())
		case affiliation.IsBanned():
			message = i18n.Localf("$nickname{%s} has been added to the ban list.", accountAddress.String())
		}

		if message != "" {
			c.saveAndDisplayMessage(accountAddress.String(), message, time.Now(), data.OccupantInformationChanged)
		}
	})
}

// roomDisableEvent MUST be called from the UI thread
func (c *roomViewConversation) roomDisableEvent() {
	addRoomDisableClass(c.chatTextView)
}

// roomEnableEvent MUST be called from the UI thread
func (c *roomViewConversation) roomEnableEvent() {
	removeRoomDisableClass(c.chatTextView)
}

// onSelfOccupantBanned MUST be called from the UI thread
func (c *roomViewConversation) onSelfOccupantBanned() {
	message := messageForSelfOccupantBanned()
	c.updateNotificationMessage(message)
	c.disableSendCapabilities()
}

// onSelfOccupantKicked MUST be called from the UI thread
func (c *roomViewConversation) onSelfOccupantKicked() {
	message := messageForSelfOccupantExpelled()
	c.updateNotificationMessage(message)
	c.disableSendCapabilities()
}

// onSelfOccupantVoiceRevoked MUST be called from the UI thread
func (c *roomViewConversation) onSelfOccupantVoiceRevoked() {
	message := messageForSelfOccupantVisitor()
	c.updateNotificationMessage(message)
	c.disableSendCapabilities()
}

func (c *roomViewConversation) selfOccupantJoinedEvent(nickname string, r data.Role) {
	go c.waitForDiscussionHistoryAndJoin(r, func() {
		c.enableSendCapabilitiesIfHasVoice(r.HasVoice())
		c.handleOccupantJoinedRoom(nickname)
	})
	c.disconnectMessagePrinted = false
}

// waitForDiscussionHistoryAndJoin MUST NOT be called from the UI thread
// The given argument `onDiscussionHistoryPrinted` WILL be called from the UI thread
func (c *roomViewConversation) waitForDiscussionHistoryAndJoin(r data.Role, onDiscussionHistoryPrinted func()) {
	<-c.historyPrinted
	doInUIThread(onDiscussionHistoryPrinted)
	c.selfOccupantJoined <- true
}

func (c *roomViewConversation) occupantUpdatedEvent(nickname string, r data.Role) {
	if c.selfOccupantNickname() == nickname {
		doInUIThread(func() {
			c.enableSendCapabilitiesIfHasVoice(r.HasVoice())
		})
	}
}

func (c *roomViewConversation) occupantLeftEvent(nickname string) {
	doInUIThread(func() {
		c.saveAndDisplayMessage(nickname, messageForSomeoneWhoLeftTheRoom(nickname), time.Now(), data.Left)
	})
}

func (c *roomViewConversation) occupantJoinedEvent(nickname string) {
	c.handleOccupantJoinedRoom(nickname)
}

func (c *roomViewConversation) messageEvent(tp, nickname, message string, timestamp time.Time) {
	switch tp {
	case "live":
		// We don't really care about self-incoming messages because
		// we already have those messages in the conversation's textview
		if c.selfOccupantNickname() != nickname {
			doInUIThread(func() {
				c.displayNewMessage(nickname, message, timestamp)
			})
		}
	case "delayed":
		doInUIThread(func() {
			c.displayNewMessage(nickname, message, timestamp)
		})
	default:
		c.log.WithField("type", tp).Warn("Unknown message event type")
	}
}

func (c *roomViewConversation) discussionHistoryEvent(dh *data.DiscussionHistory) {
	c.handleDiscussionHistory(dh)
	c.historyPrinted <- true
}

func (c *roomViewConversation) reopenRoomWindowEvent(dh *data.DiscussionHistory) {
	c.handleDiscussionHistory(dh)
}

func (c *roomViewConversation) handleDiscussionHistory(dh *data.DiscussionHistory) {
	if len(dh.GetHistory()) > 0 {
		doInUIThread(func() {
			c.clearBuffer()
			c.displayDiscussionHistory(dh.GetHistory())
			c.displayDivider()
		})
	}
}

// displayDiscussionHistory MUST be called from the UI thread
func (c *roomViewConversation) displayDiscussionHistory(history []*data.DelayedMessages) {
	for _, dm := range history {
		c.displayDiscussionHistoryDate(dm.GetDate())
		c.displayDiscussionHistoryMessages(dm.GetMessages())
	}
}

func (c *roomViewConversation) messageForbiddenEvent() {
	doInUIThread(func() {
		message := messageForSelfOccupantForbidden()
		c.displayErrorMessage(message)
	})
}

func (c *roomViewConversation) messageNotAcceptableEvent() {
	doInUIThread(func() {
		message := messageForSelfOccupantNotAcceptable()
		c.displayErrorMessage(message)
	})
}

func (c *roomViewConversation) subjectUpdatedEvent(nickname, subject string) {
	doInUIThread(func() {
		c.saveAndDisplayMessage(nickname, messageFromRoomSubjectUpdate(nickname, subject), time.Now(), data.Subject)
	})
}

func (c *roomViewConversation) subjectReceivedEvent(subject string) {
	go func() {
		<-c.selfOccupantJoined
		doInUIThread(func() {
			c.saveAndDisplayMessage("", messageForRoomSubject(subject), time.Now(), data.Subject)
		})
	}()
}

func (c *roomViewConversation) loggingEnabledEvent() {
	doInUIThread(func() {
		c.saveAndDisplayMessage("", messageForRoomPubliclyLogged(), time.Now(), data.Warning)
	})
}

func (c *roomViewConversation) loggingDisabledEvent() {
	doInUIThread(func() {
		c.saveAndDisplayMessage("", messageForRoomNotPubliclyLogged(), time.Now(), data.Warning)
	})
}

func (c *roomViewConversation) nonAnonymousRoomEvent() {
	doInUIThread(func() {
		c.saveAndDisplayMessage("", messageForRoomCanSeeRealJid(), time.Now(), data.RoomConfigurationChanged)
	})
}

func (c *roomViewConversation) semiAnonymousRoomEvent() {
	doInUIThread(func() {
		c.saveAndDisplayMessage("", messageForRoomCanNotSeeRealJid(), time.Now(), data.RoomConfigurationChanged)
	})
}

func (c *roomViewConversation) roomAnonymityChangedEvent(anonymityLevel string) {
	switch anonymityLevel {
	case "semi":
		c.semiAnonymousRoomEvent()
	case "no":
		c.nonAnonymousRoomEvent()
	default:
		c.log.Warn("room anonymity level unsupported")
	}
}

func (c *roomViewConversation) roomConfigChangedEvent(changes roomConfigChangedTypes, discoInfo data.RoomDiscoInfo) {
	doInUIThread(func() {
		messages := getRoomConfigUpdatedFriendlyMessages(changes, discoInfo)
		for _, m := range messages {
			c.saveAndDisplayMessage("", m, time.Now(), data.RoomConfigurationChanged)
		}
	})
}

func (c *roomViewConversation) selfOccupantRemovedEvent(nickname string) {
	c.occupantRemovedEvent(nickname)
	doInUIThread(func() {
		message := messageForSelfOccupantRoomConfiguration()
		c.updateNotificationMessage(message)
		c.disableSendCapabilities()
	})
}

func (c *roomViewConversation) occupantRemovedEvent(nickname string) {
	doInUIThread(func() {
		c.saveAndDisplayMessage(nickname, messageForMembersOnlyRoom(nickname), time.Now(), data.OccupantInformationChanged)
	})
}

func (c *roomViewConversation) selfOccupantRemovedOnAffiliationChangeEvent(nickname string) {
	c.occupantRemovedOnAffiliationChangeEvent(nickname)
	doInUIThread(func() {
		c.updateNotificationMessage(messageForSelfOccupantRemovedOnAffiliationChange())
		c.disableSendCapabilities()
	})
}

func (c *roomViewConversation) occupantRemovedOnAffiliationChangeEvent(nickname string) {
	doInUIThread(func() {
		c.saveAndDisplayMessage(nickname, messageForOccupantRemovedOnAffiliationChange(nickname), time.Now(), data.OccupantInformationChanged)
	})
}

// disableSendCapabilities MUST be called from the UI thread
func (c *roomViewConversation) disableSendCapabilities() {
	c.canSendMessages = false
	c.updateSendCapabilities()
}

// enableSendCapabilitiesIfHasVoice MUST be called from the UI thread
func (c *roomViewConversation) enableSendCapabilitiesIfHasVoice(hasVoice bool) {
	c.canSendMessages = hasVoice
	c.updateSendCapabilities()
}

// updateSendCapabilities MUST be called from the UI thread
func (c *roomViewConversation) updateSendCapabilities() {
	if c.canSendMessages {
		c.enableEntryAndSendButton()
		c.notificationBox.Hide()
	} else {
		c.disableEntryAndSendButton()
		c.notificationBox.Show()
	}
}

// getWrittenMessage MUST be called from the UI thread
func (c *roomViewConversation) getWrittenMessage() string {
	b := c.getMessageTextBuffer()
	starts, ends := b.GetBounds()
	return b.GetText(starts, ends, false)
}

// clearTypedMessage MUST be called from the UI thread
func (c *roomViewConversation) clearTypedMessage() {
	b := c.getMessageTextBuffer()
	b.SetText("")
}

// disableEntryAndSendButton MUST be called from the UI thread
func (c *roomViewConversation) disableEntryAndSendButton() {
	c.enableOrDisableFields(false)
}

// enableEntryAndSendButton MUST be called from the UI thread
func (c *roomViewConversation) enableEntryAndSendButton() {
	c.enableOrDisableFields(true)
}

// enableOrDisableFields MUST be called from the UI thread
func (c *roomViewConversation) enableOrDisableFields(v bool) {
	c.messageTextView.SetEditable(v)
	c.sendButton.SetSensitive(v)
	c.messageTextView.SetVisible(v)
	c.messageView.SetVisible(v)
}

// beforeSendingMessage MUST be called from the UI thread
func (c *roomViewConversation) beforeSendingMessage() {
	c.disableEntryAndSendButton()
}

// onSendMessageFinish MUST be called from the UI thread
func (c *roomViewConversation) onSendMessageFinish() {
	c.clearTypedMessage()
	if c.canSendMessages {
		c.enableEntryAndSendButton()
		c.messageTextView.GrabFocus()
	}
}

// onSendMessageFailed MUST be called from the UI thread
func (c *roomViewConversation) onSendMessageFailed(err error) {
	c.log.WithError(err).Error("Failed to send the message to all occupants")

	message := messageNotSent()
	c.displayErrorMessage(message)
}

// onKeyPress MUST be called from the UI thread
func (c *roomViewConversation) onKeyPress(_ gtki.Widget, ev gdki.Event) bool {
	if isNormalEnter(g.gdk.EventKeyFrom(ev)) {
		c.sendMessage()
		return true
	}

	return false
}

// onSendMessage MUST be called from the UI thread
func (c *roomViewConversation) onSendMessage() {
	c.sendMessage()
}

// sendMessage MUST be called from the UI thread
func (c *roomViewConversation) sendMessage() {
	if !c.canSendMessages {
		c.log.Warn("Trying to send a message to all occupants without having voice")
		return
	}

	c.beforeSendingMessage()
	defer c.onSendMessageFinish()

	m := c.getWrittenMessage()
	if m == "" {
		return
	}

	err := c.account.session.SendMUCMessage(c.roomID.String(), c.account.Account(), m)
	if err != nil {
		c.onSendMessageFailed(err)
		return
	}

	c.displayNewMessage(c.selfOccupantNickname(), m, time.Now())

	doALittleBitLater(func() {
		scrollToBottom(c.chatScrolledWindow)
	})
}

// onEdgeReached MUST be called from the UI thread
func (c *roomViewConversation) onEdgeReached(_ gtki.ScrolledWindow, pos int) {
	if pos == bottomPositionValue {
		c.maxAdjustment = c.chatScrolledWindow.GetVAdjustment().GetValue()
	}
}

func (c *roomViewConversation) onAdjustmentChanged() {
	if c.currentAdjustment == c.maxAdjustment {
		doALittleBitLater(func() {
			scrollToBottom(c.chatScrolledWindow)
		})
	}
}

func (c *roomViewConversation) updateCurrentAdjustmentValue() {
	c.currentAdjustment = c.chatScrolledWindow.GetVAdjustment().GetValue()
}

func (c *roomViewConversation) updateNotificationMessage(m string) {
	c.messageBoxNotification.updateMessage(m)
}
