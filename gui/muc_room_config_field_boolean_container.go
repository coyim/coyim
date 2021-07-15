package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldBooleanContainer struct {
	fields []hasRoomConfigFormField

	widget  gtki.Box     `gtk-widget:"room-config-field-box"`
	content gtki.ListBox `gtk-widget:"room-config-boolean-fields-content"`
}

func newRoomConfigFormFieldBooleanContainer(f []hasRoomConfigFormField) hasRoomConfigFormField {
	field := &roomConfigFormFieldBooleanContainer{
		fields: f,
	}

	field.initBuilder()
	field.initDefaults()

	return field
}

func (fc *roomConfigFormFieldBooleanContainer) initBuilder() {
	builder := newBuilder("MUCRoomConfigFormFieldBooleanContainer")
	panicOnDevError(builder.bindObjects(fc))
}

func (fc *roomConfigFormFieldBooleanContainer) initDefaults() {
	fc.content.Add(fc.fields[0].fieldWidget())
	for _, f := range fc.fields[1:] {
		fc.content.Add(createSeparator(gtki.HorizontalOrientation))
		fc.content.Add(f.fieldWidget())
	}
}

func (fc *roomConfigFormFieldBooleanContainer) fieldWidget() gtki.Widget {
	return fc.widget
}

// refreshContent MUST NOT be called from the UI thread
func (fc *roomConfigFormFieldBooleanContainer) refreshContent() {}

// updateFieldValue MUST be called from the UI thread
func (fc *roomConfigFormFieldBooleanContainer) updateFieldValue() {
	for _, f := range fc.fields {
		f.updateFieldValue()
	}
}

// isValid implements the hasRoomConfigFormField interface
func (fc *roomConfigFormFieldBooleanContainer) isValid() bool {
	return true
}

// showValidationErrors implements the hasRoomConfigFormField interface
func (fc *roomConfigFormFieldBooleanContainer) showValidationErrors() {}

// fieldKey implements the hasRoomConfigFormField interface
func (fc *roomConfigFormFieldBooleanContainer) fieldKey() muc.RoomConfigFieldType {
	return muc.RoomConfigFieldUnexpected
}

func createSeparator(orientation gtki.Orientation) gtki.Separator {
	sep, _ := g.gtk.SeparatorNew(orientation)
	return sep
}
