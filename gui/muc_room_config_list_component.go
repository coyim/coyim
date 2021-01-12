package gui

import (
	"errors"

	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomConfigListComponent struct {
	u                       *gtkUI
	list                    gtki.TreeView
	listModel               gtki.TreeStore
	listModelColumns        []glibi.Type
	addButton, removeButton gtki.Button
	onAdd                   func()
	items                   [][]string
}

func (u *gtkUI) newMUCRoomConfigListComponent(list gtki.TreeView, listModelColumns []glibi.Type, addButton, removeButton gtki.Button, onAdd func()) *mucRoomConfigListComponent {
	cl := &mucRoomConfigListComponent{
		u:                u,
		list:             list,
		listModelColumns: listModelColumns,
		addButton:        addButton,
		removeButton:     removeButton,
		onAdd:            onAdd,
	}

	cl.initListModel()
	cl.initDefaults()

	return cl
}

func (cl *mucRoomConfigListComponent) initListModel() {
	lm, _ := g.gtk.TreeStoreNew(cl.listModelColumns...)

	cl.listModel = lm
	cl.list.SetModel(cl.listModel)
}

func (cl *mucRoomConfigListComponent) initDefaults() {
	enableListWidget(cl.addButton)
	disableListWidget(cl.removeButton)

	cl.addButton.Connect("clicked", cl.onAddClicked)
	cl.removeButton.Connect("clicked", cl.onRemoveClicked)

	selection, err := cl.list.GetSelection()
	if err == nil {
		selection.Connect("changed", cl.onSelectionChanged)
	}
}

func (cl *mucRoomConfigListComponent) onAddClicked() {
	if cl.onAdd != nil {
		cl.onAdd()
	}
}

func (cl *mucRoomConfigListComponent) onRemoveClicked() {
	if iter, err := cl.getSelectedRow(); err == nil {
		cl.listModel.Remove(iter)
	}
}

func (cl *mucRoomConfigListComponent) onSelectionChanged() {
	if _, err := cl.getSelectedRow(); err != nil {
		disableListWidget(cl.removeButton)
	} else {
		enableListWidget(cl.removeButton)
	}
}

func (cl *mucRoomConfigListComponent) getSelectedRow() (gtki.TreeIter, error) {
	selection, err := cl.list.GetSelection()
	if err != nil {
		return nil, err
	}

	_, iter, ok := selection.GetSelected()
	if !ok {
		return nil, errors.New("no row selected")
	}

	return iter, nil
}

func (cl *mucRoomConfigListComponent) addListItems(items [][]string) {
	for _, it := range items {
		if len(items) > 0 {
			cl.addListItem(it...)
		}
	}
}

func (cl *mucRoomConfigListComponent) addListItem(cells ...string) {
	li := cl.listModel.Append(nil)

	for i, c := range cells {
		_ = cl.listModel.SetValue(li, i, c)
	}

	cl.items = append(cl.items, cells)
}

func enableListWidget(w gtki.Widget) {
	setListWidgetSensitive(w, true)
}

func disableListWidget(w gtki.Widget) {
	setListWidgetSensitive(w, false)
}

func setListWidgetSensitive(w gtki.Widget, v bool) {
	w.SetSensitive(v)
}
