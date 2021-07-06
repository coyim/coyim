package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigSummaryOccupantField struct {
	affiliation            data.Affiliation
	occupantsByAffiliation func(data.Affiliation) []*muc.RoomOccupantItem

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

func newRoomConfigSummaryOccupantField(label string, affiliation data.Affiliation, occupantsByAffiliation func(data.Affiliation) []*muc.RoomOccupantItem) hasRoomConfigFormField {
	field := &roomConfigSummaryOccupantField{
		affiliation:            affiliation,
		occupantsByAffiliation: occupantsByAffiliation,
	}

	field.initBuilder()

	field.fieldLabel.SetText(label)
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

func (f *roomConfigSummaryOccupantField) handleFieldValue() {
	f.listModel, _ = g.gtk.ListStoreNew(glibi.TYPE_STRING)
	f.fieldListValues.SetModel(f.listModel)

	occupants := f.occupantsByAffiliation(f.affiliation)
	setLabelText(f.fieldValueLabel, summaryTotalPositionsText(len(occupants)))
	f.fieldListValueButton.SetVisible(len(occupants) > 0)

	f.initListContent()
}

func (f *roomConfigSummaryOccupantField) initListContent() {
	f.listModel.Clear()

	for _, o := range f.occupantsByAffiliation(f.affiliation) {
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
