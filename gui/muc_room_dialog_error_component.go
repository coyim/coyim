package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/gotk3adapter/gtki"
)

// dialogErrorComponent is shown after an error occurred.
type dialogErrorComponent struct {
	title   string
	header  string
	message string

	dialog       gtki.Dialog `gtk-widget:"room-error-dialog"`
	errorTitle   gtki.Label  `gtk-widget:"room-error-dialog-title"`
	errorMessage gtki.Label  `gtk-widget:"room-error-dialog-message"`
}

func createDialogErrorComponent(title, header, message string) *dialogErrorComponent {
	d := &dialogErrorComponent{
		title:   title,
		header:  header,
		message: message,
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
		"on_ok": d.onOkClicked,
	})
}

func (d *dialogErrorComponent) onOkClicked() {
	d.dialog.Destroy()
}

func (d *dialogErrorComponent) onRetry() {
	d.dialog.Destroy()
}

func (d *dialogErrorComponent) show() {
	d.dialog.Show()
}

func (d *dialogErrorComponent) updateMessageBasedOnError(err error) {
	switch err {
	case session.ErrRemoveOwnerAffiliation:
		d.errorMessage.SetText(i18n.Local("You can't change your own position because you are the only owner for this room. Every room must have at least one owner."))
	default:
		d.errorMessage.SetText(i18n.Local("An unknown error occurred during the process. Please try again later."))
	}
}
