package gui

import (
	"fmt"
	"time"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3adapter/pangoi"
)

type roomViewConversation struct {
	view             gtki.Box      `gtk-widget:"roomConversation"`
	roomChatTextView gtki.TextView `gtk-widget:"roomChatTextView"`

	tags *mucStyleTags
}

type mucStyleTags struct {
	table gtki.TextTagTable
}

func getTimestamp() string {
	return time.Now().Format("15:04:05")
}

func (v *roomViewConversation) getStyleTags() *mucStyleTags {
	if v.tags == nil {
		v.tags = v.newStyleTags()
	}
	return v.tags
}

func (v *roomViewConversation) newStyleTags() *mucStyleTags {
	// TODO: for now we are using a default styles, but we can improve it
	// if we define a structure with a predefined colors pallete based on kind
	// of messages to show like entering a room, leaving the room, incoming
	// message, etc
	t := &mucStyleTags{}

	t.table, _ = g.gtk.TextTagTableNew()

	leftRoomTag, _ := g.gtk.TextTagNew("leftRoomText")
	_ = leftRoomTag.SetProperty("foreground", "#EE0000")
	_ = leftRoomTag.SetProperty("style", pangoi.STYLE_ITALIC)

	timestampTag, _ := g.gtk.TextTagNew("timestampText")
	_ = timestampTag.SetProperty("foreground", "#AAB7B8")
	_ = timestampTag.SetProperty("style", pangoi.STYLE_NORMAL)

	cset := v.u.currentColorSet()

	warningTag, _ := g.gtk.TextTagNew("warning")
	_ = warningTag.SetProperty("foreground", cset.warningForeground)

	t.table.Add(leftRoomTag)
	t.table.Add(timestampTag)
	t.table.Add(warningTag)

	return t
}

func (t *mucStyleTags) createTextBuffer() gtki.TextBuffer {
	buf, _ := g.gtk.TextBufferNew(t.table)
	return buf
}

func newRoomViewConversation() *roomViewConversation {
	c := &roomViewConversation{}

	builder := newBuilder("MUCRoomConversation")
	panicOnDevError(builder.bindObjects(c))

	t := c.getStyleTags()
	c.roomChatTextView.SetBuffer(t.createTextBuffer())

	return c
}

func (v *roomViewConversation) addLineToChatText(timestamp, text string) {
	buf, _ := v.roomChatTextView.GetBuffer()
	i := buf.GetEndIter()

	buf.Insert(i, fmt.Sprintf("%s %s\n", timestamp, text))
}

func (v *roomViewConversation) addLineToChatTextUsingTagID(text string, tag string) {
	buf, _ := v.roomChatTextView.GetBuffer()

	charCount := buf.GetCharCount()

	t := fmt.Sprintf("[%s] ", getTimestamp())
	v.addLineToChatText(t, text)

	oldIterEnd := buf.GetIterAtOffset(charCount)
	offsetTimestamp := buf.GetIterAtOffset(charCount + len(t))
	offsetText := buf.GetIterAtOffset(charCount + len(t) + 1)
	newIterEnd := buf.GetEndIter()

	buf.ApplyTagByName("timestampText", oldIterEnd, offsetTimestamp)
	buf.ApplyTagByName(tag, offsetText, newIterEnd)
}

func (v *roomViewConversation) showOccupantLeftRoom(nickname jid.Resource) {
	msg := i18n.Localf("%s left the room", nickname)
	v.addLineToChatTextUsingTagID(msg, "leftRoomText")
}
