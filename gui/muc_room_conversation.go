package gui

import (
	"fmt"
	"time"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3adapter/pangoi"
)

type roomViewConversation struct {
	tags gtki.TextTagTable

	view               gtki.Box            `gtk-widget:"roomConversation"`
	roomChatTextView   gtki.TextView       `gtk-widget:"roomChatTextView"`
	roomChatScrollView gtki.ScrolledWindow `gtk-widget:"roomChatScrollView"`

	log coylog.Logger
}

func (v *roomView) newRoomViewConversation() *roomViewConversation {
	c := &roomViewConversation{}

	c.initBuilder()
	c.initSubscribers(v)
	c.initTags(v)

	return c
}

func (c *roomViewConversation) initBuilder() {
	builder := newBuilder("MUCRoomConversation")
	panicOnDevError(builder.bindObjects(c))
}

func (c *roomViewConversation) initSubscribers(v *roomView) {
	v.subscribeAll("conversation", roomViewEventObservers{
		"occupantLeftEvent": func(ei roomViewEventInfo) {
			c.displayNotificationWhenOccupantLeftTheRoom(ei["nickname"])
		},
		"occupantJoinedEvent": func(ei roomViewEventInfo) {
			c.displayNotificationWhenOccupantJoinedRoom(ei["nickname"])
		},
		"messageReceivedEvent": func(ei roomViewEventInfo) {
			c.displayNewLiveMessage(
				ei["nickname"],
				ei["subject"],
				ei["message"],
			)
		},
		"loggingEnabledEvent": func(roomViewEventInfo) {
			msg := i18n.Local("This room is now publicly logged, meaning that everything you and the others in the room say or do can be made public on a website.")
			v.conv.displayWarningMessage(msg)
		},
		"loggingDisabledEvent": func(roomViewEventInfo) {
			msg := i18n.Local("This room is no longer publicly logged.")
			v.conv.displayWarningMessage(msg)
		},
	})
}

func (c *roomViewConversation) initTags(v *roomView) {
	c.tags = c.newMUCTableStyleTags(v.u)
	c.roomChatTextView.SetBuffer(c.createTextBuffer())
}

func getTimestamp() string {
	return time.Now().Format("15:04:05")
}

func (c *roomViewConversation) createConversationTag(name string, properties map[string]interface{}) gtki.TextTag {
	tag, _ := g.gtk.TextTagNew(name)
	for attribute, value := range properties {
		_ = tag.SetProperty(attribute, value)
	}
	return tag
}

func (c *roomViewConversation) createWarningTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag("warning", map[string]interface{}{
		"foreground": cs.warningForeground,
	})
}

func (c *roomViewConversation) createLeftRoomTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag("leftRoomText", map[string]interface{}{
		"foreground": cs.someoneLeftForeground,
		"style":      pangoi.STYLE_ITALIC,
	})
}

func (c *roomViewConversation) createJoinedRoomTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag("joinedRoomText", map[string]interface{}{
		"foreground": cs.someoneJoinedForeground,
		"style":      pangoi.STYLE_ITALIC,
	})
}

func (c *roomViewConversation) createTimestampTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag("timestampText", map[string]interface{}{
		"foreground": cs.timestampForeground,
		"style":      pangoi.STYLE_NORMAL,
	})
}

func (c *roomViewConversation) createNicknameTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag("nicknameText", map[string]interface{}{
		"foreground": cs.nicknameForeground,
		"style":      pangoi.STYLE_NORMAL,
	})
}

func (c *roomViewConversation) createSubjectTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag("subjectText", map[string]interface{}{
		"foreground": cs.subjectForeground,
		"style":      pangoi.STYLE_ITALIC,
	})
}

func (c *roomViewConversation) createMessageTag(cs mucColorSet) gtki.TextTag {
	return c.createConversationTag("messageText", map[string]interface{}{
		"foreground": cs.messageForeground,
		"style":      pangoi.STYLE_NORMAL,
	})
}

func (c *roomViewConversation) newMUCTableStyleTags(u *gtkUI) gtki.TextTagTable {
	table, _ := g.gtk.TextTagTableNew()
	cs := u.currentMUCColorSet()

	tags := []func(mucColorSet) gtki.TextTag{
		c.createWarningTag,
		c.createLeftRoomTag,
		c.createJoinedRoomTag,
		c.createTimestampTag,
		c.createNicknameTag,
		c.createSubjectTag,
		c.createMessageTag,
	}

	for _, t := range tags {
		table.Add(t(cs))
	}

	return table
}

func (c *roomViewConversation) createTextBuffer() gtki.TextBuffer {
	buf, _ := g.gtk.TextBufferNew(c.tags)
	return buf
}

func (c *roomViewConversation) addNewLine() {
	buf, _ := c.roomChatTextView.GetBuffer()
	i := buf.GetEndIter()

	buf.Insert(i, "\n")
}

func (c *roomViewConversation) addTimestamp() {
	buf, _ := c.roomChatTextView.GetBuffer()
	i := buf.GetEndIter()
	t := fmt.Sprintf("[%s] ", getTimestamp())

	buf.InsertWithTagByName(i, t, "timestampText")
}

func (c *roomViewConversation) addTextToChatTextUsingTagID(text string, tagName string) {
	buf, _ := c.roomChatTextView.GetBuffer()
	i := buf.GetEndIter()

	if tagName != "" {
		buf.InsertWithTagByName(i, text, tagName)
	} else {
		buf.Insert(i, text)
	}
}

func (c *roomViewConversation) addLineTextToChatTextUsingTagID(text string, tagName string) {
	c.addTimestamp()
	c.addTextToChatTextUsingTagID(text, tagName)
	c.addNewLine()
}

// displayNotificationWhenOccupantJoinedRoom MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenOccupantJoinedRoom(nickname string) {
	text := i18n.Localf("%s joined the room", nickname)
	c.addLineTextToChatTextUsingTagID(text, "joinedRoomText")
}

// displayNotificationWhenOccupantLeftTheRoom MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenOccupantLeftTheRoom(nickname string) {
	text := i18n.Localf("%s left the room", nickname)
	c.addLineTextToChatTextUsingTagID(text, "leftRoomText")
}

// displayNewLiveMessage MUST be called from the UI thread
func (c *roomViewConversation) displayNewLiveMessage(nickname, subject, message string) {
	c.addTimestamp()
	nicknameText := i18n.Localf("%s: ", nickname)
	c.addTextToChatTextUsingTagID(nicknameText, "nicknameText")
	if subject != "" {
		subjectText := i18n.Localf("[%s] ", subject)
		c.addTextToChatTextUsingTagID(subjectText, "subjectText")
	}
	c.addTextToChatTextUsingTagID(message, "messageText")
	c.addNewLine()
}

// displayWarningMessage MUST be called from the UI thread
func (c *roomViewConversation) displayWarningMessage(message string) {
	c.addLineTextToChatTextUsingTagID(message, "warning")
}
