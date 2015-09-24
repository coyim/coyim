package main

import (
	"github.com/twstrike/go-gtk/gtk"
)

const DESTROY_SIGNAL = "destroy"

type textBox struct {
	view *gtk.TextView
	buf  *gtk.TextBuffer
	iter *gtk.TextIter
}

func main() {
	gtk.Init(nil)
	window := startNewWindow("Conversation with Fan")

	conv_info := newReadOnlyTextBox()
	conv_info.write("Started at 10:55:33")
	conv_info.write("\nConversation has been active for 00:11:23")
	conv_info.write("\nConversation will be logged")
	conv_info.write("\nOTR enabled")

	window.Add(conv_info.view)

	window.ShowAll()
	gtk.Main()
}

func startNewWindow(title string) *gtk.Window {
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle(title)
	window.Connect(DESTROY_SIGNAL, gtk.MainQuit)
	window.SetSizeRequest(600, 480)
	return window
}

func newReadOnlyTextBox() textBox {
	view := gtk.NewTextView()
	view.SetEditable(false)
	view.SetCursorVisible(false)

	buf := view.GetBuffer()

	var iter gtk.TextIter
	buf.GetStartIter(&iter)

	return textBox{view, buf, &iter}
}

func (textBox textBox) write(text string) {
	textBox.buf.Insert(textBox.iter, text)
}
