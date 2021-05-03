package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigInfoPage struct {
	*roomConfigPageBase
	roomDescriptionBuffer gtki.TextBuffer
	roomLanguageComponent *languageSelectorComponent

	roomTitle            gtki.Entry        `gtk-widget:"room-title"`
	roomDescription      gtki.TextView     `gtk-widget:"room-description"`
	roomLanguageCombobox gtki.ComboBoxText `gtk-widget:"room-language-combobox"`
	roomLanguageEntry    gtki.Entry        `gtk-widget:"room-language-entry"`
	roomPersistent       gtki.Switch       `gtk-widget:"room-persistent"`
	roomPublic           gtki.Switch       `gtk-widget:"room-public"`
}

func (c *mucRoomConfigComponent) newRoomConfigInfoPage() mucRoomConfigPage {
	p := &roomConfigInfoPage{}
	p.roomConfigPageBase = c.newConfigPage(pageConfigInfo, "MUCRoomConfigPageInfo", p, nil)

	p.roomDescriptionBuffer, _ = g.gtk.TextBufferNew(nil)
	p.roomDescription.SetBuffer(p.roomDescriptionBuffer)

	p.roomLanguageComponent = c.u.createLanguageSelectorComponent(p.roomLanguageEntry, p.roomLanguageCombobox)

	p.initDefaultValues()

	return p
}

func (p *roomConfigInfoPage) initDefaultValues() {
	setEntryText(p.roomTitle, p.form.GetStringValue(muc.ConfigFieldRoomName))
	setTextViewText(p.roomDescription, p.form.GetStringValue(muc.ConfigFieldRoomDescription))
	p.roomLanguageComponent.setLanguage(p.form.GetStringValue(muc.ConfigFieldLanguage))
	setSwitchActive(p.roomPersistent, p.form.Persistent)
	setSwitchActive(p.roomPublic, p.form.Public)
}

func (p *roomConfigInfoPage) collectData() {
	p.form.UpdateFieldValue(muc.ConfigFieldRoomName, getEntryText(p.roomTitle))
	p.form.UpdateFieldValue(muc.ConfigFieldRoomDescription, getTextViewText(p.roomDescription))
	p.form.UpdateFieldValue(muc.ConfigFieldLanguage, p.roomLanguageComponent.currentLanguage())
	p.form.Persistent = getSwitchActive(p.roomPersistent)
	p.form.Public = getSwitchActive(p.roomPublic)
}
