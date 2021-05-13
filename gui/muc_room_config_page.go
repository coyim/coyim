package gui

import (
	"fmt"

	"github.com/coyim/coyim/i18n"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

const (
	pageConfigInfo        = "info"
	pageConfigAccess      = "access"
	pageConfigPermissions = "permissions"
	pageConfigOccupants   = "occupants"
	pageConfigOthers      = "others"
	pageConfigSummary     = "summary"
)

var roomConfigPagesFields map[string][]muc.RoomConfigFieldType

func initMUCRoomConfigPages() {
	roomConfigPagesFields = map[string][]muc.RoomConfigFieldType{
		pageConfigInfo: {
			muc.RoomConfigFieldName,
			muc.RoomConfigFieldDescription,
			muc.RoomConfigFieldLanguage,
			muc.RoomConfigFieldIsPublic,
			muc.RoomConfigFieldIsPersistent,
		},
		pageConfigAccess: {
			muc.RoomConfigFieldIsPasswordProtected,
			muc.RoomConfigFieldPassword,
			muc.RoomConfigFieldIsMembersOnly,
			muc.RoomConfigFieldAllowInvites,
		},
		pageConfigPermissions: {
			muc.RoomConfigFieldPresenceBroadcast,
			muc.RoomConfigFieldIsModerated,
			muc.RoomConfigFieldCanChangeSubject,
			muc.RoomConfigFieldAllowPrivateMessages,
		},
		pageConfigOccupants: {
			muc.RoomConfigFieldOwners,
			muc.RoomConfigFieldAdmins,
			muc.RoomConfigFieldMembers,
		},
		pageConfigOthers: {
			muc.RoomConfigFieldMaxOccupantsNumber,
			muc.RoomConfigFieldMaxHistoryFetch,
			muc.RoomConfigFieldEnableLogging,
			muc.RoomConfigFieldPubsub,
			muc.RoomConfigFieldWhoIs,
		},
	}
}

type mucRoomConfigPage interface {
	pageView() gtki.Overlay
	pageTitle() string
	isInvalid() bool
	showValidationErrors()
	collectData()
	refresh()
	notifyError(string)
	onConfigurationApply()
	onConfigurationApplyError()
}

type roomConfigPageBase struct {
	u    *gtkUI
	form *muc.RoomConfigForm

	title string

	page              gtki.Overlay `gtk-widget:"room-config-page-overlay"`
	content           gtki.Box     `gtk-widget:"room-config-page-content"`
	notificationsArea gtki.Box     `gtk-widget:"notifications-box"`

	notifications  *notificationsComponent
	loadingOverlay *loadingOverlayComponent
	doAfterRefresh *callbacksSet

	log coylog.Logger
}

func (c *mucRoomConfigComponent) newConfigPage(pageID, pageTemplate string, page interface{}, signals map[string]interface{}) *roomConfigPageBase {
	p := &roomConfigPageBase{
		u:              c.u,
		title:          configPageDisplayTitle(pageID),
		loadingOverlay: c.u.newLoadingOverlayComponent(),
		doAfterRefresh: newCallbacksSet(),
		form:           c.form,
		log: c.log.WithFields(log.Fields{
			"page":     pageID,
			"template": pageTemplate,
		}),
	}

	builder := newBuilder("MUCRoomConfigPage")
	panicOnDevError(builder.bindObjects(p))

	p.notifications = c.u.newNotificationsComponent()
	p.loadingOverlay = c.u.newLoadingOverlayComponent()
	p.notificationsArea.Add(p.notifications.contentBox())

	p.page.AddOverlay(p.loadingOverlay.getOverlay())
	p.page.SetHExpand(true)
	p.page.SetVExpand(true)

	builder = newBuilder(pageTemplate)
	panicOnDevError(builder.bindObjects(page))
	builder.ConnectSignals(signals)

	pc, err := builder.GetObject(fmt.Sprintf("room-config-%s-page", pageID))
	if err != nil {
		panic(fmt.Sprintf("developer error: the ID for \"%s\" page doesn't exists", pageID))
	}

	pageContent := pc.(gtki.Box)
	pageContent.SetHExpand(false)
	p.content.Add(pageContent)

	mucStyles.setRoomConfigPageStyle(pageContent)

	return p
}

// pageTitle implements the "mucRoomConfigPage" interface
func (p *roomConfigPageBase) pageTitle() string {
	return p.title
}

// pageView implements the "mucRoomConfigPage" interface
func (p *roomConfigPageBase) pageView() gtki.Overlay {
	return p.page
}

// isInvalid implements the "mucRoomConfigPage" interface
func (p *roomConfigPageBase) isInvalid() bool {
	return false
}

// validate implements the "mucRoomConfigPage" interface
func (p *roomConfigPageBase) showValidationErrors() {
}

// Nothing to do, just implement the "mucRoomConfigPage" interface
func (p *roomConfigPageBase) collectData() {}

// refresh MUST be called from the UI thread
func (p *roomConfigPageBase) refresh() {
	p.page.ShowAll()
	p.hideLoadingOverlay()
	p.clearErrors()
	p.doAfterRefresh.invokeAll()
}

// clearErrors MUST be called from the ui thread
func (p *roomConfigPageBase) clearErrors() {
	p.notifications.clearErrors()
}

// notifyError MUST be called from the ui thread
func (p *roomConfigPageBase) notifyError(m string) {
	p.notifications.notifyOnError(m)
}

// onConfigurationApply MUST be called from the ui thread
func (p *roomConfigPageBase) onConfigurationApply() {
	p.showLoadingOverlay(i18n.Local("Saving room configuration"))
}

// onConfigurationApplyError MUST be called from the ui thread
func (p *roomConfigPageBase) onConfigurationApplyError() {
	p.hideLoadingOverlay()
}

// showLoadingOverlay MUST be called from the ui thread
func (p *roomConfigPageBase) showLoadingOverlay(m string) {
	p.loadingOverlay.setSolid()
	p.loadingOverlay.showWithMessage(m)
}

// hideLoadingOverlay MUST be called from the ui thread
func (p *roomConfigPageBase) hideLoadingOverlay() {
	p.loadingOverlay.hide()
}
