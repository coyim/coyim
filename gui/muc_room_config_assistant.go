package gui

import (
	"fmt"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

type roomConfigAssistant struct {
	u       *gtkUI
	account *account
	roomID  jid.Bare

	roomConfigSidebar   *roomConfigAssistantSidebar
	roomConfigComponent *mucRoomConfigComponent

	autoJoin  bool
	onSuccess func(autoJoin bool)
	onCancel  func()

	currentPageIndex int
	currentPage      mucRoomConfigPage

	assistant          gtki.Assistant `gtk-widget:"room-config-assistant"`
	infoPageBox        gtki.Box       `gtk-widget:"room-config-info-page"`
	accessPageBox      gtki.Box       `gtk-widget:"room-config-access-page"`
	permissionsPageBox gtki.Box       `gtk-widget:"room-config-permissions-page"`
	occupantsPageBox   gtki.Box       `gtk-widget:"room-config-occupants-page"`
	othersPageBox      gtki.Box       `gtk-widget:"room-config-others-page"`
	summaryPageBox     gtki.Box       `gtk-widget:"room-config-summary-page"`

	infoPage        mucRoomConfigPage
	accessPage      mucRoomConfigPage
	permissionsPage mucRoomConfigPage
	occupantsPage   mucRoomConfigPage
	othersPage      mucRoomConfigPage
	summaryPage     mucRoomConfigPage

	log coylog.Logger
}

func (u *gtkUI) newRoomConfigAssistant(ca *account, roomID jid.Bare, form *muc.RoomConfigForm, autoJoin bool, onSuccess func(autoJoin bool), onCancel func()) *roomConfigAssistant {
	rc := &roomConfigAssistant{
		u:        u,
		account:  ca,
		roomID:   roomID,
		autoJoin: autoJoin,
		log: u.log.WithFields(log.Fields{
			"room":  roomID,
			"where": "configureRoomAssistant",
		}),
	}

	rc.onSuccess = func(aj bool) {
		if onSuccess != nil {
			onSuccess(aj)
		}
	}

	rc.onCancel = func() {
		if onCancel != nil {
			onCancel()
		}
	}

	rc.initBuilder()
	rc.initRoomConfigComponent(form)
	rc.initRoomConfigPages()
	rc.initDefaults()
	rc.initSidebar()

	return rc
}

func (rc *roomConfigAssistant) initBuilder() {
	b := newBuilder("MUCRoomConfigAssistant")
	panicOnDevError(b.bindObjects(rc))

	b.ConnectSignals(map[string]interface{}{
		"on_page_changed": rc.onPageChanged,
		"on_cancel":       rc.onCancelClicked,
		"on_apply":        rc.onApplyClicked,
	})
}

func (rc *roomConfigAssistant) initRoomConfigComponent(form *muc.RoomConfigForm) {
	rc.roomConfigComponent = rc.u.newMUCRoomConfigComponent(rc.account, rc.roomID, form, rc.autoJoin, rc.assistant)
	rc.roomConfigComponent.setCurrentPage = rc.assistant.SetCurrentPage
}

func (rc *roomConfigAssistant) initRoomConfigPages() {
	rc.infoPage = rc.roomConfigComponent.getConfigPage(roomConfigInformationPageIndex)
	rc.accessPage = rc.roomConfigComponent.getConfigPage(roomConfigAccessPageIndex)
	rc.permissionsPage = rc.roomConfigComponent.getConfigPage(roomConfigPermissionsPageIndex)
	rc.occupantsPage = rc.roomConfigComponent.getConfigPage(roomConfigOccupantsPageIndex)
	rc.othersPage = rc.roomConfigComponent.getConfigPage(roomConfigOthersPageIndex)
	rc.summaryPage = rc.roomConfigComponent.getConfigPage(roomConfigSummaryPageIndex)

	rc.infoPageBox.Add(rc.infoPage.pageView())
	rc.accessPageBox.Add(rc.accessPage.pageView())
	rc.permissionsPageBox.Add(rc.permissionsPage.pageView())
	rc.occupantsPageBox.Add(rc.occupantsPage.pageView())
	rc.othersPageBox.Add(rc.othersPage.pageView())
	rc.summaryPageBox.Add(rc.summaryPage.pageView())

	rc.currentPageIndex = roomConfigInformationPageIndex
	rc.currentPage = rc.infoPage
}

func (rc *roomConfigAssistant) initDefaults() {
	rc.infoPageBox.SetHExpand(true)
	rc.accessPageBox.SetHExpand(true)
	rc.permissionsPageBox.SetHExpand(true)
	rc.occupantsPageBox.SetHExpand(true)
	rc.othersPageBox.SetHExpand(true)
	rc.summaryPageBox.SetHExpand(true)

	rc.assistant.SetTitle(i18n.Localf("Configuration for room [%s]", rc.roomID.String()))
}

func (rc *roomConfigAssistant) initSidebar() {
	rc.roomConfigSidebar = rc.newRoomConfigAssistantSidebar()
	updateSidebarContent(rc.assistant, rc.roomConfigSidebar.getContent())
}

func (rc *roomConfigAssistant) refreshButtonLabels() {
	buttons := getButtonsForAssistantHeader(rc.assistant)

	buttons.updateButtonLabelByName("last", i18n.Local("Summary"))
	buttons.updateButtonLabelByName("apply", i18n.Local("Create Room"))
}

func (rc *roomConfigAssistant) onPageChanged() {
	if rc.canChangeOfPage() {
		rc.pageByIndex(rc.currentPageIndex).collectData()
		rc.roomConfigSidebar.selectOption(rc.assistant.GetCurrentPage())
		rc.updateContentPage(rc.assistant.GetCurrentPage())
	}
}

func (rc *roomConfigAssistant) updateAssistantPage(indexPage int) {
	if rc.canChangeOfPage() {
		rc.pageByIndex(rc.currentPageIndex).collectData()
		rc.updateContentPage(indexPage)
		return
	}
	rc.roomConfigSidebar.selectOption(rc.currentPageIndex)
}

func (rc *roomConfigAssistant) canChangeOfPage() bool {
	previousPage := rc.pageByIndex(rc.currentPageIndex)
	if previousPage.isNotValid() {
		rc.assistant.SetCurrentPage(rc.currentPageIndex)
		rc.currentPage.showValidationErrors()
		return false
	}

	return true
}

func (rc *roomConfigAssistant) updateContentPage(indexPage int) {
	rc.currentPageIndex = indexPage
	rc.currentPage = rc.pageByIndex(rc.currentPageIndex)
	rc.assistant.SetCurrentPage(rc.currentPageIndex)
	rc.currentPage.refresh()
	rc.refreshButtonLabels()
	removeActionArea(rc.assistant)
}

func (rc *roomConfigAssistant) enable() {
	rc.assistant.SetSensitive(true)
}

func (rc *roomConfigAssistant) disable() {
	rc.assistant.SetSensitive(false)
}

func (rc *roomConfigAssistant) onCancelClicked() {
	cv := newRoomConfigAssistantCancelView(rc)
	cv.show()
}

func (rc *roomConfigAssistant) onCancelError(err error) {
	dr := createDialogErrorComponent(
		"Cancel room settings",
		"We were unable to cancel the room configuration",
		"An error occurred while trying to cancel the configuration of the room.",
	)

	dr.show()
}

func (rc *roomConfigAssistant) onApplyClicked() {
	rc.disable()
	rc.currentPage.onConfigurationApply()

	rc.roomConfigComponent.submitConfigurationForm(
		rc.onApplySuccess,
		rc.onApplyError,
	)
}

func (rc *roomConfigAssistant) onApplySuccess() {
	rc.onSuccess(rc.roomConfigComponent.autoJoin)
	rc.destroyAssistant()
}

func (rc *roomConfigAssistant) destroyAssistant() {
	doInUIThread(rc.assistant.Destroy)
}

func (rc *roomConfigAssistant) onApplyError(err error) {
	rc.enable()
	rc.currentPage.onConfigurationApplyError()
	rc.currentPage.notifyError(rc.roomConfigComponent.friendlyConfigErrorMessage(err))
}

func (rc *roomConfigAssistant) pageByIndex(p int) mucRoomConfigPage {
	switch p {
	case roomConfigInformationPageIndex:
		return rc.infoPage
	case roomConfigAccessPageIndex:
		return rc.accessPage
	case roomConfigPermissionsPageIndex:
		return rc.permissionsPage
	case roomConfigOccupantsPageIndex:
		return rc.occupantsPage
	case roomConfigOthersPageIndex:
		return rc.othersPage
	case roomConfigSummaryPageIndex:
		return rc.summaryPage
	default:
		panic(fmt.Sprintf("developer error: unsupported room config assistant page \"%d\"", p))
	}
}

func (rc *roomConfigAssistant) show() {
	rc.assistant.ShowAll()
}
