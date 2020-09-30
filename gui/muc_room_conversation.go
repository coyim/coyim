package gui

import (
	"time"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewConversation struct {
	tags gtki.TextTagTable

	view gtki.Box      `gtk-widget:"roomConversation"`
	text gtki.TextView `gtk-widget:"roomChatTextView"`

	log coylog.Logger
}

func (v *roomView) newRoomViewConversation() *roomViewConversation {
	c := &roomViewConversation{}

	c.initBuilder()
	c.initSubscribers(v)
	c.initTagsAndTextBuffer(v)

	return c
}

func (c *roomViewConversation) initBuilder() {
	builder := newBuilder("MUCRoomConversation")
	panicOnDevError(builder.bindObjects(c))
}

func (c *roomViewConversation) initSubscribers(v *roomView) {
	v.subscribeAll("conversation", roomViewEventObservers{
		"occupantLeftEvent": func(ei roomViewEventInfo) {
			doInUIThread(func() {
				c.displayNotificationWhenOccupantLeftTheRoom(ei["nickname"])
			})
		},
		"occupantJoinedEvent": func(ei roomViewEventInfo) {
			doInUIThread(func() {
				c.displayNotificationWhenOccupantJoinedRoom(ei["nickname"])
			})
		},
		"messageReceivedEvent": func(ei roomViewEventInfo) {
			doInUIThread(func() {
				c.displayNewLiveMessage(
					ei["nickname"],
					ei["subject"],
					ei["message"],
				)
			})
		},
		"loggingEnabledEvent": func(roomViewEventInfo) {
			doInUIThread(func() {
				c.displayWarningMessage(i18n.Local("This room is now publicly logged, meaning that everything you and the others in the room say or do can be made public on a website."))
			})
		},
		"loggingDisabledEvent": func(roomViewEventInfo) {
			doInUIThread(func() {
				c.displayWarningMessage(i18n.Local("This room is no longer publicly logged."))
			})
		},
	})
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
