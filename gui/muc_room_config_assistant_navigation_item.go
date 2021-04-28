package gui

import "github.com/coyim/gotk3adapter/gtki"

type roomConfigAssistantNavigationItem struct {
	label  gtki.Label      `gtk-widget:"room-config-assistant-navigation-item-label"`
	boxRow gtki.ListBoxRow `gtk-widget:"room-config-assistant-navigation-item-row"`
}

func newRoomConfigAssistantNavigationItem(lbl string) *roomConfigAssistantNavigationItem {
	itm := &roomConfigAssistantNavigationItem{}

	b := newBuilder("MUCRoomConfigAssistantNavigationItem")
	panicOnDevError(b.bindObjects(itm))

	itm.label.SetText(lbl)

	return itm
}
