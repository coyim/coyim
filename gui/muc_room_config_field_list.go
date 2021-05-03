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

	list gtki.ComboBox `gtk-widget:"room-config-field-list"`

	optionsModel gtki.ListStore
	options      map[string]int
}

func newRoomConfigFormFieldList(f muc.HasRoomConfigFormField) hasRoomConfigFormField {
	field := &roomConfigFormFieldList{}
	field.roomConfigFormField = newRoomConfigFormField(f, "MUCRoomConfigFormFieldList")

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
	f.options = map[string]int{}

	if lf, ok := f.field.Value().(muc.ConfigListSingleField); ok {
		for index, o := range lf.Options() {
			iter := f.optionsModel.Append()

			_ = f.optionsModel.SetValue(iter, roomConfigFieldListOptionValueIndex, o)
			_ = f.optionsModel.SetValue(iter, roomConfigFieldListOptionLabelIndex, configOptionToFriendlyMessage(o))

			f.options[o] = index
		}

		f.activateOption(lf.CurrentValue())
	}
}

// activateOption MUST be called from the UI thread
func (f *roomConfigFormFieldList) activateOption(o string) {
	if index, ok := f.options[o]; ok {
		f.list.SetActive(index)
		return
	}
}

// fieldValue MUST be called from the UI thread
func (f *roomConfigFormFieldList) fieldValue() interface{} {
	for o, index := range f.options {
		if index == f.list.GetActive() {
			return o
		}
	}
	return nil
}
