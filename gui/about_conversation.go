package gui

import "github.com/twstrike/gotk3adapter/gtki"

const destroySignal = "destroy"

type aboutConversationWin struct {
	win gtki.Window
	txt textBox
}

type textBox struct {
	view gtki.TextView
	buf  gtki.TextBuffer
	iter gtki.TextIter
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
	view, _ := g.gtk.TextViewNew()
	view.SetEditable(false)
	view.SetCursorVisible(false)

	buf, _ := view.GetBuffer()
	iter := buf.GetStartIter()

	return textBox{view, buf, iter}
}

func (box textBox) write(text string) {
	box.buf.Insert(box.iter, text)
}
