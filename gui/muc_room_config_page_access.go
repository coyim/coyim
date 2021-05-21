package gui

type roomConfigAccessPage struct {
	*roomConfigPageBase
}

func (c *mucRoomConfigComponent) newRoomConfigAccessPage() mucRoomConfigPage {
	p := &roomConfigAccessPage{}
	return p
}

func (p *roomConfigAccessPage) isValid() bool {
	return p.roomConfigPageBase.isValid()
}
