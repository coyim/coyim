package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigFormFieldList struct {
	field *muc.RoomConfigFormField

	widget       gtki.Box      `gtk-widget:"room-config-field-box"`
	label        gtki.Label    `gtk-widget:"room-config-field-label"`
	list         gtki.ComboBox `gtk-widget:"room-config-field-list"`
	optionsModel gtki.ListStore
	options      map[string]int
}

func newRoomConfigFormFieldList(field *muc.RoomConfigFormField) *roomConfigFormFieldList {
	f := &roomConfigFormFieldList{field: field}

	f.initBuilder()
	f.initDefaults()
	return f
}

func (f *roomConfigFormFieldList) initBuilder() {
	builder := newBuilder("MUCRoomConfigFormFieldList")
	panicOnDevError(builder.bindObjects(f))

	f.optionsModel, _ = g.gtk.ListStoreNew(glibi.TYPE_STRING, glibi.TYPE_STRING)
	f.list.SetModel(f.optionsModel)
}

func (f *roomConfigFormFieldList) initDefaults() {
	f.optionsModel.Clear()
	f.label.SetText(f.field.Label)

	f.options = make(map[string]int)
	sf := f.field.Value.(muc.ConfigListSingleField)
	for index, o := range sf.Options() {
		iter := f.optionsModel.Append()
		_ = f.optionsModel.SetValue(iter, configWhoisOptionValueIndex, o)
		_ = f.optionsModel.SetValue(iter, configWhoisOptionLabelIndex, configOptionToFriendlyMessage(o))
		f.options[o] = index
	}

	f.activateOption(sf.CurrentValue())
}

func (f *roomConfigFormFieldList) activateOption(o string) {
	if index, ok := f.options[o]; ok {
		f.list.SetActive(index)
		return
	}
}

func (f *roomConfigFormFieldList) fieldWidget() gtki.Widget {
	return f.widget
}

func (f *roomConfigFormFieldList) fieldName() string {
	return f.field.Name
}

func (f *roomConfigFormFieldList) fieldLabel() string {
	return f.field.Label
}

func (f *roomConfigFormFieldList) fieldValue() interface{} {
	for o, index := range f.options {
		if index == f.list.GetActive() {
			return o
		}
	}
	return nil
}

func (f *roomConfigFormFieldList) refreshContent() {}
