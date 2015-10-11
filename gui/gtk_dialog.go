package gui

import "github.com/twstrike/go-gtk/gtk"

type dialog struct {
	title    string
	position gtk.WindowPosition
	vbox     []createable
	id       string
}

func (wr *widgetRegistry) dialogDestroy(id string) {
	wr.reg[id].(*gtk.Dialog).Destroy()
}

func (wr *widgetRegistry) dialogShowAll(id string) {
	wr.reg[id].(*gtk.Dialog).ShowAll()
}

func (d dialog) create(reg *widgetRegistry) gtk.IWidget {
	dialog := gtk.NewDialog()
	dialog.SetTitle(d.title)
	dialog.SetPosition(d.position)
	vbox := dialog.GetVBox()

	for _, item := range d.vbox {
		vbox.Add(item.create(reg))
	}

	reg.register(d.id, dialog)

	return dialog
}
