package gui

import "github.com/twstrike/go-gtk/gtk"

type button struct {
	text      string
	onClicked func()
	id        string
}

func (b button) create(reg *widgetRegistry) gtk.IWidget {
	button := gtk.NewButtonWithLabel(b.text)
	if b.onClicked != nil {
		button.Connect("clicked", func() { b.onClicked() })
	}
	reg.register(b.id, button)
	return button
}
