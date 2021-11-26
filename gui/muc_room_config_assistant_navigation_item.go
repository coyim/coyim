package gui

import "github.com/coyim/gotk3adapter/gtki"

type roomConfigAssistantNavigationItem struct {
	page *roomConfigPage

	label gtki.Label      `gtk-widget:"room-config-assistant-navigation-item-label"`
	row   gtki.ListBoxRow `gtk-widget:"room-config-assistant-navigation-item-row"`
}

func (rcn *roomConfigAssistantNavigation) newRoomConfigAssistantNavigationItem(page *roomConfigPage) *roomConfigAssistantNavigationItem {
	itm := &roomConfigAssistantNavigationItem{
		page: page,
	}

	b := newBuilder("MUCRoomConfigAssistantNavigationItem")
	panicOnDevError(b.bindObjects(itm))

	itm.label.SetText(page.title)

	return itm
}

func (itm *roomConfigAssistantNavigationItem) pageID() mucRoomConfigPageID {
	return itm.page.pageID
}

// disable MUST be called from the UI thread
func (itm *roomConfigAssistantNavigationItem) disable() {
	disableField(itm.row)
}

// enable MUST be called from the UI thread
func (itm *roomConfigAssistantNavigationItem) enable() {
	enableField(itm.row)
}
