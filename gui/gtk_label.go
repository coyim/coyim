package gui

import "github.com/twstrike/gotk3/gtk"

type label struct {
	text string
}

func (l label) create(reg *widgetRegistry) (gtk.IWidget, error) {
	return gtk.LabelNew(l.text)
}
