package gui

import "github.com/coyim/gotk3adapter/gtki"

func (c *roomViewConversation) getMessageTextBuffer() gtki.TextBuffer {
	b, _ := c.messageTextView.GetBuffer()
	return b
}

func (c *roomViewConversation) getTextBuffer() gtki.TextBuffer {
	b, _ := c.chatTextView.GetBuffer()
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
