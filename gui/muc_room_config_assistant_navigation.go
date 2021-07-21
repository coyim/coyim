package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

const noRowIndex = -1

type roomConfigAssistantNavigation struct {
	assistant *roomConfigAssistant

	content    gtki.Box     `gtk-widget:"room-config-assistant-navigation-content"`
	navigation gtki.ListBox `gtk-widget:"room-config-assistant-navigation-list"`

	items []*roomConfigAssistantNavigationItem
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
	for _, p := range rcn.assistant.allPages() {
		itm := rcn.newRoomConfigAssistantNavigationItem(p)
		rcn.items = append(rcn.items, itm)
		rcn.navigation.Add(itm.row)
		rcn.navigation.Add(itm.divider)
	}
}

// onRowSelected MUST be called from the UI thread
func (rcn *roomConfigAssistantNavigation) onRowSelected(_ gtki.ListBox, r gtki.ListBoxRow) {
	// Every time a row is selected, we check if it's not a divider or a normal item.
	// Normal items will always be in even positions, and rows dividers will be in odd positions.
	if idx := getListBoxRowIndex(r); idx%2 == 0 {
		rcn.assistant.updateAssistantPage(mucRoomConfigPageID(idx - (idx / 2)))
	}
}

// selectPageByIndex MUST be called from the UI thread
func (rcn *roomConfigAssistantNavigation) selectPageByIndex(pageID mucRoomConfigPageID) {
	// Since we have a divider for each navigation item, every time we want to select
	// a page by its index, we should calculate the right index to avoid a wrong behavior
	row := rcn.navigation.GetRowAtIndex(int(pageID) * 2)
	rowIndex := getListBoxRowIndex(row)
	currentRowIndex := getListBoxRowIndex(rcn.navigation.GetSelectedRow())

	if rowIndex != noRowIndex && rowIndex != currentRowIndex {
		rcn.navigation.SelectRow(row)
	}
}

func getListBoxRowIndex(r gtki.ListBoxRow) int {
	if r != nil {
		return r.GetIndex()
	}
	return noRowIndex
}
