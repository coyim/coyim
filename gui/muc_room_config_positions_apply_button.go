package gui

import "github.com/coyim/gotk3adapter/gtki"

type roomConfigPositionsWithApplyButton struct {
	*roomConfigPositions
	applyButton gtki.Button
}

func newRoomConfigPositionsWithApplyButton(applyButton gtki.Button, options roomConfigPositionsOptions) hasRoomConfigFormField {
	rcpb := &roomConfigPositionsWithApplyButton{
		roomConfigPositions: newRoomConfigPositions(options),
		applyButton:         applyButton,
	}

	rcpb.loadUIDefinition()

	return rcpb
}

func (rcpb *roomConfigPositionsWithApplyButton) setUIBuilder(b *builder) {
	rcpb.builder = b
}

func (rcpb *roomConfigPositionsWithApplyButton) connectUISignals(b *builder) {
	b.ConnectSignals(map[string]interface{}{
		"on_jid_edited": rcpb.onOccupantJidEdited,
	})
}

func (rcpb *roomConfigPositionsWithApplyButton) loadUIDefinition() {
	buildUserInterface("MUCRoomConfigFieldPositions", rcpb.roomConfigPositions, rcpb.setUIBuilder, rcpb.connectUISignals)
}

// onOccupantJidEdited MUST be called from the UI thread
func (rcpb *roomConfigPositionsWithApplyButton) onOccupantJidEdited(crt gtki.CellRendererText, path string, newValue string) {
	rcpb.roomConfigPositions.onOccupantJidEdited(crt, path, newValue)
	rcpb.enableOrDisableApplyButton()
}

// enableOrDisableApplyButton MUST be called from the UI thread
func (rcpb *roomConfigPositionsWithApplyButton) enableOrDisableApplyButton() {
	rcpb.applyButton.SetSensitive(rcpb.hasListChanged())
}
