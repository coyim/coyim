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
			muc.RoomConfigFieldWhoIs,
			muc.RoomConfigFieldIsModerated,
			muc.RoomConfigFieldCanChangeSubject,
			muc.RoomConfigFieldAllowPrivateMessages,
			muc.RoomConfigFieldPresenceBroadcast,
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
		},
	}
}

func getPageIndexBasedOnPageID(pageID string) int {
	switch pageID {
	case pageConfigInfo:
		return roomConfigInformationPageIndex
	case pageConfigAccess:
		return roomConfigAccessPageIndex
	case pageConfigPermissions:
		return roomConfigPermissionsPageIndex
	case pageConfigOccupants:
		return roomConfigOccupantsPageIndex
	case pageConfigOthers:
		return roomConfigOthersPageIndex
	}
	return roomConfigSummaryPageIndex
}

type mucRoomConfigPage interface {
	pageView() gtki.Overlay
	pageTitle() string
	isValid() bool
	showValidationErrors()
	collectData()
	refresh()
	notifyError(string)
	onConfigurationApply()
	onConfigurationApplyError()
}

type roomConfigPageBase struct {
	u      *gtkUI
	form   *muc.RoomConfigForm
	fields []hasRoomConfigFormField

	title          string
	pageID         string
	setCurrentPage func(indexPage int)
	fieldsContent  gtki.Box

	page              gtki.Overlay `gtk-widget:"room-config-page-overlay"`
	header            gtki.Label   `gtk-widget:"room-config-page-header-label"`
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
		setCurrentPage: c.setCurrentPage,
		title:          configPageDisplayTitle(pageID),
		pageID:         pageID,
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

	fieldsContent, err := builder.GetObject("room-config-fields-content")
	if err != nil {
		panic(fmt.Sprintf("developer error: the ID for \"%s\" page doesn't exists", pageID))
	}

	p.fieldsContent = fieldsContent.(gtki.Box)
	p.initDefaults()

	mucStyles.setRoomConfigPageStyle(pageContent)

	return p
}

func (p *roomConfigPageBase) initDefaults() {
	switch p.pageID {
	case pageConfigSummary:
		p.initSummary()
		return
	case pageConfigOthers:
		p.initIntroPage()
		p.initKnownFields()
		p.initUnknownFields()
		return
	}
	p.initIntroPage()
	p.initKnownFields()
}

func (p *roomConfigPageBase) initIntroPage() {
	intro := configPageDisplayIntro(p.pageID)
	if intro == "" {
		p.header.SetVisible(false)
		return
	}
	p.header.SetText(intro)
}

func (p *roomConfigPageBase) initKnownFields() {
	if knownFields, ok := roomConfigPagesFields[p.pageID]; ok {
		booleanFields := []*roomConfigFormFieldBoolean{}
		for _, kf := range knownFields {
			if knownField, ok := p.form.GetKnownField(kf); ok {
				field, err := roomConfigFormFieldFactory(kf, roomConfigFieldsTexts[kf], knownField.ValueType())
				if err != nil {
					p.log.WithError(err).Error("Room configuration form field not supported")
					continue
				}
				if f, ok := field.(*roomConfigFormFieldBoolean); ok {
					booleanFields = append(booleanFields, f)
					continue
				}
				p.addField(field)
			}
		}
		if len(booleanFields) > 0 {
			p.addField(newRoomConfigFormFieldBooleanContainer(booleanFields))
		}
	}
}

func (p *roomConfigPageBase) initUnknownFields() {
	booleanFields := []*roomConfigFormFieldBoolean{}
	for _, ff := range p.form.GetUnknownFields() {
		field, err := roomConfigFormUnknownFieldFactory(newRoomConfigFieldTextInfo(ff.Label, ff.Description), ff.ValueType())
		if err != nil {
			p.log.WithError(err).Error("Room configuration form field not supported")
			continue
		}
		if f, ok := field.(*roomConfigFormFieldBoolean); ok {
			booleanFields = append(booleanFields, f)
			continue
		}
		p.addField(field)
	}
	if len(booleanFields) > 0 {
		p.addField(newRoomConfigFormFieldBooleanContainer(booleanFields))
	}
}

func (p *roomConfigPageBase) initSummary() {
	p.initSummaryFields(pageConfigInfo)
}

func (p *roomConfigPageBase) initSummaryFields(pageID string) {
	p.addField(newRoomConfigFormFieldLinkButton(pageID, p.setCurrentPage))
	fields := []*roomConfigSummaryField{}
	for _, kf := range roomConfigPagesFields[pageID] {
		knownField, _ := p.form.GetKnownField(kf)
		fields = append(fields, newRoomConfigSummaryField(kf, roomConfigFieldsTexts[kf], knownField.ValueType()))
	}
	p.addField(newRoomConfigSummaryFieldContainer(fields))
}

func (p *roomConfigPageBase) addField(field hasRoomConfigFormField) {
	p.fields = append(p.fields, field)
	p.fieldsContent.Add(field.fieldWidget())
	p.doAfterRefresh.add(field.refreshContent)
}

// pageTitle implements the "mucRoomConfigPage" interface
func (p *roomConfigPageBase) pageTitle() string {
	return p.title
}

// pageView implements the "mucRoomConfigPage" interface
func (p *roomConfigPageBase) pageView() gtki.Overlay {
	return p.page
}

// isValid implements the "mucRoomConfigPage" interface
func (p *roomConfigPageBase) isValid() bool {
	isValid := true
	for _, f := range p.fields {
		if !f.isValid() {
			f.showValidationErrors()
			isValid = false
		}
	}
	return isValid
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
