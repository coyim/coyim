package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigSummaryField struct {
	*roomConfigFormField
	fieldTexts roomConfigFieldTextInfo

	widget                gtki.Box            `gtk-widget:"room-config-field-box"`
	field                 gtki.ListBoxRow     `gtk-widget:"room-config-field"`
	fieldLabel            gtki.Label          `gtk-widget:"room-config-field-label"`
	fieldValue            gtki.Label          `gtk-widget:"room-config-field-value"`
	fieldTextMultiContent gtki.ScrolledWindow `gtk-widget:"room-config-field-text-area"`
	fieldTextMultiValue   gtki.TextView       `gtk-widget:"room-config-field-text-area-value"`
}

func newRoomConfigSummaryField(fieldType muc.RoomConfigFieldType, fieldTexts roomConfigFieldTextInfo, fieldTypeValue interface{}) *roomConfigSummaryField {
	field := &roomConfigSummaryField{fieldTexts: fieldTexts}

	field.initBuilder()
	field.initDefaults()
	field.handleFieldValue(fieldType, fieldTypeValue)

	return field
}

func (f *roomConfigSummaryField) initBuilder() {
	builder := newBuilder("MUCRoomConfigSummaryField")
	panicOnDevError(builder.bindObjects(f))
}

func (f *roomConfigSummaryField) initDefaults() {
	f.fieldLabel.SetText(f.fieldTexts.displayLabel)
}

func (f *roomConfigSummaryField) fieldWidget() gtki.Widget {
	return f.field
}

func (f *roomConfigSummaryField) handleFieldValue(fieldType muc.RoomConfigFieldType, fieldTypeValue interface{}) {
	switch v := fieldTypeValue.(type) {
	case *muc.RoomConfigFieldTextValue:
		f.handleTextFieldValue(fieldType, v.Text())
	case *muc.RoomConfigFieldTextMultiValue:
		f.handleTextMultiFieldValue(fieldType, v.Text())
	case *muc.RoomConfigFieldBooleanValue:
		f.handleTextFieldValue(fieldType, summaryYesOrNoText(v.Boolean()))
	}
}

func (f *roomConfigSummaryField) handleTextFieldValue(ft muc.RoomConfigFieldType, value string) {
	switch ft {
	case muc.RoomConfigFieldLanguage:
		setLabelText(f.fieldValue, supportedLanguageDescription(value))
	case muc.RoomConfigFieldPassword:
		setLabelText(f.fieldValue, summaryPasswordText(value == ""))
	}
	setLabelText(f.fieldValue, summaryAssignedValueText(value))
}

func (f *roomConfigSummaryField) handleTextMultiFieldValue(ft muc.RoomConfigFieldType, value string) {
	if value != "" {
		setTextViewText(f.fieldTextMultiValue, summaryAssignedValueText(value))
		f.fieldTextMultiContent.Show()
		f.fieldValue.SetVisible(false)
		return
	}

	setLabelText(f.fieldValue, summaryAssignedValueText(value))
	f.fieldTextMultiContent.Hide()
	f.fieldValue.SetVisible(true)
}

func summaryPasswordText(v bool) string {
	if v {
		return i18n.Local("**********")
	}
	return i18n.Local("Not assigned")
}

func summaryYesOrNoText(v bool) string {
	if v {
		return i18n.Local("Yes")
	}
	return i18n.Local("No")
}

func summaryAssignedValueText(label string) string {
	if label != "" {
		return label
	}
	return i18n.Local("Not assigned")
}
