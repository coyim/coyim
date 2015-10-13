package gui

import "github.com/gotk3/gotk3/gtk"

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
		t, _ := w.GetText()
		return t
	}

	return ""
}

func (e entry) create(reg *widgetRegistry) (gtk.IWidget, error) {
	entry, err := gtk.EntryNew()
	if err != nil {
		return nil, err
	}

	entry.SetText(e.text)
	entry.SetEditable(e.editable)
	entry.SetVisibility(e.visibility)
	reg.register(e.id, entry)

	return entry, nil
}
