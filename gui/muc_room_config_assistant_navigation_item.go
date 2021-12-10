package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type navigationItemIconName string

func (n navigationItemIconName) String() string {
	return string(n)
}

// TODO: All these icon names SHOULD be reviewed.
// The current icon names are being used in order to test how it looks on differents SO.
const (
	basicInformationIconName navigationItemIconName = "goa-account-msn-symbolic"
	accessIconName           navigationItemIconName = "dialog-password-symbolic"
	permissionsIconName      navigationItemIconName = "system-switch-user-symbolic"
	positionsIconName        navigationItemIconName = "contact-new-symbolic"
	otherIconName            navigationItemIconName = "system-run-symbolic"
	sumaryIconName           navigationItemIconName = "view-list-bullet-symbolic"
)

var assistantNavigationIconByPage = map[mucRoomConfigPageID]navigationItemIconName{
	roomConfigInformationPageIndex: basicInformationIconName,
	roomConfigAccessPageIndex:      accessIconName,
	roomConfigPermissionsPageIndex: permissionsIconName,
	roomConfigPositionsPageIndex:   positionsIconName,
	roomConfigOthersPageIndex:      otherIconName,
	roomConfigSummaryPageIndex:     sumaryIconName,
}

type roomConfigAssistantNavigationItem struct {
	page *roomConfigPage

	row   gtki.ListBoxRow `gtk-widget:"room-config-assistant-navigation-item-row"`
	icon  gtki.Image      `gtk-widget:"room-config-assistant-navigation-item-icon"`
	label gtki.Label      `gtk-widget:"room-config-assistant-navigation-item-label"`
}

func (rcn *roomConfigAssistantNavigation) newRoomConfigAssistantNavigationItem(page *roomConfigPage) *roomConfigAssistantNavigationItem {
	itm := &roomConfigAssistantNavigationItem{
		page: page,
	}

	b := newBuilder("MUCRoomConfigAssistantNavigationItem")
	panicOnDevError(b.bindObjects(itm))

	itm.label.SetText(page.title)
	itm.icon.SetFromIconName(string(assistantNavigationIconByPage[page.pageID]), gtki.ICON_SIZE_MENU)

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
