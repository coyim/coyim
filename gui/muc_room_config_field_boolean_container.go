package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldBooleanContainer struct {
	fields []*roomConfigFormFieldBoolean

	widget  gtki.Box     `gtk-widget:"room-config-field-box"`
	content gtki.ListBox `gtk-widget:"room-config-boolean-fields-content"`
}

func newRoomConfigFormFieldBooleanContainer(f []*roomConfigFormFieldBoolean) hasRoomConfigFormField {
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
	fc.content.Add(fc.fields[0].widget)
	for _, f := range fc.fields[1:] {
		fc.content.Add(createSeparator(gtki.HorizontalOrientation))
		fc.content.Add(f.widget)
	}
}

func (fc *roomConfigFormFieldBooleanContainer) fieldWidget() gtki.Widget {
	return fc.widget
}

// refreshContent MUST NOT be called from the UI thread
func (fc *roomConfigFormFieldBooleanContainer) refreshContent() {}

// collectFieldValue MUST be called from the UI thread
func (fc *roomConfigFormFieldBooleanContainer) collectFieldValue() {
	for _, f := range fc.fields {
		f.collectFieldValue()
	}
}

func createSeparator(orientation gtki.Orientation) gtki.Separator {
	sep, _ := g.gtk.SeparatorNew(orientation)
	return sep
}
