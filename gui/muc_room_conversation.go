package gui

import (
	"time"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewConversation struct {
	tags    gtki.TextTagTable
	roomID  jid.Bare
	session access.Session

	occupantID func() (jid.Full, error)

	view    gtki.Box      `gtk-widget:"roomConversation"`
	text    gtki.TextView `gtk-widget:"roomChatTextView"`
	newText gtki.Entry    `gtk-widget:"textConversation"`

	log coylog.Logger
}

func (v *roomView) newRoomViewConversation(s access.Session) *roomViewConversation {
	c := &roomViewConversation{
		roomID:     v.roomID(),
		occupantID: v.occupantID,
		session:    s,
	}

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
	c.newText.SetEditable(false)
	text, _ := c.newText.GetText()
	c.newText.SetText("")
	c.newText.SetEditable(true)

	return text
}

func (c *roomViewConversation) onSendMessage() {
	text := c.getTypedMessage()
	if text == "" {
		return
	}

	n, err := c.occupantID()
	if err != nil {
		//TODO: Show a friendly message to the user
		return
	}

	err = c.session.SendMUCMessage(c.roomID.String(), n.String(), text)
	if err != nil {
		//TODO: Show a friendly message to the user
		c.log.WithError(err).Warn("Failed to send the message")
	}
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

func (c *roomViewConversation) getTextBuffer() gtki.TextBuffer {
	b, _ := c.text.GetBuffer()
	return b
}

func (c *roomViewConversation) addNewLine() {
	c.addText("\n")
}

func (c *roomViewConversation) addText(text string) {
	b := c.getTextBuffer()
	b.Insert(b.GetEndIter(), text)
}

func (c *roomViewConversation) addTextWithTag(text string, tag string) {
	b := c.getTextBuffer()
	b.InsertWithTagByName(b.GetEndIter(), text, tag)
}

func (c *roomViewConversation) addTextLineWithTimestamp(text string, tag string) {
	c.displayTimestamp()
	c.addTextWithTag(text, tag)
	c.addNewLine()
}

// displayTimestamp MUST be called from the UI thread
func (c *roomViewConversation) displayTimestamp() {
	c.addTextWithTag(i18n.Localf("[%s] ", getTimestamp()), "timestamp")
}

// displayNotificationWhenOccupantJoinedRoom MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenOccupantJoinedRoom(nickname string) {
	c.addTextLineWithTimestamp(i18n.Localf("%s joined the room", nickname), "joinedRoom")
}

// displayNotificationWhenOccupantLeftTheRoom MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenOccupantLeftTheRoom(nickname string) {
	c.addTextLineWithTimestamp(i18n.Localf("%s left the room", nickname), "leftRoom")
}

// displayNickname MUST be called from the UI thread
func (c *roomViewConversation) displayNickname(nickname string) {
	c.addTextWithTag(i18n.Localf("%s: ", nickname), "nickname")
}

// displayRoomSubject MUST be called from the UI thread
func (c *roomViewConversation) displayRoomSubject(subject string) {
	c.addTextWithTag(i18n.Localf("[%s] ", subject), "subject")
}

// displayMessage MUST be called from the UI thread
func (c *roomViewConversation) displayMessage(message string) {
	c.addTextWithTag(message, "message")
}

// displayNewLiveMessage MUST be called from the UI thread
func (c *roomViewConversation) displayNewLiveMessage(nickname, subject, message string) {
	c.displayTimestamp()

	c.displayNickname(nickname)
	if subject != "" {
		c.displayRoomSubject(subject)
	}
	c.displayMessage(message)

	c.addNewLine()
}

// displayWarningMessage MUST be called from the UI thread
func (c *roomViewConversation) displayWarningMessage(message string) {
	c.addTextLineWithTimestamp(message, "warning")
}

func getTimestamp() string {
	return time.Now().Format("15:04:05")
}
