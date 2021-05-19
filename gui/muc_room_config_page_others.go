package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigOthersPage struct {
	*roomConfigPageBase
	fields []hasRoomConfigFormField

	roomUnknowFieldsBox gtki.Box `gtk-widget:"room-config-unknow-fields-box"`
}

func (c *mucRoomConfigComponent) newRoomConfigOthersPage() mucRoomConfigPage {
	p := &roomConfigOthersPage{}
	p.roomConfigPageBase = c.newConfigPage(pageConfigOthers, "MUCRoomConfigPageOthers", p, nil)

	return p
}

// isInvalid MUST be called from the UI thread
func (p *roomConfigOthersPage) isValid() bool {
	return p.roomConfigPageBase.isValid()
}

// collectData MUST be called from the UI thread
func (p *roomConfigOthersPage) collectData() {
	for _, f := range p.fields {
		f.collectFieldValue()
	}
}

// addField MUST BE called from the UI thread
func (p *roomConfigOthersPage) addField(f hasRoomConfigFormField) {
	p.fields = append(p.fields, f)
	p.roomUnknowFieldsBox.PackStart(f.fieldWidget(), true, false, 0)
}
