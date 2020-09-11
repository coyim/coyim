package gui

import (
	"fmt"
	"time"

	"github.com/coyim/coyim/coylog"
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
	_ = leftRoomTag.SetProperty("foreground", "#731629")
	_ = leftRoomTag.SetProperty("style", pangoi.STYLE_ITALIC)

	joinedRoomTag, _ := g.gtk.TextTagNew("joinedRoomText")
	_ = joinedRoomTag.SetProperty("foreground", cset.mucSomeoneJoinedForeground)
	_ = joinedRoomTag.SetProperty("style", pangoi.STYLE_ITALIC)

	timestampTag, _ := g.gtk.TextTagNew("timestampText")
	_ = timestampTag.SetProperty("foreground", "#AAB7B8")
	_ = timestampTag.SetProperty("style", pangoi.STYLE_NORMAL)

	nicknameTag, _ := g.gtk.TextTagNew("nicknameText")
	_ = nicknameTag.SetProperty("foreground", "#395BA3")
	_ = nicknameTag.SetProperty("style", pangoi.STYLE_NORMAL)

	messageTag, _ := g.gtk.TextTagNew("messageText")
	_ = messageTag.SetProperty("foreground", "#000000")
	_ = messageTag.SetProperty("style", pangoi.STYLE_NORMAL)

	t.table.Add(warningTag)
	t.table.Add(leftRoomTag)
	t.table.Add(joinedRoomTag)
	t.table.Add(timestampTag)
	t.table.Add(nicknameTag)
	t.table.Add(messageTag)

	return t
}

func (t *mucStyleTags) createTextBuffer() gtki.TextBuffer {
	buf, _ := g.gtk.TextBufferNew(t.table)
	return buf
}

func (u *gtkUI) newRoomViewConversation() *roomViewConversation {
	c := &roomViewConversation{}

	builder := newBuilder("MUCRoomConversation")
	panicOnDevError(builder.bindObjects(c))

	t := c.getStyleTags(u)
	c.roomChatTextView.SetBuffer(t.createTextBuffer())

	return c
}

func (v *roomViewConversation) addTextToChat(text string) {
	buf, _ := v.roomChatTextView.GetBuffer()
	i := buf.GetEndIter()

	buf.Insert(i, fmt.Sprintf("%s\n", text))
}

func (v *roomViewConversation) displayNotificationWhenOccupantJoinedRoom(nickname string) {
	text := fmt.Sprintf("%s joined the room", nickname)
	v.addLineToChatTextUsingTagID(text, "joinedRoomText")
}

func (v *roomViewConversation) showOccupantLeftRoom(nickname string) {
	text := fmt.Sprintf("%s left the room", nickname)
	v.addLineToChatTextUsingTagID(text, "leftRoomText")
}

func (v *roomViewConversation) showLiveMessageInTheRoom(nickname, subject, message string) {
	text := fmt.Sprintf("%s: %s", nickname, message)
	if len(subject) > 0 {
		text = fmt.Sprintf("%s: [%s] %s", nickname, subject, message)
	}
	v.addLineToChatTextUsingTagID(text, "messageText")
}

func (v *roomViewConversation) addLineToChatText(timestamp, text string) {
	buf, _ := v.roomChatTextView.GetBuffer()
	i := buf.GetEndIter()

	newtext := fmt.Sprintf("%s %s\n", timestamp, text)
	buf.Insert(i, newtext)
}

func (v *roomViewConversation) addLineToChatTextUsingTagID(text string, tag string) {
	buf, _ := v.roomChatTextView.GetBuffer()

	charCount := buf.GetCharCount()

	t := fmt.Sprintf("[%s]", getTimestamp())
	v.addLineToChatText(t, text)

	oldIterEnd := buf.GetIterAtOffset(charCount)
	offsetTimestamp := buf.GetIterAtOffset(charCount + len(t) + 1)
	newIterEnd := buf.GetEndIter()

	buf.ApplyTagByName("timestampText", oldIterEnd, offsetTimestamp)
	buf.ApplyTagByName(tag, offsetTimestamp, newIterEnd)
}
