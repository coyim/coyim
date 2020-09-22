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
	view               gtki.Box            `gtk-widget:"roomConversation"`
	roomChatTextView   gtki.TextView       `gtk-widget:"roomChatTextView"`
	roomChatScrollView gtki.ScrolledWindow `gtk-widget:"roomChatScrollView"`

	tags *mucStyleTags
	log  coylog.Logger
}

func (v *roomView) newRoomViewConversation() *roomViewConversation {
	c := &roomViewConversation{}

	builder := newBuilder("MUCRoomConversation")
	panicOnDevError(builder.bindObjects(c))

	t := c.getStyleTags(v.u)
	c.roomChatTextView.SetBuffer(t.createTextBuffer())

	v.subscribe("conversation", occupantLeft, func(ei roomViewEventInfo) {
		c.displayNotificationWhenOccupantLeftTheRoom(ei.nickname)
	})

	v.subscribe("conversation", occupantJoined, func(ei roomViewEventInfo) {
		c.displayNotificationWhenOccupantJoinedRoom(ei.nickname)
	})

	v.subscribe("conversation", messageReceived, func(ei roomViewEventInfo) {
		c.displayNewLiveMessage(
			ei.nickname,
			ei.subject,
			ei.message,
		)
	})

	v.subscribe("conversation", loggingEnabled, func(roomViewEventInfo) {
		msg := i18n.Local("This room is now publicly logged, meaning that everything you and the others in the room say or do can be made public on a website.")
		v.conv.displayWarningMessage(msg)
	})

	v.subscribe("conversation", loggingDisabled, func(roomViewEventInfo) {
		msg := i18n.Local("This room is no longer publicly logged.")
		v.conv.displayWarningMessage(msg)
	})

	return c
}

type mucStyleTags struct {
	table gtki.TextTagTable
}

func getTimestamp() string {
	return time.Now().Format("15:04:05")
}

func (v *roomViewConversation) getStyleTags(u *gtkUI) *mucStyleTags {
	if v.tags == nil {
		v.tags = v.newStyleTags(u)
	}
	return v.tags
}

func (v *roomViewConversation) newStyleTags(u *gtkUI) *mucStyleTags {
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

func (v *roomViewConversation) addNewLine() {
	buf, _ := v.roomChatTextView.GetBuffer()
	i := buf.GetEndIter()

	buf.Insert(i, "\n")
}

func (v *roomViewConversation) addTimestamp() {
	buf, _ := v.roomChatTextView.GetBuffer()
	i := buf.GetEndIter()
	t := fmt.Sprintf("[%s] ", getTimestamp())

	buf.InsertWithTagByName(i, t, "timestampText")
}

func (v *roomViewConversation) addTextToChatTextUsingTagID(text string, tagName string) {
	buf, _ := v.roomChatTextView.GetBuffer()
	i := buf.GetEndIter()

	if len(tagName) > 0 {
		buf.InsertWithTagByName(i, text, tagName)
	} else {
		buf.Insert(i, text)
	}
}

func (v *roomViewConversation) addLineTextToChatTextUsingTagID(text string, tagName string) {
	v.addTimestamp()
	v.addTextToChatTextUsingTagID(text, tagName)
	v.addNewLine()
}

// displayNotificationWhenOccupantJoinedRoom MUST be called from the UI thread
func (v *roomViewConversation) displayNotificationWhenOccupantJoinedRoom(nickname string) {
	text := fmt.Sprintf("%s joined the room", nickname)
	v.addLineTextToChatTextUsingTagID(text, "joinedRoomText")
}

// displayNotificationWhenOccupantLeftTheRoom MUST be called from the UI thread
func (v *roomViewConversation) displayNotificationWhenOccupantLeftTheRoom(nickname string) {
	text := fmt.Sprintf("%s left the room", nickname)
	v.addLineTextToChatTextUsingTagID(text, "leftRoomText")
}

// displayNewLiveMessage MUST be called from the UI thread
func (v *roomViewConversation) displayNewLiveMessage(nickname, subject, message string) {
	v.addTimestamp()
	nicknameText := fmt.Sprintf("%s: ", nickname)
	v.addTextToChatTextUsingTagID(nicknameText, "nicknameText")
	if len(subject) > 0 {
		subjectText := fmt.Sprintf("[%s] ", subject)
		v.addTextToChatTextUsingTagID(subjectText, "subjectText")
	}
	v.addTextToChatTextUsingTagID(message, "messageText")
	v.addNewLine()
}

// displayWarningMessage MUST be called from the UI thread
func (v *roomViewConversation) displayWarningMessage(message string) {
	v.addLineTextToChatTextUsingTagID(message, "warning")
}
