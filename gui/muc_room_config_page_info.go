package gui

import (
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
	p.roomConfigPageBase = c.newConfigPage("info", "MUCRoomConfigPageInfo", p, nil)

	p.roomDescriptionBuffer, _ = g.gtk.TextBufferNew(nil)
	p.roomDescription.SetBuffer(p.roomDescriptionBuffer)

	p.roomLanguageComponent = c.u.createLanguageSelectorComponent(p.roomLanguageEntry, p.roomLanguageCombobox)

	p.initDefaultValues()

	return p
}

func (p *roomConfigInfoPage) initDefaultValues() {
	setEntryText(p.roomTitle, p.form.Title)
	setTextViewText(p.roomDescription, p.form.Description)
	p.roomLanguageComponent.setLanguage(p.form.Language)
	setSwitchActive(p.roomPersistent, p.form.Persistent)
	setSwitchActive(p.roomPublic, p.form.Public)
}

func (p *roomConfigInfoPage) collectData() {
	p.form.Title = getEntryText(p.roomTitle)
	p.form.Description = getTextViewText(p.roomDescription)
	p.form.Language = p.roomLanguageComponent.currentLanguage()
	p.form.Persistent = getSwitchActive(p.roomPersistent)
	p.form.Public = getSwitchActive(p.roomPublic)
}
