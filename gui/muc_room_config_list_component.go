package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

const jidColumnIndex = 0

type mucRoomConfigListComponent struct {
	u                 *gtkUI
	list              gtki.TreeView
	listModel         gtki.ListStore
	addButton         gtki.Button
	removeButton      gtki.Button
	removeButtonLabel gtki.Label
	onAdd             func()
	onNoItems         func()
}

func (u *gtkUI) newMUCRoomConfigListComponent(list gtki.TreeView, addButton, removeButton gtki.Button, removeButtonLabel gtki.Label, onAdd, onNoItems func()) *mucRoomConfigListComponent {
	cl := &mucRoomConfigListComponent{
		u:                 u,
		list:              list,
		addButton:         addButton,
		removeButton:      removeButton,
		removeButtonLabel: removeButtonLabel,
		onAdd:             onAdd,
		onNoItems:         onNoItems,
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
	cl.addButton.SetSensitive(true)
	cl.removeButton.SetSensitive(false)

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

	currentIndex := len(selectedRows) - 1
	for currentIndex >= 0 {
		if iter, err := cl.listModel.GetIter(selectedRows[currentIndex]); err == nil {
			cl.listModel.Remove(iter)
		}
		currentIndex -= 1
	}

	if _, ok := cl.listModel.GetIterFirst(); !ok {
		cl.onNoItems()
	}
}

// onSelectionChanged MUST be called from the UI thread
func (cl *mucRoomConfigListComponent) onSelectionChanged() {
	selection, _ := cl.list.GetSelection()
	selectedRows := selection.GetSelectedRows(cl.listModel)

	totalItems := len(selectedRows)
	cl.removeButton.SetSensitive(totalItems > 0)

	l := i18n.Local("Remove")
	if totalItems > 1 {
		l = i18n.Local("Remove all")
	}
	cl.removeButtonLabel.SetText(l)
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
