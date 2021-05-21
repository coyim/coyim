package gui

type roomConfigPermissionsPage struct {
	*roomConfigPageBase
}

func (c *mucRoomConfigComponent) newRoomConfigPermissionsPage() mucRoomConfigPage {
	p := &roomConfigPermissionsPage{}
	return p
}

const (
	configWhoisOptionValueIndex int = iota
	configWhoisOptionLabelIndex
)
