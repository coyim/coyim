package gui

import "github.com/gotk3/gotk3/gtk"

type hbox struct {
	content   []creatable
	fromRight bool
}

func (h hbox) getId() string {
	return ""
}

func (h hbox) create(reg *widgetRegistry) (gtk.IWidget, error) {
	b, e := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1)
	if e != nil {
		return nil, e
	}

	b.SetHomogeneous(false)
	b.SetBorderWidth(5)

	for _, item := range h.content {
		i, _ := item.create(reg)
		if h.fromRight {
			b.PackEnd(i, false, false, 2)
		} else {
			b.PackStart(i, false, false, 2)
		}
	}

	return b, nil
}
