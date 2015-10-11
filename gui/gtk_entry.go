package gui

import "github.com/twstrike/go-gtk/gtk"

type entry struct {
	text       string
	editable   bool
	visibility bool
	id         string
}

func (wr *widgetRegistry) getText(id string) string {
	w := wr.reg[id]

	switch w := w.(type) {
	case *gtk.Entry:
		return w.GetText()
	}

	return ""
}

func (e entry) create(reg *widgetRegistry) gtk.IWidget {
	entry := gtk.NewEntry()
	entry.SetText(e.text)
	entry.SetEditable(e.editable)
	entry.SetVisibility(e.visibility)
	reg.register(e.id, entry)

	return entry
}
