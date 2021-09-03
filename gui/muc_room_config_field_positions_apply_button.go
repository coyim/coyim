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

	rcpb.onListChanged.add(rcpb.enableOrDisableApplyButton)
	rcpb.onRefreshContentLists.add(rcpb.enableOrDisableApplyButton)

	return rcpb
}

// enableOrDisableApplyButton MUST be called from the UI thread
func (rcpb *roomConfigPositionsWithApplyButton) enableOrDisableApplyButton() {
	rcpb.applyButton.SetSensitive(rcpb.hasListChanged())
}
