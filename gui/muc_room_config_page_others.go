package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigOthersPage struct {
	*roomConfigPageBase

	configOthersBox     gtki.Box          `gtk-widget:"room-config-others-page"`
	roomMaxHistoryFetch gtki.ComboBoxText `gtk-widget:"room-maxhistoryfetch"`
	roomMaxOccupants    gtki.ComboBoxText `gtk-widget:"room-maxoccupants"`
	roomEnableLoggin    gtki.Switch       `gtk-widget:"room-enablelogging"`
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
	setSwitchActive(p.roomEnableLoggin, p.form.Logged)

	p.initializeAvailableOptions(p.roomMaxHistoryFetch, p.form.MaxHistoryFetch.Options())
	p.initializeAvailableOptions(p.roomMaxOccupants, p.form.MaxOccupantsNumber.Options())
}

func (p *roomConfigOthersPage) initializeAvailableOptions(combo gtki.ComboBoxText, options []string) {
	combo.SetSensitive(true)
	for _, o := range options {
		combo.AppendText(configOptionToFriendlyMessage(o))
	}
	combo.SetActive(0)
}

func (p *roomConfigOthersPage) collectData() {
	p.form.Logged = getSwitchActive(p.roomEnableLoggin)
}
