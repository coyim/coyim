package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigSummaryOccupantField struct {
	affiliation data.Affiliation
	occupants   map[data.Affiliation][]*muc.RoomOccupantItem

	widget                    gtki.Box        `gtk-widget:"room-config-field-box"`
	field                     gtki.ListBoxRow `gtk-widget:"room-config-field"`
	fieldLabel                gtki.Label      `gtk-widget:"room-config-field-label"`
	fieldValueLabel           gtki.Label      `gtk-widget:"room-config-field-value"`
	fieldListContent          gtki.Box        `gtk-widget:"room-config-field-list-content"`
	fieldListValueButton      gtki.Button     `gtk-widget:"room-config-field-list-button"`
	fieldListValueButtonImage gtki.Image      `gtk-widget:"room-config-field-list-button-image"`
	fieldListValues           gtki.TreeView   `gtk-widget:"room-config-field-list-values-tree"`

	listModel gtki.ListStore
}

func newRoomConfigSummaryOccupantField(affiliation data.Affiliation, occupants map[data.Affiliation][]*muc.RoomOccupantItem) hasRoomConfigFormField {
	field := &roomConfigSummaryOccupantField{
		affiliation: affiliation,
		occupants:   occupants,
	}

	field.initBuilder()
	field.initDefaults()
	field.handleFieldValue()

	return field
}

func (f *roomConfigSummaryOccupantField) initBuilder() {
	builder := newBuilder("MUCRoomConfigSummaryField")
	panicOnDevError(builder.bindObjects(f))
	builder.ConnectSignals(map[string]interface{}{
		"on_show_list": f.onShowList,
	})
}

func (f *roomConfigSummaryOccupantField) initDefaults() {
	fieldLabel := i18n.Local("Owners")
	switch {
	case f.affiliation.IsAdmin():
		fieldLabel = i18n.Local("Administrators")
	case f.affiliation.IsBanned():
		fieldLabel = i18n.Local("Banned")
	}
	setLabelText(f.fieldLabel, fieldLabel)
}

func (f *roomConfigSummaryOccupantField) handleFieldValue() {
	f.listModel, _ = g.gtk.ListStoreNew(glibi.TYPE_STRING)
	f.fieldListValues.SetModel(f.listModel)

	setLabelText(f.fieldValueLabel, summaryListTotalText(len(f.occupants[f.affiliation])))
	f.fieldListValueButton.SetVisible(len(f.occupants[f.affiliation]) > 0)

	f.initListContent()
}

func (f *roomConfigSummaryOccupantField) initListContent() {
	f.listModel.Clear()

	for _, o := range f.occupants[f.affiliation] {
		iter := f.listModel.Append()
		f.listModel.SetValue(iter, 0, configOptionToFriendlyMessage(o.Jid.String(), o.Jid.String()))
	}
}

func (f *roomConfigSummaryOccupantField) onShowList() {
	summaryListHideOrShow(f.fieldListValues, f.fieldListValueButtonImage, f.fieldListContent)
}

func (f *roomConfigSummaryOccupantField) fieldWidget() gtki.Widget {
	return f.widget
}

func (f *roomConfigSummaryOccupantField) refreshContent() {
	f.handleFieldValue()
}

func (f *roomConfigSummaryOccupantField) updateFieldValue() {}

func (f *roomConfigSummaryOccupantField) isValid() bool {
	return true
}

func (f *roomConfigSummaryOccupantField) showValidationErrors() {}
