package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type navigationItemIconName string

func (n navigationItemIconName) String() string {
	return string(n)
}

const (
	basicInformationIconName     navigationItemIconName = "room_config_basic_information"
	basicInformationIconNameDark navigationItemIconName = "room_config_basic_information_dark"
	accessIconName               navigationItemIconName = "room_config_access"
	accessIconNameDark           navigationItemIconName = "room_config_access_dark"
	permissionsIconName          navigationItemIconName = "room_config_permissions"
	permissionsIconNameDark      navigationItemIconName = "room_config_permissions_dark"
	positionsIconName            navigationItemIconName = "room_config_positions"
	positionsIconNameDark        navigationItemIconName = "room_config_positions_dark"
	otherIconName                navigationItemIconName = "room_config_others"
	otherIconNameDark            navigationItemIconName = "room_config_others_dark"
	sumaryIconName               navigationItemIconName = "room_config_summary"
	sumaryIconNameDark           navigationItemIconName = "room_config_summary_dark"
)

var assistantIconSet = assistantNavigationIconByPage

var assistantNavigationIconByPage = map[mucRoomConfigPageID]navigationItemIconName{
	roomConfigInformationPageIndex: basicInformationIconName,
	roomConfigAccessPageIndex:      accessIconName,
	roomConfigPermissionsPageIndex: permissionsIconName,
	roomConfigPositionsPageIndex:   positionsIconName,
	roomConfigOthersPageIndex:      otherIconName,
	roomConfigSummaryPageIndex:     sumaryIconName,
}

var assistantNavigationDarkIconByPage = map[mucRoomConfigPageID]navigationItemIconName{
	roomConfigInformationPageIndex: basicInformationIconNameDark,
	roomConfigAccessPageIndex:      accessIconNameDark,
	roomConfigPermissionsPageIndex: permissionsIconNameDark,
	roomConfigPositionsPageIndex:   positionsIconNameDark,
	roomConfigOthersPageIndex:      otherIconNameDark,
	roomConfigSummaryPageIndex:     sumaryIconNameDark,
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
	itm.icon.SetFromPixbuf(getMUCIconPixbuf(assistantIconSet[page.pageID].String()))

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
