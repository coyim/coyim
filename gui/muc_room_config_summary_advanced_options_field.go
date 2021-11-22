package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigSummaryAdvancedOptionsField struct {
	advancedOptions []hasRoomConfigFormField

	widget                    gtki.Box     `gtk-widget:"room-config-field-box"`
	fieldListValueButtonImage gtki.Image   `gtk-widget:"room-config-field-list-button-image"`
	advancedOptionsContent    gtki.Box     `gtk-widget:"room-config-field-advanced-options-box"`
	advancedOptionsList       gtki.ListBox `gtk-widget:"room-config-field-advanced-options-content-box"`
}

func newAdvancedOptionSummaryField(advancedOptions []hasRoomConfigFormField) hasRoomConfigFormField {
	field := &roomConfigSummaryAdvancedOptionsField{
		advancedOptions: advancedOptions,
	}

	field.initBuilder()
	field.initDefaults()

	return field
}

func (f *roomConfigSummaryAdvancedOptionsField) initBuilder() {
	builder := newBuilder("MUCRoomConfigSummaryAdvancedOptionsField")
	panicOnDevError(builder.bindObjects(f))
	builder.ConnectSignals(map[string]interface{}{
		"on_show_list": f.showAdvancedOptions,
	})
}

func (f *roomConfigSummaryAdvancedOptionsField) fieldWidgets() (widgets []gtki.Widget) {
	for _, ff := range f.advancedOptions {
		widgets = append(widgets, ff.fieldWidget())
	}
	return
}

func (f *roomConfigSummaryAdvancedOptionsField) initDefaults() {
	fields := addSeparatorsBetweenWidgets(f.fieldWidgets())
	for _, ff := range fields {
		f.advancedOptionsList.Add(ff)
	}
}

func (f *roomConfigSummaryAdvancedOptionsField) showAdvancedOptions() {
	if f.advancedOptionsContent.IsVisible() {
		f.fieldListValueButtonImage.SetFromIconName("pan-down-symbolic", gtki.ICON_SIZE_MENU)
		f.advancedOptionsContent.Hide()
	} else {
		f.fieldListValueButtonImage.SetFromIconName("pan-up-symbolic", gtki.ICON_SIZE_MENU)
		f.advancedOptionsContent.Show()
	}
}

func (f *roomConfigSummaryAdvancedOptionsField) fieldWidget() gtki.Widget {
	return f.widget
}

func (f *roomConfigSummaryAdvancedOptionsField) refreshContent() {
	for _, field := range f.advancedOptions {
		field.refreshContent()
	}
}

func (f *roomConfigSummaryAdvancedOptionsField) updateFieldValue() {}

func (f *roomConfigSummaryAdvancedOptionsField) isValid() bool {
	return true
}

func (f *roomConfigSummaryAdvancedOptionsField) showValidationErrors() {}

// fieldKey implements the hasRoomConfigFormField interface
func (f *roomConfigSummaryAdvancedOptionsField) fieldKey() muc.RoomConfigFieldType {
	return muc.RoomConfigFieldUnexpected
}

// focusWidget implements the hasRoomConfigFormField interface
func (f *roomConfigSummaryAdvancedOptionsField) focusWidget() focusable {
	return f.widget
}
