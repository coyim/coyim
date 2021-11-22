package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigSummaryField struct {
	fieldTexts roomConfigFieldTextInfo
	fieldType  muc.RoomConfigFieldType
	fieldValue muc.HasRoomConfigFormFieldValue

	widget                    gtki.Box        `gtk-widget:"room-config-field-box"`
	field                     gtki.ListBoxRow `gtk-widget:"room-config-field"`
	fieldLabel                gtki.Label      `gtk-widget:"room-config-field-label"`
	fieldValueLabel           gtki.Label      `gtk-widget:"room-config-field-value"`
	fieldTextMultiValue       gtki.Label      `gtk-widget:"room-config-field-text-multi-value"`
	fieldListContent          gtki.Box        `gtk-widget:"room-config-field-list-content"`
	fieldListValueButton      gtki.Button     `gtk-widget:"room-config-field-list-button"`
	fieldListValueButtonImage gtki.Image      `gtk-widget:"room-config-field-list-button-image"`
	fieldListValues           gtki.TreeView   `gtk-widget:"room-config-field-list-values-tree"`

	listModel gtki.ListStore
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
	builder.ConnectSignals(map[string]interface{}{
		"on_show_list": f.onShowList,
	})
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
	case *muc.RoomConfigFieldListMultiValue:
		f.handleListMultiFieldValue(v.Value())
	}
}

func (f *roomConfigSummaryField) handleTextFieldValue(value string) {
	setLabelText(f.fieldValueLabel, summaryAssignedValueText(value))

	switch f.fieldType {
	case muc.RoomConfigFieldDescription:
		f.handleTextMultiFieldValue(value)
	case muc.RoomConfigFieldLanguage:
		setLabelText(f.fieldValueLabel, supportedLanguageDescription(value))
	case muc.RoomConfigFieldPassword:
		setLabelText(f.fieldValueLabel, summaryPasswordText(value))
	}
}

// handleTextMultiFieldValue MUST be called from the UI thread
func (f *roomConfigSummaryField) handleTextMultiFieldValue(value string) {
	if value != "" {
		setLabelText(f.fieldTextMultiValue, summaryAssignedValueText(value))
		f.fieldTextMultiValue.Show()
		f.fieldValueLabel.Hide()
		return
	}

	setLabelText(f.fieldValueLabel, summaryAssignedValueText(value))
	f.fieldTextMultiValue.Hide()
	f.fieldValueLabel.Show()
}

func (f *roomConfigSummaryField) handleListMultiFieldValue(value []string) {
	f.listModel, _ = g.gtk.ListStoreNew(glibi.TYPE_STRING)
	f.fieldListValues.SetModel(f.listModel)

	totalItemsLabel := summaryListTotalText(len(value))
	if f.fieldType == muc.RoomConfigFieldPresenceBroadcast {
		totalItemsLabel = summaryRolesListText(len(value))
	}
	setLabelText(f.fieldValueLabel, totalItemsLabel)
	f.fieldListValueButton.SetVisible(len(value) > 0)

	f.initListContent(value)
}

func (f *roomConfigSummaryField) initListContent(items []string) {
	f.listModel.Clear()

	for _, j := range items {
		iter := f.listModel.Append()
		f.listModel.SetValue(iter, 0, configOptionToFriendlyMessage(j, j))
	}
}

func (f *roomConfigSummaryField) onShowList() {
	summaryListHideOrShow(f.fieldListValues, f.fieldListValueButtonImage, f.fieldListContent)
}

func (f *roomConfigSummaryField) fieldWidget() gtki.Widget {
	return f.widget
}

func (f *roomConfigSummaryField) refreshContent() {
	f.handleFieldValue()
}

func (f *roomConfigSummaryField) updateFieldValue() {}

func (f *roomConfigSummaryField) isValid() bool {
	return true
}

func (f *roomConfigSummaryField) showValidationErrors() {}

// fieldKey implements the hasRoomConfigFormField interface
func (f *roomConfigSummaryField) fieldKey() muc.RoomConfigFieldType {
	return muc.RoomConfigFieldUnexpected
}

func summaryPasswordText(v string) string {
	if v != "" {
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

func summaryListTotalText(total int) string {
	switch {
	case total == 1:
		return i18n.Local("One result")
	case total > 0:
		return i18n.Localf("%d results", total)
	}
	return i18n.Local("None")
}

func summaryRolesListText(total int) string {
	switch {
	case total == 1:
		return i18n.Local("One role")
	case total > 0:
		return i18n.Localf("%d roles", total)
	}
	return i18n.Local("No roles")
}

func summaryTotalPositionsText(total int) string {
	switch {
	case total == 1:
		return i18n.Local("One account")
	case total > 0:
		return i18n.Localf("%d accounts", total)
	}
	return i18n.Local("No accounts")
}

func summaryListHideOrShow(list gtki.TreeView, toggleButtonImage gtki.Image, container gtki.Box) {
	if list.IsVisible() {
		toggleButtonImage.SetFromIconName("pan-down-symbolic", gtki.ICON_SIZE_MENU)
		container.SetVisible(false)
	} else {
		toggleButtonImage.SetFromIconName("pan-up-symbolic", gtki.ICON_SIZE_MENU)
		container.SetVisible(true)
	}
}

// focusWidget implements the hasRoomConfigFormField interface
func (f *roomConfigSummaryField) focusWidget() focusable {
	return f.widget
}
