package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigSummaryFieldContainer struct {
	fields []hasRoomConfigFormField

	widget  gtki.Box     `gtk-widget:"room-config-field-box"`
	content gtki.ListBox `gtk-widget:"room-config-fields-content"`
}

func newRoomConfigSummaryFieldContainer(f []hasRoomConfigFormField) hasRoomConfigFormField {
	field := &roomConfigSummaryFieldContainer{
		fields: f,
	}

	field.initBuilder()
	field.initDefaults()

	return field
}

func (fc *roomConfigSummaryFieldContainer) initBuilder() {
	builder := newBuilder("MUCRoomConfigSummaryFieldContainer")
	panicOnDevError(builder.bindObjects(fc))
}

func (fc *roomConfigSummaryFieldContainer) initDefaults() {
	fc.content.Add(fc.fields[0].fieldWidget())
	for _, f := range fc.fields[1:] {
		fc.content.Add(createSeparator(gtki.HorizontalOrientation))
		fc.content.Add(f.fieldWidget())
	}
}

func (fc *roomConfigSummaryFieldContainer) fieldWidget() gtki.Widget {
	return fc.widget
}

// refreshContent MUST NOT be called from the UI thread
func (fc *roomConfigSummaryFieldContainer) refreshContent() {
	for _, f := range fc.fields {
		f.refreshContent()
	}
}

// updateFieldValue MUST be called from the UI thread
func (fc *roomConfigSummaryFieldContainer) updateFieldValue() {
	for _, f := range fc.fields {
		f.updateFieldValue()
	}
}

// isValid implements the hasRoomConfigFormField interface
func (fc *roomConfigSummaryFieldContainer) isValid() bool {
	return true
}

// showValidationErrors implements the hasRoomConfigFormField interface
func (fc *roomConfigSummaryFieldContainer) showValidationErrors() {}
