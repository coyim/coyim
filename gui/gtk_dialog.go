package gui

import "github.com/gotk3/gotk3/gtk"

type dialog struct {
	title       string
	position    gtk.WindowPosition
	content     []creatable
	id          string
	defaultSize []int
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

	if d.defaultSize != nil {
		dialog.SetDefaultSize(d.defaultSize[0], d.defaultSize[1])
	}

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

func (d dialog) createWithDefault(wr *widgetRegistry, buttonId string) (gtk.IWidget, error) {
	dialog, e := d.create(wr)
	if e != nil {
		return nil, e
	}

	button := wr.reg[buttonId].(*gtk.Button)
	button.SetCanDefault(true)
	dialog.(*gtk.Dialog).SetDefault(button)

	for _, item := range d.content {
		widget := wr.reg[item.getId()]
		switch w := widget.(type) {
		case *gtk.Entry:
			w.SetActivatesDefault(true)
		}
	}

	return dialog, nil
}
