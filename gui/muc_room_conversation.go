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

type roomViewConversation struct {
	tags                 gtki.TextTagTable
	roomID               jid.Bare
	account              *account
	canSendMessages      bool
	selfOccupantNickname func() string

	view                  gtki.Box            `gtk-widget:"room-conversation"`
	chatTextView          gtki.TextView       `gtk-widget:"chat-text-view"`
	messageView           gtki.Box            `gtk-widget:"message-view"`
	messageScrolledWindow gtki.ScrolledWindow `gtk-widget:"message-scrolled-window"`
	messageTextView       gtki.TextView       `gtk-widget:"message-text-view"`
	sendButton            gtki.Button         `gtk-widget:"message-send-button"`
	sendButtonIcon        gtki.Image          `gtk-widget:"message-send-icon"`
	notificationBox       gtki.Box            `gtk-widget:"notification-box"`

	log coylog.Logger
}

func (v *roomView) newRoomViewConversation() *roomViewConversation {
	c := &roomViewConversation{
		roomID:               v.roomID(),
		account:              v.account,
		selfOccupantNickname: v.room.SelfOccupantNickname,
	}

	c.log = c.account.log.WithFields(log.Fields{
		"who":  "roomViewConversation",
		"room": c.roomID,
	})

	c.initBuilder()
	c.initDefaults(v)
	c.initSubscribers(v)
	c.initTagsAndTextBuffers(v)

	return c
}

func (c *roomViewConversation) initBuilder() {
	builder := newBuilder("MUCRoomConversation")
	panicOnDevError(builder.bindObjects(c))

	builder.ConnectSignals(map[string]interface{}{
		"on_send_message": c.onSendMessage,
		"on_key_press":    c.onKeyPress,
	})

	updateWithStyle(c.messageScrolledWindow, providerWithStyle("scrolledwindow", style{
		"border": "none",
	}))
}

func (c *roomViewConversation) initDefaults(v *roomView) {
	c.sendButtonIcon.SetFromPixbuf(getMUCIconPixbuf("send"))

	voiceNotification := newRoomVoiceNotification()
	c.notificationBox.Add(voiceNotification.widget())

	c.disableEntryAndSendButton()
	if v.room.SelfOccupant().HasVoice() {
		c.enableSendCapabilitiesIfHasVoice(true)
	}
}

func (c *roomViewConversation) initSubscribers(v *roomView) {
	v.subscribe("conversation", func(ev roomViewEvent) {
		switch t := ev.(type) {
		case occupantLeftEvent:
			c.occupantLeftEvent(t.nickname)
		case occupantJoinedEvent:
			c.occupantJoinedEvent(t.nickname)
		case occupantSelfJoinedEvent:
			c.occupantSelfJoinedEvent(t.role)
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
		case nonAnonymousRoomEvent:
			c.nonAnonymousRoomEvent()
		case semiAnonymousRoomEvent:
			c.semiAnonymousRoomEvent()
		case roomConfigChangedEvent:
			c.roomConfigChangedEvent(t.changes, t.config)
		case selfOccupantRemovedEvent:
			c.selfOccupantRemovedEvent(t.nickname)
		case occupantRemovedEvent:
			c.occupantRemovedEvent(t.nickname)
		}
	})
}

func (c *roomViewConversation) occupantSelfJoinedEvent(r data.Role) {
	doInUIThread(func() {
		c.enableSendCapabilitiesIfHasVoice(r.HasVoice())
	})
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
		c.displayNotificationWhenOccupantLeftTheRoom(nickname)
	})
}

func (c *roomViewConversation) occupantJoinedEvent(nickname string) {
	doInUIThread(func() {
		c.displayNotificationWhenOccupantJoinedRoom(nickname)
	})
}

