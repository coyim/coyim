package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/gotk3adapter/gtki"
)

type destroyDialogError struct {
	dialog       gtki.Dialog `gtk-widget:"room-destory-error-dialog"`
	errorTitle   gtki.Label  `gtk-widget:"room-destroy-error-title"`
	errorMessage gtki.Label  `gtk-widget:"room-destroy-error-message"`

	retry func()
}

func createDestroyDialogError(retry func()) *destroyDialogError {
	d := &destroyDialogError{
		retry: retry,
	}

	d.initBuilder()
	d.initDefaults()

	return d
}

func (d *destroyDialogError) initDefaults() {
	mucStyles.setLabelBoldStyle(d.errorTitle)
	d.errorTitle.SetLabel(i18n.Local("The room couldn't be destroyed"))
}

func (d *destroyDialogError) initBuilder() {
	builder := newBuilder("MUCRoomDestroyDialogError")
	panicOnDevError(builder.bindObjects(d))

	builder.ConnectSignals(map[string]interface{}{
		"on_cancel": d.onCancel,
		"on_retry":  d.onRetry,
	})
}

func (d *destroyDialogError) onCancel() {
	d.dialog.Destroy()
}

func (d *destroyDialogError) onRetry() {
	go d.retry()
	d.dialog.Destroy()
}

func (d *destroyDialogError) show() {
	d.dialog.Show()
}

func (d *destroyDialogError) updateErrorMessage(err error) {
	msg := ""
	switch err {
	case session.ErrDestroyRoomInvalidIQResponse, session.ErrDestroyRoomNoResult:
		msg = i18n.Local("We were able to connect to the room service, " +
			"but we received an invalid response from it. Please try again later.")
	case session.ErrDestroyRoomForbidden:
		msg = i18n.Local("You don't have the permission to destroy this room. " +
			"Please contact one of the room owners.")
	case session.ErrDestroyRoomDoesntExist:
		msg = i18n.Local("We couldn't find the room.")
	default:
		msg = i18n.Local("An unknown error occurred during the process. Please try again later.")
	}

	d.errorMessage.SetText(msg)
}
