package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldAdvancedOptionsContainer struct {
	fields []hasRoomConfigFormField

	widget  gtki.Box `gtk-widget:"room-config-field-box"`
	content gtki.Box `gtk-widget:"room-config-field-advanced-options-content-box"`
}

func newRoomConfigFormFieldAdvancedOptionsContainer(f []hasRoomConfigFormField) hasRoomConfigFormField {
	field := &roomConfigFormFieldAdvancedOptionsContainer{
		fields: f,
	}

	field.initBuilder()
	field.initDefaults()

	return field
}

func (fc *roomConfigFormFieldAdvancedOptionsContainer) initBuilder() {
	builder := newBuilder("MUCRoomConfigFormFieldAdvancedOptionsContainer")
	panicOnDevError(builder.bindObjects(fc))
}

func (fc *roomConfigFormFieldAdvancedOptionsContainer) initDefaults() {
	for _, f := range fc.fields {
		fc.content.Add(f.fieldWidget())
	}
}

func (fc *roomConfigFormFieldAdvancedOptionsContainer) fieldWidget() gtki.Widget {
	return fc.widget
}

// refreshContent MUST NOT be called from the UI thread
func (fc *roomConfigFormFieldAdvancedOptionsContainer) refreshContent() {}

// updateFieldValue MUST be called from the UI thread
func (fc *roomConfigFormFieldAdvancedOptionsContainer) updateFieldValue() {
	for _, f := range fc.fields {
		f.updateFieldValue()
	}
}

// isValid implements the hasRoomConfigFormField interface
func (fc *roomConfigFormFieldAdvancedOptionsContainer) isValid() bool {
	return true
}

// showValidationErrors implements the hasRoomConfigFormField interface
func (fc *roomConfigFormFieldAdvancedOptionsContainer) showValidationErrors() {}
