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
	tags *mucStyleTags

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
		occupantLeft: func(ei roomViewEventInfo) {
			c.displayNotificationWhenOccupantLeftTheRoom(ei.nickname)
		},
		occupantJoined: func(ei roomViewEventInfo) {
			c.displayNotificationWhenOccupantJoinedRoom(ei.nickname)
		},
		messageReceived: func(ei roomViewEventInfo) {
			c.displayNewLiveMessage(
				ei.nickname,
				ei.subject,
				ei.message,
			)
		},
		loggingEnabled: func(roomViewEventInfo) {
			msg := i18n.Local("This room is now publicly logged, meaning that everything you and the others in the room say or do can be made public on a website.")
			v.conv.displayWarningMessage(msg)
		},
		loggingDisabled: func(roomViewEventInfo) {
			msg := i18n.Local("This room is no longer publicly logged.")
			v.conv.displayWarningMessage(msg)
		},
	})
}

func (c *roomViewConversation) initTags(v *roomView) {
	t := c.getStyleTags(v.u)
	c.roomChatTextView.SetBuffer(t.createTextBuffer())
}

// TODO: I don't think I see the need for the mucStyleTags structure
// It doesn't really help us understand the code
type mucStyleTags struct {
	table gtki.TextTagTable
}

func getTimestamp() string {
	return time.Now().Format("15:04:05")
}

func (c *roomViewConversation) getStyleTags(u *gtkUI) *mucStyleTags {
	if c.tags == nil {
		c.tags = c.newStyleTags(u)
	}
	return c.tags
}

func (c *roomViewConversation) newStyleTags(u *gtkUI) *mucStyleTags {
	// TODO: for now we are using a default styles, but we can improve it
	// if we define a structure with a predefined colors pallete based on kind
	// of messages to show like entering a room, leaving the room, incoming
	// message, etc
	t := &mucStyleTags{}

	t.table, _ = g.gtk.TextTagTableNew()
	cset := u.currentColorSet()

	warningTag, _ := g.gtk.TextTagNew("warning")
	_ = warningTag.SetProperty("foreground", cset.warningForeground)

	leftRoomTag, _ := g.gtk.TextTagNew("leftRoomText")
	_ = leftRoomTag.SetProperty("foreground", cset.mucSomeoneLeftForeground)
	_ = leftRoomTag.SetProperty("style", pangoi.STYLE_ITALIC)

	joinedRoomTag, _ := g.gtk.TextTagNew("joinedRoomText")
	_ = joinedRoomTag.SetProperty("foreground", cset.mucSomeoneJoinedForeground)
	_ = joinedRoomTag.SetProperty("style", pangoi.STYLE_ITALIC)

	timestampTag, _ := g.gtk.TextTagNew("timestampText")
	_ = timestampTag.SetProperty("foreground", cset.mucTimestampForeground)
	_ = timestampTag.SetProperty("style", pangoi.STYLE_NORMAL)

	nicknameTag, _ := g.gtk.TextTagNew("nicknameText")
	_ = nicknameTag.SetProperty("foreground", cset.mucNicknameForeground)
	_ = nicknameTag.SetProperty("style", pangoi.STYLE_NORMAL)

	subjectTag, _ := g.gtk.TextTagNew("subjectText")
	_ = subjectTag.SetProperty("foreground", cset.mucSubjectForeground)
	_ = subjectTag.SetProperty("style", pangoi.STYLE_ITALIC)

	messageTag, _ := g.gtk.TextTagNew("messageText")
	_ = messageTag.SetProperty("foreground", cset.mucMessageForeground)
	_ = messageTag.SetProperty("style", pangoi.STYLE_NORMAL)

	t.table.Add(warningTag)
	t.table.Add(leftRoomTag)
	t.table.Add(joinedRoomTag)
	t.table.Add(timestampTag)
	t.table.Add(nicknameTag)
	t.table.Add(subjectTag)
	t.table.Add(messageTag)

	return t
}

func (t *mucStyleTags) createTextBuffer() gtki.TextBuffer {
	buf, _ := g.gtk.TextBufferNew(t.table)
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
	// TODO: i18n
	text := fmt.Sprintf("%s joined the room", nickname)
	c.addLineTextToChatTextUsingTagID(text, "joinedRoomText")
}

// displayNotificationWhenOccupantLeftTheRoom MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenOccupantLeftTheRoom(nickname string) {
	// TODO: i18n
	text := fmt.Sprintf("%s left the room", nickname)
	c.addLineTextToChatTextUsingTagID(text, "leftRoomText")
}

// displayNewLiveMessage MUST be called from the UI thread
func (c *roomViewConversation) displayNewLiveMessage(nickname, subject, message string) {
	c.addTimestamp()
	// TODO: i18n
	nicknameText := fmt.Sprintf("%s: ", nickname)
	c.addTextToChatTextUsingTagID(nicknameText, "nicknameText")
	if subject != "" {
		// TODO: i18n
		subjectText := fmt.Sprintf("[%s] ", subject)
		c.addTextToChatTextUsingTagID(subjectText, "subjectText")
	}
	c.addTextToChatTextUsingTagID(message, "messageText")
	c.addNewLine()
}

// displayWarningMessage MUST be called from the UI thread
func (c *roomViewConversation) displayWarningMessage(message string) {
	c.addLineTextToChatTextUsingTagID(message, "warning")
}
