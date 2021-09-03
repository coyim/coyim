package gui

import "github.com/coyim/gotk3adapter/gtki"

type roomConfigPositionsWithApplyButton struct {
	*roomConfigFieldPositions
	applyButton gtki.Button
}

func newRoomConfigPositionsWithApplyButton(applyButton gtki.Button, options roomConfigPositionsOptions) hasRoomConfigFormField {
	rcpb := &roomConfigPositionsWithApplyButton{
		roomConfigFieldPositions: newRoomConfigFieldPositions(options),
		applyButton:              applyButton,
	}

	rcpb.connectUISignals()
	rcpb.initPositionsLists(options.parentWindow)

	return rcpb
}

func (rcpb *roomConfigPositionsWithApplyButton) connectUISignals() {
	rcpb.builder.ConnectSignals(map[string]interface{}{
		"on_jid_edited": rcpb.onOccupantJidEdited,
	})
}

func (rcpb *roomConfigPositionsWithApplyButton) initPositionsLists(parent gtki.Window) {
	// [ps] This should use the parent implementation to avoid repetition
	rcpb.positionsListController = newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:      rcpb.positionsAddButton,
		removeOccupantButton:   rcpb.positionsRemoveButton,
		removeOccupantLabel:    rcpb.positionsRemoveLabel,
		occupantsTreeView:      rcpb.positionsList,
		parentWindow:           parent,
		addOccupantDialogTitle: getFieldTextByAffiliation(rcpb.affiliation).dialogTitle,
		addOccupantDescription: getFieldTextByAffiliation(rcpb.affiliation).dialogDescription,
		onListUpdated:          rcpb.refreshContentLists,
	})

	rcpb.addItemsToListController()
}

// refreshContentLists MUST be called from the UI thread
func (rcpb *roomConfigPositionsWithApplyButton) refreshContentLists() {
	rcpb.roomConfigFieldPositions.refreshContentLists()
	rcpb.enableOrDisableApplyButton()
}

// onOccupantJidEdited MUST be called from the UI thread
func (rcpb *roomConfigPositionsWithApplyButton) onOccupantJidEdited(cell gtki.CellRendererText, path string, newValue string) {
	rcpb.roomConfigFieldPositions.onOccupantJidEdited(cell, path, newValue)
	rcpb.enableOrDisableApplyButton()
}

// enableOrDisableApplyButton MUST be called from the UI thread
func (rcpb *roomConfigPositionsWithApplyButton) enableOrDisableApplyButton() {
	rcpb.applyButton.SetSensitive(rcpb.hasListChanged())
}
