package gui

import "github.com/gotk3/gotk3/gtk"

type label struct {
	text      string
	wrapLines bool
}

func (l label) getID() string {
	return ""
}

func (l label) create(reg *widgetRegistry) (gtk.IWidget, error) {
	label, e := gtk.LabelNew(l.text)

	if l.wrapLines {
		label.SetLineWrap(true)
	}

	return label, e
}
