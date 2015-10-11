package gui

import "github.com/twstrike/go-gtk/gtk"

type widgetRegistry struct {
	reg map[string]gtk.IWidget
}

func createWidgetRegistry() *widgetRegistry {
	v := widgetRegistry{
		reg: make(map[string]gtk.IWidget),
	}
	return &v
}

func (wr *widgetRegistry) register(id string, w gtk.IWidget) {
	if id != "" {
		wr.reg[id] = w
	}
}

type createable interface {
	create(reg *widgetRegistry) gtk.IWidget
}
