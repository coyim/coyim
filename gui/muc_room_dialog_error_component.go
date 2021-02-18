package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/gotk3adapter/gtki"
)

// dialogErrorComponent is shown after an error occurred.
// It displays an  error message and gives an option to invoke
// the previously executed function through the `retry` callback.
type dialogErrorComponent struct {
	title   string
	header  string
	message string

	dialog       gtki.Dialog `gtk-widget:"room-error-dialog"`
	errorTitle   gtki.Label  `gtk-widget:"room-error-dialog-title"`
	errorMessage gtki.Label  `gtk-widget:"room-error-dialog-message"`

	retry func()
}

func createDialogErrorComponent(title, header, message string, retry func()) *dialogErrorComponent {
	d := &dialogErrorComponent{
		title:   title,
		header:  header,
		message: message,
		retry:   retry,
	}

	d.initBuilder()
	d.initDefaults()

	return d
}

func (d *dialogErrorComponent) initDefaults() {
	mucStyles.setLabelBoldStyle(d.errorTitle)

	d.dialog.SetTitle(d.title)
	d.errorTitle.SetText(d.header)
	d.errorMessage.SetText(d.message)
}

func (d *dialogErrorComponent) initBuilder() {
	builder := newBuilder("MUCRoomDialogErrorComponent")
	panicOnDevError(builder.bindObjects(d))

	builder.ConnectSignals(map[string]interface{}{
		"on_cancel": d.onCancel,
		"on_retry":  d.onRetry,
	})
}

func (d *dialogErrorComponent) onCancel() {
	d.dialog.Destroy()
}

func (d *dialogErrorComponent) onRetry() {
	go d.retry()
	d.dialog.Destroy()
}

func (d *dialogErrorComponent) show() {
	d.dialog.Show()
}

func (d *dialogErrorComponent) updateMessageForDestroyError(err error) {
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
	case session.ErrRemoveOwnerAffiliation:
		msg = i18n.Local("You can't change your own position because you are the only owner for this room. Every room must have at least one owner.")
	default:
		msg = i18n.Local("An unknown error occurred during the process. Please try again later.")
	}

	d.errorMessage.SetText(msg)
}
