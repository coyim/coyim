package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldListEntry struct {
	*roomConfigFormField
	value *muc.RoomConfigFieldListValue

	list  gtki.ComboBox `gtk-widget:"room-config-field-list"`
	entry gtki.Entry    `gtk-widget:"room-config-field-list-entry"`

	optionsModel gtki.ListStore
	options      map[string]int
}

func newRoomConfigFormFieldListEntry(fieldInfo roomConfigFieldTextInfo, value *muc.RoomConfigFieldListValue) hasRoomConfigFormField {
	field := &roomConfigFormFieldListEntry{value: value}
	field.roomConfigFormField = newRoomConfigFormField(fieldInfo, "MUCRoomConfigFormFieldListEntry")

	panicOnDevError(field.builder.bindObjects(field))

	field.optionsModel, _ = g.gtk.ListStoreNew(
		// the option value
		glibi.TYPE_STRING,
		// the option display label
		glibi.TYPE_STRING,
	)

	field.list.SetModel(field.optionsModel)

	field.initOptions()

	return field
}

func (f *roomConfigFormFieldListEntry) initOptions() {
	f.options = map[string]int{}

	for index, o := range f.value.Options() {
		iter := f.optionsModel.Append()

		_ = f.optionsModel.SetValue(iter, roomConfigFieldListOptionValueIndex, o.Value)
		_ = f.optionsModel.SetValue(iter, roomConfigFieldListOptionLabelIndex, configOptionToFriendlyMessage(o.Value, o.Label))

		f.options[o.Value] = index
	}

	f.activateOption(f.value.Selected())
}

// activateOption MUST be called from the UI thread
func (f *roomConfigFormFieldListEntry) activateOption(v string) {
	f.entry.SetText(v)
}

// collectFieldValue MUST be called from the UI thread
func (f *roomConfigFormFieldListEntry) collectFieldValue() {
	f.value.SetSelected("")
	for o, index := range f.options {
		if index == f.list.GetActive() {
			f.value.SetSelected(o)
			return
		}
	}
}
