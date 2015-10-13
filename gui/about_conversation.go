package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

const DESTROY_SIGNAL = "destroy"

type aboutConversationWin struct {
	win *gtk.Window
	txt textBox
}

type textBox struct {
	view *gtk.TextView
	buf  *gtk.TextBuffer
	iter *gtk.TextIter
}

func newAboutConversationWindow(title string) aboutConversationWin {
	window, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle(title)
	window.Connect(DESTROY_SIGNAL, gtk.MainQuit)
	window.SetSizeRequest(600, 480)

	return aboutConversationWin{window, newReadOnlyTextBox()}
}

func (about aboutConversationWin) add(box textBox) {
	about.win.Add(box.view)
}

func (about aboutConversationWin) write(text string) {
	about.txt.write(text)
}

func (about aboutConversationWin) render() {
	about.add(about.txt)
	about.win.ShowAll()
}

func newReadOnlyTextBox() textBox {
	view, _ := gtk.TextViewNew()
	view.SetEditable(false)
	view.SetCursorVisible(false)

	buf, _ := view.GetBuffer()
	iter := buf.GetStartIter()

	return textBox{view, buf, iter}
}

func (box textBox) write(text string) {
	box.buf.Insert(box.iter, text)
}
