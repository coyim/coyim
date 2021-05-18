package gui

import (
	"github.com/coyim/coyim/session/muc"
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

	p.initDefaultValues()
	return p
}

func (p *roomConfigOthersPage) initDefaultValues() {
	for _, f := range p.form.Fields {
		field, err := roomConfigFormFieldFactory(muc.RoomConfigFieldUnexpected, newRoomConfigFieldTextInfo(f.Label, f.Description), f.ValueType())
		if err != nil {
			p.log.WithField("field", f.Name).WithError(err).Error("Room configuration form field not supported")
			continue
		}

		p.addField(field)
		p.doAfterRefresh.add(field.refreshContent)
	}
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
