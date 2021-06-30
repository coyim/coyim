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
	u                  *gtkUI
	account            *account
	roomID             jid.Bare
	roomConfigScenario roomConfigScenario

	roomConfigComponent *mucRoomConfigComponent
	navigation          *roomConfigAssistantNavigation

	autoJoin                        bool
	doAfterConfigSaved              func(autoJoin bool) // doAfterConfigSaved will be called from the UI thread
	doAfterConfigCanceled           func()              // doAfterConfigCanceled will be called from the UI thread
	doNotAskForConfirmationOnCancel bool

	currentPageIndex mucRoomConfigPageID
	currentPage      *roomConfigPage
	parentWindow     gtki.Window

	assistant gtki.Assistant `gtk-widget:"room-config-assistant"`

	log coylog.Logger
}

func (u *gtkUI) newRoomConfigAssistant(data *roomConfigData) *roomConfigAssistant {
	rc := &roomConfigAssistant{
		u:                               u,
		account:                         data.account,
		roomID:                          data.roomID,
		roomConfigScenario:              data.roomConfigScenario,
		autoJoin:                        data.autoJoinRoomAfterSaved,
		doAfterConfigSaved:              data.doAfterConfigSaved,
		doAfterConfigCanceled:           data.doAfterConfigCanceled,
		doNotAskForConfirmationOnCancel: data.doNotAskForConfirmationOnCancel,
		parentWindow:                    data.parentWindow,
		log: u.log.WithFields(log.Fields{
			"room":  data.roomID,
			"where": "configureRoomAssistant",
		}),
	}

	rc.initBuilder()
	rc.initRoomConfigComponent(data.configForm)
	rc.initRoomConfigPages()
	rc.initDefaults()
	rc.initSidebarNavigation()

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
	rc.roomConfigComponent = rc.u.newMUCRoomConfigComponent(rc.account, rc.roomID, form, rc.autoJoin, rc.setCurrentPage, rc.assistant)
}

func (rc *roomConfigAssistant) setCurrentPage(pageID mucRoomConfigPageID) {
	doInUIThread(func() { rc.assistant.SetCurrentPage(int(pageID)) })
}

func (rc *roomConfigAssistant) initRoomConfigPages() {
	assignedDefaultCurrentPage := false

	for _, p := range rc.roomConfigComponent.pages {
		ap := newRoomConfigAssistantPage(p)

		rc.assistant.AppendPage(ap.page)
		rc.assistant.SetPageTitle(ap.page, configPageDisplayTitle(p.pageID))
		rc.assistant.SetPageComplete(ap.page, true)

		if p.pageID == roomConfigSummaryPageIndex {
			rc.assistant.SetPageType(ap.page, gtki.ASSISTANT_PAGE_CONFIRM)
		}

		if !assignedDefaultCurrentPage {
			assignedDefaultCurrentPage = true
			rc.currentPageIndex = p.pageID
			rc.currentPage = p
		}
	}
}

func (rc *roomConfigAssistant) initDefaults() {
	rc.assistant.SetTitle(i18n.Localf("Configuration for room [%s]", rc.roomID))
	if rc.parentWindow != nil {
		rc.assistant.SetTransientFor(rc.parentWindow)
	}

	removeMarginFromAssistantPages(rc.assistant)
}

func (rc *roomConfigAssistant) initSidebarNavigation() {
	rc.navigation = rc.newRoomConfigAssistantNavigation()
	setAssistantSidebarContent(rc.assistant, rc.navigation.content)
}

// refreshButtonLabels MUST be called from the UI thread
func (rc *roomConfigAssistant) refreshButtonLabels() {
	buttons := getButtonsForAssistantHeader(rc.assistant)

	buttons.updateButtonLabelByName("last", i18n.Local("Summary"))
	buttons.updateButtonLabelByName("apply", rc.applyLabelBasedOnCurrentScenario())
}

func (rc *roomConfigAssistant) applyLabelBasedOnCurrentScenario() string {
	if rc.roomConfigScenario == roomConfigScenarioCreate {
		return i18n.Local("Create Room")
	}
	return i18n.Local("Update Configuration")
}

