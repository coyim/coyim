package gui

import (
	"fmt"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
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

	doAfterConfigSaved              func(autoJoin bool) // doAfterConfigSaved will be called from the UI thread
	doAfterConfigCanceled           func()              // doAfterConfigCanceled will be called from the UI thread
	doNotAskForConfirmationOnCancel bool

	currentPageIndex mucRoomConfigPageID
	currentPage      *roomConfigPage
	assistantPages   map[mucRoomConfigPageID]*roomConfigAssistantPage
	assistantButtons assistantButtons
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
	rc.initRoomConfigComponent(data)
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

func (rc *roomConfigAssistant) initRoomConfigComponent(data *roomConfigData) {
	rc.roomConfigComponent = rc.u.newMUCRoomConfigComponent(rc.account, data, rc.setCurrentPage, rc.assistant)

	rc.roomConfigComponent.onValidationErrors.add(func() {
		doInUIThread(rc.disableNavigation)
	})

	rc.roomConfigComponent.onNoValidationErrors.add(func() {
		doInUIThread(rc.enableNavigation)
	})
}

func (rc *roomConfigAssistant) initRoomConfigPages() {
	rc.assistantPages = map[mucRoomConfigPageID]*roomConfigAssistantPage{}
	assignedDefaultCurrentPage := false

	for _, p := range rc.roomConfigComponent.pages {
		ap := newRoomConfigAssistantPage(p)
		rc.assistantPages[p.pageID] = ap

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

	rc.assistantButtons = getButtonsForAssistantHeader(rc.assistant)
	removeMarginFromAssistantPages(rc.assistant)

	mucStyles.setRoomConfigSummaryStyle(rc.assistant)
}

func (rc *roomConfigAssistant) initSidebarNavigation() {
	rc.navigation = rc.newRoomConfigAssistantNavigation()
	setAssistantSidebarContent(rc.assistant, rc.navigation.content)
}

// refreshButtonLabels MUST be called from the UI thread
func (rc *roomConfigAssistant) refreshButtonLabels() {
	rc.assistantButtons.updateLastButtonLabel(i18n.Local("Summary"))
	rc.assistantButtons.updateApplyButtonLabel(rc.applyLabelBasedOnCurrentScenario())
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
		rc.setCurrentPage(rc.currentPageIndex)
		rc.navigation.selectPageByIndex(rc.currentPageIndex)
		rc.disableNavigation()
	}
}

// canChangePage MUST be called from the UI thread
func (rc *roomConfigAssistant) canChangePage() bool {
	previousPage := rc.pageByIndex(rc.currentPageIndex)
	return previousPage.isValid()
}

// updateContentPage MUST be called from the UI thread
func (rc *roomConfigAssistant) updateContentPage(pageID mucRoomConfigPageID) {
	rc.setCurrentPage(pageID)
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

// enableNavigation MUST be called from the UI thread
func (rc *roomConfigAssistant) enableNavigation() {
	rc.navigation.enableNavigation()
	rc.setPageComplete(rc.currentPageIndex, true)
	rc.assistantButtons.enableNavigation()
}

// disableNavigation MUST be called from the UI thread
func (rc *roomConfigAssistant) disableNavigation() {
	rc.navigation.disableNavigation()
	rc.setPageComplete(rc.currentPageIndex, false)
	rc.assistantButtons.disableNavigationButNotCancel()
}

// setPageComplete MUST be called from the UI thread
func (rc *roomConfigAssistant) setPageComplete(pageID mucRoomConfigPageID, setting bool) {
	if p, ok := rc.assistantPages[pageID]; ok {
		rc.assistant.SetPageComplete(p.page, setting)
	}
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
		title:        i18n.Local("Cancel room settings"),
		header:       i18n.Local("We were unable to cancel the room configuration"),
		message:      i18n.Local("An error occurred while trying to cancel the configuration of the room."),
		parentWindow: rc.assistant,
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
		rc.doAfterConfigSaved(rc.roomConfigComponent.data.autoJoinRoomAfterSaved)
	}
}

// destroyAssistant MUST be called from the UI thread
func (rc *roomConfigAssistant) destroyAssistant() {
	rc.assistant.Destroy()
}

// onApplyError MUST be called from the UI thread
func (rc *roomConfigAssistant) onApplyError(sfe *muc.SubmitFormError) {
	rc.enableAssistant()
	rc.currentPage.onConfigurationApplyError()

	if sfe.Error() == session.ErrBadRequestResponse {
		rc.onBadRequestError(sfe)
	}

	rc.currentPage.notifyError(rc.roomConfigComponent.friendlyConfigErrorMessage(sfe.Error()))
}

// onBadRequestError MUST be called from the UI thread
func (rc *roomConfigAssistant) onBadRequestError(sfe *muc.SubmitFormError) {
	pageID := getPageBasedOnField(sfe.Field())
	rc.setCurrentPage(pageID)

	rc.disableNavigation()

	for _, f := range rc.currentPage.fields {
		if f.fieldKey() == sfe.Field() {
			f.showValidationErrors()
		}
	}
}

// setCurrentPage MUST be called from the UI thread
func (rc *roomConfigAssistant) setCurrentPage(pageID mucRoomConfigPageID) {
	rc.currentPageIndex = pageID
	rc.currentPage = rc.pageByIndex(rc.currentPageIndex)
	rc.assistant.SetCurrentPage(rc.currentPageIndex.index())
}

func getPageBasedOnField(field muc.RoomConfigFieldType) mucRoomConfigPageID {
	for pageID, fields := range roomConfigPagesFields {
		for _, f := range fields {
			if f == field {
				return pageID
			}
		}
	}
	return roomConfigOthersPageIndex
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
