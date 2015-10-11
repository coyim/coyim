package gui

import "github.com/twstrike/go-gtk/gtk"

type label struct {
	text string
}

func (l label) create(reg *widgetRegistry) gtk.IWidget {
	return gtk.NewLabel(l.text)
}
