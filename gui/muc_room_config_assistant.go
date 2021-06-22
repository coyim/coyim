package gui

import (
	"fmt"
	"sync"

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

	roomConfigComponent *mucRoomConfigComponent
	navigation          *roomConfigAssistantNavigation

	autoJoin  bool
	onSuccess func(autoJoin bool)
	onCancel  func()

	currentPageIndex mucRoomConfigPageID
	currentPage      *roomConfigPage

	assistant gtki.Assistant `gtk-widget:"room-config-assistant"`

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
	assignDefaultCurrentPageOnce := sync.Once{}

	for _, p := range rc.roomConfigComponent.pages {
		ap := newRoomConfigAssistantPage(p)

		rc.assistant.AppendPage(ap.page)
		rc.assistant.SetPageTitle(ap.page, configPageDisplayTitle(p.pageID))
		rc.assistant.SetPageComplete(ap.page, true)

		if p.pageID == roomConfigSummaryPageIndex {
			rc.assistant.SetPageType(ap.page, gtki.ASSISTANT_PAGE_CONFIRM)
		}

		assignDefaultCurrentPageOnce.Do(func() {
			rc.currentPageIndex = p.pageID
			rc.currentPage = p
		})
	}
}

func (rc *roomConfigAssistant) initDefaults() {
	rc.assistant.SetTitle(i18n.Localf("Configuration for room [%s]", rc.roomID))
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

	rc.roomConfigComponent.configureRoom(
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

func (rc *roomConfigAssistant) allPages() []*roomConfigPage {
	return rc.roomConfigComponent.pages
}

func (rc *roomConfigAssistant) pageByIndex(pageID mucRoomConfigPageID) *roomConfigPage {
	if page, ok := rc.roomConfigComponent.getConfigPage(pageID); ok {
		return page
	}
	panic(fmt.Sprintf("developer error: unsupported room config assistant page \"%d\"", pageID))
}
