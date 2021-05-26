package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigSummaryField struct {
	fieldTexts roomConfigFieldTextInfo
	fieldType  muc.RoomConfigFieldType
	fieldValue muc.HasRoomConfigFormFieldValue

	widget                gtki.Box            `gtk-widget:"room-config-field-box"`
	field                 gtki.ListBoxRow     `gtk-widget:"room-config-field"`
	fieldLabel            gtki.Label          `gtk-widget:"room-config-field-label"`
	fieldValueLabel       gtki.Label          `gtk-widget:"room-config-field-value"`
	fieldTextMultiContent gtki.ScrolledWindow `gtk-widget:"room-config-field-text-area"`
	fieldTextMultiValue   gtki.TextView       `gtk-widget:"room-config-field-text-area-value"`
}

func newRoomConfigSummaryField(fieldType muc.RoomConfigFieldType, fieldTexts roomConfigFieldTextInfo, value muc.HasRoomConfigFormFieldValue) hasRoomConfigFormField {
	field := &roomConfigSummaryField{
		fieldTexts: fieldTexts,
		fieldType:  fieldType,
		fieldValue: value,
	}

	field.initBuilder()
	field.initDefaults()
	field.handleFieldValue()

	return field
}

func (f *roomConfigSummaryField) initBuilder() {
	builder := newBuilder("MUCRoomConfigSummaryField")
	panicOnDevError(builder.bindObjects(f))
}

func (f *roomConfigSummaryField) initDefaults() {
	f.fieldLabel.SetText(f.fieldTexts.summaryLabel)
}

func (f *roomConfigSummaryField) handleFieldValue() {
	switch v := f.fieldValue.(type) {
	case *muc.RoomConfigFieldTextValue:
		f.handleTextFieldValue(v.Text())
	case *muc.RoomConfigFieldTextMultiValue:
		f.handleTextMultiFieldValue(v.Text())
	case *muc.RoomConfigFieldBooleanValue:
		f.handleTextFieldValue(summaryYesOrNoText(v.Boolean()))
	case *muc.RoomConfigFieldListValue:
		f.handleTextFieldValue(configOptionToFriendlyMessage(v.Selected(), v.Selected()))
	}
}

func (f *roomConfigSummaryField) handleTextFieldValue(value string) {
	switch f.fieldType {
	case muc.RoomConfigFieldLanguage:
		setLabelText(f.fieldValueLabel, supportedLanguageDescription(value))
	case muc.RoomConfigFieldPassword:
		setLabelText(f.fieldValueLabel, summaryPasswordText(value == ""))
	}
	setLabelText(f.fieldValueLabel, summaryAssignedValueText(value))
}

func (f *roomConfigSummaryField) handleTextMultiFieldValue(value string) {
	if value != "" {
		setTextViewText(f.fieldTextMultiValue, summaryAssignedValueText(value))
		f.fieldTextMultiContent.Show()
		f.fieldValueLabel.SetVisible(false)
		return
	}

	setLabelText(f.fieldValueLabel, summaryAssignedValueText(value))
	f.fieldTextMultiContent.Hide()
	f.fieldValueLabel.SetVisible(true)
}

func (f *roomConfigSummaryField) fieldWidget() gtki.Widget {
	return f.widget
}

func (f *roomConfigSummaryField) refreshContent() {
	f.handleFieldValue()
}

func (f *roomConfigSummaryField) collectFieldValue() {}

func (f *roomConfigSummaryField) isValid() bool {
	return true
}

func (f *roomConfigSummaryField) showValidationErrors() {}

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
