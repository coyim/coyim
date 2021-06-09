package gui

import (
	"strconv"

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
	field.list.SetIDColumn(roomConfigFieldListOptionValueIndex)
	field.list.SetEntryTextColumn(roomConfigFieldListOptionLabelIndex)

	field.initOptions()

	return field
}

func (f *roomConfigFormFieldListEntry) initOptions() {
	for _, o := range f.value.Options() {
		iter := f.optionsModel.Append()

		_ = f.optionsModel.SetValue(iter, roomConfigFieldListOptionValueIndex, o.Value)
		_ = f.optionsModel.SetValue(iter, roomConfigFieldListOptionLabelIndex, configOptionToFriendlyMessage(o.Value, o.Label))
	}

	f.activateOption(f.value.Selected())
}

// activateOption MUST be called from the UI thread
func (f *roomConfigFormFieldListEntry) activateOption(v string) {
	f.entry.SetText(v)
}

// updateFieldValue MUST be called from the UI thread
func (f *roomConfigFormFieldListEntry) updateFieldValue() {
	f.value.SetSelected(f.currentValue())
}

// isValid implements the hasRoomConfigFormField interface
func (f *roomConfigFormFieldListEntry) isValid() bool {
	v := f.currentValue()
	if v != "" {
		_, err := strconv.Atoi(v)
		return err == nil
	}
	return true
}

// showValidationErrors implements the hasRoomConfigFormField interface
func (f *roomConfigFormFieldListEntry) showValidationErrors() {
	f.entry.GrabFocus()
}

func (f *roomConfigFormFieldListEntry) currentValue() string {
	iter, err := f.list.GetActiveIter()
	if err == nil {
		return getStringValueFromModel(f.optionsModel, iter, roomConfigFieldListOptionValueIndex)
	}
	return getEntryText(f.entry)
}
