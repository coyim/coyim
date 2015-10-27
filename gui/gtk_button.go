package gui

import "github.com/gotk3/gotk3/gtk"

type button struct {
	text       string
	onClicked  func()
	onActivate func()
	id         string
	focused    bool
}

func (b button) getId() string {
	return b.id
}

func (b button) create(reg *widgetRegistry) (gtk.IWidget, error) {
	button, e := gtk.ButtonNewWithLabel(b.text)
	if e != nil {
		return nil, e
	}

	if b.onClicked != nil {
		button.Connect("clicked", b.onClicked)
	}

	if b.onActivate != nil {
		button.Connect("activate", b.onActivate)
	}

	if b.focused {
		button.GrabFocus()
	}

	reg.register(b.id, button)

	return button, e
}
