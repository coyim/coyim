package gui

import "github.com/coyim/gotk3adapter/gtki"

type roomConfigAssistantNavigationItem struct {
	label   gtki.Label      `gtk-widget:"room-config-assistant-navigation-item-label"`
	row     gtki.ListBoxRow `gtk-widget:"room-config-assistant-navigation-item-row"`
	divider gtki.Separator  `gtk-widget:"room-config-assistant-navigation-item-separator"`
}

func (rcn *roomConfigAssistantNavigation) newRoomConfigAssistantNavigationItem(lbl string) *roomConfigAssistantNavigationItem {
	itm := &roomConfigAssistantNavigationItem{}

	b := newBuilder("MUCRoomConfigAssistantNavigationItem")
	panicOnDevError(b.bindObjects(itm))

	itm.label.SetText(lbl)

	return itm
}
