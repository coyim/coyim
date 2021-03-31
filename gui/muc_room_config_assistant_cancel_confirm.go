package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigAssistantCancelView struct {
	builder *builder

	rc *roomConfigAssistant

	dialog gtki.Dialog `gtk-widget:"confirmation-dialog"`
	header gtki.Label  `gtk-widget:"confirmation-dialog-header"`

	log coylog.Logger
}

func newRoomConfigAssistantCancelView(rc *roomConfigAssistant) *roomConfigAssistantCancelView {
	d := &roomConfigAssistantCancelView{
		rc: rc,
	}

	d.initBuilder()

	return d
}

func (cv *roomConfigAssistantCancelView) initBuilder() {
	cv.builder = newBuilder("MUCRoomConfirmationRoomDialog")
	panicOnDevError(cv.builder.bindObjects(cv))

	cv.builder.ConnectSignals(map[string]interface{}{
		"on_no_clicked":  cv.onNoConfirmation,
		"on_yes_clicked": cv.onYesConfirmation,
	})

	mucStyles.setRoomCreationDialogErrorComponentHeaderStyle(cv.header)
}

func (cv *roomConfigAssistantCancelView) onNoConfirmation() {
	cv.close()
}

func (cv *roomConfigAssistantCancelView) onYesConfirmation() {
	cv.close()
	cv.rc.destroyAssistant()
	cv.rc.onCancel()
	cv.rc.roomConfigComponent.cancelConfiguration(cv.rc.onCancelError)
}

// show MUST be called from the UI thread
func (cv *roomConfigAssistantCancelView) show() {
	cv.dialog.Show()
}

// close MUST be called from the UI thread
func (cv *roomConfigAssistantCancelView) close() {
	cv.dialog.Destroy()
}
