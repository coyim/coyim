package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigAssistantCancelView struct {
	builder *builder

	rc *roomConfigAssistant

	dialog  gtki.Dialog `gtk-widget:"confirmation-room-dialog"`
	cancel  gtki.Button `gtk-widget:"cancel-button"`
	confirm gtki.Button `gtk-widget:"destroy-room-button"`

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
		"on_no_confirmation_button":  cv.onNoConfirmation,
		"on_yes_confirmation_button": cv.onYesConfirmation,
	})
}

func (cv *roomConfigAssistantCancelView) onNoConfirmation() {
	cv.close()
}

func (cv *roomConfigAssistantCancelView) onYesConfirmation() {
	cv.close()
	cv.rc.disable()
	cv.rc.currentPage.onConfigurationCancel()
	cv.rc.roomConfigComponent.cancelConfiguration(cv.rc.onCancelSuccess, cv.rc.onCancelError)
}

// show MUST be called from the UI thread
func (cv *roomConfigAssistantCancelView) show() {
	cv.dialog.Show()
}

// close MUST be called from the UI thread
func (cv *roomConfigAssistantCancelView) close() {
	cv.dialog.Destroy()
}
