package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

const (
	roomConfigFieldListOptionValueIndex int = iota
	roomConfigFieldListOptionLabelIndex
)

type roomConfigFormFieldList struct {
	*roomConfigFormField
	value *muc.RoomConfigFieldListValue

	list gtki.ComboBox `gtk-widget:"room-config-field-list"`

	optionsModel gtki.ListStore
}

func newRoomConfigFormFieldList(ft muc.RoomConfigFieldType, fieldInfo roomConfigFieldTextInfo, value *muc.RoomConfigFieldListValue, onShowValidationErrors func(), onHideValidationErrors func()) hasRoomConfigFormField {
	field := &roomConfigFormFieldList{value: value}
	field.roomConfigFormField = newRoomConfigFormField(ft, fieldInfo, "MUCRoomConfigFormFieldList", onShowValidationErrors, onHideValidationErrors)

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

func (f *roomConfigFormFieldList) initOptions() {
	for _, o := range f.value.Options() {
		iter := f.optionsModel.Append()

		_ = f.optionsModel.SetValue(iter, roomConfigFieldListOptionValueIndex, o.Value)
		_ = f.optionsModel.SetValue(iter, roomConfigFieldListOptionLabelIndex, configOptionToFriendlyMessage(o.Value, o.Label))
	}

	f.activateOption(f.value.Selected())
}

// activateOption MUST be called from the UI thread
func (f *roomConfigFormFieldList) activateOption(o string) {
	iter, ok := f.optionsModel.GetIterFirst()
	idx := 0
	for ok {
		if getStringValueFromModel(f.optionsModel, iter, roomConfigFieldListOptionValueIndex) == o {
			f.list.SetActive(idx)
			return
		}
		idx++
		ok = f.optionsModel.IterNext(iter)
	}
}

// updateFieldValue MUST be called from the UI thread
func (f *roomConfigFormFieldList) updateFieldValue() {
	iter, _ := f.list.GetActiveIter()
	f.value.SetSelected(getStringValueFromModel(f.optionsModel, iter, roomConfigFieldListOptionValueIndex))
}

func getStringValueFromModel(model gtki.ListStore, iter gtki.TreeIter, columnID int) string {
	ov, _ := model.GetValue(iter, columnID)
	s, _ := ov.GetString()

	return s
}
