package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

const (
	roomConfigSummaryOccupantJidIndex int = iota
)

type roomConfigSummaryOccupantField struct {
	retrieveOccupantList func() muc.RoomOccupantItemList

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

func newRoomConfigSummaryOccupantField(label string, retrieveOccupantList func() muc.RoomOccupantItemList) hasRoomConfigFormField {
	field := &roomConfigSummaryOccupantField{
		retrieveOccupantList: retrieveOccupantList,
	}

	field.initBuilder()
	field.initOccupantsModel()

	field.fieldLabel.SetText(label)

	return field
}

func (f *roomConfigSummaryOccupantField) initBuilder() {
	builder := newBuilder("MUCRoomConfigSummaryField")
	panicOnDevError(builder.bindObjects(f))
	builder.ConnectSignals(map[string]interface{}{
		"on_show_list": f.onShowList,
	})
}

func (f *roomConfigSummaryOccupantField) initOccupantsModel() {
	f.listModel, _ = g.gtk.ListStoreNew(
		// occupant jid
		glibi.TYPE_STRING,
	)

	f.fieldListValues.SetModel(f.listModel)
}

// handleFieldValue MUST be called from the UI thread
func (f *roomConfigSummaryOccupantField) handleFieldValue() {
	occupants := f.retrieveOccupantList()
	setLabelText(f.fieldValueLabel, summaryTotalPositionsText(len(occupants)))
	f.fieldListValueButton.SetVisible(len(occupants) > 0)

	f.printOccupantsView()
}

// refreshOccupantsView MUST be called from the UI thread
func (f *roomConfigSummaryOccupantField) printOccupantsView() {
	f.listModel.Clear()

	for _, o := range f.retrieveOccupantList() {
		iter := f.listModel.Append()
		f.listModel.SetValue(iter, roomConfigSummaryOccupantJidIndex, o.Jid.String())
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
