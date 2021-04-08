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

	assistantSidebar    *roomConfigAssistantSidebar
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
	rc.assistant.SetTitle(i18n.Localf("Configuration for room [%s]", rc.roomID.String()))

	for _, b := range rc.allPagesBoxes() {
		b.SetHExpand(true)
	}
}

func (rc *roomConfigAssistant) initSidebar() {
	rc.assistantSidebar = rc.newRoomConfigAssistantSidebar()
	setAssistantSidebarContent(rc.assistant, rc.assistantSidebar.box)
}

// refreshButtonLabels MUST be called from the UI thread
func (rc *roomConfigAssistant) refreshButtonLabels() {
	buttons := getButtonsForAssistantHeader(rc.assistant)

	buttons.updateButtonLabelByName("last", i18n.Local("Summary"))
	buttons.updateButtonLabelByName("apply", i18n.Local("Create Room"))
}

// hideBottomActionArea MUST be called from the UI thread
func (rc *roomConfigAssistant) hideBottomActionArea() {
	if actionArea, ok := getBottomActionAreaFromAssistant(rc.assistant); ok {
		actionArea.Hide()
	}
}

// onPageChanged MUST be called from the UI thread
func (rc *roomConfigAssistant) onPageChanged() {
	rc.updateAssistantPage(rc.assistant.GetCurrentPage())
}

// updateAssistantPage MUST be called from the UI thread
func (rc *roomConfigAssistant) updateAssistantPage(indexPage int) {
	if rc.canChangePage() {
		rc.pageByIndex(rc.currentPageIndex).collectData()
		rc.updateContentPage(indexPage)
		rc.assistantSidebar.selectOptionByIndex(indexPage)
	} else {
		rc.assistantSidebar.selectOptionByIndex(rc.currentPageIndex)
	}
}

// canChangePage MUST be called from the UI thread
func (rc *roomConfigAssistant) canChangePage() bool {
	previousPage := rc.pageByIndex(rc.currentPageIndex)
	if previousPage.isNotValid() {
		rc.assistant.SetCurrentPage(rc.currentPageIndex)
		rc.currentPage.showValidationErrors()
		return false
	}
	return true
}

// updateContentPage MUST be called from the UI thread
func (rc *roomConfigAssistant) updateContentPage(indexPage int) {
	rc.currentPageIndex = indexPage
	rc.currentPage = rc.pageByIndex(rc.currentPageIndex)

	rc.assistant.SetCurrentPage(rc.currentPageIndex)
	rc.currentPage.refresh()

	rc.refreshButtonLabels()
	rc.hideBottomActionArea()
}

// enableAssistant MUST be called from the UI thread
func (rc *roomConfigAssistant) enableAssistant() {
	rc.assistant.SetSensitive(true)
}

// disableAssistant MUST be called from the UI thread
func (rc *roomConfigAssistant) disableAssistant() {
	rc.assistant.SetSensitive(false)
}

// onCancelClicked MUST be called from the UI thread
func (rc *roomConfigAssistant) onCancelClicked() {
	cv := rc.newRoomConfigAssistantCancelView()
	cv.show()
}

// cancelConfiguration MUST be called from the UI thread
func (rc *roomConfigAssistant) cancelConfiguration() {
	rc.destroyAssistant()
	rc.onCancel()
	rc.roomConfigComponent.cancelConfiguration(rc.onCancelError)
}

// onCancelError MUST be called from the UI thread
func (rc *roomConfigAssistant) onCancelError(err error) {
	dr := createDialogErrorComponent(
		i18n.Local("Cancel room settings"),
		i18n.Local("We were unable to cancel the room configuration"),
		i18n.Local("An error occurred while trying to cancel the configuration of the room."),
	)

	dr.setParent(rc.assistant)
	dr.show()
}

// onApplyClicked MUST be called from the UI thread
func (rc *roomConfigAssistant) onApplyClicked() {
	rc.disableAssistant()
	rc.currentPage.onConfigurationApply()

	rc.roomConfigComponent.submitConfigurationForm(
		rc.onApplySuccess,
		rc.onApplyError,
	)
}

// onApplySuccess MUST be called from the UI thread
func (rc *roomConfigAssistant) onApplySuccess() {
	rc.onSuccess(rc.roomConfigComponent.autoJoin)
	rc.destroyAssistant()
}

// destroyAssistant MUST be called from the UI thread
func (rc *roomConfigAssistant) destroyAssistant() {
	rc.assistant.Destroy()
}

// onApplyError MUST be called from the UI thread
func (rc *roomConfigAssistant) onApplyError(err error) {
	rc.enableAssistant()
	rc.currentPage.onConfigurationApplyError()
	rc.currentPage.notifyError(rc.roomConfigComponent.friendlyConfigErrorMessage(err))
}

// show MUST be called from the UI thread
func (rc *roomConfigAssistant) showAssistant() {
	rc.assistant.ShowAll()
}

func (rc *roomConfigAssistant) allPagesBoxes() []gtki.Box {
	return []gtki.Box{
		rc.infoPageBox,
		rc.accessPageBox,
		rc.permissionsPageBox,
		rc.occupantsPageBox,
		rc.othersPageBox,
		rc.summaryPageBox,
	}
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
