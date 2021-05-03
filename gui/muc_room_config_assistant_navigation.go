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
	rcn.initNavigationDividers()

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
		itm := rcn.newRoomConfigAssistantNavigationItem(p.pageTitle())
		rcn.items = append(rcn.items, itm)
		rcn.navigation.Add(itm.row)
		rcn.navigation.Add(itm.divider)
	}
}

func (rcn *roomConfigAssistantNavigation) initNavigationDividers() {
	for idx := range rcn.dividerIndexes() {
		divider := rcn.navigation.GetRowAtIndex(idx)
		if divider != nil {
			divider.SetProperty("activatable", false)
			divider.SetProperty("selectable", false)
		}
	}
}

// onRowSelected MUST be called from the UI thread
func (rcn *roomConfigAssistantNavigation) onRowSelected(_ gtki.ListBox, r gtki.ListBoxRow) {
	// Every time a row is selected, we check if it's not a divider or a normal item.
	// Normal items will always be in even positions, and rows dividers will be in odd positions.
	if idx := getListBoxRowIndex(r); idx%2 == 0 {
		rcn.assistant.updateAssistantPage(idx - (idx / 2))
	}
}

// selectPageByIndex MUST be called from the UI thread
func (rcn *roomConfigAssistantNavigation) selectPageByIndex(idx int) {
	// Since we have a divider for each navigation item, every time we want to select
	// a page by its index, we should calculate the right index to avoid a wrong behavior
	row := rcn.navigation.GetRowAtIndex(idx * 2)
	rowIndex := getListBoxRowIndex(row)
	currentRowIndex := getListBoxRowIndex(rcn.navigation.GetSelectedRow())

	if rowIndex != noRowIndex && rowIndex != currentRowIndex {
		rcn.navigation.SelectRow(row)
	}
}

func (rcn *roomConfigAssistantNavigation) dividerIndexes() map[int]bool {
	indexes := map[int]bool{}

	// Each row has a divider near to it. We want the divider and not the row itself.
	for idx := 1; idx < len(rcn.assistant.allPages())*2; idx += 2 {
		indexes[idx] = true
	}

	return indexes
}

func getListBoxRowIndex(r gtki.ListBoxRow) int {
	if r != nil {
		return r.GetIndex()
	}
	return noRowIndex
}
