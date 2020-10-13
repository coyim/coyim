package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

type roomViewConversation struct {
	tags            gtki.TextTagTable
	roomID          jid.Bare
	account         *account
	canSendMessages bool

	view                  gtki.Box            `gtk-widget:"room-conversation"`
	chatTextView          gtki.TextView       `gtk-widget:"chat-text-view"`
	messageView           gtki.Box            `gtk-widget:"message-view"`
	messageScrolledWindow gtki.ScrolledWindow `gtk-widget:"message-scrolled-window"`
	messageTextView       gtki.TextView       `gtk-widget:"message-text-view"`
	sendButton            gtki.Button         `gtk-widget:"message-send-button"`

	subject string

	log coylog.Logger
}

func (v *roomView) newRoomViewConversation() *roomViewConversation {
	c := &roomViewConversation{
		roomID:  v.roomID(),
		account: v.account,
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

	messageScrolledWindoStyle := providerWithStyle("scrolledwindow", style{
		"border": "none",
	})

	updateWithStyle(c.messageScrolledWindow, messageScrolledWindoStyle)
}

func (c *roomViewConversation) initDefaults(v *roomView) {
	c.subject = v.room.Subject
	c.disableEntryAndSendButton()
}

func (c *roomViewConversation) initSubscribers(v *roomView) {
	v.subscribe("conversation", func(ev roomViewEvent) {
		switch t := ev.(type) {
		case occupantLeftEvent:
			c.occupantLeftEvent(t.nickname)
		case occupantJoinedEvent:
			c.occupantJoinedEvent(t.nickname)
		case messageEvent:
			c.messageEvent(v.room, t.tp, t.nickname, t.message)
		case messageForbidden:
			c.messageForbiddenEvent()
		case messageNotAcceptable:
			c.messageNotAcceptable()
		case subjectEvent:
			if c.subject != "" && c.subject != t.subject {
				c.subjectEvent(i18n.Localf("The room subject was changed to: %s", t.subject))
				return
			}
			c.subject = t.subject
			c.subjectEvent(i18n.Localf("The room subject is: %s", t.subject))
		case loggingEnabledEvent:
			c.loggingEnabledEvent()
		case loggingDisabledEvent:
			c.loggingDisabledEvent()
		case occupantUpdatedEvent:
			c.enableSendCapabilitiesIfHasVoice(v.room.SelfOccupant())
		}
	})
}

func (c *roomViewConversation) enableSendCapabilitiesIfHasVoice(occupant *muc.Occupant) {
	if occupant != nil && occupant.Role != nil && occupant.Role.HasVoice() {
		c.canSendMessages = true
		c.enableEntryAndSendButton()
		return
	}

	c.canSendMessages = false
	c.disableEntryAndSendButton()
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

func (c *roomViewConversation) messageEvent(r *muc.Room, tp, nickname, message string) {
	// Display of self messages is disabled because the own typed messages
	// are displayed automatically on chat screen
	if r.SelfOccupantNickname() == nickname {
		return
	}

	doInUIThread(func() {
		switch tp {
		case "received":
			c.displayNewLiveMessage(nickname, message)
		default:
			c.log.WithField("type", tp).Warn("Unknow message event type")
		}
	})
}

func (c *roomViewConversation) messageForbiddenEvent() {
	doInUIThread(func() {
		c.displayErrorMessage(i18n.Local("You are forbidden to send messages to this room."))
	})
}

func (c *roomViewConversation) messageNotAcceptable() {
	doInUIThread(func() {
		c.displayErrorMessage(i18n.Local("Your messages to this room aren't accepted."))
	})
}

func (c *roomViewConversation) subjectEvent(subject string) {
	doInUIThread(func() {
		c.displayRoomSubject(subject)
	})
}

func (c *roomViewConversation) loggingEnabledEvent() {
	doInUIThread(func() {
		c.displayWarningMessage(i18n.Local("This room is now publicly logged, meaning that everything you and the others in the room say or do can be made public on a website."))
	})
}

func (c *roomViewConversation) loggingDisabledEvent() {
	doInUIThread(func() {
		c.displayWarningMessage(i18n.Local("This room is no longer publicly logged."))
	})
}

func (c *roomViewConversation) getTypedMessage() string {
	b := c.getMessageTextBuffer()
	starts, ends := b.GetBounds()
	return b.GetText(starts, ends, false)
}

func (c *roomViewConversation) clearTypedMessage() {
	b := c.getMessageTextBuffer()
	b.SetText("")
}

func (c *roomViewConversation) disableEntryAndSendButton() {
	c.messageTextView.SetEditable(false)
	c.sendButton.SetSensitive(false)
	c.messageTextView.SetVisible(false)
	c.messageView.SetVisible(false)
}

func (c *roomViewConversation) enableEntryAndSendButton() {
	c.messageTextView.SetEditable(true)
	c.sendButton.SetSensitive(true)
	c.messageTextView.SetVisible(true)
	c.messageView.SetVisible(true)
}

func (c *roomViewConversation) beforeSendingMessage() {
	c.disableEntryAndSendButton()
}

func (c *roomViewConversation) onSendMessageFinish() {
	c.enableEntryAndSendButton()
	c.clearTypedMessage()
}

func (c *roomViewConversation) onSendMessageFailed(err error) {
	c.log.WithError(err).Error("failed to send the message")
	doInUIThread(func() {
		c.displayErrorMessage(i18n.Local("The message couldn't be sent, please try again"))
	})
}

func (c *roomViewConversation) onKeyPress(_ gtki.Widget, ev gdki.Event) bool {
	evk := g.gdk.EventKeyFrom(ev)
	ret := false

	if isNormalEnter(evk) {
		c.onSendMessage()
		ret = true
	}

	return ret
}

func (c *roomViewConversation) onSendMessage() {
	if c.canSendMessages {
		c.sendMessage()
	}
}

func (c *roomViewConversation) sendMessage() {
	c.beforeSendingMessage()
	defer c.onSendMessageFinish()

	message := c.getTypedMessage()
	if message == "" {
		return
	}

	err := c.account.session.SendMUCMessage(c.roomID.String(), c.account.Account(), message)
	if err != nil {
		c.onSendMessageFailed(err)
		return
	}

	c.displayTextLineWithTimestamp(message, "message")
}
