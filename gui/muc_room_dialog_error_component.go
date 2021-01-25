package gui

import "github.com/coyim/gotk3adapter/gtki"

type retryFunction func()

type dialogErrorComponent struct {
	builder *builder
	title   string
	message string

	dialog       gtki.Dialog `gtk-widget:"room-error-dialog"`
	errorMessage gtki.Label  `gtk-widget:"title-error-message"`

	retry retryFunction
}

func createDialogErrorComponent(title, message string, cb func()) *dialogErrorComponent {
	d := &dialogErrorComponent{
		title:   title,
		message: message,
		retry:   cb,
	}

	d.initBuilderAndSignals()
	d.initDefaults()

	return d
}

func (d *dialogErrorComponent) initDefaults() {
	d.dialog.SetTitle(d.title)
	d.errorMessage.SetText(d.message)
}

func (d *dialogErrorComponent) initBuilderAndSignals() {
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