// hideBottomActionArea MUST be called from the UI thread
func (rc *roomConfigAssistant) hideBottomActionArea() {
	if actionArea, ok := getBottomActionAreaFromAssistant(rc.assistant); ok {
		actionArea.Hide()
	}
}

// onPageChanged MUST be called from the UI thread
func (rc *roomConfigAssistant) onPageChanged() {
	rc.updateAssistantPage(mucRoomConfigPageID(rc.assistant.GetCurrentPage()))
}

// updateAssistantPage MUST be called from the UI thread
func (rc *roomConfigAssistant) updateAssistantPage(pageID mucRoomConfigPageID) {
	if rc.canChangePage() {
		rc.pageByIndex(rc.currentPageIndex).updateFieldValues()
		rc.updateContentPage(pageID)
		rc.navigation.selectPageByIndex(pageID)
	} else {
		rc.navigation.selectPageByIndex(rc.currentPageIndex)
	}
}

// canChangePage MUST be called from the UI thread
func (rc *roomConfigAssistant) canChangePage() bool {
	previousPage := rc.pageByIndex(rc.currentPageIndex)
	if !previousPage.isValid() {
		rc.assistant.SetCurrentPage(int(rc.currentPageIndex))
		return false
	}
	return true
}

// updateContentPage MUST be called from the UI thread
func (rc *roomConfigAssistant) updateContentPage(indexPage mucRoomConfigPageID) {
	rc.currentPageIndex = indexPage
	rc.currentPage = rc.pageByIndex(rc.currentPageIndex)

	rc.assistant.SetCurrentPage(int(rc.currentPageIndex))
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
	if rc.doNotAskForConfirmationOnCancel {
		rc.onCancel()
	} else {
		cv := rc.newRoomConfigAssistantCancelView()
		cv.show()
	}
}

// cancelConfiguration MUST be called from the UI thread
func (rc *roomConfigAssistant) cancelConfiguration() {
	rc.onCancel()
}

// onCancel MUST be called from the UI thread
func (rc *roomConfigAssistant) onCancel() {
	rc.destroyAssistant()

	if rc.doAfterConfigCanceled != nil {
		rc.doAfterConfigCanceled()
	}
}

// onCancelError MUST be called from the UI thread
func (rc *roomConfigAssistant) onCancelError(err error) {
	dr := createDialogErrorComponent(dialogErrorOptions{
		title:   i18n.Local("Cancel room settings"),
		header:  i18n.Local("We were unable to cancel the room configuration"),
		message: i18n.Local("An error occurred while trying to cancel the configuration of the room."),
	})

	dr.setParent(rc.assistant)
	dr.show()
}

// onApplyClicked MUST be called from the UI thread
func (rc *roomConfigAssistant) onApplyClicked() {
	rc.disableAssistant()
	rc.currentPage.onConfigurationApply()

	rc.roomConfigComponent.configureRoom(
		rc.onApplySuccess,
		rc.onApplyError,
	)
}

// onApplySuccess MUST be called from the UI thread
func (rc *roomConfigAssistant) onApplySuccess() {
	rc.destroyAssistant()

	if rc.doAfterConfigSaved != nil {
		rc.doAfterConfigSaved(rc.roomConfigComponent.autoJoin)
	}
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

func (rc *roomConfigAssistant) allPages() []*roomConfigPage {
	return rc.roomConfigComponent.pages
}

func (rc *roomConfigAssistant) pageByIndex(pageID mucRoomConfigPageID) *roomConfigPage {
	if page, ok := rc.roomConfigComponent.getConfigPage(pageID); ok {
		return page
	}
	panic(fmt.Sprintf("developer error: unsupported room config assistant page \"%d\"", pageID))
}

// canConfigureRoom implements the "markAsComponentToConfigureARoom" interface
func (rc *roomConfigAssistant) canConfigureRoom() bool {
	return true
}

// launchRoomConfigView implements the "markAsComponentToConfigureARoom" interface
func (rc *roomConfigAssistant) launchRoomConfigView() {
	rc.showAssistant()
}