func (c *roomViewConversation) messageEvent(tp, nickname, message string, timestamp time.Time) {
	switch tp {
	case "live":
		// We don't really care about self-incoming messages because
		// we already have those messages in the conversation's textview
		if c.selfOccupantNickname() != nickname {
			doInUIThread(func() {
				c.displayLiveMessage(nickname, message, timestamp)
			})
		}
	case "delayed":
		doInUIThread(func() {
			c.displayDelayedMessage(nickname, message, timestamp)
		})
	default:
		c.log.WithField("type", tp).Warn("Unknow message event type")
	}
}

func (c *roomViewConversation) discussionHistoryEvent(dh *data.DiscussionHistory) {
	doInUIThread(func() {
		for _, dm := range dh.GetHistory() {
			c.displayDiscussionHistoryDate(dm.GetDate())
			c.displayDiscussionHistoryMessages(dm.GetMessages())
		}
		c.displayDivider()
	})
}

func (c *roomViewConversation) messageForbiddenEvent() {
	doInUIThread(func() {
		c.displayErrorMessage(i18n.Local("You are forbidden to send messages to this room."))
	})
}

func (c *roomViewConversation) messageNotAcceptableEvent() {
	doInUIThread(func() {
		c.displayErrorMessage(i18n.Local("Your messages to this room aren't accepted."))
	})
}

func (c *roomViewConversation) subjectUpdatedEvent(nickname, subject string) {
	doInUIThread(func() {
		c.displayNewInfoMessage(getDisplayRoomSubjectForNickname(nickname, subject))
	})
}

func (c *roomViewConversation) subjectReceivedEvent(subject string) {
	doInUIThread(func() {
		c.displayNewInfoMessage(getDisplayRoomSubject(subject))
	})
}

func (c *roomViewConversation) loggingEnabledEvent() {
	doInUIThread(func() {
		c.displayWarningMessage(i18n.Local("This room is now publicly logged."))
	})
}

func (c *roomViewConversation) loggingDisabledEvent() {
	doInUIThread(func() {
		c.displayWarningMessage(i18n.Local("This room is not publicly logged anymore."))
	})
}

func (c *roomViewConversation) nonAnonymousRoomEvent() {
	doInUIThread(func() {
		c.displayNewConfigurationMessage(i18n.Local("Your real JID can now be seen by anyone."))
	})
}

func (c *roomViewConversation) semiAnonymousRoomEvent() {
	doInUIThread(func() {
		c.displayNewConfigurationMessage(i18n.Local("Your real JID can now be seen only by moderators."))
	})
}

func (c *roomViewConversation) roomConfigChangedEvent(changes roomConfigChangedTypes, config data.RoomConfig) {
	doInUIThread(func() {
		messages := getRoomConfigUpdatedFriendlyMessages(changes, config)
		for _, m := range messages {
			c.displayNewConfigurationMessage(m)
		}
	})
}

func (c *roomViewConversation) selfOccupantRemovedEvent(nickname string) {
	doInUIThread(func() {
		c.displayWarningMessage(i18n.Local("You have been thrown out of this room because it's now a members only room."))
		c.disableSendCapabilities()
	})
}

func (c *roomViewConversation) occupantRemovedEvent(nickname string) {
	doInUIThread(func() {
		c.displayWarningMessage(i18n.Localf("\"%s\" was removed from this room.", nickname))
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
	c.displayErrorMessage(i18n.Local("The message couldn't be sent, please try again"))
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

	c.displayLiveMessage(c.selfOccupantNickname(), m, time.Now())
}

func getDisplayRoomSubjectForNickname(nickname, subject string) string {
	if nickname == "" {
		return i18n.Localf("Someone has updated the room subject to: \"%s\"", subject)
	}

	return i18n.Localf("%s updated the room subject to \"%s\"", nickname, subject)
}

func getDisplayRoomSubject(subject string) string {
	if subject == "" {
		return i18n.Local("The room does not have a subject")
	}

	return i18n.Localf("The room subject is \"%s\"", subject)
}
