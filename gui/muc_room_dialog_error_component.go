package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/gotk3adapter/gtki"
)

type retryFunction func()

type dialogErrorComponent struct {
	builder *builder
	title   string
	header  string
	message string

	dialog       gtki.Dialog `gtk-widget:"room-error-dialog"`
	errorTitle   gtki.Label  `gtk-widget:"room-error-dialog-title"`
	errorMessage gtki.Label  `gtk-widget:"room-error-dialog-message"`

	retry retryFunction
}

func createDialogErrorComponent(title, header, message string, cb func()) *dialogErrorComponent {
	d := &dialogErrorComponent{
		title:   title,
		header:  header,
		message: message,
		retry:   cb,
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
	d.builder = newBuilder("MUCRoomDialogErrorComponent")

	panicOnDevError(d.builder.bindObjects(d))

	d.builder.ConnectSignals(map[string]interface{}{
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
