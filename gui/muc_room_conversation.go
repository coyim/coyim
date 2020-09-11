package gui

import (
	"fmt"
	"strings"
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

type tagUtil struct {
	tags map[string]string
	keys []string
}

func (v *roomViewConversation) newTagUtil() *tagUtil {
	tu := &tagUtil{}
	tu.tags = make(map[string]string)
	return tu
}

func (tu *tagUtil) add(tagName, text string) {
	tu.keys = append(tu.keys, tagName)
	tu.tags[tagName] = text
}

func (tu *tagUtil) addFromTagUtil(fromtu *tagUtil) {
	for _, key := range fromtu.keys {
		tu.add(key, fromtu.tags[key])
	}
}

func (tu *tagUtil) getText() string {
	var text string
	for _, key := range tu.keys {
		text = text + tu.tags[key]
	}
	return text
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

	messageTag, _ := g.gtk.TextTagNew("messageText")
	_ = messageTag.SetProperty("foreground", cset.mucMessageForeground)
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

func (v *roomViewConversation) addLineToChatText(text string) {
	buf, _ := v.roomChatTextView.GetBuffer()
	i := buf.GetEndIter()

	newtext := fmt.Sprintf("%s\n", text)
	buf.Insert(i, newtext)
}

func (v *roomViewConversation) addLineToChatTextUsingTagID(text string, tag string) {
	buf, _ := v.roomChatTextView.GetBuffer()

	charCount := buf.GetCharCount()

	t := fmt.Sprintf("[%s]", getTimestamp())
	v.addLineToChatText(text)

	oldIterEnd := buf.GetIterAtOffset(charCount)
	offsetTimestamp := buf.GetIterAtOffset(charCount + len(t) + 1)
	newIterEnd := buf.GetEndIter()

	buf.ApplyTagByName("timestampText", oldIterEnd, offsetTimestamp)
	buf.ApplyTagByName(tag, offsetTimestamp, newIterEnd)
}

func (v *roomViewConversation) addLineToChatTextUsingTagUtils(tu *tagUtil) {
	newtu := v.newTagUtil()
	t := fmt.Sprintf("[%s] ", getTimestamp())
	newtu.add("timestampText", t)
	newtu.addFromTagUtil(tu)
	text := newtu.getText()

	buf, _ := v.roomChatTextView.GetBuffer()
	charCount := buf.GetCharCount()

	v.addLineToChatText(text)

	for _, tagName := range newtu.keys {
		tagText := newtu.tags[tagName]
		pos := strings.Index(text, tagText)
		iterFrom := buf.GetIterAtOffset(charCount + pos)
		iterTo := buf.GetIterAtOffset(charCount + pos + len(tagText))
		buf.ApplyTagByName(tagName, iterFrom, iterTo)
	}
}

// displayNotificationWhenOccupantJoinedRoom MUST be called from the UI thread
func (v *roomViewConversation) displayNotificationWhenOccupantJoinedRoom(nickname string) {
	text := fmt.Sprintf("%s joined the room", nickname)
	v.addLineToChatTextUsingTagID(text, "joinedRoomText")
}

// displayNotificationWhenOccupantLeftTheRoom MUST be called from the UI thread
func (v *roomViewConversation) displayNotificationWhenOccupantLeftTheRoom(nickname string) {
	text := fmt.Sprintf("%s left the room", nickname)
	v.addLineToChatTextUsingTagID(text, "leftRoomText")
}

// displayNewLiveMessage MUST be called from the UI thread
func (v *roomViewConversation) displayNewLiveMessage(nickname, subject, message string) {
	tu := v.newTagUtil()
	nicknameText := fmt.Sprintf("%s: ", nickname)
	tu.add("nicknameText", nicknameText)
	if len(subject) > 0 {
		subjectText := fmt.Sprintf("[%s] ", subject)
		tu.add("messageText", subjectText)
	}
	tu.add("messageText", message)
	v.addLineToChatTextUsingTagUtils(tu)
}

// displayWarningMessage MUST be called from the UI thread
func (v *roomViewConversation) displayWarningMessage(message string) {
	v.addLineToChatTextUsingTagID(message, "warning")
}
