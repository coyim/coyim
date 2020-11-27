package gui

import (
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
}

func (cl *mucRoomConfigListComponent) onAddClicked() {
	cl.onAdd()
}

func (cl *mucRoomConfigListComponent) onRemoveClicked() {}

func (cl *mucRoomConfigListComponent) addListItem(cells ...string) {
	li := cl.listModel.Append(nil)

	for i, c := range cells {
		_ = cl.listModel.SetValue(li, i, c)
	}
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
