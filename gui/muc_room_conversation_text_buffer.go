package gui

import "github.com/coyim/gotk3adapter/gtki"

// initTagsAndTextBuffers MUST be called from the UI thread
func (c *roomViewConversation) initTagsAndTextBuffers() {
	c.tags = c.newConversationTags()

	cb, _ := g.gtk.TextBufferNew(c.tags.table)
	c.chatTextView.SetBuffer(cb)

	mb, _ := g.gtk.TextBufferNew(nil)
	c.messageTextView.SetBuffer(mb)
}

// getMessageTextBuffer MUST be called from the UI thread
func (c *roomViewConversation) getMessageTextBuffer() gtki.TextBuffer {
	b, _ := c.messageTextView.GetBuffer()
	return b
}

// getTextBuffer MUST be called from the UI thread
func (c *roomViewConversation) getTextBuffer() gtki.TextBuffer {
	b, _ := c.chatTextView.GetBuffer()
	return b
}

const newLineText = "\n"

// addNewLine MUST be called from the UI thread
func (c *roomViewConversation) addNewLine() {
	if c.hasContentDisplayed() {
		c.addText(newLineText)
	}
}

// hasContentDisplayed MUST be called from the UI thread
func (c *roomViewConversation) hasContentDisplayed() bool {
	b := c.getTextBuffer()
	return b.GetText(b.GetStartIter(), b.GetEndIter(), false) != ""
}

// addText MUST be called from the UI thread
func (c *roomViewConversation) addText(text string) {
	b := c.getTextBuffer()
	b.Insert(b.GetEndIter(), text)
}

// addTextWithTag MUST be called from the UI thread
func (c *roomViewConversation) addTextWithTag(text string, tag conversationTag) {
	b := c.getTextBuffer()
	b.InsertWithTagByName(b.GetEndIter(), text, string(tag))
}

// clearBuffer MUST be called from the UI thread
func (c *roomViewConversation) clearBuffer() {
	b := c.getTextBuffer()
	b.Delete(b.GetStartIter(), b.GetEndIter())
}
