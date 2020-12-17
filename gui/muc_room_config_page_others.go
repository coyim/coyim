package gui

import "github.com/coyim/gotk3adapter/gtki"

type roomConfigOthersPage struct {
	*roomConfigPageBase

	configOthersBox     gtki.Box        `gtk-widget:"room-config-others-page"`
	roomMaxHistoryFetch gtki.SpinButton `gtk-widget:"room-maxhistoryfetch"`
	roomMaxOccupants    gtki.SpinButton `gtk-widget:"room-maxoccupants"`
	roomEnableLoggin    gtki.Switch     `gtk-widget:"room-enablelogging"`
}

func (c *mucRoomConfigComponent) newRoomConfigOthersPage() mucRoomConfigPage {
	p := &roomConfigOthersPage{}

	builder := newBuilder("MUCRoomConfigPageOthers")
	panicOnDevError(builder.bindObjects(p))

	p.roomConfigPageBase = c.newConfigPage(p.configOthersBox)

	p.initDefaultValues()

	return p
}

func (p *roomConfigOthersPage) initDefaultValues() {
	setEntryText(p.roomMaxHistoryFetch, p.form.MaxHistoryFetch)
	setEntryText(p.roomMaxOccupants, p.form.MaxOccupantsNumber.CurrentValue())
	setSwitchActive(p.roomEnableLoggin, p.form.Logged)
}

func (p *roomConfigOthersPage) collectData() {
	p.form.MaxHistoryFetch = getEntryText(p.roomMaxHistoryFetch)
	p.form.MaxOccupantsNumber.UpdateValue(getEntryText(p.roomMaxOccupants))
	p.form.Logged = getSwitchActive(p.roomEnableLoggin)
}
