package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldListEntry struct {
	*roomConfigFormField
	value *muc.RoomConfigFieldListValue

	fieldOptions gtki.ComboBox `gtk-widget:"room-config-field-list"`
	entry        gtki.Entry    `gtk-widget:"room-config-field-list-entry"`

	optionsModel gtki.ListStore
}

func newRoomConfigFormFieldListEntry(ft muc.RoomConfigFieldType, fieldInfo roomConfigFieldTextInfo, value *muc.RoomConfigFieldListValue, onShowValidationErrors func(), onHideValidationErrors func()) hasRoomConfigFormField {
	field := &roomConfigFormFieldListEntry{value: value}
	field.roomConfigFormField = newRoomConfigFormField(ft, fieldInfo, "MUCRoomConfigFormFieldListEntry", onShowValidationErrors, onHideValidationErrors)

	panicOnDevError(field.builder.bindObjects(field))
	field.builder.ConnectSignals(map[string]interface{}{
		"on_field_entry_change": field.onFieldEntryChanged,
	})

	field.optionsModel, _ = g.gtk.ListStoreNew(
		// the option value
		glibi.TYPE_STRING,
		// the option display label
		glibi.TYPE_STRING,
	)

	field.fieldOptions.SetModel(field.optionsModel)
	field.fieldOptions.SetIDColumn(roomConfigFieldListOptionValueIndex)
	field.fieldOptions.SetEntryTextColumn(roomConfigFieldListOptionLabelIndex)

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

func (f *roomConfigFormFieldListEntry) onFieldEntryChanged() {
	if f.isValid() {
		f.hideValidationErrors()
		return
	}

	f.showValidationErrors()
}

// updateFieldValue MUST be called from the UI thread
func (f *roomConfigFormFieldListEntry) updateFieldValue() {
	f.value.SetSelected(f.currentValue())
}

// isValid is the implementation for hasRoomConfigFormField interface
func (f *roomConfigFormFieldListEntry) isValid() bool {
	v := f.currentValue()

	validator, ok := roomConfigFieldValidator[f.field]
	if ok {
		return validator(v)
	}

	return true
}

const entryDefaultCursorPosition = -1

// refreshContent is the implementation for hasRoomConfigFormField interface
func (f *roomConfigFormFieldListEntry) refreshContent() {
	f.entry.SetPosition(entryDefaultCursorPosition)
}

func (f *roomConfigFormFieldListEntry) currentValue() string {
	iter, err := f.fieldOptions.GetActiveIter()
	if err == nil {
		return getStringValueFromModel(f.optionsModel, iter, roomConfigFieldListOptionValueIndex)
	}
	return getEntryText(f.entry)
}

// focusWidget is the implementation for hasRoomConfigFormField interface
func (f *roomConfigFormFieldListEntry) focusWidget() focusable {
	return f.entry
}
