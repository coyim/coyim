package gui

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

const jidColumnIndex = 0

type mucRoomConfigListComponent struct {
	u                       *gtkUI
	list                    gtki.TreeView
	listModel               gtki.ListStore
	addButton, removeButton gtki.Button
	onAdd                   func()
	onNoItems               func()
}

func (u *gtkUI) newMUCRoomConfigListComponent(list gtki.TreeView, addButton, removeButton gtki.Button, onAdd, onNoItems func()) *mucRoomConfigListComponent {
	cl := &mucRoomConfigListComponent{
		u:            u,
		list:         list,
		addButton:    addButton,
		removeButton: removeButton,
		onAdd:        onAdd,
		onNoItems:    onNoItems,
	}

	cl.initListModel()
	cl.initDefaults()

	return cl
}

func (cl *mucRoomConfigListComponent) initListModel() {
	lm, _ := g.gtk.ListStoreNew(
		// jid
		glibi.TYPE_STRING,
	)

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

// onAddClicked MUST be called from the UI thread
func (cl *mucRoomConfigListComponent) onAddClicked() {
	if cl.onAdd != nil {
		cl.onAdd()
	}
}

// onRemoveClicked MUST be called from the UI thread
func (cl *mucRoomConfigListComponent) onRemoveClicked() {
	selection, _ := cl.list.GetSelection()
	selectedRows := selection.GetSelectedRows(cl.listModel)

	for i := len(selectedRows) - 1; i >= 0; i-- {
		iter, _ := cl.listModel.GetIter(selectedRows[i])
		cl.listModel.Remove(iter)
	}
}

// onSelectionChanged MUST be called from the UI thread
func (cl *mucRoomConfigListComponent) onSelectionChanged() {
	selection, _ := cl.list.GetSelection()
	selectedRows := selection.GetSelectedRows(cl.listModel)

	if len(selectedRows) > 0 {
		enableListWidget(cl.removeButton)
		return
	}
	disableListWidget(cl.removeButton)
}

// addListItems MUST be called from the UI thread
func (cl *mucRoomConfigListComponent) addListItems(jids []string) {
	for _, v := range jids {
		if cl.canBeAdded(v) {
			li := cl.listModel.Append()
			_ = cl.listModel.SetValue(li, jidColumnIndex, v)
		}
	}
}

func (cl *mucRoomConfigListComponent) canBeAdded(jid string) bool {
	iter, ok := cl.listModel.GetIterFirst()

	for ok {
		v, _ := cl.listModel.GetValue(iter, jidColumnIndex)
		s, _ := v.GetString()

		if s == jid {
			return false
		}

		ok = cl.listModel.IterNext(iter)
	}

	return true
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
