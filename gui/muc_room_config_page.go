package gui

import (
	"github.com/coyim/coyim/i18n"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

var roomConfigPagesFields map[int][]muc.RoomConfigFieldType

func initMUCRoomConfigPages() {
	roomConfigPagesFields = map[int][]muc.RoomConfigFieldType{
		roomConfigInformationPageIndex: {
			muc.RoomConfigFieldName,
			muc.RoomConfigFieldDescription,
			muc.RoomConfigFieldLanguage,
			muc.RoomConfigFieldIsPublic,
			muc.RoomConfigFieldIsPersistent,
		},
		roomConfigAccessPageIndex: {
			muc.RoomConfigFieldIsPasswordProtected,
			muc.RoomConfigFieldPassword,
			muc.RoomConfigFieldIsMembersOnly,
			muc.RoomConfigFieldAllowInvites,
		},
		roomConfigPermissionsPageIndex: {
			muc.RoomConfigFieldWhoIs,
			muc.RoomConfigFieldIsModerated,
			muc.RoomConfigFieldCanChangeSubject,
			muc.RoomConfigFieldAllowPrivateMessages,
			muc.RoomConfigFieldPresenceBroadcast,
		},
		roomConfigOccupantsPageIndex: {
			muc.RoomConfigFieldOwners,
			muc.RoomConfigFieldAdmins,
			muc.RoomConfigFieldMembers,
		},
		roomConfigOthersPageIndex: {
			muc.RoomConfigFieldMaxOccupantsNumber,
			muc.RoomConfigFieldMaxHistoryFetch,
			muc.RoomConfigFieldEnableLogging,
			muc.RoomConfigFieldPubsub,
		},
	}
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

	title               string
	pageID              int
	roomConfigComponent *mucRoomConfigComponent

	page                gtki.Overlay     `gtk-widget:"room-config-page-overlay"`
	header              gtki.Label       `gtk-widget:"room-config-page-header-label"`
	content             gtki.Box         `gtk-widget:"room-config-page-content"`
	notificationsArea   gtki.Box         `gtk-widget:"notifications-box"`
	autojoinContent     gtki.Box         `gtk-widget:"room-config-autojoin-content"`
	autojoinCheckButton gtki.CheckButton `gtk-widget:"room-config-autojoin"`

	notifications  *notificationsComponent
	loadingOverlay *loadingOverlayComponent
	doAfterRefresh *callbacksSet

	log coylog.Logger
}

func (c *mucRoomConfigComponent) newConfigPage(pageID int) *roomConfigPageBase {
	p := &roomConfigPageBase{
		u:                   c.u,
		roomConfigComponent: c,
		title:               configPageDisplayTitle(pageID),
		pageID:              pageID,
		loadingOverlay:      c.u.newLoadingOverlayComponent(),
		doAfterRefresh:      newCallbacksSet(),
		form:                c.form,
		log: c.log.WithFields(log.Fields{
			"page": pageID,
		}),
	}

	p.initBuilder()
	p.initDefaults()
	mucStyles.setRoomConfigPageStyle(p.content)

	return p
}

func (p *roomConfigPageBase) initBuilder() {
	builder := newBuilder("MUCRoomConfigPage")
	panicOnDevError(builder.bindObjects(p))
	builder.ConnectSignals(map[string]interface{}{
		"on_autojoin_toggled": func() {
			p.roomConfigComponent.updateAutoJoin(p.autojoinCheckButton.GetActive())
		},
	})

	p.notifications = p.u.newNotificationsComponent()
	p.loadingOverlay = p.u.newLoadingOverlayComponent()
	p.notificationsArea.Add(p.notifications.contentBox())
}

func (p *roomConfigPageBase) initDefaults() {
	p.initIntroPage()
	switch p.pageID {
	case roomConfigSummaryPageIndex:
		p.initSummary()
		return
	case roomConfigOccupantsPageIndex:
		p.initOccupants()
		return
	case roomConfigOthersPageIndex:
		p.initKnownFields()
		p.initUnknownFields()
		return
	}
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
	p.initSummaryFields(roomConfigInformationPageIndex)
	p.initSummaryFields(roomConfigAccessPageIndex)
	p.initSummaryFields(roomConfigPermissionsPageIndex)
	p.initSummaryFields(roomConfigOthersPageIndex)
	p.autojoinContent.Show()
}

func (p *roomConfigPageBase) initSummaryFields(pageID int) {
	p.addField(newRoomConfigFormFieldLinkButton(pageID, p.roomConfigComponent.setCurrentPage))
	fields := []hasRoomConfigFormField{}
	for _, kf := range roomConfigPagesFields[pageID] {
		if knownField, ok := p.form.GetKnownField(kf); ok {
			fields = append(fields, newRoomConfigSummaryField(kf, roomConfigFieldsTexts[kf], knownField.ValueType()))
		}
	}
	if pageID == roomConfigOthersPageIndex {
		for _, ff := range p.form.GetUnknownFields() {
			fields = append(fields, newRoomConfigSummaryField(muc.RoomConfigFieldUnexpected, newRoomConfigFieldTextInfo(ff.Label, ff.Description), ff.ValueType()))
		}
	}
	p.addField(newRoomConfigSummaryFieldContainer(fields))
}

func (p *roomConfigPageBase) initOccupants() {
	p.addField(newRoomConfigOccupants(&data.OwnerAffiliation{}, p.form.GetRoomOccupants()))
	p.content.Add(createSeparator(gtki.HorizontalOrientation))
	p.addField(newRoomConfigOccupants(&data.AdminAffiliation{}, p.form.GetRoomOccupants()))
	p.content.Add(createSeparator(gtki.HorizontalOrientation))
	p.addField(newRoomConfigOccupants(&data.OutcastAffiliation{}, p.form.GetRoomOccupants()))
}

func (p *roomConfigPageBase) addField(field hasRoomConfigFormField) {
	p.fields = append(p.fields, field)
	p.content.Add(field.fieldWidget())
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
func (p *roomConfigPageBase) collectData() {
	for _, f := range p.fields {
		f.collectFieldValue()
	}
}

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
