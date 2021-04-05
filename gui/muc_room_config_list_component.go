package gui

import (
	"errors"

	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

const jidColumnIndex = 0

type mucRoomConfigListComponent struct {
	u                       *gtkUI
	list                    gtki.TreeView
	listModel               gtki.TreeStore
	addButton, removeButton gtki.Button
	onAdd                   func()
	onNoItems               func()
	jidList                 []string
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
	lm, _ := g.gtk.TreeStoreNew(
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
	iter, err := cl.getSelectedRow()
	if err != nil {
		return
	}

	jidValue, _ := cl.listModel.GetValue(iter, jidColumnIndex)
	selectedJid, _ := jidValue.GetString()

	copy := []string{}
	for _, jid := range cl.jidList {
		if jid != selectedJid {
			copy = append(copy, jid)
		}
	}

	cl.jidList = copy
	cl.redraw()

	if len(cl.jidList) == 0 && cl.onNoItems != nil {
		cl.onNoItems()
	}
}

// onSelectionChanged MUST be called from the UI thread
func (cl *mucRoomConfigListComponent) onSelectionChanged() {
	if _, err := cl.getSelectedRow(); err != nil {
		disableListWidget(cl.removeButton)
	} else {
		enableListWidget(cl.removeButton)
	}
}

// getSelectedRow MUST be called from the UI thread
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

// addListItems MUST be called from the UI thread
func (cl *mucRoomConfigListComponent) addListItems(jids []string) {
	for _, v := range jids {
		if cl.canBeAdded(v) {
			cl.jidList = append(cl.jidList, v)
		}
	}
	cl.redraw()
}

func (cl *mucRoomConfigListComponent) canBeAdded(jid string) bool {
	for _, cj := range cl.jidList {
		if cj == jid {
			return false
		}
	}
	return true
}

// redraw MUST be called from the UI thread
func (cl *mucRoomConfigListComponent) redraw() {
	cl.listModel.Clear()

	for _, jid := range cl.jidList {
		li := cl.listModel.Append(nil)
		_ = cl.listModel.SetValue(li, jidColumnIndex, jid)
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
