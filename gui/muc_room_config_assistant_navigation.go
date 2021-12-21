package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

const noRowIndex = -1

type roomConfigAssistantNavigation struct {
	assistant *roomConfigAssistant

	content    gtki.Box     `gtk-widget:"room-config-assistant-navigation-content"`
	navigation gtki.ListBox `gtk-widget:"room-config-assistant-navigation-list"`

	items   []*roomConfigAssistantNavigationItem
	iconSet navigationIconMapper
}

func (rc *roomConfigAssistant) newRoomConfigAssistantNavigation() *roomConfigAssistantNavigation {
	rcn := &roomConfigAssistantNavigation{
		assistant: rc,
	}

	rcn.initBuilder()
	rcn.initNavigationItems()

	return rcn
}

func (rcn *roomConfigAssistantNavigation) initBuilder() {
	b := newBuilder("MUCRoomConfigAssistantNavigation")
	panicOnDevError(b.bindObjects(rcn))

	b.ConnectSignals(map[string]interface{}{
		"on_row_selected": rcn.onRowSelected,
	})
}

func (rcn *roomConfigAssistantNavigation) initNavigationItems() {
	rcn.iconSet = rcn.assistant.u.currentIconSet()
	for _, p := range rcn.assistant.allPages() {
		itm := rcn.newRoomConfigAssistantNavigationItem(p)
		rcn.items = append(rcn.items, itm)
		rcn.navigation.Add(itm.row)
	}
}

// disableNavigation MUST be called from the UI thread
func (rcn *roomConfigAssistantNavigation) disableNavigation() {
	rcn.forEachItem(func(itm *roomConfigAssistantNavigationItem) {
		if rcn.assistant.currentPageIndex != itm.pageID() {
			itm.disable()
		}
	})
}

// enableNavigation MUST be called from the UI thread
func (rcn *roomConfigAssistantNavigation) enableNavigation() {
	rcn.forEachItem(func(itm *roomConfigAssistantNavigationItem) {
		itm.enable()
	})
}

// forEachItem MUST be called from the UI thread
func (rcn *roomConfigAssistantNavigation) forEachItem(fn func(*roomConfigAssistantNavigationItem)) {
	for _, itm := range rcn.items {
		fn(itm)
	}
}

// onRowSelected MUST be called from the UI thread
func (rcn *roomConfigAssistantNavigation) onRowSelected(_ gtki.ListBox, r gtki.ListBoxRow) {
	rcn.assistant.updateAssistantPage(mucRoomConfigPageID(getListBoxRowIndex(r)))
}

// selectPageByIndex MUST be called from the UI thread
func (rcn *roomConfigAssistantNavigation) selectPageByIndex(pageID mucRoomConfigPageID) {
	row := rcn.navigation.GetRowAtIndex(int(pageID))
	rowIndex := getListBoxRowIndex(row)
	currentRowIndex := getListBoxRowIndex(rcn.navigation.GetSelectedRow())

	if rowIndex != noRowIndex && rowIndex != currentRowIndex {
		rcn.navigation.SelectRow(row)
		doALittleBitLater(row.GrabFocus)
	}

}

// doALittleBitLater will execute the function given when all current
// tasks in the UI thread has been managed.
func doALittleBitLater(f func()) {
	doInUIThread(f)
}

func getListBoxRowIndex(r gtki.ListBoxRow) int {
	if r != nil {
		return r.GetIndex()
	}
	return noRowIndex
}
