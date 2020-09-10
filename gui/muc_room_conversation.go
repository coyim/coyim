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

type mucStyleTags struct {
	table gtki.TextTagTable
}

type messageType int

const (
	mtLeftRoom messageType = iota
	mtLiveMessage
	mtWarning
)

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
	_ = leftRoomTag.SetProperty("foreground", cset.warningForeground)
	_ = leftRoomTag.SetProperty("style", pangoi.STYLE_ITALIC)

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

func (v *roomViewConversation) showOccupantLeftRoom(nickname string) {
	v.showMessageInChatRoom(nickname, "", "left the room", mtLeftRoom)
}

func (v *roomViewConversation) showMessageInChatRoom(nickname, subject, message string, mt messageType) {
	buf, _ := v.roomChatTextView.GetBuffer()
	c := buf.GetCharCount()

	t := fmt.Sprintf("[%s]", getTimestamp())

	n := fmt.Sprintf("%s", nickname)
	if mt == mtLiveMessage {
		n = fmt.Sprintf("%s:", nickname)
	}

	switch mt {
	case mtWarning:
		txt := i18n.Localf("%s %s", t, message)
		v.addTextToChat(txt)
		c = v.applyTagByNameAndOffset(buf, "timestampText", t, c)
		c = v.applyTagByNameAndOffset(buf, "warning", t, c)
	case mtLeftRoom:
		txt := i18n.Localf("%s %s %s", t, n, message)
		v.addTextToChat(txt)
		c = v.applyTagByNameAndOffset(buf, "timestampText", t, c)
		c = v.applyTagByNameAndOffset(buf, "leftRoomText", n, c)
		c = v.applyTagByNameAndOffset(buf, "leftRoomText", message, c)
	case mtLiveMessage:
		txt := i18n.Localf("%s %s %s ", t, n, message)
		v.addTextToChat(txt)
		c = v.applyTagByNameAndOffset(buf, "timestampText", t, c)
		c = v.applyTagByNameAndOffset(buf, "nicknameText", n, c)
		c = v.applyTagByNameAndOffset(buf, "messageText", message, c)
	default:
		v.log.WithField("messageType", mt).Debug("message type not controlled")
		return
	}

	// TODO: This call should be better if we connect this to the signal
	// "size-allocate" of the roomChatTextView, for now it's ok here
	// because this is the only function that add text to the textview
	scrollToBottom(v.roomChatScrollView)
}

func (v *roomViewConversation) applyTagByNameAndOffset(b gtki.TextBuffer, tagName, text string, initialPos int) int {
	finalPos := initialPos + 1 + len(text)
	beginIter := b.GetIterAtOffset(initialPos)
	endIter := b.GetIterAtOffset(finalPos)
	b.ApplyTagByName(tagName, beginIter, endIter)

	return finalPos
}
