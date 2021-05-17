package gui

type roomConfigPermissionsPage struct {
	*roomConfigPageBase
}

func (c *mucRoomConfigComponent) newRoomConfigPermissionsPage() mucRoomConfigPage {
	p := &roomConfigPermissionsPage{}
	p.roomConfigPageBase = c.newConfigPage(pageConfigPermissions, "MUCRoomConfigPagePermissions", p, nil)
	return p
}

const (
	configWhoisOptionValueIndex int = iota
	configWhoisOptionLabelIndex
)
