package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

type roomViewConversation struct {
	tags    gtki.TextTagTable
	roomID  jid.Bare
	account *account

	view       gtki.Box      `gtk-widget:"roomConversation"`
	text       gtki.TextView `gtk-widget:"roomChatTextView"`
	newText    gtki.Entry    `gtk-widget:"textConversation"`
	sendButton gtki.Button   `gtk-widget:"button-send"`

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
	c.initSubscribers(v)
	c.initTagsAndTextBuffer(v)

	return c
}

func (c *roomViewConversation) initBuilder() {
	builder := newBuilder("MUCRoomConversation")
	panicOnDevError(builder.bindObjects(c))

	builder.ConnectSignals(map[string]interface{}{
		"on_send_message": c.onSendMessage,
		"on_key_press":    c.onKeyPress,
	})
}

func (c *roomViewConversation) initSubscribers(v *roomView) {
	v.subscribe("conversation", func(ev roomViewEvent) {
		switch t := ev.(type) {
		case occupantLeftEvent:
			c.occupantLeftEvent(t.nickname)
		case occupantJoinedEvent:
			c.occupantJoinedEvent(t.nickname)
		case messageEvent:
			c.messageEvent(t.tp, t.nickname, t.subject, t.message)
		case loggingEnabledEvent:
			c.loggingEnabledEvent()
		case loggingDisabledEvent:
			c.loggingDisabledEvent()
		}
	})
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

func (c *roomViewConversation) messageEvent(tp, nickname, subject, message string) {
	doInUIThread(func() {
		switch tp {
		case "received":
			c.displayNewLiveMessage(nickname, subject, message)
		default:
			c.log.WithField("type", tp).Warn("Unknow message event type")
		}
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
	text, _ := c.newText.GetText()
	return text
}

func (c *roomViewConversation) clearTypedMessage() {
	c.newText.SetText("")
}

func (c *roomViewConversation) disableEntryAndSendButton() {
	c.newText.SetEditable(false)
	c.sendButton.SetSensitive(false)
}

func (c *roomViewConversation) enableEntryAndSendButton() {
	c.newText.SetEditable(true)
	c.sendButton.SetSensitive(true)
}

func (c *roomViewConversation) beforeSendingMessage() {
	c.disableEntryAndSendButton()
}

func (c *roomViewConversation) onSendMessageFinish() {
	c.enableEntryAndSendButton()
	c.clearTypedMessage()
}

func (c *roomViewConversation) onSendMessageFailed(err error) {
	c.log.WithError(err).Warn("failed to send the message")
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
	c.beforeSendingMessage()
	defer c.onSendMessageFinish()

	message := c.getTypedMessage()
	if message == "" {
		return
	}

	err := c.account.session.SendMUCMessage(c.roomID.String(), c.account.Account(), message)
	if err != nil {
		c.onSendMessageFailed(err)
	}
}
