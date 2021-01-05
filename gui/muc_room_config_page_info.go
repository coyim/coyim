package gui

import "github.com/coyim/gotk3adapter/gtki"

type roomConfigInfoPage struct {
	*roomConfigPageBase
	roomDescriptionBuffer gtki.TextBuffer

	roomTitle       gtki.Entry    `gtk-widget:"room-title"`
	roomDescription gtki.TextView `gtk-widget:"room-description"`
	roomLanguage    gtki.Entry    `gtk-widget:"room-language"`
	roomPersistent  gtki.Switch   `gtk-widget:"room-persistent"`
	roomPublic      gtki.Switch   `gtk-widget:"room-public"`
}

func (c *mucRoomConfigComponent) newRoomConfigInfoPage() mucRoomConfigPage {
	p := &roomConfigInfoPage{}
	p.roomConfigPageBase = c.newConfigPage("info", "MUCRoomConfigPageInfo", p, nil)

	p.roomDescriptionBuffer, _ = g.gtk.TextBufferNew(nil)
	p.roomDescription.SetBuffer(p.roomDescriptionBuffer)

	p.initDefaultValues()

	return p
}

func (p *roomConfigInfoPage) initDefaultValues() {
	setEntryText(p.roomTitle, p.form.Title)
	setTextViewText(p.roomDescription, p.form.Description)
	setEntryText(p.roomLanguage, p.form.Language)
	setSwitchActive(p.roomPersistent, p.form.Persistent)
	setSwitchActive(p.roomPublic, p.form.Public)
}

func (p *roomConfigInfoPage) collectData() {
	p.form.Title = getEntryText(p.roomTitle)
	p.form.Description = getTextViewText(p.roomDescription)
	p.form.Language = getEntryText(p.roomLanguage)
	p.form.Persistent = getSwitchActive(p.roomPersistent)
	p.form.Public = getSwitchActive(p.roomPublic)
}
