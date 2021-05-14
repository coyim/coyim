package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigInfoPage struct {
	*roomConfigPageBase
	roomDescriptionBuffer gtki.TextBuffer

	roomDescription            gtki.TextView `gtk-widget:"room-description"`
	roomLanguageFieldContainer gtki.Box      `gtk-widget:"room-config-language-field"`
	roomPersistent             gtki.Switch   `gtk-widget:"room-persistent"`
	roomPublic                 gtki.Switch   `gtk-widget:"room-public"`
}

func (c *mucRoomConfigComponent) newRoomConfigInfoPage() mucRoomConfigPage {
	p := &roomConfigInfoPage{}
	p.roomConfigPageBase = c.newConfigPage(pageConfigInfo, "MUCRoomConfigPageInfo", p, nil)

	p.roomDescriptionBuffer, _ = g.gtk.TextBufferNew(nil)
	p.roomDescription.SetBuffer(p.roomDescriptionBuffer)

	p.initDefaultValues()

	return p
}

func (p *roomConfigInfoPage) initDefaultValues() {
	setTextViewText(p.roomDescription, p.form.Description)
	setSwitchActive(p.roomPersistent, p.form.Persistent)
	setSwitchActive(p.roomPublic, p.form.Public)
}

func (p *roomConfigInfoPage) collectData() {
	p.form.Description = getTextViewText(p.roomDescription)
	p.form.Persistent = getSwitchActive(p.roomPersistent)
	p.form.Public = getSwitchActive(p.roomPublic)
}
