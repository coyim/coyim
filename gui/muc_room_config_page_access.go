package gui

type roomConfigAccessPage struct {
	*roomConfigPageBase
}

func (c *mucRoomConfigComponent) newRoomConfigAccessPage() mucRoomConfigPage {
	p := &roomConfigAccessPage{}
	p.roomConfigPageBase = c.newConfigPage(pageConfigAccess, "MUCRoomConfigPageAccess", p, nil)

	return p
}

func (p *roomConfigAccessPage) isInvalid() bool {
	return p.roomConfigPageBase.isInvalid()
}
