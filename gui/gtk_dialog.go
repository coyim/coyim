package gui

import "github.com/twstrike/gotk3/gtk"

type dialog struct {
	title    string
	position gtk.WindowPosition
	content  []createable
	id       string
}

func (wr *widgetRegistry) dialogDestroy(id string) {
	wr.reg[id].(*gtk.Dialog).Destroy()
}

func (wr *widgetRegistry) dialogShowAll(id string) {
	wr.reg[id].(*gtk.Dialog).ShowAll()
}

func (d dialog) create(reg *widgetRegistry) (gtk.IWidget, error) {
	dialog, e := gtk.DialogNew()
	if e != nil {
		return nil, e
	}

	dialog.SetTitle(d.title)
	dialog.SetPosition(d.position)

	//TODO: error
	content, _ := dialog.GetContentArea()

	for _, item := range d.content {
		//TODO: error
		i, _ := item.create(reg)
		content.Add(i)
	}

	reg.register(d.id, dialog)

	return dialog, nil
}
