package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigAssistant struct {
	u                   *gtkUI
	roomConfigComponent *mucRoomConfigComponent

	assistant                gtki.Assistant `gtk-widget:"room-config-assistant"`
	roomConfigInfoBox        gtki.Box       `gtk-widget:"room-config-info-page"`
	roomConfigAccessBox      gtki.Box       `gtk-widget:"room-config-access-page"`
	roomConfigPermissionsBox gtki.Box       `gtk-widget:"room-config-permissions-page"`
	roomConfigOccupantsBox   gtki.Box       `gtk-widget:"room-config-occupants-page"`
	roomConfigOthersBox      gtki.Box       `gtk-widget:"room-config-others-page"`
	roomConfigSummaryBox     gtki.Box       `gtk-widget:"room-config-summary-page"`

	roomConfigInfoPage        *mucRoomConfigPage
	roomConfigAccessPage      *mucRoomConfigPage
	roomConfigPermissionsPage *mucRoomConfigPage
	roomConfigOccupantsPage   *mucRoomConfigPage
	roomConfigOthersPage      *mucRoomConfigPage
	roomConfigSummaryPage     *mucRoomConfigPage

	log coylog.Logger
}

func (u *gtkUI) newRoomConfigAssistant(form *muc.RoomConfigForm) *roomConfigAssistant {
	rc := &roomConfigAssistant{u: u}

	rc.initBuilder()
	rc.initRoomConfigComponent(form)
	rc.initRoomConfigPages()
	rc.initDefaults()

	return rc
}

func (rc *roomConfigAssistant) initBuilder() {
	b := newBuilder("MUCRoomConfigAssistant")
	panicOnDevError(b.bindObjects(rc))

	b.ConnectSignals(map[string]interface{}{
		"on_cancel":       rc.onCancel,
		"on_page_changed": rc.onPageChanged,
	})
}

func (rc *roomConfigAssistant) initRoomConfigComponent(form *muc.RoomConfigForm) {
	rc.roomConfigComponent = rc.u.newMUCRoomConfigComponent(form)
}

func (rc *roomConfigAssistant) initRoomConfigPages() {
	rc.roomConfigInfoPage = rc.roomConfigComponent.getConfigPage("information")
	rc.roomConfigAccessPage = rc.roomConfigComponent.getConfigPage("access")
	rc.roomConfigPermissionsPage = rc.roomConfigComponent.getConfigPage("permissions")
	rc.roomConfigOccupantsPage = rc.roomConfigComponent.getConfigPage("occupants")
	rc.roomConfigOthersPage = rc.roomConfigComponent.getConfigPage("others")
	rc.roomConfigSummaryPage = rc.roomConfigComponent.getConfigPage("summary")

	rc.roomConfigInfoBox.Add(rc.roomConfigInfoPage.getPageView())
	rc.roomConfigAccessBox.Add(rc.roomConfigAccessPage.getPageView())
	rc.roomConfigPermissionsBox.Add(rc.roomConfigPermissionsPage.getPageView())
	rc.roomConfigOccupantsBox.Add(rc.roomConfigOccupantsPage.getPageView())
	rc.roomConfigOthersBox.Add(rc.roomConfigOthersPage.getPageView())
	rc.roomConfigSummaryBox.Add(rc.roomConfigSummaryPage.getPageView())
}

func (rc *roomConfigAssistant) initDefaults() {
	rc.roomConfigInfoBox.SetHExpand(true)
	rc.roomConfigAccessBox.SetHExpand(true)
	rc.roomConfigPermissionsBox.SetHExpand(true)
	rc.roomConfigOccupantsBox.SetHExpand(true)
	rc.roomConfigOthersBox.SetHExpand(true)
	rc.roomConfigSummaryBox.SetHExpand(true)
}

func (rc *roomConfigAssistant) onCancel() {
	rc.assistant.Destroy()
}

func (rc *roomConfigAssistant) onPageChanged(_ gtki.Assistant, p gtki.Widget) {
	rc.assistant.SetPageComplete(p, true)

	switch rc.assistant.GetCurrentPage() {
	case roomConfigInformationPage:
		// TODO: Add implementation for "basic information" step
	case roomConfigAccessPage:
		// TODO: Add implementation for "access" step
	case roomConfigPermissionsPage:
		// TODO: Add implementation for "permissions" step
	case roomConfigOccupantsPage:
		// TODO: Add implementation for "occupants" step
	case roomConfigOthersPage:
		// TODO: Add implementation for "other settings" step
	case roomConfigSummaryPage:
		// TODO: Add implementation for "summary" step
	}
}

func (rc *roomConfigAssistant) show() {
	rc.assistant.ShowAll()
}
