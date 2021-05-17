package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigInfoPage struct {
	*roomConfigPageBase
	roomDescriptionBuffer gtki.TextBuffer
}

func (c *mucRoomConfigComponent) newRoomConfigInfoPage() mucRoomConfigPage {
	p := &roomConfigInfoPage{}
	p.roomConfigPageBase = c.newConfigPage(pageConfigInfo, "MUCRoomConfigPageInfo", p, nil)
	return p
}
