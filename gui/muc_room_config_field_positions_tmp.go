package gui

import "github.com/coyim/gotk3adapter/gtki"

// [ps] All this file MUST doesn't make sense and we should remove the
// use of the `roomConfigPositionsField` struct.

type roomConfigPositionsField struct {
	*roomConfigFieldPositions
}

func newRoomConfigPositionsField(options roomConfigPositionsOptions) hasRoomConfigFormField {
	rcpf := &roomConfigPositionsField{
		newRoomConfigFieldPositions(options),
	}

	rcpf.connectUISignals()
	rcpf.initPositionsLists(options.parentWindow)

	return rcpf
}

func (rcpf *roomConfigPositionsField) connectUISignals() {
	rcpf.builder.ConnectSignals(map[string]interface{}{
		"on_jid_edited": rcpf.onOccupantJidEdited,
	})
}

func (rcpf *roomConfigPositionsField) initPositionsLists(parent gtki.Window) {
	rcpf.positionsListController = newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:      rcpf.positionsAddButton,
		removeOccupantButton:   rcpf.positionsRemoveButton,
		removeOccupantLabel:    rcpf.positionsRemoveLabel,
		occupantsTreeView:      rcpf.positionsList,
		parentWindow:           parent,
		addOccupantDialogTitle: getFieldTextByAffiliation(rcpf.affiliation).dialogTitle,
		addOccupantDescription: getFieldTextByAffiliation(rcpf.affiliation).dialogDescription,
		onListUpdated:          rcpf.refreshContentLists,
	})

	rcpf.addItemsToListController()
}
