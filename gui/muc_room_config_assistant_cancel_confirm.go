package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigAssistantCancelView struct {
	assistant *roomConfigAssistant

	dialog gtki.Dialog `gtk-widget:"confirmation-dialog"`
	header gtki.Label  `gtk-widget:"confirmation-dialog-header"`
}

func (rc *roomConfigAssistant) newRoomConfigAssistantCancelView() *roomConfigAssistantCancelView {
	cv := &roomConfigAssistantCancelView{
		assistant: rc,
	}

	builder := newBuilder("MUCRoomConfirmationRoomDialog")
	panicOnDevError(builder.bindObjects(cv))

	builder.ConnectSignals(map[string]interface{}{
		"on_no_clicked":  cv.close,
		"on_yes_clicked": cv.onYesClicked,
	})

	cv.dialog.SetTransientFor(rc.assistant)

	return cv
}

// onYesClicked MUST be called from the UI thread
func (cv *roomConfigAssistantCancelView) onYesClicked() {
	cv.close()
	cv.assistant.cancelConfiguration()
}

// show MUST be called from the UI thread
func (cv *roomConfigAssistantCancelView) show() {
	cv.dialog.Show()
}

// close MUST be called from the UI thread
func (cv *roomConfigAssistantCancelView) close() {
	cv.dialog.Destroy()
}
