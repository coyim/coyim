package gui

import "github.com/twstrike/gotk3/gtk"

type button struct {
	text      string
	onClicked func()
	id        string
}

func (b button) create(reg *widgetRegistry) (gtk.IWidget, error) {
	button, e := gtk.ButtonNewWithLabel(b.text)
	if e != nil {
		return nil, e
	}

	if b.onClicked != nil {
		button.Connect("clicked", b.onClicked)
	}

	reg.register(b.id, button)

	return button, e
}
