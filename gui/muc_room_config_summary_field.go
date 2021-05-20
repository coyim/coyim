package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigSummaryField struct {
	*roomConfigFormField
	fieldTexts roomConfigFieldTextInfo

	widget     gtki.Box        `gtk-widget:"room-config-field-box"`
	field      gtki.ListBoxRow `gtk-widget:"room-config-field"`
	fieldLabel gtki.Label      `gtk-widget:"room-config-field-label"`
}

func newRoomConfigSummaryField(fieldTexts roomConfigFieldTextInfo, fieldTypeValue interface{}) *roomConfigSummaryField {
	field := &roomConfigSummaryField{fieldTexts: fieldTexts}

	field.initBuilder()
	field.initDefaults()

	return field
}

func (f *roomConfigSummaryField) initBuilder() {
	builder := newBuilder("MUCRoomConfigSummaryField")
	panicOnDevError(builder.bindObjects(f))
}

func (f *roomConfigSummaryField) initDefaults() {
	f.fieldLabel.SetText(f.fieldTexts.displayLabel)
}

func (f *roomConfigSummaryField) fieldWidget() gtki.Widget {
	return f.field
}
